package main

import (
	"encoding/json"
	"fmt"
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

	err = ch.ExchangeDeclare(
		"task_exchange", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to register a exchange")

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

	if len(os.Args) > 5 && os.Args[5] != "all" {
		fmt.Println("LOOKING UP", os.Args[5])
		err = ch.Publish(
			"task_exchange", // exchange
			os.Args[5],      // routing key
			false,           // mandatory
			false,
			amqp.Publishing{
				ReplyTo:       q.Name,
				CorrelationId: uuid.String(),
				ContentType:   "application/json",
				Body:          body,
			})
		failOnError(err, "Failed to publish a message")
	} else {
		err = ch.Publish(
			"task_exchange", // exchange
			"",              // routing key
			false,           // mandatory
			false,
			amqp.Publishing{
				ReplyTo:       q.Name,
				CorrelationId: uuid.String(),
				ContentType:   "application/json",
				Body:          body,
			})
		failOnError(err, "Failed to publish a message")
	}

	for d := range msgs {
		if uuid.String() == d.CorrelationId {
			log.Printf(" [x] Resp from server %s", d.Body)
			break
		}
		log.Printf(" [x] from server but invalid %s", d.Body)
	}

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) []byte {
	m := message(args[1], args[2], args[3], args[4])

	if len(args) > 6 {
		n := strings.Split(args[6], ".")
		m.NewName = n[0]
		m.NewExtension = n[1]
	}

	fmt.Println(m)
	b, _ := json.Marshal(m)
	return b
}

func message(user, name, body, t string) communication.Message {
	n := strings.Split(name, ".")
	mss := communication.Message{}
	mss.Name = n[0]
	mss.Extension = n[1]
	mss.Body = []byte(body)
	mss.UserID = user
	mss.Type = t
	return mss
}

func pop(args []string) (string, []string) {
	return args[0], args[1:]
}
