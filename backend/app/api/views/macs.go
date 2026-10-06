package views

import (
	"context"
	"net/http"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// MACs returns the MAC workload numbers as a tree (backend/pkg/macs). Like
// Codes, this is static reference data that changes only with a new build,
// so the response is cacheable.
func (h Handlers) MACs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if _, err := web.QueryParams(r.URL.RawQuery); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.MACs(ctx)
	if err != nil {
		return viewError(ctx, err)
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	return web.RespondOK(ctx, w, value)
}
