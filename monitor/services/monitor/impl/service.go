package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/zhiburt/openbox/monitor/services/monitor"
	"github.com/zhiburt/openbox/openbox/monitor/services/qcommunicator"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	mntr "github.com/zhiburt/openbox/openbox/monitor/services/monitor"
)

// service implements the Order Service
type service struct {
	repository mntr.Repository
	logger     log.Logger
	qfs        qcommunicator.QFileSystem
}

// NewService creates and returns a new Order service instance
func NewService(rep mntr.Repository, logger log.Logger, qfs qcommunicator.QFileSystem) mntr.Service {
	return &service{
		repository: rep,
		logger:     logger,
		qfs:        qfs,
	}
}

// Create makes an order
func (s *service) Create(ctx context.Context, file mntr.File) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	level.Debug(logger).Log("file", fmt.Sprint(file))

	if f, ok := getFile(&file); ok {
		mss, err := s.qfs.PushFile(ctx, f.OwnerID, f.Name, f.Body)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", mntr.ErrInvalidParams
		}

		level.Info(logger).Log("responce from queue service ", mss)
		f.ServerID = mss
	}

	file.Status = "just_created"
	file.CreatedOn = time.Now().Unix()

	id, err := s.repository.CreateFile(ctx, file)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", mntr.ErrRepository
	}
	level.Info(logger).Log("created file with ID", id)

	return id, nil
}

// GetByID returns an order given by id
func (s *service) GetByID(ctx context.Context, id, userid string) (mntr.File, error) {
	logger := log.With(s.logger, "method", "GetByID")

	file, err := s.repository.GetFileByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return file, mntr.ErrCommunication
	}

	return file, nil
}

// GetByOwner returns an order given by id
func (s *service) GetByOwner(ctx context.Context, owner string) ([]mntr.File, error) {
	logger := log.With(s.logger, "method", "GetByOwner")
	files, err := s.repository.GetFilesByOwner(ctx, owner)
	if err != nil {
		level.Error(logger).Log("err", err)
		return files, mntr.ErrRepository
	}

	err = getBodies(s.qfs, files...)
	if err != nil {
		level.Error(logger).Log("err", err)
		return files, mntr.ErrCommunication
	}

	return files, nil
}

// ChangeName returns an order given by id
func (s *service) ChangeName(ctx context.Context, id, newname string) error {
	logger := log.With(s.logger, "method", "ChangeName")

	if err := s.repository.ChangeFileName(ctx, id, newname); err != nil {
		level.Error(logger).Log("err", err)
		return mntr.ErrRepository
	}
	return nil
}

// ChangeName returns an order given by id
func (s *service) ChangeBody(ctx context.Context, id string, body []byte) error {
	logger := log.With(s.logger, "method", "ChangeBody")

	if err := s.repository.ChangeFileBody(ctx, id, body); err != nil {
		level.Error(logger).Log("err", err)
		return mntr.ErrRepository
	}
	return nil
}

// ChangeName returns an order given by id
func (s *service) RemoveByID(ctx context.Context, id string) error {
	logger := log.With(s.logger, "method", "ChangeBody")

	if err := s.repository.RemoveFileByID(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return mntr.ErrRepository
	}
	return nil
}

func getBodies(qs qcommunicator.QFileSystem, f ...monitor.File) error {
	for i := 0; i < len(f); i++ {
		if err := getBody(qs, &f[i]); err != nil {
			return err
		}
	}

	return nil
}

func getBody(qfs qcommunicator.QFileSystem, f *monitor.File) error {
	if f.IsFolder {
		if f.Files != nil {
			for i := 0; i < len(f.Files); i++ {
				fmt.Println("---- RECURSIVE", f)
				getBody(qfs, &f.Files[i])
			}
		}

		return nil
	}

	if f.ServerID == "" {
		fmt.Println("SOMETHING WRONG", f)
		return nil
	}

	b, err := qfs.GetFileBody(context.TODO(), f.OwnerID, f.Name, f.ServerID)
	if err != nil {
		return err
	}

	fmt.Println("---- GETTING BODY", b)
	f.Body = b

	return nil
}

func getFile(f *monitor.File) (*monitor.File, bool) {
	if f.IsFolder {
		if f.Files == nil {
			return nil, false
		}

		return getFile(&f.Files[0])
	}

	return f, true
}
