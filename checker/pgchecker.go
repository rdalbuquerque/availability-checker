package checker

import (
	"availability-checker/credentialprovider"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/sys/windows/svc/mgr"
)

type PostgresChecker struct {
	Server             string
	Port               string
	CredentialProvider credentialprovider.CredentialProvider
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
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Printf("Error opening connection: %v\n", err)
		return false, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return false, err
	}

	return true, nil
}

func (c *PostgresChecker) Fix() error {
	serviceName := "postgresql-x64-15"

	// Connect to the Service Control Manager
	manager, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer manager.Disconnect()

	// Open the service by name
	service, err := manager.OpenService(serviceName)
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
