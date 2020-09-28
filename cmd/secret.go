package cmd

import (
	"fmt"
	"os"

	"github.com/mvisonneau/strongbox/rand"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-input"
	cli "github.com/urfave/cli/v2"
)

// SecretRead ..
func SecretRead(ctx *cli.Context) (int, error) {
	if ctx.String("secret") == "" || ctx.String("key") == "" {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()
	fmt.Println(v.Decrypt(s.ReadSecretKey(ctx.String("secret"), ctx.String("key"))))

	return 0, nil
}

// SecretWrite ..
func SecretWrite(ctx *cli.Context) (int, error) {
	if ctx.String("secret") == "" ||
		ctx.String("key") == "" ||
		(ctx.String("value") == "" && !ctx.Bool("masked_value") && ctx.Int("random") == 0) ||
		(ctx.String("value") != "" && ctx.Bool("masked_value")) ||
		(ctx.String("value") != "" && ctx.Int("random") != 0) ||
		(ctx.Bool("masked_value") && ctx.Int("random") != 0) {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, fmt.Errorf("invalid arguments provided")
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

	s.WriteSecretKey(ctx.String("secret"), ctx.String("key"), secret)

	return 0, nil
}

// SecretList ..
func SecretList(ctx *cli.Context) (int, error) {
	s.Load()
	switch ctx.NArg() {
	case 1:
		s.ListSecrets(ctx.Args().First())
	default:
		s.ListSecrets("")
	}

	return 0, nil
}

// SecretDelete ..
func SecretDelete(ctx *cli.Context) (int, error) {
	if ctx.String("secret") == "" {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()

	if ctx.String("key") == "" {
		s.DeleteSecret(ctx.String("secret"))
	} else {
		s.DeleteSecretKey(ctx.String("secret"), ctx.String("key"))
	}

	return 0, nil
}

// SecretRotateFrom ..
func SecretRotateFrom(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()
	s.RotateFromOldTransitKey(ctx.Args().First())

	return 0, nil
}

// SecretGetPath ..
func SecretGetPath(ctx *cli.Context) (int, error) {
	s.Load()
	fmt.Println(s.VaultSecretPath())

	return 0, nil
}

// SecretSetPath ..
func SecretSetPath(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()
	s.SetVaultSecretPath(ctx.Args().First())

	return 0, nil
}
