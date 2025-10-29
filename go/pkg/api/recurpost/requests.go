package recurpost

import "time"

// Common
type Pagination struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// Auth
type TokenRequest struct {
	GrantType    string `json:"grant_type"`              // password | refresh_token
	Username     string `json:"username,omitempty"`      // for password
	Password     string `json:"password,omitempty"`      // for password
	RefreshToken string `json:"refresh_token,omitempty"` // for refresh_token
}

// Posts
type CreatePostRequest struct {
	Text        string     `json:"text"`
	MediaIDs    []string   `json:"media_ids,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	AccountIDs  []string   `json:"account_ids,omitempty"` // profiles to publish to
}

type UpdatePostRequest struct {
	Text        *string    `json:"text,omitempty"`
	MediaIDs    []string   `json:"media_ids,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	Status      *string    `json:"status,omitempty"` // draft|scheduled|...
}

type ListPostsRequest struct {
	Pagination
	Status string `json:"status,omitempty"` // draft|scheduled|published|failed
}

// Media
type UploadMediaRequest struct {
    Filename string `json:"filename"`
    // Provide one of URL or Bytes (as base64 externally when wiring HTTP)
    URL   string `json:"url,omitempty"`
    Bytes []byte `json:"-"`
}

// User login
type UserLoginRequest struct {
    EmailID string `json:"emailid"`
    PassKey string `json:"pass_key"`
}

// Connect Social Account URLs
type ConnectSocialAccountURLsRequest struct {
    EmailID string `json:"emailid"`
    PassKey string `json:"pass_key"`
}

// Social account list
type SocialAccountListRequest struct {
    EmailID string `json:"emailid"`
}

// Library list
type LibraryListRequest struct {
    EmailID string `json:"emailid"`
    PassKey string `json:"pass_key"`
}

// History data
type HistoryDataRequest struct {
    EmailID            string `json:"emailid"`
    PassKey            string `json:"pass_key"`
    ID                 string `json:"id,omitempty"` // social account id
    IsGetVideoUpdates  string `json:"is_get_video_updates,omitempty"` // "true" or "false"
}

// Add content in library
type AddContentInLibraryRequest struct {
    EmailID    string `json:"emailid"`
    PassKey    string `json:"pass_key"`
    ID         string `json:"id"` // library id
    Message    string `json:"message"`
    FBMessage  string `json:"fb_message,omitempty"`
    TWMessage  string `json:"tw_message,omitempty"`
    LNMessage  string `json:"ln_message,omitempty"`
}
