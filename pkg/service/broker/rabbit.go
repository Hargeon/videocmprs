// Package broker uses for implementation message broker
package broker

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Rabbit represent rabbitmq client
type Rabbit struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue

	dbURL string
}

// NewRabbit initialize Rabbit
func NewRabbit(dbURL string) *Rabbit {
	return &Rabbit{
		dbURL: dbURL,
	}
}

// Connect to rabbit, create chan, declare queue
func (r *Rabbit) Connect(queueName string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(r.dbURL)

	if err != nil {
		return nil, err
	}

	r.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	r.ch = ch

	q, err := ch.QueueDeclare(
		queueName,
		true, // message will not lose if rabbit crashed
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, err
	}

	r.q = q

	return conn, err
}

// Publish body to rabbit
func (r *Rabbit) Publish(body []byte) error {
	return r.ch.Publish(
		"",
		r.q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // message will not lose if rabbit crashed
			ContentType:  "application/json",
			Body:         body,
		})
}

// Consume return channel which read messages from rabbit
func (r *Rabbit) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.q.Name,
		"",
		false, // needs to mark a message was processed
		false,
		false,
		false,
		nil)

	return msgs, err
}

func (r *Rabbit) Ping() error {
	rabbitDsn := fmt.Sprintf(r.dbURL)
	conn, err := amqp.Dial(rabbitDsn)

	if err != nil {
		return err
	}

	defer conn.Close()

	return nil
}
