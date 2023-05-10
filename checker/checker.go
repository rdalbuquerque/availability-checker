package checker

import (
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/vertica/vertica-sql-go"
)

type Checker interface {
	Check() (bool, error)
	Name() string
}

type CheckResult struct {
	Name   string
	Status bool
}
