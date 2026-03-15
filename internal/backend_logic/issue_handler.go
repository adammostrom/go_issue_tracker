package handlers

/*
Handles HTTP requests: parses input, writes responses.

Calls business logic functions.

Should not know SQL.
*/
import (
	"fmt"
	"issuetracker/internal/models"
	"time"
)

// kind of like an interface
type IssueEndpoint struct{}

// TODO: Think about how to handle duplicate issues, if every issue gets a unique number, there will never be duplicates

func (s IssueEndpoint) CreateNewIssue(title string, desc string, id int64) models.Issue {

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
		Active:      false,
	}

	fmt.Printf("Issue: %s - Created at: %v\n", title, timestamp)
	return issue
}
