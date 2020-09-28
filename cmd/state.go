package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// State : Handles state information
type State struct {
	Vault struct {
		TransitKey string
		KV         struct {
			Path    string
			Version int
		}
	}
	Secrets map[string]map[string]string
	Config  *StateConfig `yaml:"-"`
}

// StateConfig handles state client configuration
type StateConfig struct {
	Path string
}

func getStateClient(sc *StateConfig) *State {
	return &State{
		Config: sc,
	}
}

// Init : Generates an empty state file at the configured state file location
func (s *State) Init() {
	log.Infof("Creating an empty state file at %v", s.Config.Path)
	s.SetVaultTransitKey("default")
	s.SetVaultKVPath("secret/")
	s.SetVaultKVVersion(2)
	s.save()
}

// SetVaultTransitKey : Update state file with a Vault/TransitKey value
func (s *State) SetVaultTransitKey(value string) {
	s.Vault.TransitKey = value
	s.save()
}

// VaultTransitKey : Returns the value of the configured Vault/TransitKey
func (s *State) VaultTransitKey() string {
	return s.Vault.TransitKey
}

// SetVaultKVPath : Update state file with a Vault/Secret/Path value
func (s *State) SetVaultKVPath(value string) {
	s.Vault.KV.Path = value
	s.save()
}

// VaultKVPath : Returns the value of the configured Vault/Secret/Path
func (s *State) VaultKVPath() string {
	if s.Vault.KV.Path == "" {
		return "secret/"
	}
	return s.Vault.KV.Path
}

// SetVaultKVVersion : Update state file with a Vault/Secret/Version value
func (s *State) SetVaultKVVersion(version int) {
	s.Vault.KV.Version = version
	s.save()
}

// VaultKVVersion : Returns the value of the configured Vault/Secret/Version
func (s *State) VaultKVVersion() int {
	if s.Vault.KV.Version == 0 {
		return 1
	}
	return s.Vault.KV.Version
}

// Load : Loads the statefile content in memory
func (s *State) Load() {
	if s.Config.Path == "" {
		log.Fatal("State file must be defined")
	}

	log.Debugf("Loading from statefile: %v", s.Config.Path)
	if _, err := os.Stat(s.Config.Path); os.IsNotExist(err) {
		log.Fatalf("State file not found at location: %s, use 'strongbox init' to generate an empty one.\n", s.Config.Path)
	}

	filename, _ := filepath.Abs(s.Config.Path)
	data, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		log.Fatal("Error: State file not found, create a new one using : 'strongbox init'")
	}

	err = yaml.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Debugf("Loaded Transit Key: %v", s.VaultTransitKey())
	log.Debugf("Loaded KV Path: %v", s.VaultKVPath())
	log.Debugf("Loaded KV Version: %v", s.VaultKVVersion())
	log.Debugf("Loaded Secrets: %#v", s.Secrets)
}

// Status : Returns information about statefile content
func (s *State) Status() {
	fmt.Println("[STRONGBOX STATE]")
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Transit Key", s.VaultTransitKey()})
	table.Append([]string{"KV Path", s.VaultKVPath()})
	table.Append([]string{"KV Version", strconv.Itoa(s.VaultKVVersion())})
	table.Append([]string{"Secrets #", fmt.Sprintf("%v", len(s.Secrets))})
	table.Render()
}

// ListSecrets : List the secrets, safely stored into the statefile
func (s *State) ListSecrets(secret string) {
	log.Debug("Rendering local secrets list")

	if secret == "" {
		for k, l := range s.Secrets {
			table := tablewriter.NewWriter(os.Stdout)
			fmt.Printf("[%v]\n", k)
			for m, n := range l {
				table.Append([]string{m, n})
			}
			table.Render()
		}
	} else {
		if s.Secrets[secret] == nil {
			fmt.Printf("No secret '%v' found\n", secret)
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		for k, l := range s.Secrets[secret] {
			table.Append([]string{k, l})
		}
		table.Render()
	}
}

// WriteSecretKey : Add or Update a key value within a secret
func (s *State) WriteSecretKey(secret, key, value string) {
	if s.Secrets == nil {
		s.Secrets = map[string]map[string]string{}
	}

	if s.Secrets[secret] == nil {
		s.Secrets[secret] = map[string]string{}
	}

	s.Secrets[secret][key] = value
	s.save()
}

// ReadSecretKey : Read the value of a SecretKey
func (s *State) ReadSecretKey(secret, key string) string {
	if s.Secrets == nil || s.Secrets[secret] == nil {
		fmt.Printf("No secret '%v' found\n", secret)
		os.Exit(1)
	}

	if s.Secrets[secret][key] == "" {
		fmt.Printf("No key '%v' found in secret '%v'\n", key, secret)
		os.Exit(1)
	}

	return s.Secrets[secret][key]
}

// DeleteSecret : Delete a secret from the statefile based on its name
func (s *State) DeleteSecret(secret string) {
	if s.Secrets == nil || s.Secrets[secret] == nil {
		fmt.Printf("No secret '%v' found\n", secret)
		os.Exit(1)
	}

	delete(s.Secrets, secret)
	s.save()
	fmt.Println("Secret deleted!")
}

// DeleteSecretKey : Delete a secret:key from the statefile based on the secret and key names
func (s *State) DeleteSecretKey(secret, key string) {
	if s.Secrets == nil || s.Secrets[secret] == nil {
		fmt.Printf("No secret '%v' found\n", secret)
		os.Exit(1)
	}

	if s.Secrets[secret][key] == "" {
		fmt.Printf("No key '%v' found in secret '%v'\n", key, secret)
		os.Exit(1)
	}

	delete(s.Secrets[secret], key)
	s.save()
	fmt.Println("Key deleted!")
}

// RotateFromOldTransitKey : Replace locally ciphered values with new transit key
func (s *State) RotateFromOldTransitKey(key string) {
	transitKey := s.VaultTransitKey()
	if transitKey == key {
		log.Fatalf("%v is already the currently configured key, can't rotate with same key", key)
	}

	s.SetVaultTransitKey(key)
	secrets := make(map[string]map[string]string)
	for k, l := range s.Secrets {
		if secrets[k] == nil {
			secrets[k] = make(map[string]string)
		}
		for m, n := range l {
			secrets[k][m] = v.Decipher(n)
		}
	}

	s.SetVaultTransitKey(transitKey)

	for k, l := range secrets {
		for m, n := range l {
			s.WriteSecretKey(k, m, v.Cipher(n))
		}
	}
	fmt.Printf("Rotated secrets from '%v' to '%v'\n", key, transitKey)
}

// save : write the statefile onto the disk
func (s *State) save() {
	log.Debugf("Saving state file at %v", s.Config.Path)
	var output bytes.Buffer
	y := yaml.NewEncoder(&output)
	y.SetIndent(2)
	if err := y.Encode(&s); err != nil {
		log.Fatalf("Error: %v", err)
	}

	filename, err := filepath.Abs(s.Config.Path)
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(filename, output.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
}
