package impl

import (
	"context"
	"time"

	"github.com/openbox/monitor/services/qservice"
	comm "github.com/openbox/monitor/services/qservice/communication"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	mntr "github.com/openbox/monitor/services/monitor"
)

// service implements the Order Service
type service struct {
	repository mntr.Repository
	logger     log.Logger
	queue      qservice.QueueService
}

// NewService creates and returns a new Order service instance
func NewService(rep mntr.Repository, logger log.Logger, qs qservice.QueueService) mntr.Service {
	return &service{
		repository: rep,
		logger:     logger,
		queue:      qs,
	}
}

// Create makes an order
func (s *service) Create(ctx context.Context, file mntr.File) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	file.Status = "just_created"
	file.CreatedOn = time.Now().Unix()

	id, err := s.repository.CreateFile(ctx, file)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", mntr.ErrRepository
	}
	level.Info(logger).Log("created file with ID", id)

	mss, err := comm.Marshal(comm.NewMessage(comm.CREATE, id, file.Name, file.Body))
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", mntr.ErrInvalidParams
	}

	mss, err = s.queue.Send(ctx, mss, "")
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", mntr.ErrCommunication
	}
	level.Info(logger).Log("responce from queue service ", mss)

	return id, nil
}

// GetByID returns an order given by id
func (s *service) GetByID(ctx context.Context, id, userid string) (mntr.File, error) {
	logger := log.With(s.logger, "method", "GetByID")

	file, err := s.repository.GetFileByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return file, mntr.ErrRepository
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
