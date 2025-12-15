package notion

// Minimal response shapes we need.

type pageRef struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type queryResponse struct {
	Results []pageRef `json:"results"`
}

type queryResponseFull struct {
	Results    []pageObject `json:"results"`
	HasMore    bool         `json:"has_more"`
	NextCursor string       `json:"next_cursor"`
}

type pageObject struct {
	ID         string                   `json:"id"`
	URL        string                   `json:"url"`
	Properties map[string]propertyValue `json:"properties"`
}

type propertyValue struct {
	Type string `json:"type"`

	Title    []richText    `json:"title,omitempty"`
	RichText []richText    `json:"rich_text,omitempty"`
	Number   *float64      `json:"number,omitempty"`
	Checkbox *bool         `json:"checkbox,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Select   *selectValue  `json:"select,omitempty"`
	Date     *dateValue    `json:"date,omitempty"`
	Relation []relationRef `json:"relation,omitempty"`
	// include other fields as needed over time
}

type richText struct {
	PlainText string `json:"plain_text"`
}

type selectValue struct {
	Name string `json:"name"`
}

type dateValue struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type relationRef struct {
	ID string `json:"id"`
}
