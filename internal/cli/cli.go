package cli

import (
	"bufio"
	"fmt"
	"issuetracker/internal/models"
	"os"
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
	DeleteLogsFromIssue(id int) error
}

type CommandLine struct {
	issueService IssueServiceInterface
}

func NewCLI(s IssueServiceInterface) *CommandLine {

	cli := &CommandLine{
		issueService: s,
	}
	return cli
}

type Command struct {
	name        string
	operation   func(args []string)
	subcommands map[string]*Command
}

func (s *CommandLine) BuildCommands() map[string]*Command {
	return map[string]*Command{
		"list": {
			name:      "list",
			operation: s.listCmd,
		},
		"get": {
			name:      "get",
			operation: s.getCmd,
		},
		// TODO: started, idle, completed
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
				"progress": {
					name:      "progress",
					operation: s.setProgressCmd,
				},
			},
		},
		"create": {
			name:      "create",
			operation: s.createCmd,
		},
		"modify": {
			name:      "modify",
			operation: s.modifyCmd,
		},
		// TODO: Patch
		// TODO: Delete
		// TODO: Create issue
		"log": {
			name: "log",
			operation: func(args []string) {
				fmt.Println("Available subcommands: get, create, delete")
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
				"delete": {
					name:      "delete",
					operation: s.deleteLogsCmd,
				},
			},
		},
	}
}

func (s *CommandLine) Run(cmds map[string]*Command, args []string) {
	s.dispatch(cmds, args)
}

func (s *CommandLine) PrintCommands(cmds map[string]*Command, depth int) {
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

func (s *CommandLine) dispatch(cmds map[string]*Command, args []string) {
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

func (s *CommandLine) setProgress(id int, status models.ProgressStatus) error {

	req := models.UpdateIssueRequest{
		Progress: &status,
	}
	err := s.ModifyIssue(id, req)
	if err != nil {
		return err
	}
	fmt.Printf("Progress status changed successfully to: %s", req.Progress.String())
	return nil

}

func (s *CommandLine) modifyIssueHelper(id int) error {

	reader := bufio.NewReader(os.Stdin)

	req := models.UpdateIssueRequest{}

	title, err := readValidated(reader, "Title: ", models.ValidateTitle, true)
	if err != nil {
		return err
	}
	if title != "" {
		req.Title = &title
	}

	ext_ref, err := readValidated(reader, "External Ref: ", models.ValidateExternalRefWrapper, true)
	if err != nil {
		return err
	}
	if ext_ref != "" {
		req.External_Ref = &ext_ref
	}

	desc, err := readValidated(reader, "Description: ", models.ValidateDescription, true)
	if err != nil {
		return err
	}
	if desc != "" {
		req.Description = &desc
	}

	return s.ModifyIssue(id, req)
}

func (s *CommandLine) ModifyIssue(id int, upd_req models.UpdateIssueRequest) error {

	err := s.issueService.PatchIssue(id, upd_req)
	if err != nil {
		return err
	}
	fmt.Println("Issue modified successfully with values: ")
	fmt.Printf("upd_req: %v\n", &upd_req)
	return nil
}

func (s *CommandLine) GetIssues(status models.IssueStatus) ([]models.Issue, error) {
	issues, err := s.issueService.GetAllIssues(status)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (s *CommandLine) GetIssue(id int) (*models.Issue, error) {
	issue, err := s.issueService.GetIssueByID(id)
	if err != nil || issue == nil {
		return nil, err
	}
	return issue, nil
}

// TODO: 2026-04-27: Split into separate functions that are each called here by this function, with each of the having validation at the end.

func (s *CommandLine) CreateIssue() error {
	var issue_request models.CreateIssueRequest

	reader := bufio.NewReader(os.Stdin)

	title, err := readValidated(reader, "Title: ", models.ValidateTitle, false)
	if err != nil {
		return err
	}
	issue_request.Title = title

	externalRef, err := readValidated(reader, "External Reference: ", models.ValidateExternalRefWrapper, true)
	if err != nil {
		return err
	}
	if externalRef == "" {
		issue_request.External_Ref = nil // IMPORTANT
	} else {
		issue_request.External_Ref = &externalRef
	}

	description, err := readValidated(reader, "Description: ", models.ValidateDescription, false)
	if err != nil {
		return err
	}
	issue_request.Description = description

	issue, err := s.issueService.CreateNewIssue(issue_request)
	if err != nil {
		fmt.Printf("Failed to create issue: %s\n", err)
		return err
	}
	fmt.Println("Issue successfully created")
	printIssue(issue)
	return nil
}

func (s *CommandLine) CreateLogEntry(id int, entry string) error {
	err := s.issueService.AddLogEntry(id, entry)
	if err != nil {
		fmt.Printf("Failed to create log entry: %s", err)
		return err
	}
	return nil
}

func (s *CommandLine) DeleteLogsForIssue(id int) error {
	return s.issueService.DeleteLogsFromIssue(id)
}

func readValidated(reader *bufio.Reader, prompt string, validate func(string) error, allowEmpty bool) (string, error) {

	for {
		fmt.Print(prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)

		if input == "" {
			allowEmpty = true
			if allowEmpty {
				return "", nil // For modifying, allow empty input
			}
			continue // else, for create issue, force validation
		}

		if err := validate(input); err != nil {
			fmt.Println(err)
			continue // Ask again
		}

		return input, nil
	}
}

// TODO: 2026-04-14: come back to update this. Somewhat clunky.
func printIssue(issue *models.Issue) error {
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
	fmt.Printf("Progrss:            %s\n", issue.Progress.String())
	fmt.Println("-------------------------")

	return nil
}

func simplePrintIssue(issue *models.Issue, issueService IssueServiceInterface) error {
	if issue == nil {
		return nil
	}

	fmt.Println(simplePrintString(issue, issueService))
	return nil
}

func simplePrintString(i *models.Issue, issueService IssueServiceInterface) string {
	logs, err := issueService.GetLogsFromIssue(int(i.Internal_ID))
	if err != nil {
		return fmt.Sprint(err)
	}
	created := logs[0]

	ext_ref_print := deref(i.External_Ref, "null")

	return fmt.Sprintf(
		"%s - %s %d %s | %s",
		progressSymbol(i.Progress),
		created,
		i.Internal_ID,
		layoutDistancePrint(ext_ref_print, models.EXTERNAL_MAX),
		layoutDistancePrint(i.Title, models.TITLE_MAX),
	)
}

func deref(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}
func layoutDistancePrint(param string, max int) string {
	dist_param := max - len(param)
	distance := strings.Repeat(" ", dist_param)

	return param + distance
}

func progressSymbol(p models.ProgressStatus) string {
	switch p {
	case models.Idle:
		return "[ ]"
	case models.Started:
		return "[/]"
	case models.Finished:
		return "[X]"
	default:
		return "[-]"
	}
}
