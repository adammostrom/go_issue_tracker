package db

// THIS PACKAGE (FILE) ONLY TALKS TO THE DATABASE, NOTHING ELSE, NO LOGIC
// That’s all it should do. No HTTP logic, no request parsing,
// Only knows SQL and persistence.
// Exposes methods like AddIssue(issue) or GetAllIssues().

import (
	"database/sql"
	"issuetracker/internal/models"
)

type IssueDBConn struct {
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
func NewIssueConn(db *sql.DB) *IssueDBConn {
	return &IssueDBConn{db: db} // exports the reference to the created instance
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
func (s *IssueDBConn) GetAllIssues() ([]models.Issue, error) {
	rows, err := s.db.Query("SELECT * FROM Issues") // Come back to update based on sql db
	if err != nil {
		return nil, err
	}
	// 	rows.Close() -> frees DB resources, it is required or connections will leak

	/*
		DEFER
		Its good because:
			- You declare cleanup next to acquisition
			- No try/finally
			- No forgotten cleanup in long functions

	*/
	defer rows.Close()

	var issues []models.Issue

	for rows.Next() {
		var i models.Issue
		rows.Scan(&i.Internal_id, &i.Title, &i.Description, &i.Log, &i.Active) // Scans the database (sql tables) for the values. Internal ID should only be between the database and the backend logic, not the user.
		issues = append(issues, i)
	}
	return issues, nil
}

// Add an issue
// Should all fields be required? Or just the name of the issue?
func (s *IssueDBConn) AddIssue(issue models.Issue) error {
	//TODO:  Generate internal ID

	_, err := s.db.Exec("INSERT INTO Issues(internal_id, name, description, log, resolved) VALUES ($1, $2, $3)", issue.Title, issue.Description, issue.Log, issue.Active)
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
func (s *IssueDBConn) GetIssueByID(id int) (models.Issue, error) {

	var issue models.Issue

	err := s.db.QueryRow(
		"SELECT * FROM Issues WHERE internal_id = $1", id).Scan(&issue.Internal_id, &issue.Title, &issue.Description, &issue.Log, &issue.Active)

	if err != nil {
		return models.Issue{}, err // return empty issue, interpret it higher up (empty issue = no issue found)
	}

	return issue, nil
}
