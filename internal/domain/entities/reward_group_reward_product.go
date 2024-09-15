package entities

// RewardGroupRewardProduct represents the association between a RewardGroup and a RewardItem.
type RewardGroupRewardProduct struct {
	RewardGroupID string `json:"reward_group_id"` // Foreign key referencing the RewardGroup
	ProductID     string `json:"product_id"`      // Foreign key referencing the RewardItem
}
