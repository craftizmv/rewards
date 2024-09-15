package repository

import (
	. "github.com/craftizmv/rewards/internal/domain/entities"
)

// RewardRepository defines the interface for reward-related data operations
type RewardRepository interface {
	GetRewardGroupByID(id int64) (*RewardGroup, error)
	GetRewardItemIDsFromRewardGroup(rewardGroupID int64) ([]int64, error)
	GetProductIDsFromRewardGroup(rewardGroupID int64) ([]int64, error)
	InsertRewardGroupRewardItem(rewardGroupID, rewardItemID int64) error
	InsertRewardGroupRewardItemsBatch(rewardGroupID int64, rewardItemIDs []int64, batchSize int) error
	UpdateOrderRewardItemsBatch(orderRewardItems []*OrderRewardItem, batchSize int) error
	GetRewardGroupIDByOrderID(orderID int64) ([]int64, error)
	DeleteRewardGroupByOrderID(orderID int64, rewardGroupID int64) error
	DeleteRewardItemsByOrderID(orderID int64) error
	DeleteRewardItemsByRewardGroupID(rewardGroupID int64) error
}
