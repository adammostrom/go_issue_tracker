package services

/*
Handles HTTP requests: parses input, writes responses. SHOULD NOT KNOW HTTP

Calls business logic functions.

Should not know SQL.

The service only deals with Go data structures and business rules.
*/
import (
	"fmt"
	"issuetracker/internal/models"
	"log"
	"strings"
	"time"
)

type DatabaseInterface interface {
	GetIssue(id int) (*models.Issue, error)
	GetIssues() ([]models.Issue, error)
	ModifyIssue(fields []interface{}, query string, id int) error
	CreateIssue(issue *models.Issue) (*models.Issue, error)
	CreateLogEntry(id int64, logEntry models.LogEntry) error
}

// IssueService accepts a database connection in order to delegate tasks downwards
// 2026-03-24: Depends on interface instead of direct connection
type IssueService struct {
	db_layer DatabaseInterface
}

// Constructor that accepts an issue service and returns a router pointer.
func NewIssueService(db_layer DatabaseInterface) *IssueService {
	return &IssueService{
		db_layer: db_layer,
	}
}

// TODO: Think about how to handle duplicate issues, if every issue gets a unique number, there will never be duplicates

func (s *IssueService) CreateNewIssue(external_ref string, title string, desc string) (models.Issue, error) {

	timestamp := time.Now().Format("2006-01-02-15:04")
	var logEntries = []models.LogEntry{
		{Timestamp: timestamp, Entry: "Issue created"},
	}

	// Internal ID generated at db insert
	issue := models.Issue{
		External_Ref: external_ref,
		Title:        title,
		Description:  desc,
		Log:          logEntries,
		Active:       true,
	}

	issue_updated, err := s.db_layer.CreateIssue(&issue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Issue: %s - Created at: %v\n", title, timestamp)
	return *issue_updated, nil
}

// Update to integrate the database instead of a slice of issues
func (s *IssueService) GetAllIssues() []models.Issue {
	issues, err := s.db_layer.GetIssues()
	if err != nil {
		log.Fatal(err)
	}
	return issues
}

func (s *IssueService) GetIssueByID(id int) (*models.Issue, error) {
	issue, err := s.db_layer.GetIssue(id)

	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *IssueService) PatchIssue(id int, upd_req models.UpdateIssueRequest) error {

	query := "UPDATE issues SET "
	updated_fields := []interface{}{}
	i := 1

	// Log entry
	timestamp := time.Now().Format("2006-01-02-15:04")
	var logEntries = []models.LogEntry{}

	if upd_req.Title != nil {
		updated_fields = append(updated_fields, *upd_req.Title)
		query += fmt.Sprintf("title=$%d,", i)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "Title changed to: " + *upd_req.Title})
		i++
	}
	if upd_req.External_Ref != nil {
		updated_fields = append(updated_fields, *upd_req.External_Ref)
		query += fmt.Sprintf("external_ref=$%d,", i)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "External Reference changed to: " + *upd_req.External_Ref})
		i++
	}
	if upd_req.Description != nil {
		updated_fields = append(updated_fields, *upd_req.Description)
		query += fmt.Sprintf("description=$%d,", i)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "Description changed to: " + *upd_req.Description})
		i++
	}
	if upd_req.Active != nil {
		updated_fields = append(updated_fields, *upd_req.Active)
		query += fmt.Sprintf("active=$%d,", i)
		str := fmt.Sprintf("%t", *upd_req.Active)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "Active status changed to: " + str})
		i++
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id=%d", id)

	fmt.Printf("Sent %s to db with fields: \n %s\n", query, updated_fields)

	err := s.db_layer.ModifyIssue(updated_fields, query, id)

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range logEntries {
		err := s.db_layer.CreateLogEntry(int64(id), entry)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
