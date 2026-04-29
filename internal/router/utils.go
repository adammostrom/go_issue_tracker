package router

import (
	"fmt"
	"issuetracker/internal/models"
	"strconv"
	"strings"
)

func parseIDfromPath(path string) (int64, error) {

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[2], 10, 64)

	if err != nil {
		return id, err
	}

	if id < 0 {
		err := fmt.Errorf("Negative ID not allowed")
		return id, err
	}

	return id, nil
}

// Helpers

func issueToIssueResponse(issue models.Issue) models.IssueResponse {

	resp := models.IssueResponse{
		Internal_ID:  issue.Internal_ID,
		External_Ref: deref(issue.External_Ref, "null"),
		Title:        issue.Title,
		Description:  issue.Description,
		Active:       issue.Active,
	}
	return resp
}

// Handling null values for external ref.
func deref(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}
