package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tmpDir string
	state  State
)

// mainWrapper is necessary to be able to leverage "defer"
// as os.Exit does not honour it and is required by *testing.M
func mainWrapper(m *testing.M) int {
	var err error
	if tmpDir, err = ioutil.TempDir(".", ".test"); err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func getTestStateClient() *State {
	var err error
	var tmpFile *os.File
	if tmpFile, err = ioutil.TempFile(tmpDir, "state"); err != nil {
		log.Fatal(err)
	}
	return getStateClient(&StateConfig{
		Path: tmpFile.Name(),
	})
}

func TestStateInit(t *testing.T) {
	s := getTestStateClient()
	s.Init()
	assert.Equal(t, "secret/", s.VaultKVPath())
}

func TestStateSetVaultTransitKey(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultTransitKey("foo")
	s.Load()
	assert.Equal(t, "foo", s.VaultTransitKey())
}

func TestStateVaultTransitKey(t *testing.T) {
	s := getTestStateClient()
	s.Vault.TransitKey = "foo"
	assert.Equal(t, "foo", s.VaultTransitKey())
}

func TestStateSetVaultKVPath(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultKVPath("secret/foo/")
	s.Load()
	assert.Equal(t, "secret/foo/", s.VaultKVPath())
}

func TestStateVaultKVPath(t *testing.T) {
	s := getTestStateClient()
	s.Vault.KV.Path = "secret/foo/"
	assert.Equal(t, "secret/foo/", s.VaultKVPath())
}

func TestStateLoad(t *testing.T) {
	s := getTestStateClient()
	s.SetVaultTransitKey("foo")
	s.SetVaultKVPath("secret/foo/")
	assert.Equal(t, "foo", s.VaultTransitKey())
	assert.Equal(t, "secret/foo/", s.VaultKVPath())
}

func TestStateWriteSecretKey(t *testing.T) {
	s := getTestStateClient()
	s.Init()
	s.WriteSecretKey("foo", "bar", "sensitive")
	s.Load()

	require.Contains(t, s.Secrets, "foo")
	require.Contains(t, s.Secrets["foo"], "bar")
	assert.Equal(t, "sensitive", s.Secrets["foo"]["bar"])
}
