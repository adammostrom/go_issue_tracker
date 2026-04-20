package router

import (
	"encoding/json"
	"fmt"
	"issuetracker/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeIssueService struct {
	called bool
	status models.IssueStatus
	issues []models.Issue
	err    error
}

/*
	type IssueServiceInterface interface {
		CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error)
		GetAllIssues(status models.IssueStatus) ([]models.Issue, error)
		GetIssueByID(id int) (*models.Issue, error)
		DeleteIssue(id int) error
		PatchIssue(id int, upd_req models.UpdateIssueRequest) error
		GetLogsFromIssue(id int) ([]models.LogEntry, error)
		AddLogEntry(id int, entry string) error
	}
*/
func (f *fakeIssueService) CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error) {
	return nil, nil
}

func (f *fakeIssueService) GetAllIssues(status models.IssueStatus) ([]models.Issue, error) {
	f.called = true
	f.status = status
	fmt.Println("fakeIssueService: GetAllIssues CALLED OK")
	return f.issues, f.err
}

func (f *fakeIssueService) GetIssueByID(id int) (*models.Issue, error) {
	return nil, nil
}

func (f *fakeIssueService) DeleteIssue(id int) error {
	return nil
}

func (f *fakeIssueService) PatchIssue(id int, upd_req models.UpdateIssueRequest) error {
	return nil
}

func (f *fakeIssueService) GetLogsFromIssue(id int) ([]models.LogEntry, error) {
	return nil, nil
}

func (f *fakeIssueService) AddLogEntry(id int, entry string) error {
	return nil
}

// HAPPY PATH
// run with "go test ./internal/router/ -v"
// go test looks for files named *_test.go"
func TestGetIssuesHandler_OK(t *testing.T) {

	// Fake service, (dependency injection)
	service := &fakeIssueService{
		issues: []models.Issue{
			{Internal_ID: 1, Title: "testGetAll"},
		},
	}

	// Inject fake service into router
	router := &Router{issueService: service}

	// GET /issues
	req := httptest.NewRequest(http.MethodGet, "/issues", nil)

	// Instead of sending data over network, stores status code and response body
	w := httptest.NewRecorder()

	// Same as the HTTP request
	router.getIssuesHandler(w, req)

	// Get the status code and response body
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d\n", res.StatusCode)
	}
	// Flow correctness
	if !service.called {
		t.Fatalf("service not called\n")
	}
	if service.status != models.StatusDefault {
		t.Fatalf("expected default status, got %v\n", service.status)
	}

	var body []models.IssueResponse
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response\n")
	}
	if len(body) != 1 {
		t.Fatalf("expected 1 issue, got %d\n", len(body))
	}
}
