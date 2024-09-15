package consumers

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/iancoleman/strcase"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
	"time"
)

var orderConfirmedMessages []string

type OrderConfirmedConsumer[T any] struct {
	*BaseConsumer
	handler func(queue string, msg amqp.Delivery, dependencies T) error
	ctx     context.Context
}

func NewOrderConfirmedConsumer[T any](ctx context.Context, cfg *queue.RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, handler func(queue string, msg amqp.Delivery, dependencies T) error) IConsumer[T] {
	return &OrderConfirmedConsumer[T]{
		ctx: ctx,
		BaseConsumer: &BaseConsumer{
			cfg:  cfg,
			conn: conn,
			log:  log,
		},
		handler: handler,
	}
}

func (c *OrderConfirmedConsumer[T]) ConsumeMessage(msg interface{}, dependencies T) error {
	ch, err := c.conn.Channel()
	if err != nil {
		c.log.Error("Error in opening channel to consume message")
		return err
	}

	defer ch.Close()

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)

	err = ch.ExchangeDeclare(
		snakeTypeName, // exchange name
		c.cfg.Kind,    // type of exchange - we have used topic type
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

	orderConfirmedQueue := fmt.Sprintf("%s_%s", snakeTypeName, "order_confirmed")
	q, err := ch.QueueDeclare(
		orderConfirmedQueue, // name
		false,               // durable
		false,               // delete when unused
		true,                // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring queue to consume message")
		return err
	}

	err = ch.QueueBind(
		q.Name,              // queue name
		orderConfirmedQueue, // routing key
		snakeTypeName,       // exchange
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
		for {
			select {
			case <-c.ctx.Done():
				defer func(ch *amqp.Channel) {
					err := ch.Close()
					if err != nil {
						c.log.Errorf("failed to close channel for queue: %s", q.Name)
					}
				}(ch)
				c.log.Infof("channel closed for queue: %s", q.Name)
				return

			case delivery, ok := <-deliveries:
				if !ok {
					c.log.Errorf("NOT OK deliveries channel closed for queue: %s", q.Name)
					return
				}

				err := c.handler(q.Name, delivery, dependencies)
				if err != nil {
					c.log.Error(err.Error())
				}

				orderConfirmedMessages = append(orderConfirmedMessages, snakeTypeName)

				// Cannot use defer inside a for loop
				time.Sleep(1 * time.Millisecond)

				err = delivery.Ack(false)
				if err != nil {
					c.log.Errorf("We didn't get an ack for delivery: %v", string(delivery.Body))
				}
			}
		}
	}()

	c.log.Infof("Waiting for messages in queue :%s. To exit press CTRL+C", q.Name)

	return nil
}

func (c *OrderConfirmedConsumer[T]) IsConsumed(msg interface{}) bool {
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

		isConsumed = linq.From(orderConfirmedMessages).Contains(snakeTypeName)

		timeOutExpired = time.Now().Sub(startTime) > timeOutTime
	}
}
