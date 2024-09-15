package entities

import (
	"errors"
	"time"
)

// TODO : ReVisit the reward eligibility criteria if needed here.

// RewardGroup represents a reward that may contain multiple reward items
type RewardGroup struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CampaignID int64      `json:"campaign_id"`
}

// NewRewardGroup creates a new RewardGroup with necessary validations
func NewRewardGroup(rewardGroupID int64, name string, expiresAt *time.Time) (*RewardGroup, error) {
	if rewardGroupID <= 0 {
		return nil, errors.New("reward ID must be valid")
	}
	if name == "" {
		return nil, errors.New("reward name must not be empty")
	}

	if expiresAt != nil && expiresAt.Before(time.Now()) {
		return nil, errors.New("reward cannot be expired upon creation")
	}

	reward := &RewardGroup{
		ID:        rewardGroupID,
		Name:      name,
		ExpiresAt: expiresAt,
	}

	return reward, nil
}
