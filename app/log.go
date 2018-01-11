package app

import (
  "fmt"
  "os"

	log "github.com/sirupsen/logrus"
)

func configureLogging() error {
  level, _ := log.ParseLevel(cfg.Log.Level)
  log.SetLevel(level)

  formatter := &log.TextFormatter{
    FullTimestamp: true,
  }
  log.SetFormatter(formatter)

	switch cfg.Log.Format {
    case "text":
  	case "json":
  		log.SetFormatter(&log.JSONFormatter{})
    default:
      fmt.Sprintln("Invalid log format '%v'", cfg.Log.Format)
      os.Exit(1)
  }

  log.SetOutput(os.Stdout)

  return nil
}
