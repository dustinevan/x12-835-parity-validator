package service

import (
	"sort"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/comparison"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/difference"
)

// Summary counts one comparison run's rows by outcome. SetAside counts
// keys (comparison rows), not the MAC claims sharing them; Duplicated
// counts keys too, each one MAC claim.
type Summary struct {
	Equal       int
	Different   int
	Missing     int
	Duplicated  int
	Extra       int
	SetAside    int
	Differences int
}

// CompareSubmissions matches a vendor submission's claims to the MAC
// submission's by match key and builds the comparison and difference rows
// under the compared set:
//
//	set_aside   the key occurs more than once in the MAC submission; nothing
//	            is compared, whatever the vendor sent for it
//	missing     the key is on one MAC claim and on no vendor claim
//	duplicated  the key is on one MAC claim and on more than one vendor
//	            claim, so no copy can be chosen; nothing is compared
//	extra       the key is in the vendor submission only (one row per key,
//	            however many vendor claims carry it)
//	equal / different  one claim on each side, compared element by element
//
// Rows come out in MAC file order, then extras in vendor file order, so
// the same inputs always give the same rows. The compare stage's SQL
// (lib/core/query CompareQuery) applies the same rule per workload.
func CompareSubmissions(set CompareSet, vendorPrefix string, macClaims, vendorClaims []claim.Claim) ([]comparison.Comparison, []difference.Difference, Summary, error) {
	group := func(claims []claim.Claim) (map[string][]claim.Claim, []string) {
		byKey := map[string][]claim.Claim{}
		var order []string
		for _, c := range claims {
			k := string(c.MatchKey)
			if _, ok := byKey[k]; !ok {
				order = append(order, k)
			}
			byKey[k] = append(byKey[k], c)
		}
		return byKey, order
	}
	macByKey, macOrder := group(macClaims)
	vendorByKey, vendorOrder := group(vendorClaims)

	var comparisons []comparison.Comparison
	var differences []difference.Difference
	var summary Summary
	for _, k := range macOrder {
		m := macByKey[k]
		v := vendorByKey[k]
		var row comparison.Comparison
		switch {
		case len(m) > 1:
			row = comparisonRow(vendorPrefix, m[0])
			row.Outcome = comparison.OutcomeSetAside
			summary.SetAside++
		case len(v) == 0:
			row = comparisonRow(vendorPrefix, m[0])
			row.Outcome = comparison.OutcomeMissing
			summary.Missing++
		case len(v) > 1:
			row = comparisonRow(vendorPrefix, m[0])
			row.Outcome = comparison.OutcomeDuplicated
			summary.Duplicated++
		default:
			var diffs []difference.Difference
			var err error
			row, diffs, err = ComparePair(set, vendorPrefix, m[0], v[0])
			if err != nil {
				return nil, nil, Summary{}, err
			}
			if row.Outcome == comparison.OutcomeEqual {
				summary.Equal++
			} else {
				summary.Different++
			}
			differences = append(differences, diffs...)
		}
		comparisons = append(comparisons, row)
	}
	for _, k := range vendorOrder {
		if _, ok := macByKey[k]; ok {
			continue
		}
		v := vendorByKey[k][0]
		row := comparisonRow(vendorPrefix, v)
		row.MACClaimStatusCode = ""
		row.MACClaimPaymentAmount = ""
		row.VendorClaimStatusCode = v.ClaimStatusCode
		row.VendorClaimPaymentAmount = v.ClaimPaymentAmount
		row.Outcome = comparison.OutcomeExtra
		summary.Extra++
		comparisons = append(comparisons, row)
	}
	summary.Differences = len(differences)
	return comparisons, differences, summary, nil
}

// ComparePair compares one MAC claim with the one vendor claim that shares
// its match key under the compared set and returns the comparison row
// (outcome equal or different, the vendor's status and payment, the
// difference count, and the tags) and the pair's difference rows.
// CompareSubmissions and the compare stage's per-workload stream
// (lib/core/control/submission) both compare their pairs with it, so there
// is one pair logic.
func ComparePair(set CompareSet, vendorPrefix string, mac, vendor claim.Claim) (comparison.Comparison, []difference.Difference, error) {
	diffs, err := DiffClaims(set, mac, vendor)
	if err != nil {
		return comparison.Comparison{}, nil, err
	}
	row := comparisonRow(vendorPrefix, mac)
	row.VendorClaimStatusCode = vendor.ClaimStatusCode
	row.VendorClaimPaymentAmount = vendor.ClaimPaymentAmount
	row.DifferenceCount = len(diffs)
	row.Tags = tags(diffs)
	row.Outcome = comparison.OutcomeEqual
	if len(diffs) > 0 {
		row.Outcome = comparison.OutcomeDifferent
	}
	var differences []difference.Difference
	for _, d := range diffs {
		differences = append(differences, difference.Difference{
			SubmissionPrefix:  vendorPrefix,
			MatchKey:          mac.MatchKey,
			Element:           d.Element,
			Occurrence:        d.Occurrence,
			LineControlNumber: d.LineControlNumber,
			MACValue:          d.MACValue,
			VendorValue:       d.VendorValue,
			AmountDifference:  d.AmountDifference,
		})
	}
	return row, differences, nil
}

// comparisonRow fills the key and MAC-side columns from one claim.
func comparisonRow(vendorPrefix string, c claim.Claim) comparison.Comparison {
	return comparison.Comparison{
		SubmissionPrefix:              vendorPrefix,
		MatchKey:                      c.MatchKey,
		MACWorkloadID:                 c.MACWorkloadID,
		ClaimType:                     c.ClaimType,
		MedicareBeneficiaryIdentifier: c.MedicareBeneficiaryIdentifier,
		DateOfService:                 c.DateOfService,
		TotalClaimChargeAmount:        c.TotalClaimChargeAmount,
		MACClaimStatusCode:            c.ClaimStatusCode,
		MACClaimPaymentAmount:         c.ClaimPaymentAmount,
		ClaimFrequencyCode:            c.ClaimFrequencyCode,
		Adjudication:                  c.Adjudication,
		Jurisdiction:                  c.Jurisdiction,
		Part:                          c.Part,
	}
}

// tags are the distinct differing elements, sorted.
func tags(diffs []Difference) []string {
	seen := map[string]bool{}
	var out []string
	for _, d := range diffs {
		if !seen[d.Element] {
			seen[d.Element] = true
			out = append(out, d.Element)
		}
	}
	sort.Strings(out)
	return out
}
