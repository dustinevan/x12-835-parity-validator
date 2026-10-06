package views

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	viewsctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/views"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/web"
)

// Front answers the Overview: the summary cube under a selection. Each
// filter is a comma-separated list of values; absent means every value. A
// list may hold an empty value (`adjudication=,inpatient`), the
// not-applicable bucket the cube writes as ”.
func (h Handlers) Front(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q, err := web.QueryParams(r.URL.RawQuery, "submission", "jurisdiction", "part", "workload", "claim_type", "adjudication", "mac_status", "vendor_status", "frequency", "paid_band", "service_month", "payment_band", "outcome")
	prefix, submissionErr := submissionParam(q, "submission")
	jurisdiction, jurisdictionErr := valueListParam(q, "jurisdiction", codeValue(4), "letters and digits")
	part, partErr := valueListParam(q, "part", oneOf("A", "B", "DME"), "A, B, or DME")
	workload, workloadErr := valueListParam(q, "workload", codeValue(20), "MAC workload IDs")
	claimType, claimTypeErr := valueListParam(q, "claim_type", oneOf(service.Institutional, service.Professional, service.DME), "institutional, professional, or dme")
	adjudication, adjudicationErr := valueListParam(q, "adjudication", oneOf("inpatient", "outpatient"), "inpatient or outpatient")
	macStatus, macStatusErr := valueListParam(q, "mac_status", codeValue(2), "claim status codes")
	vendorStatus, vendorStatusErr := valueListParam(q, "vendor_status", codeValue(2), "claim status codes")
	frequency, frequencyErr := valueListParam(q, "frequency", codeValue(1), "claim frequency codes")
	paidBand, paidBandErr := valueListParam(q, "paid_band", oneOf("zero", "under 100", "under 1000", "under 10000", "10000 and more"), "paid bands")
	serviceMonth, serviceMonthErr := valueListParam(q, "service_month", monthValue, "months as YYYY-MM")
	paymentBand, paymentBandErr := valueListParam(q, "payment_band", oneOf("match", "up to 1.00", "up to 1%", "up to 5%", "over 5%"), "payment bands")
	outcome, outcomeErr := valueListParam(q, "outcome", oneOf(viewsctrl.OutcomeEqual, viewsctrl.OutcomeDifferent, viewsctrl.OutcomeMissing, viewsctrl.OutcomeExtra, viewsctrl.OutcomeSetAside, viewsctrl.OutcomeDuplicated), "outcomes")
	if err = errors.Join(err, submissionErr, jurisdictionErr, partErr, workloadErr, claimTypeErr, adjudicationErr, macStatusErr, vendorStatusErr, frequencyErr, paidBandErr, serviceMonthErr, paymentBandErr, outcomeErr); err != nil {
		return badRequest(err)
	}
	value, err := h.Views.Front(ctx, viewsctrl.FrontQuery{
		Submission:   prefix,
		Jurisdiction: jurisdiction,
		Part:         part,
		Workload:     workload,
		ClaimType:    claimType,
		Adjudication: adjudication,
		MACStatus:    macStatus,
		VendorStatus: vendorStatus,
		Frequency:    frequency,
		PaidBand:     paidBand,
		ServiceMonth: serviceMonth,
		PaymentBand:  paymentBand,
		Outcome:      outcome,
	})
	if err != nil {
		return viewError(ctx, err)
	}
	return respond(ctx, w, value)
}

// valueListParam reads a comma-separated list of values; absent is nil
// (every value). An empty value is allowed in a list: it is the
// not-applicable bucket. "~" is allowed too: a value no row carries, which a
// page sends when every value of a dimension is unchecked, so the selection
// matches nothing rather than everything. Every other value must satisfy
// ok; the message names the parameter and the kind of value it takes, never
// the value.
func valueListParam(q url.Values, name string, ok func(string) bool, kind string) ([]string, error) {
	if !q.Has(name) {
		return nil, nil
	}
	values := strings.Split(q.Get(name), ",")
	for _, v := range values {
		if v != "" && v != "~" && !ok(v) {
			return nil, fmt.Errorf("%s must be comma-separated %s", name, kind)
		}
	}
	return values, nil
}

// codeValue accepts letters and digits up to n characters: the shape of a
// code the file carries.
func codeValue(n int) func(string) bool {
	return func(v string) bool {
		return len(v) <= n && !strings.ContainsFunc(v, func(r rune) bool {
			return !(r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z')
		})
	}
}

// oneOf accepts the listed values.
func oneOf(allowed ...string) func(string) bool {
	return func(v string) bool {
		for _, a := range allowed {
			if v == a {
				return true
			}
		}
		return false
	}
}

var monthPattern = regexp.MustCompile(`^[0-9]{4}-(0[1-9]|1[0-2])$`)

// monthValue accepts YYYY-MM.
func monthValue(v string) bool {
	return monthPattern.MatchString(v)
}
