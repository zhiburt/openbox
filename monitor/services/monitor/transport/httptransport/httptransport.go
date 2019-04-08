package httptransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openbox/monitor/services/monitor/transport/endpoints"
)

func NewService(points endpoints.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(errEncoder),
	}
	r.Methods("POST").Path("/files").Handler(kithttp.NewServer(
		points.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/files/{id}").Handler(kithttp.NewServer(
		points.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.CreateRequest
	if e := json.NewDecoder(r.Body).Decode(&req.File); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("cannot find ID in your request, please read documentation before")
	}
	return endpoints.GetByIDRequest{ID: id}, nil
}

func decodeChangeNameRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.ChangeNameRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		errEncoder(ctx, e, w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(struct {
		err error
	}{err: err})

	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}