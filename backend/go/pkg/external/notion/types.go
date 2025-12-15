package notion

// Minimal response shapes we need.

type pageRef struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type queryResponse struct {
	Results []pageRef `json:"results"`
}
