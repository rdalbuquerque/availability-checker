package credentialprovider

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type HcpVaultCredentialProvider struct {
	client *api.Client
}

func (v *HcpVaultCredentialProvider) Authenticate() error {
	hcpvaultAddr := os.Getenv("HCPVAULT_ADDR")
	if hcpvaultAddr == "" {
		return errors.New("environment variable 'HCPVAULT_ADDR' not set")
	}
	client, err := api.NewClient(&api.Config{
		Address: hcpvaultAddr, // Replace with your Vault server address
	})
	if err != nil {
		return err
	}

	// Assuming token-based authentication
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return errors.New("environment variable VAULT_TOKEN not set")
	}

	client.SetToken(token)
	client.SetNamespace("admin")

	v.client = client
	return nil
}

func (v *HcpVaultCredentialProvider) GetCredentials(checkerType string) (user, password string, err error) {
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
	pass, passOk := data["pwd"].(string)

	if !userOk || !passOk {
		return "", "", errors.New("user and password not found in secret data")
	}

	return user, pass, nil
}
