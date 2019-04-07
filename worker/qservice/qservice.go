package qservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type (
	Job func(d Delivery) error

	QueueService interface {
		Handle(context.Context, Job) error
		Close() error
	}
)

type queueService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	myq  amqp.Queue
	dq   amqp.Queue
}

var ErrConnection = errors.New("Connection Error")
var ErrUnexpected = errors.New("Something bad was hapened")

const defaultqueue = "tasks"

func NewQueueService(login, password, network, exchange, servername string) (QueueService, error) {
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
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, ErrUnexpected
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, ErrUnexpected
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
		return nil, ErrUnexpected
	}

	err = ch.QueueBind(
		q.Name,     // queue name
		servername, // routing key
		exchange,   // exchange
		false,
		nil)
	if err != nil {
		return nil, ErrUnexpected
	}

	//second queue

	q1, err := ch.QueueDeclare(
		defaultqueue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, ErrUnexpected
	}

	err = ch.QueueBind(
		q1.Name,  // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, ErrUnexpected
	}

	return &queueService{conn, ch, q, q1}, nil
}

func (q *queueService) Handle(ctx context.Context, j Job) error {
	mmsgs, err := q.ch.Consume(
		q.myq.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return ErrUnexpected
	}
	dmsgs, err := q.ch.Consume(
		q.dq.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return ErrUnexpected
	}

	for {
		select {
		case d := <-mmsgs:
			j(NewDelivery(d, q.ch))
		case d := <-dmsgs:
			j(NewDelivery(d, q.ch))
		case <-ctx.Done():
			return nil
		}
	}
}

func (q *queueService) Close() error {
	if err := q.ch.Close(); err != nil {
		return err
	}

	return q.conn.Close()
}
