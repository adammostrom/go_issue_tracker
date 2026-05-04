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

func ptrBool(v bool) *bool {
	return &v
}

func ParseStatus(s string) (*bool, error) {

	s = cleanString(s)

	if s == "" {
		return nil, nil
	}

	switch s {
	case "active":
		return ptrBool(true), nil
	case "inactive":
		return ptrBool(false), nil
	default:
		return nil, fmt.Errorf("invalid status: %s  (valid: active, inactive)", s)
	}
}

func cleanString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSuffix(s, `\n`)
	s = strings.TrimSpace(s)

	return s
}
