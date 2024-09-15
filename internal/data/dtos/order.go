package dtos

// OrderDTO represents the data structure for an order request to check reward eligibility.
type OrderDTO struct {
	OrderID    int64   `json:"order_id"`    // Unique identifier for the order
	OrderValue float64 `json:"order_value"` // Total value of the order
	Quantity   int     `json:"quantity"`    // Number of items in the order
}
