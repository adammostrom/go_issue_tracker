package cli

import (
	"bufio"
	"fmt"
	"issuetracker/internal/models"
	"os"
	"strconv"
	"strings"
)

type IssueServiceInterface interface {
	CreateNewIssue(req models.CreateIssueRequest) (*models.Issue, error)
	GetAllIssues(status models.IssueStatus) ([]models.Issue, error)
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

func (s *CommandLineInterface) GetIssues(status models.IssueStatus) ([]models.Issue, error) {
	issues, err := s.issueService.GetAllIssues(status)
	if err != nil {
		fmt.Println("Could not find issues")
		return nil, err

	}

	fmt.Printf("Found %d issues\n", len(issues))
	for i := range issues {
		s.printIssue(&issues[i])
	}
	return issues, nil
}

func (s *CommandLineInterface) GetIssue(id int) (*models.Issue, error) {
	issue, err := s.issueService.GetIssueByID(id)
	if err != nil || issue == nil {
		fmt.Printf("Could not find issue with id: %d\n", id)
		return nil, err
	}
	s.printIssue(issue)
	return issue, nil
}

func (s *CommandLineInterface) CreateIssue() error {
	var issue_request models.CreateIssueRequest

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read title: %s\n", title)
		return err
	}
	issue_request.Title = strings.TrimSpace(title)

	fmt.Print("External Reference: ")
	external_ref, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read external ref : %s\n", external_ref)
		return err
	}
	issue_request.External_Ref = strings.TrimSpace(external_ref)

	fmt.Print("Description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read description : %s\n", desc)
		return err
	}
	issue_request.Description = strings.TrimSpace(desc)

	issue, err := s.issueService.CreateNewIssue(issue_request)
	if err != nil {
		fmt.Printf("Failed to create issue: %s\n", err)
		return err
	}
	fmt.Println("Issue successfully created")
	s.printIssue(issue)
	return nil
}

func (c *CommandLineInterface) Run(args []string) {

	if len(args) < 1 {
		fmt.Println("expected subcommand")
		return
	}
	switch args[0] {
	case "list":
		c.listCmd(args[1:])
	case "show":
		if len(args) < 2 {
			fmt.Printf("No valid argument provided for SHOW subcommand: %s\n", args)
			return
		}
		c.showCmd(args[1])
	case "create":
		c.CreateIssue()
	}

}

func (c *CommandLineInterface) showCmd(arg string) {

	id, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Printf("invalid id: %s\n", arg)
		return
	}
	issue, err := c.GetIssue(id)
	if err != nil {
		return
	}
	fmt.Printf("single issue found: %s\n", issue.Title)
}

func (c *CommandLineInterface) listCmd(args []string) {

	status := models.StatusDefault

	if len(args) > 0 {
		parsed, err := models.ParseStatus(args[0])
		if err != nil {
			fmt.Println("invalid status")
			return
		}
		status = parsed
	}
	c.GetIssues(status)
}

func (s *CommandLineInterface) printIssue(issue *models.Issue) error {
	if issue == nil {
		return nil
	}

	str := fmt.Sprintf("%t", issue.Active)

	fmt.Printf("ID: %d\n", issue.Internal_ID)
	fmt.Printf("Title: %s\n", issue.Title)
	fmt.Printf("External Reference: %s\n", issue.External_Ref)
	fmt.Printf("Description: %s\n", issue.Description)
	fmt.Printf("Active Status: %s\n", str)

	return nil
}
