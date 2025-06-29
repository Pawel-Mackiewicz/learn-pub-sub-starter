package pubsub

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	QueueTypeDurable   SimpleQueueType = "durable"
	QueueTypeTransient SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {

	chConn, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error while creating channel: %v", err)
	}
	var durable bool
	var autoDelete bool
	var exclusive bool

	switch queueType {
	case QueueTypeDurable:
		durable = true
	case QueueTypeTransient:
		autoDelete = true
		exclusive = true
	default:
		return nil, amqp.Queue{}, fmt.Errorf("wrong `queueType` selected")
	}

	table := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}
	queue, err := chConn.QueueDeclare(queueName, durable, autoDelete, exclusive, false, table)
	if err != nil {
		_ = chConn.Close()
		return nil, amqp.Queue{}, fmt.Errorf("error when declaring queue: %v", err)
	}

	err = chConn.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error when binding queue: %v", err)
	}

	return chConn, queue, nil
}
