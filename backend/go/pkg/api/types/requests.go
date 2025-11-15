package types

import (
	"time"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CreateContentRequest represents the request to create new content
type CreateContentRequest struct {
	Title       string                 `json:"title" validate:"required"`
	Content     string                 `json:"content"`
	ContentType string                 `json:"content_type" validate:"required,oneof=article video podcast"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateContentRequest represents the request to update content
type UpdateContentRequest struct {
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Status   string                 `json:"status"`
	Metadata map[string]interface{} `json:"metadata"`
}

// PublishContentRequest represents the request to publish content
type PublishContentRequest struct {
	PublishAt *time.Time `json:"publish_at,omitempty"`
}

// CreateSubscriptionRequest represents the request to create a subscription
type CreateSubscriptionRequest struct {
	PlanID       string `json:"plan_id" validate:"required"`
	PaymentToken string `json:"payment_token" validate:"required"`
}

// UpdateSubscriptionRequest represents the request to update a subscription
type UpdateSubscriptionRequest struct {
	PlanID string `json:"plan_id" validate:"required"`
}

// CancelSubscriptionRequest represents the request to cancel a subscription
type CancelSubscriptionRequest struct {
	Feedback string `json:"feedback,omitempty"`
}

// CreateTeamRequest represents the request to create a team
type CreateTeamRequest struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description"`
	MemberEmails []string `json:"member_emails"`
}

// InviteUserRequest represents the request to invite a user
type InviteUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member viewer"`
}

// UpdateUserRoleRequest represents the request to update a user's role
type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer"`
}

// SearchRequest represents a search query
type SearchRequest struct {
	Query   string   `json:"query"`
	Filters []string `json:"filters,omitempty"`
	Limit   int      `json:"limit,omitempty"`
	Offset  int      `json:"offset,omitempty"`
}

// AnalyticsQuery represents a query for analytics data
type AnalyticsQuery struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Metrics    []string  `json:"metrics"`
	Dimensions []string  `json:"dimensions,omitempty"`
	Filter     string    `json:"filter,omitempty"`
}

// WebhookRequest represents an incoming webhook
type WebhookRequest struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}
