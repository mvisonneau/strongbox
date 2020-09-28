package main

import (
	"os"

	"github.com/mvisonneau/strongbox/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
