package app

import (
  "testing"

	log "github.com/sirupsen/logrus"
)

func TestConfigureLogging(t *testing.T) {
  cfg.Log.Level = "debug"
  cfg.Log.Format = "text"
  configureLogging()

  if log.GetLevel() != log.DebugLevel {
    t.Fatalf("Expected log.Level to be debug but got %s", log.GetLevel())
  }
}
