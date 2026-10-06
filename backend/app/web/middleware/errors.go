package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// ErrorResponse is the body of every error the handlers return. The
// request ID is a string because a nanosecond timestamp does not fit a
// JavaScript number exactly.
type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Error     string `json:"error"`
}

// statusClientClosedRequest is recorded in the log line, never sent, when
// the client closed the request before a response (nginx uses 499 too).
const statusClientClosedRequest = 499

// Errors turns an error returned by a handler into the response:
//
//   - a web.RequestError answers with its status and its message;
//   - any other error answers 500 "Internal Server Error", and the error
//     itself goes to the log only;
//   - a cancelled request answers nothing, since no one is waiting;
//   - a failed response write is passed up to the Logger to record.
//
// Server errors (500 and up) are recorded as the log line's error. A client
// error's message goes to the client only: it can name what the client
// sent, and the log stays free of request values by construction.
func Errors() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := handler(ctx, w, r)
			if err == nil {
				return nil
			}
			if _, ok := web.AsResponseWriteError(err); ok {
				return err
			}
			if errors.Is(err, context.Canceled) {
				web.SetStatusCode(ctx, statusClientClosedRequest)
				web.AddMessage(ctx, "request cancelled before a response: the client closed it or the server is shutting down")
				return nil
			}

			status := http.StatusInternalServerError
			message := http.StatusText(http.StatusInternalServerError)
			if re, ok := web.AsRequestError(err); ok {
				status = re.Status
				message = re.Error()
			}
			if status >= http.StatusInternalServerError {
				web.SetError(ctx, err)
			}
			return web.Respond(ctx, w, ErrorResponse{
				RequestID: web.GetValues(ctx).RequestID,
				Error:     message,
			}, status)
		}
	}
}
