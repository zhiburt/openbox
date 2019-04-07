package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/openbox/worker/commands"
	"github.com/openbox/worker/communication"
	"github.com/openbox/worker/qservice"

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

		qs.Handle(context.Background(), handlefunction(fs))
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

func handlefunction(fs filesystem.Filesystem) qservice.Job {
	return func(d qservice.Delivery) error {
		m := &communication.Message{}
		json.Unmarshal(d.Body(), m)

		log.Printf("Received a message: %s\n", d.Body())

		command, err := commands.NewCommand(fs, *m)
		if err != nil {
			log.Println("[error] with creating command", err)
			return err
		}

		mss, err := command(fs, *m)
		if err != nil {
			log.Println("[error] with command", err)
			return err
		}

		if mss == nil {
			mss = []byte(servername)
		}

		return d.Reply("text/plain", mss)
	}
}
