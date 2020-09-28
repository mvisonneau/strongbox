package cmd

import (
	"fmt"
	"strconv"

	cli "github.com/urfave/cli/v2"
)

// KVGetPath ..
func KVGetPath(_ *cli.Context) (int, error) {
	s.Load()
	fmt.Println(s.VaultKVPath())

	return 0, nil
}

// KVSetPath ..
func KVSetPath(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()
	s.SetVaultKVPath(ctx.Args().First())

	return 0, nil
}

// KVGetVersion ..
func KVGetVersion(_ *cli.Context) (int, error) {
	s.Load()
	fmt.Println(s.VaultKVVersion())

	return 0, nil
}

// KVSetVersion ..
func KVSetVersion(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()

	version, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return 1, err
	}

	if version != 1 && version != 2 {
		return 1, fmt.Errorf("KV version must be either 1 or 2, got %d", version)
	}

	s.SetVaultKVVersion(version)

	return 0, nil
}
