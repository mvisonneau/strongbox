package app

import (
  "os"
  "os/exec"
  "testing"

	log "github.com/sirupsen/logrus"
)

func TestConfigureLoggingFatalText(t *testing.T) {
  cfg.Log.Level = "fatal"
  cfg.Log.Format = "text"
  configureLogging()

  if log.GetLevel() != log.FatalLevel {
    t.Fatalf("Expected log.Level to be 'fatal' but got %s", log.GetLevel())
  }
}

func TestConfigureLoggingDefault(t *testing.T) {
  cfg.Log.Level = "fatal"
  cfg.Log.Format = "default"

  if os.Getenv("BE_CRASHER") == "1" {
      configureLogging()
      return
  }
  cmd := exec.Command(os.Args[0], "-test.run=TestConfigureLoggingDefault")
  cmd.Env = append(os.Environ(), "BE_CRASHER=1")
  err := cmd.Run()
  if e, ok := err.(*exec.ExitError); ok && !e.Success() {
      return
  }
  t.Fatalf("Process ran with err %v, wanted exit status 1", err)
}
