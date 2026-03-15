package services

import (
	"encoding/json"
	"issuetracker/internal/models"
	"net/http"
)

type IssueService struct {
	Issues *[]models.Issue // Pointer reference to the temporary storage
}

func (s *IssueService) MainRouter(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		s.getIssues(w, r)
		return
	}

	if r.Method == http.MethodPost {
		s.createIssue(w, r)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *IssueService) getIssues(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(*s.Issues) // Fetches the IssueService slice reference

}

func (s *IssueService) createIssue(w http.ResponseWriter, r *http.Request) {

	var issue models.Issue

	err := json.NewDecoder(r.Body).Decode(&issue)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	*s.Issues = append(*s.Issues, issue)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issue)
}
