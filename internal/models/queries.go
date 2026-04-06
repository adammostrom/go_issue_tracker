package models

import (
	"fmt"
	"strings"
)

type IssueStatus int

const (
	StatusUknown   IssueStatus = 0
	StatusActive   IssueStatus = 1
	StatusInactive IssueStatus = 2
	StatusDefault  IssueStatus = 3
)

func ParseStatus(s string) (IssueStatus, error) {

	s = cleanString(s)
	switch s {
	case "active":
		return StatusActive, nil
	case "inactive":
		return StatusInactive, nil
	case "":
		return StatusDefault, nil
	default:
		return StatusUknown, fmt.Errorf("invalid status: %s", s)
	}
}

func cleanString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSuffix(s, `\n`)
	s = strings.TrimSpace(s)

	return s
}
