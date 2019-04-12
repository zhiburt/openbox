package qcommunicator

import (
	"context"
	"fmt"

	comm "github.com/zhiburt/openbox/monitor/services/qcommunicator/communication"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/zhiburt/openbox/monitor/services/qservice"
)

type QFileSystem interface {
	GetFileBody(ctx context.Context, userid, name, to string) ([]byte, error)
	PushFile(ctx context.Context, userid, name string, body []byte) (string, error)
}

func NewQFileSystem(logger log.Logger, queue qservice.QueueService) QFileSystem {
	return &defaultQFileSystem{
		logger: logger,
		queue:  queue,
	}
}

type defaultQFileSystem struct {
	logger log.Logger
	queue  qservice.QueueService
}

func (qfs *defaultQFileSystem) GetFileBody(ctx context.Context, userid, name, to string) ([]byte, error) {
	mss, err := comm.Marshal(comm.NewMessage(comm.LOOKUP, userid, name, nil))
	if err != nil {
		return nil, fmt.Errorf("cannot marshal message %v", err)
	}

	fmt.Println("SEND TO ... ", to)

	mss, err = qfs.queue.Send(context.Background(), mss, to)
	if err != nil {
		return nil, fmt.Errorf("cannot send message %v", err)
	}

	return mss, nil
}

func (qfs *defaultQFileSystem) PushFile(ctx context.Context, userid, name string, body []byte) (string, error) {
	mss, err := comm.Marshal(comm.NewMessage(comm.CREATE, userid, name, body))
	if err != nil {
		level.Error(qfs.logger).Log("err", err)
		return "", err
	}

	mss, err = qfs.queue.Send(ctx, mss, "")
	if err != nil {
		level.Error(qfs.logger).Log("err", err)
		return "", err
	}

	return string(mss), nil
}
