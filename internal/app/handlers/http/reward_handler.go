package http

import (
	"github.com/craftizmv/rewards/internal/app/usecase"
	"github.com/craftizmv/rewards/internal/data/dtos"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RewardHandler struct {
	useCase usecase.RewardUseCase
	log     logger.ILogger
}

func NewRewardHandler(usecase usecase.RewardUseCase, logger logger.ILogger) *RewardHandler {
	return &RewardHandler{
		useCase: usecase,
		log:     logger,
	}
}

func (h *RewardHandler) CheckRewardEligibility(c echo.Context) error {
	// call useCase layer etc.
	reqBody := new(dtos.OrderDTO)

	if err := c.Bind(reqBody); err != nil {
		h.log.Errorf("Error binding request body: %v", err)
		return SendResponse(c, http.StatusBadRequest, "Bad request")
	}

	if status, err := h.useCase.CheckRewardEligibility(reqBody); err != nil {
		return SendResponse(c, http.StatusInternalServerError, "could not check, please try again")
	} else if status {
		eligibilityResponse := dtos.RewardEligibilityResponse{
			Eligible: true,
			Message:  "reward is eligible",
		}
		return SendResponseWithData(c, http.StatusOK, "", eligibilityResponse)
	}

	return SendResponse(c, http.StatusOK, "OK")
}
