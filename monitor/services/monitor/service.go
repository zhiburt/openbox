package monitor

import (
	"context"
	"errors"
)

// Service describes the Monitor service for some FileSystem.
type Service interface {
	Create(ctx context.Context, f File) (string, error)
	GetByID(ctx context.Context, id, userid string) (File, error)
	GetByOwner(ctx context.Context, owner string) ([]File, error)
	ChangeName(ctx context.Context, id, newname string) error
	ChangeBody(ctx context.Context, id string, body []byte) error
	RemoveByID(ctx context.Context, id string) error
}

// File represents an file
type File struct {
	ID        string `json:"id,omitempty"`
	OwnerID   string `json:"owner_id"`
	Name      string `json:"name"`
	Body      []byte `json:"body"`
	Status    string `json:"status"`
	CreatedOn int64  `json:"created_on,omitempty"`
	IsFolder  bool   `json:"is_folder,omitempty"`
	Files     []File `json:"files,omitempty"`
}

/*
curl --header "Content-Type: application/json" --request POST --data '{"name":"maxim.cpp","owner_id":"helloworld"}' http://localhost:8082/files

*/

// Repository describes the persistence on order model
type Repository interface {
	CreateFile(ctx context.Context, f File) (string, error)
	GetFileByID(ctx context.Context, id string) (File, error)
	GetFilesByOwner(ctx context.Context, owner string) ([]File, error)
	ChangeFileName(ctx context.Context, id, newname string) error
	ChangeFileBody(ctx context.Context, id string, body []byte) error
	RemoveFileByID(ctx context.Context, id string) error
}

var ErrRepository = errors.New("error happend in repository")
var ErrCommunication = errors.New("error connect with communication with services")
var ErrInvalidParams = errors.New("try to figure out why")
