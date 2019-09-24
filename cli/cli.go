package cli

import (
	"time"

	"github.com/mvisonneau/strongbox/cmd"
	"github.com/urfave/cli"
)

// Init : Generates CLI configuration for the application
func Init(version *string) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "strongbox"
	app.Version = *version
	app.Usage = "safely manage Hashicorp Vault secrets at rest"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "state,s",
			EnvVar: "STRONGBOX_STATE",
			Usage:  "load state from `FILE`",
			Value:  ".strongbox_state.yml",
		},
		cli.StringFlag{
			Name:   "vault-addr",
			EnvVar: "VAULT_ADDR",
			Usage:  "vault endpoint",
		},
		cli.StringFlag{
			Name:   "vault-token",
			EnvVar: "VAULT_TOKEN",
			Usage:  "vault token",
		},
		cli.StringFlag{
			Name:   "vault-role-id",
			EnvVar: "VAULT_ROLE_ID",
			Usage:  "vault role id",
		},
		cli.StringFlag{
			Name:   "vault-secret-id",
			EnvVar: "VAULT_SECRET_ID",
			Usage:  "vault secret id",
		},
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "STRONGBOX_LOG_LEVEL",
			Usage:  "log level (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "STRONGBOX_LOG_FORMAT",
			Usage:  "log format (json,text)",
			Value:  "text",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "transit",
			Usage: "perform actions on transit key/backend",
			Subcommands: []cli.Command{
				{
					Name:      "use",
					Usage:     "configure a transit key to use",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    cmd.Execute,
				},
				{
					Name:      "info",
					Usage:     "get information about the currently used transit key",
					ArgsUsage: " ",
					Action:    cmd.Execute,
				},
				{
					Name:      "list",
					Usage:     "list available transit keys",
					ArgsUsage: " ",
					Action:    cmd.Execute,
				},
				{
					Name:      "create",
					Usage:     "create and use a transit key",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    cmd.Execute,
				},
				{
					Name:      "delete",
					Usage:     "delete an existing transit key from Vault",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    cmd.Execute,
				},
			},
		},
		{
			Name:  "secret",
			Usage: "perform actions on secrets (locally)",
			Subcommands: []cli.Command{
				{
					Name:      "write",
					Usage:     "write a secret",
					ArgsUsage: "<secret> -k <key> [-v <value> or -r <string_length> or -V]",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "key,k",
							Usage: "key name",
						},
						cli.StringFlag{
							Name:  "value,v",
							Usage: "sensitive value of the key to encrypt",
						},
						cli.BoolFlag{
							Name:  "masked_value,V",
							Usage: "sensitive value of the key to encrypt (stdin)",
						},
						cli.IntFlag{
							Name:  "random,r",
							Usage: "automatically generates a string of this length",
						},
					},
					Action: cmd.Execute,
				},
				{
					Name:      "read",
					Usage:     "read secret value",
					ArgsUsage: "<secret> -k <key>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "key,k",
							Usage: "key name",
						},
					},
					Action: cmd.Execute,
				},
				{
					Name:      "delete",
					Usage:     "delete secret",
					ArgsUsage: "<secret>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "key,k",
							Usage: "key name",
						},
					},
					Action: cmd.Execute,
				},
				{
					Name:      "list",
					Usage:     "list all managed secrets",
					ArgsUsage: " ",
					Action:    cmd.Execute,
				},
				{
					Name:      "rotate-from",
					Usage:     "rotate local secrets encryption from an old transit key",
					ArgsUsage: "<old_vault_transit_key>",
					Action:    cmd.Execute,
				},
			},
		},
		{
			Name:      "get-secret-path",
			Usage:     "display the currently used vault secret path in the statefile",
			ArgsUsage: " ",
			Action:    cmd.Execute,
		},
		{
			Name:      "set-secret-path",
			Usage:     "update the vault secret path in the statefile",
			ArgsUsage: "<secret_path>",
			Action:    cmd.Execute,
		},
		{
			Name:      "init",
			Usage:     "Create a empty state file at configured location",
			ArgsUsage: " ",
			Action:    cmd.Execute,
		},
		{
			Name:      "status",
			Usage:     "display current status",
			ArgsUsage: " ",
			Action:    cmd.Execute,
		},
		{
			Name:      "plan",
			Usage:     "compare local version with vault cluster",
			ArgsUsage: " ",
			Action:    cmd.Execute,
		},
		{
			Name:      "apply",
			Usage:     "synchronize vault managed secrets",
			ArgsUsage: " ",
			Action:    cmd.Execute,
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": time.Now(),
	}

	return
}
