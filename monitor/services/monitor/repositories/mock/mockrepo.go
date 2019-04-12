package mock

import (
	"context"

	"github.com/zhiburt/openbox/monitor/services/monitor"
)

// "go.mongodb.org/mongo-driver/mongo"

type MockRepo struct {
	id    string
	files []monitor.File
	err   error
}

func NewRepository(id string, err error, files ...monitor.File) monitor.Repository {
	return &MockRepo{id, files, err}
}

func (rep *MockRepo) CreateFile(ctx context.Context, f monitor.File) (string, error) {
	return rep.id, rep.err
}

func (rep *MockRepo) GetFileByID(ctx context.Context, id string) (monitor.File, error) {
	return rep.files[0], rep.err
}

func (rep *MockRepo) GetFilesByOwner(ctx context.Context, id string) ([]monitor.File, error) {
	return rep.files, rep.err
}

func (rep *MockRepo) ChangeFileName(ctx context.Context, id, newname string) error {
	return rep.err
}

func (rep *MockRepo) ChangeFileBody(ctx context.Context, id string, b []byte) error {
	return rep.err
}

func (rep *MockRepo) RemoveFileByID(ctx context.Context, id string) error {
	return rep.err
}
