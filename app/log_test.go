package app

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestConfigureLoggingFatalText(t *testing.T) {
	configureLogging("fatal", "text")

	if log.GetLevel() != log.FatalLevel {
		t.Fatalf("Expected log.Level to be 'fatal' but got %s", log.GetLevel())
	}
}

func TestConfigureLoggingDefault(t *testing.T) {
	err := configureLogging("fatal", "default")

	if err == nil {
		t.Fatal("Expected function to return an error, got nil")
	}
}

func TestConfigureLoggingJson(t *testing.T) {
	err := configureLogging("debug", "json")

	if err != nil {
		t.Fatalf("Function is not expected to return an error, got '%s'", err.Error())
	}
}

func TestConfigureLoggingInvalidLogFormat(t *testing.T) {
	err := configureLogging("foo", "default")

	if err == nil {
		t.Fatal("Expected function to return an error, got nil")
	}
}
