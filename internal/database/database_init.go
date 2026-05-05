package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema_issues.sql
var schemaFS embed.FS

// const schemaFile = "internal/database/schema_issues.sql"
const dbFileName = "issuedb.db"

func Open() (*sql.DB, error) {
	dbPath, firstRun, err := resolveDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if firstRun {
		if err := initSchema(db); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

func resolveDBPath() (path string, firstRun bool, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", false, err
	}

	exePath, _ = filepath.EvalSymlinks(exePath)
	baseDir := filepath.Dir(exePath)

	dataDir := filepath.Join(baseDir, ".issuetracker")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", false, err
	}

	dbPath := filepath.Join(dataDir, dbFileName)

	_, err = os.Stat(dbPath)
	firstRun = os.IsNotExist(err)

	return dbPath, firstRun, nil
}

func initSchema(db *sql.DB) error {
	schema, err := schemaFS.ReadFile("schema_issues.sql")
	if err != nil {
		return fmt.Errorf("read embedded schema: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("applying schema: %w", err)
	}

	return nil
}
