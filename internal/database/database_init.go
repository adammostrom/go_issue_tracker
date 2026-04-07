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

func NewDB(path string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", path)
	panic_mode(err)

	err = db.Ping()
	if err != nil {
		fmt.Printf("Database not reached with error: %s\n", err)
		panic(err)
	}
	return db, nil
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
func InitDB() *sql.DB {

	path := DB_FOLDER + "/issues.db"

	os.MkdirAll(DB_FOLDER, 0755)

	firstInit := !DBExists(path)

	db, err := NewDB(path)
	panic_mode(err)

	// Now schema runs only once
	if firstInit {
		err = InitSchema(db)
		panic_mode(err)
		fmt.Println("Database initiated successfully")
		fmt.Printf("Created DB instance with name: %s", DB_NAME+"db\n")

	}

	return db
}

func DBExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
