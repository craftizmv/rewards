package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// RewardItemType represents the type of the reward or reward
type RewardItemType string

const (
	RewardTypeProduct  RewardItemType = "product"  // A physical product
	RewardTypeDiscount RewardItemType = "discount" // A discount or coupon
	RewardTypeVoucher  RewardItemType = "voucher"  // A reward voucher or code
)

// RewardConditions is a custom type to hold conditions in a flexible structure like JSON
type RewardConditions map[string]interface{}

// RewardItem represents the structure for reward items in the reward system.
type RewardItem struct {
	RewardItemID     string           `json:"reward_item_id"`              // Unique identifier for the reward item
	Type             RewardItemType   `json:"type"`                        // Type of reward (Discount, Product, Voucher)
	ItemID           *string          `json:"item_id,omitempty"`           // ItemID if the reward is related to a product
	DiscountAmount   *float64         `json:"discount_amount,omitempty"`   // Discount amount (only for Discount type)
	VoucherCode      *string          `json:"voucher_code,omitempty"`      // Voucher code (only for Voucher type)
	ProductID        *string          `json:"product_id,omitempty"`        // ProductID if the reward item is a product
	ExpirationDate   *time.Time       `json:"expiration_date,omitempty"`   // Expiration date of the reward
	RewardConditions RewardConditions `json:"reward_conditions,omitempty"` // Flexible structure to store reward conditions
	IsActive         bool             `json:"is_active"`                   // Indicates if the reward is currently active
	Metadata         json.RawMessage  `json:"metadata,omitempty"`          // Any additional metadata as a JSON object
}

// NewRewardItem A constructor function to create a new reward item with default values
func NewRewardItem(rewardItemID string, rewardType RewardItemType) *RewardItem {
	return &RewardItem{
		RewardItemID: rewardItemID,
		Type:         rewardType,
		IsActive:     true, // default to active
	}
}

// SetDiscount Example method to set the discount amount for Discount reward type
func (r *RewardItem) SetDiscount(discount float64) {
	if r.Type == RewardTypeDiscount {
		r.DiscountAmount = &discount
	}
}

// SetVoucherCode Example method to set the voucher code for Voucher reward type
func (r *RewardItem) SetVoucherCode(code string) {
	if r.Type == RewardTypeVoucher {
		r.VoucherCode = &code
	}
}

// SetProductID Example method to set a product ID for Product reward type
func (r *RewardItem) SetProductID(productID string) {
	if r.Type == RewardTypeProduct {
		r.ProductID = &productID
	}
}

// AddCondition Method to add a reward condition
func (r *RewardItem) AddCondition(key string, value interface{}) {
	if r.RewardConditions == nil {
		r.RewardConditions = make(RewardConditions)
	}
	r.RewardConditions[key] = value
}

// Method to deactivate a reward item
func (r *RewardItem) Deactivate() {
	r.IsActive = false
}

// Validate checks the RewardItem against business invariants
func (r *RewardItem) Validate() error {
	// Check if reward is active
	if !r.IsActive {
		return errors.New("reward item is inactive")
	}

	// Check if the reward has expired (if expiration date is set)
	if r.ExpirationDate != nil && time.Now().After(*r.ExpirationDate) {
		return errors.New("reward item is expired")
	}

	// Validate based on the type of reward
	switch r.Type {
	case RewardTypeDiscount:
		if r.DiscountAmount == nil || *r.DiscountAmount <= 0 {
			return errors.New("discount amount must be set and greater than 0 for discount type reward")
		}
	case RewardTypeProduct:
		if r.ProductID == nil && r.ItemID == nil {
			return errors.New("product or item ID must be set for product type reward")
		}
	case RewardTypeVoucher:
		if r.VoucherCode == nil || *r.VoucherCode == "" {
			return errors.New("voucher code must be set for voucher type reward")
		}
	default:
		return fmt.Errorf("invalid reward type: %s", r.Type)
	}

	// Validate any custom conditions set for the reward (optional)
	if err := r.validateConditions(); err != nil {
		return err
	}

	// All validations passed
	return nil
}

// validateConditions checks the reward-specific conditions (optional implementation)
func (r *RewardItem) validateConditions() error {
	if r.RewardConditions != nil {
		// Example condition: check if a minimum purchase amount is required
		if minPurchase, ok := r.RewardConditions["min_purchase_amount"]; ok {
			if minPurchaseAmount, ok := minPurchase.(float64); !ok || minPurchaseAmount <= 0 {
				return errors.New("invalid minimum purchase amount condition")
			}
		}

		// Add other condition validations as needed
		// Example: check if a condition "valid_on_weekends" is true
		if validOnWeekends, ok := r.RewardConditions["valid_on_weekends"]; ok {
			if validOnWeekendsBool, ok := validOnWeekends.(bool); ok && validOnWeekendsBool {
				// Only allow the reward on weekends (Saturday, Sunday)
				weekday := time.Now().Weekday()
				if weekday != time.Saturday && weekday != time.Sunday {
					return errors.New("reward can only be used on weekends")
				}
			}
		}
	}

	// All conditions passed
	return nil
}
