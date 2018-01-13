package app

import (
	"testing"
)

func TestCli(t *testing.T) {
	c := Cli()
	if c.Name != "strongbox" {
		t.Fatalf("Expected c.Name to be strongbox, got '%v'", c.Name)
	}
}
