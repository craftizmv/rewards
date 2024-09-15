package entities

import (
	"time"
)

// OrderStatus defines possible statuses for an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCanceled   OrderStatus = "canceled"
	OrderStatusReturned   OrderStatus = "returned"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type RewardStatus string

const (
	RewardStatusNone      = RewardStatus("none")
	RewardStatusAllocated = RewardStatus("allocated")
	RewardStatusShipped   = RewardStatus("shipped")
	RewardStatusDelivered = RewardStatus("delivered")
	RewardStatusCancelled = RewardStatus("cancelled")
)

// Order represents an order in the domain
type Order struct {
	ID               int64
	CustomerID       int64
	Items            []OrderItem
	TotalPrice       float64
	Status           OrderStatus
	ReturnWindowTime time.Time
	RewardStatus     RewardStatus
}

// IsComplete checks if the order has been already completed.
func (o *Order) IsComplete() bool {
	if o.Status == OrderStatusDelivered && time.Now().After(o.ReturnWindowTime) {
		return true
	}

	if o.Status == OrderStatusCanceled || o.Status == OrderStatusReturned || o.Status == OrderStatusRefunded {
		return true
	}

	return false
}

func (o *Order) IsOrderEffectivelyRolledBack() bool {
	if o.Status == OrderStatusCanceled || o.Status == OrderStatusReturned || o.Status == OrderStatusRefunded {
		return true
	}

	return false
}

// isModifiable checks if an order can be modified
func (o *Order) isModifiable() bool {
	// handle special condition for order delivered and in the return window.
	if o.Status == OrderStatusDelivered && time.Now().Before(o.ReturnWindowTime) {
		return true
	}

	return o.Status != OrderStatusDelivered &&
		o.Status != OrderStatusCanceled &&
		o.Status != OrderStatusReturned &&
		o.Status != OrderStatusRefunded
}

// calculateTotalPrice computes the total price of an order
func calculateTotalPrice(items []OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
