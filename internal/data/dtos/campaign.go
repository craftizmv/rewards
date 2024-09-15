package dtos

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// Status represents the possible statuses of a campaign.
type Status string

const (
	Active Status = "active"
	Paused Status = "paused"
	Ended  Status = "ended"
)

// CampaignDTO represents the data structure of a Campaign received from the Campaign Service.
type CampaignDTO struct {
	ID                   uuid.UUID           `json:"id"`
	RewardGroupID        int64               `json:"reward_group_id"`
	Name                 string              `json:"name"`
	StartDate            time.Time           `json:"start_date"`
	EndDate              time.Time           `json:"end_date"`
	Status               Status              `json:"status"`
	Budget               float64             `json:"budget"`
	AllocatedRewards     int                 `json:"allocated_rewards"`
	TotalEligibleRewards int                 `json:"total_eligible_rewards"`
	TargetAudience       string              `json:"target_audience"`
	EligibilityCriteria  EligibilityCriteria // Criteria that must be met to redeem the reward
}

// EligibilityCriteria defines conditions that must be met to redeem the reward
type EligibilityCriteria struct {
	MinimumPurchaseAmount float64
	// Additional criteria can be added here
}
