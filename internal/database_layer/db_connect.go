package db

import (
	"database/sql"
	"fmt"
)

// OpenDB opens the DB and returns *sql.DB (connection pool)
func OpenDB(connection string) (*sql.DB, error) {
	dsn := fmt.Sprintf(connection)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Force a real connection to test
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
