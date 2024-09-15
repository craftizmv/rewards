package dtos

type UserDetail struct {
	UserID        string       `json:"user_id"`
	UserName      string       `json:"user_name"`
	Email         string       `json:"email"`
	Location      UserLocation `json:"location"`
	PhoneNumber   *string      `json:"phone_number,omitempty"`   // Optional contact number for delivery
	DeliveryNotes *string      `json:"delivery_notes,omitempty"` // Optional delivery instructions
}

// UserLocation represents a delivery location for home delivery
type UserLocation struct {
	StreetAddress string   `json:"street_address"`           // Main address (e.g., "123 Main St")
	AddressLine2  *string  `json:"address_line_2,omitempty"` // Optional second address line (e.g., "Apt 4B")
	City          string   `json:"city"`                     // City or locality
	State         string   `json:"state"`                    // State or region
	PostalCode    string   `json:"postal_code"`              // Postal or ZIP code
	Country       string   `json:"country"`                  // Country name or code (e.g., "US", "Germany")
	Latitude      *float64 `json:"latitude,omitempty"`       // Geolocation: latitude
	Longitude     *float64 `json:"longitude,omitempty"`      // Geolocation: longitude
}
