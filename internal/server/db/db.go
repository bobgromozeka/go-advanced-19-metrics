package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var pgxConnection *pgx.Conn

// Connect connects to postgres DB using specified dsn.
func Connect(dsn string) error {
	if conn, connErr := pgx.Connect(context.Background(), dsn); connErr != nil {
		return connErr
	} else {
		pgxConnection = conn
	}

	return nil
}

// Connection return postgres connection.
func Connection() *pgx.Conn {
	return pgxConnection
}
