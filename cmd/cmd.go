package cmd

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
	cli "github.com/urfave/cli/v2"
)

// Init ..
func Init(_ *cli.Context) (int, error) {
	s.Init()
	return 0, nil
}

// Status ..
func Status(_ *cli.Context) (int, error) {
	s.Load()
	s.Status()
	v.Status()
	return 0, nil
}

// Plan ..
func Plan(_ *cli.Context) (int, error) {
	return run("plan")
}

// Apply ..
func Apply(_ *cli.Context) (int, error) {
	return run("apply")
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
			local[k][m] = v.Decipher(n)
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
