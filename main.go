package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"availability-checker/pkg/checker"
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/database"
	"availability-checker/pkg/k8s"
	"availability-checker/pkg/server"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Checkers []struct {
		Type   string
		URL    string `yaml:",omitempty"`
		Server string `yaml:"server,omitempty"`
		Port   string `yaml:"port,omitempty"`
	}
}

func main() {
	credProvider, err := credentialProviderAuth()
	if err != nil {
		log.Fatalf("Error authenticating credential provider: %v", err)
	}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	var config Config
	yaml.Unmarshal(data, &config)

	k8sclient, err := k8s.NewK8sClient()
	if err != nil {
		log.Fatalf("Error creating k8s client: %v", err)
	}
	checkers := make([]checker.Checker, len(config.Checkers))
	for i, confChecker := range config.Checkers {
		switch confChecker.Type {
		case "http":
			checkers[i] = &checker.HttpChecker{URL: confChecker.URL}
		case "postgres":
			checkers[i] = &checker.PostgresChecker{
				Server:             confChecker.Server,
				Port:               confChecker.Port,
				DBConnection:       &database.SQLDBConnection{},
				CredentialProvider: credProvider,
				K8sClient:          *k8sclient,
			}
		case "mysql":
			checkers[i] = &checker.MySQLChecker{
				Server:             confChecker.Server,
				Port:               confChecker.Port,
				DBConnection:       &database.SQLDBConnection{},
				CredentialProvider: credProvider,
				K8sClient:          *k8sclient,
			}
		}
	}

	serverInstance := server.NewServer(checkers, "template.gotmpl")

	go serverInstance.StartChecking()

	http.ListenAndServe(":8080", serverInstance)
}

func credentialProviderAuth() (credentialprovider.CredentialProvider, error) {
	var credProvider credentialprovider.CredentialProvider
	if os.Getenv("AZURE_KEYVAULT") != "" {
		credProvider = &credentialprovider.AzureKeyVaultCredentialProvider{}
	} else if os.Getenv("HCPVAULT_ADDR") != "" {
		credProvider = &credentialprovider.HcpVaultCredentialProvider{}
	} else {
		return nil, fmt.Errorf("no credential provider environment variable set")
	}
	err := credProvider.Authenticate()
	if err != nil {
		return nil, err
	}
	return credProvider, nil
}
