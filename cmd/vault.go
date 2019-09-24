package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
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
		return nil, fmt.Errorf("Vault endpoint must be defined")
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
			os.Exit(1)
		}

		if r.Auth == nil {
			log.Fatalf("no auth info returned with provided approle credentials")
			os.Exit(1)
		} else {
			v.SetToken(r.Auth.ClientToken)
		}
	} else {
		home, _ := homedir.Dir()
		f, err := ioutil.ReadFile(home + "/.vault-token")
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
		os.Exit(1)
	}

	if d == nil {
		log.Fatalf("The configured transit key doesn't seem to exists : %v", s.Vault.TransitKey)
		os.Exit(1)
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
		os.Exit(1)
	}

	fmt.Println("Transit key created successfully")
}

// ListTransitKeys : List available transit keys from Vault
func (v *Vault) ListTransitKeys() {
	d, err := v.Client.Logical().List("transit/keys")
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
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
	var p = make(map[string]interface{})
	p["deletion_allowed"] = "true"
	_, err := v.Client.Logical().Write("transit/keys/"+key+"/config", p)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
	}

	_, err = v.Client.Logical().Delete("transit/keys/" + key)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
	}
	color.Green("=> Deleted transit key '%v' from Vault", key)
}

// ListSecrets : Do what it says
func (v *Vault) ListSecrets() {
	log.Debugf("Listing secrets in Vault SecretPath: %v", s.Vault.SecretPath)
	d, err := v.Client.Logical().List(s.Vault.SecretPath)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
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
		os.Exit(1)
	}

	d, err := v.Client.Logical().List(s.Vault.SecretPath)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
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

// Encrypt : Encrypt a value using the TransitKey
func (v *Vault) Encrypt(value string) string {
	payload := make(map[string]interface{})
	payload["plaintext"] = base64.StdEncoding.EncodeToString([]byte(value))
	d, err := v.Client.Logical().Write("transit/encrypt/"+s.Vault.TransitKey, payload)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
	}

	return d.Data["ciphertext"].(string)
}

// Decrypt : Decrypt a value using the TransitKey
func (v *Vault) Decrypt(value string) string {
	payload := make(map[string]interface{})
	payload["ciphertext"] = value
	d, err := v.Client.Logical().Write("transit/decrypt/"+s.Vault.TransitKey, payload)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
	}
	output, _ := base64.StdEncoding.DecodeString(d.Data["plaintext"].(string))
	return string(output)
}

// WriteSecret : Write a secret into Vault
func (v *Vault) WriteSecret(secret string, payload map[string]interface{}) {
	_, err := v.Client.Logical().Write(s.Vault.SecretPath+secret, payload)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
	}
	color.Green("=> Added/Updated secret '%v' and managed keys", secret)
}

// DeleteSecret : DeleteSecret a secret from Vault
func (v *Vault) DeleteSecret(secret string) {
	_, err := v.Client.Logical().Delete(s.Vault.SecretPath + secret)
	if err != nil {
		log.Fatalf("Vault error: %v", err)
		os.Exit(1)
	}
	color.Green("=> Deleted secret '%v' and its underlying keys", secret)
}

// DeleteSecretKey : Delete a key of a secret from Vault
func (v *Vault) DeleteSecretKey(secret, key string) {
	//TODO: Implement
	color.Green("=> Deleted secret:key %v:%v", secret, key)
}
