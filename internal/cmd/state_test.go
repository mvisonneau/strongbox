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

	if s.VaultKVPath() != "secret/" {
		t.Fatalf("Expected s.VaultKVPath() to return secret/, got '%v'", s.VaultKVPath())
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

func TestStateSetVaultKVPath(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultKVPath("secret/foo/")
	s.Load()

	if s.VaultKVPath() != "secret/foo/" {
		t.Fatalf("Expected s.SetVaultKVPath('secret/foo/') to set s.VaultKVPath() to secret/foo/, got '%v'", s.VaultKVPath())
	}
}

func TestStateVaultKVPath(t *testing.T) {
	s := getTestStateClient()
	s.Vault.KV.Path = "secret/foo/"
	if s.VaultKVPath() != "secret/foo/" {
		t.Fatalf("Expected s.VaultKVPath() to return secret/foo/, got '%v'", s.VaultKVPath())
	}
}

func TestStateLoad(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultTransitKey("foo")
	s.SetVaultKVPath("secret/foo/")

	if s.VaultTransitKey() != "foo" {
		t.Fatalf("Expected s.VaultTransitKey() to return foo, got '%v'", s.VaultTransitKey())
	}

	if s.VaultKVPath() != "secret/foo/" {
		t.Fatalf("Expected s.VaultKVPath() to return secret/foo/, got '%v'", s.VaultKVPath())
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
