package entities

import (
	"errors"
)

// RewardGroupRewardItem represents the association between a RewardGroup and a RewardItem.
type RewardGroupRewardItem struct {
	RewardGroupID string `json:"reward_group_id"` // Foreign key referencing the RewardGroup
	RewardItemID  string `json:"reward_item_id"`  // Foreign key referencing the RewardItem
}

// NewRewardGroupRewardItem Constructor to create a new RewardGroupRewardItem
func NewRewardGroupRewardItem(rewardGroupID, rewardItemID string) *RewardGroupRewardItem {
	return &RewardGroupRewardItem{
		RewardGroupID: rewardGroupID,
		RewardItemID:  rewardItemID,
	}
}

// Validate checks the RewardGroupRewardItem against business invariants
func (rgri *RewardGroupRewardItem) Validate() error {
	// Check if RewardGroupID and RewardItemID are set
	if rgri.RewardGroupID == "" {
		return errors.New("reward group ID must be set")
	}
	if rgri.RewardItemID == "" {
		return errors.New("reward item ID must be set")
	}

	// All validations passed
	return nil
}
