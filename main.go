package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mvisonneau/strongbox/app"
)

var version = "<devel>"

func main() {
	start := time.Now()
	app.Cli(version).Run(os.Args)
	log.Debugf("Executed in %s, exiting..", time.Since(start))
}
