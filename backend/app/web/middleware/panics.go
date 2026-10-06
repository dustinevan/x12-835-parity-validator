package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// Panics recovers a panic in the handler, records it with its stack for the
// log line, and returns a 500 for the Errors middleware to write.
func Panics() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					web.SetPanicError(ctx, fmt.Errorf("panic: %v\n%s", rec, debug.Stack()))
					err = web.InternalServerErrorf("%s", http.StatusText(http.StatusInternalServerError))
				}
			}()
			return handler(ctx, w, r)
		}
	}
}
