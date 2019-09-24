package cmd

import (
	"testing"
)

var state State

func getTestStateClient() *State {
	return getStateClient(&StateConfig{
		Path: "/tmp/test.yml",
	})
}

func TestStateInit(t *testing.T) {
	s := getTestStateClient()
	s.Init()

	if s.VaultSecretPath() != "secret/" {
		t.Fatalf("Expected s.VaultSecretPath() to return secret/, got '%v'", s.VaultSecretPath())
	}
}

func TestStateSetVaultTransitKey(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultTransitKey("foo")
	s.Load()

	if s.VaultTransitKey() != "foo" {
		t.Fatalf("Expected s.SetVaultTransitKey('foo') to set s.Vault.TransitKey to foo, got '%v'", s.VaultTransitKey())
	}
}

func TestStateVaultTransitKey(t *testing.T) {
	s := getTestStateClient()
	s.Vault.TransitKey = "foo"
	if s.VaultTransitKey() != "foo" {
		t.Fatalf("Expected s.VaultTransitKey() to return foo, got '%v'", s.VaultTransitKey())
	}
}

func TestStateSetVaultSecretPath(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultSecretPath("secret/foo/")
	s.Load()

	if s.VaultSecretPath() != "secret/foo/" {
		t.Fatalf("Expected s.SetVaultSecretPath('secret/foo/') to set s.Vault.SecretPath to secret/foo/, got '%v'", s.VaultSecretPath())
	}
}

func TestStateVaultSecretPath(t *testing.T) {
	s := getTestStateClient()
	s.Vault.SecretPath = "secret/foo/"
	if s.VaultSecretPath() != "secret/foo/" {
		t.Fatalf("Expected s.VaultSecretPath() to return secret/foo/, got '%v'", s.VaultSecretPath())
	}
}

func TestStateLoad(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultTransitKey("foo")
	s.SetVaultSecretPath("secret/foo/")

	if s.VaultTransitKey() != "foo" {
		t.Fatalf("Expected s.VaultTransitKey() to return foo, got '%v'", s.VaultTransitKey())
	}

	if s.VaultSecretPath() != "secret/foo/" {
		t.Fatalf("Expected s.VaultSecretPath() to return secret/foo/, got '%v'", s.VaultSecretPath())
	}
}

func TestStateWriteSecretKey(t *testing.T) {
	s := getTestStateClient()
	s.Init()
	s.WriteSecretKey("foo", "bar", "sensitive")

	s.Load()
	if s.Secrets["foo"]["bar"] != "sensitive" {
		t.Fatalf("Expected s.Secrets['foo']['bar] to equal 'sensitive', got '%v'", s.Secrets["foo"]["bar"])
	}
}
