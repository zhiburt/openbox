FROM golang:alpine

ADD . /go/src/github.com/openbox/monitor
WORKDIR /go/src/github.com/openbox/monitor

ENTRYPOINT ["go", "run", "main.go"]