package models

// Keep it simple for now
// Capital = Exported (public fields)
// THe internal ID is autoassigned but user can still see it for referals

type Issue struct {
	Internal_id int // Private
	Name        string
	Description string
	Log         []LogEntry
	Resolved    bool
}

type LogEntry struct {
	Timestamp string // Change to time package later
	Entry     string
}
