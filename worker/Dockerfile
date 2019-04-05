FROM golang:alpine

ADD . /go/src/github.com/openbox/worker
WORKDIR /go/src/github.com/openbox/worker   

ENV PORT=3001

CMD ["go", "run", "main.go"]`