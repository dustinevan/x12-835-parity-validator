// Package diff answers the file comparison endpoints:
// POST /diff compares two named input files and writes one difference
// file; GET /diffs/{filename} returns a difference file as stored. Both are
// served at the root and under /api/.
package diff

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

type Handlers struct {
	// Service is shared by pointer because it holds the lock that allows
	// one diff at a time; a copy would have its own lock.
	Service         *service.Service
	Outputs         sys.FileSystem
	OutputDirectory string
}

func NewHandlers(diffService *service.Service, outputs sys.FileSystem, outputDirectory string) Handlers {
	return Handlers{Service: diffService, Outputs: outputs, OutputDirectory: outputDirectory}
}

type diffRequest struct {
	Original  string `json:"original"`
	Candidate string `json:"candidate"`
}

// Diff compares the two named files and answers with the difference file's
// name and its number of differences. A second diff while one runs is 409,
// an input that is not a comparable 835 is 422, and a missing file is 404.
func (h Handlers) Diff(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	request, err := web.Decode[diffRequest](w, r, 4096)
	if err != nil {
		return web.BadRequestErrorf("invalid diff request: %w", err)
	}
	result, err := h.Service.Diff(ctx, request.Original, request.Candidate)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBusy):
			return web.NewRequestError(err, http.StatusConflict)
		case errors.Is(err, service.ErrInput):
			return web.NewRequestError(err, http.StatusUnprocessableEntity)
		case errors.Is(err, fs.ErrNotExist):
			return web.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, context.Canceled):
			return err
		}
		web.AddLogOnlyError(ctx, err)
		return web.InternalServerErrorf("unable to compare the files")
	}
	return web.RespondOK(ctx, w, result)
}

// File returns one difference file. The name must be a 64-character hex
// SHA-256 with .json, as Diff writes it; anything else is 404.
func (h Handlers) File(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filename := r.PathValue("filename")
	id := strings.TrimSuffix(filename, ".json")
	if _, err := hex.DecodeString(id); err != nil || len(id) != 64 || filename != id+".json" {
		return web.NotFoundErrorf("not found")
	}
	data, err := h.Outputs.Read(ctx, path.Join(h.OutputDirectory, filename))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return web.NotFoundErrorf("not found")
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		web.AddLogOnlyError(ctx, err)
		return web.InternalServerErrorf("unable to read difference file")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	web.SetStatusCode(ctx, http.StatusOK)
	if _, err := w.Write(data); err != nil {
		return web.NewResponseWriteError(err)
	}
	return nil
}
