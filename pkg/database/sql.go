package database

import "database/sql"

type SQLDBConnection struct {
	*sql.DB
}

func (s *SQLDBConnection) Open(driverName, dataSourceName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

func (s *SQLDBConnection) Ping() error {
	return s.DB.Ping()
}

func (s *SQLDBConnection) Close() error {
	return s.DB.Close()
}
