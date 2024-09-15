package proxies

import (
	"errors"
	"fmt"
	. "github.com/craftizmv/rewards/internal/app/contracts"
	"github.com/craftizmv/rewards/internal/data/dtos"
	"time"
)

// ShippingProxy is a proxy that adds additional functionality before delegating to the actual shipper
type ShippingProxy struct {
	shipper Shipper
}

// NewShippingProxy creates a new instance of the ShippingProxy with the injected shipper (e.g., LogiDeli)
func NewShippingProxy(shipper Shipper) *ShippingProxy {
	return &ShippingProxy{
		shipper: shipper,
	}
}

// ShipItem ShipRewardItem adds logging and retries, then delegates the actual shipping operation to the shipper
func (p *ShippingProxy) ShipItem(itemID int64, shipmentDetail *dtos.UserDetail) (*dtos.ShipmentResponse, error) {
	fmt.Printf("Starting shipping process for Order %d...\n", itemID)

	// Add retry logic (retry 3 times)
	for i := 0; i < 3; i++ {
		shipmentResponse, err := p.shipper.ShipItem(itemID, shipmentDetail)
		if err == nil {
			fmt.Printf("Shipping successful for Order %d with Tracking ID: %s\n", itemID, shipmentResponse)
			return shipmentResponse, nil
		}
		fmt.Printf("Attempt %d to ship Order %d failed: %v. Retrying...\n", i+1, itemID, err)
		time.Sleep(1 * time.Second)
	}

	// TODO: We need to send an event here in RabbitMQ/Kafka to re-receive and retry this.

	return nil, fmt.Errorf("failed to ship Order %d after 3 attempts", itemID)
}

// ShipItems performs the shipping of multiple items with retries and logs the result
func (p *ShippingProxy) ShipItems(itemIDs []int64, shipmentDetail *dtos.UserDetail) (*dtos.ShipmentResponse, error) {
	// Retry logic, logging, etc.
	for i := 0; i < 3; i++ {
		shipmentResponse, err := p.shipper.ShipItems(itemIDs, shipmentDetail)
		if err == nil {
			fmt.Printf("successfully shipped multiple items. Tracking ID: %s\n", shipmentResponse)
			return shipmentResponse, nil
		} else {
			if shipmentResponse.Error != nil {
				if *shipmentResponse.Error != dtos.ErrorSystemFailure {
					// then return break from retrying ...
					return shipmentResponse, errors.New("error in supplied data, shipper could not ship")
				}
			}
		}

		fmt.Printf("Attempt %d to shipfailed. Retrying...\n", i+1)
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("failed to ship after 3 attempts")
}

// GetShipmentStatus adds logging and then delegates the actual status check to the shipper
func (p *ShippingProxy) GetShipmentStatus(shipmentID string) (string, error) {
	fmt.Printf("Querying shipment status for Tracking ID: %s...\n", shipmentID)
	status, err := p.shipper.GetShipmentStatus(shipmentID)
	if err != nil {
		fmt.Printf("Failed to retrieve shipment status for %s: %v\n", shipmentID, err)
		return "", err
	}
	fmt.Printf("Shipment %s is currently: %s\n", shipmentID, status)

	// TODO : we will need to map these shipment status to domain shipment status
	return status, nil
}
