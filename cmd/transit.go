package cmd

import cli "github.com/urfave/cli/v2"

// TransitUse ..
func TransitUse(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()
	s.SetVaultTransitKey(ctx.Args().First())

	return 0, nil
}

// TransitInfo ..
func TransitInfo(ctx *cli.Context) (int, error) {
	s.Load()
	v.GetTransitInfo()

	return 0, nil
}

// TransitList ..
func TransitList(ctx *cli.Context) (int, error) {
	v.ListTransitKeys()

	return 0, nil
}

// TransitCreate ..
func TransitCreate(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	s.Load()
	v.CreateTransitKey(ctx.Args().First())
	s.SetVaultTransitKey(ctx.Args().First())

	return 0, nil
}

// TransitDelete ..
func TransitDelete(ctx *cli.Context) (int, error) {
	if ctx.NArg() != 1 {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, nil
	}
	v.DeleteTransitKey(ctx.Args().First())

	return 0, nil
}
