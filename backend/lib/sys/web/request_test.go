package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQueryParams(t *testing.T) {
	tests := []struct {
		rawQuery string
		wantErr  string
	}{
		{"submission=mac", ""},
		{"", ""},
		{"submission=mac&foo=1", `unknown parameter "foo"`},
		{"submission=a&submission=b", `parameter "submission" given more than once`},
		{"submission=%zz", "query string does not parse"},
	}
	for _, tt := range tests {
		_, err := QueryParams(tt.rawQuery, "submission")
		got := ""
		if err != nil {
			got = err.Error()
		}
		if got != tt.wantErr {
			t.Errorf("QueryParams(%q) error = %q, want %q", tt.rawQuery, got, tt.wantErr)
		}
	}
}

func TestDecode(t *testing.T) {
	type request struct {
		Original string `json:"original"`
	}
	tests := []struct {
		body    string
		wantErr bool
	}{
		{`{"original":"a.835"}`, false},
		{`{"original":"a.835","x":1}`, true},
		{`{"original":"a.835"}{}`, true},
		{`{"original":"` + strings.Repeat("a", 100) + `"}`, true},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/diff", strings.NewReader(tt.body))
		got, err := Decode[request](httptest.NewRecorder(), r, 64)
		if (err != nil) != tt.wantErr {
			t.Errorf("Decode(%q) error = %v, want error %v", tt.body, err, tt.wantErr)
		}
		if !tt.wantErr && got.Original != "a.835" {
			t.Errorf("Decode(%q) = %+v, want original a.835", tt.body, got)
		}
	}
}

func TestGetIPAddress(t *testing.T) {
	tests := []struct {
		forwarded string
		want      string
	}{
		{"203.0.113.7, 10.0.0.1", "203.0.113.7"},
		{"not an address", "192.0.2.1"},
		{"", "192.0.2.1"},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "192.0.2.1:5000"
		r.Header.Set("X-Forwarded-For", tt.forwarded)
		if got := GetIPAddress(r); got != tt.want {
			t.Errorf("GetIPAddress(X-Forwarded-For %q) = %q, want %q", tt.forwarded, got, tt.want)
		}
	}
}
