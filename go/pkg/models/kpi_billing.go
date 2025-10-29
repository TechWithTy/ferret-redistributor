package models

import (
	"time"

	"github.com/google/uuid"
)

// KPITier represents a success-based pricing tier
type KPITier struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	KPI               string    `json:"kpi"` // e.g., "mrr_growth", "user_retention"
	TargetValue       float64   `json:"target_value"`
	SuccessThreshold  float64   `json:"success_threshold"` // e.g., 1.2 for 120% of target
	BasePrice         int64     `json:"base_price"`        // in cents
	SuccessMultiplier float64   `json:"success_multiplier"` // e.g., 1.1 for 10% bonus
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// KPIScore represents a KPI measurement for an organization
type KPIScore struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	KPI         string    `json:"kpi"`
	Value       float64   `json:"value"`
	TargetValue float64   `json:"target_value"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	CreatedAt   time.Time `json:"created_at"`
}

// SuccessBasedBilling represents the success-based billing configuration for an org
type SuccessBasedBilling struct {
	ID                string    `json:"id"`
	OrgID             string    `json:"org_id"`
	CurrentTierID     string    `json:"current_tier_id"`
	BasePlanPrice     int64     `json:"base_plan_price"` // in cents
	CurrentMultiplier float64   `json:"current_multiplier"`
	NextBillingDate   time.Time `json:"next_billing_date"`
	LastKPICalculation *time.Time `json:"last_kpi_calculation,omitempty"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// NewKPITier creates a new KPI tier
func NewKPITier(name, description, kpi string, targetValue, successThreshold float64, basePrice int64) *KPITier {
	now := time.Now().UTC()
	return &KPITier{
		ID:                uuid.New().String(),
		Name:              name,
		Description:       description,
		KPI:               kpi,
		TargetValue:       targetValue,
		SuccessThreshold:  successThreshold,
		BasePrice:         basePrice,
		SuccessMultiplier: 1.0, // Default to no bonus
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// CalculateSuccessMultiplier calculates the success multiplier based on KPI performance
func (t *KPITier) CalculateSuccessMultiplier(actualValue float64) float64 {
	if actualValue >= (t.TargetValue * t.SuccessThreshold) {
		return t.SuccessMultiplier
	}
	return 1.0 // Default multiplier if target not met
}

// RecordKPIScore records a new KPI measurement
func RecordKPIScore(orgID, kpi string, value, targetValue float64, periodStart, periodEnd time.Time) *KPIScore {
	return &KPIScore{
		ID:          uuid.New().String(),
		OrgID:       orgID,
		KPI:         kpi,
		Value:       value,
		TargetValue: targetValue,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		CreatedAt:   time.Now().UTC(),
	}
}

// CalculateSuccessBasedPrice calculates the adjusted price based on KPI performance
func CalculateSuccessBasedPrice(basePrice int64, kpiScore *KPIScore, tier *KPITier) int64 {
	if kpiScore == nil || tier == nil {
		return basePrice
	}

	multiplier := tier.CalculateSuccessMultiplier(kpiScore.Value)
	return int64(float64(basePrice) * multiplier)
}

// UpdateBillingForSuccess updates the billing based on KPI performance
func (sb *SuccessBasedBilling) UpdateBillingForSuccess(kpiScore *KPIScore, tier *KPITier) {
	if kpiScore == nil || tier == nil {
		return
	}

	sb.CurrentMultiplier = tier.CalculateSuccessMultiplier(kpiScore.Value)
	sb.UpdatedAt = time.Now().UTC()
}

// Common KPIs for success-based billing
const (
	KPIMRRGrowth       = "mrr_growth"
	KPIUserRetention   = "user_retention"
	KPIFeatureAdoption = "feature_adoption"
	KPIEngagement     = "engagement_score"
)

// DefaultKPITiers returns a set of default KPI tiers
func DefaultKPITiers() []*KPITier {
	now := time.Now().UTC()
	return []*KPITier{
		{
			ID:                uuid.New().String(),
			Name:              "Growth Accelerator",
			Description:       "Rewards for MRR growth",
			KPI:               KPIMRRGrowth,
			TargetValue:       10000,  // $100 MRR
			SuccessThreshold:  1.2,    // 120% of target
			BasePrice:         5000,   // $50
			SuccessMultiplier: 1.1,    // 10% bonus
			IsActive:          true,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		{
			ID:                uuid.New().String(),
			Name:              "Retention Champion",
			Description:       "Rewards for high user retention",
			KPI:               KPIUserRetention,
			TargetValue:       0.8,    // 80% retention
			SuccessThreshold:  1.1,    // 88% retention
			BasePrice:         3000,   // $30
			SuccessMultiplier: 1.15,   // 15% bonus
			IsActive:          true,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
	}
}
