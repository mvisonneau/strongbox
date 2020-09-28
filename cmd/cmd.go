package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-input"
	cli "github.com/urfave/cli/v2"

	"github.com/mvisonneau/strongbox/rand"
)

// Execute all commands
func Execute(ctx *cli.Context) (int, error) {
	if err := configure(ctx); err != nil {
		return 1, err
	}

	log.Debugf("Function: %v", ctx.Command.FullName())

	switch ctx.Command.FullName() {
	case "transit use":
		if ctx.NArg() != 1 {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		s.Load()
		s.SetVaultTransitKey(ctx.Args().First())
	case "transit info":
		s.Load()
		v.GetTransitInfo()
	case "transit list":
		v.ListTransitKeys()
	case "transit create":
		if ctx.NArg() != 1 {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		s.Load()
		v.CreateTransitKey(ctx.Args().First())
		s.SetVaultTransitKey(ctx.Args().First())
	case "transit delete":
		if ctx.NArg() != 1 {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		v.DeleteTransitKey(ctx.Args().First())
	case "secret write":
		if ctx.NArg() != 1 ||
			ctx.String("key") == "" ||
			(ctx.String("value") == "" && !ctx.Bool("masked_value") && ctx.Int("random") == 0) ||
			(ctx.String("value") != "" && ctx.Bool("masked_value")) ||
			(ctx.String("value") != "" && ctx.Int("random") != 0) ||
			(ctx.Bool("masked_value") && ctx.Int("random") != 0) {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}

		s.Load()

		var secret string
		if ctx.Bool("masked_value") {
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
		} else if ctx.String("value") != "" {
			secret = v.Encrypt(ctx.String("value"))
		} else if ctx.Int("random") != 0 {
			secret = v.Encrypt(rand.String(ctx.Int("random")))
		} else {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}

		s.WriteSecretKey(ctx.Args().First(), ctx.String("key"), secret)
	case "secret read":
		if ctx.NArg() != 1 || ctx.String("key") == "" {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		s.Load()
		fmt.Println(v.Decrypt(s.ReadSecretKey(ctx.Args().First(), ctx.String("key"))))
	case "secret delete":
		if ctx.NArg() != 1 {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		s.Load()
		if ctx.String("key") == "" {
			s.DeleteSecret(ctx.Args().First())
		} else {
			s.DeleteSecretKey(ctx.Args().First(), ctx.String("key"))
		}
	case "secret list":
		s.Load()
		switch ctx.NArg() {
		case 1:
			s.ListSecrets(ctx.Args().First())
		default:
			s.ListSecrets("")
		}
	case "secret rotate-from":
		if ctx.NArg() != 1 {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		s.Load()
		s.RotateFromOldTransitKey(ctx.Args().First())
	case "get-secret-path":
		s.Load()
		fmt.Println(s.VaultSecretPath())
	case "set-secret-path":
		if ctx.NArg() != 1 {
			if err := cli.ShowSubcommandHelp(ctx); err != nil {
				return 1, err
			}
			return 1, nil
		}
		s.Load()
		s.SetVaultSecretPath(ctx.Args().First())
	case "init":
		s.Init()
	case "status":
		s.Load()
		s.Status()
		v.Status()
	case "plan":
		return run("plan")
	case "apply":
		return run("apply")
	default:
		log.Fatalf("Function %v not implemented yet", ctx.Command.FullName())
	}

	return 0, nil
}

func run(action string) (int, error) {
	s.Load()

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
		return 1, err
	}

	if d != nil {
		if keys, ok := d.Data["keys"]; ok {
			for _, k := range keys.([]interface{}) {
				if remote[k.(string)] == nil {
					remote[k.(string)] = make(map[string]string)
				}

				l, err := v.Client.Logical().Read(s.Vault.SecretPath + k.(string))
				if err != nil {
					return 1, err
				}

				for m, n := range l.Data {
					remote[k.(string)][m] = n.(string)
				}
			}
		}
	}

	eq := reflect.DeepEqual(local, remote)
	if eq {
		color.Green("Nothing to do! Local state and remote Vault config are in synctx.")
		return 0, nil
	}

	return reconcile(local, remote, action)
}

func reconcile(local map[string]map[string]string, remote map[string]map[string]string, action string) (int, error) {
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
		return 1, fmt.Errorf("No action specified")
	}

	return 0, nil
}
