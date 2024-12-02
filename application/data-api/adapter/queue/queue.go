package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type queue struct {
	conn    *amqp.Connection
	queue   amqp.Queue
	channel *amqp.Channel
}

func New(host, port, name string) (*queue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &queue{conn: conn, channel: ch, queue: q}, nil
}

func (q *queue) Send(msg string) error {
	err := q.channel.PublishWithContext(
		context.Background(),
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return err
	}

	return nil
}

func (q *queue) Close() error {
	chErr := q.channel.Close()
	connErr := q.conn.Close()

	if chErr != nil {
		return chErr
	}
	if connErr != nil {
		return connErr
	}

	return nil
}
