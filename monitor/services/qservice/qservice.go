package qservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"

	"github.com/streadway/amqp"
)

type (
	job func(d Delivery) error

	QueueService interface {
		Send(ctx context.Context, m []byte, to string) ([]byte, error)
		Close() error
	}
)

type queueService struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	myq      amqp.Queue
	exchange string
	msgs     <-chan amqp.Delivery
}

var ErrConnection = errors.New("Connection Error")
var ErrUnexpected = errors.New("Something bad was hapened")

const defaultqueue = "tasks"

func NewQueueService(login, password, network, exchange string) (QueueService, error) {
	connstr := fmt.Sprintf("amqp://%s:%s@%s:5672/", login, password, network)
	conn, err := amqp.Dial(connstr)
	if err != nil {
		log.Printf("[failed] try to connection to %s error %s", connstr, err)
		return nil, ErrConnection
	}
	log.Printf("[success] try to connection to %s", connstr)

	ch, err := conn.Channel()
	if err != nil {
		return nil, ErrConnection
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, ErrConnection
	}

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, ErrConnection
	}

	mmsgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, ErrUnexpected
	}

	return &queueService{conn, ch, q, exchange, mmsgs}, nil
}

func (q *queueService) Close() error {
	if err := q.ch.Close(); err != nil {
		return err
	}

	return q.conn.Close()
}

func (q *queueService) Send(ctx context.Context, m []byte, to string) ([]byte, error) {
	uuid, _ := uuid.NewV4()

	q.ch.Publish(
		"task_exchange", // exchange
		to,              // routing key
		false,           // mandatory
		false,
		amqp.Publishing{
			ReplyTo:       q.myq.Name,
			CorrelationId: uuid.String(),
			ContentType:   "application/json",
			Body:          m,
		})

	for d := range q.msgs {
		if d.CorrelationId == uuid.String() {
			return d.Body, nil
		}
	}

	return nil, ErrUnexpected
}
