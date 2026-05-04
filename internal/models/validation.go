package models

import (
	"fmt"
	"strings"
)

const TITLE_MIN = 2
const TITLE_MAX = 40
const EXTERNAL_MIN = 0
const EXTERNAL_MAX = 10
const DESCR_MIN = 0
const DESCR_MAX = 200

func (i *Issue) ValidateIssue() error {
	if err := ValidateExternalRef(i.External_Ref); err != nil {
		return err
	}
	if err := ValidateTitle(i.Title); err != nil {
		return err
	}
	if err := ValidateDescription(i.Description); err != nil {
		return err
	}
	return nil
}

func ValidateExternalRefWrapper(s string) error {
	if s == "" {
		return nil // or treat as optional
	}

	// *value = get the value at the address
	// &value = get the address from the value

	return ValidateExternalRef(&s)
}
func ValidateExternalRef(extRef *string) error {
	if extRef == nil {
		return nil
	}
	if len(*extRef) < EXTERNAL_MIN || len(*extRef) > EXTERNAL_MAX {
		return fmt.Errorf("External Reference must be between %d and %d . Provided was: %d\n", EXTERNAL_MIN, EXTERNAL_MAX, len(*extRef))
	}
	return nil
}

func ValidateTitle(title string) error {
	if len(title) < TITLE_MIN || len(title) > TITLE_MAX {
		return fmt.Errorf("Title must be between %d and %d. Provided was: %d \n", TITLE_MIN, TITLE_MAX, len(title))
	}
	return nil
}

func ValidateDescription(description string) error {
	if len(description) < DESCR_MIN || len(description) > DESCR_MAX {
		return fmt.Errorf("Description must be between %d and %d. Provided was: %d \n", DESCR_MIN, DESCR_MAX, len(description))
	}
	return nil
}

func (l *LogEntry) ValidateEntry() error {
	if len(l.Entry) < DESCR_MIN || len(l.Entry) > DESCR_MAX {
		return fmt.Errorf("A log entry must be between %d and %d. Provided was: %d \n", DESCR_MIN, DESCR_MAX, len(l.Entry))
	}
	return nil
}

func (p *ProgressStatus) IsValidProgress() bool {
	switch *p {
	case Idle:
		return true
	case Started:
		return true
	case Finished:
		return true
	default:
		return false
	}
}

func (p *ProgressStatus) String() string {
	switch *p {
	case Idle:
		return "idle"
	case Started:
		return "started"
	case Finished:
		return "finished"
	default:
		return "unknown"
	}
}

func ParseProgressStatus(s string) (ProgressStatus, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "idle":
		return Idle, nil
	case "started":
		return Started, nil
	case "completed":
		return Finished, nil
	default:
		return -1, fmt.Errorf("invalid progress: %s (valid: idle, started, completed)", s)
	}
}
