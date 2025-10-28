package models

import (
    "time"
)

type Organization struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
    ID          string    `json:"id"`
    OrgID       string    `json:"org_id"`
    Email       string    `json:"email"`
    DisplayName string    `json:"display_name"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Team struct {
    ID        string    `json:"id"`
    OrgID     string    `json:"org_id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type TeamMember struct {
    TeamID    string    `json:"team_id"`
    UserID    string    `json:"user_id"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type SocialAccount struct {
    ID          string    `json:"id"`
    OrgID       string    `json:"org_id"`
    TeamID      *string   `json:"team_id,omitempty"`
    Platform    string    `json:"platform"`
    Handle      string    `json:"handle"`
    ExternalID  string    `json:"external_id"`
    AuthKind    string    `json:"auth_kind"`
    AuthMeta    any       `json:"auth_meta"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ContentItem struct {
    ID           string    `json:"id"`
    OrgID        string    `json:"org_id"`
    Title        string    `json:"title"`
    Body         string    `json:"body"`
    CanonicalURL string    `json:"canonical_url"`
    MediaURL     string    `json:"media_url"`
    Metadata     any       `json:"metadata"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type Campaign struct {
    ID          string    `json:"id"`
    OrgID       string    `json:"org_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ScheduledPost struct {
    ID              string     `json:"id"`
    OrgID           string     `json:"org_id"`
    CampaignID      *string    `json:"campaign_id,omitempty"`
    ContentID       *string    `json:"content_id,omitempty"`
    SocialAccountID *string    `json:"social_account_id,omitempty"`
    Platform        string     `json:"platform"`
    Caption         *string    `json:"caption,omitempty"`
    Hashtags        *string    `json:"hashtags,omitempty"`
    ScheduledAt     time.Time  `json:"scheduled_at"`
    Status          string     `json:"status"`
    ExternalID      *string    `json:"external_id,omitempty"`
    PublishedAt     *time.Time `json:"published_at,omitempty"`
    Metadata        any        `json:"metadata"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}

type AIGeneration struct {
    ID            string    `json:"id"`
    OrgID         string    `json:"org_id"`
    UserID        *string   `json:"user_id,omitempty"`
    Model         string    `json:"model"`
    Prompt        string    `json:"prompt"`
    Parameters    any       `json:"parameters"`
    OutputText    string    `json:"output_text"`
    OutputJSON    any       `json:"output_json"`
    ContentItemID *string   `json:"content_item_id,omitempty"`
    Status        string    `json:"status"`
    ErrorMessage  string    `json:"error_message"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type AIVariant struct {
    ID            string    `json:"id"`
    GenerationID  string    `json:"generation_id"`
    ContentItemID *string   `json:"content_item_id,omitempty"`
    VariantIndex  int       `json:"variant_index"`
    Payload       any       `json:"payload"`
    CreatedAt     time.Time `json:"created_at"`
}

type Experiment struct {
    ID         string    `json:"id"`
    OrgID      string    `json:"org_id"`
    Name       string    `json:"name"`
    Hypothesis string    `json:"hypothesis"`
    StartedAt  time.Time `json:"started_at"`
    EndedAt    *time.Time `json:"ended_at,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type ExperimentArm struct {
    ID           string    `json:"id"`
    ExperimentID string    `json:"experiment_id"`
    VariantID    *string   `json:"variant_id,omitempty"`
    Weight       float64   `json:"weight"`
    CreatedAt    time.Time `json:"created_at"`
}

type PostOutcome struct {
    ID             string    `json:"id"`
    ScheduledPostID string   `json:"scheduled_post_id"`
    Platform       string    `json:"platform"`
    ExternalID     string    `json:"external_id"`
    Impressions    int64     `json:"impressions"`
    Reach          int64     `json:"reach"`
    Likes          int64     `json:"likes"`
    Comments       int64     `json:"comments"`
    Shares         int64     `json:"shares"`
    Clicks         int64     `json:"clicks"`
    Saves          int64     `json:"saves"`
    Conversions    int64     `json:"conversions"`
    CollectedAt    time.Time `json:"collected_at"`
    Metadata       any       `json:"metadata"`
}

type TrendMetric struct {
    OrgID       string    `json:"org_id"`
    Source      string    `json:"source"`
    Dimension   string    `json:"dimension"`
    Metric      string    `json:"metric"`
    BucketStart time.Time `json:"bucket_start"`
    BucketEnd   time.Time `json:"bucket_end"`
    Value       float64   `json:"value"`
    Meta        any       `json:"meta"`
    CreatedAt   time.Time `json:"created_at"`
}

type AppMetric struct {
    Name       string         `json:"name"`
    Value      float64        `json:"value"`
    Attributes map[string]any `json:"attributes"`
    RecordedAt time.Time      `json:"recorded_at"`
}

type MarketplacePost struct {
    ID            string    `json:"id"`
    OrgID         string    `json:"org_id"`
    SellerUserID  string    `json:"seller_user_id"`
    ContentItemID *string   `json:"content_item_id,omitempty"`
    Title         string    `json:"title"`
    Description   string    `json:"description"`
    PriceCents    int       `json:"price_cents"`
    Currency      string    `json:"currency"`
    Status        string    `json:"status"`
    Metadata      any       `json:"metadata"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type MarketplaceTransaction struct {
    ID           string    `json:"id"`
    PostID       string    `json:"post_id"`
    BuyerUserID  string    `json:"buyer_user_id"`
    AmountCents  int       `json:"amount_cents"`
    Currency     string    `json:"currency"`
    Status       string    `json:"status"`
    CreatedAt    time.Time `json:"created_at"`
}
