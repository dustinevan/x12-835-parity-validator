package main

import (
	"fmt"
	"strings"
)

// PLB is one provider-level adjustment. A positive amount reduces the
// payment; interest (L6) is negative because it adds to the payment.
type PLB struct {
	Code    string
	Ref     string
	Amount  Money
	ClaimID int // the claim an interest item belongs to; -1 for provider-level items
}

// Payment is one ST..SE transaction: one payee, one payment date.
type Payment struct {
	Workload int
	Payee    int
	Date     string
	Trace    string
	Claims   []Claim
	PLBs     []PLB
}

// total is BPR02: the claim payments less the provider-level adjustments.
func (p Payment) total() Money {
	var t Money
	for _, c := range p.Claims {
		t += c.Paid
	}
	for _, a := range p.PLBs {
		t -= a.Amount
	}
	return t
}

// floorAtZero carries a negative balance forward (PLB FB) so the payment
// is never negative.
func floorAtZero(p Payment) Payment {
	if t := p.total(); t < 0 {
		p.PLBs = append(p.PLBs, PLB{Code: "FB", Ref: "BAL" + p.Date, Amount: t, ClaimID: -1})
	}
	return p
}

// File is one interchange: one workload, several payments.
type File struct {
	Path     string
	Workload int
	ISA13    int
	GS06     int
	Date     string
	Time     string
	STBase   int
	Payments []Payment
}

// macPayments groups the MAC claims into one payment per payee per cycle
// and draws the provider-level adjustments.
func (g generator) macPayments(w world) []Payment {
	var out []Payment
	for _, c := range w.Claims {
		n := len(out)
		if n == 0 || out[n-1].Workload != c.Workload || out[n-1].Payee != c.Payee || out[n-1].Date != w.Cycles[c.Cycle] {
			out = append(out, Payment{Workload: c.Workload, Payee: c.Payee, Date: w.Cycles[c.Cycle]})
			n++
		}
		out[n-1].Claims = append(out[n-1].Claims, c)
		if c.Interest > 0 {
			out[n-1].PLBs = append(out[n-1].PLBs, PLB{Code: "L6", Ref: c.ICN, Amount: -c.Interest, ClaimID: c.ID})
		}
	}
	for i := range out {
		p := out[i]
		p.Trace = fmt.Sprintf("%010d", 918400+i)
		// About 8 percent of payments carry a recovery or a forwarding balance.
		if t := p.total(); t > 1000 && g.chance(0.08) {
			code := "WO"
			if g.chance(0.3) {
				code = "FB"
			}
			ref := fmt.Sprintf("%s25%03d9%05d", w.Workloads[p.Workload].Region, g.between(1, 365), g.intn(100000))
			amount := Money(float64(t) * (0.02 + 0.28*g.float()))
			p.PLBs = append(p.PLBs, PLB{Code: code, Ref: ref, Amount: amount, ClaimID: -1})
		}
		out[i] = floorAtZero(p)
	}
	return out
}

// macFiles puts each workload's payments for one cycle into one file.
func macFiles(w world) []File {
	var out []File
	for _, p := range w.Payments {
		n := len(out)
		if n == 0 || out[n-1].Workload != p.Workload || out[n-1].Payments[0].Date != p.Date {
			wl := w.Workloads[p.Workload]
			out = append(out, File{
				Path:     fmt.Sprintf("mac/%s-%s.835", wl.Number, p.Date),
				Workload: p.Workload,
				ISA13:    100000001 + n,
				GS06:     1 + n,
				Date:     addDays(p.Date, -2),
				Time:     "0600",
				STBase:   1,
			})
			n++
		}
		out[n-1].Payments = append(out[n-1].Payments, p)
	}
	return out
}

// segment joins elements with * and drops trailing empty elements.
func segment(elems ...string) string {
	n := len(elems)
	for n > 1 && elems[n-1] == "" {
		n--
	}
	return strings.Join(elems[:n], "*")
}

const receiverID = "SYNCLEARING01"

// render writes one file as X12 835 bytes, one segment per line.
func render(ws []Workload, f File) []byte {
	w := ws[f.Workload]
	sender := "SYNMAC" + w.Number
	var b strings.Builder
	put := func(s string) {
		b.WriteString(s)
		b.WriteString("~\n")
	}
	put(fmt.Sprintf("ISA*00*          *00*          *ZZ*%-15s*ZZ*%-15s*%s*%s*^*00501*%09d*0*T*:", sender, receiverID, f.Date[2:], f.Time, f.ISA13))
	put(segment("GS", "HP", sender, receiverID, f.Date, f.Time, fmt.Sprint(f.GS06), "X", "005010X221A1"))
	for i, p := range f.Payments {
		st := fmt.Sprintf("%04d", f.STBase+i)
		segs := transaction(w, p, st)
		for _, s := range segs {
			put(s)
		}
		put(segment("SE", fmt.Sprint(len(segs)+1), st))
	}
	put(segment("GE", fmt.Sprint(len(f.Payments)), fmt.Sprint(f.GS06)))
	put(fmt.Sprintf("IEA*1*%09d", f.ISA13))
	return []byte(b.String())
}

// transaction returns the segments from ST through the last PLB.
func transaction(w Workload, p Payment, st string) []string {
	payee := w.Payees[p.Payee]
	total := p.total()
	segs := []string{segment("ST", "835", st)}
	if total > 0 {
		segs = append(segs, segment("BPR", "I", total.String(), "C", "ACH", "CCP", "01", w.DFI, "DA", w.Account, "1"+w.TaxID, "", "01", payee.DFI, "DA", payee.Account, p.Date))
	} else {
		segs = append(segs, segment("BPR", "H", "0.00", "C", "NON", "", "", "", "", "", "", "", "", "", "", "", p.Date))
	}
	segs = append(segs,
		segment("TRN", "1", p.Trace, "1"+w.TaxID),
		segment("DTM", "405", p.Date),
		segment("N1", "PR", w.Name, "PI", w.Number),
		segment("N3", "1 SYNTHETIC PAYER PLAZA"),
		segment("N4", "EXAMPLE CITY", w.State, w.Zip),
		segment("PER", "BL", "SYNTHETIC EDI HELP DESK", "TE", "8005550100"),
		segment("N1", "PE", payee.Name, "XX", payee.NPI),
		segment("N3", payee.Street),
		segment("N4", "EXAMPLE CITY", w.State, payee.Zip),
		segment("REF", "TJ", payee.TaxID),
		segment("LX", "1"),
	)
	for _, c := range p.Claims {
		segs = append(segs, claimSegments(c)...)
	}
	for _, a := range p.PLBs {
		segs = append(segs, segment("PLB", payee.NPI, p.Date[:4]+"1231", a.Code+":"+a.Ref, a.Amount.String()))
	}
	return segs
}

func weightString(w int64) string {
	return fmt.Sprintf("%d.%04d", w/10000, w%10000)
}

func filingIndicator(c Claim) string {
	if c.Type == "institutional" {
		return "MA"
	}
	return "MB"
}

func claimSegments(c Claim) []string {
	// CLP09, the claim frequency code, is sent on institutional claims
	// only (the CMS 835 flat file marks it intermediary-only); CLP11 and
	// CLP12, the DRG and its weight, on inpatient claims.
	clp := []string{"CLP", c.PatientCtl, c.Status, c.Charge.String(), c.Paid.String(), c.PatientResp.String(), filingIndicator(c), c.ICN, c.Facility, c.Frequency}
	if c.inpatient() && c.DRG != "" {
		clp = append(clp, "", c.DRG, weightString(c.Weight))
	}
	segs := []string{segment(clp...)}
	segs = append(segs, casSegments(c.Adjs)...)
	segs = append(segs, segment("NM1", "QC", "1", "SYNTHETIC", c.Patient.First, c.Patient.Middle, "", "", "MI", c.Patient.MBI))
	if c.Rendering.NPI != "" && !c.OmitRendering {
		segs = append(segs, segment("NM1", "82", "1", "SYNTHETIC", c.Rendering.Name, "", "", "", "XX", c.Rendering.NPI))
	}
	if c.CrossoverShown {
		segs = append(segs, segment("NM1", "TT", "2", c.Patient.Insurer, "", "", "", "", "PI", c.Patient.InsurerID))
	}
	if c.inpatient() && !c.Reversal && c.DRGAmount != 0 {
		mia := make([]string, 25)
		mia[0] = "MIA"
		mia[1] = fmt.Sprint(c.Days)
		if c.OutlierAmount != 0 {
			mia[2] = c.OutlierAmount.String()
		}
		mia[4] = c.DRGAmount.String()
		for i, r := range c.Remarks {
			switch {
			case i == 0:
				mia[5] = r
			case i <= 4:
				mia[19+i] = r
			}
		}
		segs = append(segs, segment(mia...))
	} else if len(c.Remarks) > 0 {
		moa := []string{"MOA", "", ""}
		moa = append(moa, c.Remarks...)
		segs = append(segs, segment(moa...))
	}
	if c.OrigICN != "" {
		segs = append(segs, segment("REF", "F8", c.OrigICN))
	}
	if c.Type == "institutional" {
		segs = append(segs, segment("DTM", "232", c.From), segment("DTM", "233", c.Thru))
	}
	segs = append(segs, segment("DTM", "050", c.Received))
	if c.Interest > 0 && !c.OmitInterest {
		segs = append(segs, segment("AMT", "I", c.Interest.String()))
	}
	if c.inpatient() {
		segs = append(segs, segment("QTY", "CA", fmt.Sprint(c.Days)))
	}
	for _, ln := range c.Lines {
		segs = append(segs, lineSegments(c, ln)...)
	}
	return segs
}

func procedure(ln Line) string {
	return strings.Join(append([]string{ln.Qual, ln.Code}, ln.Mods...), ":")
}

func unitsString(ln Line) string {
	if ln.OmitUnits {
		return ""
	}
	return fmt.Sprint(ln.Units)
}

func lineSegments(c Claim, ln Line) []string {
	rev := ""
	if c.Type == "institutional" {
		rev = ln.Rev
	}
	segs := []string{
		segment("SVC", procedure(ln), ln.Charge.String(), ln.Paid.String(), rev, unitsString(ln)),
		segment("DTM", "472", ln.Date),
	}
	segs = append(segs, casSegments(ln.Adjs)...)
	if !ln.OmitID {
		segs = append(segs, segment("REF", "6R", ln.ID))
	}
	if ln.Allowed != 0 && !ln.OmitAllowed {
		segs = append(segs, segment("AMT", "B6", ln.Allowed.String()))
	}
	for _, r := range ln.Remarks {
		segs = append(segs, segment("LQ", "HE", r))
	}
	return segs
}

// casSegments writes adjustments as CAS segments: consecutive adjustments
// in one group share a segment, up to six per segment.
func casSegments(adjs []Adj) []string {
	var out []string
	for i := 0; i < len(adjs); {
		s := "CAS*" + adjs[i].Group
		n := 0
		j := i
		for j < len(adjs) && adjs[j].Group == adjs[i].Group && n < 6 {
			if n > 0 {
				s += "*"
			}
			s += "*" + adjs[j].Reason + "*" + adjs[j].Amount.String()
			n++
			j++
		}
		out = append(out, s)
		i = j
	}
	return out
}
