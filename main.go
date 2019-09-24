package main

import (
	"os"

	"github.com/mvisonneau/strongbox/cli"
)

var version = ""

func main() {
	cli.Init(&version).Run(os.Args)
}
