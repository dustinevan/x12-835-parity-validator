// Package middleware wraps every API handler.
package middleware

import (
	"time"

	"go.uber.org/zap"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// Stack returns the middleware every API route uses, outermost first.
func Stack(log *zap.SugaredLogger, timeout time.Duration) []web.Middleware {
	return []web.Middleware{
		Timeout(timeout), // 1. Request deadline (before Logger, so the logged context is inside it)
		Logger(log),      // 2. Start the request values; write one log line at the end
		Errors(),         // 3. Turn a returned error into the JSON error response
		Panics(),         // 4. Recover a panic into a 500 (inside Errors, so it gets a response)
	}
}
