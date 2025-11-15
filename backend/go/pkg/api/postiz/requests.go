package postiz

// Upload from URL
type UploadFromURLRequest struct {
	URL string `json:"url"`
}

// Posts list
type PostsListRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Customer  string `json:"customer,omitempty"`
}

// Create/Update posts
type CreateUpdatePostRequest struct {
	Type  string      `json:"type"` // draft|schedule|now
	Date  string      `json:"date,omitempty"`
	Posts []PostInput `json:"posts,omitempty"`
}

type PostInput struct {
	Integration IntegrationInput `json:"integration"`
	Value       []PostContent    `json:"value"`
	Group       string           `json:"group,omitempty"`
	Settings    map[string]any   `json:"settings,omitempty"`
}

type IntegrationInput struct {
	ID string `json:"id"`
}

type PostContent struct {
	Content string     `json:"content"`
	ID      string     `json:"id,omitempty"`
	Image   []MediaDto `json:"image,omitempty"`
}

type MediaDto struct {
	ID   string `json:"id,omitempty"`
	Path string `json:"path,omitempty"`
}

// Generate video
type GenerateVideoRequest struct {
	Type         string         `json:"type"`   // image-text-slides|veo3
	Output       string         `json:"output"` // vertical|horizontal
	CustomParams map[string]any `json:"customParams"`
}

// Video function (e.g., loadVoices)
type VideoFunctionRequest struct {
	FunctionName string `json:"functionName"`
	Identifier   string `json:"identifier"`
}
