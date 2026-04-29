package cli

import (
	"fmt"
	"issuetracker/internal/models"
	"strconv"
)

func (s *CommandLine) setActiveCmd(args []string) {
	s.setStatusCmd(args, true)
}

func (s *CommandLine) setInactiveCmd(args []string) {
	s.setStatusCmd(args, false)
}

func (s *CommandLine) setStatusCmd(args []string, status bool) {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return
	}

	req := models.UpdateIssueRequest{
		Active: &status,
	}

	if err := s.ModifyIssue(id, req); err != nil {
		fmt.Printf("Failed to update issue: %s\n", err)
	}
}

func (s *CommandLine) modifyCmd(args []string) {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return
	}

	err = s.modifyIssueHelper(id)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *CommandLine) getLogCmd(args []string) {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
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

func (s *CommandLine) setProgressCmd(args []string) {
	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return
	}
	if len(args) < 2 {
		fmt.Println("Expected subcommand: idle, started, completed!")
	}

	statusStr := args[1]

	status, err := models.ParseProgressStatus(statusStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	if status.IsValidProgress() {
		err = s.setProgress(id, status)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Not a valid progress status")
		return
	}

}

func (s *CommandLine) deleteLogsCmd(args []string) {

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return
	}

	err = s.DeleteLogsForIssue(id)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// Get one issue
func (s *CommandLine) getCmd(args []string) {

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return
	}

	issue, err := s.GetIssue(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	printIssue(issue)
}

// Get all issues
func (s *CommandLine) listCmd(args []string) {

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
		simplePrintIssue(&issue, s.issueService)
	}
}

func (s *CommandLine) createCmd(args []string) {

	err := s.CreateIssue()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (s *CommandLine) createLogCmd(args []string) {

	id, err := getIdFromInput(args)
	if err != nil || id == -1 {
		return
	}

	if len(args) < 2 {
		fmt.Println("Expected entry as arguments to command!")
		return
	}

	entry := args[1]

	err = s.CreateLogEntry(id, entry)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Log entry created successfully")
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
