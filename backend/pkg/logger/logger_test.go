package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestNew checks the shape of one line: JSON, the service and build fields,
// the message and its key-value pairs, and a timestamp with a colon in its
// timezone, which the RFC3339Nano encoder writes and the default does not.
func TestNew(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, "validator", "abc123")
	log.Infow("startup", "status", "database ready")
	log.Debugw("hidden")
	if err := log.Sync(); err != nil {
		t.Fatalf("Sync() = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("lines = %d, want 1 (debug is not written):\n%s", len(lines), buf.String())
	}
	var line map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &line); err != nil {
		t.Fatalf("the line is not JSON: %v\n%s", err, lines[0])
	}
	tests := []struct {
		key  string
		want string
	}{
		{"level", "info"},
		{"msg", "startup"},
		{"status", "database ready"},
		{"service", "validator"},
		{"build", "abc123"},
	}
	for _, tt := range tests {
		if got, _ := line[tt.key].(string); got != tt.want {
			t.Errorf("line[%q] = %q, want %q", tt.key, got, tt.want)
		}
	}
	ts, _ := line["ts"].(string)
	if len(ts) < 6 || ts[len(ts)-3] != ':' && ts[len(ts)-1] != 'Z' {
		t.Errorf("ts = %q, want an RFC3339 time with a colon in its zone or Z", ts)
	}
}
