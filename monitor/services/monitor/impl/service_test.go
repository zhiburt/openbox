package impl

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/zhiburt/openbox/monitor/services/qcommunicator"

	"github.com/zhiburt/openbox/openbox/monitor/services/monitor"

	"github.com/zhiburt/openbox/openbox/monitor/services/monitor/repositories/mock"

	"github.com/go-kit/kit/log"
	qmock "github.com/zhiburt/openbox/openbox/monitor/services/qcommunicator/mock"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		qfs         qcommunicator.QFileSystem
		repo        monitor.Repository
		file        monitor.File
		expectedStr string
		expectedErr error
	}{
		{
			repo:        mock.NewRepository("id__for__file", nil, monitor.File{}),
			qfs:         qmock.NewQcommunicatorMock(nil, nil),
			file:        monitor.File{},
			expectedErr: nil,
			expectedStr: "id__for__file",
		},
		{
			repo:        mock.NewRepository("", fmt.Errorf("some_error"), monitor.File{}),
			qfs:         qmock.NewQcommunicatorMock(nil, nil),
			file:        monitor.File{},
			expectedErr: monitor.ErrRepository,
			expectedStr: "",
		},
		{
			repo:        mock.NewRepository("", nil, monitor.File{}),
			qfs:         qmock.NewQcommunicatorMock(nil, fmt.Errorf("some_error")),
			file:        monitor.File{},
			expectedErr: monitor.ErrInvalidParams,
			expectedStr: "",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			serv := NewService(c.repo, log.NewNopLogger(), c.qfs)

			id, err := serv.Create(context.Background(), c.file)
			if err != c.expectedErr {
				t.Errorf("expected erorr compared with %v\nbut was %v", c.expectedErr, err)
			}
			if id != c.expectedStr {
				t.Errorf("expected id compared with %v\nbut was %v", c.expectedStr, id)
			}
		})
	}
}

func TestGetByOwner(t *testing.T) {
	cases := []struct {
		qfs           qcommunicator.QFileSystem
		repo          monitor.Repository
		owner         string
		expectedFiles []monitor.File
		expectedErr   error
	}{
		{
			repo:          mock.NewRepository("", nil, monitor.File{OwnerID: "maxim", ServerID: "NOT_EMPTY"}),
			qfs:           qmock.NewQcommunicatorMock([]byte("hello world"), nil),
			owner:         "maxim",
			expectedErr:   nil,
			expectedFiles: []monitor.File{{OwnerID: "maxim", Body: []byte("hello world"), ServerID: "NOT_EMPTY"}},
		},
		{
			repo:          mock.NewRepository("", nil, monitor.File{OwnerID: "maxim"}),
			qfs:           qmock.NewQcommunicatorMock([]byte("hello world"), nil),
			owner:         "maxim",
			expectedErr:   nil,
			expectedFiles: []monitor.File{{OwnerID: "maxim", Body: nil}},
		},
		{
			repo:          mock.NewRepository("", nil, monitor.File{OwnerID: "maxim"}),
			qfs:           qmock.NewQcommunicatorMock(nil, nil),
			owner:         "maxim",
			expectedErr:   nil,
			expectedFiles: []monitor.File{{OwnerID: "maxim", Body: nil}},
		},
		{
			repo:          mock.NewRepository("", nil, monitor.File{}),
			qfs:           qmock.NewQcommunicatorMock(nil, nil),
			owner:         "",
			expectedErr:   nil,
			expectedFiles: []monitor.File{{}},
		},
		{
			repo:          mock.NewRepository("", nil, monitor.File{OwnerID: "maxim", ServerID: "NOT_EMPTY"}),
			qfs:           qmock.NewQcommunicatorMock(nil, fmt.Errorf("some_error")),
			owner:         "maxim",
			expectedErr:   monitor.ErrCommunication,
			expectedFiles: []monitor.File{monitor.File{OwnerID: "maxim", ServerID: "NOT_EMPTY"}},
		},
		{
			repo:          mock.NewRepository("", fmt.Errorf("some_error"), monitor.File{}),
			qfs:           qmock.NewQcommunicatorMock(nil, nil),
			owner:         "",
			expectedErr:   monitor.ErrRepository,
			expectedFiles: []monitor.File{{}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			serv := NewService(c.repo, log.NewNopLogger(), c.qfs)

			files, err := serv.GetByOwner(context.Background(), c.owner)
			if err != c.expectedErr {
				t.Errorf("expected erorr compared with %v\nbut was %v", c.expectedErr, err)
			}
			if !reflect.DeepEqual(files, c.expectedFiles) {
				t.Errorf("expected files compared with %v\nbut was %v", c.expectedFiles, files)
			}
		})
	}
}
