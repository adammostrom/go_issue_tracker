package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "issues"
const SCHEMA_FILE = "internal/database/schema_issues.sql"
const DB_FOLDER = ".issuetracker"

// SQLITE

func panic_mode(err error) {
	if err != nil {
		panic(err)
	}
}

// Basically Run the DB
func OpenDB() (*sql.DB, error) {
	path := DB_FOLDER + "/" + "issues.db"

	// check existence
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("database not initialized")
	}

	return sql.Open("sqlite3", path)
}

// Loads and executes the schema
func InitSchema(db *sql.DB) error {
	// TODO: Change to not read the schema every time
	schema, err := os.ReadFile(SCHEMA_FILE)
	panic_mode(err)

	_, err = db.Exec(string(schema))

	fmt.Printf("Schema: %s initiated.\n", SCHEMA_FILE)
	return err
}

// Should initiate the DB only if it doesnt exist
// Check if the issue.db exists, if not, load and execute the schema
// via the InitSchema function
func InitDB() (*sql.DB, error) {
	path := DB_FOLDER + "/" + DB_NAME + ".db"

	err := os.MkdirAll(DB_FOLDER, 0755)
	if err != nil {
		return nil, err
	}

	// todo: add config for swapping out the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = InitSchema(db)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database initialized")
	return db, nil
}

func DBExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
