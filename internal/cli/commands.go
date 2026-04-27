package cli

import (
	"fmt"
	"issuetracker/internal/models"
	"strconv"
)

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

func (s *CommandLineInterface) deleteLogsCmd(args []string) {

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
	err = s.DeleteLogsForIssue(id)
	if err != nil {
		fmt.Println(err)
		return
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
		s.simplePrintIssue(&issue)
	}
}

func (s *CommandLineInterface) createCmd(args []string) {

	err := s.CreateIssue()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (s *CommandLineInterface) createLogCmd(args []string) {
	if len(args) < 2 {
		fmt.Println("Expected id and entry as arguments to command!")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("invalid id: %s\n", args[0])
		fmt.Println(err)
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
