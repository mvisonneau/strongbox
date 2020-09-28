package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// State : Handles state information
type State struct {
	Vault struct {
		TransitKey string
		SecretPath string
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
	s.SetVaultSecretPath("secret/")
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

// SetVaultSecretPath : Update state file with a Vault/SecretPath value
func (s *State) SetVaultSecretPath(value string) {
	s.Vault.SecretPath = value
	s.save()
}

// VaultSecretPath : Returns the value of the configured Vault/SecretPath
func (s *State) VaultSecretPath() string {
	return s.Vault.SecretPath
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

	log.Debugf("Loaded TransitKey: %v", s.Vault.TransitKey)
	log.Debugf("Loaded SecretPath: %v", s.Vault.SecretPath)
	log.Debugf("Loaded Secrets: %#v", s.Secrets)
}

// Status : Returns information about statefile content
func (s *State) Status() {
	fmt.Println("[STATE]")
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"TransitKey", s.Vault.TransitKey})
	table.Append([]string{"SecretPath", s.Vault.SecretPath})
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

// RotateFromOldTransitKey : Replace local encrypted values with new transit key
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
			secrets[k][m] = v.Decrypt(n)
		}
	}

	s.SetVaultTransitKey(transitKey)

	for k, l := range secrets {
		for m, n := range l {
			s.WriteSecretKey(k, m, v.Encrypt(n))
		}
	}
	fmt.Printf("Rotated secrets from '%v' to '%v'\n", key, transitKey)
}

// save : write the statefile onto the disk
func (s *State) save() {
	log.Debugf("Saving state file at %v", s.Config.Path)
	output, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	filename, err := filepath.Abs(s.Config.Path)
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(filename, output, 0600); err != nil {
		log.Fatal(err)
	}
}
