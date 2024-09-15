package queue

import (
	"context"
	"github.com/craftizmv/rewards/config"
	"github.com/craftizmv/rewards/internal/app/usecase"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/publisher"
	"github.com/craftizmv/rewards/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderDeliveryBase struct {
	Log          logger.ILogger
	Cfg          *config.Config
	ConnRabbitmq *amqp.Connection
	Ctx          context.Context
	GiftUseCases usecase.RewardUseCase
	Publisher    publisher.IPublisher
}
