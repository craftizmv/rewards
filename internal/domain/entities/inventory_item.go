package entities

import "time"

type InventoryItem struct {
	ID                int       `json:"id"`
	ProductID         int       `json:"product_id"`
	UpcID             string    `json:"upc_id"`
	QuantityInStock   int       `json:"quantity_in_stock"`
	WarehouseLocation string    `json:"warehouse_location"`
	DateReceived      time.Time `json:"date_received"`
	ExpirationDate    time.Time `json:"expiration_date"`
	Status            string    `json:"status"`
}
