package fal

// Response represents the Fal queue response payload.
type Response struct {
	Status          string  `json:"status,omitempty"`
	RequestID       string  `json:"request_id,omitempty"`
	ResponseURL     string  `json:"response_url,omitempty"`
	StatusURL       string  `json:"status_url,omitempty"`
	CancelURL       string  `json:"cancel_url,omitempty"`
	Images          []Image `json:"images,omitempty"`
	Timings         Timings `json:"timings,omitempty"`
	Seed            int     `json:"seed,omitempty"`
	HasNSFWConcepts []bool  `json:"has_nsfw_concepts,omitempty"`
	Prompt          string  `json:"prompt,omitempty"`
	QueuePosition   int     `json:"queue_position,omitempty"`
}

// Image is an individual generated asset from Fal.
type Image struct {
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"content_type"`
}

// Timings captures Fal latency stats.
type Timings struct {
	Inference float64 `json:"inference"`
}
