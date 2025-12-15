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
	Type  string     `json:"type"`
	Title []richText `json:"title"`
	// include other fields as needed over time
}

type richText struct {
	PlainText string `json:"plain_text"`
}
