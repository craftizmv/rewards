package dtos

import (
	"time"
)

// ErrorType represents different types of errors that could occur during shipment
type ErrorType string

const (
	ErrorInvalidAddress   ErrorType = "Invalid Address"
	ErrorOutOfServiceArea ErrorType = "Out of Service Area"
	ErrorSystemFailure    ErrorType = "System Failure"
)

// ShipmentResponse represents the response after trying to create a shipment
type ShipmentResponse struct {
	IsShippingPossible    bool       `json:"is_shipping_possible"`              // Indicates if shipping is possible
	Cost                  float64    `json:"cost,omitempty"`                    // The shipping cost, if applicable
	TentativeShipmentDate *time.Time `json:"tentative_shipment_date,omitempty"` // Tentative date when shipment will happen
	ConfirmationID        *string    `json:"confirmation_id,omitempty"`         // The confirmation shipping ID, if available
	Error                 *ErrorType `json:"error,omitempty"`                   // Describes the type of error, if any
}
