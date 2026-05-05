package services

import (
	"fmt"
	"issuetracker/internal/models"
	"testing"
)

type MockDB struct {
	GetIssuesFn    func(filter models.IssueFilter) ([]models.Issue, error)
	CreateIssueFn  func(issue *models.Issue) (*models.Issue, error)
	ExtRefExistsFn func(ref *string) (bool, error)
}

func (m *MockDB) GetIssues(filter models.IssueFilter) ([]models.Issue, error) {
	return m.GetIssuesFn(filter)
}

func (m *MockDB) CreateIssue(issue *models.Issue) (*models.Issue, error) {
	return m.CreateIssueFn(issue)
}

func (m *MockDB) ExtRefExists(ref *string) (bool, error) {
	return m.ExtRefExistsFn(ref)
}

// stub unused methods (required to satisfy interface)
func (m *MockDB) GetIssue(id int) (*models.Issue, error) { return nil, nil }
func (m *MockDB) ModifyIssue(fields []interface{}, query string, id int) error {
	return nil
}
func (m *MockDB) CreateLogEntry(id int64, logEntry models.LogEntry) error {
	return nil
}
func (m *MockDB) DeleteLogs(id int) error  { return nil }
func (m *MockDB) DeleteIssue(id int) error { return nil }
func (m *MockDB) GetLogs(id int) ([]models.LogEntry, error) {
	return nil, nil
}

func TestCreateNewIssue_Success(t *testing.T) {
	mock := &MockDB{
		ExtRefExistsFn: func(ref *string) (bool, error) {
			return false, nil
		},
		CreateIssueFn: func(issue *models.Issue) (*models.Issue, error) {
			issue.Internal_ID = 1
			return issue, nil
		},
	}

	service := NewIssueService(mock)

	req := models.CreateIssueRequest{
		Title: "TEST_SUCCESS",
	}

	issue, err := service.CreateNewIssue(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if issue.Internal_ID != 1 {
		t.Fatalf("expected ID 1, got %d", issue.Internal_ID)
	}

	if issue.Progress != models.Idle {
		t.Fatalf("expected Idle progress")
	}

	fmt.Print("No errors\n")

}

func TestCreateNewIssue_ExtRefExists(t *testing.T) {
	ref := "ABC"

	mock := &MockDB{
		ExtRefExistsFn: func(ref *string) (bool, error) {
			return true, nil
		},
	}

	service := NewIssueService(mock)

	req := models.CreateIssueRequest{
		Title:        "test",
		External_Ref: &ref,
	}

	_, err := service.CreateNewIssue(req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
