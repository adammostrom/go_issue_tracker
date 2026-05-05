package models

// Keep it simple for now
// Capital = Exported (public fields)
// THe internal ID is autoassigned but user can still see it for referals

// 2026-04-24: Adding status
type ProgressStatus int

const (
	Idle     ProgressStatus = 0
	Started  ProgressStatus = 1
	Finished ProgressStatus = 2
)

type Issue struct {
	Internal_ID  int64   // postgres generated
	External_Ref *string // Unique ID put by user
	Title        string
	Description  string
	Log          []LogEntry
	Active       bool
	Progress     ProgressStatus
}

type LogEntry struct {
	Timestamp string // Change to time package later
	Entry     string
}

// Issuerequest currently hides the log
type CreateIssueRequest struct {
	External_Ref *string `json:"external_ref"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
}

// For returning an issue request
type IssueResponse struct {
	Internal_ID  int64          `json:"internal_id"`
	External_Ref string         `json:"external_ref"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Active       bool           `json:"active"`
	Progress     ProgressStatus `json:"progress`
}

// For updating (changing) a single issue. Nil values = field not sent in request, non-nil = fields sent in request, hence they were modified
type UpdateIssueRequest struct {
	External_Ref *string         `json:"external_ref"`
	Title        *string         `json:"title"`
	Description  *string         `json:"description"`
	Active       *bool           `json:"active"`
	Progress     *ProgressStatus `json:"progress"`
}

type IssueFilter struct {
	Active   *bool
	Progress *ProgressStatus
	Created  *string
}
