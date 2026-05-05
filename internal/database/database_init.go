package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const SCHEMA_FILE = "internal/database/schema_issues.sql"
const DB_FILE_NAME = "issuetracker_sqlite3.db"

// SQLITE

func panic_mode(err error) {
	if err != nil {
		panic(err)
	}
}

// Basically Run the DB
func OpenDB() (*sql.DB, error) {
	dbDir := getBaseDirectoryOfExecutable()
	if dbDir == "" {
		log.Fatal("Could not retrieve base directory")
	}

	dbPath := filepath.Join(dbDir, DB_FILE_NAME)

	// check existence
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database not initialized")
	}

	return sql.Open("sqlite3", dbPath)
}

//go:embed schema_issues.sql
var schemaFS embed.FS

// Loads and executes the schema
func InitSchema(db *sql.DB) error {
	// Use schemaFS to embed the schema into the binary
	schema, err := schemaFS.ReadFile(SCHEMA_FILE)
	panic_mode(err)

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Schema: %s initiated.\n", SCHEMA_FILE)
	return err
}

// Should initiate the DB only if it doesnt exist
// Check if the issue.db exists, if not, load and execute the schema
// via the InitSchema function
func InitDB() (*sql.DB, error) {

	dbDir := getBaseDirectoryOfExecutable()
	if dbDir == "" {
		log.Fatal("Could not retrieve base directory")
	}

	if err := os.Mkdir(dbDir, 0755); err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(dbDir, "issuetracker_sqlite3.db")

	// todo: add config for swapping out the database
	db, err := sql.Open("sqlite3", dbPath)
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

func getBaseDirectoryOfExecutable() string {
	// Get the base directory of the executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exePath, _ = filepath.EvalSymlinks(exePath)
	baseDir := filepath.Dir(exePath)

	// Create the data directory
	dataDir := filepath.Join(baseDir, ".issuetracker")
	return dataDir
}
