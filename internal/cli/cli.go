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
	GetAllIssues(filter models.IssueFilter) ([]models.Issue, error)
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
	operation   func(args []string) error
	description string
	subcommands map[string]*Command
}

func (s *CommandLine) BuildCommands() map[string]*Command {
	return map[string]*Command{
		"list": {
			name:        "list",
			description: "List all issues",
			operation:   s.getAllCmd,
		},
		"get": {
			name:        "get",
			description: "Get issue <id>",
			operation:   s.getIssueCmd,
		},
		"create": {
			name:        "create",
			description: "Create a new issue",
			operation:   s.createCmd,
		},
		"modify": {
			name:        "modify",
			description: "Modify issue <id>",
			operation:   s.modifyCmd,
		},
		"delete": {
			name:        "delete",
			description: "Delete issue <id>",
			operation:   s.deleteCmd,
		},
		"set": {
			name:        "set",
			description: "Set issue state: active | inactive | progress",
			operation: func(args []string) error {
				fmt.Println("Subcommands: active, inactive, progress")
				return nil
			},
			subcommands: map[string]*Command{
				"active": {
					name:        "active",
					description: "Set issue <id> active",
					operation:   s.setActiveCmd,
				},
				"inactive": {
					name:        "inactive",
					description: "Set issue <id> inactive",
					operation:   s.setInactiveCmd,
				},
				"progress": {
					name:        "progress",
					description: "Set progress <id> <idle|started|completed>",
					operation:   s.setProgressCmd,
				},
			},
		},
		"log": {
			name:        "log",
			description: "Manage logs: get | create | delete",
			operation: func(args []string) error {
				fmt.Println("Subcommands: get, create, delete")
				return nil
			},
			subcommands: map[string]*Command{
				"get": {
					name:        "get",
					description: "Get logs for <id>",
					operation:   s.getLogCmd,
				},
				"create": {
					name:        "create",
					description: "Create log <id> <entry>",
					operation:   s.createLogCmd,
				},
				"delete": {
					name:        "delete",
					description: "Delete logs for <id>",
					operation:   s.deleteLogsCmd,
				},
			},
		},
	}
}

func (s *CommandLine) Run(cmds map[string]*Command, args []string) {
	s.dispatch(cmds, args)
}

func (s *CommandLine) PrintCommandUsage(cmds map[string]*Command) {
	fmt.Printf("Usage: issuetracker <COMMAND> <SUBCOMMAND> \n")
	s.PrintCommands(cmds, 0)
}

const PRINT_DISTANCE = 15

func (s *CommandLine) PrintCommands(cmds map[string]*Command, depth int) {

	for name, cmd := range cmds {
		// indent based on depth
		for i := 0; i < depth; i++ {
			fmt.Print("  ")
		}
		distance := strings.Repeat(" ", PRINT_DISTANCE-len(name))
		fmt.Println("    " + name + distance + cmd.description)

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

func (s *CommandLine) modifyIssueHelper(req *models.UpdateIssueRequest) error {

	reader := bufio.NewReader(os.Stdin)

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

	return nil
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
	fmt.Printf("External Reference: %s\n", deref(issue.External_Ref, "null"))
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
	created := logs[0].Timestamp

	ext_ref_print := deref(i.External_Ref, "null")

	return fmt.Sprintf(
		"    -%d-     %s  | %s | %s | %s ",
		i.Internal_ID,
		progressSymbol(i.Progress),
		created,
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
