package entities

import (
	"errors"
	"fmt"
	"time"
)

// ShipmentStatus defines the various statuses a shipment can have
type ShipmentStatus string

// Enum-like constants for different shipment statuses
const (
	StatusPending   ShipmentStatus = "Pending"
	StatusShipped   ShipmentStatus = "Shipped"
	StatusInTransit ShipmentStatus = "In Transit"
	StatusDelivered ShipmentStatus = "Delivered"
	StatusFailed    ShipmentStatus = "Failed"
	StatusCanceled  ShipmentStatus = "Canceled"
	StatusReturned  ShipmentStatus = "Returned"
)

// Shipment represents the structure of a shipment entity
type Shipment struct {
	ShipmentID       string         `json:"shipment_id"`
	OrderID          string         `json:"order_id"`
	Carrier          string         `json:"carrier"`
	TrackingID       string         `json:"tracking_id"`
	Status           ShipmentStatus `json:"status"`
	ShippedAt        *time.Time     `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time     `json:"delivered_at,omitempty"`
	EstimatedArrival *time.Time     `json:"estimated_arrival,omitempty"`
	LastUpdated      time.Time      `json:"last_updated"`
}

// Validate checks the Shipment against business invariants
func (s *Shipment) Validate() error {
	// 1. Ensure essential fields are not empty
	if s.ShipmentID == "" {
		return errors.New("shipment ID cannot be empty")
	}
	if s.OrderID == "" {
		return errors.New("order ID cannot be empty")
	}
	if s.Carrier == "" {
		return errors.New("carrier cannot be empty")
	}
	if s.Status == "" {
		return errors.New("status cannot be empty")
	}

	// 2. Check for valid status transitions
	if err := s.validateStatusTransition(); err != nil {
		return err
	}

	// 3. Ensure timestamps are valid based on the current status
	if s.Status == StatusShipped || s.Status == StatusInTransit || s.Status == StatusDelivered {
		if s.ShippedAt == nil {
			return errors.New("shipped_at must be set for shipped, in transit, or delivered statuses")
		}
	}
	if s.Status == StatusDelivered {
		if s.DeliveredAt == nil {
			return errors.New("delivered_at must be set for delivered status")
		}
		// 4. Ensure DeliveredAt is not earlier than ShippedAt
		if s.DeliveredAt.Before(*s.ShippedAt) {
			return errors.New("delivered_at cannot be earlier than shipped_at")
		}
	}

	// 5. Ensure LastUpdated timestamp is present
	if s.LastUpdated.IsZero() {
		return errors.New("last_updated cannot be zero")
	}

	// If all validations pass
	return nil
}

// validateStatusTransition ensures valid transitions between statuses
func (s *Shipment) validateStatusTransition() error {
	// Allowed transitions based on business logic
	allowedTransitions := map[ShipmentStatus][]ShipmentStatus{
		StatusPending:   {StatusShipped, StatusCanceled},
		StatusShipped:   {StatusInTransit, StatusDelivered, StatusFailed},
		StatusInTransit: {StatusDelivered, StatusReturned, StatusFailed},
		StatusDelivered: {}, // No further transitions allowed after delivery
		StatusFailed:    {}, // Failed shipments cannot transition to any other status
		StatusCanceled:  {}, // Canceled shipments cannot transition to any other status
		StatusReturned:  {}, // Returned shipments cannot transition to any other status
	}

	// Ensure that status transitions follow the allowed path
	allowedNextStatuses, ok := allowedTransitions[s.Status]
	if !ok {
		return fmt.Errorf("invalid current shipment status: %s", s.Status)
	}

	for _, nextStatus := range allowedNextStatuses {
		if s.Status == nextStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s", s.Status)
}
