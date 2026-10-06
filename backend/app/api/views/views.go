// Package views answers the read-only endpoints in docs/API.md. Each
// handler reads its query parameters (params.go), calls one
// lib/core/control/views method, and writes the result as JSON. A malformed
// or unknown parameter is 400; a submission or claim that does not exist is
// 404.
package views

import (
	"context"
	"errors"
	"net/http"
	"strings"

	viewsctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/views"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

type Handlers struct {
	Views viewsctrl.Controller
}

func NewHandlers(views viewsctrl.Controller) Handlers {
	return Handlers{Views: views}
}

// Submissions lists the MAC submission and the summarized vendor submissions.
func (h Handlers) Submissions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if _, err := web.QueryParams(r.URL.RawQuery); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Submissions(ctx)
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

// Codes returns the embedded CARC and RARC reference lists. Unlike the
// other views, this is static reference data that changes only with a new
// build, so the response is cacheable.
func (h Handlers) Codes(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if _, err := web.QueryParams(r.URL.RawQuery); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Codes(ctx)
	if err != nil {
		return viewError(ctx, err)
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	return web.RespondOK(ctx, w, value)
}

func (h Handlers) Measures(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "submission")
	prefix, submissionErr := submissionParam(q, "submission")
	if err = errors.Join(err, submissionErr); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Measures(ctx, prefix)
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

func (h Handlers) Differences(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "submission", "claim_type", "workload")
	prefix, submissionErr := submissionParam(q, "submission")
	claimType, claimTypeErr := claimTypeParam(q)
	workload, workloadErr := workloadParam(q)
	if err = errors.Join(err, submissionErr, claimTypeErr, workloadErr); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Differences(ctx, prefix, claimType, workload)
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

// Agreement counts the MAC claims by their outcome in two vendor submissions,
// named by the slots v1 and v2.
func (h Handlers) Agreement(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "v1", "v2")
	v1, v1Err := submissionParam(q, "v1")
	v2, v2Err := submissionParam(q, "v2")
	if err = errors.Join(err, v1Err, v2Err); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Agreement(ctx, v1, v2)
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

// Claims is the combination filter: every tag in all, at least one in any,
// none in none, plus the column filters.
func (h Handlers) Claims(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "submission", "all", "any", "none", "claim_type", "workload", "outcome", "other", "other_outcome", "sort", "page")
	prefix, submissionErr := submissionParam(q, "submission")
	all, allErr := tagListParam(q, "all")
	anyTags, anyErr := tagListParam(q, "any")
	none, noneErr := tagListParam(q, "none")
	claimType, claimTypeErr := claimTypeParam(q)
	workload, workloadErr := workloadParam(q)
	outcome, outcomeErr := outcomeParam(q, "outcome", false, viewsctrl.OutcomeEqual, viewsctrl.OutcomeDifferent, viewsctrl.OutcomeMissing, viewsctrl.OutcomeDuplicated, viewsctrl.OutcomeExtra, viewsctrl.OutcomeSetAside)
	other, otherOutcome, otherErr := otherParams(q)
	sort, sortErr := sortParam(q)
	page, pageErr := pageParam(q)
	if err = errors.Join(err, submissionErr, allErr, anyErr, noneErr, claimTypeErr, workloadErr, outcomeErr, otherErr, sortErr, pageErr); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Claims(ctx, viewsctrl.ClaimsQuery{
		Submission:   prefix,
		All:          all,
		Any:          anyTags,
		None:         none,
		ClaimType:    claimType,
		Workload:     workload,
		Outcome:      outcome,
		Other:        other,
		OtherOutcome: otherOutcome,
		Sort:         sort,
		Page:         page,
	})
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

func (h Handlers) Claim(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "submission", "match_key")
	prefix, submissionErr := submissionParam(q, "submission")
	matchKey, keyErr := service.ParseMatchKey(q.Get("match_key"))
	if err = errors.Join(err, submissionErr, keyErr); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Claim(ctx, prefix, matchKey)
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

func (h Handlers) Unmatched(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "submission", "outcome")
	prefix, submissionErr := submissionParam(q, "submission")
	outcome, outcomeErr := outcomeParam(q, "outcome", true, viewsctrl.OutcomeMissing, viewsctrl.OutcomeExtra, viewsctrl.OutcomeDuplicated, viewsctrl.OutcomeSetAside)
	if err = errors.Join(err, submissionErr, outcomeErr); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Unmatched(ctx, prefix, outcome)
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

// respond writes a view. Views read tables that change as submissions load,
// so browsers must not cache them.
func respond(ctx context.Context, w http.ResponseWriter, value any) error {
	w.Header().Set("Cache-Control", "no-store")
	return web.RespondOK(ctx, w, value)
}

// badRequest answers 400 with every parameter problem, one per line. The
// ErrInput suffix is internal and left out.
func badRequest(err error) error {
	return web.BadRequestErrorf("%s", strings.ReplaceAll(err.Error(), ": "+service.ErrInput.Error(), ""))
}

// viewError answers a controller error: bad input is 400, a submission or
// claim that does not exist is 404, and anything else is 500 with the cause
// in the log only.
func viewError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrInput):
		return web.BadRequestErrorf("%s", strings.TrimSuffix(err.Error(), ": "+service.ErrInput.Error()))
	case errors.Is(err, db.ErrNotFound):
		return web.NotFoundErrorf("not found")
	case errors.Is(err, context.Canceled):
		return err
	}
	web.AddLogOnlyError(ctx, err)
	return web.InternalServerErrorf("unable to read the database")
}
