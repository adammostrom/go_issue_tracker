package router

import (
	"encoding/json"
	"fmt"
	"issuetracker/internal/models"
	"net/http"
	"strings"
)

/*
The router/handler is responsible for:
- reading request body
- parsing JSON
- calling the service
- writing JSON response
*/

type IssueServiceInterface interface {
	CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error)
	GetAllIssues() ([]models.Issue, error)
	GetIssueByID(id int) (*models.Issue, error)
	DeleteIssue(id int) error
	PatchIssue(id int, upd_req models.UpdateIssueRequest) error
	GetLogsFromIssue(id int) (*[]models.LogEntry, error)
}

type Router struct {
	issueService IssueServiceInterface
}

// Constructor that accepts an issue service interface and returns a router pointer.
func NewRouter(s IssueServiceInterface) *Router {
	return &Router{
		issueService: s,
	}
}

func (s *Router) AllRouting(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	switch len(parts) {
	case 1:
		http.Error(w, "bad request", http.StatusBadRequest)
	case 2:
		if r.Method == http.MethodGet {
			s.getIssuesHandler(w, r)
			return
		} else if r.Method == http.MethodPost {
			s.createIssueHandler(w, r)
			return
		}
	case 3:
		if r.Method == http.MethodGet {
			s.getSingleIssueHandler(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			s.deleteSingleIssueHandler(w, r)
			return
		}
		if r.Method == http.MethodPatch {
			s.PatchIssueHandler(w, r)
		}
	case 4: // /issues/{id}/logs ---> / issues / id / logs
		if r.Method == http.MethodGet {
			s.GetLogsFromIssueHandler(w, r)
			return
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Gets a single issue from the database.
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

func (s *Router) getIssuesHandler(w http.ResponseWriter, r *http.Request) {
	/*
		query := r.URL.Query()
		resolved := query.Get("resolved")
		search := query.Get("search")
	*/
	issues, err := s.issueService.GetAllIssues()
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	var response []models.IssueResponse

	for _, issue := range issues {
		response = append(response, issueToIssueResponse(issue))
	}

	json.NewEncoder(w).Encode(response) // Fetches the IssueService slice reference (returns the whole slice)

	w.WriteHeader(http.StatusOK)

}

func (s *Router) createIssueHandler(w http.ResponseWriter, r *http.Request) {

	var req models.CreateIssueRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	issue, err_post := s.issueService.CreateNewIssue(req)

	if err_post != nil {
		http.Error(w, err_post.Error(), http.StatusBadRequest)
	}

	resp := issueToIssueResponse(*issue)

	json.NewEncoder(w).Encode(resp)

	w.WriteHeader(http.StatusCreated)

}

func (s *Router) deleteSingleIssueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.issueService.DeleteIssue(int(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Router) PatchIssueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var upd_req models.UpdateIssueRequest

	err_decode := json.NewDecoder(r.Body).Decode(&upd_req)
	if err_decode != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err_patch := s.issueService.PatchIssue(int(id), upd_req)
	if err_patch != nil {
		http.Error(w, err_patch.Error(), http.StatusBadRequest)
	}

	updated, err_updated := s.issueService.GetIssueByID(int(id))
	if err_updated != nil {
		http.Error(w, err_updated.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(issueToIssueResponse(*updated))
	fmt.Printf("*** TESTING - PATCHED ***\n ID: %d\n TITLE: %s\n EXTERNAL REF: %s\n", updated.Internal_ID, updated.Title, updated.External_Ref)

	w.WriteHeader(http.StatusOK)
}

func (s *Router) GetLogsFromIssueHandler(w http.ResponseWriter, r *http.Request) {

	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logs, err := s.issueService.GetLogsFromIssue(int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
	for log, _ := range logs {
		json.NewEncoder(w).Encode((log))
	}

	w.WriteHeader(http.StatusOK)

}
