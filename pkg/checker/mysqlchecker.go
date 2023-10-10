package checker

import (
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/database"
	"availability-checker/pkg/k8s"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLChecker struct {
	Server             string
	Port               string
	DBConnection       database.DBConnection
	CredentialProvider credentialprovider.CredentialProvider
	K8sClient          k8s.K8sClient
}

func (c *MySQLChecker) Name() string {
	return fmt.Sprintf("MySQL: %s:%s", c.Server, c.Port)
}

func (c *MySQLChecker) Check() (bool, error) {
	user, pwd, err := c.CredentialProvider.GetCredentials("mysql")
	if err != nil {
		return false, fmt.Errorf("error getting credentials: %v", err)
	}

	if user == "" || pwd == "" {
		return false, errors.New("empty username or password")
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pwd, c.Server, c.Port)

	err = c.DBConnection.Open("mysql", connectionString)
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

func (c *MySQLChecker) Fix() error {
	return c.K8sClient.ScaleDeploymentToDesiredReplicas("default", "mysql", 1)
}

func (c *MySQLChecker) IsFixable() bool {
	return true
}
