package views

import (
	"context"
	"fmt"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/cell"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/maccell"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/tagcell"
)

// The Overview: counts and sums over the summary cube (cell, tag_cell)
// under a selection, so a reader's clicks never group the comparison rows.
// FrontQuery is the selection: for each dimension, the values kept, an
// empty list keeping every value; ServiceMonth holds YYYY-MM texts.
type FrontQuery struct {
	Submission   string
	Jurisdiction []string
	Part         []string
	Workload     []string
	ClaimType    []string
	Adjudication []string
	MACStatus    []string
	VendorStatus []string
	Frequency    []string
	PaidBand     []string
	ServiceMonth []string
	PaymentBand  []string
	Outcome      []string
}

// SumsView is the cube's counts and sums over a set of cells. Amounts are
// exact decimal strings. Claims counts every comparison row; MACClaims the
// rows that have a MAC claim, every outcome but extra.
type SumsView struct {
	Claims          int64  `json:"claims"`
	MACClaims       int64  `json:"mac_claims"`
	MACPaid         string `json:"mac_paid"`
	VendorPaid      string `json:"vendor_paid"`
	Billed          string `json:"billed"`
	NetDifference   string `json:"net_difference"`
	GrossDifference string `json:"gross_difference"`
	SameStatus      int64  `json:"same_status"`
}

type OutcomeSumsView struct {
	Outcome string `json:"outcome"`
	SumsView
}

type ValueSumsView struct {
	Value string `json:"value"`
	SumsView
}

type DimensionView struct {
	Name   string          `json:"name"`
	Values []ValueSumsView `json:"values"`
}

type StatusCellView struct {
	MACStatus    string `json:"mac_status"`
	VendorStatus string `json:"vendor_status"`
	Claims       int64  `json:"claims"`
	MACPaid      string `json:"mac_paid"`
	VendorPaid   string `json:"vendor_paid"`
}

type TagSumsView struct {
	service.Tag
	Claims  int64  `json:"claims"`
	Dollars string `json:"dollars"`
}

// LegacySumsView is the legacy cube's figures over a set of cells: the MAC
// claims (each claim once, a key the MAC repeats counting each of its
// claims), what the MAC paid on them (CLP04), and what was billed (CLP03).
type LegacySumsView struct {
	Claims int64  `json:"claims"`
	Paid   string `json:"paid"`
	Billed string `json:"billed"`
}

type LegacyValueView struct {
	Value string `json:"value"`
	LegacySumsView
}

type LegacyDimensionView struct {
	Name   string            `json:"name"`
	Values []LegacyValueView `json:"values"`
}

// LegacyView is the legacy side's own figures, from the legacy cube
// (mac_cell): the totals and the nine dimensions the MAC claim carries,
// under the selection on those dimensions only. The vendor status,
// outcome, and payment band filters do not apply to it.
type LegacyView struct {
	Totals     LegacySumsView        `json:"totals"`
	Dimensions []LegacyDimensionView `json:"dimensions"`
}

type FrontResponse struct {
	Submission   SubmissionView    `json:"submission"`
	Totals       SumsView          `json:"totals"`
	Outcomes     []OutcomeSumsView `json:"outcomes"`
	Dimensions   []DimensionView   `json:"dimensions"`
	StatusMatrix []StatusCellView  `json:"status_matrix"`
	Tags         []TagSumsView     `json:"tags"`
	Legacy       LegacyView        `json:"legacy"`
}

func legacySumsView(s maccell.Sums) LegacySumsView {
	return LegacySumsView{Claims: s.Claims, Paid: s.Paid, Billed: s.Billed}
}

func sumsView(s cell.Sums) SumsView {
	return SumsView{
		Claims:          s.Claims,
		MACClaims:       s.MACClaims,
		MACPaid:         s.MACPaid,
		VendorPaid:      s.VendorPaid,
		Billed:          s.Billed,
		NetDifference:   s.NetDifference,
		GrossDifference: s.GrossDifference,
		SameStatus:      s.SameStatus,
	}
}

// serviceMonths turns YYYY-MM texts into first-of-month dates. The handler
// has checked the form.
func serviceMonths(months []string) ([]time.Time, error) {
	var out []time.Time
	for _, m := range months {
		t, err := time.Parse("2006-01", m)
		if err != nil {
			return nil, fmt.Errorf("service month %q: %w", m, service.ErrInput)
		}
		out = append(out, t.UTC())
	}
	return out, nil
}

// Front answers the Overview under a selection. totals, outcomes, the
// status matrix (pairs only), and the tags are added under the whole
// selection. Each dimension's values are added under the selection on
// every other dimension and not its own, so a reader sees what each value
// would contribute; every dimension is listed, in the API's order, with
// the values present under that selection in value order. legacy is the
// MAC submission's own figures from mac_cell, written by load, under the
// selection on the nine dimensions the MAC claim carries, its dimensions
// added by the same rule; it is the same for every vendor submission.
// Nothing here is a rate: counts and sums only. It reads cell and
// tag_cell, written by summarize, so the submission must be summarized.
func (v Controller) Front(ctx context.Context, q FrontQuery) (FrontResponse, error) {
	d, _, macPrefix, err := v.vendorSubmission(ctx, q.Submission)
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}
	months, err := serviceMonths(q.ServiceMonth)
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}
	f := cell.Filter{
		SubmissionPrefix: d.Prefix,
		Jurisdiction:     q.Jurisdiction,
		Part:             q.Part,
		Workload:         q.Workload,
		ClaimType:        q.ClaimType,
		Adjudication:     q.Adjudication,
		MACStatus:        q.MACStatus,
		VendorStatus:     q.VendorStatus,
		Frequency:        q.Frequency,
		PaidBand:         q.PaidBand,
		ServiceMonth:     months,
		PaymentBand:      q.PaymentBand,
		Outcome:          q.Outcome,
	}
	cells := cell.NewTableAccess(v.db)

	totals, err := cells.Totals(ctx, f)
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}
	outcomes, err := cells.Outcomes(ctx, f)
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}
	matrix, err := cells.StatusMatrix(ctx, f)
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}
	tags, err := tagcell.NewTableAccess(v.db).Totals(ctx, tagcell.Filter{
		SubmissionPrefix: f.SubmissionPrefix,
		Jurisdiction:     f.Jurisdiction,
		Part:             f.Part,
		Workload:         f.Workload,
		ClaimType:        f.ClaimType,
		Adjudication:     f.Adjudication,
		MACStatus:        f.MACStatus,
		VendorStatus:     f.VendorStatus,
		Frequency:        f.Frequency,
		PaidBand:         f.PaidBand,
		ServiceMonth:     f.ServiceMonth,
		PaymentBand:      f.PaymentBand,
		Outcome:          f.Outcome,
	})
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}

	out := FrontResponse{
		Submission:   submissionView(d),
		Totals:       sumsView(totals),
		Outcomes:     make([]OutcomeSumsView, 0, len(outcomes)),
		Dimensions:   make([]DimensionView, 0, len(cell.Dimensions)),
		StatusMatrix: make([]StatusCellView, 0, len(matrix)),
		Tags:         make([]TagSumsView, 0, len(tags)),
	}
	for _, o := range outcomes {
		out.Outcomes = append(out.Outcomes, OutcomeSumsView{Outcome: o.Outcome, SumsView: sumsView(o.Sums)})
	}
	for _, dimension := range cell.Dimensions {
		values, err := cells.Facet(ctx, f, dimension)
		if err != nil {
			return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
		}
		dv := DimensionView{Name: dimension, Values: make([]ValueSumsView, 0, len(values))}
		for _, value := range values {
			dv.Values = append(dv.Values, ValueSumsView{Value: value.Value, SumsView: sumsView(value.Sums)})
		}
		out.Dimensions = append(out.Dimensions, dv)
	}
	for _, c := range matrix {
		out.StatusMatrix = append(out.StatusMatrix, StatusCellView{
			MACStatus:    c.MACStatus,
			VendorStatus: c.VendorStatus,
			Claims:       c.Claims,
			MACPaid:      c.MACPaid,
			VendorPaid:   c.VendorPaid,
		})
	}
	for _, t := range tags {
		out.Tags = append(out.Tags, TagSumsView{Tag: service.TagOf(t.Tag), Claims: t.Claims, Dollars: t.Dollars})
	}

	lf := maccell.Filter{
		SubmissionPrefix: macPrefix,
		Jurisdiction:     f.Jurisdiction,
		Part:             f.Part,
		Workload:         f.Workload,
		ClaimType:        f.ClaimType,
		Adjudication:     f.Adjudication,
		MACStatus:        f.MACStatus,
		Frequency:        f.Frequency,
		PaidBand:         f.PaidBand,
		ServiceMonth:     f.ServiceMonth,
	}
	legacy := maccell.NewTableAccess(v.db)
	legacyTotals, err := legacy.Totals(ctx, lf)
	if err != nil {
		return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
	}
	out.Legacy = LegacyView{Totals: legacySumsView(legacyTotals), Dimensions: make([]LegacyDimensionView, 0, len(maccell.Dimensions))}
	for _, dimension := range maccell.Dimensions {
		values, err := legacy.Facet(ctx, lf, dimension)
		if err != nil {
			return FrontResponse{}, fmt.Errorf("views: reading front: %w", err)
		}
		dv := LegacyDimensionView{Name: dimension, Values: make([]LegacyValueView, 0, len(values))}
		for _, value := range values {
			dv.Values = append(dv.Values, LegacyValueView{Value: value.Value, LegacySumsView: legacySumsView(value.Sums)})
		}
		out.Legacy.Dimensions = append(out.Legacy.Dimensions, dv)
	}
	return out, nil
}
