package proxies

import "github.com/craftizmv/rewards/internal/data/dtos"

type UserProxy struct {
}

func NewUserProxy() *UserProxy {
	return &UserProxy{}
}

// TODO: handler error cases
func (u UserProxy) GetUserDetails(userID string) *dtos.UserDetail {
	addr2 := "addr2"
	return &dtos.UserDetail{
		UserID:   "abc123",
		UserName: "MV",
		Email:    "abc@gmail.com",
		Location: dtos.UserLocation{
			StreetAddress: "add1",
			AddressLine2:  &addr2,
			City:          "bangalore",
			State:         "karnataka",
			PostalCode:    "5600",
			Country:       "INDIA",
			Latitude:      nil,
			Longitude:     nil,
		},
		PhoneNumber:   nil,
		DeliveryNotes: nil,
	}
}
