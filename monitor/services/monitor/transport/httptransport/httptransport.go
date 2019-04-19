package httptransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zhiburt/openbox/monitor/services/monitor/transport/endpoints"
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

	r.Methods("GET").Path("/files/owner/{id}").Handler(kithttp.NewServer(
		points.GetByOwner,
		decodeGetByUserIDRequest,
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

	fmt.Println("---have got file", req.File)

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

func decodeGetByUserIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("cannot find ID in your request, please read documentation before")
	}
	return endpoints.GetByUserIDRequest{ID: id}, nil
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
	if err := getError(response); err != nil {
		errEncoder(ctx, err, w)
		return nil
	}

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

func getError(response interface{}) error {
	if _, ok := reflect.TypeOf(response).FieldByName("Err"); ok {
		r := reflect.ValueOf(response).FieldByName("Err")
		if err, ok := r.Interface().(error); ok {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
