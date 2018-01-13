package app

import (
	"testing"
)

var state State

func TestStateInit(t *testing.T) {
	cfg.State.Location = "/tmp/test.yml"
	state.Init()

	if state.VaultSecretPath() != "secret/" {
		t.Fatalf("Expected state.VaultSecretPath() to return secret/, got '%v'", state.VaultSecretPath())
	}
}

func TestStateSetVaultTransitKey(t *testing.T) {
	cfg.State.Location = "/tmp/test.yml"
	state.SetVaultTransitKey("foo")
	state.Load()

	if state.VaultTransitKey() != "foo" {
		t.Fatalf("Expected state.SetVaultTransitKey('foo') to set state.Vault.TransitKey to foo, got '%v'", state.VaultTransitKey())
	}
}

func TestStateVaultTransitKey(t *testing.T) {
	state.Vault.TransitKey = "foo"
	if state.VaultTransitKey() != "foo" {
		t.Fatalf("Expected state.VaultTransitKey() to return foo, got '%v'", state.VaultTransitKey())
	}
}

func TestStateSetVaultSecretPath(t *testing.T) {
	cfg.State.Location = "/tmp/test.yml"
	state.SetVaultSecretPath("secret/foo/")
	state.Load()

	if state.VaultSecretPath() != "secret/foo/" {
		t.Fatalf("Expected state.SetVaultSecretPath('secret/foo/') to set state.Vault.SecretPath to secret/foo/, got '%v'", state.VaultSecretPath())
	}
}

func TestStateVaultSecretPath(t *testing.T) {
	state.Vault.SecretPath = "secret/foo/"
	if state.VaultSecretPath() != "secret/foo/" {
		t.Fatalf("Expected state.VaultSecretPath() to return secret/foo/, got '%v'", state.VaultSecretPath())
	}
}

func TestStateLoad(t *testing.T) {
	cfg.State.Location = "/tmp/test.yml"
	state.SetVaultTransitKey("foo")
	state.SetVaultSecretPath("secret/foo/")

	if state.VaultTransitKey() != "foo" {
		t.Fatalf("Expected state.VaultTransitKey() to return foo, got '%v'", state.VaultTransitKey())
	}

	if state.VaultSecretPath() != "secret/foo/" {
		t.Fatalf("Expected state.VaultSecretPath() to return secret/foo/, got '%v'", state.VaultSecretPath())
	}
}

func TestStateWriteSecretKey(t *testing.T) {
	cfg.State.Location = "/tmp/test.yml"
	state.Init()
	state.WriteSecretKey("foo", "bar", "sensitive")

	state.Load()
	if state.Secrets["foo"]["bar"] != "sensitive" {
		t.Fatalf("Expected state.Secrets['foo']['bar] to equal 'sensitive', got '%v'", state.Secrets["foo"]["bar"])
	}
}
