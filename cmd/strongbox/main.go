package main

import (
	"os"

	"github.com/mvisonneau/strongbox/internal/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
