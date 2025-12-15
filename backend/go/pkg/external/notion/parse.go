package notion

import (
	"strings"
	"time"
)

func TitlePlainText(props map[string]propertyValue, titlePropName string) string {
	if props == nil {
		return ""
	}
	p, ok := props[titlePropName]
	if !ok {
		return ""
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) != "title" {
		return ""
	}
	var b strings.Builder
	for _, rt := range p.Title {
		if strings.TrimSpace(rt.PlainText) == "" {
			continue
		}
		b.WriteString(rt.PlainText)
	}
	return strings.TrimSpace(b.String())
}

func SelectName(props map[string]propertyValue, propName string) string {
	if props == nil {
		return ""
	}
	p, ok := props[propName]
	if !ok {
		return ""
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) != "select" {
		return ""
	}
	if p.Select == nil {
		return ""
	}
	return strings.TrimSpace(p.Select.Name)
}

func NumberValue(props map[string]propertyValue, propName string) (float64, bool) {
	if props == nil {
		return 0, false
	}
	p, ok := props[propName]
	if !ok {
		return 0, false
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) != "number" {
		return 0, false
	}
	if p.Number == nil {
		return 0, false
	}
	return *p.Number, true
}

func DateStartTime(props map[string]propertyValue, propName string) (time.Time, bool) {
	if props == nil {
		return time.Time{}, false
	}
	p, ok := props[propName]
	if !ok {
		return time.Time{}, false
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) != "date" {
		return time.Time{}, false
	}
	if p.Date == nil || strings.TrimSpace(p.Date.Start) == "" {
		return time.Time{}, false
	}
	// Notion can return date-only (YYYY-MM-DD) or RFC3339 datetime.
	s := strings.TrimSpace(p.Date.Start)
	if len(s) == len("2006-01-02") {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func RelationIDs(props map[string]propertyValue, propName string) []string {
	if props == nil {
		return nil
	}
	p, ok := props[propName]
	if !ok {
		return nil
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) != "relation" {
		return nil
	}
	out := make([]string, 0, len(p.Relation))
	for _, r := range p.Relation {
		if strings.TrimSpace(r.ID) == "" {
			continue
		}
		out = append(out, r.ID)
	}
	return out
}


