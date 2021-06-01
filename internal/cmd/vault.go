package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"

	s5 "github.com/mvisonneau/s5/cipher"
	s5Vault "github.com/mvisonneau/s5/cipher/vault"
)

// Vault : Handles a Vault API Client
type Vault struct {
	Client *api.Client
}

// VaultConfig handles Vault configuration
type VaultConfig struct {
	Address  string
	Token    string
	RoleID   string
	SecretID string
}

func getVaultClient(vc *VaultConfig) (*Vault, error) {
	log.Debug("Creating Vault client..")
	v, err := api.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating Vault client: %v", err)
	}

	if vc.Address == "" {
		return nil, fmt.Errorf("Vault address must be defined")
	}

	if len(vc.Token) > 0 {
		v.SetToken(vc.Token)
	} else if len(vc.RoleID) > 0 && len(vc.SecretID) > 0 {
		data := map[string]interface{}{
			"role_id":   vc.RoleID,
			"secret_id": vc.SecretID,
		}

		r, err := v.Logical().Write("auth/approle/login", data)
		if err != nil {
			log.Fatalf("Can't authenticate against vault using provided approle credentials: %v", err)
		}

		if r.Auth == nil {
			log.Fatalf("no auth info returned with provided approle credentials")
		}

		v.SetToken(r.Auth.ClientToken)
	} else {
		home, _ := homedir.Dir()
		f, err := ioutil.ReadFile(filepath.Clean(home + "/.vault-token"))
		if err != nil {
			return nil, fmt.Errorf("Vault token is not defined (VAULT_TOKEN, (--vault-role-id and --vault-secret-id) or ~/.vault-token)")
		}

		v.SetToken(string(f))
	}

	return &Vault{v}, nil
}

// GetTransitInfo : Fetch some information from Vault about the configured TransitKey
func (v *Vault) GetTransitInfo() {
	d, err := v.Client.Logical().Read("transit/keys/" + s.Vault.TransitKey)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}

	if d == nil {
		log.Fatalf("The configured transit key doesn't seem to exists : %v", s.Vault.TransitKey)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	for k, l := range d.Data {
		table.Append([]string{k, fmt.Sprintf("%v", l)})
	}
	table.Render()
}

// CreateTransitKey : Create a new transit key in Vault
func (v *Vault) CreateTransitKey(key string) {
	_, err := v.Client.Logical().Write("transit/keys/"+key, make(map[string]interface{}))
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}

	fmt.Println("Transit key created successfully")
}

// ListTransitKeys : List available transit keys from Vault
func (v *Vault) ListTransitKeys() {
	d, err := v.Client.Logical().List("transit/keys")
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}

	if d != nil {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Key"})
		for _, l := range d.Data["keys"].([]interface{}) {
			table.Append([]string{l.(string)})
		}
		table.Render()
	}
}

// DeleteTransitKey : Delete a transit key from Vault
func (v *Vault) DeleteTransitKey(key string) {
	p := make(map[string]interface{})
	p["deletion_allowed"] = "true"
	_, err := v.Client.Logical().Write("transit/keys/"+key+"/config", p)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}

	_, err = v.Client.Logical().Delete("transit/keys/" + key)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}
	color.Green("=> Deleted transit key '%v' from Vault", key)
}

// ListSecrets : Do what it says
func (v *Vault) ListSecrets() {
	log.Debugf("Listing secrets in Vault KV Path: %v", s.VaultKVPath())
	d, err := v.Client.Logical().List(s.VaultKVPath())
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key"})
	for _, l := range d.Data["keys"].([]interface{}) {
		table.Append([]string{l.(string)})
	}
	table.Render()
}

// Status : Return information about Vault API endpoint/cluster
func (v *Vault) Status() {
	vh, err := v.Client.Sys().Health()
	fmt.Println("[VAULT]")
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}

	listPath := s.VaultKVPath()
	if s.VaultKVVersion() == 2 {
		listPath = s.VaultKVPath() + "metadata"
	}

	d, err := v.Client.Logical().List(listPath)
	if err != nil {
		log.Fatalf("vault error: %v", err)
	}

	secretsCount := 0
	if d != nil {
		if keys, ok := d.Data["keys"]; ok {
			secretsCount = len(keys.([]interface{}))
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Sealed", fmt.Sprintf("%v", vh.Sealed)})
	table.Append([]string{"Cluster Version", vh.Version})
	table.Append([]string{"Cluster ID", vh.ClusterID})
	table.Append([]string{"Secrets #", fmt.Sprintf("%v", secretsCount)})
	table.Render()
}

// Cipher : Cipher a value using the TransitKey
func (v *Vault) Cipher(value string) string {
	s5Engine := s5Vault.Client{
		Client: v.Client,
		Config: &s5Vault.Config{
			Key: s.Vault.TransitKey,
		},
	}

	cipheredValue, err := s5Engine.Cipher(value)
	if err != nil {
		log.Fatal(err)
	}

	return s5.GenerateOutput(cipheredValue)
}

// Decipher : Decipher a value using the TransitKey
func (v *Vault) Decipher(value string) string {
	s5Engine := s5Vault.Client{
		Client: v.Client,
		Config: &s5Vault.Config{
			Key: s.Vault.TransitKey,
		},
	}

	parsedInput, err := s5.ParseInput(value)
	if err != nil {
		log.Fatal(err)
	}

	decipheredValue, err := s5Engine.Decipher(parsedInput)
	if err != nil {
		log.Fatal(err)
	}

	return decipheredValue
}

// WriteSecret : Write a secret into Vault
func (v *Vault) WriteSecret(secret string, data map[string]interface{}) {
	queryPath := s.VaultKVPath() + secret
	var payload map[string]interface{}
	if s.Vault.KV.Version == 2 {
		queryPath = s.VaultKVPath() + "data/" + secret
		payload = make(map[string]interface{})
		payload["data"] = data
	} else {
		payload = data
	}

	_, err := v.Client.Logical().Write(queryPath, payload)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}
	color.Green("=> Added/Updated secret '%v' and managed keys", secret)
}

// DeleteSecret : DeleteSecret a secret from Vault
func (v *Vault) DeleteSecret(secret string) {
	_, err := v.Client.Logical().Delete(s.VaultKVPath() + secret)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
	}
	color.Green("=> Deleted secret '%v' and its underlying keys", secret)
}

// DeleteSecretKey : Delete a key of a secret from Vault
func (v *Vault) DeleteSecretKey(secret, key string) {
	// TODO: Implement!
	color.Yellow("=> [NOT IMPLEMENTED] - Deleted secret:key %v:%v", secret, key)
}
