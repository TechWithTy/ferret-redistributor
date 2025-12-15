package notion

import (
	"strings"
	"time"
)

// NextRun computes the next scheduled run time given a previous run time and recurrence name.
// The recurrence names match the select options we added in Notion.
func NextRun(prev time.Time, recurrence string) (time.Time, bool) {
	r := strings.TrimSpace(recurrence)
	if prev.IsZero() {
		return time.Time{}, false
	}
	switch r {
	case "", "None":
		return time.Time{}, false
	case "Daily":
		return prev.Add(24 * time.Hour), true
	case "Every 7 days":
		return prev.Add(7 * 24 * time.Hour), true
	case "Every 2 weeks":
		return prev.Add(14 * 24 * time.Hour), true
	case "Monthly":
		return prev.AddDate(0, 1, 0), true
	default:
		return time.Time{}, false
	}
}


