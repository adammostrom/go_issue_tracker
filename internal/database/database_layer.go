package database

// THIS PACKAGE (FILE) ONLY TALKS TO THE DATABASE, NOTHING ELSE, NO LOGIC
// That’s all it should do. No HTTP logic, no request parsing,
// Only knows SQL and persistence.
// Exposes methods like AddIssue(issue) or GetAllIssues().

import (
	"database/sql"
	"fmt"
	"issuetracker/internal/models"
	"log"
)

type DatabaseConnection struct {
	db *sql.DB
}

/*
Constructor function

This is the standard Go constructor pattern.
Go doesn’t have classes. This is how you “build” objects.

- Capital N → exported (visible outside the package) (in IssueDBConn)
- Return type: pointer to IssueDBConn
*/

// Its job is only one thing: talk to SQL.
func NewDatabaseConnection(db *sql.DB) *DatabaseConnection {
	return &DatabaseConnection{db: db} // exports the reference to the created instance
}

// the db argument, which is a pointer to sql.DB, gets assigned to the IssueDBConn struct db.
// meaning that the sql.DB pointer is reachable via the IssueDBConn db field.
// & = address, returns pointer of type IssueDBConn

/*
(s *DeviceStore) = method reciever -> method belongs to DeviceStore.
"s" is like "this" or "self".
its a pointer so it operates on the real store

method GetAll (capital G = exported)
returns []models.Device -> slice of Device structs
*/

// This struct wraps the database connection so that methods like this can exist:
// “Given an Issue struct, store it in the database.”
func (s *DatabaseConnection) QueryAllIssues() ([]models.Issue, error) {
	rows, err := s.db.Query("SELECT * FROM Issues")
	if err != nil {
		panic_mode(err)
		return nil, err
	}

	defer rows.Close()

	var issues []models.Issue

	for rows.Next() {
		var i models.Issue
		rows.Scan(&i.Internal_id, &i.Title, &i.Description, &i.Active) // Skip log for now
		issues = append(issues, i)
		fmt.Printf("Issue found: %d - %s - %s \n", i.Internal_id, i.Title, i.Active)
	}
	return issues, nil
}

// Add an issue
// Should all fields be required? Or just the name of the issue?
func (s *DatabaseConnection) AddIssue(issue models.Issue) error {

	stmt := `INSERT INTO Issues(title, description, active) VALUES (?, ?, ?)`
	res, err := s.db.Exec(stmt, issue.Title, issue.Description, 1)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// MOSTLY FOR DEBUGGING
	fmt.Printf("Inserted value - id: %d, title: %s, description: %s, active %s \n", id, issue.Title, issue.Description, issue.Active)
	return err
}

// Returns device based on serial number
/*
QueryRow:

does not return rows iterator

does not need Close()

executes immediately

Scan triggers the query

It’s a tiny, elegant shortcut for “exactly one row expected.”


*/
func (s *DatabaseConnection) GetIssueByID(id int) (models.Issue, error) {

	var issue models.Issue

	err := s.db.QueryRow(
		"SELECT * FROM Issues WHERE id = $1", id).Scan(&issue.Internal_id, &issue.Title, &issue.Description, &issue.Active)

	if err != nil {
		return models.Issue{}, err // return empty issue, interpret it higher up (empty issue = no issue found)
	}

	return issue, nil
}

func (s *DatabaseConnection) UpdateIssue(fields []interface{}, query string, id int) error {

	res, err := s.db.Exec(query, fields...)
	if err != nil {
		log.Fatal(err)
	}
	// MOSTLY FOR DEBUGGING
	fmt.Printf("res: %v\n", res)
	return err
}
