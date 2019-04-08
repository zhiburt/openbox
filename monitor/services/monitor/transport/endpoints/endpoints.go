package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/openbox/monitor/services/monitor"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Create     endpoint.Endpoint
	GetByID    endpoint.Endpoint
	GetByOwner endpoint.Endpoint
	ChangeName endpoint.Endpoint
	ChangeBody endpoint.Endpoint
	RemoveByID endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s monitor.Service) Endpoints {
	return Endpoints{
		Create:     makeCreateEndpoint(s),
		GetByID:    makeGetByIDEndpoint(s),
		GetByOwner: makeGetByOwnerEndpoint(s),
		ChangeName: makeChangeNameEndpoint(s),
		ChangeBody: makeChangeBodyEndpoint(s),
		RemoveByID: makeRemoveByIDEndpoint(s),
	}
}

func makeCreateEndpoint(s monitor.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.File)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s monitor.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		f, err := s.GetByID(ctx, req.ID, req.OwnerID)
		return GetByIDResponse{File: f, Err: err}, nil
	}
}

func makeGetByOwnerEndpoint(s monitor.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByUserIDRequest)
		f, err := s.GetByOwner(ctx, req.ID)
		return GetByUserIDResponse{Files: f, Err: err}, nil
	}
}

func makeChangeNameEndpoint(s monitor.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeNameRequest)
		err := s.ChangeName(ctx, req.ID, req.Name)
		return ChangeNameResponse{Err: err}, nil
	}
}

func makeChangeBodyEndpoint(s monitor.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeBodyRequest)
		err := s.ChangeBody(ctx, req.ID, req.Body)
		return ChangeBodyResponse{Err: err}, nil
	}
}

func makeRemoveByIDEndpoint(s monitor.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RemoveByIDRequest)
		err := s.RemoveByID(ctx, req.ID)
		return RemoveByIDResponse{Err: err}, nil
	}
}
