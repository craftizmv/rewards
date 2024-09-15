package entities

import (
	"errors"
	"time"
)

// OrderRewardItem represents the association between an Order and a RewardItem.
type OrderRewardItem struct {
	OrderID       int64      `json:"order_id"`                // Foreign key referencing the Order
	RewardItemID  int64      `json:"reward_item_id"`          // Foreign key referencing the RewardItem
	ShipmentID    *int64     `json:"shipment_id,omitempty"`   // Foreign key referencing the Shipment (if applicable)
	AllocatedDate time.Time  `json:"allocated_date"`          // Date when the reward was allocated to the order
	IsRedeemed    *bool      `json:"is_redeemed"`             // Indicates if the reward item has been redeemed
	RedeemedDate  *time.Time `json:"redeemed_date,omitempty"` // Date when the reward was redeemed, if applicable
}

// Validate checks the OrderRewardItem against business invariants
func (ori *OrderRewardItem) Validate() error {
	// 1. Ensure that OrderID and RewardItemID are set
	if ori.OrderID == 0 {
		return errors.New("order ID must be set")
	}
	if ori.RewardItemID == 0 {
		return errors.New("reward item ID must be set")
	}

	// 2. Validate the AllocatedDate (it should not be in the future)
	if ori.AllocatedDate.After(time.Now()) {
		return errors.New("allocated date cannot be in the future")
	}

	// 3. If the item is marked as redeemed, RedeemedDate must be set
	if ori.IsRedeemed != nil && *ori.IsRedeemed && ori.RedeemedDate == nil {
		return errors.New("redeemed date must be set when the reward is marked as redeemed")
	}

	// 4. If the item is not redeemed, RedeemedDate should be nil
	if ori.IsRedeemed != nil && !*ori.IsRedeemed && ori.RedeemedDate != nil {
		return errors.New("redeemed date should be nil if the item is not redeemed")
	}

	// 5. Validate ShipmentID (optional check)
	// If ShipmentID is set, ensure it's valid (positive number)
	if ori.ShipmentID != nil && *ori.ShipmentID <= 0 {
		return errors.New("shipment ID, if present, must be a valid positive integer")
	}

	// If all validations pass, return nil
	return nil
}
