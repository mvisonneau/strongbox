package app

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"

	"github.com/mvisonneau/strongbox/config"
	"github.com/mvisonneau/strongbox/rand"
)

var cfg config.Config
var s State
var start time.Time
var v Vault

func execute(c *cli.Context) error {
	start = time.Now()
	err := configureLogging(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	log.Debugf("Function: %v", c.Command.FullName())

	switch c.Command.FullName() {
	case "transit use":
		if c.NArg() != 1 {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		s.Load()
		s.SetVaultTransitKey(c.Args().First())
	case "transit info":
		s.Load()
		v.ConfigureClient()
		v.GetTransitInfo()
	case "transit list":
		v.ConfigureClient()
		v.ListTransitKeys()
	case "transit create":
		if c.NArg() != 1 {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		s.Load()
		v.ConfigureClient()
		v.CreateTransitKey(c.Args().First())
		s.SetVaultTransitKey(c.Args().First())
	case "transit delete":
		if c.NArg() != 1 {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		v.ConfigureClient()
		v.DeleteTransitKey(c.Args().First())
	case "secret write":
		if c.NArg() != 1 ||
			c.String("key") == "" ||
			(c.String("value") == "" && !c.Bool("masked_value") && c.Int("random") == 0) ||
			(c.String("value") != "" && c.Bool("masked_value")) ||
			(c.String("value") != "" && c.Int("random") != 0) ||
			(c.Bool("masked_value") && c.Int("random") != 0) {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}

		s.Load()
		v.ConfigureClient()

		var secret string
		if c.Bool("masked_value") {
			ui := &input.UI{
				Writer: os.Stdout,
				Reader: os.Stdin,
			}

			value, err := ui.Ask("Sensitive", &input.Options{
				Required:    true,
				Mask:        true,
				MaskDefault: true,
			})

			if err != nil {
				log.Fatal(err)
			}

			secret = v.Encrypt(value)
		} else if c.String("value") != "" {
			secret = v.Encrypt(c.String("value"))
		} else if c.Int("random") != 0 {
			secret = v.Encrypt(rand.String(c.Int("random")))
		} else {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}

		s.WriteSecretKey(c.Args().First(), c.String("key"), secret)
	case "secret read":
		if c.NArg() != 1 || c.String("key") == "" {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		s.Load()
		v.ConfigureClient()
		fmt.Println(v.Decrypt(s.ReadSecretKey(c.Args().First(), c.String("key"))))
	case "secret delete":
		if c.NArg() != 1 {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		s.Load()
		if c.String("key") == "" {
			s.DeleteSecret(c.Args().First())
		} else {
			s.DeleteSecretKey(c.Args().First(), c.String("key"))
		}
	case "secret list":
		s.Load()
		switch c.NArg() {
		case 1:
			s.ListSecrets(c.Args().First())
		default:
			s.ListSecrets("")
		}
	case "secret rotate-from":
		if c.NArg() != 1 {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		s.Load()
		v.ConfigureClient()
		s.RotateFromOldTransitKey(c.Args().First())
	case "get-secret-path":
		s.Load()
		fmt.Println(s.VaultSecretPath())
	case "set-secret-path":
		if c.NArg() != 1 {
			cli.ShowSubcommandHelp(c)
			return cli.NewExitError("", 1)
		}
		s.Load()
		s.SetVaultSecretPath(c.Args().First())
	case "init":
		s.Init()
	case "status":
		s.Load()
		v.ConfigureClient()
		s.Status()
		v.Status()
	case "plan":
		run("plan")
	case "apply":
		run("apply")
	default:
		log.Fatalf("Function %v not implemented yet", c.Command.FullName())
	}

	return nil
}

func run(action string) error {
	s.Load()
	v.ConfigureClient()

	// Fetch local values
	local := make(map[string]map[string]string)
	for k, l := range s.Secrets {
		if local[k] == nil {
			local[k] = make(map[string]string)
		}
		for m, n := range l {
			local[k][m] = v.Decrypt(n)
		}
	}

	// Fetch remote values
	remote := make(map[string]map[string]string)
	d, err := v.Client.Logical().List(s.Vault.SecretPath)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	if d != nil {
		if keys, ok := d.Data["keys"]; ok {
			for _, k := range keys.([]interface{}) {
				if remote[k.(string)] == nil {
					remote[k.(string)] = make(map[string]string)
				}

				l, err := v.Client.Logical().Read(s.Vault.SecretPath + k.(string))
				if err != nil {
					return exit(cli.NewExitError(err.Error(), 1))
				}

				for m, n := range l.Data {
					remote[k.(string)][m] = n.(string)
				}
			}
		}
	}

	eq := reflect.DeepEqual(local, remote)
	if eq {
		color.Green("Nothing to do! Local state and remote Vault config are in sync.")
	} else {
		reconcile(local, remote, action)
	}

	return nil
}

func reconcile(local map[string]map[string]string, remote map[string]map[string]string, action string) error {
	var addSecret, deleteSecret []string
	addSecretKey := make(map[string][]string)
	deleteSecretKey := make(map[string][]string)
	addSecretKeyCount := 0
	deleteSecretKeyCount := 0

	for kl, vl := range local {
		foundSecret := false
		for kr, vr := range remote {
			if kl == kr {
				foundSecret = true
				for klk := range vl {
					foundSecretKey := false
					for krk := range vr {
						if klk == krk {
							foundSecretKey = true
							break
						}
					}

					if !foundSecretKey {
						if addSecretKey[kl] == nil {
							addSecretKey[kl] = make([]string, 0)
						}
						addSecretKey[kl] = append(addSecretKey[kl], klk)
						addSecretKeyCount++
					}
				}
				break
			}
		}

		if !foundSecret {
			addSecret = append(addSecret, kl)
			addSecretKey[kl] = make([]string, 0)
			for klk := range vl {
				addSecretKey[kl] = append(addSecretKey[kl], klk)
				addSecretKeyCount++
			}
		}
	}

	for kr, vr := range remote {
		foundSecret := false
		for kl, vl := range local {
			if kr == kl {
				foundSecret = true
				for krk := range vr {
					foundSecretKey := false
					for klk := range vl {
						if krk == klk {
							foundSecretKey = true
							break
						}
					}

					if !foundSecretKey {
						fmt.Println(kr, krk)
						if deleteSecretKey[kr] == nil {
							deleteSecretKey[kr] = make([]string, 0)
						}
						deleteSecretKey[kr] = append(deleteSecretKey[kr], krk)
						deleteSecretKeyCount++
					}
				}
				break
			}
		}

		if !foundSecret {
			deleteSecret = append(deleteSecret, kr)
			deleteSecretKey[kr] = make([]string, 0)
			for krk := range vr {
				deleteSecretKey[kr] = append(deleteSecretKey[kr], krk)
				deleteSecretKeyCount++
			}
		}
	}

	switch action {
	case "plan":
		if (len(addSecret) > 0) || (addSecretKeyCount > 0) {
			color.Green("Add/Update: %v secret(s) and %v key(s)", len(addSecret), addSecretKeyCount)
			for k, l := range addSecretKey {
				for _, m := range l {
					color.Green("=> %v:%v", k, m)
				}
			}
		}
		if (len(deleteSecret) > 0) || (deleteSecretKeyCount > 0) {
			color.Red("Remove: %v secret(s) and %v key(s)", len(deleteSecret), deleteSecretKeyCount)
			for k, l := range deleteSecretKey {
				for _, m := range l {
					color.Red("=> %v:%v", k, m)
				}
			}
		}
	case "apply":
		for _, k := range addSecret {
			payload := make(map[string]interface{})
			for m, n := range local[k] {
				payload[m] = n
			}
			v.WriteSecret(k, payload)
		}
		for k := range addSecretKey {
			payload := make(map[string]interface{})
			for m, n := range local[k] {
				payload[m] = n
			}
			v.WriteSecret(k, payload)
		}
		for _, k := range deleteSecret {
			v.DeleteSecret(k)
		}
		for k := range deleteSecretKey {
			payload := make(map[string]interface{})
			for m, n := range local[k] {
				payload[m] = n
			}
			v.WriteSecret(k, payload)
		}
	default:
		return exit(cli.NewExitError("No action specified", 1))
	}

	return nil
}

func exit(err error) error {
	log.Debugf("Executed in %s, exiting..", time.Since(start))
	return err
}
