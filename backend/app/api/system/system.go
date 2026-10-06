// Package system answers the service's own endpoints. GET /health is 204
// once the service is listening; startup has already checked the database
// and its migrations.
package system

import (
	"context"
	"net/http"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// Health answers 204 with no body.
func Health(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	return web.RespondNoContent(ctx, w)
}
