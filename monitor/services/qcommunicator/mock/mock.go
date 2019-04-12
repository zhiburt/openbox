package mock

import (
	"context"

	"github.com/openbox/monitor/services/qcommunicator"
)

type qcommunicatorMock struct {
	resultBytes []byte
	resultError error
}

func NewQcommunicatorMock(b []byte, err error) qcommunicator.QFileSystem {
	return qcommunicatorMock{b, err}
}

func (mock qcommunicatorMock) GetFileBody(ctx context.Context, userid, name, to string) ([]byte, error) {
	return mock.resultBytes, mock.resultError
}

func (mock qcommunicatorMock) PushFile(ctx context.Context, userid, name string, body []byte) (string, error) {
	return string(mock.resultBytes), mock.resultError
}
