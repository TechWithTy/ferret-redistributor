package models

import (
	"time"

	"github.com/google/uuid"
)

// AnalyticsEvent represents a user interaction or system event
type AnalyticsEvent struct {
	ID           string    `json:"id"`
	EventType    string    `json:"event_type"`  // view, click, share, etc.
	EventValue   string    `json:"event_value"` // URL, button ID, etc.
	SessionID    string    `json:"session_id"`
	UserID       string    `json:"user_id,omitempty"`
	AnonymousID  string    `json:"anonymous_id"`
	Referrer     string    `json:"referrer,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	ScreenSize   string    `json:"screen_size,omitempty"`
	Metadata     JSONMap   `json:"metadata,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ContentAnalytics aggregates analytics for content pieces
type ContentAnalytics struct {
	ContentID     string    `json:"content_id"`
	ViewCount     int64     `json:"view_count"`
	UniqueViews   int64     `json:"unique_views"`
	AvgTimeOnPage float64   `json:"avg_time_on_page"`
	BounceRate    float64   `json:"bounce_rate"`
	Shares        int64     `json:"shares"`
	Likes         int64     `json:"likes"`
	Comments      int64     `json:"comments"`
	Date          time.Time `json:"date"`
}

// UserEngagement tracks user engagement metrics
type UserEngagement struct {
	UserID           string    `json:"user_id"`
	SessionCount     int64     `json:"session_count"`
	AvgSessionTime   float64   `json:"avg_session_time"`
	PagesPerSession  float64   `json:"pages_per_session"`
	LastActiveAt     time.Time `json:"last_active_at"`
	TotalInteractions int64     `json:"total_interactions"`
}

// Funnel represents a conversion funnel
type Funnel struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Steps       []FunnelStep `json:"steps"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// FunnelStep represents a step in a conversion funnel
type FunnelStep struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Order    int    `json:"order"`
	Required bool   `json:"required"`
}

// FunnelResult represents the results of a funnel analysis
type FunnelResult struct {
	FunnelID     string             `json:"funnel_id"`
	Steps        []FunnelStepResult `json:"steps"`
	StartDate    time.Time         `json:"start_date"`
	EndDate      time.Time         `json:"end_date"`
	TotalUsers   int64             `json:"total_users"`
	ConversionRate float64         `json:"conversion_rate"`
}

type FunnelStepResult struct {
	StepName    string  `json:"step_name"`
	Visitors    int64   `json:"visitors"`
	DropOff     int64   `json:"drop_off"`
	DropOffRate float64 `json:"drop_off_rate"`
}

// NewAnalyticsEvent creates a new analytics event
func NewAnalyticsEvent(eventType, eventValue, sessionID string) *AnalyticsEvent {
	return &AnalyticsEvent{
		ID:        uuid.New().String(),
		EventType: eventType,
		EventValue: eventValue,
		SessionID: sessionID,
		CreatedAt: time.Now().UTC(),
	}
}

// TrackView tracks a view event for content
func TrackView(contentID, userID, sessionID string) *AnalyticsEvent {
	event := NewAnalyticsEvent("view", contentID, sessionID)
	event.UserID = userID
	return event
}

// TrackClick tracks a click event
func TrackClick(elementID, userID, sessionID string) *AnalyticsEvent {
	return NewAnalyticsEvent("click", elementID, sessionID)
}

// TrackConversion tracks a conversion event
func TrackConversion(goalName, userID, sessionID string, value float64) *AnalyticsEvent {
	event := NewAnalyticsEvent("conversion", goalName, sessionID)
	event.UserID = userID
	if event.Metadata == nil {
		event.Metadata = make(JSONMap)
	}
	event.Metadata["value"] = value
	return event
}
