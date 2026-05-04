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
	GetIssues(filter models.IssueFilter) ([]models.Issue, error)
	ModifyIssue(fields []interface{}, query string, id int) error
	CreateIssue(issue *models.Issue) (*models.Issue, error)
	CreateLogEntry(id int64, logEntry models.LogEntry) error
	DeleteLogs(id int) error
	DeleteIssue(id int) error
	ExtRefExists(ref *string) (bool, error)
	GetLogs(id int) ([]models.LogEntry, error)
	GetIssueByRef(reference *string) (*models.Issue, error)
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

func (s *IssueService) CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error) {

	timestamp := time.Now().Format("2006-01-02 15:04") // 2026-05-04: Removed HH:MM, no need for that resolution
	var logEntries = []models.LogEntry{
		{Timestamp: timestamp, Entry: "Issue created"},
	}

	issue := &models.Issue{
		Title:        req.Title,
		External_Ref: req.External_Ref,
		Description:  req.Description,
		Log:          logEntries,
		Active:       true,
		Progress:     models.Idle,
	}

	if err := issue.ValidateIssue(); err != nil {
		return nil, err
	}
	exists, err := s.db_layer.ExtRefExists(issue.External_Ref)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("External ref: %s already exists.", *issue.External_Ref)
	}

	return s.db_layer.CreateIssue(issue)
}

// TODO: 2026-04-05 Come back to this and refactor/update so  its more centralized and not a new function for each filtering query
func (s *IssueService) GetAllIssues(filter models.IssueFilter) ([]models.Issue, error) {

	issues, err := s.db_layer.GetIssues(filter)
	if err != nil {
		return nil, err
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

	// Check if issue by this ID exists
	_, err := s.GetIssueByID(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

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
		if err := models.ValidateExternalRefWrapper(*upd_req.External_Ref); err != nil {
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
	if upd_req.Progress != nil {
		updated_fields = append(updated_fields, *upd_req.Progress)
		query += fmt.Sprintf("progress=$%d", i)
		logEntries = append(logEntries, models.LogEntry{Timestamp: timestamp, Entry: "Progress changed to: " + upd_req.Progress.String()})
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id=%d", id)

	err = s.db_layer.ModifyIssue(updated_fields, query, id)

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

func (s *IssueService) GetLogsFromIssue(id int) ([]models.LogEntry, error) {

	// Check if issue by this ID exists
	_, err := s.GetIssueByID(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	logs, err := s.db_layer.GetLogs(id)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// TODO: Change to CreateLogEntry
func (s *IssueService) AddLogEntry(id int, entry string) error {

	_, err := s.GetIssueByID(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	timestamp := time.Now().Format("2006-01-02 15:04") // 2026-05-04: Removed HH:MM, no need for that resolution
	var logEntry = models.LogEntry{
		Timestamp: timestamp, Entry: entry,
	}
	err = logEntry.ValidateEntry()
	if err != nil {
		return err
	}

	return s.db_layer.CreateLogEntry(int64(id), logEntry)

}

func (s *IssueService) DeleteLogsFromIssue(id int) error {
	_, err := s.GetIssueByID(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return s.db_layer.DeleteLogs(id)
}

func (s *IssueService) GetIssueByRef(reference *string) (*models.Issue, error) {
	issue, err := s.db_layer.GetIssueByRef(reference)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
