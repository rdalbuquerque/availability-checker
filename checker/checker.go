package checker

import (
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/vertica/vertica-sql-go"
)

type Checker interface {
	Check() (bool, error)
	Name() string
	Fix() error
	IsFixable() bool
}

type CheckResult struct {
	Name        string
	Status      bool
	LastChecked time.Time
	IsFixable   bool
}
