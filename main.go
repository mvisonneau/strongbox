package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mvisonneau/strongbox/app"
)

func main() {
	start := time.Now()
	app.Cli().Run(os.Args)
	log.Debugf("Executed in %s, exiting..", time.Since(start))
}
