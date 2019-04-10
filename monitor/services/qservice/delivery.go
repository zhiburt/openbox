package qservice

import "github.com/streadway/amqp"

type Delivery interface {
	Reply(string, []byte) error
	Body() []byte
	Ack(bool) error
}

func NewDelivery(d amqp.Delivery, ch *amqp.Channel) Delivery {
	return delivery{d, ch}
}

type delivery struct {
	d  amqp.Delivery
	ch *amqp.Channel
}

func (d delivery) Reply(t string, body []byte) error {
	return d.ch.Publish(
		"",          // exchange
		d.d.ReplyTo, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   t,
			CorrelationId: d.d.CorrelationId,
			Body:          body,
		})
}

func (d delivery) Body() []byte {
	return d.d.Body
}

func (d delivery) Ack(b bool) error {
	return d.d.Ack(b)
}
