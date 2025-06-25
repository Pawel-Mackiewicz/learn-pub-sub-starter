package pubsub

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type simpleQueueType string

const (
	QueueTypeDurable   simpleQueueType = "durable"
	QueueTypeTransient simpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType simpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {

	chConn, err := conn.Channel()
	if err != nil {
		fmt.Errorf("error while creating channel: %v", err)
	}
	var durable bool
	var autoDelete bool
	var exclusive bool
	noWait := false

	if queueType == QueueTypeDurable {
		durable = true
	} else if queueType == QueueTypeTransient {
		autoDelete = true
		exclusive = true
	} else {
		return nil, amqp.Queue{}, fmt.Errorf("wrong `queueType` selected")
	}
	queue, err := chConn.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error when declaring queue: %v", err)
	}
	err = chConn.QueueBind(queueName, key, exchange, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error when binding queue: %v", err)
	}

	return chConn, queue, nil
}
