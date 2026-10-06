package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"time"
)

// Money is an amount in cents. Every amount in a remittance is an exact
// decimal with two places, so the generator never computes money in floats
// after a value has been drawn.
type Money int64

func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	return fmt.Sprintf("%s%d.%02d", sign, int64(m)/100, int64(m)%100)
}

// percentOf returns p percent of m, rounded half away from zero to the cent.
func percentOf(m Money, p int64) Money {
	x := int64(m) * p
	if x >= 0 {
		return Money((x + 50) / 100)
	}
	return -Money((-x + 50) / 100)
}

// Adj is one claim adjustment: a group code, a reason code, and an amount.
type Adj struct {
	Group  string
	Reason string
	Amount Money
}

// Denial says why a line (or an inpatient claim) was denied and which
// remark explains it.
type Denial struct {
	Group  string
	Reason string
	Remark string
}

// Line is one service line: what the provider submitted, the pricing
// decisions for it, and the adjudicated result.
type Line struct {
	Qual       string
	Code       string
	Mods       []string
	Rev        string
	Class      string
	Units      int
	Date       string
	UnitCharge Money
	UnitFee    Money
	UrbanFee   Money // ambulance mileage: the urban per-mile fee; UnitFee is the rural fee on rural trips
	ID         string

	Denial   Denial
	Packaged bool
	Reduced  bool // multiple-procedure reduction: allowed at 50 percent
	Exempt   bool // no deductible or coinsurance (clinical laboratory fee schedule)

	OmitID      bool
	OmitAllowed bool
	OmitUnits   bool

	Charge  Money
	Allowed Money
	Paid    Money
	Adjs    []Adj
	Remarks []string
}

// Renderer is a rendering practitioner inside a billing group.
type Renderer struct {
	Name string
	NPI  string
}

// Patient is a synthetic Medicare beneficiary.
type Patient struct {
	First        string
	Middle       string
	MBI          string
	Supplemental bool
	Insurer      string
	InsurerID    string
}

// Claim is one claim as submitted, the decisions made about it, and the
// adjudicated result. The MAC side and each vendor side hold their own copy.
type Claim struct {
	ID       int // creation index
	DrawID   int // the claim whose vendor draws apply: the original for a reversal, the first of a key pair
	Workload int
	Payee    int
	Juris    string
	Part     string // A, B, or DME: the workload's part
	Type     string // professional, dme, or institutional
	BillType string // institutional bill type, for example 131
	Facility string // CLP08: place of service or the first two digits of the bill type

	PatientCtl string
	Patient    Patient
	Rendering  Renderer
	Received   string
	From       string
	Thru       string
	Month      int
	Rural      bool
	Lines      []Line

	Deductible      Money // Part B deductible still unmet when the claim was processed
	AltDeductible   Money // deductible a vendor applies when it credits this claim instead of another
	PartADeductible bool
	DRG             string
	Weight          int64 // DRG relative weight in ten-thousandths
	BaseRate        Money
	Days            int
	InpatientCharge Money
	Denial          Denial // inpatient whole-claim denial
	Unprocessable   bool
	Crossover       bool
	Interest        Money
	Pick            int
	Cycle           int
	Reversal        bool
	Correction      bool
	NetAdjust       bool // a vendor's correction sent as a fresh claim: no REF*F8, frequency 1 (pattern 9)
	OrigID          int
	Pair            string // resub-first, resub-second, distinct-first, distinct-second

	SeqPerLine    bool
	SeqTruncate   bool  // sequestration truncated to the cent instead of rounded
	FeeScale      int64 // percent added to every fee (a different locality or an older schedule); 0 is none
	DropCrossover bool
	OmitInterest  bool
	OmitRendering bool
	ExtraAlert    bool
	RemarkVariant int  // 1 swaps the first two claim remarks, 2 drops MA01
	RemarkLevel   int  // 1 moves the first line remark up to the claim, 2 moves the claim's first remark down to the first line
	ZeroCAS       bool // a zero-amount CAS*CO*45 on the first paid line

	// What the MAC side decided beyond paid or denied.
	Frequency      string // CLP09 on institutional claims: 1, or 7 and 8 on a chain
	Secondary      bool   // Medicare secondary payer: a prior payer paid part (CLP02 2, CAS*OA*23)
	NotOurs        bool   // sent to the wrong contractor (CLP02 23, every line OA*109)
	Outlier        bool   // inpatient cost outlier paid (MIA02)
	OutlierDropped bool   // a vendor that does not pay the outlier
	OutlierAmount  Money

	Status         string
	Charge         Money
	Paid           Money
	PatientResp    Money
	DRGAmount      Money
	Adjs           []Adj
	Remarks        []string
	CrossoverShown bool
	ICN            string
	OrigICN        string
}

// clone returns a copy that shares no slices with c.
func (c Claim) clone() Claim {
	lines := make([]Line, len(c.Lines))
	for i, ln := range c.Lines {
		ln.Mods = append([]string(nil), ln.Mods...)
		ln.Adjs = append([]Adj(nil), ln.Adjs...)
		ln.Remarks = append([]string(nil), ln.Remarks...)
		lines[i] = ln
	}
	c.Lines = lines
	c.Adjs = append([]Adj(nil), c.Adjs...)
	c.Remarks = append([]string(nil), c.Remarks...)
	return c
}

func (c Claim) inpatient() bool {
	return c.BillType == "111"
}

// dateOfService is the key date: the statement start for institutional
// claims, otherwise the earliest line date.
func (c Claim) dateOfService() string {
	if c.Type == "institutional" {
		return c.From
	}
	d := ""
	for _, ln := range c.Lines {
		if d == "" || ln.Date < d {
			d = ln.Date
		}
	}
	return d
}

// Payee is a billing provider. Prices is its charge list, fixed when the
// payee is created.
type Payee struct {
	Name      string
	Kind      string
	NPI       string
	TaxID     string
	Street    string
	Zip       string
	DFI       string
	Account   string
	Weight    float64
	Habit     string // how it numbers CLP01: unique, patient, or episode
	Code      string
	Renderers []Renderer
	Prices    map[string]price
	Markup    float64 // hospital charge markup over cost
	IPRatio   float64 // inpatient payment ÷ charge
}

type price struct {
	Charge Money
	Fee    Money
	Urban  Money
}

// Workload is one MAC contract workload: a jurisdiction and a part.
type Workload struct {
	Number   string
	Name     string
	Juris    string
	Part     string
	Locality int64 // percent applied to national fees
	TaxID    string
	DFI      string
	Account  string
	Region   string
	State    string
	Zip      string
	Payees   []Payee
}

// service is one billable code with its synthetic national fee (professional
// and DME) or typical cost (institutional).
type service struct {
	Qual     string
	Code     string
	Mod      string
	Rev      string
	Class    string
	Cents    Money
	Weight   float64
	MaxUnits int
}

// Synthetic fee schedule. Values are plausible magnitudes, not published fees.
var clinicFirst = []service{
	{"HC", "99213", "", "", "em", 9200, 0.40, 1},
	{"HC", "99214", "", "", "em", 12844, 0.40, 1},
	{"HC", "99215", "", "", "em", 18110, 0.07, 1},
	{"HC", "99203", "", "", "em", 11020, 0.05, 1},
	{"HC", "99204", "", "", "em", 16800, 0.08, 1},
}

var clinicMore = []service{
	{"HC", "93000", "", "", "diag", 1652, 0.25, 1},
	{"HC", "96372", "", "", "diag", 1480, 0.15, 1},
	{"HC", "36415", "", "", "lab", 300, 0.20, 1},
	{"HC", "81002", "", "", "lab", 350, 0.10, 1},
	{"HC", "20610", "", "", "surg", 7190, 0.12, 1},
	{"HC", "11102", "", "", "surg", 9520, 0.08, 1},
	{"HC", "17000", "", "", "surg", 6200, 0.10, 1},
}

var surgeryServices = []service{
	{"HC", "29881", "", "", "surg", 52000, 0.15, 1},
	{"HC", "64721", "", "", "surg", 48000, 0.12, 1},
	{"HC", "26055", "", "", "surg", 39000, 0.12, 1},
	{"HC", "66984", "", "", "surg", 54000, 0.12, 1},
	{"HC", "64483", "", "", "surg", 21000, 0.18, 1},
	{"HC", "11042", "", "", "surg", 11800, 0.13, 1},
	{"HC", "10060", "", "", "surg", 12500, 0.08, 1},
	{"HC", "12001", "", "", "surg", 10500, 0.10, 1},
}

var labServices = []service{
	{"HC", "80053", "", "", "lab", 1056, 0.22, 1},
	{"HC", "85025", "", "", "lab", 777, 0.22, 1},
	{"HC", "80061", "", "", "lab", 1330, 0.12, 1},
	{"HC", "83036", "", "", "lab", 953, 0.10, 1},
	{"HC", "84443", "", "", "lab", 1664, 0.10, 1},
	{"HC", "82947", "", "", "lab", 400, 0.06, 1},
	{"HC", "87086", "", "", "lab", 802, 0.06, 1},
	{"HC", "36415", "", "", "lab", 300, 0.08, 1},
	{"HC", "P9603", "", "", "lab", 110, 0.04, 20},
}

var ambulanceBase = []service{
	{"HC", "A0427", "", "", "amb", 52000, 0.30, 1},
	{"HC", "A0429", "", "", "amb", 44000, 0.40, 1},
	{"HC", "A0428", "", "", "amb", 26500, 0.30, 1},
}

// Ground mileage: urban and rural per-mile fees (synthetic).
var mileage = service{"HC", "A0425", "", "", "mile", 1210, 1, 60}

const urbanMileFee Money = 805

var therapyServices = []service{
	{"HC", "97110", "", "", "therapy", 2980, 0.35, 4},
	{"HC", "97530", "", "", "therapy", 3510, 0.25, 4},
	{"HC", "97140", "", "", "therapy", 2690, 0.15, 3},
	{"HC", "97112", "", "", "therapy", 3350, 0.12, 3},
	{"HC", "97116", "", "", "therapy", 2890, 0.08, 3},
	{"HC", "97161", "", "", "therapy", 10300, 0.05, 1},
}

var oxygenServices = []service{
	{"HC", "E1390", "RR", "", "oxygen", 16850, 0.55, 1},
	{"HC", "E0431", "RR", "", "oxygen", 2260, 0.25, 1},
	{"HC", "E1392", "RR", "", "oxygen", 5163, 0.10, 1},
	{"HC", "K0738", "RR", "", "oxygen", 5163, 0.10, 1},
}

var cpapServices = []service{
	{"HC", "E0601", "RR", "", "rental", 8810, 0.20, 1},
	{"HC", "A7038", "NU", "", "supply", 410, 0.20, 2},
	{"HC", "A7035", "NU", "", "supply", 2700, 0.12, 1},
	{"HC", "A7034", "NU", "", "supply", 8800, 0.16, 1},
	{"HC", "A7031", "NU", "", "supply", 3600, 0.16, 1},
	{"HC", "A4604", "NU", "", "supply", 2600, 0.16, 1},
}

var diabeticServices = []service{
	{"HC", "A4239", "NU", "", "supply", 21000, 0.45, 1},
	{"HC", "A4253", "NU", "", "supply", 900, 0.40, 3},
	{"HC", "A4259", "NU", "", "supply", 600, 0.15, 2},
}

var generalDMEServices = []service{
	{"HC", "E0260", "RR", "", "rental", 11035, 0.30, 1},
	{"HC", "E0601", "RR", "", "rental", 8810, 0.20, 1},
	{"HC", "A4253", "NU", "", "supply", 900, 0.25, 3},
	{"HC", "A7034", "NU", "", "supply", 8800, 0.25, 1},
}

// Institutional lines carry a typical cost; the hospital's markup sets the
// charge and its payment ratio sets the fee.
var outpatientFirst = []service{
	{"HC", "99284", "", "0450", "visit", 30000, 0.35, 1},
	{"HC", "99283", "", "0450", "visit", 20000, 0.25, 1},
	{"HC", "71046", "", "0320", "visit", 9000, 0.12, 1},
	{"HC", "74177", "", "0350", "visit", 45000, 0.10, 1},
	{"HC", "29881", "", "0360", "visit", 380000, 0.05, 1},
	{"HC", "96374", "", "0260", "visit", 25000, 0.08, 1},
	{"HC", "G0378", "", "0762", "visit", 60000, 0.05, 1},
}

var outpatientMore = []service{
	{"HC", "80053", "", "0300", "lab", 6000, 0.20, 1},
	{"HC", "85025", "", "0300", "lab", 4000, 0.20, 1},
	{"HC", "36415", "", "0300", "lab", 1500, 0.15, 1},
	{"HC", "93005", "", "0730", "pack", 8000, 0.15, 1},
	{"NU", "0250", "", "", "pack", 5000, 0.15, 1},
	{"NU", "0270", "", "", "pack", 4000, 0.10, 1},
	{"HC", "71046", "", "0320", "visit", 9000, 0.05, 1},
}

var rhcServices = []service{
	{"NU", "0521", "", "", "visit", 22000, 1, 1},
	{"HC", "36415", "", "0300", "pack", 1500, 0.5, 1},
}

var esrdServices = []service{
	{"NU", "0821", "", "", "visit", 40000, 1, 14},
	{"HC", "85025", "", "0300", "lab", 4000, 0.5, 1},
}

var orfServices = []service{
	{"HC", "97110", "", "0420", "visit", 9000, 0.4, 4},
	{"HC", "97530", "", "0420", "visit", 10000, 0.3, 4},
	{"HC", "97140", "", "0420", "visit", 8000, 0.3, 3},
}

var snfServices = []service{
	{"NU", "0022", "", "", "pps", 1500000, 1, 1},
	{"NU", "0120", "", "", "zero", 2500000, 1, 1},
	{"NU", "0250", "", "", "zero", 300000, 1, 1},
	{"NU", "0420", "", "", "zero", 400000, 1, 1},
}

var hhaServices = []service{
	{"NU", "0023", "", "", "pps", 220000, 1, 1},
	{"NU", "0551", "", "", "zero", 150000, 1, 1},
	{"NU", "0421", "", "", "zero", 120000, 1, 1},
	{"NU", "0431", "", "", "zero", 90000, 1, 1},
}

var hospiceServices = []service{
	{"NU", "0651", "", "", "pps", 22000, 1, 31},
	{"NU", "0551", "", "", "zero", 20000, 1, 1},
}

// DRG table with synthetic relative weights. The first twelve are the
// severity-boundary DRGs where a grouper version change moves the claim to
// the paired DRG (pattern 7, variantGrouperVersion in vendor.go).
type drg struct {
	Code      string
	Weight    int64
	AltCode   string
	AltWeight int64
}

var drgs = []drg{
	{"0871", 19500, "0872", 10500},
	{"0872", 10500, "0871", 19500},
	{"0291", 13500, "0292", 9000},
	{"0292", 9000, "0291", 13500},
	{"0190", 11000, "0191", 8900},
	{"0194", 9600, "0193", 13200},
	{"0690", 8000, "0689", 11500},
	{"0683", 9500, "0682", 14800},
	{"0378", 10200, "0377", 17600},
	{"0603", 8700, "0602", 13900},
	{"0470", 19000, "0469", 29000},
	{"0065", 10000, "0064", 17900},
	{"0247", 21000, "", 0},
	{"0460", 39000, "", 0},
	{"0392", 7600, "", 0},
	{"0313", 7200, "", 0},
	{"0310", 6100, "", 0},
	{"0812", 8100, "", 0},
}

// Codes for which one jurisdiction's local coverage determinations deny
// some lines (pattern 8, variantLocalCoverage in vendor.go, applies here).
var lcdCodes = map[string]bool{
	"64483": true, "20610": true, "97110": true, "97530": true, "97140": true,
	"97112": true, "11042": true, "17000": true, "93000": true, "83036": true,
	"84443": true, "80061": true, "29881": true, "66984": true, "11102": true,
	"96372": true, "97116": true, "10060": true, "12001": true, "64721": true,
	"26055": true, "82947": true, "87086": true, "81002": true, "97161": true,
	"A0428": true, "E0601": true, "A7034": true, "A4239": true, "E0260": true,
}

// generator holds the single seeded source of randomness. Every drawing
// operation is a method on it; plain functions take the drawn values.
type generator struct {
	r *rand.Rand
}

func newGenerator(seed int64) generator {
	return generator{r: rand.New(rand.NewPCG(uint64(seed), 0x835))}
}

func (g generator) float() float64 {
	return g.r.Float64()
}

func (g generator) intn(n int) int {
	return g.r.IntN(n)
}

func (g generator) between(lo, hi int) int {
	return lo + g.r.IntN(hi-lo+1)
}

func (g generator) chance(p float64) bool {
	return g.r.Float64() < p
}

func (g generator) lognormal(median, sigma float64) float64 {
	return median * math.Exp(sigma*g.r.NormFloat64())
}

// choose returns an index with probability proportional to its weight.
func (g generator) choose(weights []float64) int {
	total := 0.0
	for _, w := range weights {
		total += w
	}
	x := g.r.Float64() * total
	for i, w := range weights {
		if x < w {
			return i
		}
		x -= w
	}
	return len(weights) - 1
}

func (g generator) service(list []service) service {
	weights := make([]float64, len(list))
	for i, s := range list {
		weights[i] = s.Weight
	}
	return list[g.choose(weights)]
}

func (g generator) digits(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + g.r.IntN(10))
	}
	return string(b)
}

// npi returns an invented NPI in the 1999 range with a valid Luhn check digit.
func (g generator) npi() string {
	base := "1999" + g.digits(5)
	return base + npiCheckDigit(base)
}

// npiCheckDigit computes the NPI check digit: the Luhn check over the
// 80840 prefix and the nine base digits.
func npiCheckDigit(base string) string {
	sum := 24 // the doubled-digit sum of the 80840 prefix
	for i := 0; i < len(base); i++ {
		d := int(base[i] - '0')
		if i%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return fmt.Sprint((10 - sum%10) % 10)
}

// mbiLetters excludes S, L, O, I, B, and Z, as the MBI format requires.
const mbiLetters = "ACDEFGHJKMNPQRTUVWXY"
const mbiAlnum = "0123456789" + mbiLetters

// mbi returns an invented Medicare Beneficiary Identifier in the published
// format: C A AN N A AN N A A N N.
func (g generator) mbi() string {
	pick := func(set string) byte { return set[g.r.IntN(len(set))] }
	b := []byte{
		byte('1' + g.r.IntN(9)),
		pick(mbiLetters),
		pick(mbiAlnum),
		pick("0123456789"),
		pick(mbiLetters),
		pick(mbiAlnum),
		pick("0123456789"),
		pick(mbiLetters),
		pick(mbiLetters),
		pick("0123456789"),
		pick("0123456789"),
	}
	return string(b)
}

// quantile draws from a distribution given by (cumulative probability,
// value) points, interpolating linearly between them.
func (g generator) quantile(points [][2]float64) float64 {
	u := g.r.Float64()
	for i := 1; i < len(points); i++ {
		if u <= points[i][0] {
			lo, hi := points[i-1], points[i]
			return lo[1] + (u-lo[0])/(hi[0]-lo[0])*(hi[1]-lo[1])
		}
	}
	return points[len(points)-1][1]
}

// Professional allowed ÷ charge by service (REALISM.md section 2.3).
var profRatio = [][2]float64{
	{0, 0.10}, {0.10, 0.157}, {0.25, 0.205}, {0.50, 0.350}, {0.75, 0.426},
	{0.90, 0.507}, {0.95, 0.591}, {1, 0.70},
}

var firstNames = []string{
	"ALICE", "BRUNO", "CAROL", "DMITRI", "ELENA", "FRANK", "GRACE", "HECTOR",
	"IRENE", "JAMAL", "KAREN", "LUIS", "MARTA", "NATHAN", "OLGA", "PEDRO",
	"QUINN", "ROSA", "SAMUEL", "TERESA", "ULRICH", "VERA", "WALTER", "XENIA",
	"YUSUF", "ZELDA",
}

var insurers = []struct{ Name, ID string }{
	{"SYNTHETIC SUPPLEMENTAL INSURER A", "SYNSUP01"},
	{"SYNTHETIC SUPPLEMENTAL INSURER B", "SYNSUP02"},
	{"SYNTHETIC MEDIGAP PLAN C", "SYNSUP03"},
}

func (g generator) patient() Patient {
	p := Patient{
		First:  firstNames[g.intn(len(firstNames))],
		Middle: string(mbiLetters[g.intn(len(mbiLetters))]),
		MBI:    g.mbi(),
	}
	if g.chance(0.35) {
		ins := insurers[g.intn(len(insurers))]
		p.Supplemental = true
		p.Insurer = ins.Name
		p.InsurerID = ins.ID
	}
	return p
}

// jurisdiction is one MAC jurisdiction: its code, its share of the national
// workload, the workload numbers gen writes as the payer identifier
// (N1*PR N104), and the state its payer and payees sit in.
type jurisdiction struct {
	Code     string  // JE, JF, JH, JJ, JK, JL, JM, JN, J5, J6, J8, J15; JA to JD for DME
	Share    float64 // the jurisdiction's share of the national A/B (or DME) workload
	Locality int64   // percent applied to national fees (synthetic)
	State    string
	Zip      string
	NumberA  string // the Part A workload number used here; the one DME number for a DME jurisdiction
	NumberB  string // the Part B workload number used here; "" for DME
}

// Every A/B jurisdiction, largest first. The workload numbers are real:
// CMS's contractor numbers for one state's Part A and Part B workload in
// each jurisdiction, copied from backend/pkg/macs/macs.go (from CMS's
// change requests for each jurisdiction's award); gen cannot import that
// package, since it is its own module. The shares are each jurisdiction's
// percent of the national Part A/B workload as of 2024-09-30 from CMS's
// jurisdiction pages (docs/testdata/REALISM.md section 1.3). Everything
// else about a jurisdiction here is invented: the locality percent, and
// the state's payer and payees.
var abJurisdictions = []jurisdiction{
	{"JH", 0.139, 96, "TX", "75001", "04411", "04412"},
	{"JK", 0.118, 112, "NY", "12201", "13201", "13202"},
	{"JL", 0.117, 104, "PA", "17101", "12501", "12502"},
	{"JE", 0.109, 108, "CA", "95814", "01111", "01112"},
	{"JN", 0.081, 99, "FL", "32301", "09101", "09102"},
	{"JM", 0.080, 95, "NC", "27601", "11501", "11502"},
	{"J6", 0.075, 101, "IL", "62701", "06101", "06102"},
	{"JF", 0.075, 103, "WA", "98501", "02401", "02402"},
	{"JJ", 0.062, 94, "GA", "30301", "10211", "10212"},
	{"J5", 0.054, 93, "MO", "65101", "05301", "05302"},
	{"J15", 0.047, 97, "OH", "43215", "15201", "15202"},
	{"J8", 0.043, 98, "MI", "48901", "08201", "08202"},
}

// The four DME MAC jurisdictions with their numbers (the same sources) and
// their shares of the national DME workload (REALISM.md section 1.3).
var dmeJurisdictions = []jurisdiction{
	{"JC", 0.360, 100, "SC", "29201", "18003", ""},
	{"JD", 0.248, 100, "CA", "95814", "19003", ""},
	{"JA", 0.209, 100, "NY", "12201", "16013", ""},
	{"JB", 0.183, 100, "IL", "62701", "17013", ""},
}

// The jurisdictions whose local coverage decisions deny some lines the
// others pay (pattern 8 lives here): a mid-sized one and a small one.
var lcdJurisdictions = map[string]bool{"JN": true, "JJ": true}

type payeeKind struct {
	Kind   string
	Label  string
	Median float64
}

var partBKinds = []payeeKind{
	{"clinic", "FAMILY CLINIC", 1.0},
	{"clinic", "INTERNAL MEDICINE", 1.0},
	{"clinic", "PRIMARY CARE", 0.8},
	{"group", "MULTISPECIALTY GROUP", 2.5},
	{"surgery", "SURGICAL ASSOCIATES", 0.6},
	{"lab", "CLINICAL LABORATORY", 1.8},
	{"ambulance", "AMBULANCE SERVICE", 1.0},
	{"therapy", "PHYSICAL THERAPY", 1.0},
}

var partAKinds = []payeeKind{
	{"hospital", "GENERAL HOSPITAL", 1},
	{"hospital", "REGIONAL MEDICAL CENTER", 1},
	{"cah", "CRITICAL ACCESS HOSPITAL", 1},
	{"snf", "SKILLED NURSING FACILITY", 1},
	{"hha", "HOME HEALTH AGENCY", 1},
	{"hospice", "HOSPICE", 1},
	{"esrd", "DIALYSIS CENTER", 1},
	{"rhc", "RURAL HEALTH CLINIC", 1},
	{"orf", "OUTPATIENT REHABILITATION", 1},
}

var dmeKinds = []payeeKind{
	{"oxygen", "OXYGEN SUPPLY", 1.2},
	{"cpap", "SLEEP SUPPLY", 1.0},
	{"diabetic", "DIABETIC SUPPLY", 1.0},
	{"dme", "MEDICAL EQUIPMENT", 0.8},
}

// Institutional claim subtypes (REALISM.md section 1.2): bill type, the
// payee kind that bills it, and its share of institutional claims.
var institutionalMix = []struct {
	BillType string
	Kind     string
	Share    float64
}{
	{"111", "hospital", 0.047},
	{"211", "snf", 0.016},
	{"321", "hha", 0.043},
	{"811", "hospice", 0.031},
	{"131", "hospital", 0.863 * 0.55},
	{"851", "cah", 0.863 * 0.10},
	{"711", "rhc", 0.863 * 0.15},
	{"721", "esrd", 0.863 * 0.05},
	{"741", "orf", 0.863 * 0.15},
}

func servicesFor(kind string) []service {
	switch kind {
	case "clinic", "group":
		out := append([]service(nil), clinicFirst...)
		out = append(out, clinicMore...)
		if kind == "group" {
			out = append(out, surgeryServices...)
		}
		return out
	case "surgery":
		return surgeryServices
	case "lab":
		return labServices
	case "ambulance":
		return append(append([]service(nil), ambulanceBase...), mileage)
	case "therapy":
		return therapyServices
	case "oxygen":
		return oxygenServices
	case "cpap":
		return cpapServices
	case "diabetic":
		return diabeticServices
	case "dme":
		return generalDMEServices
	case "hospital", "cah":
		return append(append([]service(nil), outpatientFirst...), outpatientMore...)
	case "rhc":
		return rhcServices
	case "esrd":
		return esrdServices
	case "orf":
		return orfServices
	case "snf":
		return snfServices
	case "hha":
		return hhaServices
	case "hospice":
		return hospiceServices
	}
	return nil
}

// payee creates one billing provider and draws its charge list.
func (g generator) payee(w Workload, k payeeKind, n int) Payee {
	p := Payee{
		Name:    fmt.Sprintf("SYNTHETIC %s %s%d", k.Label, w.Juris, n),
		Kind:    k.Kind,
		NPI:     g.npi(),
		TaxID:   "99" + g.digits(7),
		Street:  fmt.Sprintf("%d SYNTHETIC AVE", 100+g.intn(900)),
		Zip:     w.Zip,
		DFI:     "99" + g.digits(7),
		Account: g.digits(10),
		Weight:  g.lognormal(k.Median, 0.8),
		Code:    fmt.Sprintf("%s%s%d", w.Juris, w.Part, n),
		Prices:  map[string]price{},
	}
	switch h := g.float(); {
	case h < 0.6:
		p.Habit = "unique"
	case h < 0.9:
		p.Habit = "patient"
	default:
		p.Habit = "episode"
	}
	switch k.Kind {
	case "clinic", "group", "surgery", "therapy":
		count := g.between(2, 6)
		if k.Kind == "group" {
			count = g.between(6, 14)
		}
		for i := 0; i < count; i++ {
			p.Renderers = append(p.Renderers, Renderer{Name: firstNames[g.intn(len(firstNames))], NPI: g.npi()})
		}
	}
	// Hospital markup: mode about 2.4, mean about 3.4 (REALISM.md section 2.3).
	p.Markup = 1.4 + g.lognormal(1.3, 0.7)
	p.IPRatio = clamp(g.lognormal(0.20, 0.35), 0.08, 0.40)
	opRatio := clamp(g.lognormal(0.145, 0.45), 0.05, 0.40)
	if k.Kind == "cah" {
		opRatio = clamp(g.lognormal(0.30, 0.3), 0.15, 0.60)
	}
	for _, s := range servicesFor(k.Kind) {
		if _, done := p.Prices[s.Code]; done {
			continue
		}
		p.Prices[s.Code] = g.price(w, p, s, opRatio)
	}
	return p
}

// price draws one code's charge and fee for a payee.
func (g generator) price(w Workload, p Payee, s service, opRatio float64) price {
	switch s.Class {
	case "pps", "zero":
		charge := Money(math.Round(float64(s.Cents)*g.lognormal(1, 0.2)/100) * 100)
		pr := price{Charge: charge}
		if s.Class == "pps" {
			pr.Fee = Money(math.Round(float64(charge) * clamp(g.lognormal(0.6, 0.15), 0.3, 0.9)))
		}
		return pr
	}
	if p.Kind == "hospital" || p.Kind == "cah" || p.Kind == "rhc" || p.Kind == "esrd" || p.Kind == "orf" {
		charge := Money(math.Round(float64(s.Cents)*p.Markup/2.8*g.lognormal(1, 0.25)/100) * 100)
		if charge < 1000 {
			charge = 1000
		}
		ratio := clamp(opRatio*g.lognormal(1, 0.3), 0.03, 0.9)
		return price{Charge: charge, Fee: Money(math.Round(float64(charge) * ratio))}
	}
	fee := percentOf(s.Cents, w.Locality)
	ratio := g.quantile(profRatio)
	if w.Part == "DME" {
		ratio = clamp(g.lognormal(0.426, 0.35), 0.20, 0.80)
	}
	if s.Class == "amb" || s.Class == "mile" {
		ratio = clamp(g.lognormal(0.45, 0.2), 0.3, 0.7)
	}
	// Providers set whole-dollar charges at or above the fee.
	charge := Money(math.Ceil(float64(fee)/ratio/100) * 100)
	if s.Class == "mile" {
		charge = Money(math.Ceil(float64(fee)/ratio/25) * 25)
	}
	pr := price{Charge: charge, Fee: fee}
	if s.Class == "mile" {
		pr.Urban = percentOf(urbanMileFee, w.Locality)
	}
	return pr
}

func clamp(x, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, x))
}

func (g generator) workloads() []Workload {
	var out []Workload
	add := func(j jurisdiction, part string, number string, name string, kinds []payeeKind) {
		w := Workload{
			Number:   number,
			Name:     name,
			Juris:    j.Code,
			Part:     part,
			Locality: j.Locality,
			TaxID:    "98" + g.digits(7),
			DFI:      "98" + g.digits(7),
			Account:  g.digits(10),
			Region:   fmt.Sprintf("%d%d", 1+len(out)/9, 1+len(out)%9),
			State:    j.State,
			Zip:      j.Zip,
		}
		for i, k := range kinds {
			w.Payees = append(w.Payees, g.payee(w, k, i+1))
		}
		out = append(out, w)
	}
	for _, j := range abJurisdictions {
		add(j, "A", j.NumberA, "SYNTHETIC MAC "+j.Code+" PART A", partAKinds)
		add(j, "B", j.NumberB, "SYNTHETIC MAC "+j.Code+" PART B", partBKinds)
	}
	for _, j := range dmeJurisdictions {
		add(j, "DME", j.NumberA, "SYNTHETIC DME MAC "+j.Code, dmeKinds)
	}
	return out
}

const dateLayout = "20060102"

func ymd(t time.Time) string {
	return t.Format(dateLayout)
}

func parseYMD(s string) time.Time {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		panic(fmt.Sprintf("gen: parsing date %q: %v", s, err))
	}
	return t
}

func addDays(s string, n int) string {
	return ymd(parseYMD(s).AddDate(0, 0, n))
}

// Dates of service run October 1 through November 30, 2025 (the two months
// the generated set covers); serviceDays is how many days that is.
var serviceStart = time.Date(2025, time.October, 1, 0, 0, 0, 0, time.UTC)

const serviceDays = 61

// weekday moves a Saturday or Sunday date of service to the Monday after,
// or to the Friday before when that Monday is past the last date of
// service, so every date stays inside the two months.
func weekday(dos time.Time) time.Time {
	next := dos
	for next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
		next = next.AddDate(0, 0, 1)
	}
	if next.Before(serviceStart.AddDate(0, 0, serviceDays)) {
		return next
	}
	for dos.Weekday() == time.Saturday || dos.Weekday() == time.Sunday {
		dos = dos.AddDate(0, 0, -1)
	}
	return dos
}

// paymentCycles returns the MAC payment dates: every Friday from
// October 17, 2025 through January 30, 2026, so a claim is paid a few weeks
// after its date of service and each workload writes a file a week.
func paymentCycles() []string {
	var out []string
	last := time.Date(2026, time.January, 30, 0, 0, 0, 0, time.UTC)
	for t := time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC); !t.After(last); t = t.AddDate(0, 0, 7) {
		out = append(out, ymd(t))
	}
	return out
}

// cycleFor returns the first payment cycle at least 14 days after receipt.
func cycleFor(cycles []string, received string) int {
	earliest := addDays(received, 14)
	for i, c := range cycles {
		if c >= earliest {
			return i
		}
	}
	return len(cycles) - 1
}

// world is everything generated from the seed before any vendor output is built.
type world struct {
	Workloads []Workload
	Cycles    []string
	Claims    []Claim      // MAC claims in MAC order, adjudicated
	Payments  []Payment    // MAC payments
	Extras    [][]Claim    // per vendor: vendor-only claims, adjudicated
	Draws     []vendorDraw // per vendor
}

// vendorDraw holds a vendor's draws, indexed by Claim.ID.
type vendorDraw struct {
	Gap       [][]float64   // [claim][variant]: the variant applies while this is below its pattern's rate
	Regress   [][][]float64 // [submission][claim][variant]: a fixed variant comes back while this is below the regression rate
	Missing   []float64     // [claim]: the lowest values are the claims the vendor does not return (pattern 11)
	Duplicate []float64     // [claim]: the claim is sent twice while this is below the rate (pattern 6)
	Chain     []float64     // [claim]: a chain's original; its chain is reported as a net adjustment while this is below the rate (pattern 9)
}

// patientPool keeps each jurisdiction's beneficiaries so some see several providers.
type patientPool struct {
	Juris    []string
	Patients [][]Patient
}

func (g generator) poolPatient(pool patientPool, juris string) (patientPool, Patient) {
	for i, j := range pool.Juris {
		if j != juris {
			continue
		}
		if len(pool.Patients[i]) > 0 && g.chance(0.4) {
			return pool, pool.Patients[i][g.intn(len(pool.Patients[i]))]
		}
		p := g.patient()
		pool.Patients[i] = append(pool.Patients[i], p)
		return pool, p
	}
	p := g.patient()
	pool.Juris = append(pool.Juris, juris)
	pool.Patients = append(pool.Patients, []Patient{p})
	return pool, p
}

// Line-count distributions, capped at six lines.
var profLineCounts = []float64{0.50, 0.22, 0.11, 0.07, 0.06, 0.04}
var dmeLineCounts = []float64{0.45, 0.25, 0.10, 0.08, 0.07, 0.05}
var opLineCounts = []float64{0.12, 0.15, 0.18, 0.20, 0.17, 0.18}

// Denial reasons for lines (REALISM.md section 4.2). Unprocessable
// (CO-16) is handled separately at claim level.
var lineDenials = []struct {
	Denial Denial
	Share  float64
}{
	{Denial{"CO", "97", "M15"}, 0.20},
	{Denial{"CO", "18", ""}, 0.12},
	{Denial{"CO", "50", "N115"}, 0.08},
	{Denial{"CO", "151", "N640"}, 0.04},
	{Denial{"CO", "96", "N130"}, 0.05},
	{Denial{"PR", "96", "N130"}, 0.03},
	{Denial{"CO", "109", "N418"}, 0.05},
	{Denial{"CO", "B7", "N95"}, 0.06},
	{Denial{"CO", "4", "N519"}, 0.06},
	{Denial{"CO", "29", ""}, 0.03},
	{Denial{"CO", "27", "N30"}, 0.05},
	{Denial{"CO", "236", "N20"}, 0.05},
}

var unprocessableRemarks = []string{"M51", "N290", "MA63", "N382"}

func (g generator) lineDenial() Denial {
	weights := make([]float64, len(lineDenials))
	for i, d := range lineDenials {
		weights[i] = d.Share
	}
	return lineDenials[g.choose(weights)].Denial
}

// deductibleRemaining draws the unmet Part B deductible for a claim in the
// given month (REALISM.md section 2.4: met early in the year), out of the
// 2025 deductible, $257, the year of the dates of service.
func (g generator) deductibleRemaining(month int) Money {
	p := 0.01
	switch month {
	case 1:
		p = 0.35
	case 2:
		p = 0.15
	case 3:
		p = 0.06
	}
	if !g.chance(p) {
		return 0
	}
	if g.chance(0.5) {
		return 25700
	}
	return Money(g.between(10, 257) * 100)
}

// newLine prices one submitted line at the payee's charge list.
func newLine(p Payee, s service, units int, date string) Line {
	pr := p.Prices[s.Code]
	ln := Line{
		Qual:       s.Qual,
		Code:       s.Code,
		Rev:        s.Rev,
		Class:      s.Class,
		Units:      units,
		Date:       date,
		UnitCharge: pr.Charge,
		UnitFee:    pr.Fee,
		UrbanFee:   pr.Urban,
		Exempt:     s.Class == "lab",
		Packaged:   s.Class == "pack",
	}
	if s.Mod != "" {
		ln.Mods = []string{s.Mod}
	}
	return ln
}

// baseClaim draws one claim: type, workload, payee, patient, dates, lines,
// and the reference system's decisions about it.
func (g generator) baseClaim(ws []Workload, cycles []string, pool patientPool, id int) (patientPool, Claim) {
	var wi int
	var typ string
	switch x := g.float(); {
	case x < 0.77:
		typ = "professional"
		wi = 2*g.choose(jurisShares(abJurisdictions)) + 1
	case x < 0.83:
		typ = "dme"
		wi = 2*len(abJurisdictions) + g.choose(jurisShares(dmeJurisdictions))
	default:
		typ = "institutional"
		wi = 2 * g.choose(jurisShares(abJurisdictions))
	}
	w := ws[wi]
	c := Claim{ID: id, DrawID: id, Workload: wi, Juris: w.Juris, Part: w.Part, Type: typ, OrigID: -1, Pick: g.intn(1 << 30)}

	if typ == "institutional" {
		shares := make([]float64, len(institutionalMix))
		for i, m := range institutionalMix {
			shares[i] = m.Share
		}
		m := institutionalMix[g.choose(shares)]
		c.BillType = m.BillType
		c.Facility = m.BillType[:2]
		var cand []int
		var weights []float64
		for i, p := range w.Payees {
			if p.Kind == m.Kind {
				cand = append(cand, i)
				weights = append(weights, p.Weight)
			}
		}
		c.Payee = cand[g.choose(weights)]
	} else {
		weights := make([]float64, len(w.Payees))
		for i, p := range w.Payees {
			weights[i] = p.Weight
		}
		c.Payee = g.choose(weights)
	}
	p := w.Payees[c.Payee]

	pool, c.Patient = g.poolPatient(pool, w.Juris)
	c.Crossover = c.Patient.Supplemental

	dos := serviceStart.AddDate(0, 0, g.intn(serviceDays))
	if p.Kind != "ambulance" && p.Kind != "hospital" {
		dos = weekday(dos)
	}
	c.Month = int(dos.Month())
	date := ymd(dos)

	switch typ {
	case "professional":
		c = g.professionalLines(c, p, date)
		c.Facility = "11"
		switch p.Kind {
		case "lab":
			c.Facility = "81"
		case "ambulance":
			c.Facility = "41"
		case "surgery":
			c.Facility = "22"
		}
		if len(p.Renderers) > 0 {
			c.Rendering = p.Renderers[g.intn(len(p.Renderers))]
			// Some MAC claims carry no rendering provider loop at all.
			c.OmitRendering = g.chance(0.12)
		}
		c.Received = addDays(date, g.between(2, 20))
	case "dme":
		c = g.dmeLines(c, p, date)
		c.Facility = "12"
		c.Received = addDays(date, g.between(2, 20))
	default:
		c = g.institutionalLines(c, p, dos)
		c.Frequency = "1"
		c.Received = addDays(c.Thru, g.between(3, 20))
	}
	c.Cycle = cycleFor(cycles, c.Received)
	c.PatientCtl = g.patientControl(p, c)

	// Part B deductible applies to professional, DME, and outpatient claims.
	partB := typ != "institutional" || c.BillType == "131" || c.BillType == "851" || c.BillType == "711" || c.BillType == "721" || c.BillType == "741"
	if partB {
		c.Deductible = g.deductibleRemaining(c.Month)
	}
	c.AltDeductible = Money(g.between(20, 257) * 100)
	c = g.outcome(c, partB)
	return pool, c
}

func jurisShares(js []jurisdiction) []float64 {
	out := make([]float64, len(js))
	for i, j := range js {
		out[i] = j.Share
	}
	return out
}

func (g generator) patientControl(p Payee, c Claim) string {
	switch p.Habit {
	case "patient":
		return fmt.Sprintf("%s-%s", p.Code, c.Patient.MBI[7:])
	case "episode":
		return fmt.Sprintf("%s-%s-%02d", p.Code, c.Patient.MBI[7:], c.Month)
	}
	return fmt.Sprintf("%s-%06d", p.Code, 100000+g.intn(900000))
}

func (g generator) professionalLines(c Claim, p Payee, date string) Claim {
	n := 1 + g.choose(profLineCounts)
	switch p.Kind {
	case "clinic", "group":
		c.Lines = append(c.Lines, newLine(p, g.service(clinicFirst), 1, date))
		more := clinicMore
		if p.Kind == "group" {
			more = append(append([]service(nil), clinicMore...), surgeryServices...)
		}
		for len(c.Lines) < n {
			c.Lines = addDistinct(c.Lines, newLine(p, g.service(more), 1, date))
			if g.chance(0.1) {
				break
			}
		}
	case "surgery":
		n = 1 + g.choose([]float64{0.45, 0.40, 0.15})
		for len(c.Lines) < n {
			c.Lines = addDistinct(c.Lines, newLine(p, g.service(surgeryServices), 1, date))
			if g.chance(0.05) {
				break
			}
		}
	case "lab":
		n = 1 + g.choose([]float64{0.30, 0.25, 0.20, 0.12, 0.08, 0.05})
		for len(c.Lines) < n {
			s := g.service(labServices)
			units := 1
			if s.Code == "P9603" {
				units = g.between(4, s.MaxUnits)
			}
			c.Lines = addDistinct(c.Lines, newLine(p, s, units, date))
			if g.chance(0.05) {
				break
			}
		}
	case "ambulance":
		c.Rural = g.chance(0.35)
		c.Lines = append(c.Lines, newLine(p, g.service(ambulanceBase), 1, date))
		miles := int(clamp(math.Round(g.lognormal(9, 0.7)), 1, float64(mileage.MaxUnits)))
		ln := newLine(p, mileage, miles, date)
		if !c.Rural {
			ln.UnitFee = ln.UrbanFee
		}
		c.Lines = append(c.Lines, ln)
	case "therapy":
		n = 1 + g.choose([]float64{0.35, 0.35, 0.20, 0.10})
		for len(c.Lines) < n {
			s := g.service(therapyServices)
			c.Lines = addDistinct(c.Lines, newLine(p, s, g.between(1, s.MaxUnits), date))
			if g.chance(0.05) {
				break
			}
		}
	}
	// An E/M visit billed with a procedure carries modifier 25.
	surgical := 0
	for _, ln := range c.Lines {
		if ln.Class == "surg" {
			surgical++
		}
	}
	if surgical > 0 && c.Lines[0].Class == "em" {
		c.Lines[0].Mods = []string{"25"}
	}
	// Multiple-procedure reduction: every surgical line but the highest-priced is allowed at 50 percent.
	if surgical >= 2 {
		top := -1
		for i, ln := range c.Lines {
			if ln.Class == "surg" && (top == -1 || ln.UnitFee > c.Lines[top].UnitFee) {
				top = i
			}
		}
		for i := range c.Lines {
			if c.Lines[i].Class == "surg" && i != top {
				c.Lines[i].Reduced = true
			}
		}
	}
	return c
}

// addDistinct appends a line unless the claim already has that code.
func addDistinct(lines []Line, ln Line) []Line {
	for _, l := range lines {
		if l.Code == ln.Code {
			return lines
		}
	}
	return append(lines, ln)
}

func (g generator) dmeLines(c Claim, p Payee, date string) Claim {
	n := 1 + g.choose(dmeLineCounts)
	list := servicesFor(p.Kind)
	for len(c.Lines) < n {
		s := g.service(list)
		c.Lines = addDistinct(c.Lines, newLine(p, s, g.between(1, s.MaxUnits), date))
		if g.chance(0.08) {
			break
		}
	}
	return c
}

func (g generator) institutionalLines(c Claim, p Payee, dos time.Time) Claim {
	date := ymd(dos)
	c.From = date
	c.Thru = date
	monthly := func() {
		first := time.Date(dos.Year(), dos.Month(), 1, 0, 0, 0, 0, time.UTC)
		c.From = ymd(first)
		c.Thru = ymd(first.AddDate(0, 1, -1))
	}
	switch c.BillType {
	case "111":
		d := drgs[g.intn(len(drgs))]
		if g.chance(0.5) {
			d = drgs[g.intn(12)]
		}
		c.DRG = d.Code
		c.Weight = d.Weight
		c.Days = int(clamp(math.Round(g.lognormal(4, 0.5)), 1, 20))
		c.Thru = addDays(date, c.Days)
		c.BaseRate = percentOf(680000, 100+(int64(g.intn(40))-15))
		allowed := drgAmount(c.BaseRate, c.Weight)
		c.InpatientCharge = Money(math.Round(float64(allowed)/p.IPRatio/100) * 100)
		c.PartADeductible = g.chance(0.35)
		// A few stays run far past the DRG's cost and earn an outlier payment.
		c.Outlier = g.chance(0.08)
		if c.Outlier {
			c.InpatientCharge = Money(math.Round(float64(c.InpatientCharge)*g.lognormal(2.4, 0.2)/100) * 100)
		}
		return c
	case "131", "851":
		n := 1 + g.choose(opLineCounts)
		c.Lines = append(c.Lines, newLine(p, g.service(outpatientFirst), 1, date))
		for len(c.Lines) < n {
			c.Lines = addDistinct(c.Lines, newLine(p, g.service(outpatientMore), 1, date))
			if g.chance(0.05) {
				break
			}
		}
	case "711":
		c.Lines = append(c.Lines, newLine(p, rhcServices[0], 1, date))
		if g.chance(0.4) {
			c.Lines = append(c.Lines, newLine(p, rhcServices[1], 1, date))
		}
	case "721":
		monthly()
		c.Lines = append(c.Lines, newLine(p, esrdServices[0], g.between(12, 14), c.From))
		if g.chance(0.6) {
			c.Lines = append(c.Lines, newLine(p, esrdServices[1], 1, c.From))
		}
	case "741":
		monthly()
		n := g.between(1, 3)
		for len(c.Lines) < n {
			s := g.service(orfServices)
			c.Lines = addDistinct(c.Lines, newLine(p, s, g.between(4, 12), c.From))
		}
	case "211", "321":
		monthly()
		list := snfServices
		if c.BillType == "321" {
			list = hhaServices
			c.Thru = addDays(c.From, 29)
		}
		n := g.between(2, len(list))
		for i := 0; i < n; i++ {
			c.Lines = append(c.Lines, newLine(p, list[i], 1, c.From))
		}
	case "811":
		monthly()
		days := int(parseYMD(c.Thru).Sub(parseYMD(c.From)).Hours()/24) + 1
		c.Lines = append(c.Lines, newLine(p, hospiceServices[0], days, c.From))
		if g.chance(0.5) {
			c.Lines = append(c.Lines, newLine(p, hospiceServices[1], g.between(2, 8), c.From))
		}
	}
	// Part A prospective-payment bills have no patient cost sharing here.
	if c.BillType == "211" || c.BillType == "321" || c.BillType == "811" {
		for i := range c.Lines {
			c.Lines[i].Exempt = true
		}
	}
	return c
}

// drgAmount is the base rate times the DRG weight, rounded to the cent.
func drgAmount(base Money, weight int64) Money {
	return Money((int64(base)*weight + 5000) / 10000)
}

// outcome draws the reference system's outcome (REALISM.md section 3) and
// records it as decisions on the claim and its lines.
func (g generator) outcome(c Claim, partB bool) Claim {
	shares := []float64{0.84, 0.06, 0.05, 0.05} // paid, partial, denied, zero-pay
	switch {
	case c.Type == "dme":
		shares = []float64{0.81, 0.07, 0.08, 0.04}
	case c.Type == "institutional":
		shares = []float64{0.90, 0.03, 0.04, 0.03}
	}
	o := g.choose(shares)
	if c.inpatient() {
		if o == 2 {
			c.Denial = Denial{"CO", "50", "N115"}
		}
		return c
	}
	switch o {
	case 1:
		if len(c.Lines) >= 2 {
			i := 1 + g.intn(len(c.Lines)-1)
			c.Lines[i].Denial = g.lineDenial()
		}
	case 2:
		if c.Type != "institutional" && g.chance(0.18) {
			c.Unprocessable = true
			remark := unprocessableRemarks[g.intn(len(unprocessableRemarks))]
			for i := range c.Lines {
				c.Lines[i].Denial = Denial{"CO", "16", remark}
			}
			break
		}
		d := g.lineDenial()
		if c.Type == "institutional" {
			d = []Denial{{"CO", "50", "N115"}, {"CO", "96", "N130"}, {"CO", "29", ""}}[g.intn(3)]
		}
		for i := range c.Lines {
			c.Lines[i].Denial = d
		}
	case 3:
		if partB {
			c.Deductible = 25700
		}
	}
	// Oxygen past the rental cap is denied (pattern 5, variantOxygenMonth in
	// vendor.go, toggles this).
	for i, ln := range c.Lines {
		if ln.Class == "oxygen" && ln.Denial.Reason == "" && g.chance(0.10) {
			c.Lines[i].Denial = Denial{"CO", "119", "N640"}
		}
	}
	// Two jurisdictions' local coverage determinations deny some lines
	// the others pay (pattern 8 is built on these).
	if lcdJurisdictions[c.Juris] && c.Part == "B" {
		for i, ln := range c.Lines {
			if lcdCodes[ln.Code] && ln.Denial.Reason == "" && g.chance(0.16) {
				c.Lines[i].Denial = Denial{"CO", "50", "N115"}
			}
		}
	}
	// A few claims have another payer ahead of Medicare (processed as
	// secondary, CLP02 2), and a very few went to the wrong contractor and
	// were forwarded (CLP02 23, nothing paid).
	if !c.inpatient() && !c.Unprocessable {
		switch {
		case g.chance(0.003):
			c.NotOurs = true
			for i := range c.Lines {
				c.Lines[i].Denial = Denial{"OA", "109", ""}
			}
		case g.chance(0.03):
			c.Secondary = true
		}
	}
	return c
}

// generateWorld draws the full MAC side, the vendor-only extras, and every
// vendor draw, in a fixed order, so one seed always gives one world.
func generateWorld(cfg config) world {
	g := newGenerator(cfg.Seed)
	w := world{Workloads: g.workloads(), Cycles: paymentCycles()}

	chains := max(1, cfg.Claims*7/1000)
	pairs := max(2, cfg.Claims*4/1000)
	base := cfg.Claims - 2*chains - pairs
	var pool patientPool
	var claims []Claim
	for i := 0; i < base; i++ {
		var c Claim
		pool, c = g.baseClaim(w.Workloads, w.Cycles, pool, len(claims))
		claims = append(claims, c)
	}

	// Key collisions: two distinct claims sharing all six key fields. Odd
	// pairs are exact resubmissions (the second denied as a duplicate);
	// even pairs are distinct same-day services (return trips, repeat draws).
	used := make([]bool, len(claims))
	for k := 0; k < pairs; k++ {
		resub := k%2 == 1
		i := g.pairCandidate(w.Workloads, claims, used, resub)
		used[i] = true
		first := claims[i]
		second := first.clone()
		second.ID = len(claims)
		second.DrawID = first.ID
		if resub {
			first.Pair, second.Pair = "resub-first", "resub-second"
			second.Received = addDays(first.Received, g.between(5, 12))
			for j := range second.Lines {
				second.Lines[j].Denial = Denial{"CO", "18", ""}
			}
		} else {
			first.Pair, second.Pair = "distinct-first", "distinct-second"
			if g.chance(0.5) {
				second.PatientCtl = g.patientControl(w.Workloads[first.Workload].Payees[first.Payee], first)
			}
		}
		second.Cycle = cycleFor(w.Cycles, second.Received)
		second.Pick = g.intn(1 << 30)
		claims[i] = first
		claims = append(claims, second)
		used = append(used, true)
	}

	// Adjustment chains: the original is paid in one cycle; a later cycle
	// carries the reversal (CLP02 22, amounts negated, REF*F8 naming the
	// original). A replacement (the provider's frequency code 7) also
	// carries the corrected claim; a void (frequency code 8) is the reversal
	// alone. Institutional claims show the frequency in CLP09.
	for k := 0; k < chains; k++ {
		i := g.chainCandidate(claims, used, len(w.Cycles))
		used[i] = true
		orig := claims[i]
		later := min(orig.Cycle+g.between(2, 5), len(w.Cycles)-1)
		void := k%3 == 2

		rev := orig.clone()
		rev.ID = len(claims)
		rev.DrawID = orig.ID
		rev.Reversal = true
		rev.OrigID = orig.ID
		rev.Cycle = later
		rev.Received = addDays(w.Cycles[later], -20)
		if rev.Type == "institutional" {
			rev.Frequency = "7"
			if void {
				rev.Frequency = "8"
			}
		}
		claims = append(claims, rev)
		used = append(used, true)
		if void {
			continue
		}

		corr := orig.clone()
		corr.ID = len(claims)
		corr.DrawID = corr.ID
		corr.Correction = true
		corr.OrigID = orig.ID
		corr.Cycle = later
		corr.Received = rev.Received
		corr.Lines[0].Units++
		corr.Pick = g.intn(1 << 30)
		if corr.Type == "institutional" {
			corr.Frequency = "7"
		}
		claims = append(claims, corr)
		used = append(used, true)
	}

	// Interest on a few clean paid claims (REALISM.md section 4.3).
	interest := max(3, cfg.Claims*3/1000)
	for k := 0; k < interest; k++ {
		i := g.interestCandidate(claims, used)
		used[i] = true
		claims[i].Interest = Money(clamp(math.Round(g.lognormal(40, 1.4)), 1, 4000))
	}

	for i := range claims {
		claims[i] = assignLineIDs(claims[i])
		claims[i] = adjudicate(claims[i])
	}
	w.Claims = macOrder(w.Workloads, claims)
	w.Payments = g.macPayments(w)

	// Vendor-only claims, the most any submission will need.
	for v := range cfg.Vendors {
		prof := profiles[v%len(profiles)]
		n := 0
		for d := range cfg.Submissions {
			n = max(n, countFor(cfg.Claims, prof.Extra[stage(d)]))
		}
		var extras []Claim
		for i := 0; i < n; i++ {
			var c Claim
			var fresh patientPool
			_, c = g.baseClaim(w.Workloads, w.Cycles, fresh, -1)
			c.PatientCtl = fmt.Sprintf("%s-X%d%02d", w.Workloads[c.Workload].Payees[c.Payee].Code, v+1, i+1)
			c = assignLineIDs(c)
			extras = append(extras, adjudicate(c))
		}
		w.Extras = append(w.Extras, extras)
	}

	for range cfg.Vendors {
		w.Draws = append(w.Draws, g.vendorDraw(len(claims), len(cfg.Submissions)))
	}
	return w
}

func (g generator) vendorDraw(claims, submissions int) vendorDraw {
	var d vendorDraw
	for i := 0; i < claims; i++ {
		row := make([]float64, len(variants))
		for j := range row {
			row[j] = g.float()
		}
		d.Gap = append(d.Gap, row)
	}
	for k := 0; k < submissions; k++ {
		var rows [][]float64
		for i := 0; i < claims; i++ {
			row := make([]float64, len(variants))
			for j := range row {
				row[j] = g.float()
			}
			rows = append(rows, row)
		}
		d.Regress = append(d.Regress, rows)
	}
	for i := 0; i < claims; i++ {
		d.Missing = append(d.Missing, g.float())
	}
	for i := 0; i < claims; i++ {
		d.Duplicate = append(d.Duplicate, g.float())
	}
	for i := 0; i < claims; i++ {
		d.Chain = append(d.Chain, g.float())
	}
	return d
}

// pairCandidate picks an unused paid claim for a key collision: an
// ambulance or lab claim for distinct same-day services, a clinic claim for
// a resubmission.
func (g generator) pairCandidate(ws []Workload, claims []Claim, used []bool, resub bool) int {
	var cand []int
	for i, c := range claims {
		if used[i] || c.Type != "professional" || c.Unprocessable || c.Cycle >= len(paymentCycles())-1 {
			continue
		}
		kind := ws[c.Workload].Payees[c.Payee].Kind
		if resub && kind == "clinic" || !resub && (kind == "ambulance" || kind == "lab") {
			cand = append(cand, i)
		}
	}
	if len(cand) == 0 {
		for i := range claims {
			if !used[i] && claims[i].Type == "professional" {
				cand = append(cand, i)
			}
		}
	}
	return cand[g.intn(len(cand))]
}

func (g generator) chainCandidate(claims []Claim, used []bool, cycles int) int {
	// Adjustment bills are more common on institutional claims, so half the
	// chains are drawn from institutional claims when any qualify.
	institutional := g.chance(0.5)
	var cand, any []int
	for i, c := range claims {
		if used[i] || c.inpatient() || c.Unprocessable || c.Cycle > cycles-6 || len(c.Lines) == 0 {
			continue
		}
		denied := false
		for _, ln := range c.Lines {
			if ln.Denial.Reason != "" {
				denied = true
			}
		}
		if denied {
			continue
		}
		any = append(any, i)
		if !institutional || c.Type == "institutional" {
			cand = append(cand, i)
		}
	}
	if len(cand) == 0 {
		cand = any
	}
	return cand[g.intn(len(cand))]
}

func (g generator) interestCandidate(claims []Claim, used []bool) int {
	var cand []int
	for i, c := range claims {
		if used[i] || c.Type == "institutional" || c.Unprocessable || c.Deductible == 25700 {
			continue
		}
		denied := false
		for _, ln := range c.Lines {
			if ln.Denial.Reason != "" {
				denied = true
			}
		}
		if !denied {
			cand = append(cand, i)
		}
	}
	return cand[g.intn(len(cand))]
}

// assignLineIDs gives each line its REF*6R line item control number.
func assignLineIDs(c Claim) Claim {
	c.Lines = append([]Line(nil), c.Lines...)
	base := c.ID
	if c.Reversal {
		base = c.OrigID
	}
	for i := range c.Lines {
		if c.ID < 0 {
			c.Lines[i].ID = fmt.Sprintf("%s-%d", c.PatientCtl, i+1)
			continue
		}
		c.Lines[i].ID = fmt.Sprintf("L%05d%d", base, i+1)
	}
	return c
}

// macOrder sorts claims into the MAC's output order and assigns ICNs.
func macOrder(ws []Workload, claims []Claim) []Claim {
	out := append([]Claim(nil), claims...)
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Workload != b.Workload {
			return a.Workload < b.Workload
		}
		if a.Cycle != b.Cycle {
			return a.Cycle < b.Cycle
		}
		if a.Payee != b.Payee {
			return a.Payee < b.Payee
		}
		if a.Received != b.Received {
			return a.Received < b.Received
		}
		return a.ID < b.ID
	})
	counters := make([]int, len(ws))
	for i := range out {
		c := out[i]
		if c.Reversal {
			continue
		}
		counters[c.Workload]++
		julian := parseYMD(c.Received).YearDay()
		out[i].ICN = fmt.Sprintf("%s26%03d%06d", ws[c.Workload].Region, julian, counters[c.Workload])
	}
	return linkChains(out)
}

// linkChains copies the original's ICN onto its reversal (CLP07, and REF*F8
// naming the claim undone) and onto its correction (REF*F8).
func linkChains(claims []Claim) []Claim {
	icn := map[int]string{}
	for _, c := range claims {
		if !c.Reversal && c.ID >= 0 {
			icn[c.ID] = c.ICN
		}
	}
	for i, c := range claims {
		if c.Reversal {
			claims[i].ICN = icn[c.OrigID]
			claims[i].OrigICN = icn[c.OrigID]
		}
		if c.Correction && !c.NetAdjust {
			claims[i].OrigICN = icn[c.OrigID]
		}
	}
	return claims
}

// countFor turns a per-claim rate into a count, at least one.
func countFor(claims int, rate float64) int {
	return max(1, int(math.Round(float64(claims)*rate)))
}
