package connection

import "github.com/jackc/pgx"

var config = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "forum",
	User:     "postgres",
	Password: "docker",
}

func InitDBConnection() (*pgx.ConnPool, error) {
	var err error
	connection, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 100,
		})
	return connection, err
}