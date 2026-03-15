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
