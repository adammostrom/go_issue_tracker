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
	"strings"
	"time"
)

type DatabaseInterface interface {
	GetIssue(id int) (*models.Issue, error)
	GetIssues() ([]models.Issue, error)
	ModifyIssue(fields []interface{}, query string, id int) error
	CreateIssue(issue *models.Issue) (*models.Issue, error)
	CreateLogEntry(id int64, logEntry models.LogEntry) error
	DeleteIssue(id int) error
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

func (s *IssueService) CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error) {

	timestamp := time.Now().Format("2006-01-02-15:04")
	var logEntries = []models.LogEntry{
		{Timestamp: timestamp, Entry: "Issue created"},
	}

	// Internal ID generated at db insert
	issue := &models.Issue{
		External_Ref: req.External_Ref,
		Title:        req.Title,
		Description:  req.Description,
		Log:          logEntries,
		Active:       true,
	}
	if err := issue.ValidateIssue(); err != nil {
		return nil, err
	}

	issue_updated, err := s.db_layer.CreateIssue(issue)
	if err != nil {
		return nil, err
	}
	// FOR TESTING
	fmt.Printf("Issue: %s - Created at: %v\n", req.Title, timestamp)
	return issue_updated, nil
}

func (s *IssueService) GetAllIssues() ([]models.Issue, error) {
	issues, err := s.db_layer.GetIssues()
	if err != nil {
		return nil, err // Works because a slice is a pointer to an array
	}
	return issues, nil
}

func (s *IssueService) GetIssueByID(id int) (*models.Issue, error) {
	issue, err := s.db_layer.GetIssue(id)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *IssueService) DeleteIssue(id int) error {
	err := s.db_layer.DeleteIssue(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *IssueService) PatchIssue(id int, upd_req models.UpdateIssueRequest) error {

	query := "UPDATE issues SET "
	updated_fields := []interface{}{}
	i := 1

	// Log entry
	timestamp := time.Now().Format("2006-01-02-15:04")
	var logEntries = []models.LogEntry{}

	if upd_req.Title != nil {
		if err := models.ValidateTitle(*upd_req.Title); err != nil {
			return err
		}
		updated_fields = append(updated_fields, *upd_req.Title)
		query += fmt.Sprintf("title=$%d,", i)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "Title changed to: " + *upd_req.Title})
		i++
	}
	if upd_req.External_Ref != nil {
		if err := models.ValidateExternalRef(*upd_req.External_Ref); err != nil {
			return err
		}
		updated_fields = append(updated_fields, *upd_req.External_Ref)
		query += fmt.Sprintf("external_ref=$%d,", i)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "External Reference changed to: " + *upd_req.External_Ref})
		i++
	}
	if upd_req.Description != nil {
		if err := models.ValidateDescription(*upd_req.Description); err != nil {
			return err
		}
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

	err := s.db_layer.ModifyIssue(updated_fields, query, id)

	if err != nil {
		return err
	}

	for _, entry := range logEntries {
		err := s.db_layer.CreateLogEntry(int64(id), entry)
		if err != nil {
			fmt.Printf(err.Error()) // Maybe come back to this... /2026-03-31
		}
	}

	return nil
}

func (s *IssueService) GetLogsFromIssue(id int) (*[]models.LogEntry, error) {
	logs, err := s.db_layer.GetLogs(id)
	if err != nil {
		return nil, err
	}
	return logs

}

// TODO, IMPLEMENT:
