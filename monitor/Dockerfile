FROM golang:alpine

ADD . /go/src/github.com/zhiburt/openbox/monitor
WORKDIR /go/src/github.com/zhiburt/openbox/openbox/monitor

ENTRYPOINT ["go", "run", "main.go"]