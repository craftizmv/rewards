package proxies

type OrderProxy struct {
}

func NewOrderProxy() *OrderProxy {
	return &OrderProxy{}
}

func (p OrderProxy) UpdateOrderRewardStatus(orderID int64, rewardID int64, status string) (bool, error) {
	return true, nil
}
