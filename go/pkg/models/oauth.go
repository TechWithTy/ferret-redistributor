package models

import (
	"time"

)

// OAuthProvider represents an OAuth provider configuration
type OAuthProvider struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Provider     string    `json:"provider" gorm:"uniqueIndex:idx_oauth_provider_name"` // e.g., "google", "facebook", "github"
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"-"` // Never expose in API responses
	AuthURL      string    `json:"auth_url"`
	TokenURL     string    `json:"token_url"`
	UserInfoURL  string    `json:"user_info_url"`
	Scopes       []string  `json:"scopes" gorm:"type:text[]"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OAuthAccount links a user account to an OAuth provider
type OAuthAccount struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Provider       string    `json:"provider"` // Matches OAuthProvider.Provider
	ProviderUserID string    `json:"provider_user_id" gorm:"index"`
	AccessToken    string    `json:"-" gorm:"type:text"`
	RefreshToken   string    `json:"-" gorm:"type:text"`
	TokenType      string    `json:"token_type"`
	ExpiresAt      time.Time `json:"expires_at"`
	Scope          string    `json:"scope"`
	Email          string    `json:"email"`
	ProfileData    JSONB     `json:"profile_data" gorm:"type:jsonb"` // Raw profile data from provider
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// PhoneVerification stores verification codes for phone numbers
type PhoneVerification struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PhoneNumber string    `json:"phone_number" gorm:"index"`
	Code        string    `json:"-"` // The verification code (hashed)
	Attempts    int       `json:"attempts" gorm:"default:0"`
	ExpiresAt   time.Time `json:"expires_at"`
	Verified    bool      `json:"verified" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	VerifiedAt  time.Time `json:"verified_at,omitempty"`
}

// SocialAccountLink represents a link between a user's account and their social media accounts
type SocialAccountLink struct {
	ID                string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID            string     `json:"user_id" gorm:"type:uuid;not null;index"`
	Provider          string     `json:"provider"` // e.g., "twitter", "linkedin", "instagram"
	ExternalID        string     `json:"external_id" gorm:"index"`
	Username          string     `json:"username"`
	DisplayName       string     `json:"display_name"`
	ProfileURL        string     `json:"profile_url"`
	AvatarURL         string     `json:"avatar_url"`
	AccessToken       string     `json:"-"` // Encrypted
	AccessTokenSecret string     `json:"-"` // For OAuth 1.0a
	RefreshToken      string     `json:"-"` // For OAuth 2.0
	TokenExpiry       *time.Time `json:"token_expiry,omitempty"`
	IsPrimary         bool       `json:"is_primary" gorm:"default:false"`
	Metadata          JSONB      `json:"metadata" gorm:"type:jsonb"` // Additional provider-specific data
	LastSyncedAt      *time.Time `json:"last_synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// JSONB is a helper type for JSONB fields
type JSONB map[string]interface{}
