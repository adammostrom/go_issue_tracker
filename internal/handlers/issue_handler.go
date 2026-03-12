package handlers

import (
	"fmt"
	"issuetracker/internal/models"
	"time"
)

// kind of like an interface
type IssueEndpoint struct{}

// TODO: Think about how to handle duplicate issues, if every issue gets a unique number, there will never be duplicates

func (s IssueEndpoint) CreateNewIssue(name string, desc string) models.Issue {

	timestamp := time.Now().Format("2006-01-02-15:04")
	var log = []models.LogEntry{
		{Timestamp: timestamp, Entry: "Issue created"},
	}

	// Internal ID generated at db insert
	issue := models.Issue{
		Name:        name,
		Description: desc,
		Log:         log,
		Resolved:    false,
	}

	fmt.Printf("Issue: %s - Created at: %v\n", name, timestamp)
	return issue
}
