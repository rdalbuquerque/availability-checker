package database

type DBConnection interface {
	Open(driverName, dataSourceName string) error
	Close() error
	Ping() error
}
