package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"availability-checker/pkg/checker"
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/database"
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
	credProvider := &credentialprovider.AzureKeyVaultCredentialProvider{}
	err := credProvider.Authenticate()
	if err != nil {
		log.Fatalf("Error authenticating credential provider: %v", err)
	}
	data, _ := ioutil.ReadFile("config.yaml")
	var config Config
	yaml.Unmarshal(data, &config)

	checkers := make([]checker.Checker, len(config.Checkers))
	for i, confChecker := range config.Checkers {
		switch confChecker.Type {
		case "http":
			checkers[i] = &checker.HttpChecker{URL: confChecker.URL}
		case "postgres":
			checkers[i] = &checker.PostgresChecker{Server: confChecker.Server, Port: confChecker.Port, DBConnection: &database.SQLDBConnection{}, CredentialProvider: credProvider}
		case "mysql":
			checkers[i] = &checker.MySQLChecker{Server: confChecker.Server, Port: confChecker.Port, DBConnection: &database.SQLDBConnection{}, CredentialProvider: credProvider}
		}
	}

	serverInstance := server.NewServer(checkers, "template.gotmpl")

	go serverInstance.StartChecking()

	http.ListenAndServe(":8080", serverInstance)
}
