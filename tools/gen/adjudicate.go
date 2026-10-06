package main

// adjudicate prices a claim from its submitted lines and its decisions and
// returns the claim with every result filled in. The reference and each
// vendor call it with their own decisions, so every side balances:
// SVC02 - SVC03 equals the line's adjustments, and CLP03 - CLP04 equals
// every adjustment on the claim and its lines.
func adjudicate(c Claim) Claim {
	c = c.clone()
	c.Adjs = nil
	c.Remarks = nil
	c.Charge, c.Paid, c.PatientResp, c.DRGAmount, c.OutlierAmount = 0, 0, 0, 0, 0
	c.CrossoverShown = false

	if c.inpatient() {
		c = priceInpatient(c)
	} else {
		c = priceLines(c)
	}
	if c.Reversal {
		c = negate(c)
		c.Status = "22"
		return c
	}
	c.Remarks = claimRemarks(c)
	c.CrossoverShown = c.Crossover && !c.DropCrossover && c.PatientResp > 0 && c.Paid > 0
	// A claim whose balance goes on to a supplemental payer is reported as
	// processed as primary and forwarded (CLP02 19).
	if c.CrossoverShown && c.Status == "1" {
		c.Status = "19"
	}
	// Medicare secondary: a prior payer paid part of the allowed amount, so
	// the claim-level OA*23 carries what that payer paid and Medicare pays
	// the rest.
	if c.Secondary && c.Paid > 0 && !c.inpatient() {
		prior := percentOf(c.Paid, 60)
		c.Adjs = append(c.Adjs, Adj{"OA", "23", prior})
		c.Paid -= prior
		c.Status = "2"
	}
	c = applyRemarkLevel(c)
	if c.ZeroCAS {
		for i, ln := range c.Lines {
			if ln.Paid > 0 {
				c.Lines[i].Adjs = append(c.Lines[i].Adjs, Adj{"CO", "45", 0})
				break
			}
		}
	}
	return c
}

func priceInpatient(c Claim) Claim {
	charge := c.InpatientCharge
	c.Charge = charge
	if c.Denial.Reason != "" {
		c.Adjs = []Adj{{c.Denial.Group, c.Denial.Reason, charge}}
		c.Status = "4"
		if c.Denial.Group == "PR" {
			c.PatientResp = charge
		}
		return c
	}
	allowed := drgAmount(c.BaseRate, c.Weight)
	// A cost outlier adds to the DRG amount; it is reported in MIA02.
	if c.Outlier && !c.OutlierDropped {
		c.OutlierAmount = percentOf(allowed, 30)
	}
	// The 2025 Part A inpatient deductible, $1,676: the year of the dates of
	// service (REALISM.md section 2.4).
	var ded Money
	if c.PartADeductible {
		ded = min(167600, allowed)
	}
	seq := percentOf(allowed+c.OutlierAmount-ded, 2)
	c.DRGAmount = allowed
	c.Paid = allowed + c.OutlierAmount - ded - seq
	c.PatientResp = ded
	c.Adjs = nonzero([]Adj{{"CO", "45", charge - allowed - c.OutlierAmount}, {"CO", "253", seq}, {"PR", "1", ded}})
	c.Status = "1"
	return c
}

func priceLines(c Claim) Claim {
	remaining := c.Deductible
	before := make([]Money, len(c.Lines))
	writeOff := make([]Money, len(c.Lines))
	reduction := make([]Money, len(c.Lines))
	deductible := make([]Money, len(c.Lines))
	coinsurance := make([]Money, len(c.Lines))
	denied := 0
	for i := range c.Lines {
		ln := c.Lines[i]
		ln.Adjs, ln.Remarks = nil, nil
		ln.Allowed, ln.Paid = 0, 0
		ln.Charge = ln.UnitCharge * Money(ln.Units)
		switch {
		case ln.Denial.Reason != "":
			denied++
			ln.Adjs = []Adj{{ln.Denial.Group, ln.Denial.Reason, ln.Charge}}
			if ln.Denial.Remark != "" {
				ln.Remarks = []string{ln.Denial.Remark}
			}
		case ln.Packaged:
			ln.Adjs = []Adj{{"CO", "97", ln.Charge}}
		default:
			unitFee := ln.UnitFee
			if c.FeeScale != 0 {
				unitFee = percentOf(unitFee, 100+c.FeeScale)
			}
			fee := unitFee * Money(ln.Units)
			allowed := fee
			if ln.Reduced {
				reduction[i] = fee - percentOf(fee, 50)
				allowed = fee - reduction[i]
			}
			writeOff[i] = ln.Charge - fee
			if !ln.Exempt {
				deductible[i] = min(allowed, remaining)
				remaining -= deductible[i]
				coinsurance[i] = percentOf(allowed-deductible[i], 20)
			}
			before[i] = allowed - deductible[i] - coinsurance[i]
			ln.Allowed = allowed
		}
		c.Lines[i] = ln
	}

	seq := sequestration(before, c.SeqPerLine, c.SeqTruncate)
	for i := range c.Lines {
		ln := c.Lines[i]
		if ln.Denial.Reason == "" && !ln.Packaged {
			ln.Adjs = nonzero([]Adj{
				{"CO", "45", writeOff[i]},
				{"CO", "59", reduction[i]},
				{"CO", "253", seq[i]},
				{"PR", "1", deductible[i]},
				{"PR", "2", coinsurance[i]},
			})
			ln.Paid = before[i] - seq[i]
		}
		c.Charge += ln.Charge
		c.Paid += ln.Paid
		for _, a := range ln.Adjs {
			if a.Group == "PR" {
				c.PatientResp += a.Amount
			}
		}
		c.Lines[i] = ln
	}
	c.Status = "1"
	if denied == len(c.Lines) {
		c.Status = "4"
	}
	if c.NotOurs {
		c.Status = "23"
	}
	return c
}

// sequestration returns the 2 percent reduction for each line. The
// reference rounds once for the claim and spreads the cents across lines;
// perLine rounds each line on its own, and truncate drops the fraction of
// a cent instead of rounding it (pattern 2, both variants).
func sequestration(before []Money, perLine, truncate bool) []Money {
	seq := make([]Money, len(before))
	if perLine {
		for i, b := range before {
			seq[i] = percentOf(b, 2)
		}
		return seq
	}
	var total Money
	for _, b := range before {
		total += b
	}
	target := percentOf(total, 2)
	if truncate {
		target = Money(int64(total) * 2 / 100)
	}
	// Give each line the whole cents of its share, then hand the remaining
	// cents, one each, to the lines with the largest leftover fractions
	// (the earlier line wins a tie).
	leftover := make([]int64, len(before))
	var given Money
	for i, b := range before {
		x := int64(b) * 2
		seq[i] = Money(x / 100)
		leftover[i] = x % 100
		if b == 0 {
			leftover[i] = -1
		}
		given += seq[i]
	}
	for given < target {
		best := -1
		for i := range leftover {
			if leftover[i] >= 0 && (best == -1 || leftover[i] > leftover[best]) {
				best = i
			}
		}
		if best == -1 {
			break
		}
		seq[best]++
		leftover[best] = -1
		given++
	}
	return seq
}

func nonzero(adjs []Adj) []Adj {
	var out []Adj
	for _, a := range adjs {
		if a.Amount != 0 {
			out = append(out, a)
		}
	}
	return out
}

// negate turns an adjudicated claim into its reversal: every amount changes sign.
func negate(c Claim) Claim {
	c.Charge, c.Paid, c.PatientResp, c.DRGAmount, c.OutlierAmount = -c.Charge, -c.Paid, -c.PatientResp, -c.DRGAmount, -c.OutlierAmount
	for i, a := range c.Adjs {
		c.Adjs[i].Amount = -a.Amount
	}
	for i := range c.Lines {
		ln := c.Lines[i]
		ln.Charge, ln.Allowed, ln.Paid = -ln.Charge, -ln.Allowed, -ln.Paid
		for j, a := range ln.Adjs {
			ln.Adjs[j].Amount = -a.Amount
		}
		c.Lines[i] = ln
	}
	c.Interest = 0
	return c
}

// claimRemarks returns the claim-level remark codes (MOA, or MIA for inpatient).
func claimRemarks(c Claim) []string {
	if c.Unprocessable {
		return []string{"MA130"}
	}
	denied := c.Denial.Reason != ""
	for _, ln := range c.Lines {
		if ln.Denial.Reason != "" {
			denied = true
		}
	}
	var r []string
	if (denied || c.PatientResp != 0) && c.RemarkVariant != 2 {
		r = append(r, "MA01")
	}
	if c.Crossover && !c.DropCrossover && c.PatientResp > 0 && c.Paid > 0 {
		r = append(r, "MA18")
	}
	if c.ExtraAlert {
		r = append(r, "N130")
	}
	if c.RemarkVariant == 1 && len(r) >= 2 {
		r[0], r[1] = r[1], r[0]
	}
	return r
}

// applyRemarkLevel moves one remark between the claim and a line: level 1
// lifts the first line remark up to the claim, level 2 pushes the claim's
// first remark down to the first line (pattern 10, T3.34). The set of
// remarks is the same; only where each one sits changes.
func applyRemarkLevel(c Claim) Claim {
	switch c.RemarkLevel {
	case 1:
		for i, ln := range c.Lines {
			if len(ln.Remarks) > 0 {
				c.Remarks = append(c.Remarks, ln.Remarks[0])
				c.Lines[i].Remarks = ln.Remarks[1:]
				break
			}
		}
	case 2:
		if len(c.Remarks) > 0 && len(c.Lines) > 0 && !c.inpatient() {
			c.Lines[0].Remarks = append([]string{c.Remarks[0]}, c.Lines[0].Remarks...)
			c.Remarks = c.Remarks[1:]
		}
	}
	return c
}
