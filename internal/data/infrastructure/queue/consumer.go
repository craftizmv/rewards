package queue

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/iancoleman/strcase"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
	"time"
)

//go:generate mockery --name IConsumer
type IConsumer[T any] interface {
	ConsumeMessage(msg interface{}, dependencies T) error
	IsConsumed(msg interface{}) bool
}

var consumedMessages []string

type Consumer[T any] struct {
	cfg     *RabbitMQConfig
	conn    *amqp.Connection
	log     logger.ILogger
	handler func(queue string, msg amqp.Delivery, dependencies T) error
	ctx     context.Context
}

func (c Consumer[T]) ConsumeMessage(msg interface{}, dependencies T) error {
	ch, err := c.conn.Channel()
	if err != nil {
		c.log.Error("Error in opening channel to consume message")
		return err
	}

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)

	err = ch.ExchangeDeclare(
		snakeTypeName, // name
		c.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring exchange to consume message")
		return err
	}

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s_%s", snakeTypeName, "queue"), // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring queue to consume message")
		return err
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil)
	if err != nil {
		c.log.Error("Error in binding queue to consume message")
		return err
	}

	deliveries, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		c.log.Error("Error in consuming message")
		return err
	}

	go func() {

		select {
		case <-c.ctx.Done():
			defer func(ch *amqp.Channel) {
				err := ch.Close()
				if err != nil {
					c.log.Errorf("failed to close channel closed for for queue: %s", q.Name)
				}
			}(ch)
			c.log.Infof("channel closed for for queue: %s", q.Name)
			return

		case delivery, ok := <-deliveries:
			{
				if !ok {
					c.log.Errorf("NOT OK deliveries channel closed for queue: %s", q.Name)
					return
				}

				err := c.handler(q.Name, delivery, dependencies)
				if err != nil {
					c.log.Error(err.Error())
				}

				consumedMessages = append(consumedMessages, snakeTypeName)

				// Cannot use defer inside a for loop
				time.Sleep(1 * time.Millisecond)

				err = delivery.Ack(false)
				if err != nil {
					c.log.Errorf("We didn't get a ack for delivery: %v", string(delivery.Body))
				}
			}
		}
	}()

	c.log.Infof("Waiting for messages in queue :%s. To exit press CTRL+C", q.Name)

	return nil
}

func (c Consumer[T]) IsConsumed(msg interface{}) bool {
	timeOutTime := 20 * time.Second
	startTime := time.Now()
	timeOutExpired := false
	isConsumed := false

	for {
		if timeOutExpired {
			return false
		}
		if isConsumed {
			return true
		}

		time.Sleep(time.Second * 2)

		typeName := reflect.TypeOf(msg).Name()
		snakeTypeName := strcase.ToSnake(typeName)

		isConsumed = linq.From(consumedMessages).Contains(snakeTypeName)

		timeOutExpired = time.Now().Sub(startTime) > timeOutTime
	}
}

func NewConsumer[T any](ctx context.Context, cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, handler func(queue string, msg amqp.Delivery, dependencies T) error) IConsumer[T] {
	return &Consumer[T]{ctx: ctx, cfg: cfg, conn: conn, log: log, handler: handler}
}
