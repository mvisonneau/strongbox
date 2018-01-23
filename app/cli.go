package app

import (
	"fmt"

	"github.com/urfave/cli"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}
}

// Cli : Generates cli configuration for the application
func Cli(version string) (c *cli.App) {
	c = cli.NewApp()
	c.Name = "strongbox"
	c.Version = version
	c.Usage = "Securely store secrets at rest with Hashicorp Vault"
	c.EnableBashCompletion = true

	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "state,s",
			EnvVar:      "STRONGBOX_STATE",
			Usage:       "load state from `FILE`",
			Value:       ".strongbox_state.yml",
			Destination: &cfg.State.Location,
		},
		cli.StringFlag{
			Name:        "vault-addr",
			EnvVar:      "VAULT_ADDR",
			Usage:       "vault endpoint",
			Destination: &cfg.Vault.Address,
		},
		cli.StringFlag{
			Name:        "vault-token",
			EnvVar:      "VAULT_TOKEN",
			Usage:       "vault token",
			Destination: &cfg.Vault.Token,
		},
		cli.StringFlag{
			Name:        "vault-role-id",
			EnvVar:      "VAULT_ROLE_ID",
			Usage:       "vault role id",
			Destination: &cfg.Vault.RoleID,
		},
		cli.StringFlag{
			Name:        "vault-secret-id",
			EnvVar:      "VAULT_SECRET_ID",
			Usage:       "vault secret id",
			Destination: &cfg.Vault.SecretID,
		},
		cli.StringFlag{
			Name:        "log-level",
			EnvVar:      "STRONGBOX_LOG_LEVEL",
			Usage:       "log level (debug,info,warn,fatal,panic)",
			Value:       "info",
			Destination: &cfg.Log.Level,
		},
		cli.StringFlag{
			Name:        "log-format",
			EnvVar:      "STRONGBOX_LOG_FORMAT",
			Usage:       "log format (json,text)",
			Value:       "text",
			Destination: &cfg.Log.Format,
		},
	}

	c.Commands = []cli.Command{
		{
			Name:  "transit",
			Usage: "perform actions on transit key/backend",
			Subcommands: []cli.Command{
				{
					Name:      "use",
					Usage:     "configure a transit key to use",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    execute,
				},
				{
					Name:      "info",
					Usage:     "get information about the currently used transit key",
					ArgsUsage: " ",
					Action:    execute,
				},
				{
					Name:      "list",
					Usage:     "list available transit keys",
					ArgsUsage: " ",
					Action:    execute,
				},
				{
					Name:      "create",
					Usage:     "create and use a transit key",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    execute,
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
					Action: execute,
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
					Action: execute,
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
					Action: execute,
				},
				{
					Name:      "list",
					Usage:     "list all managed secrets",
					ArgsUsage: " ",
					Action:    execute,
				},
				{
					Name:      "rotate-from",
					Usage:     "rotate local secrets encryption from an old transit key",
					ArgsUsage: "<old_vault_transit_key>",
					Action:    execute,
				},
			},
		},
		{
			Name:      "get-secret-path",
			Usage:     "display the currently used vault secret path in the statefile",
			ArgsUsage: " ",
			Action:    execute,
		},
		{
			Name:      "set-secret-path",
			Usage:     "update the vault secret path in the statefile",
			ArgsUsage: "<secret_path>",
			Action:    execute,
		},
		{
			Name:      "init",
			Usage:     "Create a empty state file at configured location",
			ArgsUsage: " ",
			Action:    execute,
		},
		{
			Name:      "status",
			Usage:     "display current status",
			ArgsUsage: " ",
			Action:    execute,
		},
		{
			Name:      "plan",
			Usage:     "compare local version with vault cluster",
			ArgsUsage: " ",
			Action:    execute,
		},
		{
			Name:      "apply",
			Usage:     "synchronize vault managed secrets",
			ArgsUsage: " ",
			Action:    execute,
		},
	}

	return
}
