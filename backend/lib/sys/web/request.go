package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Decode reads exactly one JSON object of type T from a request body of at
// most limit bytes. Unknown fields, a larger body, and anything after the
// object are errors.
func Decode[T any](w http.ResponseWriter, r *http.Request, limit int64) (T, error) {
	var v T
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, limit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		return v, fmt.Errorf("web: decoding request body: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return v, fmt.Errorf("web: decoding request body: expected one JSON object")
	}
	return v, nil
}

// QueryParams parses a query string. It rejects a query string that does
// not parse (url.URL.Query silently drops such pairs), a parameter not in
// allowed, and a parameter given more than once, and reports every problem
// at once. Messages name parameters, never their values.
func QueryParams(rawQuery string, allowed ...string) (url.Values, error) {
	q, err := url.ParseQuery(rawQuery)
	if err != nil {
		return q, fmt.Errorf("query string does not parse")
	}
	var errs []error
	for name, values := range q {
		known := false
		for _, a := range allowed {
			if name == a {
				known = true
			}
		}
		switch {
		case !known:
			errs = append(errs, fmt.Errorf("unknown parameter %q", name))
		case len(values) > 1:
			errs = append(errs, fmt.Errorf("parameter %q given more than once", name))
		}
	}
	return q, errors.Join(errs...)
}

// GenerateRequestID returns the current Unix time in nanoseconds: unique
// enough for one node, and it sorts by arrival.
func GenerateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// GetIPAddress returns the client address: the first X-Forwarded-For entry
// (set by the platform's proxy), then X-Real-IP, then the connection's
// remote address. A header value that is not an IP address is ignored, so
// the log field only ever holds an address.
func GetIPAddress(r *http.Request) string {
	forwarded, _, _ := strings.Cut(r.Header.Get("X-Forwarded-For"), ",")
	if addr, err := netip.ParseAddr(strings.TrimSpace(forwarded)); err == nil {
		return addr.String()
	}
	if addr, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Real-IP"))); err == nil {
		return addr.String()
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
