package main

import (
	"fmt"
	"sort"
	"strings"
)

// The vendor side is built from a catalog of difference patterns: ways a
// replacement adjudication system's output can differ from the MAC's. Each
// pattern has a name, the taxonomy entries it corresponds to
// (docs/taxonomy), and one or more variants, each a change to the
// decisions on a claim in the population where that pattern would show. A
// vendor's profile gives each pattern a rate per submission, falling from
// one submission to the next, so later submissions carry fewer planted
// differences. Every planted difference is recorded in truth.jsonl with the
// pattern behind it.

// pattern is one catalog entry.
type pattern struct {
	Name     string
	Taxonomy []string
}

// catalog is numbered from 1; index 0 is unused. Patterns 6 (the vendor
// sending a claim twice), 9 (a chain reported as a net adjustment), and 11
// (a claim finalized in a later cycle) are decided per submission in
// buildSubmission; the rest are variants applied per claim.
var catalog = [...]pattern{
	{},
	{"fee schedule and locality pricing", []string{"T2.10", "T4.18"}},
	{"sequestration rounding", []string{"T2.09", "T3.14"}},
	{"deductible accumulator state", []string{"T3.09"}},
	{"multiple-procedure reduction", []string{"T3.19", "T4.29"}},
	{"NCCI and MUE edits", []string{"T3.20", "T4.30", "T4.31"}},
	{"duplicate-claim detection", []string{"T3.21", "T1.06"}},
	{"inpatient grouping and outliers", []string{"T2.19", "T2.24", "T2.25"}},
	{"MAC-local edits and coverage decisions", []string{"T2.01", "T3.07"}},
	{"replacement and void chains", []string{"T1.10", "T1.13"}},
	{"remark and placement conventions", []string{"T3.34", "T3.36", "T2.53"}},
	{"payment cycle timing", []string{"T1.01"}},
	{"front-end strictness", []string{"T1.05", "T2.03"}},
}

// profile describes one vendor: the rate of each pattern at the first,
// second, and third-or-later submission, plus how it numbers and groups
// its output. Rates are per claim in the pattern's slice, times each
// variant's weight. Missing, Duplicate, Chain, and Extra are per claim of
// the whole set.
type profile struct {
	Rates     [len(catalog)][3]float64
	Missing   [3]float64 // pattern 11: claims the vendor finalizes in a later cycle
	Duplicate [3]float64 // pattern 6: claims the vendor sends twice
	Chain     [3]float64 // pattern 9: chains reported as one net adjustment
	Extra     [3]float64 // claims only the vendor has
	ICN       string
	ISABase   int
	Time      string
	ByMonth   bool // payments per payee per received month; otherwise one payment per payee
	STBase    int
	Trace     string
}

// The vendor-1 profile plants the pricing patterns (1, 2, 4, 7); the
// vendor-2 profile the edit and history patterns (3, 5, 6, 8); both the
// convention, timing, and front-end patterns (9 to 12). Each later
// submission drops some and keeps the rest. Vendor 2 also carries a little
// of pattern 2 (the truncation variant), so its files carry cent
// differences too, which the comparison counts as the same amount
// (matchCents). The profiles are synthetic and describe no real system.
var profiles = []profile{
	{
		Rates: [len(catalog)][3]float64{
			1: {0.18, 0.08, 0.02}, 2: {0.25, 0.08, 0}, 4: {0.45, 0.15, 0}, 7: {0.40, 0.15, 0},
			10: {0.14, 0.10, 0.06}, 12: {0.015, 0.006, 0},
		},
		Missing:   [3]float64{0.004, 0.002, 0.001},
		Duplicate: [3]float64{0, 0, 0},
		Chain:     [3]float64{0.5, 0.2, 0},
		Extra:     [3]float64{0.003, 0.002, 0.001},
		ICN:       "CA",
		ISABase:   700000000,
		Time:      "0415",
		STBase:    1,
		Trace:     "CA",
	},
	{
		Rates: [len(catalog)][3]float64{
			2: {0.06, 0.03, 0.01}, 3: {0.10, 0.06, 0.02}, 5: {0.25, 0.15, 0.04}, 6: {0.6, 0.4, 0.1},
			8: {0.9, 0.5, 0.1}, 10: {0.14, 0.12, 0.08}, 12: {0.01, 0.01, 0.004},
		},
		Missing:   [3]float64{0.006, 0.003, 0.002},
		Duplicate: [3]float64{0.005, 0.003, 0},
		Chain:     [3]float64{0.4, 0.2, 0.1},
		Extra:     [3]float64{0.004, 0.003, 0.002},
		ICN:       "B",
		ISABase:   800000000,
		Time:      "2130",
		ByMonth:   true,
		STBase:    9001,
		Trace:     "88",
	},
}

// stage maps a submission to its column in the rate schedule.
func stage(d int) int {
	return min(d, 2)
}

// regressRate is the chance that a pattern dropped in an earlier submission
// comes back in a later one.
const regressRate = 0.08

// variant is one way a pattern shows on a claim. Apply returns the claim
// with the vendor's decision and the taxonomy IDs it produces (the first is
// the one a difference is recorded under), or nil IDs when the claim is
// outside the variant's slice. Weight scales the pattern's rate for this
// variant. RemarkOnly variants change only remarks, so they cannot show on
// a reversal.
type variant struct {
	Pattern    int
	Weight     float64
	RemarkOnly bool
	Apply      func(c Claim) (Claim, []string)
}

var variants = []variant{
	// 1. Fee schedule and locality pricing.
	{1, 1, false, variantLocality},
	{1, 1, false, variantAmbulanceRural},
	{1, 0.3, false, variantOutpatientPackaging},
	// 2. Sequestration rounding.
	{2, 1, false, variantSequestrationPerLine},
	{2, 1, false, variantSequestrationTruncated},
	// 3. Deductible accumulator state.
	{3, 1, false, variantDeductibleElsewhere},
	{3, 0.3, false, variantLabCostSharing},
	// 4. Multiple-procedure reduction.
	{4, 1, false, variantMultipleProcedure},
	// 5. NCCI and MUE edits.
	{5, 1, false, variantBundlingEdit},
	{5, 0.5, false, variantUnitLimit},
	{5, 1, false, variantOxygenMonth},
	// 6. Duplicate-claim detection (the other half, a claim sent twice, is per submission).
	{6, 1, false, variantDuplicateOccurrence},
	// 7. Inpatient grouping and outliers.
	{7, 1, false, variantGrouperVersion},
	{7, 1, false, variantOutlierDropped},
	// 8. MAC-local edits and coverage decisions.
	{8, 1, false, variantLocalCoverage},
	{8, 0.4, false, variantDifferentReason},
	{8, 0.15, false, variantLiableParty},
	// 10. Remark and placement conventions.
	{10, 1, true, variantRemarkChoice},
	{10, 0.3, true, variantRemarkLevel},
	{10, 0.15, false, variantZeroAdjustment},
	{10, 0.2, false, variantRenderingToggle},
	{10, 0.3, true, variantCrossoverMissing},
	{10, 1, true, variantInterestLevel},
	{10, 0.05, true, variantAlertRemark},
	{10, 0.03, false, variantReceivedDate},
	{10, 0.03, false, variantAllowedMissing},
	{10, 0.03, false, variantUnitsMissing},
	{10, 0.1, false, variantStatementDate},
	{10, 0.1, false, variantRevenueCode},
	// 12. Front-end strictness.
	{12, 1, false, variantUnprocessable},
}

// Pattern 1: Part B lines priced at a different locality or an older fee
// schedule than the MAC's, 1 to 5 percent apart.
func variantLocality(c Claim) (Claim, []string) {
	if c.Type == "institutional" || c.Paid <= 0 {
		return c, nil
	}
	c = c.clone()
	c.FeeScale = []int64{-4, -3, -2, -1, 1, 2, 3, 5}[c.Pick%8]
	return c, []string{"T2.10", "T4.18"}
}

// Pattern 1: ambulance mileage priced at the urban rate for a rural pickup.
func variantAmbulanceRural(c Claim) (Claim, []string) {
	if !c.Rural {
		return c, nil
	}
	c = c.clone()
	for i, ln := range c.Lines {
		if ln.Class == "mile" {
			c.Lines[i].UnitFee = ln.UrbanFee
		}
	}
	return c, []string{"T4.56", "T2.10"}
}

// Pattern 1: outpatient packaging differs on one line of a 13x bill.
func variantOutpatientPackaging(c Claim) (Claim, []string) {
	if c.BillType != "131" {
		return c, nil
	}
	var cand []int
	for i, ln := range c.Lines {
		if i > 0 && ln.UnitFee > 0 && ln.Denial.Reason == "" {
			cand = append(cand, i)
		}
	}
	if len(cand) == 0 {
		return c, nil
	}
	c = c.clone()
	i := cand[c.Pick%len(cand)]
	c.Lines[i].Packaged = !c.Lines[i].Packaged
	return c, []string{"T4.32", "T4.40", "T2.28"}
}

// Pattern 2: sequestration rounded per line instead of per claim.
func variantSequestrationPerLine(c Claim) (Claim, []string) {
	if len(c.Lines) < 2 || c.Paid <= 0 {
		return c, nil
	}
	c = c.clone()
	c.SeqPerLine = true
	return c, []string{"T2.09", "T3.14", "T3.16", "T4.19"}
}

// Pattern 2: the fraction of a cent truncated instead of rounded.
func variantSequestrationTruncated(c Claim) (Claim, []string) {
	if c.inpatient() || c.Paid <= 0 || c.SeqPerLine {
		return c, nil
	}
	c = c.clone()
	c.SeqTruncate = true
	return c, []string{"T2.09", "T3.14"}
}

// Pattern 3: the deductible on one side only: the vendor applies it where
// the MAC treated it as met on an earlier claim, or treats it as met where
// the MAC applied it.
func variantDeductibleElsewhere(c Claim) (Claim, []string) {
	if c.Type != "professional" {
		return c, nil
	}
	c = c.clone()
	if c.Deductible > 0 {
		c.Deductible = 0
	} else {
		c.Deductible = c.AltDeductible
	}
	return c, []string{"T3.09", "T2.11"}
}

// Pattern 3: cost sharing applied to a clinical laboratory line that is
// exempt from it.
func variantLabCostSharing(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		if ln.Class == "lab" && ln.Exempt && ln.Allowed > 0 {
			c = c.clone()
			c.Lines[i].Exempt = false
			return c, []string{"T3.11", "T4.37"}
		}
	}
	return c, nil
}

// Pattern 4: the multiple-procedure reduction on a different line, or on none.
func variantMultipleProcedure(c Claim) (Claim, []string) {
	reduced, top := -1, -1
	for i, ln := range c.Lines {
		if ln.Class != "surg" {
			continue
		}
		if ln.Reduced && reduced == -1 {
			reduced = i
		}
		if !ln.Reduced {
			top = i
		}
	}
	if reduced == -1 {
		return c, nil
	}
	c = c.clone()
	if c.Pick%2 == 0 || top == -1 {
		for i := range c.Lines {
			c.Lines[i].Reduced = false
		}
	} else {
		c.Lines[reduced].Reduced = false
		c.Lines[top].Reduced = true
	}
	return c, []string{"T4.29", "T3.19", "T2.54"}
}

// Pattern 5: a line the MAC denied as bundled is paid by the vendor, or a
// second paid line is bundled into the first.
func variantBundlingEdit(c Claim) (Claim, []string) {
	if c.inpatient() {
		return c, nil
	}
	for i, ln := range c.Lines {
		if ln.Denial.Reason == "97" {
			c = c.clone()
			c.Lines[i].Denial = Denial{}
			return c, []string{"T4.30", "T3.20", "T4.51"}
		}
	}
	paid := 0
	for _, ln := range c.Lines {
		if ln.Denial.Reason == "" && !ln.Packaged && ln.UnitFee > 0 {
			paid++
		}
	}
	if paid < 2 || c.Pick%2 == 1 {
		return c, nil
	}
	c = c.clone()
	seen := 0
	for i, ln := range c.Lines {
		if ln.Denial.Reason == "" && !ln.Packaged && ln.UnitFee > 0 {
			seen++
			if seen == 2 {
				c.Lines[i].Denial = Denial{"CO", "97", "M15"}
				break
			}
		}
	}
	return c, []string{"T4.30", "T3.20", "T4.51"}
}

// Pattern 5: a medically unlikely edit applied to a multi-unit line the
// MAC paid.
func variantUnitLimit(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		if ln.Units > 2 && ln.Denial.Reason == "" && ln.UnitFee > 0 && ln.Class != "mile" && ln.Class != "pps" && ln.Class != "zero" {
			c = c.clone()
			c.Lines[i].Denial = Denial{"CO", "151", "N362"}
			return c, []string{"T4.31", "T3.20"}
		}
	}
	return c, nil
}

// Pattern 5: the oxygen rental month counted differently: paid on one
// side, denied past the cap on the other.
func variantOxygenMonth(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		if ln.Class != "oxygen" {
			continue
		}
		switch ln.Denial.Reason {
		case "119":
			c = c.clone()
			c.Lines[i].Denial = Denial{}
			return c, []string{"T4.46", "T3.22"}
		case "":
			c = c.clone()
			c.Lines[i].Denial = Denial{"CO", "119", "N640"}
			return c, []string{"T4.46", "T3.22"}
		}
	}
	return c, nil
}

// Pattern 6: the duplicate check picks a different occurrence of a key
// collision: the vendor pays the resubmission the MAC denied, or denies the
// one the MAC paid.
func variantDuplicateOccurrence(c Claim) (Claim, []string) {
	switch c.Pair {
	case "resub-first", "distinct-second":
		c = c.clone()
		for i := range c.Lines {
			c.Lines[i].Denial = Denial{"CO", "18", ""}
		}
		if c.Pair == "resub-first" {
			return c, []string{"T1.07", "T3.21"}
		}
		return c, []string{"T4.53", "T3.21"}
	case "resub-second":
		c = c.clone()
		for i := range c.Lines {
			c.Lines[i].Denial = Denial{}
		}
		return c, []string{"T1.07", "T3.21"}
	}
	return c, nil
}

// Pattern 7: a different DRG grouper version moves a boundary DRG to its pair.
func variantGrouperVersion(c Claim) (Claim, []string) {
	if !c.inpatient() || c.Denial.Reason != "" {
		return c, nil
	}
	for _, d := range drgs {
		if d.Code == c.DRG && d.AltCode != "" {
			c = c.clone()
			c.DRG = d.AltCode
			c.Weight = d.AltWeight
			return c, []string{"T2.19", "T2.24"}
		}
	}
	return c, nil
}

// Pattern 7: the cost outlier the MAC paid is not paid.
func variantOutlierDropped(c Claim) (Claim, []string) {
	if !c.inpatient() || !c.Outlier || c.Denial.Reason != "" {
		return c, nil
	}
	c = c.clone()
	c.OutlierDropped = true
	return c, []string{"T2.25", "T2.24"}
}

// Pattern 8: the vendor pays lines the MAC denied under a jurisdiction's
// local coverage determination.
func variantLocalCoverage(c Claim) (Claim, []string) {
	if !lcdJurisdictions[c.Juris] || c.Part != "B" {
		return c, nil
	}
	c = c.clone()
	found := false
	for i, ln := range c.Lines {
		if lcdCodes[ln.Code] && ln.Denial == (Denial{"CO", "50", "N115"}) {
			c.Lines[i].Denial = Denial{}
			found = true
		}
	}
	if !found {
		return c, nil
	}
	return c, []string{"T2.01", "T3.37", "T4.51"}
}

// Pattern 8: the same denial under a different reason code.
func variantDifferentReason(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		if ln.Denial == (Denial{"CO", "50", "N115"}) {
			c = c.clone()
			c.Lines[i].Denial = Denial{"CO", "96", "N130"}
			return c, []string{"T3.07"}
		}
	}
	return c, nil
}

// Pattern 8: the same reason assigned to the other liable party.
func variantLiableParty(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		switch ln.Denial.Reason {
		case "96", "50", "27":
			c = c.clone()
			if ln.Denial.Group == "CO" {
				c.Lines[i].Denial.Group = "PR"
			} else {
				c.Lines[i].Denial.Group = "CO"
			}
			return c, []string{"T3.04", "T2.12"}
		}
	}
	return c, nil
}

var remarkSynonyms = map[string]string{
	"M15":  "N19",
	"N20":  "N19",
	"N115": "N386",
	"N640": "N362",
}

// Pattern 10: remark codes ordered or chosen differently.
func variantRemarkChoice(c Claim) (Claim, []string) {
	if c.Type != "professional" || len(c.Remarks) == 0 {
		return c, nil
	}
	c = c.clone()
	if len(c.Remarks) >= 2 && c.Pick%2 == 0 {
		c.RemarkVariant = 1
		return c, []string{"T3.36", "T2.27"}
	}
	for i, ln := range c.Lines {
		if s, ok := remarkSynonyms[ln.Denial.Remark]; ok {
			c.Lines[i].Denial.Remark = s
			return c, []string{"T3.31"}
		}
	}
	if c.Remarks[0] == "MA01" && c.Pick%3 == 0 {
		c.RemarkVariant = 2
		return c, []string{"T3.30", "T2.27"}
	}
	return c, nil
}

// Pattern 10: a remark at claim level on one side and line level on the other.
func variantRemarkLevel(c Claim) (Claim, []string) {
	if c.inpatient() || len(c.Lines) == 0 {
		return c, nil
	}
	hasLine := false
	for _, ln := range c.Lines {
		if ln.Denial.Remark != "" {
			hasLine = true
		}
	}
	c = c.clone()
	switch {
	case hasLine && c.Pick%2 == 0:
		c.RemarkLevel = 1
	case len(c.Remarks) > 0 || c.PatientResp > 0:
		c.RemarkLevel = 2
	default:
		return c, nil
	}
	return c, []string{"T3.34"}
}

// Pattern 10: a zero-amount adjustment on one side only.
func variantZeroAdjustment(c Claim) (Claim, []string) {
	if c.Paid <= 0 || c.inpatient() {
		return c, nil
	}
	c = c.clone()
	c.ZeroCAS = true
	return c, []string{"T3.28", "T2.53"}
}

// Pattern 10: the rendering provider loop sent where the MAC leaves it
// out, or left out where the MAC sends it.
func variantRenderingToggle(c Claim) (Claim, []string) {
	if c.Rendering.NPI == "" {
		return c, nil
	}
	c = c.clone()
	c.OmitRendering = !c.OmitRendering
	return c, []string{"T2.44", "T2.53"}
}

// Pattern 10: the crossover remark and carrier missing.
func variantCrossoverMissing(c Claim) (Claim, []string) {
	if !c.CrossoverShown {
		return c, nil
	}
	c = c.clone()
	c.DropCrossover = true
	return c, []string{"T3.33", "T1.21"}
}

// Pattern 10: interest reported only at provider level (PLB), not on the claim.
func variantInterestLevel(c Claim) (Claim, []string) {
	if c.Interest == 0 {
		return c, nil
	}
	c = c.clone()
	c.OmitInterest = true
	return c, []string{"T5.24", "T2.32"}
}

// Pattern 10: an alert remark on one side only.
func variantAlertRemark(c Claim) (Claim, []string) {
	if c.inpatient() {
		return c, nil
	}
	c = c.clone()
	c.ExtraAlert = true
	return c, []string{"T3.32"}
}

// Pattern 10: the received date off by a day.
func variantReceivedDate(c Claim) (Claim, []string) {
	c = c.clone()
	c.Received = addDays(c.Received, 1)
	return c, []string{"T2.38"}
}

// Pattern 10: the line allowed amount not reported.
func variantAllowedMissing(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		if ln.Allowed > 0 {
			c = c.clone()
			c.Lines[i].OmitAllowed = true
			return c, []string{"T4.21"}
		}
	}
	return c, nil
}

// Pattern 10: a unit count of one left out.
func variantUnitsMissing(c Claim) (Claim, []string) {
	for i, ln := range c.Lines {
		if ln.Units == 1 {
			c = c.clone()
			c.Lines[i].OmitUnits = true
			return c, []string{"T4.15"}
		}
	}
	return c, nil
}

// Pattern 10: the statement through date off by a day.
func variantStatementDate(c Claim) (Claim, []string) {
	if c.Type != "institutional" {
		return c, nil
	}
	c = c.clone()
	c.Thru = addDays(c.Thru, 1)
	return c, []string{"T2.37"}
}

// Pattern 10: a revenue code off by one in its last digit.
func variantRevenueCode(c Claim) (Claim, []string) {
	if c.Type != "institutional" {
		return c, nil
	}
	for i, ln := range c.Lines {
		if ln.Rev != "" {
			c = c.clone()
			last := "1"
			if strings.HasSuffix(ln.Rev, "1") {
				last = "0"
			}
			c.Lines[i].Rev = ln.Rev[:3] + last
			return c, []string{"T4.12"}
		}
	}
	return c, nil
}

// Pattern 12: the vendor returns as unprocessable a claim the MAC adjudicated.
func variantUnprocessable(c Claim) (Claim, []string) {
	if c.Type == "institutional" || c.Unprocessable || c.NotOurs || c.Pair != "" {
		return c, nil
	}
	for _, ln := range c.Lines {
		if ln.Denial.Reason != "" {
			return c, nil
		}
	}
	c = c.clone()
	c.Unprocessable = true
	remark := unprocessableRemarks[c.Pick%len(unprocessableRemarks)]
	for i := range c.Lines {
		c.Lines[i].Denial = Denial{"CO", "16", remark}
	}
	return c, []string{"T2.03", "T1.05"}
}

// variantActive reports whether a variant applies to a claim at submission
// d: its draw is under the submission's rate, or it applied at an earlier
// submission and this submission's regression draw brings it back.
func variantActive(prof profile, v variant, draw vendorDraw, d, id, vi int) bool {
	rates := prof.Rates[v.Pattern]
	u := draw.Gap[id][vi]
	if u < rates[stage(d)]*v.Weight {
		return true
	}
	for e := 0; e < d; e++ {
		if u < rates[stage(e)]*v.Weight {
			return draw.Regress[d][id][vi] < regressRate
		}
	}
	return false
}

// planted is one variant that changed a claim's output: the pattern it
// belongs to, the taxonomy ID a difference is recorded under, and the
// elements it changed.
type planted struct {
	Pattern  string
	Taxonomy string
	IDs      []string
	Keys     map[string]bool
}

// vendorClaim applies the vendor's active variants to one reference claim,
// one at a time, and keeps a variant only when it changes the output. It
// returns the adjudicated vendor claim and what was planted, in order.
func vendorClaim(prof profile, draw vendorDraw, d int, c Claim) (Claim, []planted) {
	cur := c.clone()
	if c.Reversal {
		// Decide on the positive claim, then reverse it, so the vendor's
		// reversal mirrors the vendor's original.
		cur.Reversal = false
		cur = adjudicate(cur)
	}
	var out []planted
	for vi, v := range variants {
		if c.Reversal && v.RemarkOnly {
			continue
		}
		if !variantActive(prof, v, draw, d, c.DrawID, vi) {
			continue
		}
		next, ids := v.Apply(cur)
		if ids == nil {
			continue
		}
		next = adjudicate(next)
		changed := differences(claimFields(cur), claimFields(next))
		if len(changed) == 0 {
			continue
		}
		p := planted{Pattern: catalog[v.Pattern].Name, Taxonomy: ids[0], IDs: ids, Keys: map[string]bool{}}
		for _, df := range changed {
			p.Keys[df.Element+"|"+df.Line] = true
		}
		out = append(out, p)
		cur = next
	}
	if c.Reversal {
		cur.Reversal = true
		cur = adjudicate(cur)
	}
	return cur, out
}

// matchCents is the validator's five-cent rule (AmountMatchCents in
// backend/lib/core/service/compare.go): two dollar amounts, both sent,
// whose difference is this many cents or less either way are the same
// amount, and the comparison records no difference for them. gen keeps its
// own copy because it imports nothing from the validator.
const matchCents = 5

// field is one comparable value in a claim's output. A dollar amount also
// carries its cents, for the five-cent rule.
type field struct {
	Element string
	Line    string
	Value   string
	IsMoney bool
	Amount  Money
}

// text is a field compared exactly: a code, an identifier, a date, a count.
func text(element, line, value string) field {
	return field{Element: element, Line: line, Value: value}
}

// money is a field holding a dollar amount, compared under the five-cent
// rule.
func money(element, line string, amount Money) field {
	return field{Element: element, Line: line, Value: amount.String(), IsMoney: true, Amount: amount}
}

// claimFields lists the values a comparison sees in one claim, leaving out
// the system-assigned ones (CLP07, and REF*F8, which echoes an ICN) and
// REF*6R, which pairs the lines rather than being compared.
func claimFields(c Claim) []field {
	f := []field{
		text("CLP01", "", c.PatientCtl),
		text("CLP02", "", c.Status),
		money("CLP03", "", c.Charge),
		money("CLP04", "", c.Paid),
		money("CLP05", "", c.PatientResp),
		text("CLP06", "", filingIndicator(c)),
		text("CLP08", "", c.Facility),
		text("NM1*QC", "", c.Patient.MBI),
		text("DTM*050", "", c.Received),
	}
	if c.Type == "institutional" {
		f = append(f, text("CLP09", "", c.Frequency))
	}
	if c.inpatient() {
		f = append(f,
			text("CLP11", "", c.DRG),
			text("CLP12", "", weightString(c.Weight)),
			money("MIA04", "", c.DRGAmount),
			text("MIA", "", strings.Join(c.Remarks, " ")),
		)
		if c.OutlierAmount != 0 {
			f = append(f, money("MIA02", "", c.OutlierAmount))
		}
	} else {
		f = append(f, text("MOA", "", strings.Join(c.Remarks, " ")))
	}
	for _, a := range c.Adjs {
		f = append(f, money("CAS*"+a.Group+"*"+a.Reason, "", a.Amount))
	}
	if c.Rendering.NPI != "" && !c.OmitRendering {
		f = append(f, text("NM1*82", "", c.Rendering.NPI))
	}
	if c.CrossoverShown {
		f = append(f, text("NM1*TT", "", c.Patient.InsurerID))
	}
	if c.Type == "institutional" {
		f = append(f, text("DTM*232", "", c.From), text("DTM*233", "", c.Thru))
	}
	if c.Interest > 0 && !c.OmitInterest {
		f = append(f, money("AMT*I", "", c.Interest))
	}
	for _, ln := range c.Lines {
		f = append(f,
			text("SVC01", ln.ID, procedure(ln)),
			money("SVC02", ln.ID, ln.Charge),
			money("SVC03", ln.ID, ln.Paid),
			text("SVC05", ln.ID, unitsString(ln)),
			text("DTM*472", ln.ID, ln.Date),
			text("LQ", ln.ID, strings.Join(ln.Remarks, " ")),
		)
		if c.Type == "institutional" {
			f = append(f, text("SVC04", ln.ID, ln.Rev))
		}
		if ln.Allowed != 0 && !ln.OmitAllowed {
			f = append(f, money("AMT*B6", ln.ID, ln.Allowed))
		}
		for _, a := range ln.Adjs {
			f = append(f, money("CAS*"+a.Group+"*"+a.Reason, ln.ID, a.Amount))
		}
	}
	return f
}

// difference is one element that differs between the two sides, with the
// pattern it was planted by.
type difference struct {
	Element  string `json:"element"`
	Line     string `json:"line,omitempty"`
	MAC      string `json:"mac"`
	Vendor   string `json:"vendor"`
	Pattern  string `json:"pattern,omitempty"`
	Taxonomy string `json:"taxonomy,omitempty"`
}

// differences lists every value that differs between two claims' fields,
// exactly: what a variant changed in a claim's output.
func differences(mac, vendor []field) []difference {
	key := func(f field) string { return f.Element + "|" + f.Line }
	vendorValue := map[string]string{}
	for _, f := range vendor {
		vendorValue[key(f)] = f.Value
	}
	seen := map[string]bool{}
	var out []difference
	for _, f := range mac {
		seen[key(f)] = true
		if v := vendorValue[key(f)]; v != f.Value {
			out = append(out, difference{Element: f.Element, Line: f.Line, MAC: f.Value, Vendor: v})
		}
	}
	for _, f := range vendor {
		if !seen[key(f)] && f.Value != "" {
			out = append(out, difference{Element: f.Element, Line: f.Line, MAC: "", Vendor: f.Value})
		}
	}
	return out
}

// recordedDifferences lists the differences the comparison records between
// two claims: every difference, less the dollar amounts both sides sent
// that are within matchCents of each other. An amount sent on one side
// only stays a difference.
func recordedDifferences(mac, vendor []field) []difference {
	key := func(element, line string) string { return element + "|" + line }
	macField := map[string]field{}
	for _, f := range mac {
		macField[key(f.Element, f.Line)] = f
	}
	vendorField := map[string]field{}
	for _, f := range vendor {
		vendorField[key(f.Element, f.Line)] = f
	}
	var out []difference
	for _, df := range differences(mac, vendor) {
		m, onMAC := macField[key(df.Element, df.Line)]
		v, onVendor := vendorField[key(df.Element, df.Line)]
		if onMAC && onVendor && m.IsMoney && v.IsMoney && max(m.Amount-v.Amount, v.Amount-m.Amount) <= matchCents {
			continue
		}
		out = append(out, df)
	}
	return out
}

// attribute names, on each difference, the first planted variant that
// changed that element.
func attribute(diffs []difference, plants []planted) []difference {
	for i, df := range diffs {
		k := df.Element + "|" + df.Line
		for _, p := range plants {
			if p.Keys[k] {
				diffs[i].Pattern = p.Pattern
				diffs[i].Taxonomy = p.Taxonomy
				break
			}
		}
	}
	return diffs
}

// claimKey is the six-field match key the comparison matches on.
type claimKey struct {
	MBI           string `json:"mbi"`
	DateOfService string `json:"date_of_service"`
	TotalCharge   string `json:"total_charge"`
	TaxID         string `json:"tax_id"`
	NPI           string `json:"npi"`
	Workload      string `json:"workload"`
}

func keyOf(ws []Workload, c Claim) claimKey {
	w := ws[c.Workload]
	p := w.Payees[c.Payee]
	return claimKey{c.Patient.MBI, c.dateOfService(), c.Charge.String(), p.TaxID, p.NPI, w.Number}
}

// truthRecord is one line of truth.jsonl: one vendor claim, its outcome as
// the comparison should find it, the patterns planted in it, the taxonomy
// entries they correspond to, and every difference the comparison should
// record. A planted amount within matchCents of the MAC's is the same
// amount to the comparison, so it is not listed: a claim whose only planted
// changes are such cents keeps its patterns and has no differences.
type truthRecord struct {
	Submission  string       `json:"submission"`
	Key         claimKey     `json:"key"`
	PatientCtl  string       `json:"clp01"`
	ClaimType   string       `json:"claim_type"`
	Outcome     string       `json:"outcome"`
	Patterns    []string     `json:"patterns"`
	Taxonomy    []string     `json:"taxonomy"`
	Differences []difference `json:"differences"`
}

func newRecord(name string, k claimKey, c Claim) truthRecord {
	return truthRecord{Submission: name, Key: k, PatientCtl: c.PatientCtl, ClaimType: c.Type, Patterns: []string{}, Taxonomy: []string{}, Differences: []difference{}}
}

// submission is one vendor's output for one submission date.
type submission struct {
	Name  string
	Files []File
	Truth []truthRecord
}

// buildSubmission produces vendor v's output for submission d and its truth.
func buildSubmission(cfg config, w world, v, d int) submission {
	prof := profiles[v%len(profiles)]
	draw := w.Draws[v]
	date := cfg.Submissions[d]
	out := submission{Name: cfg.Vendors[v] + "/" + date}
	s := stage(d)

	keyCount := map[claimKey]int{}
	for _, c := range w.Claims {
		keyCount[keyOf(w.Workloads, c)]++
	}
	chained := map[int]bool{}
	for _, c := range w.Claims {
		if c.OrigID >= 0 {
			chained[c.OrigID] = true
		}
	}
	plain := func(c Claim) bool {
		return c.Pair == "" && c.OrigID < 0 && !chained[c.ID] && c.Interest == 0 && keyCount[keyOf(w.Workloads, c)] == 1
	}

	// Pattern 11: claims the vendor finalizes in a later cycle, so they
	// are missing from the compared files. The ones received late in their
	// cycle are the ones at risk; the lowest draws among them, so a claim
	// missing in a later submission was missing in the earlier ones.
	var eligible []Claim
	for _, c := range w.Claims {
		if plain(c) && daysBetween(c.Received, w.Cycles[c.Cycle]) <= 18 {
			eligible = append(eligible, c)
		}
	}
	sort.SliceStable(eligible, func(i, j int) bool { return draw.Missing[eligible[i].ID] < draw.Missing[eligible[j].ID] })
	missing := map[int]bool{}
	for i := 0; i < countFor(cfg.Claims, prof.Missing[s]) && i < len(eligible); i++ {
		missing[eligible[i].ID] = true
	}

	// Pattern 9: a chain the vendor reports as one net adjustment: no
	// reversal, and the correction sent as a fresh claim.
	netChain := map[int]bool{}
	for _, c := range w.Claims {
		if c.OrigID >= 0 && draw.Chain[c.OrigID] < prof.Chain[s] {
			netChain[c.OrigID] = true
		}
	}

	var vendorClaims []Claim
	for _, c := range w.Claims {
		k := keyOf(w.Workloads, c)
		rec := newRecord(out.Name, k, c)
		if missing[c.ID] {
			rec.Outcome = "missing"
			rec.Patterns = []string{catalog[11].Name}
			rec.Taxonomy = []string{"T1.01"}
			out.Truth = append(out.Truth, rec)
			continue
		}
		if c.Reversal && netChain[c.OrigID] {
			rec.Outcome = "missing"
			rec.Patterns = []string{catalog[9].Name}
			rec.Taxonomy = []string{"T1.10"}
			out.Truth = append(out.Truth, rec)
			continue
		}
		vc, plants := vendorClaim(prof, draw, d, c)
		if c.Correction && netChain[c.OrigID] {
			vc.NetAdjust = true
			if vc.Type == "institutional" {
				vc.Frequency = "1"
			}
			plants = append(plants, planted{Pattern: catalog[9].Name, Taxonomy: "T1.13", IDs: []string{"T1.13"}, Keys: map[string]bool{"CLP09|": true}})
		}
		vendorClaims = append(vendorClaims, vc)
		rec.Outcome = "matched"
		if keyCount[k] > 1 {
			rec.Outcome = "set_aside"
			rec.Taxonomy = append(rec.Taxonomy, "T1.20")
		} else if plain(c) && draw.Duplicate[c.ID] < prof.Duplicate[s] {
			// Pattern 6: the same claim sent twice.
			vendorClaims = append(vendorClaims, vc)
			rec.Outcome = "duplicated"
			rec.Patterns = append(rec.Patterns, catalog[6].Name)
			rec.Taxonomy = append(rec.Taxonomy, "T1.06")
		}
		diffs := attribute(recordedDifferences(claimFields(c), claimFields(vc)), plants)
		for _, p := range plants {
			if !contains(rec.Patterns, p.Pattern) {
				rec.Patterns = append(rec.Patterns, p.Pattern)
			}
			rec.Taxonomy = append(rec.Taxonomy, p.IDs...)
		}
		for _, df := range diffs {
			if df.Element == "CLP02" {
				rec.Taxonomy = append(rec.Taxonomy, "T2.01")
			}
		}
		rec.Taxonomy = unique(rec.Taxonomy)
		if diffs != nil {
			rec.Differences = diffs
		}
		out.Truth = append(out.Truth, rec)
	}

	extras := w.Extras[v][:countFor(cfg.Claims, prof.Extra[s])]
	for _, c := range extras {
		vendorClaims = append(vendorClaims, c)
		rec := newRecord(out.Name, keyOf(w.Workloads, c), c)
		rec.Outcome = "extra"
		rec.Taxonomy = []string{"T1.02"}
		out.Truth = append(out.Truth, rec)
	}

	out.Files = vendorFiles(cfg, w, prof, v, d, vendorClaims)
	return out
}

// daysBetween returns the days from a to b, both CCYYMMDD.
func daysBetween(a, b string) int {
	return int(parseYMD(b).Sub(parseYMD(a)).Hours() / 24)
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func unique(list []string) []string {
	out := []string{}
	for _, x := range list {
		if !contains(out, x) {
			out = append(out, x)
		}
	}
	return out
}

// vendorFiles orders the vendor's claims its own way, assigns its own
// ICNs, trace numbers, and control numbers, and groups them into payments
// and files.
func vendorFiles(cfg config, w world, prof profile, v, d int, claims []Claim) []File {
	date := strings.ReplaceAll(cfg.Submissions[d], "-", "")
	month := func(c Claim) string { return c.Received[:6] }
	sort.SliceStable(claims, func(i, j int) bool {
		a, b := claims[i], claims[j]
		if a.Workload != b.Workload {
			return a.Workload < b.Workload
		}
		if a.Payee != b.Payee {
			return a.Payee < b.Payee
		}
		if prof.ByMonth {
			if month(a) != month(b) {
				return month(a) < month(b)
			}
			if a.dateOfService() != b.dateOfService() {
				return a.dateOfService() < b.dateOfService()
			}
		} else if a.PatientCtl != b.PatientCtl {
			return a.PatientCtl < b.PatientCtl
		}
		return a.Received < b.Received
	})

	julian := parseYMD(date).YearDay()
	counters := make([]int, len(w.Workloads))
	for i, c := range claims {
		if c.Reversal {
			continue
		}
		counters[c.Workload]++
		if prof.ICN == "CA" {
			claims[i].ICN = fmt.Sprintf("CA26%03d%06d", julian, counters[c.Workload]*10+c.Workload)
		} else {
			claims[i].ICN = fmt.Sprintf("B26%03d%07d", julian, c.Workload*1000000+counters[c.Workload])
		}
		claims[i].OrigICN = ""
	}
	claims = linkChains(claims)

	var payments []Payment
	for _, c := range claims {
		n := len(payments)
		if n == 0 || payments[n-1].Workload != c.Workload || payments[n-1].Payee != c.Payee ||
			prof.ByMonth && month(payments[n-1].Claims[0]) != month(c) {
			payments = append(payments, Payment{Workload: c.Workload, Payee: c.Payee, Date: date})
			n++
		}
		payments[n-1].Claims = append(payments[n-1].Claims, c)
		if c.Interest > 0 {
			payments[n-1].PLBs = append(payments[n-1].PLBs, PLB{Code: "L6", Ref: c.ICN, Amount: -c.Interest, ClaimID: c.ID})
		}
	}

	// Provider-level recoveries and balances stay with the payee: each goes
	// to that payee's largest vendor payment (the earlier one on a tie).
	for _, mp := range w.Payments {
		for _, a := range mp.PLBs {
			if a.ClaimID >= 0 || strings.HasPrefix(a.Ref, "BAL") {
				continue
			}
			best := -1
			for i, p := range payments {
				if p.Workload == mp.Workload && p.Payee == mp.Payee && (best == -1 || p.total() > payments[best].total()) {
					best = i
				}
			}
			if best >= 0 {
				payments[best].PLBs = append(payments[best].PLBs, a)
			}
		}
	}

	for i := range payments {
		if prof.Trace == "CA" {
			payments[i].Trace = fmt.Sprintf("CA%s%05d", date[2:], i+1)
		} else {
			payments[i].Trace = fmt.Sprintf("88%02d%06d", d+1, i+1)
		}
		payments[i] = floorAtZero(payments[i])
	}

	// Vendor 1 puts up to three payees in a file; vendor 2 puts a whole workload in one.
	var files []File
	payeesInFile := 0
	for i, p := range payments {
		newPayee := i == 0 || payments[i-1].Payee != p.Payee || payments[i-1].Workload != p.Workload
		n := len(files)
		if n == 0 || files[n-1].Workload != p.Workload || !prof.ByMonth && newPayee && payeesInFile == 3 {
			files = append(files, File{
				Workload: p.Workload,
				ISA13:    prof.ISABase + (d+1)*10000 + n + 1,
				GS06:     (d+1)*1000 + n + 1,
				Date:     addDays(date, -1),
				Time:     prof.Time,
				STBase:   prof.STBase,
			})
			payeesInFile = 0
		}
		if newPayee {
			payeesInFile++
		}
		files[len(files)-1].Payments = append(files[len(files)-1].Payments, p)
	}
	perWorkload := make([]int, len(w.Workloads))
	for i, f := range files {
		perWorkload[f.Workload]++
		files[i].Path = fmt.Sprintf("%s/%s/%s-%02d.835", cfg.Vendors[v], cfg.Submissions[d], w.Workloads[f.Workload].Number, perWorkload[f.Workload])
	}
	return files
}
