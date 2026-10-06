package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// Logger starts the request's values and writes one log line when the
// handler returns: Error level when anything failed, Info otherwise. The
// request ID goes back in the X-Request-ID header and in every error body,
// so a message on screen can be matched to its log line.
//
// The line holds the matched route pattern, never the URL path, query
// string, headers, or body. A client error's message is not logged either
// (see Errors), so no text a client typed reaches the log.
func Logger(log *zap.SugaredLogger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			requestID := web.GenerateRequestID()
			_, route, _ := strings.Cut(r.Pattern, " ")
			ctx = web.WithValues(ctx, web.Values{
				RequestID: requestID,
				Start:     time.Now(),
				Method:    r.Method,
				Route:     route,
				IPAddress: web.GetIPAddress(r),
			})
			w.Header().Set("X-Request-ID", requestID)

			// Errors returns only a failure to write the response.
			writeErr := handler(ctx, w, r)

			v := web.GetValues(ctx)
			fields := []any{
				"requestID", v.RequestID,
				"method", v.Method,
				"route", v.Route,
				"statusCode", v.StatusCode,
				"dur", time.Since(v.Start).String(),
				"ipAddr", v.IPAddress,
			}
			failed := false
			if v.Error != nil {
				fields = append(fields, "error", v.Error.Error())
				failed = true
			}
			if len(v.LogOnlyErrors) > 0 {
				var texts []string
				for _, err := range v.LogOnlyErrors {
					texts = append(texts, err.Error())
				}
				fields = append(fields, "logOnlyErrors", texts)
				failed = true
			}
			if writeErr != nil {
				fields = append(fields, "responseWriteError", writeErr.Error())
				failed = true
			}
			if v.PanicError != nil {
				fields = append(fields, "panicError", v.PanicError.Error())
				failed = true
			}
			if v.StatusCode == 0 {
				v.Messages = append(v.Messages, "request did not record a status code")
			}
			if len(v.Messages) > 0 {
				fields = append(fields, "messages", v.Messages)
			}

			if failed {
				log.Errorw("request", fields...)
			} else {
				log.Infow("request", fields...)
			}
			return nil
		}
	}
}
