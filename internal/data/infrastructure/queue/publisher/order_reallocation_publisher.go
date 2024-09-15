package publisher

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/iancoleman/strcase"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
)

type OrderReAllocationEventPublisher struct {
	*BasePublisher // embedded struct.
	ctx            context.Context
}

var rewardReAllocatePublishedMessages []string

func (p *OrderReAllocationEventPublisher) PublishMessage(msg interface{}) error {
	data, snakeTypeName, err := p.prepareMessage(msg)
	if err != nil {
		return err
	}

	channel, err := p.conn.Channel()
	if err != nil {
		p.log.Error("Error opening channel")
		return err
	}
	defer channel.Close()

	err = p.declareExchange(channel, snakeTypeName)
	if err != nil {
		return err
	}

	publishingMsg := p.createPublishingMessage(p.ctx, data)
	err = channel.Publish(snakeTypeName, snakeTypeName, false, false, publishingMsg)
	if err != nil {
		p.log.Error("Error publishing message")
		return err
	}

	rewardReAllocatePublishedMessages = append(rewardReAllocatePublishedMessages, snakeTypeName)
	p.log.Infof("Published message: %s", publishingMsg.Body)

	return nil
}

func (p *OrderReAllocationEventPublisher) IsPublished(msg interface{}) bool {
	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)
	return linq.From(rewardReAllocatePublishedMessages).Contains(snakeTypeName)
}

func NewPublisher(ctx context.Context, cfg *queue.RabbitMQConfig, conn *amqp.Connection, log logger.ILogger) IPublisher {
	basePublisher := &BasePublisher{cfg: cfg, conn: conn, log: log}
	return &OrderReAllocationEventPublisher{
		ctx:           ctx,
		BasePublisher: basePublisher,
	}
}
