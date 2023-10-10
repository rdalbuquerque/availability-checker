package checker

import (
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/database"
	"availability-checker/pkg/k8s"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresChecker struct {
	Server             string
	Port               string
	DBConnection       database.DBConnection
	CredentialProvider credentialprovider.CredentialProvider
	K8sClient          k8s.K8sClient
}

func (c *PostgresChecker) Name() string {
	return fmt.Sprintf("Postgres: %s:%s", c.Server, c.Port)
}

func (c *PostgresChecker) Check() (bool, error) {
	user, pwd, err := c.CredentialProvider.GetCredentials("postgres")
	if err != nil {
		return false, fmt.Errorf("error getting credentials: %v", err)
	}

	if user == "" || pwd == "" {
		return false, errors.New("empty username or password")
	}

	connectionString := fmt.Sprintf("host=%s port=%s dbname=postgres user=%s password=%s sslmode=disable connect_timeout=10", c.Server, c.Port, user, pwd)

	fmt.Printf("Connecting to postgres: %s:%s\n", c.Server, c.Port)
	err = c.DBConnection.Open("postgres", connectionString)
	if err != nil {
		fmt.Printf("Error opening connection: %v\n", err)
		return false, err
	}
	defer c.DBConnection.Close()

	err = c.DBConnection.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return false, err
	}

	return true, nil
}

func (c *PostgresChecker) Fix() error {
	return c.K8sClient.ScaleDeploymentToDesiredReplicas("default", "postgres", 1)
}

func (c *PostgresChecker) IsFixable() bool {
	return true
}
