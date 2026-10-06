// Package ui serves the browser application on every path the API does not
// take: the site built into the executable, or in development the Vite
// server. It is a plain http.Handler outside the API middleware stack;
// asset requests are not logged.
package ui

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"go.uber.org/zap"
)

// Bundled serves the built site: assets is the frontend build's dist
// directory. Client routes get the app shell; a missing asset or a
// dot-file path is 404.
func Bundled(assets fs.FS) (http.Handler, error) {
	if _, err := fs.Stat(assets, "index.html"); err != nil {
		return nil, fmt.Errorf("ui: opening bundled site: run make build or set APP_SERVER_WEB_UI_PROXY: %w", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		for _, part := range strings.Split(name, "/") {
			if strings.HasPrefix(part, ".") && part != "." {
				http.NotFound(w, r)
				return
			}
		}
		info, err := fs.Stat(assets, name)
		if err != nil || info.IsDir() {
			// Client routes use the app shell; missing assets stay 404.
			if name == "assets" || strings.HasPrefix(name, "assets/") || path.Ext(name) != "" {
				http.NotFound(w, r)
				return
			}
			name = "index.html"
		}
		w.Header().Set("Cache-Control", "no-cache")
		if strings.HasPrefix(name, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		http.ServeFileFS(w, r, assets, name)
	}), nil
}

// Proxy sends every request to the development UI server at origin, an
// HTTP or HTTPS origin with no path, query, or credentials.
func Proxy(log *zap.SugaredLogger, origin string) (http.Handler, error) {
	target, err := url.Parse(origin)
	if err != nil || target.Host == "" || (target.Scheme != "http" && target.Scheme != "https") || target.User != nil || target.RawQuery != "" || target.Fragment != "" || (target.Path != "" && target.Path != "/") {
		return nil, fmt.Errorf("ui: configuring proxy: expected an HTTP or HTTPS origin")
	}
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.SetXForwarded()
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) {
			log.Warnw("UI proxy unavailable")
			http.Error(w, "UI development server unavailable", http.StatusBadGateway)
		},
	}, nil
}
