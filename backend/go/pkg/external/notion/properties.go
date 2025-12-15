package notion

import "time"

// Property builders for Notion page create/update payloads.
// These match the Notion Public API JSON shapes.

func Title(content string) map[string]any {
	return map[string]any{
		"title": []map[string]any{
			{
				"type": "text",
				"text": map[string]any{
					"content": content,
				},
			},
		},
	}
}

func RichText(content string) map[string]any {
	return map[string]any{
		"rich_text": []map[string]any{
			{
				"type": "text",
				"text": map[string]any{
					"content": content,
				},
			},
		},
	}
}

func Select(name string) map[string]any {
	return map[string]any{
		"select": map[string]any{
			"name": name,
		},
	}
}

func URL(u string) map[string]any {
	return map[string]any{"url": u}
}

func Checkbox(v bool) map[string]any {
	return map[string]any{"checkbox": v}
}

func Number(v float64) map[string]any {
	return map[string]any{"number": v}
}

func DateTime(t time.Time) map[string]any {
	if t.IsZero() {
		return map[string]any{"date": nil}
	}
	return map[string]any{
		"date": map[string]any{
			"start": t.UTC().Format(time.RFC3339),
		},
	}
}

func Relation(pageIDs ...string) map[string]any {
	rels := make([]map[string]any, 0, len(pageIDs))
	for _, id := range pageIDs {
		if id == "" {
			continue
		}
		rels = append(rels, map[string]any{"id": id})
	}
	return map[string]any{"relation": rels}
}
