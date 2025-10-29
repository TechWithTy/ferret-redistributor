package postiz

// Integration/channel
type Integration struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Picture    string `json:"picture"`
	Disabled   bool   `json:"disabled"`
	Profile    string `json:"profile"`
	Customer   *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"customer,omitempty"`
}

type FindSlotResponse struct {
	Date string `json:"date"`
}

type UploadResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	OrganizationID string `json:"organizationId"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type PostsListResponse struct {
	Posts []Post `json:"posts"`
	Meta  *Meta  `json:"meta,omitempty"`
}

type Post struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	PublishDate string `json:"publishDate"`
	ReleaseURL  string `json:"releaseURL"`
	State       string `json:"state"`
	Integration struct {
		ID                 string `json:"id"`
		ProviderIdentifier string `json:"providerIdentifier"`
		Name               string `json:"name"`
		Picture            string `json:"picture"`
	} `json:"integration"`
}

type CreateUpdateResult struct {
	PostID      string `json:"postId"`
	Integration string `json:"integration"`
}

type DeletePostResponse struct {
	ID string `json:"id"`
}

type GeneratedVideo struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type VideoFunctionResponse struct {
	Voices []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"voices"`
}

// Meta describes optional pagination or rate-limit info if API returns it.
type Meta struct {
	NextPage   *int   `json:"nextPage,omitempty"`
	Limit      *int   `json:"limit,omitempty"`
	Remaining  *int   `json:"remaining,omitempty"`
	ResetAfter *int64 `json:"resetAfter,omitempty"`
}
