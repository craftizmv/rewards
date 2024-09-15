package helper

import (
	. "github.com/craftizmv/rewards/internal/domain/entities"
	"time"
)

// CreateOrderRewardItems creates a slice of OrderRewardItem pointers based on the current time, orderID, and a list of rewardItemIDs
func CreateOrderRewardItems(orderID int64, rewardItemIDs []int64) []*OrderRewardItem {
	// Get the current time to use as the allocated date
	now := time.Now()

	// Create a slice to hold the OrderRewardItems
	orderRewardItems := make([]*OrderRewardItem, len(rewardItemIDs))

	// Iterate over the rewardItemIDs and create an OrderRewardItem for each one
	for i, rewardItemID := range rewardItemIDs {
		orderRewardItems[i] = &OrderRewardItem{
			OrderID:       orderID,
			RewardItemID:  rewardItemID,
			ShipmentID:    nil, // No shipment associated for now
			AllocatedDate: now, // Use the current time for allocated date
			IsRedeemed:    nil, // Not redeemed for now
			RedeemedDate:  nil, // No redeemed date for now
		}
	}

	// Return the slice of OrderRewardItem pointers
	return orderRewardItems
}
