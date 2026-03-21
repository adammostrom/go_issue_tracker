package models

// Keep it simple for now
// Capital = Exported (public fields)
// THe internal ID is autoassigned but user can still see it for referals

type Issue struct {
	Internal_id  int64 // postgres generated
	External_Ref int64 // Unique ID put by user
	Title        string
	Description  string
	Log          []LogEntry
	Active       bool
}

type LogEntry struct {
	Timestamp string // Change to time package later
	Entry     string
}

// Issuerequest currently hides the log
type CreateIssueRequest struct {
	ExternalRef int64  `json:"external_ref"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// For returning an issue request
type IssueResponse struct {
	InternalID  int64  `json:"internal_id"`
	ExternalRef int64  `json:"external_ref"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// For updating (changing) a single issue. Nil values = field not sent in request, non-nil = fields sent in request, hence they were modified
type UpdateIssueRequest struct {
	ExternalRef *int64  `json:"external_ref"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Active      *bool   `json:"active"`
}
