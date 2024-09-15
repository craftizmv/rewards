package proxies

import (
	. "github.com/craftizmv/rewards/internal/data/dtos"
	"github.com/craftizmv/rewards/internal/data/infrastructure/external/mocks"
)

type InventoryProxy struct{}

func NewInventoryProxy() *InventoryProxy {
	return &InventoryProxy{}
}

// ReturnMockInventoryData - returns list of mock campaigns
func (p *InventoryProxy) ReturnMockInventoryData(productID int64) *InventoryDTO {
	return mocks.GetMockInventoryByProductID(productID)
}

func (p *InventoryProxy) BulkVerifyInventoryAvailability(itemIDs []int64) (bool, []int64) {
	return true, []int64{}
}

func (p *InventoryProxy) BlockInventoryForProducts(productIDs []int64) (bool, []int64) {
	// returns the itemIDs list, empty if not able to retrieve
	return true, []int64{}
}
