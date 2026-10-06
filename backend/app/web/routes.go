// Package web registers the service's API routes: each handler package
// under app/api is built from RouteConfig, and its methods are added to
// one router with the middleware stack. main mounts the router's handler
// and the UI.
package web

import (
	"time"

	"go.uber.org/zap"

	diffapi "github.com/dustinevan/x12-835-parity-validator/backend/app/api/diff"
	systemapi "github.com/dustinevan/x12-835-parity-validator/backend/app/api/system"
	viewsapi "github.com/dustinevan/x12-835-parity-validator/backend/app/api/views"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/web/middleware"
	viewsctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/views"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// RouteConfig holds everything the handlers need.
type RouteConfig struct {
	Log            *zap.SugaredLogger
	RequestTimeout time.Duration

	// Diff is shared by pointer: it holds the one-diff-at-a-time lock.
	Diff            *service.Service
	Outputs         sys.FileSystem
	OutputDirectory string

	Views viewsctrl.Controller
}

// RegisterRoutes returns the router with every API route. The views are
// under /api/ only; /health and the diff endpoints are also served at the
// root.
func RegisterRoutes(cfg RouteConfig) web.Router {
	router := web.NewRouter(cfg.Log)
	stack := middleware.Stack(cfg.Log, cfg.RequestTimeout)

	router.GET("/health", systemapi.Health, stack...)
	router.GET("/api/health", systemapi.Health, stack...)

	diff := diffapi.NewHandlers(cfg.Diff, cfg.Outputs, cfg.OutputDirectory)
	router.POST("/diff", diff.Diff, stack...)
	router.POST("/api/diff", diff.Diff, stack...)
	router.GET("/diffs/{filename}", diff.File, stack...)
	router.GET("/api/diffs/{filename}", diff.File, stack...)

	views := viewsapi.NewHandlers(cfg.Views)
	router.GET("/api/codes", views.Codes, stack...)
	router.GET("/api/macs", views.MACs, stack...)
	router.GET("/api/submissions", views.Submissions, stack...)
	router.GET("/api/measures", views.Measures, stack...)
	router.GET("/api/differences", views.Differences, stack...)
	router.GET("/api/agreement", views.Agreement, stack...)
	router.GET("/api/claims", views.Claims, stack...)
	router.GET("/api/claim", views.Claim, stack...)
	router.GET("/api/unmatched", views.Unmatched, stack...)
	router.GET("/api/front", views.Front, stack...)

	return router
}
