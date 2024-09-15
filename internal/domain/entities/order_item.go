package entities

import "errors"

// OrderItem represents an item in the order
type OrderItem struct {
	ProductID int64   // Unique identifier for the product
	Name      string  // Product name
	Quantity  int     // Number of units ordered
	Price     float64 // Price per unit
}

// NewOrderItem creates a new OrderItem with necessary validations
func NewOrderItem(productID int64, name string, quantity int, price float64) (*OrderItem, error) {
	if productID <= 0 {
		return nil, errors.New("product ID must be valid")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	item := &OrderItem{
		ProductID: productID,
		Name:      name,
		Quantity:  quantity,
		Price:     price,
	}

	return item, nil
}

// TotalPrice calculates the total price for the item
func (oi *OrderItem) TotalPrice() float64 {
	return float64(oi.Quantity) * oi.Price
}
