package publisher

import (
	"context"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"time"
)

type BasePublisher struct {
	cfg  *queue.RabbitMQConfig
	conn *amqp.Connection
	log  logger.ILogger
}

func (bp *BasePublisher) prepareMessage(msg interface{}) ([]byte, string, error) {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		bp.log.Error("Error marshalling message")
		return nil, "", err
	}

	typeName := reflect.TypeOf(msg).Elem().Name()
	snakeTypeName := strcase.ToSnake(typeName)
	return data, snakeTypeName, nil
}

func (bp *BasePublisher) declareExchange(channel *amqp.Channel, exchangeName string) error {
	err := channel.ExchangeDeclare(
		exchangeName, // name
		bp.cfg.Kind,  // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		bp.log.Error("Error in declaring exchange")
	}
	return err
}

func (bp *BasePublisher) createPublishingMessage(ctx context.Context, data []byte) amqp.Publishing {
	correlationId := ""
	if ctx.Value(echo.HeaderXCorrelationID) != nil {
		correlationId = ctx.Value(echo.HeaderXCorrelationID).(string)
	}

	return amqp.Publishing{
		Body:          data,
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		MessageId:     uuid.NewV4().String(),
		Timestamp:     time.Now(),
		CorrelationId: correlationId,
	}
}
