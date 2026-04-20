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
	AddLogEntry(id int, entry string) error
}

type CommandLineInterface struct {
	issueService IssueServiceInterface
}

func NewCLI(s IssueServiceInterface) *CommandLineInterface {

	cli := &CommandLineInterface{
		issueService: s,
	}
	/* 	var commands = Command{
		Call: "list",
		Op:   cli.listCmd, // ✅ safe, cli already exists
	} */
	return cli
}

type Command struct {
	name        string
	operation   func(args []string)
	subcommands map[string]*Command
}

func (s *CommandLineInterface) BuildCommands() map[string]*Command {
	return map[string]*Command{
		"list": {
			name:      "list",
			operation: s.listCmd,
		},
		"get": {
			name:      "get",
			operation: s.getCmd,
		},
		"set": {
			name: "set",
			operation: func(args []string) {
				fmt.Println("Available subcommands: active, inactive")
			},
			subcommands: map[string]*Command{
				"active": {
					name:      "active",
					operation: s.setActiveCmd,
				},
				"inactive": {
					name:      "inactive",
					operation: s.setInactiveCmd,
				},
			},
		},
		"create": {
			name:      "create",
			operation: s.createCmd,
		},
		// TODO: Patch
		// TODO: Delete
		// TODO: Create issue
		"log": {
			name: "log",
			operation: func(args []string) {
				fmt.Println("Available subcommands: get, create")
			},
			subcommands: map[string]*Command{
				"get": {
					name:      "get",
					operation: s.getLogCmd,
				},
				"create": {
					name:      "create",
					operation: s.createLogCmd,
				},
			},
		},
	}
}

func (s *CommandLineInterface) Run(cmds map[string]*Command, args []string) {
	s.dispatch(cmds, args)
}

func (s *CommandLineInterface) PrintCommands(cmds map[string]*Command, depth int) {
	for name, cmd := range cmds {
		// indent based on depth
		for i := 0; i < depth; i++ {
			fmt.Print("  ")
		}

		fmt.Println(name)

		// recurse into subcommands
		if cmd.subcommands != nil {
			s.PrintCommands(cmd.subcommands, depth+1)
		}
	}
}

func (s *CommandLineInterface) dispatch(cmds map[string]*Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Expected subcommand")
		return
	}

	current, ok := cmds[args[0]]
	if !ok {
		fmt.Println("Unknown command: ", args[0])
	}

	args = args[1:]

	for len(args) > 0 {
		if current.subcommands == nil {
			break
		}

		next, ok := current.subcommands[args[0]]
		if !ok {
			break
		}

		current = next
		args = args[1:]
	}
	if current.operation != nil {
		current.operation(args)
	} else {
		fmt.Println("Missing subcommand:", current.name)
	}

}

func (s *CommandLineInterface) setActiveCmd(args []string) {
	s.setStatusCmd(args, true)
}

func (s *CommandLineInterface) setInactiveCmd(args []string) {
	s.setStatusCmd(args, false)
}

func (s *CommandLineInterface) setStatusCmd(args []string, status bool) {
	if len(args) < 1 {
		fmt.Println("Expected id")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("invalid id: %s\n", args[0])
		return
	}

	req := models.UpdateIssueRequest{
		Active: &status,
	}

	if err := s.PatchIssue(id, req); err != nil {
		fmt.Printf("Failed to update issue: %s\n", err)
	}
}

func (s *CommandLineInterface) getLogCmd(args []string) {
	if len(args) < 1 {
		fmt.Println("Expected id as argument to command!")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("invalid id: %s\n", args[0])
		fmt.Println(err)
		return
	}
	logs, err := s.issueService.GetLogsFromIssue(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(logs) == 0 {
		fmt.Println("No logs found")
		return
	}
	for _, log := range logs {
		fmt.Printf("log: %v\n", log)
	}

}

func (s *CommandLineInterface) createLogCmd(args []string) {

}

func (s *CommandLineInterface) logCmd(args []string) {

	if len(args) < 1 {
		fmt.Println("expected subcommand: ")
	}
}

// Get one issue
func (s *CommandLineInterface) getCmd(args []string) {

	if len(args) < 1 {
		fmt.Println("Expected id as argument to command!")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("invalid id: %s\n", args[0])
		fmt.Println(err)
		return
	}
	issue, err := s.GetIssue(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.printIssue(issue)
}

// Get all issues
func (s *CommandLineInterface) listCmd(args []string) {

	status := models.StatusDefault

	if len(args) > 0 {
		parsed, err := models.ParseStatus(args[0])
		if err != nil {
			fmt.Println("invalid status")
			return
		}
		status = parsed
	}
	issues, err := s.GetIssues(status)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(issues) == 0 {
		fmt.Println("No issues found")
	}
	for _, issue := range issues {
		s.printIssue(&issue)
	}
}

func (s *CommandLineInterface) createCmd(args []string) {

	err := s.CreateIssue()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (s *CommandLineInterface) PatchIssue(id int, upd_req models.UpdateIssueRequest) error {

	err := s.issueService.PatchIssue(id, upd_req)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommandLineInterface) GetIssues(status models.IssueStatus) ([]models.Issue, error) {
	issues, err := s.issueService.GetAllIssues(status)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (s *CommandLineInterface) GetIssue(id int) (*models.Issue, error) {
	issue, err := s.issueService.GetIssueByID(id)
	if err != nil || issue == nil {
		return nil, err
	}
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

	for {
		fmt.Print("External Reference: ")
		external_ref, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("Could not read external ref : %s\n", external_ref)
			continue
		}

		err = models.ValidateExternalRef(external_ref)
		if err != nil {
			fmt.Println(err)
			continue
		}

		issue_request.External_Ref = strings.TrimSpace(external_ref)
		break
	}

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

func (s *CommandLineInterface) CreateLogEntry(id int, entry string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Entry: ")
	entry, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read entry: %s\n", entry)
		return err
	}
	err = s.issueService.AddLogEntry(id, entry)
	if err != nil {
		fmt.Printf("Failed to create log entry: %s", err)
		return err
	}
	return nil
}

// TODO: 2026-04-14: come back to update this. Somewhat clunky.
func (s *CommandLineInterface) printIssue(issue *models.Issue) error {
	if issue == nil {
		return nil
	}

	//str := fmt.Sprintf("%t", issue.Active)

	fmt.Println("-------------------------")
	fmt.Printf("ID:                 %d\n", issue.Internal_ID)
	fmt.Printf("Title:              %s\n", issue.Title)
	fmt.Printf("External Reference: %s\n", issue.External_Ref)
	fmt.Printf("Description:        %s\n", issue.Description)
	fmt.Printf("Active Status:      %t\n", issue.Active)
	fmt.Println("-------------------------")

	return nil
}
