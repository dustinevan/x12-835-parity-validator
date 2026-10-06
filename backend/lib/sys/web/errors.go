package web

import (
	"errors"
	"fmt"
	"net/http"
)

// RequestError pairs an error with the HTTP status to answer with. Its
// message is written to the client and not to the log, so it may name what
// the client sent; an internal cause goes to AddLogOnlyError instead.
type RequestError struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return RequestError{Err: err, Status: status}
}

func BadRequestErrorf(format string, args ...any) error {
	return NewRequestError(fmt.Errorf(format, args...), http.StatusBadRequest)
}

func NotFoundErrorf(format string, args ...any) error {
	return NewRequestError(fmt.Errorf(format, args...), http.StatusNotFound)
}

func InternalServerErrorf(format string, args ...any) error {
	return NewRequestError(fmt.Errorf(format, args...), http.StatusInternalServerError)
}

func (re RequestError) Error() string {
	return re.Err.Error()
}

func (re RequestError) Unwrap() error {
	return re.Err
}

// AsRequestError finds a RequestError in err's chain.
func AsRequestError(err error) (RequestError, bool) {
	var re RequestError
	ok := errors.As(err, &re)
	return re, ok
}

// ResponseWriteError means writing the response failed after it started,
// usually because the client went away. Nothing more can be sent; the
// Logger middleware records it.
type ResponseWriteError struct {
	Err error
}

func NewResponseWriteError(err error) error {
	return ResponseWriteError{Err: err}
}

func (e ResponseWriteError) Error() string {
	return "web: writing response: " + e.Err.Error()
}

func (e ResponseWriteError) Unwrap() error {
	return e.Err
}

// AsResponseWriteError finds a ResponseWriteError in err's chain.
func AsResponseWriteError(err error) (ResponseWriteError, bool) {
	var rwe ResponseWriteError
	ok := errors.As(err, &rwe)
	return rwe, ok
}
