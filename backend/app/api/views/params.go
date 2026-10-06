package views

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	viewsctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/views"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
)

// The parameter readers. Each returns the value or an error naming the
// parameter; no message includes the value the client sent.

// submissionParam reads a required submission prefix. Whether it exists is
// the controller's question (404); here it only has to be plausible text.
func submissionParam(q url.Values, name string) (string, error) {
	prefix := q.Get(name)
	if prefix == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	if len(prefix) > 512 || strings.ContainsFunc(prefix, func(r rune) bool { return r < 0x20 || r == 0x7f }) {
		return "", fmt.Errorf("%s is not a submission prefix", name)
	}
	return prefix, nil
}

func claimTypeParam(q url.Values) (string, error) {
	switch v := q.Get("claim_type"); v {
	case "", service.Institutional, service.Professional, service.DME:
		return v, nil
	}
	return "", fmt.Errorf("claim_type must be institutional, professional, or dme")
}

// workloadParam reads a MAC workload ID: letters and digits, as in N1*PR N104.
func workloadParam(q url.Values) (string, error) {
	v := q.Get("workload")
	if len(v) > 20 || strings.ContainsFunc(v, func(r rune) bool {
		return !(r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z')
	}) {
		return "", fmt.Errorf("workload must be a MAC workload ID")
	}
	return v, nil
}

func outcomeParam(q url.Values, name string, required bool, allowed ...string) (string, error) {
	v := q.Get(name)
	if v == "" && !required {
		return "", nil
	}
	for _, a := range allowed {
		if v == a {
			return v, nil
		}
	}
	return "", fmt.Errorf("%s must be one of %s", name, strings.Join(allowed, ", "))
}

// otherParams reads other and other_outcome, which go together: the
// claims' outcome in another vendor submission. Both absent is no filter.
func otherParams(q url.Values) (string, string, error) {
	if q.Get("other") == "" && q.Get("other_outcome") == "" {
		return "", "", nil
	}
	other, otherErr := submissionParam(q, "other")
	outcome, outcomeErr := outcomeParam(q, "other_outcome", true, viewsctrl.OutcomeEqual, viewsctrl.OutcomeDifferent, viewsctrl.OutcomeMissing, viewsctrl.OutcomeDuplicated, viewsctrl.OutcomeSetAside)
	if err := errors.Join(otherErr, outcomeErr); err != nil {
		return "", "", err
	}
	return other, outcome, nil
}

// sortParam defaults to the largest paid delta first.
func sortParam(q url.Values) (string, error) {
	switch v := q.Get("sort"); v {
	case "":
		return viewsctrl.SortDeltaDescending, nil
	case viewsctrl.SortDelta, viewsctrl.SortDeltaDescending, viewsctrl.SortPaid, viewsctrl.SortCharge:
		return v, nil
	}
	return "", fmt.Errorf("sort must be delta, -delta, paid, or charge")
}

func pageParam(q url.Values) (int, error) {
	v := q.Get("page")
	if v == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(v)
	if err != nil || page < 1 || page > 1_000_000 {
		return 0, fmt.Errorf("page must be a whole number from 1")
	}
	return page, nil
}

// tagListParam reads a comma-separated list of tags; absent is empty.
func tagListParam(q url.Values, name string) ([]string, error) {
	v := q.Get(name)
	if v == "" {
		return nil, nil
	}
	var tags []string
	for _, tag := range strings.Split(v, ",") {
		if !service.ValidTag(tag) {
			return nil, fmt.Errorf("%s must be comma-separated element labels such as CLP04,CAS*CO*253", name)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
