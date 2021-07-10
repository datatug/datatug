package endpoints

import (
	"context"
	"net/http"
)

// ProjectEndpoints defines project endpoints
type ProjectEndpoints interface {
	CreateProject(w http.ResponseWriter, r *http.Request)
	DeleteProject(w http.ResponseWriter, r *http.Request)
}

// GetContextFromRequest creates context for a given request
// We need a request for example to decode Firebase token
// Probably can be done better then being hard-depending on an *http.Request
var GetContextFromRequest = func(r *http.Request) (context.Context, error) {
	return r.Context(), nil
}

// RequestDTO defines an interface that should be implemented by request DTO struct
type RequestDTO interface {
	Validate() error
}

// ResponseDTO common interface for response objects
type ResponseDTO interface {
	// Validate validates response
	Validate() error
}

// VerifyRequestOptions - options for request verification
type VerifyRequestOptions interface { // TODO: move to shared Sneat package
	MinimumContentLength() int64
	MaximumContentLength() int64
	AuthenticationRequired() bool
}

var _ VerifyRequestOptions = (*VerifyRequest)(nil)

// VerifyRequest implements VerifyRequestOptions
type VerifyRequest struct { // TODO: move to shared Sneat package
	MinContentLength int64
	MaxContentLength int64
	AuthRequired     bool
}

// MinimumContentLength defines min content length
func (v VerifyRequest) MinimumContentLength() int64 {
	return v.MinContentLength
}

// MaximumContentLength defines max content length, if < 0 no limit
func (v VerifyRequest) MaximumContentLength() int64 {
	return v.MaxContentLength
}

// AuthenticationRequired specifies if authentication is mandatory
func (v VerifyRequest) AuthenticationRequired() bool {
	return v.AuthRequired
}

// Handler is responsible for creating context and call `handler()` func that should use
// provided context along with `requestDTO` that was populated from request body
// Its is exposed publicly so it can be replaced with custom implementation
type Handler = func(
	w http.ResponseWriter,
	r *http.Request,
	requestDTO RequestDTO,
	verifyOptions VerifyRequestOptions,
	statusCode int,
	handler func(ctx context.Context) (responseDTO ResponseDTO, err error),
)

// SetHandler sets handler
func SetHandler(handler Handler) {
	if handler == nil {
		panic("handler is not provided")
	}
	handle = handler
}

var handle Handler = func(w http.ResponseWriter, r *http.Request, requestDTO RequestDTO, verifyOptions VerifyRequestOptions, statusCode int, handler func(ctx context.Context) (responseDTO ResponseDTO, err error)) {
	panic("not initialized properly")
}