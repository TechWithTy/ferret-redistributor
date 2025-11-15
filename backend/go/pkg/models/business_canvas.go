package models

import (
	"time"
)

// BusinessExperiment tracks business experiments and their outcomes
type BusinessExperiment struct {
	ID          string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string                 `json:"name" gorm:"not null"`
	Hypothesis  string                 `json:"hypothesis" gorm:"not null"`
	StartDate   time.Time              `json:"start_date" gorm:"not null"`
	EndDate     *time.Time             `json:"end_date,omitempty"`
	Status      string                 `json:"status" gorm:"not null;default:'draft'"` // draft, running, completed, paused
	Metrics     map[string]interface{} `json:"metrics" gorm:"type:jsonb"`
	Results     string                 `json:"results,omitempty"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// ValueProposition defines the value offered to different customer segments
type ValueProposition struct {
	ID              string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name            string                 `json:"name" gorm:"not null"`
	Description     string                 `json:"description"`
	TargetSegment   string                 `json:"target_segment"`
	Benefits        []string               `json:"benefits" gorm:"type:text[]"`
	Metrics         map[string]interface{} `json:"metrics" gorm:"type:jsonb"`
	IsActive        bool                   `json:"is_active" gorm:"default:true"`
	OrganizationID  *string                `json:"organization_id,omitempty"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// BusinessPartnership represents business partnerships and integrations
type BusinessPartnership struct {
	ID            string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string                 `json:"name" gorm:"not null"`
	Type          string                 `json:"type" gorm:"not null"` // integration, reseller, strategic, etc.
	Status        string                 `json:"status" gorm:"not null;default:'prospect'"` // prospect, active, inactive
	StartDate     *time.Time             `json:"start_date,omitempty"`
	EndDate       *time.Time             `json:"end_date,omitempty"`
	ContactInfo   map[string]string      `json:"contact_info" gorm:"type:jsonb"`
	Terms         string                 `json:"terms,omitempty"`
	Metadata      map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt     time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// CostCenter tracks business expenses and budgets
type CostCenter struct {
	ID             string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string     `json:"name" gorm:"not null"`
	Description    string     `json:"description"`
	Budget         float64    `json:"budget" gorm:"not null;default:0"`
	OwnerID        string     `json:"owner_id" gorm:"not null"`
	ParentID       *string    `json:"parent_id,omitempty"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	OrganizationID string     `json:"organization_id" gorm:"not null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// CustomerSegment defines different customer segments for targeting
type CustomerSegment struct {
	ID              string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name            string                 `json:"name" gorm:"not null"`
	Description     string                 `json:"description"`
	Criteria        map[string]interface{} `json:"criteria" gorm:"type:jsonb"`
	SizeEstimate    *int                   `json:"size_estimate,omitempty"`
	GrowthRate      *float64               `json:"growth_rate,omitempty"`
	LTV             *float64               `json:"ltv,omitempty"` // Lifetime Value
	CAC             *float64               `json:"cac,omitempty"`  // Customer Acquisition Cost
	IsActive        bool                   `json:"is_active" gorm:"default:true"`
	OrganizationID  string                 `json:"organization_id" gorm:"not null"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// BusinessModelCanvas represents the complete business model canvas
type BusinessModelCanvas struct {
	ID                    string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name                  string                 `json:"name" gorm:"not null"`
	Description           string                 `json:"description"`
	CustomerSegments      []CustomerSegment      `json:"customer_segments,omitempty" gorm:"foreignKey:CanvasID"`
	ValuePropositions     []ValueProposition     `json:"value_propositions,omitempty" gorm:"foreignKey:CanvasID"`
	Channels              []string               `json:"channels" gorm:"type:text[]"`
	CustomerRelationships []string               `json:"customer_relationships" gorm:"type:text[]"`
	RevenueStreams        map[string]interface{} `json:"revenue_streams" gorm:"type:jsonb"`
	KeyResources          []string               `json:"key_resources" gorm:"type:text[]"`
	KeyActivities         []string               `json:"key_activities" gorm:"type:text[]"`
	KeyPartnerships       []BusinessPartnership  `json:"key_partnerships,omitempty" gorm:"foreignKey:CanvasID"`
	CostStructure         map[string]interface{} `json:"cost_structure" gorm:"type:jsonb"`
	OrganizationID        string                 `json:"organization_id" gorm:"not null"`
	IsActive              bool                   `json:"is_active" gorm:"default:true"`
	CreatedAt             time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}
