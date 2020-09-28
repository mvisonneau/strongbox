package cli

import (
	"fmt"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/mvisonneau/strongbox/cmd"
)

// Run handles the instanciation of the CLI application
func Run(version string, args []string) {
	err := NewApp(version, time.Now()).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "strongbox"
	app.Version = version
	app.Usage = "safely manage Hashicorp Vault secrets at rest"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.StringFlag{
			Name:    "state",
			Aliases: []string{"s"},
			EnvVars: []string{"STRONGBOX_STATE"},
			Usage:   "load state from `FILE`",
			Value:   ".strongbox_state.yml",
		},
		&cli.StringFlag{
			Name:    "vault-address",
			EnvVars: []string{"VAULT_ADDR"},
			Usage:   "vault endpoint",
		},
		&cli.StringFlag{
			Name:    "vault-token",
			EnvVars: []string{"VAULT_TOKEN"},
			Usage:   "vault token",
		},
		&cli.StringFlag{
			Name:    "vault-role-id",
			EnvVars: []string{"VAULT_ROLE_ID"},
			Usage:   "vault role id",
		},
		&cli.StringFlag{
			Name:    "vault-secret-id",
			EnvVars: []string{"VAULT_SECRET_ID"},
			Usage:   "vault secret id",
		},
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"STRONGBOX_LOG_LEVEL"},
			Usage:   "log level (debug,info,warn,fatal,panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "log-format",
			EnvVars: []string{"STRONGBOX_LOG_FORMAT"},
			Usage:   "log format (json,text)",
			Value:   "text",
		},
	}

	app.Commands = cli.CommandsByName{
		{
			Name:  "transit",
			Usage: "perform actions on transit key/backend",
			Subcommands: cli.CommandsByName{
				{
					Name:      "use",
					Usage:     "configure a transit key to use",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "info",
					Usage:     "get information about the currently used transit key",
					ArgsUsage: " ",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "list",
					Usage:     "list available transit keys",
					ArgsUsage: " ",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "create",
					Usage:     "create and use a transit key",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "delete",
					Usage:     "delete an existing transit key from Vault",
					ArgsUsage: "<vault_transit_key_name>",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
			},
		},
		{
			Name:  "secret",
			Usage: "perform actions on secrets (locally)",
			Subcommands: cli.CommandsByName{
				{
					Name:      "write",
					Usage:     "write a secret",
					ArgsUsage: "<secret> -k <key> [-v <value> or -r <string_length> or -V]",
					Flags: cli.FlagsByName{
						&cli.StringFlag{
							Name:  "key,k",
							Usage: "key name",
						},
						&cli.StringFlag{
							Name:  "value,v",
							Usage: "sensitive value of the key to encrypt",
						},
						&cli.BoolFlag{
							Name:  "masked_value,V",
							Usage: "sensitive value of the key to encrypt (stdin)",
						},
						&cli.IntFlag{
							Name:  "random,r",
							Usage: "automatically generates a string of this length",
						},
					},
					Action: cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "read",
					Usage:     "read secret value",
					ArgsUsage: "<secret> -k <key>",
					Flags: cli.FlagsByName{
						&cli.StringFlag{
							Name:  "key,k",
							Usage: "key name",
						},
					},
					Action: cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "delete",
					Usage:     "delete secret",
					ArgsUsage: "<secret>",
					Flags: cli.FlagsByName{
						&cli.StringFlag{
							Name:  "key,k",
							Usage: "key name",
						},
					},
					Action: cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "list",
					Usage:     "list all managed secrets",
					ArgsUsage: " ",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
				{
					Name:      "rotate-from",
					Usage:     "rotate local secrets encryption from an old transit key",
					ArgsUsage: "<old_vault_transit_key>",
					Action:    cmd.ExecWrapper(cmd.Execute),
				},
			},
		},
		{
			Name:      "get-secret-path",
			Usage:     "display the currently used vault secret path in the statefile",
			ArgsUsage: " ",
			Action:    cmd.ExecWrapper(cmd.Execute),
		},
		{
			Name:      "set-secret-path",
			Usage:     "update the vault secret path in the statefile",
			ArgsUsage: "<secret_path>",
			Action:    cmd.ExecWrapper(cmd.Execute),
		},
		{
			Name:      "init",
			Usage:     "Create a empty state file at configured location",
			ArgsUsage: " ",
			Action:    cmd.ExecWrapper(cmd.Execute),
		},
		{
			Name:      "status",
			Usage:     "display current status",
			ArgsUsage: " ",
			Action:    cmd.ExecWrapper(cmd.Execute),
		},
		{
			Name:      "plan",
			Usage:     "compare local version with vault cluster",
			ArgsUsage: " ",
			Action:    cmd.ExecWrapper(cmd.Execute),
		},
		{
			Name:      "apply",
			Usage:     "synchronize vault managed secrets",
			ArgsUsage: " ",
			Action:    cmd.ExecWrapper(cmd.Execute),
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
