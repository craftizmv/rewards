package mocks

import (
	. "github.com/craftizmv/rewards/internal/data/dtos"
	"time"
)

// GetMockInventoryByProductID returns a mock InventoryDTO based on the provided ProductID
func GetMockInventoryByProductID(productID int64) *InventoryDTO {
	// Mock data based on the productID
	switch productID {
	case 1001:
		return NewInventoryDTO(productID, 150, "Aisle 1", time.Now().Add(-24*time.Hour))
	case 1002:
		return NewInventoryDTO(productID, 200, "Aisle 2", time.Now().Add(-48*time.Hour))
	case 1003:
		return NewInventoryDTO(productID, 75, "Warehouse", time.Now().Add(-72*time.Hour))
	case 1004:
		return NewInventoryDTO(productID, 300, "Aisle 3", time.Now().Add(-96*time.Hour))
	default:
		return NewInventoryDTO(productID, 0, "Unknown", time.Now())
	}
}

// GetAllMockInventories returns a slice of all mock InventoryDTOs
func GetAllMockInventories() []InventoryDTO {
	return []InventoryDTO{
		*NewInventoryDTO(1001, 150, "Aisle 1", time.Now().Add(-24*time.Hour)),
		*NewInventoryDTO(1002, 200, "Aisle 2", time.Now().Add(-48*time.Hour)),
		*NewInventoryDTO(1003, 75, "Warehouse", time.Now().Add(-72*time.Hour)),
		*NewInventoryDTO(1004, 300, "Aisle 3", time.Now().Add(-96*time.Hour)),
	}
}
