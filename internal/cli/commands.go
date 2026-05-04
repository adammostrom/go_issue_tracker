package cli

import (
	"bufio"
	"fmt"
	"issuetracker/internal/models"
	"os"
	"strconv"
)

func (s *CommandLine) setActiveCmd(args []string) error {
	return s.setStatusCmd(args, true)
}

func (s *CommandLine) setInactiveCmd(args []string) error {
	return s.setStatusCmd(args, false)
}

func (s *CommandLine) setStatusCmd(args []string, status bool) error {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}

	req := models.UpdateIssueRequest{
		Active: &status,
	}

	if err := s.issueService.PatchIssue(id, req); err != nil {
		fmt.Printf("Failed to update issue with id: %d\n", id)
		return err
	}

	fmt.Printf("Status updated successfully to: %t\n", status)
	return nil
}

func (s *CommandLine) modifyCmd(args []string) error {

	var req = &models.UpdateIssueRequest{}

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}

	err = s.modifyIssueHelper(req)
	if err != nil {
		return err
	}

	err = s.issueService.PatchIssue(id, *req)
	fmt.Printf("Issue modified successfully with values: %v \n", *req)
	return nil
}

func (s *CommandLine) deleteCmd(args []string) error {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}
	err = s.issueService.DeleteIssue(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Issue deleted successfully")
	return nil
}

func (s *CommandLine) getLogCmd(args []string) error {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}

	logs, err := s.issueService.GetLogsFromIssue(id)
	if err != nil {
		return err
	}
	if len(logs) == 0 {
		fmt.Printf("No log entries for given id: %d\n", id)
		return nil
	}
	for _, log := range logs {
		fmt.Printf("entry: %v\n", log)
	}
	return nil
}

func (s *CommandLine) setProgressCmd(args []string) error {

	var req = &models.UpdateIssueRequest{}

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}
	if len(args) < 2 {
		fmt.Println("Expected subcommand: idle, started, completed!")
		return nil
	}

	statusStr := args[1]

	status, err := models.ParseProgressStatus(statusStr)
	if err != nil {
		return err
	}

	if status.IsValidProgress() {
		req.Progress = &status
		err = s.issueService.PatchIssue(id, *req)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Not a valid progress status")
		return err
	}

	fmt.Printf("Progress status changed successfully to: %s\n", req.Progress.String())
	return nil

}

func (s *CommandLine) deleteLogsCmd(args []string) error {

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}
	err = s.issueService.DeleteLogsFromIssue(id)
	if err != nil {
		return err
	}
	fmt.Println("Logs erased successfully")
	return nil
}

// Get one issue
func (s *CommandLine) getIssueCmd(args []string) error {

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}

	issue, err := s.issueService.GetIssueByID(id)
	if err != nil || issue == nil {
		return err
	}
	logs, err := s.issueService.GetLogsFromIssue(id)
	if err != nil {
		return err
	}
	printIssue(issue, logs)

	return nil
}

func parseArguments(args []string) (models.IssueFilter, error) {
	var f models.IssueFilter

	// with no given, return all issues
	if len(args) < 1 {
		return f, nil
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "status":
			if 1+i >= len(args) {
				return f, fmt.Errorf("status requires a value")
			}
			active_status, err := models.ParseStatus(args[i+1])
			if err != nil {
				return f, err
			}
			f.Active = active_status
			i++

		case "progress":
			if 1+i >= len(args) {
				return f, fmt.Errorf("progress requires a value")
			}
			progress, err := models.ParseProgressStatus(args[i+1])
			if err != nil {
				return f, err
			}
			f.Progress = &progress
			i++
		default:
			return f, fmt.Errorf("unknown argument: %s", args[i])
		}
	}
	return f, nil
}

// Get all issues, regardless of status or progress
func (s *CommandLine) getAllCmd(args []string) error {

	//status := models.StatusDefault

	filter, err := parseArguments(args)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	issues, err := s.issueService.GetAllIssues(filter)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		fmt.Println("No issues found")
		return nil
	}
	simplePrintIssue(issues, s.issueService)
	return nil
}

func (s *CommandLine) createCmd(args []string) error {

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
	printIssue(issue, nil)
	return nil

}

func (s *CommandLine) createLogCmd(args []string) error {

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return err
	}

	if len(args) < 2 {
		fmt.Println("Expected entry as arguments to command!")
		return err
	}

	entry := args[1]

	err = s.issueService.AddLogEntry(id, entry)
	if err != nil {
		fmt.Printf("Failed to create log entry: %s", err)
		return err
	}

	fmt.Println("Log entry created successfully")
	return nil
}

func getIdFromInput(args []string) (int, error) {
	if len(args) < 1 {
		fmt.Println("Expected id as argument to command!")
		return -1, nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("invalid id: %s\n", args[0])
		fmt.Println(err)
		return -1, err
	}
	return id, nil
}
