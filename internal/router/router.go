package router

import (
	"encoding/json"
	"fmt"
	"issuetracker/internal/models"
	services "issuetracker/internal/services"
	"log"
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

func (s *Router) AllRouting(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	switch len(parts) {
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
			return
		}
		if r.Method == http.MethodDelete {
			fmt.Printf("Tried to delete a single issue: %s\n", r.URL.Path)
			s.deleteSingleIssueHandler(w, r)
			return
		}
		if r.Method == http.MethodPatch {
			fmt.Printf("Tried to update a single issue: %s\n", r.URL.Path)
			s.updateSingleIssueHandler(w, r)
		}
	default:
		http.Error(w, "method now allowed", http.StatusMethodNotAllowed)
	}
}

// Gets a single issue from the database.
func (s *Router) getSingleIssueHandler(w http.ResponseWriter, r *http.Request) {

	id := extractIdFromUrlPath(r.URL.Path)
	if id < 0 {
		http.Error(w, "bad id ", http.StatusBadRequest)
		return
	}

	issue, err := s.issueService.GetSingleIssue(int(id))
	if err != nil {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	resp := issueToIssueResponse(*issue)
	json.NewEncoder(w).Encode(resp)
	fmt.Printf("*** TESTING - RETURNED ***\n ID: %d\n TITLE: %s\n EXTERNAL REF: %s\n", resp.Internal_ID, resp.Title, resp.External_Ref)
}

func (s *Router) getIssueshandler(w http.ResponseWriter, r *http.Request) {
	/*
		query := r.URL.Query()
		resolved := query.Get("resolved")
		search := query.Get("search")
	*/
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

	issue := s.issueService.CreateNewIssue(req.External_Ref, req.Title, req.Description)

	resp := issueToIssueResponse(issue)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resp)
}

func (s *Router) deleteSingleIssueHandler(w http.ResponseWriter, r *http.Request) {
	id := extractIdFromUrlPath(r.URL.Path)
	if id < 0 {
		http.Error(w, "bad id ", http.StatusBadRequest)
	}

	//TODO
}

func (s *Router) updateSingleIssueHandler(w http.ResponseWriter, r *http.Request) {
	id := extractIdFromUrlPath(r.URL.Path)
	if id < 0 {
		http.Error(w, "bad id ", http.StatusBadRequest)
	}
	var upd_req models.UpdateIssueRequest

	fmt.Printf("r.Body: %v\n", r.Body)
	err := json.NewDecoder(r.Body).Decode(&upd_req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	updated, err := s.issueService.PatchIssue(int(id), upd_req)
	json.NewEncoder(w).Encode(issueToIssueResponse(*updated))
	fmt.Printf("*** TESTING - PATCHED ***\n ID: %d\n TITLE: %s\n EXTERNAL REF: %s\n", updated.Internal_ID, updated.Title, updated.External_Ref)

}

// Helpers

func issueToIssueResponse(issue models.Issue) models.IssueResponse {
	resp := models.IssueResponse{
		Internal_ID:  issue.Internal_ID,
		External_Ref: issue.External_Ref,
		Title:        issue.Title,
		Description:  issue.Description,
		Active:       issue.Active,
	}
	return resp
}

func extractIdFromUrlPath(path string) int64 {
	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return id
}
