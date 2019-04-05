package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openbox/worker/communication"

	"github.com/openbox/worker/filesystem"
	"github.com/streadway/amqp"
)

func init() {
	if n := os.Getenv("NETWORKNAME"); n != "" {
		log.Println("network", n)
		net = n
	}
	if n := os.Getenv("SERVER_NAME"); n != "" {
		servername = n
	}
	if n := os.Getenv("ROOT"); n != "" {
		root = n
	}

}

var (
	net        = "localhost"
	servername = "rabbitmqworker"
	root       = "."
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@" + net + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

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

	forever := make(chan bool)

	var fs = filesystem.NewFilesystem(root)

	go func() {
		for d := range msgs {
			m := &communication.Message{}
			json.Unmarshal(d.Body, m)

			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")

			if m.Type != "" {
				log.Println("look up")
				f, err := fs.Lookup(filesystem.NewUser(m.UserID), filesystem.NewFile(m.Name, m.Extension, bytes.NewReader(m.Body)))
				if err != nil {
					log.Println("didn't find", err)
					err = ch.Publish(
						"",     // exchange
						q.Name, // routing key
						false,  // mandatory
						false,  // immediate
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          d.Body,
							ReplyTo:       d.ReplyTo,
						})
					continue
				}
				log.Println("Found")

				content, _ := ioutil.ReadAll(f.Body())
				log.Println("content", content)

				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          append(content, []byte(" "+servername)...),
					})
			} else {
				log.Println("Create")
				err = fs.Create(filesystem.NewUser(m.UserID), filesystem.NewFile(m.Name, m.Extension, bytes.NewReader(m.Body)))
				failOnError(err, "Failed to create a file")

				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(servername),
					})
			}
			log.Println("wait")
		}
		log.Printf("CLOSED CONN")
		os.Exit(1)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		forever <- true
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
