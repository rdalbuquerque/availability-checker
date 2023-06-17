package checker

import (
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/database"
	"availability-checker/pkg/winsvcmngr"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresChecker struct {
	Server             string
	Port               string
	DBConnection       database.DBConnection
	CredentialProvider credentialprovider.CredentialProvider
	WinSvcMngr         winsvcmngr.WinSvcMngr
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
	serviceName := "postgresql-x64-15"
	// Connect to the Service Control Manager
	err := c.WinSvcMngr.Connect()
	if err != nil {
		return err
	}
	defer c.WinSvcMngr.Disconnect()

	// Open the service by name
	service, err := c.WinSvcMngr.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()

	// Start the service
	err = service.Start()
	if err != nil {
		return err
	}

	return nil
}

func (c *PostgresChecker) IsFixable() bool {
	return true
}
