// Package web is the HTTP plumbing every handler shares: a handler returns
// an error instead of writing one, middleware wraps handlers, a router
// mounts them on a ServeMux, request-scoped values collect what the
// request's one log line reports, and small helpers read requests and
// write responses and errors.
package web

import (
	"context"
	"net/http"
)

// Handler handles one request. It writes a successful response itself (with
// the Respond helpers) and returns an error for anything else; the Errors
// middleware turns that error into the response.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Middleware wraps a handler with behavior that runs before and after it.
type Middleware func(handler Handler) Handler

// WrapMiddleware wraps the handler in each middleware. The first middleware
// in the list is the outermost: its code runs first and last on a request.
func WrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		handler = mw[i](handler)
	}
	return handler
}
