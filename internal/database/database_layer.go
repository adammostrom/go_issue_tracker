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
	"log"
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
	if log_err != nil {
		return nil, log_err
	}

	// Update the issue field "Internal ID" with the returned id from the database
	issue.Internal_ID = id

	// MOSTLY FOR DEBUGGING
	fmt.Printf("Inserted value - id: %d, title: %s, description: %s, active %s \n", id, issue.Title, issue.Description, fmt.Sprintf("%t", issue.Active))

	return issue, err
}

func (s *DatabaseConnection) CreateLogEntry(id int64, logEntry models.LogEntry) error {
	// Insert into logs
	log_stmt := `INSERT INTO Logs(issue_id, timestamp, entry) VALUES (?, ?, ?)`
	res, err := s.db.Exec(log_stmt, id, logEntry.Timestamp, logEntry.Entry)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.RowsAffected())

	fmt.Printf("Appended log entry: %s", logEntry.Entry)
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

	return nil
}

// TODO: implement

func (s *DatabaseConnection) ExtRefExists(ref string) bool {
	err := s.db.QueryRow(
		"SELECT * FROM Logs WHERE id = $1", id).Scan(&issue.Internal_ID, &issue.External_Ref, &issue.Title, &issue.Description, &issue.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIssueNotFound(id)
		}
		return nil, err
	}
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
