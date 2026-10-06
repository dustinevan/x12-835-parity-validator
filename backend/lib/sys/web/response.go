package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondOK sends a 200 with data as JSON.
func RespondOK(ctx context.Context, w http.ResponseWriter, data any) error {
	return Respond(ctx, w, data, http.StatusOK)
}

// RespondNoContent sends a 204 with no body.
func RespondNoContent(ctx context.Context, w http.ResponseWriter) error {
	return Respond(ctx, w, nil, http.StatusNoContent)
}

// Respond records the status for the log line and sends data as JSON, or
// only the status when data is nil. The body ends with a newline, as
// json.Encoder writes it.
func Respond(ctx context.Context, w http.ResponseWriter, data any, status int) error {
	SetStatusCode(ctx, status)
	if data == nil {
		w.WriteHeader(status)
		return nil
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("web: encoding response: %w", err)
	}
	body = append(body, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		return NewResponseWriteError(err)
	}
	return nil
}
