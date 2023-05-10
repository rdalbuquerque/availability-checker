package checker

import (
	"database/sql"
)

type SqlChecker struct {
	ConnectionString string
}

func (c *SqlChecker) Check() (bool, error) {
	db, err := sql.Open("sqlserver", c.ConnectionString)
	if err != nil {
		return false, err
	}
	defer db.Close()

	err = db.Ping()
	return err == nil, err
}

func (c *SqlChecker) Name() string {
	return "SQL Server"
}
