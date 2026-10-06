// Package views answers the read-only endpoints in docs/API.md from the
// database, reading the table packages and lib/core/query.
//
// Whole-submission totals come from the summarize stage's tables: measure
// (the measure rows, the unmatched share, the differences' claim count),
// tag_total (differences rows), and cell and tag_cell (the Overview, under
// a selection; front.go). Per-claim lookups and paged lists read
// comparison, difference, and claim. The views show the MAC submission and
// summarized vendor submissions only (shown).
package views

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/query"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/comparison"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/difference"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/measure"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/submission"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/tagtotal"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"
)

// Controller answers the read-only endpoints in docs/API.md.
// Amounts are passed through as the decimal text Postgres stores; nothing
// here does arithmetic on money.
type Controller struct {
	db db.DB
}

func NewController(database db.DB) Controller {
	return Controller{db: database}
}

// Outcomes, as stored in comparison.outcome.
const (
	OutcomeEqual      = comparison.OutcomeEqual
	OutcomeDifferent  = comparison.OutcomeDifferent
	OutcomeMissing    = comparison.OutcomeMissing
	OutcomeDuplicated = comparison.OutcomeDuplicated
	OutcomeExtra      = comparison.OutcomeExtra
	OutcomeSetAside   = comparison.OutcomeSetAside
)

// Sorts for the claim lists.
const (
	SortDelta           = query.SortDelta
	SortDeltaDescending = query.SortDeltaDescending
	SortPaid            = query.SortPaid
	SortCharge          = query.SortCharge
)

const (
	pageSize          = 25
	bucketsPerSide    = 5
	unmatchedRowLimit = 1000
)

// nullable is a string that is JSON null when empty. The data layer uses
// the empty string for NULL, so this carries that through to the API.
type nullable string

func (n nullable) MarshalJSON() ([]byte, error) {
	if n == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(n))
}

// ---------------------------------------------------------------- shapes

type KeyView struct {
	MBI           string `json:"mbi"`
	DateOfService string `json:"date_of_service"` // YYYY-MM-DD
	TotalCharge   string `json:"total_charge"`
	TaxID         string `json:"tax_id"`
	NPI           string `json:"npi"`
	Workload      string `json:"workload"`
}

type SubmissionView struct {
	Prefix      string   `json:"prefix"`
	Side        string   `json:"side"`
	Label       string   `json:"label"`
	SubmittedOn nullable `json:"submitted_on"`
	Status      string   `json:"status"`
	ClaimCount  int64    `json:"claim_count"`
}

func submissionView(d submission.Submission) SubmissionView {
	var on nullable
	if !d.SubmittedOn.IsZero() {
		on = nullable(d.SubmittedOn.Format(time.DateOnly))
	}
	return SubmissionView{
		Prefix:      d.Prefix,
		Side:        d.Side,
		Label:       d.Label,
		SubmittedOn: on,
		Status:      d.Status,
		ClaimCount:  d.ClaimCount,
	}
}

type MeasureView struct {
	ClaimType        string   `json:"claim_type"`
	Workload         string   `json:"workload"`
	MACClaims        int64    `json:"mac_claims"`
	Equal            int64    `json:"equal"`
	Different        int64    `json:"different"`
	Missing          int64    `json:"missing"`
	Extra            int64    `json:"extra"`
	SetAside         int64    `json:"set_aside"`
	Duplicated       int64    `json:"duplicated"`
	OutcomeMatchRate nullable `json:"outcome_match_rate"`
	PaidClaims       int64    `json:"paid_claims"`
	MedianDeviation  nullable `json:"median_deviation"`
	P95Deviation     nullable `json:"p95_deviation"`
}

func measureViews(ms []measure.Measure) []MeasureView {
	out := make([]MeasureView, 0, len(ms))
	for _, m := range ms {
		out = append(out, MeasureView{
			ClaimType:        m.ClaimType,
			Workload:         m.MACWorkloadID,
			MACClaims:        m.MACClaims,
			Equal:            m.EqualCount,
			Different:        m.DifferentCount,
			Missing:          m.MissingCount,
			Extra:            m.ExtraCount,
			SetAside:         m.SetAsideCount,
			Duplicated:       m.DuplicatedCount,
			OutcomeMatchRate: nullable(m.OutcomeMatchRate),
			PaidClaims:       m.PaidClaims,
			MedianDeviation:  nullable(m.MedianDeviation),
			P95Deviation:     nullable(m.P95Deviation),
		})
	}
	return out
}

type DifferenceView struct {
	Element           string   `json:"element"`
	Occurrence        int      `json:"occurrence"`
	LineControlNumber nullable `json:"line_control_number"`
	MACValue          nullable `json:"mac_value"`
	VendorValue       nullable `json:"vendor_value"`
	AmountDifference  nullable `json:"amount_difference"`
}

func differenceView(d difference.Difference) DifferenceView {
	return DifferenceView{
		Element:           d.Element,
		Occurrence:        d.Occurrence,
		LineControlNumber: nullable(d.LineControlNumber),
		MACValue:          nullable(d.MACValue),
		VendorValue:       nullable(d.VendorValue),
		AmountDifference:  nullable(d.AmountDifference),
	}
}

// ClaimView is a claim as stored: its bytes as text, CLP through its last
// service line, and the four separators of its file, which a reader needs
// to read the text back. The claim screen reads it in the browser.
type ClaimView struct {
	Segments   string         `json:"segments"`
	Delimiters DelimitersView `json:"delimiters"`
}

// DelimitersView is a file's four separators, one character each.
type DelimitersView struct {
	Element    string `json:"element"`
	Component  string `json:"component"`
	Repetition string `json:"repetition"`
	Terminator string `json:"terminator"`
}

// claimView reads a stored claim's delimiters and passes its bytes through
// as they were received.
func claimView(c claim.Claim) (ClaimView, error) {
	d, err := service.ParseDelimiters(c.Delimiters)
	if err != nil {
		return ClaimView{}, err
	}
	return ClaimView{
		Segments: string(c.Segments),
		Delimiters: DelimitersView{
			Element:    string(d.Element),
			Component:  string(d.Component),
			Repetition: string(d.Repetition),
			Terminator: string(d.Segment),
		},
	}, nil
}

// ClaimRow is one row of the claim list (/api/claims). MAC paid is null
// for an extra claim, vendor paid for a missing or set-aside one, and delta
// when either is.
type ClaimRow struct {
	MatchKey   string   `json:"match_key"`
	Key        KeyView  `json:"key"`
	ClaimType  string   `json:"claim_type"`
	MACPaid    nullable `json:"mac_paid"`
	VendorPaid nullable `json:"vendor_paid"`
	Delta      nullable `json:"delta"`
}

func keyOfRow(r query.Row) KeyView {
	return KeyView{
		MBI:           r.MedicareBeneficiaryIdentifier,
		DateOfService: r.DateOfService.Format(time.DateOnly),
		TotalCharge:   r.TotalClaimChargeAmount,
		TaxID:         r.FederalTaxpayerIdentificationNumber,
		NPI:           r.PayeeIdentificationCode,
		Workload:      r.MACWorkloadID,
	}
}

func keyOfClaim(c claim.Claim) KeyView {
	return KeyView{
		MBI:           c.MedicareBeneficiaryIdentifier,
		DateOfService: c.DateOfService.Format(time.DateOnly),
		TotalCharge:   c.TotalClaimChargeAmount,
		TaxID:         c.FederalTaxpayerIdentificationNumber,
		NPI:           c.PayeeIdentificationCode,
		Workload:      c.MACWorkloadID,
	}
}

// ---------------------------------------------------------------- submissions

type SubmissionsResponse struct {
	Submissions []SubmissionView `json:"submissions"`
}

// shown reports whether the views show a submission: the MAC submission once
// loaded, a vendor submission once summarized. This status gate keeps a
// submission that is being compared or summarized off the screens.
func shown(d submission.Submission) bool {
	return d.Side == submission.SideMAC || d.Status == submission.StatusSummarized
}

// Submissions lists every submission the views show: the MAC first, then
// summarized vendor submissions by label and date.
func (v Controller) Submissions(ctx context.Context) (SubmissionsResponse, error) {
	ds, err := submission.NewTableAccess(v.db).List(ctx)
	if err != nil {
		return SubmissionsResponse{}, fmt.Errorf("views: listing submissions: %w", err)
	}
	out := SubmissionsResponse{Submissions: make([]SubmissionView, 0, len(ds))}
	for _, d := range ds {
		if shown(d) {
			out.Submissions = append(out.Submissions, submissionView(d))
		}
	}
	return out, nil
}

// vendorSubmission returns a summarized vendor submission, the summarized one
// before it from the same vendor (zero when there is none), and the MAC
// submission's prefix.
func (v Controller) vendorSubmission(ctx context.Context, prefix string) (submission.Submission, submission.Submission, string, error) {
	ds, err := submission.NewTableAccess(v.db).List(ctx)
	if err != nil {
		return submission.Submission{}, submission.Submission{}, "", err
	}
	var found, previous submission.Submission
	macPrefix := ""
	for _, d := range ds {
		if d.Side == submission.SideMAC {
			macPrefix = d.Prefix
		}
	}
	// List orders vendors by label, then date, so the previous submission of
	// the same vendor is the nearest earlier shown row with the same label.
	for _, d := range ds {
		if d.Prefix == prefix {
			found = d
			break
		}
		if d.Side == submission.SideVendor && shown(d) {
			previous = d
		}
	}
	if found.Prefix == "" {
		return submission.Submission{}, submission.Submission{}, "", fmt.Errorf("submission %q: %w", prefix, db.ErrNotFound)
	}
	if found.Side != submission.SideVendor {
		return submission.Submission{}, submission.Submission{}, "", fmt.Errorf("submission %q is not a vendor submission: %w", prefix, service.ErrInput)
	}
	if found.Status != submission.StatusSummarized {
		return submission.Submission{}, submission.Submission{}, "", fmt.Errorf("submission %q has not been summarized: %w", prefix, service.ErrInput)
	}
	if previous.Label != found.Label {
		previous = submission.Submission{}
	}
	if macPrefix == "" {
		return submission.Submission{}, submission.Submission{}, "", fmt.Errorf("MAC submission: %w", db.ErrNotFound)
	}
	return found, previous, macPrefix, nil
}

// ---------------------------------------------------------------- measures

type MeasuresResponse struct {
	Submission       SubmissionView `json:"submission"`
	Previous         any            `json:"previous"` // SubmissionView or null
	Measures         []MeasureView  `json:"measures"`
	PreviousMeasures []MeasureView  `json:"previous_measures"`
}

// Measures returns every measure row of a vendor submission and of the same
// vendor's previous submission.
func (v Controller) Measures(ctx context.Context, prefix string) (MeasuresResponse, error) {
	d, previous, _, err := v.vendorSubmission(ctx, prefix)
	if err != nil {
		return MeasuresResponse{}, fmt.Errorf("views: reading measures: %w", err)
	}
	measures := measure.NewTableAccess(v.db)
	ms, err := measures.ListBySubmission(ctx, d.Prefix)
	if err != nil {
		return MeasuresResponse{}, fmt.Errorf("views: reading measures: %w", err)
	}
	out := MeasuresResponse{
		Submission:       submissionView(d),
		Measures:         measureViews(ms),
		PreviousMeasures: []MeasureView{},
	}
	if previous.Prefix != "" {
		pms, err := measures.ListBySubmission(ctx, previous.Prefix)
		if err != nil {
			return MeasuresResponse{}, fmt.Errorf("views: reading measures: %w", err)
		}
		out.Previous = submissionView(previous)
		out.PreviousMeasures = measureViews(pms)
	}
	return out, nil
}

// ---------------------------------------------------------------- differences

type DifferenceCount struct {
	service.Tag
	Claims         int64  `json:"claims"`
	Dollars        string `json:"dollars"`
	PreviousClaims any    `json:"previous_claims"` // int64 or null
}

type DifferencesResponse struct {
	DifferentClaims int64             `json:"different_claims"`
	Rows            []DifferenceCount `json:"rows"`
}

// Differences ranks the tags of a submission's different claims by how many
// claims carry each. Dollars is the sum of |amount_difference| of that
// element. PreviousClaims is the same count in the vendor's previous
// submission, null when there is none. It reads the summarize stage's
// measure (the different claims) and tag_total rows; an empty claimType
// or workload reads the 'all' rows.
func (v Controller) Differences(ctx context.Context, prefix, claimType, workload string) (DifferencesResponse, error) {
	d, previous, _, err := v.vendorSubmission(ctx, prefix)
	if err != nil {
		return DifferencesResponse{}, fmt.Errorf("views: counting differences: %w", err)
	}
	// 'all' marks the total rows; it is never a workload (claimType is
	// checked by the caller), so as a filter it selects no claims.
	if claimType == measure.All || workload == measure.All {
		return DifferencesResponse{Rows: []DifferenceCount{}}, nil
	}
	if claimType == "" {
		claimType = measure.All
	}
	if workload == "" {
		workload = measure.All
	}
	var differentClaims int64
	m, err := measure.NewTableAccess(v.db).Get(ctx, d.Prefix, claimType, workload)
	switch {
	case err == nil:
		differentClaims = m.DifferentCount
	case !errors.Is(err, db.ErrNotFound):
		return DifferencesResponse{}, fmt.Errorf("views: counting differences: %w", err)
	}
	tagTotals := tagtotal.NewTableAccess(v.db)
	counts, err := tagTotals.List(ctx, d.Prefix, claimType, workload)
	if err != nil {
		return DifferencesResponse{}, fmt.Errorf("views: counting differences: %w", err)
	}
	var before map[string]int64
	if previous.Prefix != "" {
		previousCounts, err := tagTotals.List(ctx, previous.Prefix, claimType, workload)
		if err != nil {
			return DifferencesResponse{}, fmt.Errorf("views: counting differences: %w", err)
		}
		before = map[string]int64{}
		for _, c := range previousCounts {
			before[c.Tag] = c.Claims
		}
	}
	out := DifferencesResponse{DifferentClaims: differentClaims, Rows: make([]DifferenceCount, 0, len(counts))}
	for _, c := range counts {
		row := DifferenceCount{Tag: service.TagOf(c.Tag), Claims: c.Claims, Dollars: c.Dollars}
		if before != nil {
			row.PreviousClaims = before[c.Tag]
		}
		out.Rows = append(out.Rows, row)
	}
	return out, nil
}

// ---------------------------------------------------------------- agreement

// AgreementCell and AgreementResponse name the two submissions v1 and v2,
// as the request does: in a cell, v1 is the claims' outcome in the
// submission requested as v1, and v2 their outcome in the one requested as
// v2.
type AgreementCell struct {
	V1       string `json:"v1"`
	V2       string `json:"v2"`
	Claims   int64  `json:"claims"`
	SameTags int64  `json:"same_tags"`
}

type AgreementResponse struct {
	V1    SubmissionView  `json:"v1"`
	V2    SubmissionView  `json:"v2"`
	Cells []AgreementCell `json:"cells"`
}

// Agreement counts the MAC claims by their outcome in two vendor
// submissions: equal, different, missing, or duplicated in each. Set-aside
// and extra keys are left out. SameTags counts, among the claims different
// in both, those whose differing elements are the same in both.
func (v Controller) Agreement(ctx context.Context, v1Prefix, v2Prefix string) (AgreementResponse, error) {
	v1, _, _, err := v.vendorSubmission(ctx, v1Prefix)
	if err != nil {
		return AgreementResponse{}, fmt.Errorf("views: counting agreement: v1: %w", err)
	}
	v2, _, _, err := v.vendorSubmission(ctx, v2Prefix)
	if err != nil {
		return AgreementResponse{}, fmt.Errorf("views: counting agreement: v2: %w", err)
	}
	cells, err := comparison.NewTableAccess(v.db).Agreement(ctx, v1.Prefix, v2.Prefix)
	if err != nil {
		return AgreementResponse{}, fmt.Errorf("views: counting agreement: %w", err)
	}
	out := AgreementResponse{V1: submissionView(v1), V2: submissionView(v2), Cells: make([]AgreementCell, 0, len(cells))}
	for _, c := range cells {
		out.Cells = append(out.Cells, AgreementCell{V1: c.V1Outcome, V2: c.V2Outcome, Claims: c.Claims, SameTags: c.SameTags})
	}
	return out, nil
}

// ---------------------------------------------------------------- claim lists

// ClaimsQuery is the combination filter. All, Any, and None are tags; an
// empty ClaimType, Workload, or Outcome means all. A non-empty
// OtherOutcome keeps the claims with that outcome in the vendor submission
// Other. Page counts from 1.
type ClaimsQuery struct {
	Submission   string
	All          []string
	Any          []string
	None         []string
	ClaimType    string
	Workload     string
	Outcome      string
	Other        string
	OtherOutcome string
	Sort         string
	Page         int
}

type Bucket struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Count int64  `json:"count"`
}

type ClaimsResponse struct {
	Claims    int64      `json:"claims"`
	Dollars   string     `json:"dollars"`
	Histogram []Bucket   `json:"histogram"`
	Rows      []ClaimRow `json:"rows"`
	Page      int        `json:"page"`
	Pages     int        `json:"pages"`
}

// Claims lists the claims that match the combination filter, 25 to a
// page. Dollars is the sum of |vendor paid - MAC paid| over those claims.
func (v Controller) Claims(ctx context.Context, q ClaimsQuery) (ClaimsResponse, error) {
	d, _, macPrefix, err := v.vendorSubmission(ctx, q.Submission)
	if err != nil {
		return ClaimsResponse{}, fmt.Errorf("views: listing claims: %w", err)
	}
	f := comparison.Filter{
		SubmissionPrefix: d.Prefix,
		ClaimType:        q.ClaimType,
		Workload:         q.Workload,
		Outcome:          q.Outcome,
		All:              q.All,
		Any:              q.Any,
		None:             q.None,
	}
	if q.OtherOutcome != "" {
		other, _, _, err := v.vendorSubmission(ctx, q.Other)
		if err != nil {
			return ClaimsResponse{}, fmt.Errorf("views: listing claims: other: %w", err)
		}
		f.OtherSubmissionPrefix = other.Prefix
		f.OtherOutcome = q.OtherOutcome
	}
	comparisons := comparison.NewTableAccess(v.db)
	totals, err := comparisons.Totals(ctx, f)
	if err != nil {
		return ClaimsResponse{}, fmt.Errorf("views: listing claims: %w", err)
	}
	buckets, err := comparisons.Histogram(ctx, f, bucketsPerSide)
	if err != nil {
		return ClaimsResponse{}, fmt.Errorf("views: listing claims: %w", err)
	}
	rows, err := query.NewClaimListQuery(v.db).List(ctx, f, macPrefix, q.Sort, pageSize, (q.Page-1)*pageSize)
	if err != nil {
		return ClaimsResponse{}, fmt.Errorf("views: listing claims: %w", err)
	}

	pages := int((totals.Claims + pageSize - 1) / pageSize)
	if pages == 0 {
		pages = 1
	}
	out := ClaimsResponse{
		Claims:    totals.Claims,
		Dollars:   totals.PaidDelta,
		Histogram: make([]Bucket, 0, len(buckets)),
		Rows:      make([]ClaimRow, 0, len(rows)),
		Page:      q.Page,
		Pages:     pages,
	}
	for _, b := range buckets {
		out.Histogram = append(out.Histogram, Bucket{From: b.From, To: b.To, Count: b.Count})
	}
	for _, r := range rows {
		out.Rows = append(out.Rows, ClaimRow{
			MatchKey:   hex.EncodeToString(r.MatchKey),
			Key:        keyOfRow(r),
			ClaimType:  r.ClaimType,
			MACPaid:    nullable(r.MACClaimPaymentAmount),
			VendorPaid: nullable(r.VendorClaimPaymentAmount),
			Delta:      nullable(r.Delta),
		})
	}
	return out, nil
}

// ---------------------------------------------------------------- one claim

type ClaimResponse struct {
	Key          KeyView          `json:"key"`
	Outcome      string           `json:"outcome"`
	MAC          any              `json:"mac"`           // ClaimView or null
	Vendor       any              `json:"vendor"`        // ClaimView or null: the first of VendorClaims
	VendorClaims []ClaimView      `json:"vendor_claims"` // every vendor claim with the key, in file order
	Differences  []DifferenceView `json:"differences"`
}

// Claim returns one key's result in a vendor submission, each claim as
// stored (ClaimView), with the differences saved for it. The MAC side
// shows its first claim in file order (several share the key when it is
// set aside); the vendor side shows its first claim and lists every one
// (several when the key is duplicated).
func (v Controller) Claim(ctx context.Context, prefix string, matchKey []byte) (ClaimResponse, error) {
	d, _, macPrefix, err := v.vendorSubmission(ctx, prefix)
	if err != nil {
		return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
	}
	c, err := comparison.NewTableAccess(v.db).Get(ctx, d.Prefix, matchKey)
	if err != nil {
		return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
	}
	claims := claim.NewTableAccess(v.db)
	macRows, err := claims.ListByMatchKey(ctx, macPrefix, matchKey)
	if err != nil {
		return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
	}
	vendorRows, err := claims.ListByMatchKey(ctx, d.Prefix, matchKey)
	if err != nil {
		return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
	}
	diffs, err := difference.NewTableAccess(v.db).ListByMatchKey(ctx, d.Prefix, matchKey)
	if err != nil {
		return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
	}

	out := ClaimResponse{
		Outcome:      c.Outcome,
		VendorClaims: make([]ClaimView, 0, len(vendorRows)),
		Differences:  make([]DifferenceView, 0, len(diffs)),
	}
	if len(macRows) > 0 {
		out.Key = keyOfClaim(macRows[0])
		mac, err := claimView(macRows[0])
		if err != nil {
			return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
		}
		out.MAC = mac
	}
	if len(vendorRows) > 0 && len(macRows) == 0 {
		out.Key = keyOfClaim(vendorRows[0])
	}
	for _, row := range vendorRows {
		vendor, err := claimView(row)
		if err != nil {
			return ClaimResponse{}, fmt.Errorf("views: reading claim: %w", err)
		}
		out.VendorClaims = append(out.VendorClaims, vendor)
	}
	if len(out.VendorClaims) > 0 {
		out.Vendor = out.VendorClaims[0]
	}
	for _, diff := range diffs {
		out.Differences = append(out.Differences, differenceView(diff))
	}
	return out, nil
}

// ---------------------------------------------------------------- unmatched

type CountByClaimType struct {
	ClaimType string `json:"claim_type"`
	Count     int64  `json:"count"`
}

type CountByWorkload struct {
	Workload string `json:"workload"`
	Count    int64  `json:"count"`
}

type UnmatchedRow struct {
	MatchKey             string  `json:"match_key"`
	Key                  KeyView `json:"key"`
	ClaimType            string  `json:"claim_type"`
	PatientControlNumber string  `json:"patient_control_number"`
	Paid                 string  `json:"paid"`
}

type UnmatchedResponse struct {
	Outcome     string             `json:"outcome"`
	Count       int64              `json:"count"`
	Share       string             `json:"share"`
	ByClaimType []CountByClaimType `json:"by_claim_type"`
	ByWorkload  []CountByWorkload  `json:"by_workload"`
	Rows        []UnmatchedRow     `json:"rows"`
}

// Unmatched lists the claims behind a submission's missing, extra,
// duplicated, or set-aside keys: the MAC's claims for missing, duplicated,
// and set aside, the vendor's for extra. Count and the breakdowns count
// those claims; share is count divided by the submission's MAC claims. Rows
// stop at 1,000.
func (v Controller) Unmatched(ctx context.Context, prefix, outcome string) (UnmatchedResponse, error) {
	d, _, macPrefix, err := v.vendorSubmission(ctx, prefix)
	if err != nil {
		return UnmatchedResponse{}, fmt.Errorf("views: listing unmatched claims: %w", err)
	}
	claimPrefix := macPrefix
	if outcome == OutcomeExtra {
		claimPrefix = d.Prefix
	}
	unmatched := query.NewUnmatchedQuery(v.db)
	byType, err := unmatched.CountByClaimType(ctx, d.Prefix, claimPrefix, outcome)
	if err != nil {
		return UnmatchedResponse{}, fmt.Errorf("views: listing unmatched claims: %w", err)
	}
	byWorkload, err := unmatched.CountByWorkload(ctx, d.Prefix, claimPrefix, outcome)
	if err != nil {
		return UnmatchedResponse{}, fmt.Errorf("views: listing unmatched claims: %w", err)
	}
	rows, err := unmatched.List(ctx, d.Prefix, claimPrefix, outcome, unmatchedRowLimit)
	if err != nil {
		return UnmatchedResponse{}, fmt.Errorf("views: listing unmatched claims: %w", err)
	}
	measures, err := measure.NewTableAccess(v.db).ListBySubmission(ctx, d.Prefix)
	if err != nil {
		return UnmatchedResponse{}, fmt.Errorf("views: listing unmatched claims: %w", err)
	}
	var macClaims int64
	for _, m := range measures {
		if m.ClaimType == measure.All && m.MACWorkloadID == measure.All {
			macClaims = m.MACClaims
		}
	}

	out := UnmatchedResponse{
		Outcome:     outcome,
		ByClaimType: make([]CountByClaimType, 0, len(byType)),
		ByWorkload:  make([]CountByWorkload, 0, len(byWorkload)),
		Rows:        make([]UnmatchedRow, 0, len(rows)),
	}
	for _, c := range byType {
		out.Count += c.Count
		out.ByClaimType = append(out.ByClaimType, CountByClaimType{ClaimType: c.Value, Count: c.Count})
	}
	for _, c := range byWorkload {
		out.ByWorkload = append(out.ByWorkload, CountByWorkload{Workload: c.Value, Count: c.Count})
	}
	out.Share = ratio(out.Count, macClaims)
	for _, c := range rows {
		out.Rows = append(out.Rows, UnmatchedRow{
			MatchKey:             hex.EncodeToString(c.MatchKey),
			Key:                  keyOfClaim(c),
			ClaimType:            c.ClaimType,
			PatientControlNumber: c.PatientControlNumber,
			Paid:                 c.ClaimPaymentAmount,
		})
	}
	return out, nil
}

// ratio writes n/d as a decimal fraction rounded to six places, without
// trailing zeros: 8/1000 is "0.008". Zero over zero is "0".
func ratio(n, d int64) string {
	if d == 0 {
		return "0"
	}
	s := big.NewRat(n, d).FloatString(6)
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

// ---------------------------------------------------------------- codes

// CodesResponse is the CARC and RARC lists, each code's entry under its
// code: { "carc": { [code]: Entry }, "rarc": { [code]: Entry } }.
type CodesResponse struct {
	CARC map[string]carc.Entry `json:"carc"`
	RARC map[string]rarc.Entry `json:"rarc"`
}

// Codes returns every entry of backend/pkg/carc and backend/pkg/rarc, keyed
// by code. Static reference data: unlike every other view, it reads no
// table and does not depend on which submissions are loaded. The maps are
// built on each call, which the one view that serves them does rarely.
func (v Controller) Codes(ctx context.Context) (CodesResponse, error) {
	carcEntries := carc.Entries()
	rarcEntries := rarc.Entries()
	out := CodesResponse{
		CARC: make(map[string]carc.Entry, len(carcEntries)),
		RARC: make(map[string]rarc.Entry, len(rarcEntries)),
	}
	for _, e := range carcEntries {
		out.CARC[e.Code] = e
	}
	for _, e := range rarcEntries {
		out.RARC[e.Code] = e
	}
	return out, nil
}
