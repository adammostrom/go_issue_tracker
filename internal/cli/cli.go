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

type CommandLineInterface struct {
	issueService IssueServiceInterface
}

func NewCLI(s IssueServiceInterface) *CommandLineInterface {

	cli := &CommandLineInterface{
		issueService: s,
	}
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

// TODO: 2026-04-27: Split into separate functions that are each called here by this function, with each of the having validation at the end.

func readTitle(reader *bufio.Reader) (string, error) {

	fmt.Print("Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read title: %s\n", title)
		return "", err // TODO: Fix to not return empty string
	}
	err = models.ValidateTitle(title)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(title), nil
}

func readExtRef(reader *bufio.Reader) (string, error) {
	fmt.Print("External Reference: ")
	external_ref, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Could not read external ref : %s\n", external_ref)
		return "", err
	}
	err = models.ValidateExternalRef(external_ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(external_ref), nil
}

func readDescription(reader *bufio.Reader) (string, error) {
	fmt.Print("Description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read description : %s\n", desc)
		return "", err
	}
	err = models.ValidateDescription(desc)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(desc), nil
}

func (s *CommandLineInterface) CreateIssue() error {
	var issue_request models.CreateIssueRequest

	reader := bufio.NewReader(os.Stdin)

	title, err := readValidated(reader, "Title: ", models.ValidateTitle)
	if err != nil {
		return err
	}
	issue_request.Title = title

	externalRef, err := readValidated(reader, "External Reference: ", models.ValidateExternalRef)
	if err != nil {
		return err
	}
	issue_request.External_Ref = externalRef

	description, err := readValidated(reader, "Description: ", models.ValidateDescription)
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
	s.printIssue(issue)
	return nil
}

func readValidated(reader *bufio.Reader, prompt string, validate func(string) error) (string, error) {

	for {
		fmt.Print(prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)

		if err := validate(input); err != nil {
			fmt.Println(err)
			continue // Ask again
		}

		return input, nil
	}
}

func (s *CommandLineInterface) CreateLogEntry(id int, entry string) error {
	err := s.issueService.AddLogEntry(id, entry)
	if err != nil {
		fmt.Printf("Failed to create log entry: %s", err)
		return err
	}
	return nil
}

func (s *CommandLineInterface) DeleteLogsForIssue(id int) error {
	return s.issueService.DeleteLogsFromIssue(id)
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
	fmt.Printf("Progrss:            %s\n", issue.Progress.String())
	fmt.Println("-------------------------")

	return nil
}

func (s *CommandLineInterface) simplePrintIssue(issue *models.Issue) error {
	if issue == nil {
		return nil
	}

	fmt.Println(simplePrintString(issue))
	return nil
}

func simplePrintString(i *models.Issue) string {
	return fmt.Sprintf(
		"%s - %d %s %s %s",
		progressSymbol(i.Progress),
		i.Internal_ID,
		layoutDistancePrint(i.External_Ref, models.EXTERNAL_MAX),
		layoutDistancePrint(i.Title, models.TITLE_MAX),
		i.Description,
	)
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
