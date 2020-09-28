package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mvisonneau/strongbox/rand"
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
	fmt.Println(v.Decipher(s.ReadSecretKey(ctx.String("secret"), ctx.String("key"))))

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
			return 1, err
		}

		secret = v.Cipher(value)
	} else if ctx.String("value") == "-" {
		read, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return 1, err
		}
		secret = v.Cipher(string(read))
	} else if ctx.String("value") != "" {
		secret = v.Cipher(ctx.String("value"))
	} else if ctx.Int("random") != 0 {
		secret = v.Cipher(rand.String(ctx.Int("random")))
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
