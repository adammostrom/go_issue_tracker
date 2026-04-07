package models

import "fmt"

const TITLE_MIN = 2
const TITLE_MAX = 30
const EXTERNAL_MIN = 2
const EXTERNAL_MAX = 20
const DESCR_MIN = 0
const DESCR_MAX = 100

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
func ValidateExternalRef(extRef string) error {
	if len(extRef) < EXTERNAL_MIN || len(extRef) > EXTERNAL_MAX {
		return fmt.Errorf("External Reference must be between %d and %d . Provided was: %s\n", EXTERNAL_MIN, EXTERNAL_MAX, extRef)
	}
	return nil
}

func ValidateTitle(title string) error {
	if len(title) < TITLE_MIN || len(title) > TITLE_MAX {
		return fmt.Errorf("Title must be between %d and %d \n", TITLE_MIN, TITLE_MAX)
	}
	return nil
}

func ValidateDescription(description string) error {
	if len(description) < DESCR_MIN || len(description) > DESCR_MAX {
		return fmt.Errorf("Description must be between %d and %d \n", DESCR_MIN, DESCR_MAX)
	}
	return nil
}

func (l *LogEntry) ValidateEntry() error {
	if len(l.Entry) < DESCR_MIN || len(l.Entry) > DESCR_MAX {
		return fmt.Errorf("A log entry must be between %d and %d\n", DESCR_MIN, DESCR_MAX)
	}
	return nil
}
