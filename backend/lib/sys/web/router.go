package web

import (
	"net/http"

	"go.uber.org/zap"
)

// Route is one method and path pattern with its handler and middleware.
type Route struct {
	Method      string
	Path        string
	Handler     Handler
	Middlewares []Middleware
}

// Router collects routes; Handler mounts them on a ServeMux.
type Router struct {
	log    *zap.SugaredLogger
	routes []Route
}

func NewRouter(log *zap.SugaredLogger) Router {
	return Router{log: log}
}

// Handle adds a route. Path is a net/http ServeMux pattern path, such as
// "/api/claim" or "/diffs/{filename}".
func (r *Router) Handle(method, path string, handler Handler, middlewares ...Middleware) {
	r.routes = append(r.routes, Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
}

// GET adds a GET route. ServeMux also answers HEAD with it.
func (r *Router) GET(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(http.MethodGet, path, handler, middlewares...)
}

// POST adds a POST route.
func (r *Router) POST(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(http.MethodPost, path, handler, middlewares...)
}

// Handler returns a ServeMux with every route registered as "METHOD path".
// A path that matches no route is 404 and a method that matches no route
// for its path is 405, both from ServeMux itself.
func (r Router) Handler() http.Handler {
	mux := http.NewServeMux()
	for _, route := range r.routes {
		handler := WrapMiddleware(route.Middlewares, route.Handler)
		mux.Handle(route.Method+" "+route.Path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// With the standard stack the Logger middleware records every
			// error and returns nil; this catches a route mounted without it.
			if err := handler(req.Context(), w, req); err != nil {
				r.log.Errorw("unhandled error", "error", err, "route", req.Pattern)
			}
		}))
	}
	return mux
}
