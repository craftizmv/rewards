package events

type AllocateReward struct {
	UserID       string `json:"user_id"`
	OrderID      int64  `json:"order_id"`
	CampaignID   int64  `json:"campaign_id"`
	RewardTypeID int64  `json:"reward_type_id"`
	OrderStatus  string `json:"order_status"`
	OrderValue   int    `json:"order_value"`
}

type ReAllocateReward struct {
	UserID        string `json:"user_id"`
	CampaignID    int64  `json:"campaign_id"`
	RewardTypeID  int64  `json:"reward_type_id"`
	RewardGroupID int64  `json:"reward_group_id"`
}

type RevokeReward struct {
	UserID       string `json:"user_id"`
	OrderID      int64  `json:"order_id"`
	OrderStatus  string `json:"order_status"`
	CampaignID   int64  `json:"campaign_id"`
	RewardTypeID int64  `json:"reward_type_id"`
}
