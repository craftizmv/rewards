package usecase

import (
	"github.com/craftizmv/rewards/internal/data/dtos"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/events"
)

type RewardUseCase interface {
	AllocateReward(allocateReward events.AllocateReward) error
	CancelReward(orderCancelledEvent events.RevokeReward) error
	ReAllocateReward(orderEvent events.ReAllocateReward) error
	CheckRewardEligibility(dto *dtos.OrderDTO) (bool, error)
}
