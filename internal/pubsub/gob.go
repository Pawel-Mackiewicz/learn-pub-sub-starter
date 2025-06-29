package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](exchange, key string, val T) error {
	ch := GetChannel()
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	msg := amqp091.Publishing{
		ContentType: "application/gob",
		Body:        data.Bytes(),
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
	if err != nil {
		return err
	}
	return nil
}
