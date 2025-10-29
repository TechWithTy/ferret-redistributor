package models

import (
	"time"

	"github.com/google/uuid"
)

// BusinessModel represents the core business model configuration
type BusinessModel struct {
	ID             string    `json:"id"`
	Name           string    `json:
ame"`
	Description    string    `json:"description"`
	PricingModel   string    `json:"pricing_model"` // freemium, subscription, pay_as_you_go, etc.
	RevenueStreams []string  `json:"revenue_streams"`
	CostStructure  []string  `json:"cost_structure"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NewBusinessModel creates a new BusinessModel with default values
func NewBusinessModel(name, description string) *BusinessModel {
	now := time.Now().UTC()
	return &BusinessModel{
		ID:             uuid.New().String(),
		Name:           name,
		Description:    description,
		PricingModel:   "freemium",
		RevenueStreams: []string{"subscription", "credits"},
		CostStructure:  []string{"hosting", "third_party_apis", "support"},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// AddRevenueStream adds a new revenue stream
func (bm *BusinessModel) AddRevenueStream(stream string) {
	for _, s := range bm.RevenueStreams {
		if s == stream {
			return
		}
	}
	bm.RevenueStreams = append(bm.RevenueStreams, stream)
	bm.UpdatedAt = time.Now().UTC()
}

// AddCost adds a new cost item to the cost structure
func (bm *BusinessModel) AddCost(cost string) {
	for _, c := range bm.CostStructure {
		if c == cost {
			return
		}
	}
	bm.CostStructure = append(bm.CostStructure, cost)
	bm.UpdatedAt = time.Now().UTC()
}
