package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/openbox/monitor/communication"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	var net = "localhost"
	if n := os.Getenv("NETWORKNAME"); n != "" {
		log.Println("network", n)
		net = n
	}

	conn, err := amqp.Dial("amqp://guest:guest@" + net + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	uuid, _ := uuid.NewV4()
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",           // exchange
		"task_queue", // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			ReplyTo:       q.Name,
			CorrelationId: uuid.String(),
			ContentType:   "text/plain",
			Body:          body,
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if uuid.String() == d.CorrelationId {
			log.Printf(" [x] Resp from server %s", d.Body)
			break
		}
	}

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) []byte {
	var name string
	var user string
	var body string
	if (len(args) < 2) || os.Args[1] == "" {
		name = "hello.clead"
		user = "maxim"
		body = "hello world"
	} else {
		name = args[2]
		user = args[1]
		body = args[3]
	}
	b, _ := json.Marshal(message(user, name, body))
	return b
}

func message(user, name, body string) communication.Message {
	n := strings.Split(name, ".")
	return communication.Message{user, []byte(body), n[0], n[1]}
}
