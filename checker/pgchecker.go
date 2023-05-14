package checker

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresChecker struct {
	ConnectionString string
}

func (c *PostgresChecker) Name() string {
	return fmt.Sprintf("Postgres: %s", c.ConnectionString)
}

func (c *PostgresChecker) Check() (bool, error) {
	fmt.Printf("Connecting to postgres with connection string: %s\n", c.ConnectionString)
	db, err := sql.Open("postgres", c.ConnectionString)
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
