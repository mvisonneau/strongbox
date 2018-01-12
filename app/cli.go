package app

import (
	"github.com/urfave/cli"
)

// Cli : Generates cli configuration for the application
func Cli() (c *cli.App) {
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
			Value:       "~/.strongbox_state.yml",
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
					Name:   "use",
					Usage:  "configure a transit key to use",
					Action: execute,
				},
				{
					Name:   "info",
					Usage:  "get information about the currently used transit key",
					Action: execute,
				},
				{
					Name:   "list",
					Usage:  "list available transit keys",
					Action: execute,
				},
				{
					Name:   "create",
					Usage:  "create and use a transit key",
					Action: execute,
				},
			},
		},
		{
			Name:  "secret",
			Usage: "perform actions on secrets (locally)",
			Subcommands: []cli.Command{
				{
					Name:   "write",
					Usage:  "write a secret",
					Action: execute,
				},
				{
					Name:   "read",
					Usage:  "read secret value",
					Action: execute,
				},
				{
					Name:   "delete",
					Usage:  "delete secret",
					Action: execute,
				},
				{
					Name:   "list",
					Usage:  "list all managed secrets",
					Action: execute,
				},
			},
		},
		{
			Name:   "get-secret-path",
			Usage:  "display the currently used vault secret path in the statefile",
			Action: execute,
		},
		{
			Name:   "set-secret-path",
			Usage:  "update the vault secret path in the statefile",
			Action: execute,
		},
		{
			Name:   "init",
			Usage:  "Create a empty state file at configured location",
			Action: execute,
		},
		{
			Name:   "status",
			Usage:  "display current status",
			Action: execute,
		},
		{
			Name:   "plan",
			Usage:  "compare local version with vault cluster",
			Action: execute,
		},
		{
			Name:   "apply",
			Usage:  "synchronize vault managed secrets",
			Action: execute,
		},
	}

	return
}
