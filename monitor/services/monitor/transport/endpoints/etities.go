package endpoints

import "github.com/openbox/monitor/services/monitor"

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	File monitor.File
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID      string
	OwnerID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	File monitor.File `json:"file"`
	Err  error        `json:"error,omitempty"`
}

// GetByUserIDRequest holds the request parameters for the GetByID method.
type GetByUserIDRequest struct {
	ID string
}

// GetByUserIDResponse holds the response values for the GetByID method.
type GetByUserIDResponse struct {
	Files []monitor.File `json:"files"`
	Err   error          `json:"error,omitempty"`
}

// ChangeNameRequest holds the request parameters for the ChangeStatus method.
type ChangeNameRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChangeNameResponse holds the response values for the ChangeStatus method.
type ChangeNameResponse struct {
	Err error `json:"error,omitempty"`
}

// ChangeBodyRequest holds the request parameters for the ChangeStatus method.
type ChangeBodyRequest struct {
	ID   string `json:"id"`
	Body []byte `json:"body"`
}

// ChangeBodyResponse holds the response values for the ChangeStatus method.
type ChangeBodyResponse struct {
	Err error `json:"error,omitempty"`
}

// RemoveByIDRequest holds the request parameters for the ChangeStatus method.
type RemoveByIDRequest struct {
	ID string `json:"id"`
}

// RemoveByIDResponse holds the response values for the ChangeStatus method.
type RemoveByIDResponse struct {
	Err error `json:"error,omitempty"`
}
