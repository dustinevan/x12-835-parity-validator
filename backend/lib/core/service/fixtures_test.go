package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/comparison"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/x12835"
)

// TestTaxonomyFixtures runs every taxonomy fixture (testdata/taxonomy/<ID>)
// through Split and CompareSubmissions under set 1 and compares the
// engine's outcome per key and its differences with the fixture's
// expected.json. testdata/taxonomy/README.md describes expected.json.
//
// A case passes when the engine agrees with the fixture on every outcome
// and every claim- and line-scope difference. The engine compares claims
// only, so file_differences (payment, provider, and envelope scope) are
// counted and logged, not compared. A case listed in
// testdata/taxonomy/KNOWN-DIFFERENCES.md passes when its disagreements are
// exactly the listed ones. Run with -v to see every case.
func TestTaxonomyFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "..", "..", "testdata", "taxonomy")
	known, err := readKnownDifferences(filepath.Join(root, "KNOWN-DIFFERENCES.md"))
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := filepath.Glob(filepath.Join(root, "T*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		t.Fatalf("no fixtures under %s", root)
	}
	seen := map[string]bool{}
	var pass, knownPass, fail, skip int
	for _, dir := range dirs {
		id := filepath.Base(dir)
		seen[id] = true
		t.Run(id, func(t *testing.T) {
			c, err := readExpectedCase(filepath.Join(dir, "expected.json"))
			if err != nil {
				fail++
				t.Fatal(err)
			}
			if !c.Representable {
				skip++
				t.Skipf("not representable: %s", c.Reason)
			}
			mismatches, compared := compareFixture(dir, c)
			result := checkKnown(mismatches, known[id])
			for _, m := range mismatches {
				t.Logf("%s: %s", m.signature, m.detail)
			}
			if len(c.FileDifferences) > 0 {
				t.Logf("%d file-level differences not compared (the engine compares claims only)", len(c.FileDifferences))
			}
			switch {
			case len(result) > 0:
				fail++
				for _, r := range result {
					t.Error(r)
				}
			case len(mismatches) > 0:
				knownPass++
				t.Logf("%d engine differences compared", compared)
				t.Logf("known: %s", strings.Join(known[id], "; "))
			default:
				pass++
				t.Logf("pass: %d claim- and line-scope differences agree", compared)
			}
		})
	}
	for id := range known {
		if !seen[id] {
			t.Errorf("KNOWN-DIFFERENCES.md lists %s, which has no fixture directory", id)
		}
	}
	t.Logf("taxonomy fixtures: %d pass, %d pass with known differences, %d fail, %d skipped", pass, knownPass, fail, skip)
}

// checkKnown compares a case's disagreements with its KNOWN-DIFFERENCES.md
// line and returns the problems: disagreements not listed, and listed ones
// that no longer occur.
func checkKnown(mismatches []mismatch, listed []string) []string {
	got := map[string]bool{}
	for _, m := range mismatches {
		got[m.signature] = true
	}
	want := map[string]bool{}
	for _, s := range listed {
		want[s] = true
	}
	var problems []string
	for _, s := range sortedKeys(got) {
		if !want[s] {
			problems = append(problems, "not in KNOWN-DIFFERENCES.md: "+s)
		}
	}
	for _, s := range sortedKeys(want) {
		if !got[s] {
			problems = append(problems, "listed in KNOWN-DIFFERENCES.md but no longer occurs (remove it): "+s)
		}
	}
	return problems
}

func sortedKeys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// readKnownDifferences reads the table rows of KNOWN-DIFFERENCES.md:
//
//	| T1.06 | outcome extra+matched vs set_aside | why ... |
//
// The second column lists the case's disagreements separated by "; ".
func readKnownDifferences(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "| T") {
			continue
		}
		cols := strings.Split(strings.Trim(line, "|"), "|")
		if len(cols) != 3 {
			return nil, fmt.Errorf("reading %s: want 3 columns: %s", path, line)
		}
		id := strings.TrimSpace(cols[0])
		why := strings.TrimSpace(cols[2])
		if why == "" {
			return nil, fmt.Errorf("reading %s: %s has no reason", path, id)
		}
		if _, dup := out[id]; dup {
			return nil, fmt.Errorf("reading %s: %s is listed twice", path, id)
		}
		for _, s := range strings.Split(cols[1], ";") {
			if s = strings.TrimSpace(strings.Trim(strings.TrimSpace(s), "`")); s != "" {
				out[id] = append(out[id], s)
			}
		}
	}
	return out, sc.Err()
}

// The expected.json convention (testdata/taxonomy/README.md). Unknown
// fields are an error, so the convention cannot drift unnoticed. A null
// value reads as "": a value is never the empty string in a fixture.
type expectedCase struct {
	TaxonomyID      string               `json:"taxonomy_id"`
	Title           string               `json:"title"`
	Representable   bool                 `json:"representable"`
	Reason          string               `json:"reason"`
	Notes           string               `json:"notes"`
	Imbalance       []expectedImbalance  `json:"imbalance"`
	SystemAssigned  []string             `json:"system_assigned_differences_ignored"`
	Claims          []expectedClaim      `json:"claims"`
	FileDifferences []expectedDifference `json:"file_differences"`
}

type expectedImbalance struct {
	Side   string `json:"side"`
	Level  string `json:"level"`
	Detail string `json:"detail"`
}

type expectedClaim struct {
	Key                  expectedKey          `json:"key"`
	PatientControlNumber string               `json:"patient_control_number"`
	Outcome              string               `json:"outcome"`
	Occurrences          expectedOccurrences  `json:"occurrences"`
	Differences          []expectedDifference `json:"differences"`
}

type expectedKey struct {
	MBI           string `json:"mbi"`
	DateOfService string `json:"date_of_service"`
	TotalCharge   string `json:"total_charge"`
	TaxID         string `json:"tax_id"`
	NPI           string `json:"npi"`
	Workload      string `json:"workload"`
}

type expectedOccurrences struct {
	MAC    int `json:"mac"`
	Vendor int `json:"vendor"`
}

type expectedDifference struct {
	Segment           string `json:"segment"`
	Element           int    `json:"element"`
	Occurrence        int    `json:"occurrence"`
	Scope             string `json:"scope"`
	Line              int    `json:"line"`
	VendorLine        int    `json:"vendor_line"`
	Transaction       int    `json:"transaction"`
	VendorTransaction int    `json:"vendor_transaction"`
	Kind              string `json:"kind"`
	SystemAssigned    bool   `json:"system_assigned"`
	MAC               string `json:"mac"`
	Vendor            string `json:"vendor"`
}

func (d expectedDifference) String() string {
	where := d.Scope
	if d.Scope == LevelLine {
		where = fmt.Sprintf("line %d/%d", d.Line, d.VendorLine)
	}
	return fmt.Sprintf("%s %s#%d element %d (%q, %q)", where, d.Segment, d.Occurrence, d.Element, d.MAC, d.Vendor)
}

func readExpectedCase(path string) (expectedCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return expectedCase{}, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var c expectedCase
	if err := dec.Decode(&c); err != nil {
		return expectedCase{}, fmt.Errorf("reading %s: %w", path, err)
	}
	return c, nil
}

// mismatch is one disagreement between the engine and a fixture. The
// signature is what KNOWN-DIFFERENCES.md lists; the detail is logged.
//
//	outcome <fixture> vs <engine>     the key's outcome differs
//	key not found (<fixture outcome>) no engine row has the fixture's key
//	engine key not in fixture (<outcome>)
//	unreported <element>              the fixture lists it; the engine does not
//	unexpected <element>              the engine reports it; the fixture does not
//	value <element>                   both report it with different values
//	not compared <what>               the fixture lists a difference in
//	                                  something the engine does not compare,
//	                                  or that set 1 leaves out
//	fixture <problem>                 expected.json does not fit its files
//	split <file>, compare             the engine returned an error
type mismatch struct {
	signature string
	detail    string
}

// compareFixture runs one fixture through the engine and lists every
// disagreement with expected.json, with the number of engine differences
// it compared.
func compareFixture(dir string, c expectedCase) ([]mismatch, int) {
	var out []mismatch
	macClaims, err := splitFixtureFile(dir, "mac.835")
	if err != nil {
		return []mismatch{{"split mac.835", err.Error()}}, 0
	}
	vendorClaims, err := splitFixtureFile(dir, "vendor.835")
	if err != nil {
		return []mismatch{{"split vendor.835", err.Error()}}, 0
	}
	rows, diffs, _, err := CompareSubmissions(compareSet1, "vendor", macClaims, vendorClaims)
	if err != nil {
		return []mismatch{{"compare", err.Error()}}, 0
	}
	rowByKey := map[string]comparison.Comparison{}
	for _, row := range rows {
		rowByKey[string(row.MatchKey)] = row
	}

	// Expected claims grouped by key, in expected.json order.
	type group struct {
		key      string
		label    string
		outcomes []string
		claims   []expectedClaim
	}
	var groups []group
	index := map[string]int{}
	for _, ec := range c.Claims {
		k := ec.Key
		key, err := MatchKey(Key{MBI: k.MBI, DateOfService: k.DateOfService, TotalCharge: k.TotalCharge, TaxID: k.TaxID, NPI: k.NPI, Workload: k.Workload})
		if err != nil {
			out = append(out, mismatch{"fixture key", err.Error()})
			continue
		}
		i, ok := index[string(key)]
		if !ok {
			i = len(groups)
			index[string(key)] = i
			label := fmt.Sprintf("%s %s %s", k.MBI, k.DateOfService, k.TotalCharge)
			groups = append(groups, group{key: string(key), label: label})
		}
		groups[i].outcomes = append(groups[i].outcomes, ec.Outcome)
		groups[i].claims = append(groups[i].claims, ec)
	}

	for _, g := range groups {
		sort.Strings(g.outcomes)
		want := strings.Join(g.outcomes, "+")
		row, ok := rowByKey[g.key]
		if !ok {
			out = append(out, mismatch{fmt.Sprintf("key not found (%s)", want), g.label})
			continue
		}
		got := row.Outcome
		if got == comparison.OutcomeEqual || got == comparison.OutcomeDifferent {
			got = "matched"
		}
		if got != want {
			out = append(out, mismatch{fmt.Sprintf("outcome %s vs %s", want, row.Outcome), g.label})
			continue
		}
		if want != "matched" {
			continue
		}
		mac, vendor := claimByKey(macClaims, g.key), claimByKey(vendorClaims, g.key)
		var engine []engineDifference
		for _, d := range diffs {
			if string(d.MatchKey) == g.key {
				engine = append(engine, engineDifference{element: d.Element, line: d.LineControlNumber, mac: d.MACValue, vendor: d.VendorValue})
			}
		}
		ms, err := compareClaimDifferences(g.claims[0].Differences, mac, vendor, engine)
		if err != nil {
			out = append(out, mismatch{"fixture claim", fmt.Sprintf("%s: %v", g.label, err)})
			continue
		}
		for _, m := range ms {
			m.detail = g.label + ": " + m.detail
			out = append(out, m)
		}
	}
	for _, row := range rows {
		if _, ok := index[string(row.MatchKey)]; !ok {
			out = append(out, mismatch{fmt.Sprintf("engine key not in fixture (%s)", row.Outcome),
				fmt.Sprintf("%s %s %s", row.MedicareBeneficiaryIdentifier, row.DateOfService.Format("20060102"), row.TotalClaimChargeAmount)})
		}
	}
	return out, len(diffs)
}

// splitFixtureFile splits one fixture 835; a missing file has no claims.
func splitFixtureFile(dir, name string) ([]claim.Claim, error) {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	claims, _, err := Split(name, name, data)
	return claims, err
}

// claimByKey returns the one claim with the key; the caller has checked
// that the key's outcome is matched, so there is exactly one.
func claimByKey(claims []claim.Claim, key string) claim.Claim {
	for _, c := range claims {
		if string(c.MatchKey) == key {
			return c
		}
	}
	return claim.Claim{}
}

type engineDifference struct {
	element, line, mac, vendor string
}

// fact is one engine difference that an expected difference implies, in
// the engine's vocabulary: its element label (CLP04, CAS*CO*253,
// CAS*CO*253*QTY, MOA, RAW*K3), the line control number it is reported
// with, and the MAC and vendor values ("" when absent).
type fact struct {
	where   string // the fixture's claim level or line, so facts from one place merge
	line    string
	element string
	mac     string
	vendor  string
}

// compareClaimDifferences maps a matched claim's expected differences to
// engine facts and matches them with the engine's differences.
func compareClaimDifferences(expected []expectedDifference, mac, vendor claim.Claim, engine []engineDifference) ([]mismatch, error) {
	macShape, err := shapeOf(mac)
	if err != nil {
		return nil, err
	}
	vendorShape, err := shapeOf(vendor)
	if err != nil {
		return nil, err
	}
	var out []mismatch
	var facts []fact
	seen := map[fact]bool{}
	for _, d := range expected {
		fs, notCompared, err := impliedFacts(d, macShape, vendorShape)
		if err != nil {
			out = append(out, mismatch{"fixture entry", fmt.Sprintf("%s: %v", d, err)})
			continue
		}
		for _, what := range notCompared {
			out = append(out, mismatch{"not compared " + what, d.String()})
		}
		for _, f := range fs {
			if !seen[f] {
				seen[f] = true
				facts = append(facts, f)
			}
		}
	}

	used := make([]bool, len(engine))
	matched := make([]bool, len(facts))
	for i, f := range facts {
		for j, e := range engine {
			if !used[j] && e.element == f.element && e.line == f.line && e.mac == f.mac && e.vendor == f.vendor {
				used[j], matched[i] = true, true
				break
			}
		}
	}
	for i, f := range facts {
		if matched[i] {
			continue
		}
		for j, e := range engine {
			if !used[j] && e.element == f.element && e.line == f.line {
				used[j], matched[i] = true, true
				out = append(out, mismatch{"value " + f.element,
					fmt.Sprintf("line %q: fixture implies (%q, %q), engine (%q, %q)", f.line, f.mac, f.vendor, e.mac, e.vendor)})
				break
			}
		}
	}
	for i, f := range facts {
		if !matched[i] {
			out = append(out, mismatch{"unreported " + f.element, fmt.Sprintf("line %q: fixture implies (%q, %q)", f.line, f.mac, f.vendor)})
		}
	}
	for j, e := range engine {
		if !used[j] {
			out = append(out, mismatch{"unexpected " + e.element, fmt.Sprintf("line %q: engine (%q, %q)", e.line, e.mac, e.vendor)})
		}
	}
	return out, nil
}

// segment is one segment as the fixture test reads it: the tag, then the
// elements, split on the element separator, as the file writes them.
type segment struct {
	elements []string // elements[0] is the tag
}

func (s segment) tag() string {
	return s.elements[0]
}

// element returns element n (1-based), or "" when the segment ends before it.
func (s segment) element(n int) string {
	if n < 1 || n >= len(s.elements) {
		return ""
	}
	return s.elements[n]
}

// segmentsOf splits a claim's bytes on the segment terminator. Whitespace
// around a segment is dropped and a blank segment is skipped, as the typed
// parser does.
func segmentsOf(data []byte, d Delimiters) []segment {
	var out []segment
	for _, raw := range bytes.Split(data, []byte{d.Segment}) {
		raw = bytes.TrimSpace(raw)
		if len(raw) == 0 {
			continue
		}
		out = append(out, segment{elements: strings.Split(string(raw), string(d.Element))})
	}
	return out
}

// claimShape is a claim's segments split into the claim level (CLP up to
// the first SVC) and its service lines (each SVC and the segments after it).
type claimShape struct {
	claim      []segment
	lines      [][]segment
	delimiters Delimiters
}

func shapeOf(c claim.Claim) (claimShape, error) {
	d, err := ParseDelimiters(c.Delimiters)
	if err != nil {
		return claimShape{}, err
	}
	s := claimShape{delimiters: d}
	for _, seg := range segmentsOf(c.Segments, d) {
		if seg.tag() == "SVC" {
			s.lines = append(s.lines, nil)
		}
		if len(s.lines) == 0 {
			s.claim = append(s.claim, seg)
			continue
		}
		s.lines[len(s.lines)-1] = append(s.lines[len(s.lines)-1], seg)
	}
	return s, nil
}

// text is a segment as the file writes it, without the terminator.
func (s claimShape) text(seg segment) string {
	return strings.Join(seg.elements, string(s.delimiters.Element))
}

// scope is the claim level or one service line of one side: its segments
// and the level the engine compares them at.
type scope struct {
	level string // LevelClaim or LevelLine
	segs  []segment
}

// scope returns the claim level or one line (1-based) of the claim; a line
// the side does not have has no segments.
func (s claimShape) scope(level string, line int) scope {
	if level == LevelClaim {
		return scope{level: LevelClaim, segs: s.claim}
	}
	if line < 1 || line > len(s.lines) {
		return scope{level: LevelLine}
	}
	return scope{level: LevelLine, segs: s.lines[line-1]}
}

// index is the position of the k-th segment (1-based) with the tag, or -1.
func (sc scope) index(tag string, k int) int {
	for i, s := range sc.segs {
		if s.tag() != tag {
			continue
		}
		k--
		if k == 0 {
			return i
		}
	}
	return -1
}

// lineControl is the line's REF*6R as the model reads it: the first one
// with a value.
func (sc scope) lineControl() string {
	for _, s := range sc.segs {
		if s.tag() == "REF" && s.element(1) == "6R" && s.element(2) != "" {
			return s.element(2)
		}
	}
	return ""
}

// The NM1 roles and the DTM qualifiers the model holds in fields of their
// own (x12835 claim.go, service_line.go).
var nm1Roles = map[string]bool{"QC": true, "74": true, "82": true, "TT": true, "PR": true}
var dateFields = map[string]map[string]bool{
	LevelClaim: {"232": true, "233": true, "050": true, "036": true},
	LevelLine:  {"472": true, "150": true, "151": true},
}

// placed reports whether the model gives the i-th segment of the scope a
// field of its own (x12835 parse.go: ClaimCLP.add, ServiceLineSVC.add).
// CLP, SVC, CAS, REF, AMT, and QTY always have one, and a line's LQ. At
// the claim level, the first NM1 of a held role, the first MIA, the first
// MOA, the first PER with a function code, and the first DTM of a held
// qualifier that carries a date have one; in a line, the first DTM of a
// line qualifier with a date. Anything else is kept raw and compared whole
// under RAW*<tag>.
func (sc scope) placed(i int) bool {
	s := sc.segs[i]
	switch s.tag() {
	case "CLP", "SVC", "CAS", "REF", "AMT", "QTY":
		return true
	case "LQ":
		return sc.level == LevelLine
	case "NM1":
		if sc.level != LevelClaim || !nm1Roles[s.element(1)] {
			return false
		}
	case "MIA", "MOA", "PER":
		if sc.level != LevelClaim {
			return false
		}
	case "DTM":
		if !dateFields[sc.level][s.element(1)] {
			return false
		}
	default:
		return false
	}
	key := placementKey(s)
	if key == "" {
		return false
	}
	for k := 0; k < i; k++ {
		if placementKey(sc.segs[k]) == key {
			return false
		}
	}
	return true
}

// placementKey is the field a segment competes for: an NM1 by its role, a
// DTM with a date by its qualifier, the MIA, the MOA, or the PER with a
// function code. The first segment with a key in a scope gets the field.
func placementKey(s segment) string {
	switch s.tag() {
	case "NM1":
		return "NM1*" + s.element(1)
	case "DTM":
		if s.element(2) == "" {
			return ""
		}
		return "DTM*" + s.element(1)
	case "MIA", "MOA":
		return s.tag()
	case "PER":
		if s.element(1) == "" {
			return ""
		}
		return "PER"
	}
	return ""
}

// rawValue is an unplaced segment as the engine compares it: its elements
// after the tag, without trailing empty ones, joined by "*".
func rawValue(s segment) string {
	parts := s.elements[1:]
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, "*")
}

// rawValues lists the unplaced segments with the tag, in file order.
func (sc scope) rawValues(tag string) []string {
	var out []string
	for i, s := range sc.segs {
		if s.tag() == tag && !sc.placed(i) {
			out = append(out, rawValue(s))
		}
	}
	return out
}

// rawRepeat counts the unplaced segments with the same tag before index.
func (sc scope) rawRepeat(index int) int {
	n := 0
	for i := 0; i < index; i++ {
		if sc.segs[i].tag() == sc.segs[index].tag() && !sc.placed(i) {
			n++
		}
	}
	return n
}

// casTriple is one (reason, amount, quantity) group of a CAS segment.
type casTriple struct {
	group, reason, amount, quantity string
}

func (c casTriple) label() string {
	return "CAS*" + c.group + "*" + c.reason
}

// casTripleAt is the t-th triple (0-based) of a CAS segment, if it has one.
func casTripleAt(s segment, t int) (casTriple, bool) {
	reason := s.element(2 + 3*t)
	if reason == "" {
		return casTriple{}, false
	}
	return casTriple{group: s.element(1), reason: reason, amount: s.element(3 + 3*t), quantity: s.element(4 + 3*t)}, true
}

// labeledValue is one value a segment gives a label, with its CAS triple
// (0 for other segments).
type labeledValue struct {
	triple int
	value  string
}

// valuesIn lists the values one placed segment gives a label, as the engine
// reads them: CAS*<group>*<reason> is the amount of each triple with that
// group and reason, and CAS*<group>*<reason>*QTY its quantity; NM1*<role>
// is NM109; AMT, DTM, QTY, and REF with that qualifier give their second
// element.
func valuesIn(s segment, label string) []labeledValue {
	tag, rest, _ := strings.Cut(label, "*")
	if s.tag() != tag {
		return nil
	}
	switch tag {
	case "CAS":
		group, after, _ := strings.Cut(rest, "*")
		reason, kind, _ := strings.Cut(after, "*")
		if s.element(1) != group {
			return nil
		}
		var out []labeledValue
		for t := 0; t < 6; t++ {
			tr, ok := casTripleAt(s, t)
			if !ok || tr.reason != reason {
				continue
			}
			value := tr.amount
			if kind == "QTY" {
				value = tr.quantity
			}
			out = append(out, labeledValue{t, value})
		}
		return out
	case "NM1":
		if s.element(1) == rest {
			return []labeledValue{{0, s.element(9)}}
		}
	default:
		if s.element(1) == rest {
			return []labeledValue{{0, s.element(2)}}
		}
	}
	return nil
}

// values lists a label's values in the scope, in file order, from the
// segments the model places.
func (sc scope) values(label string) []string {
	var out []string
	for i, s := range sc.segs {
		if !sc.placed(i) {
			continue
		}
		for _, lv := range valuesIn(s, label) {
			out = append(out, lv.value)
		}
	}
	return out
}

// repeat is how many of the label's values in the scope come before the
// one in the segment at index (triple t for CAS). The engine pairs the n-th
// value of a label on one side with the n-th on the other.
func (sc scope) repeat(index, t int, label string) int {
	count := 0
	for i, s := range sc.segs {
		if !sc.placed(i) {
			continue
		}
		for _, lv := range valuesIn(s, label) {
			if i == index && lv.triple == t {
				return count
			}
			count++
		}
	}
	return count
}

var remarkElements = map[string]map[int]bool{
	"MOA": {3: true, 4: true, 5: true, 6: true, 7: true},
	"MIA": {5: true, 20: true, 21: true, 22: true, 23: true},
}

// remarkList is the engine's value for a remark list: the codes of every
// placed segment in the scope that carries the label, in file order, joined
// by one space. The labels are MOA, MIA, LQ (LQ*HE), and LQ*<qualifier>.
func (sc scope) remarkList(label string) string {
	var codes []string
	for i, s := range sc.segs {
		if !sc.placed(i) {
			continue
		}
		switch {
		case (label == "MOA" || label == "MIA") && s.tag() == label:
			for e := 1; e < len(s.elements); e++ {
				if remarkElements[label][e] && s.element(e) != "" {
					codes = append(codes, s.element(e))
				}
			}
		case s.tag() == "LQ" && lqLabel(s) == label && s.element(2) != "":
			codes = append(codes, s.element(2))
		}
	}
	return strings.Join(codes, " ")
}

func lqLabel(s segment) string {
	if q := s.element(1); q != "HE" {
		return "LQ*" + q
	}
	return "LQ"
}

// numericLabel reports whether the engine writes the label's value as an
// exact decimal (80.00 reads "80"): the amounts, counts, weights, and
// rates of CLP, SVC, MIA, and MOA, and every CAS, AMT, and QTY value.
func numericLabel(label string) bool {
	switch label {
	case "CLP03", "CLP04", "CLP05", "CLP12", "CLP13", "SVC02", "SVC03", "SVC05", "SVC07",
		"MOA01", "MOA02", "MOA08", "MOA09":
		return true
	case "MIA05", "MIA20", "MIA21", "MIA22", "MIA23":
		return false
	}
	switch {
	case strings.HasPrefix(label, "CAS*"), strings.HasPrefix(label, "AMT*"), strings.HasPrefix(label, "QTY*"):
		return true
	case strings.HasPrefix(label, "MIA") && len(label) == 5:
		return true
	}
	return false
}

// decimalText is a numeric value as the engine writes it: parsed as an
// exact decimal and written without trailing fractional zeros. Text that
// is not a decimal comes back as it is (the engine fails such a file).
func decimalText(s string) string {
	if s == "" {
		return ""
	}
	v, err := x12835.ParseDecimal([]byte(s))
	if err != nil {
		return s
	}
	return v.String()
}

// composite writes an SVC01 or SVC06 composite as the engine does: the
// components joined by ":", without trailing empty ones.
func composite(value string, component byte) string {
	parts := strings.Split(value, string(component))
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, ":")
}

// impliedFacts maps one expected difference to the engine differences it
// implies, reading the segments it names from the two claims. It first
// checks that the fixture's values are the files' values there.
//
// The engine labels a value by element (CLP04, SVC03) or by segment and
// qualifier (AMT*B6, NM1*TT, CAS*CO*253, CAS*CO*253*QTY), pairs the n-th
// value of a qualified label on one side with the n-th on the other, and
// writes a numeric value as an exact decimal. The fixture pairs segments
// by position within their tag and holds the file's text. A segment the
// model does not place (placed) is compared whole as RAW*<tag>, the n-th
// such segment with the n-th on the other side. Set 1 (compareSet1) leaves
// some elements out; the mapping reports them as not compared. The
// mapping:
//
//	CLP n, SVC n     CLPnn, SVCnn with the fixture's values. A whole SVC on
//	                 one side implies SVCnn for each element it has. CLP07
//	                 is never listed. Left out: CLP03, CLP10, CLP14.
//	CAS              CAS*<group>*<reason> for the triple the element is in
//	                 (every triple for CAS01 or a whole segment), on each
//	                 side, with the value the engine pairs it with on the
//	                 other side; the quantity (CAS04, CAS07, ...) the same
//	                 way under *QTY.
//	AMT, DTM, QTY, REF   <tag>*<qualifier> the same way, valued by the second
//	                 element. Elements after the second: not compared.
//	                 Left out: REF*F8 and DTM*232 at the claim level,
//	                 REF*6R and every QTY in a line.
//	NM1              NM1*<role> the same way, valued by NM109. NM102 to
//	                 NM108 (names, ID qualifier): not compared. An NM1 on
//	                 one side without NM109 is invisible to the engine.
//	                 Left out: NM1*QC's NM109, the match key's MBI.
//	MOA, MIA         a remark code element (MOA03-07, MIA05, MIA20-23)
//	                 implies the joined code list, labeled MOA or MIA; any
//	                 other element n is MOAnn or MIAnn.
//	LQ               LQ (for LQ*HE) or LQ*<qualifier>: the joined code list.
//	PER              left out at the claim level.
//	kind order       not compared: the engine pairs lines by REF*6R, then
//	                 position, and claims by key, and compares neither order.
//	kind resolved_target  not compared (the engine does not read REF*F8).
//
// A fact whose two values are equal implies nothing: the engine pairs
// by label, so a change of position alone is not a difference to it. A
// fact's line is the MAC line's REF*6R when the MAC has the line, else
// the vendor line's, as the engine reports it; the claim level is "".
func impliedFacts(d expectedDifference, mac, vendor claimShape) ([]fact, []string, error) {
	switch {
	case d.Kind == "order" && d.Scope == LevelLine:
		return nil, []string{"line order"}, nil
	case d.Kind == "order":
		return nil, []string{"claim order"}, nil
	case d.Kind == "resolved_target":
		return nil, []string{"REF*F8 target"}, nil
	case d.Kind != "":
		return nil, nil, fmt.Errorf("unknown kind %q", d.Kind)
	case d.Scope != LevelClaim && d.Scope != LevelLine:
		return nil, nil, fmt.Errorf("scope %q belongs in file_differences", d.Scope)
	}

	macScope := mac.scope(d.Scope, d.Line)
	vendorScope := vendor.scope(d.Scope, d.VendorLine)
	place := fact{where: fmt.Sprintf("%s %d %d", d.Scope, d.Line, d.VendorLine)}
	if d.Scope == LevelLine {
		if d.Line > 0 {
			place.line = macScope.lineControl()
		} else {
			place.line = vendorScope.lineControl()
		}
	}

	// The segment on each side. Element 0 means the segment is on one
	// side only; otherwise both sides have it.
	macIndex, vendorIndex := macScope.index(d.Segment, d.Occurrence), vendorScope.index(d.Segment, d.Occurrence)
	haveMAC, haveVendor := macIndex >= 0, vendorIndex >= 0
	var macSeg, vendorSeg segment
	if haveMAC {
		macSeg = macScope.segs[macIndex]
	}
	if haveVendor {
		vendorSeg = vendorScope.segs[vendorIndex]
	}
	n := d.Element
	if n == 0 {
		haveMAC, haveVendor = haveMAC && d.MAC != "", haveVendor && d.Vendor != ""
		if haveMAC == haveVendor {
			return nil, nil, fmt.Errorf("element 0 needs its segment on exactly one side")
		}
		if haveMAC && mac.text(macSeg) != d.MAC || haveVendor && vendor.text(vendorSeg) != d.Vendor {
			return nil, nil, fmt.Errorf("segment text is not the file's")
		}
	} else {
		if !haveMAC || !haveVendor {
			return nil, nil, fmt.Errorf("segment not found on both sides")
		}
		if macSeg.element(n) != d.MAC || vendorSeg.element(n) != d.Vendor {
			return nil, nil, fmt.Errorf("values are not the file's (%q, %q)", macSeg.element(n), vendorSeg.element(n))
		}
	}

	var facts []fact
	var notCompared []string
	add := func(element, macValue, vendorValue string) {
		if numericLabel(element) {
			macValue, vendorValue = decimalText(macValue), decimalText(vendorValue)
		}
		if macValue == vendorValue {
			return
		}
		f := place
		f.element, f.mac, f.vendor = element, macValue, vendorValue
		facts = append(facts, f)
	}
	type side struct {
		index        int
		seg          segment
		onMAC        bool
		scope, other scope
	}
	// raw adds the fact for a side whose segment the model keeps raw: the
	// whole segment, paired with the other side's raw segment of the same
	// tag at the same repeat ("" when the other side has fewer).
	raw := func(s side) {
		other := s.other.rawValues(s.seg.tag())
		r := s.scope.rawRepeat(s.index)
		otherValue := ""
		if r < len(other) {
			otherValue = other[r]
		}
		if s.onMAC {
			add("RAW*"+s.seg.tag(), rawValue(s.seg), otherValue)
		} else {
			add("RAW*"+s.seg.tag(), otherValue, rawValue(s.seg))
		}
	}
	// sides are the sides whose segment the model places; the mapping
	// below reads those. A raw side is done here.
	var sides []side
	if haveMAC {
		s := side{macIndex, macSeg, true, macScope, vendorScope}
		if macScope.placed(macIndex) {
			sides = append(sides, s)
		} else {
			raw(s)
		}
	}
	if haveVendor {
		s := side{vendorIndex, vendorSeg, false, vendorScope, macScope}
		if vendorScope.placed(vendorIndex) {
			sides = append(sides, s)
		} else {
			raw(s)
		}
	}
	if len(sides) == 0 {
		return facts, notCompared, nil
	}
	// keyed adds a qualified label seen on one side, paired as the engine
	// pairs it: with the same label's value at the same repeat on the
	// other side ("" when the other side has fewer).
	keyed := func(s side, label, value string, triple int) {
		other := s.other.values(label)
		r := s.scope.repeat(s.index, triple, label)
		otherValue := ""
		if r < len(other) {
			otherValue = other[r]
		}
		if s.onMAC {
			add(label, value, otherValue)
		} else {
			add(label, otherValue, value)
		}
	}
	// oneSided adds an element of a segment on one side only.
	oneSided := func(s side, element, value string) {
		if s.onMAC {
			add(element, value, "")
		} else {
			add(element, "", value)
		}
	}

	switch d.Segment {
	case "CLP":
		if n == 0 || n == 7 {
			return nil, nil, fmt.Errorf("CLP element %d is never a listed difference", n)
		}
		if n == 3 || n == 10 || n == 14 {
			notCompared = append(notCompared, fmt.Sprintf("CLP%02d", n))
			break
		}
		add(fmt.Sprintf("CLP%02d", n), d.MAC, d.Vendor)

	case "SVC":
		if n > 0 {
			add(fmt.Sprintf("SVC%02d", n), d.MAC, d.Vendor)
			break
		}
		s := sides[0]
		component := vendor.delimiters.Component
		if s.onMAC {
			component = mac.delimiters.Component
		}
		for e := 1; e < len(s.seg.elements); e++ {
			v := s.seg.element(e)
			if e == 1 || e == 6 {
				v = composite(v, component)
			}
			oneSided(s, fmt.Sprintf("SVC%02d", e), v)
		}

	case "CAS":
		t, part := (n-2)/3, (n-2)%3
		for _, s := range sides {
			for tt := 0; tt < 6; tt++ {
				if n >= 2 && tt != t {
					continue
				}
				tr, ok := casTripleAt(s.seg, tt)
				if !ok {
					continue
				}
				if n < 2 || part != 2 {
					keyed(s, tr.label(), tr.amount, tt)
				}
				if n < 2 || part != 1 {
					keyed(s, tr.label()+"*QTY", tr.quantity, tt)
				}
			}
		}

	case "AMT", "DTM", "QTY", "REF":
		if n > 2 {
			notCompared = append(notCompared, fmt.Sprintf("%s%02d", d.Segment, n))
			break
		}
		for _, s := range sides {
			label := d.Segment + "*" + s.seg.element(1)
			switch {
			case d.Scope == LevelClaim && label == "REF*F8":
				notCompared = append(notCompared, "REF*F8")
			case d.Scope == LevelClaim && label == "DTM*232":
				notCompared = append(notCompared, "DTM*232")
			case d.Scope == LevelLine && label == "REF*6R":
				notCompared = append(notCompared, "REF*6R")
			case d.Scope == LevelLine && d.Segment == "QTY":
				notCompared = append(notCompared, "QTY (line)")
			default:
				keyed(s, label, s.seg.element(2), 0)
			}
		}

	case "NM1":
		if n != 0 && n != 1 && n != 9 {
			notCompared = append(notCompared, "NM1 name or ID qualifier")
			break
		}
		for _, s := range sides {
			switch {
			case s.seg.element(9) == "":
				if n != 9 {
					notCompared = append(notCompared, "NM1 without NM109")
				}
			case s.seg.element(1) == "QC":
				notCompared = append(notCompared, "NM1*QC")
			default:
				keyed(s, "NM1*"+s.seg.element(1), s.seg.element(9), 0)
			}
		}

	case "MOA", "MIA":
		remarks := remarkElements[d.Segment]
		if n > 0 && !remarks[n] {
			add(fmt.Sprintf("%s%02d", d.Segment, n), d.MAC, d.Vendor)
			break
		}
		add(d.Segment, macScope.remarkList(d.Segment), vendorScope.remarkList(d.Segment))
		if n == 0 {
			s := sides[0]
			for e := 1; e < len(s.seg.elements); e++ {
				if !remarks[e] && s.seg.element(e) != "" {
					oneSided(s, fmt.Sprintf("%s%02d", d.Segment, e), s.seg.element(e))
				}
			}
		}

	case "LQ":
		for _, s := range sides {
			label := lqLabel(s.seg)
			add(label, macScope.remarkList(label), vendorScope.remarkList(label))
		}

	case "PER":
		notCompared = append(notCompared, "PER")

	default:
		return nil, nil, fmt.Errorf("no mapping for a placed %s", d.Segment)
	}
	return facts, distinct(notCompared), nil
}

// distinct keeps the first of each repeated string, in order.
func distinct(in []string) []string {
	var out []string
	for _, s := range in {
		dup := false
		for _, o := range out {
			if o == s {
				dup = true
			}
		}
		if !dup {
			out = append(out, s)
		}
	}
	return out
}
