package helper

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func GetOrderKey(orderID int64) string {
	return "order:" + strconv.FormatInt(orderID, 10)
}

func GenerateRandomInt64() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(1<<63)) // Generate a random int64 value
	id := n.Int64()
	return id
}
