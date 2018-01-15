package app

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// Vault : Handles a Vault API Client
type Vault struct {
	Client *api.Client
}

// ConfigureClient : Instanciate and configure a vault client
func (v *Vault) ConfigureClient() {
	log.Debug("Creating Vault client..")
	var err error
	v.Client, err = api.NewClient(nil)
	if err != nil {
		log.Fatalf("Error creating Vault client: %v", err)
		os.Exit(1)
	}

	if cfg.Vault.Address == "" {
		log.Fatal("Vault endpoint must be defined")
		os.Exit(1)
	}

	if cfg.Vault.Token == "" &&
		(cfg.Vault.RoleID == "" || cfg.Vault.SecretID == "") {
		log.Fatal("Either vault-token or (vault-role-id and vault-secret-id) must be defined")
		os.Exit(1)
	}

	v.Client.SetAddress(cfg.Vault.Address)

	if cfg.Vault.Token != "" {
		v.Client.SetToken(cfg.Vault.Token)
	} else {
		data := map[string]interface{}{
			"role_id":   cfg.Vault.RoleID,
			"secret_id": cfg.Vault.SecretID,
		}

		r, err := v.Client.Logical().Write("auth/approle/login", data)
		if err != nil {
			log.Fatalf("Can't authenticate against vault using provided approle credentials: %v", err)
			os.Exit(1)
		}

		if r.Auth == nil {
			log.Fatalf("no auth info returned with provided approle credentials")
			os.Exit(1)
		} else {
			v.Client.SetToken(r.Auth.ClientToken)
		}
	}
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
	var payload map[string]interface{}
	_, err := v.Client.Logical().Write("transit/keys/"+key, payload)
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key"})
	for _, l := range d.Data["keys"].([]interface{}) {
		table.Append([]string{l.(string)})
	}
	table.Render()
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
		secretsCount = len(d.Data["keys"].([]interface{}))
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
