package credentialprovider

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

type AzureKeyVaultCredentialProvider struct {
	vaultName string
	client    keyvault.BaseClient
}

func (a *AzureKeyVaultCredentialProvider) Authenticate() error {
	client := keyvault.New()
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}
	client.Authorizer = authorizer
	a.client = client

	keyvault := os.Getenv("AZURE_KEYVAULT")
	if keyvault == "" {
		return errors.New("environment variable 'AZURE_KEYVAULT' not set")
	}
	a.vaultName = keyvault
	return nil
}

func (a *AzureKeyVaultCredentialProvider) GetCredentials(checkerType string) (user, password string, err error) {
	vaultURL := fmt.Sprintf("https://%s.vault.azure.net", a.vaultName)

	// Replace "usernameSecretName" and "passwordSecretName" with the actual secret names
	userSecretBundle, err := a.client.GetSecret(context.TODO(), vaultURL, fmt.Sprintf("%s-user", checkerType), "")
	if err != nil {
		return "", "", err
	}

	passSecretBundle, err := a.client.GetSecret(context.TODO(), vaultURL, fmt.Sprintf("%s-pwd", checkerType), "")
	if err != nil {
		return "", "", err
	}

	if userSecretBundle.Value == nil || passSecretBundle.Value == nil {
		return "", "", errors.New("user and password not found in secret data")
	}

	return *userSecretBundle.Value, *passSecretBundle.Value, nil
}
