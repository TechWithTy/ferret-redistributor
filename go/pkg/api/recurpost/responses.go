package recurpost

import "time"

// Auth
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Accounts
type Account struct {
	ID          string `json:"id"`
	Provider    string `json:"provider"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
}

// Posts
type Post struct {
	ID          string     `json:"id"`
	Text        string     `json:"text"`
	MediaIDs    []string   `json:"media_ids,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PostList struct {
	Items    []Post `json:"items"`
	NextPage *int   `json:"next_page,omitempty"`
}

// Media
type Media struct {
    ID       string `json:"id"`
    URL      string `json:"url"`
    MimeType string `json:"mime_type"`
    Width    int    `json:"width,omitempty"`
    Height   int    `json:"height,omitempty"`
}

// User login
type UserLoginResponse struct {
    Success     bool   `json:"success"`
    Message     string `json:"message,omitempty"`
    AccessToken string `json:"access_token,omitempty"`
}

// Connect Social Account URLs
type ConnectURLsResponse struct {
    URLs map[string]string `json:"urls"`
}

// Social account list
type SocialAccountListResponse struct {
    Accounts []Account `json:"accounts"`
}

// Library list
type Library struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}

type LibraryListResponse struct {
    Libraries []Library `json:"libraries"`
}

// History data (generic map to allow unknown shape)
type HistoryDataResponse struct {
    Data map[string]any `json:"data"`
}

// Add content in library
type AddContentInLibraryResponse struct {
    Success   bool   `json:"success"`
    Message   string `json:"message,omitempty"`
    ContentID string `json:"content_id,omitempty"`
}
