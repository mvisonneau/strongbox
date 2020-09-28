package cmd

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/mvisonneau/go-helpers/logger"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var start time.Time
var s *State
var v *Vault

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	}

	err := lc.Configure()
	if err != nil {
		return err
	}

	v, err = getVaultClient(&VaultConfig{
		Address:  ctx.String("vault-addr"),
		Token:    ctx.String("vault-token"),
		RoleID:   ctx.String("vault-role-id"),
		SecretID: ctx.String("vault-secret-id"),
	})

	if err != nil {
		return err
	}

	s = getStateClient(&StateConfig{
		Path: ctx.String("state"),
	})

	return nil
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

// ExecWrapper mlocks the process memory (if supported) before our `run` functions,
// and gracefully logs and exits afterwards.
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := mlock.LockMemory(); err != nil {
			return exit(1, fmt.Errorf("error locking strongbox memory: %w", err))
		}
		return exit(f(ctx))
	}
}
