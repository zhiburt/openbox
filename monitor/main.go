package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openbox/monitor/qservice"

	"github.com/openbox/monitor/communication"
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

	body := bodyFrom(os.Args)

	qs, err := qservice.NewQueueService("guest", "guest", net, "task_exchange")
	if err != nil {
		log.Println("[erorr] happend", err)
		os.Exit(1)
	}

	if len(os.Args) > 5 {
		body, err = qs.Send(context.Background(), body, os.Args[5])
	} else {
		body, err = qs.Send(context.Background(), body, "")
	}

	if err != nil {
		log.Println("[erorr] happend", err)
		os.Exit(1)
	}
	log.Printf(" [x] responce from server %s", body)
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
