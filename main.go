package main

import (
	"os"

	"github.com/mvisonneau/strongbox/app"
)

var version = "<devel>"

func main() {
	app.Cli(version).Run(os.Args)
}
