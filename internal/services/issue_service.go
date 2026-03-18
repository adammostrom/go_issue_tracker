package services

/*
Handles HTTP requests: parses input, writes responses. SHOULD NOT KNOW HTTP

Calls business logic functions.

Should not know SQL.

The service only deals with Go data structures and business rules.
*/
import (
	"fmt"
	"issuetracker/internal/database"
	"issuetracker/internal/models"
	"log"
	"time"
)

// IssueService accepts a database connection in order to delegate tasks downwards
type IssueService struct {
	db_layer *database.DatabaseConnection
}

// Constructor that accepts an issue service and returns a router pointer.
func NewIssueService(db_layer *database.DatabaseConnection) *IssueService {
	return &IssueService{
		db_layer: db_layer,
	}
}

// TODO: Think about how to handle duplicate issues, if every issue gets a unique number, there will never be duplicates

func (s *IssueService) CreateNewIssue(id int64, title string, desc string) models.Issue {

	timestamp := time.Now().Format("2006-01-02-15:04")
	var log = []models.LogEntry{
		{Timestamp: timestamp, Entry: "Issue created"},
	}

	// Internal ID generated at db insert
	issue := models.Issue{
		Internal_id: id,
		Title:       title,
		Description: desc,
		Log:         log,
		Active:      true,
	}

	s.db_layer.AddIssue(issue)
	fmt.Printf("Issue: %s - Created at: %v\n", title, timestamp)
	return issue
}

// Update to integrate the database instead of a slice of issues
func (s *IssueService) GetAllIssues() []models.Issue {
	issues, err := s.db_layer.QueryAllIssues()
	if err != nil {
		log.Fatal(err)
	}
	return issues
}

func (s *IssueService) GetSingleIssue(id int) models.Issue {
	issue, err := s.db_layer.GetIssueByID(id)
	if err != nil {
		log.Fatal(err)
	}
	return issue
}
