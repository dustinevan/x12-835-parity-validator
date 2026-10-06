package web

import (
	"context"
	"time"
)

type key int

const valuesKey key = 1

// Values is what one request's log line reports. The Logger middleware
// starts it and writes the line when the handler returns; handlers and the
// other middleware add to it through the context with the functions below.
//
// It holds the matched route pattern, never the URL path, query string,
// headers, or body, so no text a client typed reaches the log: a path
// segment such as {filename} appears as the pattern names it. It carries
// no raw path, user agent, auth claims, or SQL. Handlers here start no
// goroutines, so Values has no lock.
type Values struct {
	RequestID     string
	Start         time.Time
	Method        string
	Route         string // the ServeMux pattern path, e.g. "/api/diffs/{filename}"
	StatusCode    int
	IPAddress     string
	Error         error
	LogOnlyErrors []error
	PanicError    error
	Messages      []string
}

// WithValues starts a request's values. Only the Logger middleware calls it.
func WithValues(ctx context.Context, v Values) context.Context {
	return context.WithValue(ctx, valuesKey, &v)
}

// GetValues returns the request's values as they stand, or zero Values
// outside a request (a handler called directly in a test).
func GetValues(ctx context.Context) Values {
	if v, ok := ctx.Value(valuesKey).(*Values); ok {
		return *v
	}
	return Values{}
}

// SetStatusCode records the status written, for the log line. The Respond
// helpers and the Errors middleware call it.
func SetStatusCode(ctx context.Context, statusCode int) {
	if v, ok := ctx.Value(valuesKey).(*Values); ok {
		v.StatusCode = statusCode
	}
}

// SetError records the error a response reported. The Errors middleware
// calls it for server errors only.
func SetError(ctx context.Context, err error) {
	if v, ok := ctx.Value(valuesKey).(*Values); ok && err != nil {
		v.Error = err
	}
}

// AddLogOnlyError records a cause for the log that the client does not see,
// such as the database error behind a 500.
func AddLogOnlyError(ctx context.Context, err error) {
	if v, ok := ctx.Value(valuesKey).(*Values); ok && err != nil {
		v.LogOnlyErrors = append(v.LogOnlyErrors, err)
	}
}

// SetPanicError records a recovered panic and its stack.
func SetPanicError(ctx context.Context, err error) {
	if v, ok := ctx.Value(valuesKey).(*Values); ok && err != nil {
		v.PanicError = err
	}
}

// AddMessage records a note for the log line. Messages are server text;
// never format a request value into one.
func AddMessage(ctx context.Context, message string) {
	if v, ok := ctx.Value(valuesKey).(*Values); ok && message != "" {
		v.Messages = append(v.Messages, message)
	}
}
