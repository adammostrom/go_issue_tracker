package router

import (
	"encoding/json"
	"issuetracker/internal/models"
	"log"
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
	GetAllIssues(status models.IssueStatus) ([]models.Issue, error)
	GetIssueByID(id int) (*models.Issue, error)
	DeleteIssue(id int) error
	PatchIssue(id int, upd_req models.UpdateIssueRequest) error
	GetLogsFromIssue(id int) ([]models.LogEntry, error)
	AddLogEntry(id int, entry string) error
	DeleteLogsFromIssue(id int) error
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
		return
	// /issues/
	case 2:
		if r.Method == http.MethodGet {
			s.getIssuesHandler(w, r)
			return
		} else if r.Method == http.MethodPost {
			s.createIssueHandler(w, r)
			return
		}
	// /issues/{id}
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
	// /issues/{id}/logs
	case 4:
		if r.Method == http.MethodGet {
			s.GetLogsFromIssueHandler(w, r)
			return
		}
		if r.Method == http.MethodPost {
			s.AddLogEntryHandler(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			s.DeleteLogsFromIssueHandler(w, r)
			return
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Gets a single issue from the database.
func (s *Router) getSingleIssueHandler(w http.ResponseWriter, r *http.Request) {

	var resp models.IssueResponse

	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		log.Printf("invalid id in path: %v\n", err)
		http.Error(w, "invalid id\n", http.StatusBadRequest)
		return
	}

	issue, err := s.issueService.GetIssueByID(int(id))
	if err != nil {
		log.Printf("failed to get issue with id: %d\n", id)
		http.Error(w, "failed to get issue\n", http.StatusBadRequest)
		return
	}
	if issue == nil {
		log.Printf("Failed to get issue -> issue = %v\n", issue)
		http.Error(w, "failed to get issue\n", http.StatusBadRequest)
		return

	} else {
		resp = issueToIssueResponse(*issue)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

// TODO: 2026-04-06: Find a way to not duplicate these functions.
func (s *Router) getIssuesHandler(w http.ResponseWriter, r *http.Request) {

	statusStr := r.URL.Query().Get("status")

	status := models.StatusDefault

	if statusStr != "" {
		parsed, err := models.ParseStatus(statusStr)
		if err != nil {
			http.Error(w, "invalid status\n", http.StatusBadRequest)
			return
		}
		status = parsed
	}

	issues, err := s.issueService.GetAllIssues(status)
	if err != nil {
		log.Printf("Failed to get issues\n")
		http.Error(w, "failed to get issues\n", http.StatusBadRequest)
		return
	}

	var response []models.IssueResponse

	for _, issue := range issues {
		response = append(response, issueToIssueResponse(issue))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response) // Fetches the IssueService slice reference (returns the whole slice)

}

func (s *Router) createIssueHandler(w http.ResponseWriter, r *http.Request) {

	var req models.CreateIssueRequest

	var resp models.IssueResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("failed to decode request body: %v\n", err)
		http.Error(w, "bad request\n", http.StatusBadRequest)
		return
	}

	issue, err_post := s.issueService.CreateNewIssue(req)

	if err_post != nil {
		log.Printf("Failed to create issue -> issue = %v\n", issue)
		http.Error(w, "failed to create new issuen\n", http.StatusBadRequest)
		return
	}

	if issue == nil {
		log.Printf("Failed to create issue -> issue = %v\n", issue)
		http.Error(w, "issue not created\n", http.StatusBadRequest)
		return

	} else {
		resp = issueToIssueResponse(*issue)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

func (s *Router) deleteSingleIssueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		log.Printf("invalid id in path: %v\n", err)
		http.Error(w, "invalid id\n", http.StatusBadRequest)
		return
	}

	err = s.issueService.DeleteIssue(int(id))

	if err != nil {
		log.Printf("Error deleting issue: %v with id: %d\n", err, id)
		http.Error(w, "Failed to delete issue\n", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Router) PatchIssueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		log.Printf("invalid id in path: %v\n", err)
		http.Error(w, "invalid id\n", http.StatusBadRequest)
		return
	}

	var updReq models.UpdateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&updReq); err != nil {
		log.Printf("failed to decode request body: %v\n", err)
		http.Error(w, "bad request\n", http.StatusBadRequest)
		return
	}

	if err := s.issueService.PatchIssue(int(id), updReq); err != nil {
		log.Printf("failed to patch issue %d: %v\n", id, err)
		http.Error(w, "failed to update issue\n", http.StatusBadRequest)
		return
	}

	updated, err := s.issueService.GetIssueByID(int(id))
	if err != nil {
		log.Printf("failed to fetch updated issue %d: %v\n", id, err)
		http.Error(w, "failed to fetch updated issue\n", http.StatusInternalServerError)
		return
	}

	if updated == nil {
		log.Printf("issue %d not found after update\n", id)
		http.Error(w, "issue not found\n", http.StatusNotFound)
		return
	}

	resp := issueToIssueResponse(*updated)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Returns all the logs from a specific issue.
func (s *Router) GetLogsFromIssueHandler(w http.ResponseWriter, r *http.Request) {

	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		log.Printf("invalid id in path: %v\n", err)
		http.Error(w, "invalid id\n", http.StatusBadRequest)
		return
	}

	logs, err := s.issueService.GetLogsFromIssue(int(id))
	if err != nil {
		log.Printf("Error fetching logs: %v\n", err)
		http.Error(w, "Failed to fetch logs\n", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the whole slice at once
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		// optional: log error, cannot write another HTTP status here
	}
}

type LogRequest struct {
	Entry string `json:"entry"`
}

func (s *Router) AddLogEntryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		log.Printf("invalid id in path: %v\n", err)
		http.Error(w, "invalid id\n", http.StatusBadRequest)
		return
	}

	var req LogRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON\n", http.StatusBadRequest)
		return
	}

	err = s.issueService.AddLogEntry(int(id), req.Entry)
	if err != nil {
		log.Printf("Failed to add log entry: %s\n", req.Entry)
		http.Error(w, "failed to add entry\n", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)

}

func (s *Router) DeleteLogsFromIssueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDfromPath(r.URL.Path)
	if err != nil {
		log.Printf("invalid id in path: %v\n", err)
		http.Error(w, "invalid id\n", http.StatusBadRequest)
		return
	}

	err = s.issueService.DeleteLogsFromIssue(int(id))

	if err != nil {
		log.Printf("Error deleting logs for issue with id: %d. Error: %v\n", id, err)
		http.Error(w, "Failed to delete logs\n", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
