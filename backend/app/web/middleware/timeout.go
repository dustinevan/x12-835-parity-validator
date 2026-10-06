package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// Timeout gives each request a deadline. http.Server's WriteTimeout is
// off and the deadline lives here, where a route group could be given a
// different one.
func Timeout(timeout time.Duration) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return handler(ctx, w, r)
		}
	}
}
