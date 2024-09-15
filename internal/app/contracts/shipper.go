package contracts

import "github.com/craftizmv/rewards/internal/data/dtos"

type Shipper interface {
	// ShipItem ShipRewardItem ShipOrder Ship an order and return shipment tracking ID
	ShipItem(itemID int64, detail *dtos.UserDetail) (*dtos.ShipmentResponse, error)
	ShipItems(itemID []int64, detail *dtos.UserDetail) (*dtos.ShipmentResponse, error)

	// GetShipmentStatus Get the status of a shipment
	GetShipmentStatus(shipmentID string) (string, error)
}
