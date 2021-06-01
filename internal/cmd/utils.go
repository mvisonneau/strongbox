package cmd

import (
	"time"

	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/mvisonneau/go-helpers/logger"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	start time.Time
	s     *State
	v     *Vault
)

func configure(ctx *cli.Context) (err error) {
	start = ctx.App.Metadata["startTime"].(time.Time)

	if err = logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	}); err != nil {
		return
	}

	if v, err = getVaultClient(&VaultConfig{
		Address:  ctx.String("vault-addr"),
		Token:    ctx.String("vault-token"),
		RoleID:   ctx.String("vault-role-id"),
		SecretID: ctx.String("vault-secret-id"),
	}); err != nil {
		return
	}

	s = getStateClient(&StateConfig{
		Path: ctx.String("state"),
	})

	return
}

func exit(exitCode int, err error) cli.ExitCoder {
	defer log.WithFields(
		log.Fields{
			"execution-time": time.Since(start),
		},
	).Debug("exited..")

	if err != nil {
		log.Error(err.Error())
	}

	return cli.NewExitError("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := mlock.LockMemory(); err != nil {
			log.WithError(err).Warn("s5 requires the IPC_LOCK capability in order to secure its memory")
		}
		return exit(f(ctx))
	}
}
