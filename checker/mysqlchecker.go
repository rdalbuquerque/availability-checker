package checker

import (
	"availability-checker/credentialprovider"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLChecker struct {
	Server             string
	Port               string
	CredentialProvider credentialprovider.CredentialProvider
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

	db, err := sql.Open("mysql", connectionString)
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

func (c *MySQLChecker) Fix() error {
	return nil
}

func (c *MySQLChecker) IsFixable() bool {
	return false
}
