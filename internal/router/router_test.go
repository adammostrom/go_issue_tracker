package router

import (
	"encoding/json"
	"errors"
	"issuetracker/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeIssueService struct {
	called bool
	status models.IssueStatus
	issues []models.Issue
	issue  *models.Issue
	id_arg int
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

/*
Before writing test, ask:

What does this handler expect as input?
What does it call?
What does it return?
*/

// These simulate the service that the router is calling on. So the router.go handlers call these functions here.
func (f *fakeIssueService) CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error) {
	return nil, nil
}

func (f *fakeIssueService) GetAllIssues(status models.IssueStatus) ([]models.Issue, error) {
	f.called = true
	f.status = status
	return f.issues, f.err
}

func (f *fakeIssueService) GetIssueByID(id int) (*models.Issue, error) {
	f.called = true

	if f.err != nil {
		return nil, f.err
	}

	if f.issue == nil || f.issue.Internal_ID != int64(id) {
		return nil, nil
	}
	return f.issue, f.err
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

// INVALID STATUS:

func TestGetIssuesHandler_InvalidStatus(t *testing.T) {
	service := &fakeIssueService{}
	router := &Router{issueService: service}

	req := httptest.NewRequest(http.MethodGet, "/issues?status=INVALID", nil)
	w := httptest.NewRecorder()

	router.getIssuesHandler(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}

	if service.called {
		t.Fatal("service should NOT be called on invalid input")
	}
}

// SERVICE FAILURE, mocking a failed db initiation, normal issue request
/*
IF service returns error
THEN handler should return HTTP 400
AND stop execution

*/
func TestGetIssuesHandler_ServiceError(t *testing.T) {
	service := &fakeIssueService{
		err: errors.New("db fail"),
	}

	router := &Router{issueService: service}

	req := httptest.NewRequest(http.MethodGet, "/issues", nil)
	w := httptest.NewRecorder()

	router.getIssuesHandler(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// TODO: MORE TESTS:
// - look at your handler and identify branches:
// Every decision point = test case.

func setupGetSingleIssue(path string, service *fakeIssueService) *http.Response {

	// Fake service, (dependency injection), the issue is here placed in the service so we can "GET" it.

	service.issue = &models.Issue{

		Internal_ID:  1,
		External_Ref: "extref",
		Title:        "title",
		Description:  "description",
		Active:       true,
	}

	router := &Router{issueService: service}

	req := httptest.NewRequest(http.MethodGet, path, nil)

	w := httptest.NewRecorder()

	router.getSingleIssueHandler(w, req)

	res := w.Result()

	return res
}

// HAPPY PATH

func TestGetIssueHandler_OK(t *testing.T) {

	service := &fakeIssueService{}

	res := setupGetSingleIssue("/issues/1", service)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d\n", res.StatusCode)
	}

	// Flow correctness
	if !service.called {
		t.Fatalf("service not called\n")
	}

	var body models.IssueResponse
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response\n")
	}

	if body.External_Ref != "extref" {
		t.Fatalf("Wrong external ref or external ref is empty\n")
	}
}

// Simulate wrong id, should return 400
func TestGetSingleIssue_INVALID_ID(t *testing.T) {
	service := &fakeIssueService{}

	res := setupGetSingleIssue("/issues/a", service)

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	if service.called {
		t.Fatal("service should NOT be called on invalid input")
	}

}

func TestGetSingleIssue_ISSUE_NOT_EXIST(t *testing.T) {
	service := &fakeIssueService{}

	res := setupGetSingleIssue("/issues/999", service)

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	if !service.called {
		t.Fatal("service SHOULD be called for issue not existing")
	}
}
