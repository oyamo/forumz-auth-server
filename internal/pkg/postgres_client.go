package pkg

import (
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresClient(dsn string) (*sql.DB, error) {
	pgxConf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conn := stdlib.OpenDB(*pgxConf)

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
