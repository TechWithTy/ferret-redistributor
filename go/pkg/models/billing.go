package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BillingPlan represents a subscription plan
const (
	PlanFree       = "free"
	PlanStarter    = "starter"
	PlanGrowth     = "growth"
	PlanEnterprise = "enterprise"
)

// BillingStatus represents the status of a subscription
const (
	BillingStatusActive   = "active"
	BillingStatusPastDue  = "past_due"
	BillingStatusCanceled = "canceled"
	BillingStatusTrial    = "trial"
)

// Subscription represents a user's subscription
// @model Subscription
// @description User subscription information
// @json {
//   "type": "object",
//   "required": ["id", "org_id", "plan", "status", "current_period_end"]
// }
type Subscription struct {
	ID                 string     `json:"id"`
	OrgID              string     `json:"org_id"`
	Plan               string     `json:"plan"`
	Status             string     `json:"status"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end"`
	TrialEndsAt        *time.Time `json:	rial_ends_at,omitempty"`
	BillingEmail       string     `json:"billing_email"`
	BillingName        string     `json:"billing_name"`
	BillingAddress     string     `json:"billing_address"`
	BillingCity        string     `json:"billing_city"`
	BillingState       string     `json:"billing_state"`
	BillingPostalCode  string     `json:"billing_postal_code"`
	BillingCountry     string     `json:"billing_country"`
	TaxID              string     `json:"tax_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Credit represents a user's credit balance
// @model Credit
// @description User credit information
// @json {
//   "type": "object",
//   "required": ["id", "org_id", "balance", "currency"]
// }
type Credit struct {
	ID            string    `json:"id"`
	OrgID         string    `json:"org_id"`
	Balance       int64     `json:"balance"` // In smallest currency unit (e.g., cents)
	Currency      string    `json:"currency"` // ISO 4217 currency code
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// CreditTransaction represents a credit transaction
// @model CreditTransaction
// @description Credit transaction record
// @json {
//   "type": "object",
//   "required": ["id", "org_id", "amount", "type", "status"]
// }
type CreditTransaction struct {
	ID            string    `json:"id"`
	OrgID         string    `json:"org_id"`
	Amount        int64     `json:"amount"` // Positive for credits, negative for debits
	Type          string    `json:	ype"`    // "purchase", "usage", "refund", "adjustment"
	Status        string    `json:"status"`  // "pending", "completed", "failed", "refunded"
	Description   string    `json:"description"`
	ReferenceID   string    `json:"reference_id,omitempty"` // External reference ID
	Metadata      JSONMap   `json:"metadata,omitempty"`
	ProcessedAt   time.Time `json:"processed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// BillingEvent represents a billing-related event
// @model BillingEvent
// @description Billing event log
// @json {
//   "type": "object",
//   "required": ["id", "org_id", "type", "status"]
// }
type BillingEvent struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Type        string    `json:"type"`    // "invoice.paid", "invoice.payment_failed", etc.
	Status      string    `json:"status"`  // "pending", "processed", "failed"
	Amount      int64     `json:"amount"`  // In smallest currency unit
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Metadata    JSONMap   `json:"metadata,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// SSLCertificate represents an SSL certificate
// @model SSLCertificate
// @description SSL certificate information
// @json {
//   "type": "object",
//   "required": ["id", "org_id", "domain", "status", "expires_at"]
// }
type SSLCertificate struct {
	ID              string     `json:"id"`
	OrgID           string     `json:"org_id"`
	Domain          string     `json:"domain"`
	Status          string     `json:"status"` // "active", "pending", "expired", "revoked"
	Certificate     string     `json:"certificate,omitempty"`
	PrivateKey      string     `json:"private_key,omitempty"` // Should be encrypted at rest
	Issuer          string     `json:"issuer"`
	StartsAt        time.Time  `json:"starts_at"`
	ExpiresAt       time.Time  `json:"expires_at"`
	AutoRenew       bool       `json:"auto_renew"`
	LastRenewedAt   *time.Time `json:"last_renewed_at,omitempty"`
	ChallengeType   string     `json:"challenge_type"` // "http-01", "dns-01"
	ChallengeToken  string     `json:"challenge_token,omitempty"`
	VerificationURL string     `json:"verification_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// JSONMap is a helper type for JSON fields
type JSONMap map[string]interface{}

// NewSubscription creates a new subscription with default values
func NewSubscription(orgID, plan string) *Subscription {
	now := time.Now().UTC()
	return &Subscription{
		ID:                uuid.New().String(),
		OrgID:             orgID,
		Plan:              plan,
		Status:            BillingStatusTrial,
		CurrentPeriodEnd:  now.AddDate(0, 1, 0), // 1 month trial by default
		CancelAtPeriodEnd: false,
		TrialEndsAt:       nil,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// NewCredit initializes a new credit balance
func NewCredit(orgID string) *Credit {
	return &Credit{
		ID:            uuid.New().String(),
		OrgID:         orgID,
		Balance:       0,
		Currency:      "USD",
		LastUpdatedAt: time.Now().UTC(),
	}
}

// AddCredit adds credit to the balance
func (c *Credit) AddCredit(amount int64, description string) *CreditTransaction {
	if amount <= 0 {
		return nil
	}

	c.Balance += amount
	c.LastUpdatedAt = time.Now().UTC()

	return &CreditTransaction{
		ID:          uuid.New().String(),
		OrgID:       c.OrgID,
		Amount:      amount,
		Type:        "purchase",
		Status:      "completed",
		Description: description,
		ProcessedAt: time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
	}
}

// DeductCredit deducts credit from the balance
func (c *Credit) DeductCredit(amount int64, description string) (*CreditTransaction, error) {
	if amount <= 0 {
		return nil, nil
	}

	if c.Balance < amount {
		return nil, ErrInsufficientCredits
	}

	c.Balance -= amount
	c.LastUpdatedAt = time.Now().UTC()

	return &CreditTransaction{
		ID:          uuid.New().String(),
		OrgID:       c.OrgID,
		Amount:      -amount,
		Type:        "usage",
		Status:      "completed",
		Description: description,
		ProcessedAt: time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// IsTrialActive checks if the subscription is in trial period
func (s *Subscription) IsTrialActive() bool {
	if s.TrialEndsAt == nil {
		return false
	}
	return s.Status == BillingStatusTrial && s.TrialEndsAt.After(time.Now().UTC())
}

// WillCancelAtPeriodEnd checks if the subscription is set to cancel at period end
func (s *Subscription) WillCancelAtPeriodEnd() bool {
	return s.CancelAtPeriodEnd && s.Status == BillingStatusActive
}

// Errors
var (
	ErrInsufficientCredits = fmt.Errorf("insufficient credits")
	ErrInvalidPlan         = fmt.Errorf("invalid subscription plan")
	ErrBillingDisabled     = fmt.Errorf("billing is not enabled")
)
