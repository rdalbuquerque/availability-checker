package credentialprovider

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultCredentialProvider struct {
	client *api.Client
}

func (v *VaultCredentialProvider) Authenticate() error {
	client, err := api.NewClient(&api.Config{
		Address: "https://vault-public-vault-110e8125.d7c3b38b.z1.hashicorp.cloud:8200", // Replace with your Vault server address
	})
	if err != nil {
		return err
	}

	// Assuming token-based authentication
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return errors.New("VAULT_TOKEN environment variable not set")
	}

	client.SetToken(token)

	v.client = client
	return nil
}

func (v *VaultCredentialProvider) GetCredentials(checkerType string) (user, password string, err error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("secret/data/%s", checkerType))
	if err != nil {
		return "", "", err
	}

	if secret == nil || secret.Data == nil {
		return "", "", errors.New("no secret found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("malformed secret data")
	}

	user, userOk := data["user"].(string)
	pass, passOk := data["password"].(string)

	if !userOk || !passOk {
		return "", "", errors.New("user and password not found in secret data")
	}

	return user, pass, nil
}
