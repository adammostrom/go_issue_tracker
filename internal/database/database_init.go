package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "issues"
const SCHEMA_FILE = "internal/database/schema_issues.sql"

// SQLITE

func panic_mode(err error) {
	if err != nil {
		panic(err)
	}
}

func NewDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "issues.db")
	panic_mode(err)

	err = db.Ping()
	if err != nil {
		fmt.Printf("Database not reached with error: %s\n", err)
		panic(err)
	}
	fmt.Printf("Created DB instance with name: %s", DB_NAME+"db\n")
	return db, nil
}

// Loads and executes the schema
func InitSchema(db *sql.DB) error {
	schema, err := os.ReadFile(SCHEMA_FILE)
	panic_mode(err)

	_, err = db.Exec(string(schema))

	fmt.Printf("Schema: %s initiated.\n", SCHEMA_FILE)
	return err
}

func InitDB() *sql.DB {
	fmt.Println("HELLO")
	db, err := NewDB()
	panic_mode(err)

	err = InitSchema(db)
	panic_mode(err)

	fmt.Println("Database initiated successfully")
	return db
}

// POSTGRES
// OpenDB opens the DB and returns *sql.DB (connection pool)
/* func OpenDB(connection string) (*sql.DB, error) {
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
} */
