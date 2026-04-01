package database

// THIS PACKAGE (FILE) ONLY TALKS TO THE DATABASE, NOTHING ELSE, NO LOGIC
// That’s all it should do. No HTTP logic, no request parsing,
// Only knows SQL and persistence.
// Exposes methods like AddIssue(issue) or GetAllIssues().

import (
	"database/sql"
	"errors"
	"fmt"
	"issuetracker/internal/models"
)

type DatabaseConnection struct {
	db *sql.DB
}

/*
Constructor function

This is the standard Go constructor pattern.
Go doesn’t have classes. This is how you “build” objects.

- Return type: pointer to IssueDBConn
*/

// Its job is only one thing: talk to SQL.
func NewDatabaseConnection(db *sql.DB) *DatabaseConnection {
	return &DatabaseConnection{db: db} // exports the reference to the created instance
}

// the db argument, which is a pointer to sql.DB, gets assigned to the IssueDBConn struct db.
// meaning that the sql.DB pointer is reachable via the IssueDBConn db field.
// & = address, returns pointer of type IssueDBConn

// “Given an Issue struct, store it in the database.”
func (s *DatabaseConnection) GetIssues() ([]models.Issue, error) {
	rows, err := s.db.Query("SELECT * FROM Issues") // TODO 2026-03-31: Make a view in SQL and select from that one instead
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var issues []models.Issue

	for rows.Next() {
		var i models.Issue
		rows.Scan(&i.Internal_ID, &i.External_Ref, &i.Title, &i.Description, &i.Active) // Skip log for now
		issues = append(issues, i)
		fmt.Printf("i: %v\n", i)
	}
	return issues, nil
}

func (s *DatabaseConnection) GetLogs(id int) ([]models.LogEntry, error) {
	rows, err := s.db.Query("SELECT timestamp, entry FROM Logs WHERE issue_id = ?", id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var logs []models.LogEntry
	for rows.Next() {
		var i models.LogEntry

		/*
			rows.Scan requires pointers to the destination fields.
			Since i.Entry is a string (not *string), Go is just passing a copy,
			so the scan doesn’t actually populate the field—hence your empty string in the logs.

		*/
		rows.Scan(&i.Timestamp, &i.Entry)

		fmt.Printf("FROM DB, timestamp: %s, entry: %s\n", i.Timestamp, i.Entry)
		logs = append(logs, i)
	}
	return logs, nil
}

// Add an issue
// Should all fields be required? Or just the name of the issue?
// 2026-03-24: Returns the pointer only so the functions are the same as in service and can be mocked
func (s *DatabaseConnection) CreateIssue(issue *models.Issue) (*models.Issue, error) {

	stmt := `INSERT INTO Issues(title, external_ref, description, active) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(stmt, issue.Title, issue.External_Ref, issue.Description, 1)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	log_err := s.CreateLogEntry(id, issue.Log[0])
	fmt.Printf("LOG ENTRY From CREATE ISSUE: %s \n", issue.Log[0])
	if log_err != nil {
		return nil, log_err
	}

	// Update the issue field "Internal ID" with the returned id from the database
	issue.Internal_ID = id

	return issue, err
}

func (s *DatabaseConnection) CreateLogEntry(id int64, logEntry models.LogEntry) error {
	// Insert into logs
	log_stmt := `INSERT INTO Logs(issue_id, timestamp, entry) VALUES (?, ?, ?)`
	_, err := s.db.Exec(log_stmt, id, logEntry.Timestamp, logEntry.Entry)
	if err != nil {
		return err
	}

	fmt.Printf("CREATELOGENTRY: appended log: timestamp: %s, entry: %s", logEntry.Timestamp, logEntry.Entry)
	return err

}

func (s *DatabaseConnection) GetIssue(id int) (*models.Issue, error) {

	// For pointer referencing, initiate the struct first, otherwise pointer is nil
	issue := &models.Issue{}

	// TODO: Dont return everything (*), create a VIEW in SQL and return from that instead
	err := s.db.QueryRow(
		"SELECT * FROM Issues WHERE id = $1", id).Scan(&issue.Internal_ID, &issue.External_Ref, &issue.Title, &issue.Description, &issue.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIssueNotFound(id)
		}
		return nil, err
	}
	return issue, nil
}

func (s *DatabaseConnection) ModifyIssue(fields []interface{}, query string, id int) error {

	res, err := s.db.Exec(query, fields...)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrIssueNotFound(id)
	}
	// MOSTLY FOR DEBUGGING
	fmt.Printf("res: %v\n", res)
	return err
}

func (s *DatabaseConnection) DeleteIssue(id int) error {

	stmt := `DELETE FROM Issues WHERE id = $1;`

	res, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrIssueNotFound(id)
	}

	err = s.DeleteLogs(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseConnection) DeleteLogs(id int) error {
	stmt := `DELETE FROM Logs WHERE issue_id = $1`

	res, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrIssueNotFound(id)
	}
	fmt.Printf("Logs for issue with id: %d deleted successfully.\n", id)
	return nil
}

// TODO: implement

func (s *DatabaseConnection) ExtRefExists(ref string) (bool, error) {

	var exists bool

	err := s.db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM Issues WHERE external_ref = ?
		)
		`, ref).Scan(&exists) // <--- store the scanned value into this variable. SELECT EXISTS return true/false, SELECT 1 = "we dont care about data, only if it exists or not"

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	return exists, err
}

// Errors

type ErrorIssueNotFound struct {
	msg string
}

func (ierr ErrorIssueNotFound) Error() string {
	return ierr.msg
}

func ErrIssueNotFound(id int) error {
	return fmt.Errorf("issue with id %d not found", id)
}
