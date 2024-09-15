package consumers

import (
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BaseConsumer struct {
	cfg  *queue.RabbitMQConfig
	conn *amqp.Connection
	log  logger.ILogger
}

func (bc *BaseConsumer) DeclareExchange(exchangeName, exchangeType string) error {
	ch, err := bc.conn.Channel()
	if err != nil {
		bc.log.Error("Error opening channel for exchange declaration")
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // type of exchange
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func (bc *BaseConsumer) DeclareQueue(queueName string) (amqp.Queue, error) {
	ch, err := bc.conn.Channel()
	if err != nil {
		bc.log.Error("Error opening channel for queue declaration")
		return amqp.Queue{}, err
	}
	defer ch.Close()

	return ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func (bc *BaseConsumer) BindQueue(queueName, routingKey, exchangeName string) error {
	ch, err := bc.conn.Channel()
	if err != nil {
		bc.log.Error("Error opening channel for queue binding")
		return err
	}
	defer ch.Close()

	return ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
}
