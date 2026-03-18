package router

import (
	"encoding/json"
	"fmt"
	"issuetracker/internal/models"
	services "issuetracker/internal/services"
	"net/http"
	"strconv"
	"strings"
)

/*
So the handler is responsible for:

reading request body

parsing JSON

calling the service

writing JSON response

*/

type Router struct {
	issueService *services.IssueService
}

// Constructor that accepts an issue service and returns a router pointer.
func NewRouter(s *services.IssueService) *Router {
	return &Router{
		issueService: s,
	}
}

func (s *Router) MainDelegator(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	length := len(parts)
	switch length {
	case 1:
		http.Error(w, "bad request", http.StatusBadRequest)
	case 2:
		if r.Method == http.MethodGet {
			s.getIssueshandler(w, r)
			return
		} else if r.Method == http.MethodPost {
			s.createIssueHandler(w, r)
			return
		}
	case 3:
		if r.Method == http.MethodGet {
			s.getSingleIssueHandler(w, r)
			fmt.Printf("Tried to fetch single issue: %s", r.URL.Path)
			return
		}
	default:
		http.Error(w, "method now allowed", http.StatusMethodNotAllowed)
	}
}

// Receives a http request from mux, parses and handles the request.
func (s *Router) AllRouting(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		s.getIssueshandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		s.createIssueHandler(w, r)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Router) getSingleIssueHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	issue := s.issueService.GetSingleIssue(int(id))
	if err != nil {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	resp := issueToIssueResponse(issue)
	json.NewEncoder(w).Encode(resp)
}

func (s *Router) getIssueshandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Ask the database for the data, DONT pass any HTTP stuff (w, r)

	issues := s.issueService.GetAllIssues()

	var response []models.IssueResponse

	for _, issue := range issues {
		response = append(response, issueToIssueResponse(issue))
	}

	json.NewEncoder(w).Encode(response) // Fetches the IssueService slice reference (returns the whole slice)

}

func (s *Router) createIssueHandler(w http.ResponseWriter, r *http.Request) {

	var req models.CreateIssueRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	issue := s.issueService.CreateNewIssue(req.ExternalRef, req.Title, req.Description)

	resp := issueToIssueResponse(issue)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resp)
}

func issueToIssueResponse(issue models.Issue) models.IssueResponse {
	resp := models.IssueResponse{
		InternalID:  issue.Internal_id,
		ExternalRef: issue.External_Ref,
		Title:       issue.Title,
		Description: issue.Description,
		Active:      issue.Active,
	}
	return resp
}
