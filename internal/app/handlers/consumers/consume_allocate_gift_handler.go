package consumers

import (
	"encoding/json"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandleAllocateReward(queue string, msg amqp.Delivery, orderDeliveryBase *queue.OrderDeliveryBase) error {
	log := orderDeliveryBase.Log

	log.Infof("Message received on queue: %s with message: %s", queue, string(msg.Body))

	var event events.AllocateReward
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		return err
	}

	err = orderDeliveryBase.GiftUseCases.AllocateReward(event)
	if err != nil {
		return err
	}

	return nil
}
