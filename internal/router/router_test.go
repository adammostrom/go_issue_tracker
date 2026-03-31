package router

import (
	"issuetracker/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
func (s *Router) getSingleIssueHandler(w http.ResponseWriter, r *http.Request) {

	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	issue, err := s.issueService.GetIssueByID(int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := issueToIssueResponse(*issue)

	json.NewEncoder(w).Encode(resp)

	w.WriteHeader(http.StatusOK)

}
*/

/*
type IssueServiceInterface interface {
	CreateNewIssue(external_ref string, title string, desc string) (models.Issue, error)
	GetAllIssues() []models.Issue
	GetIssueByID(id int) (*models.Issue, error)
	DeleteIssue(id int) error
	PatchIssue(id int, upd_req models.UpdateIssueRequest) error
}

*/

// Mock Service -> NO DB
type testService struct{}

func (t *testService) GetIssueByID(id int) (*models.Issue, error) {
	return &models.Issue{
		Internal_ID: int64(id),
		Title:       "test issue",
	}, nil
}
func TestRouter(s *testService) *Router {
	return &Router{
		issueService: s,
	}
}

/*
CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error)
*/
func (t *testService) GetAllIssues() ([]models.Issue, error) {
	return []models.Issue{}, nil
}

func (t *testService) CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error) {
	return nil, nil
}

func (t *testService) DeleteIssue(id int) error {
	return nil
}

func (t *testService) PatchIssue(id int, upd_req models.UpdateIssueRequest) error {
	return nil
}

func TestGetSingleIssue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/issues/99", nil)
	w := httptest.NewRecorder()

	service := testService{}
	r := TestRouter(&service)

	r.getSingleIssueHandler(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	if w.Body.Len() == 0 {
		t.Fatalf("expected body, got empty")
	}
}
