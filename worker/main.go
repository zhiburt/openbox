package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/openbox/worker/qservice"

	"github.com/openbox/worker/communication"

	"github.com/openbox/worker/filesystem"
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
	forever := make(chan bool)

	var fs = filesystem.NewFilesystem(root)
	qs, err := qservice.NewQueueService("guest", "guest", net, "task_exchange", servername)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		var foo = func(d qservice.Delivery) error {
			m := &communication.Message{}
			json.Unmarshal(d.Body(), m)

			log.Printf("Received a message: %s\n", d.Body())

			if m.Type != "" {
				log.Println("look up")
				f, err := fs.Lookup(filesystem.NewUser(m.UserID), filesystem.NewFile(m.Name, m.Extension, bytes.NewReader(m.Body)))
				if err != nil {
					return d.Ack(false)
				}
				log.Println("Found")

				content, _ := ioutil.ReadAll(f.Body())
				log.Println("content", content)

				d.Reply("text/plain", append(content, []byte(" "+servername)...))
			} else {
				log.Println("Create")
				err = fs.Create(filesystem.NewUser(m.UserID), filesystem.NewFile(m.Name, m.Extension, bytes.NewReader(m.Body)))
				failOnError(err, "Failed to create a file")

				d.Reply("text/plain", []byte(servername))
			}
			return nil
		}

		qs.Handle(context.Background(), foo)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
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
