package checker

import (
	"time"
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
