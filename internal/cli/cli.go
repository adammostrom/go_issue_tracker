package cli

import (
	"fmt"
	"issuetracker/internal/models"
)

type IssueServiceInterface interface {
	CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error)
	GetAllIssues() ([]models.Issue, error)
	GetIssueByID(id int) (*models.Issue, error)
	DeleteIssue(id int) error
	PatchIssue(id int, upd_req models.UpdateIssueRequest) error
	GetLogsFromIssue(id int) ([]models.LogEntry, error)
}

type CommandLineInterface struct {
	issueService IssueServiceInterface
}

func NewCLI(s IssueServiceInterface) *CommandLineInterface {
	return &CommandLineInterface{
		issueService: s,
	}
}

func (s *CommandLineInterface) GetIssues() ([]models.Issue, error) {
	issues, err := s.issueService.GetAllIssues()
	if err != nil {
		fmt.Println("Could not find issues")
		return nil, err

	}

	for _, issue := range issues {
		fmt.Printf("issue found: %s\n", issue.Title)
	}
	fmt.Printf("Found %d issues\n", len(issues))
	return issues, nil
}

func (s *CommandLineInterface) GetIssue(id int) (*models.Issue, error) {
	issue, err := s.issueService.GetIssueByID(id)
	if err != nil || issue == nil {
		fmt.Printf("Could not find issue with id: %d\n", id)
		return nil, err
	}
	return issue, nil
}
