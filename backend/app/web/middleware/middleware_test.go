package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// serve runs one request through the standard stack on a real router, so
// the route pattern is set as it is in the service, and returns the
// response and the log line.
func serve(ctx context.Context, t *testing.T, target string, handler web.Handler) (*httptest.ResponseRecorder, observer.LoggedEntry) {
	t.Helper()
	core, logs := observer.New(zapcore.InfoLevel)
	log := zap.New(core).Sugar()
	router := web.NewRouter(log)
	router.GET("/api/claim", handler, Stack(log, time.Minute)...)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, target, nil).WithContext(ctx)
	router.Handler().ServeHTTP(w, r)
	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("log lines = %d, want 1", len(entries))
	}
	return w, entries[0]
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) ErrorResponse {
	t.Helper()
	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding error body %q: %v", w.Body.String(), err)
	}
	if body.RequestID == "" || body.RequestID != w.Header().Get("X-Request-ID") {
		t.Errorf("requestId = %q, X-Request-ID = %q, want the same non-empty ID", body.RequestID, w.Header().Get("X-Request-ID"))
	}
	return body
}

// TestClientErrorNotLogged: a 400 sends its message to the client; the log
// line has the route and status but no message, path, or query value.
func TestClientErrorNotLogged(t *testing.T) {
	handler := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		return web.BadRequestErrorf("unknown parameter %q", "mbi")
	}
	w, entry := serve(context.Background(), t, "/api/claim?mbi=SYNTHETIC1234", handler)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
	if got := decodeError(t, w).Error; got != `unknown parameter "mbi"` {
		t.Errorf("error = %q, want the handler's message", got)
	}
	fields := entry.ContextMap()
	if entry.Level != zapcore.InfoLevel || fields["route"] != "/api/claim" || fields["statusCode"] != int64(400) {
		t.Errorf("log = %v %v, want info with route /api/claim and status 400", entry.Level, fields)
	}
	for key, value := range fields {
		if s, ok := value.(string); ok && (strings.Contains(s, "SYNTHETIC1234") || strings.Contains(s, "mbi")) {
			t.Errorf("log field %s = %q holds request text", key, s)
		}
	}
	if _, ok := fields["error"]; ok {
		t.Errorf("log has an error field for a client error")
	}
}

// TestServerErrorLogged: an unknown error is a generic 500 to the client
// and the full error in the log.
func TestServerErrorLogged(t *testing.T) {
	handler := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		return errors.New("claim: getting: connection refused")
	}
	w, entry := serve(context.Background(), t, "/api/claim", handler)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
	if got := decodeError(t, w).Error; got != "Internal Server Error" {
		t.Errorf("error = %q, want the generic message", got)
	}
	if entry.Level != zapcore.ErrorLevel || entry.ContextMap()["error"] != "claim: getting: connection refused" {
		t.Errorf("log = %v %v, want error level with the error", entry.Level, entry.ContextMap())
	}
}

func TestPanicRecovered(t *testing.T) {
	handler := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		panic("index out of range")
	}
	w, entry := serve(context.Background(), t, "/api/claim", handler)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
	panicText, _ := entry.ContextMap()["panicError"].(string)
	if entry.Level != zapcore.ErrorLevel || !strings.HasPrefix(panicText, "panic: index out of range") {
		t.Errorf("log = %v, panicError = %q, want error level with the panic", entry.Level, panicText)
	}
}

// TestCancelledWritesNothing: when the client has gone, nothing is written
// and the log records 499 at info level.
func TestCancelledWritesNothing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	handler := func(ctx context.Context, _ http.ResponseWriter, _ *http.Request) error {
		return ctx.Err()
	}
	w, entry := serve(ctx, t, "/api/claim", handler)
	if w.Body.Len() != 0 {
		t.Errorf("body = %q, want nothing written", w.Body.String())
	}
	if entry.Level != zapcore.InfoLevel || entry.ContextMap()["statusCode"] != int64(statusClientClosedRequest) {
		t.Errorf("log = %v %v, want info with status 499", entry.Level, entry.ContextMap())
	}
}

func TestSuccessBody(t *testing.T) {
	handler := func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		return web.RespondOK(ctx, w, map[string]string{"amount": "74.45"})
	}
	w, entry := serve(context.Background(), t, "/api/claim", handler)
	if w.Code != http.StatusOK || w.Body.String() != "{\"amount\":\"74.45\"}\n" {
		t.Errorf("response = %d %q, want 200 with JSON and a newline", w.Code, w.Body.String())
	}
	if entry.Level != zapcore.InfoLevel || entry.ContextMap()["statusCode"] != int64(200) {
		t.Errorf("log = %v %v, want info with status 200", entry.Level, entry.ContextMap())
	}
}
