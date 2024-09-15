package http

import (
	"github.com/labstack/echo/v4"
)

type IRewardHandler interface {
	// CheckRewardEligibility - checks if order is eligible for the reward
	CheckRewardEligibility(c echo.Context) error
}
