package dtos

import "time"

// InventoryDTO represents the data transfer object for inventory items
type InventoryDTO struct {
	ProductID   int64     `json:"product_id"`   // unique product ID
	Quantity    int       `json:"quantity"`     // Amount of product available
	Location    string    `json:"location"`     // Storage location of the inventory
	LastUpdated time.Time `json:"last_updated"` // Timestamp of the last update
}

// NewInventoryDTO creates a new InventoryDTO
func NewInventoryDTO(productID int64, quantity int, location string, lastUpdated time.Time) *InventoryDTO {
	return &InventoryDTO{
		ProductID:   productID,
		Quantity:    quantity,
		Location:    location,
		LastUpdated: lastUpdated,
	}
}
