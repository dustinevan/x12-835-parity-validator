package macs

import "slices"

// The tree: every number in the table under its jurisdiction and Part, the
// way CMS's numbering groups them (CR 5651, "New Numbers for All MAC
// Jurisdictions", 2007). Tree places each row by its digits and the
// table's Part column:
//
//   - A DME number (Part "DME") stands alone under its jurisdiction's DME
//     Part.
//   - A home health and hospice (HHH) number, fifth digit 4, sits under its
//     Part (the table gives A for each) as its own kind.
//   - A roll-up is a number whose third digit is 0 in Part A or B. Its Part
//     comes from the table, not from its fifth digit, because CR 8303
//     numbered the JK Part A roll-up for New Hampshire and Vermont 14013.
//   - A workload, third digit not 0 and fifth digit 1 or 2, sits under the
//     roll-up of its jurisdiction and Part whose states include all of its
//     own, and directly under its Part when no roll-up does or when it has
//     no states.
//   - The J6 internal processing number 06014, and any number the rules
//     above do not place, sits under its Part with a note.
//
// Within each kind the numbers keep the table's order. The notes restate
// what the table's group comments say.

// Jurisdiction is one MAC jurisdiction in the tree.
type Jurisdiction struct {
	Name       string // "J5" through "JN" for A/B MACs, "JA" through "JD" for DME MACs
	Contractor string
	Parts      []Part // A then B for an A/B jurisdiction; DME alone for a DME jurisdiction
}

// Part is one claim Part of a jurisdiction and its numbers by kind.
type Part struct {
	Name      string   // "A", "B", "DME", or "" when the table gives none
	RollUps   []RollUp // each with the workloads under it
	Workloads []Number // the workloads under no roll-up; a DME jurisdiction's number
	HHH       []Number // home health and hospice numbers, each with its HHH region as the note
	Other     []Number // numbers the numbering does not place, each with a note
}

// RollUp is a roll-up number and the workloads under it.
type RollUp struct {
	Number
	CarriesHHH bool // the table's comment says the roll-up also carries the HHH claims
	Workloads  []Number
}

// Number is one five-digit number, its states with their names, and a note
// restating the table's comment ("" when there is none).
type Number struct {
	ID     string
	States []State
	Note   string
}

// State is a USPS code and the name it stands for.
type State struct {
	Code string
	Name string
}

// order is CMS's order of the jurisdictions: the A/B MACs, then the DME MACs.
// A jurisdiction missing from it follows these, in table order.
var order = []string{"J5", "J6", "J8", "J15", "JE", "JF", "JH", "JJ", "JK", "JL", "JM", "JN", "JA", "JB", "JC", "JD"}

// partOrder is the order of the Parts within a jurisdiction; a Part the
// table does not name ("") comes last.
var partOrder = []string{"A", "B", "DME"}

// internal lists the numbers CMS calls internal processing numbers. Each
// sits under its Part with the others the numbering does not place.
var internal = []string{"06014"}

// carriesHHH lists the roll-ups whose group comment says they also carry
// the HHH claims.
var carriesHHH = []string{"06001", "11001", "14011", "15001"}

// unplaced is the note on a number the rules do not place and the notes
// below do not explain.
const unplaced = "Its digits and Part fit no kind in CMS's numbering."

// notes restates, per number, what the table's group comments say about it:
// an HHH number's region, an internal number's use, a legacy workload, a
// state split between workloads, and a roll-up whose digits differ from
// the numbering.
var notes = map[string]string{
	"06001": "Its HHH claims are those of one Region D area: WI, MN, MI, NY, NJ, PR, and VI.",
	"14013": "A Part A roll-up, though its last digit is 3 (CR 8303).",
	"06004": "Region D",
	"11004": "Region C",
	"14014": "Region A",
	"15004": "Region B",
	"06014": "Internal processing number for the HHH claims of the other Region D area.",
	"05901": "WPS legacy workload, formerly Mutual of Omaha; CR 8059 names no states.",
	"01911": "J1 WPS legacy workload; CR 8226 names no states.",
	"04911": "Legacy Mutual of Omaha providers.",
	"12901": "Legacy WPS providers, formerly Mutual of Omaha.",
	"01112": "Northern California",
	"01182": "Southern California",
	"13202": "Downstate New York",
	"13282": "Upstate New York",
	"13292": "Queens County",
}

// stateNames spells out the USPS codes: the states, the District of
// Columbia, and the territories.
var stateNames = map[string]string{
	"AK": "Alaska",
	"AL": "Alabama",
	"AR": "Arkansas",
	"AS": "American Samoa",
	"AZ": "Arizona",
	"CA": "California",
	"CO": "Colorado",
	"CT": "Connecticut",
	"DC": "District of Columbia",
	"DE": "Delaware",
	"FL": "Florida",
	"GA": "Georgia",
	"GU": "Guam",
	"HI": "Hawaii",
	"IA": "Iowa",
	"ID": "Idaho",
	"IL": "Illinois",
	"IN": "Indiana",
	"KS": "Kansas",
	"KY": "Kentucky",
	"LA": "Louisiana",
	"MA": "Massachusetts",
	"MD": "Maryland",
	"ME": "Maine",
	"MI": "Michigan",
	"MN": "Minnesota",
	"MO": "Missouri",
	"MP": "Northern Mariana Islands",
	"MS": "Mississippi",
	"MT": "Montana",
	"NC": "North Carolina",
	"ND": "North Dakota",
	"NE": "Nebraska",
	"NH": "New Hampshire",
	"NJ": "New Jersey",
	"NM": "New Mexico",
	"NV": "Nevada",
	"NY": "New York",
	"OH": "Ohio",
	"OK": "Oklahoma",
	"OR": "Oregon",
	"PA": "Pennsylvania",
	"PR": "Puerto Rico",
	"RI": "Rhode Island",
	"SC": "South Carolina",
	"SD": "South Dakota",
	"TN": "Tennessee",
	"TX": "Texas",
	"UT": "Utah",
	"VA": "Virginia",
	"VI": "Virgin Islands",
	"VT": "Vermont",
	"WA": "Washington",
	"WI": "Wisconsin",
	"WV": "West Virginia",
	"WY": "Wyoming",
}

// stateName returns the name of a USPS code, or the code itself when the
// list does not hold it.
func stateName(code string) string {
	if name, ok := stateNames[code]; ok {
		return name
	}
	return code
}

// Tree returns every number in the table, each exactly once: the
// jurisdictions in CMS's order, each with its Parts. It builds a new value
// on each call; nothing in it is shared with the table.
func Tree() []Jurisdiction {
	var names []string
	for _, w := range Workloads {
		if !slices.Contains(names, w.Jurisdiction) {
			names = append(names, w.Jurisdiction)
		}
	}
	slices.SortStableFunc(names, func(a, b string) int {
		return rank(order, a) - rank(order, b)
	})
	out := make([]Jurisdiction, 0, len(names))
	for _, name := range names {
		out = append(out, jurisdiction(name))
	}
	return out
}

// rank is the position of s in list, or the list's length when it is not
// there, so an unlisted value sorts after the listed ones.
func rank(list []string, s string) int {
	if i := slices.Index(list, s); i >= 0 {
		return i
	}
	return len(list)
}

// jurisdiction builds one jurisdiction from its rows. Its contractor is the
// first row's; tree_test.go checks that every row agrees.
func jurisdiction(name string) Jurisdiction {
	j := Jurisdiction{Name: name}
	var parts []string
	for _, w := range Workloads {
		if w.Jurisdiction != name {
			continue
		}
		if j.Contractor == "" {
			j.Contractor = w.Contractor
		}
		if !slices.Contains(parts, w.Part) {
			parts = append(parts, w.Part)
		}
	}
	slices.SortStableFunc(parts, func(a, b string) int {
		return rank(partOrder, a) - rank(partOrder, b)
	})
	for _, p := range parts {
		j.Parts = append(j.Parts, part(name, p))
	}
	return j
}

// kind is where a row sits in its Part.
type kind int

const (
	kindWorkload kind = iota
	kindRollUp
	kindHHH
	kindDME
	kindOther
)

// kindOf sorts a row by the rules at the top of this file.
func kindOf(w Workload) kind {
	switch {
	case slices.Contains(internal, w.ID) || len(w.ID) != 5:
		return kindOther
	case w.Part == "DME":
		return kindDME
	case w.ID[4] == '4':
		return kindHHH
	case w.ID[2] == '0' && (w.Part == "A" || w.Part == "B"):
		return kindRollUp
	case w.ID[2] != '0' && (w.ID[4] == '1' || w.ID[4] == '2'):
		return kindWorkload
	}
	return kindOther
}

// part builds one Part of a jurisdiction. The roll-ups are found first, so
// each workload can be placed under its own.
func part(jurisdiction, name string) Part {
	var rows []Workload
	for _, w := range Workloads {
		if w.Jurisdiction == jurisdiction && w.Part == name {
			rows = append(rows, w)
		}
	}
	p := Part{Name: name}
	var rollUpRows []Workload
	for _, w := range rows {
		if kindOf(w) == kindRollUp {
			rollUpRows = append(rollUpRows, w)
			p.RollUps = append(p.RollUps, RollUp{Number: number(w), CarriesHHH: slices.Contains(carriesHHH, w.ID)})
		}
	}
	for _, w := range rows {
		switch kindOf(w) {
		case kindRollUp:
			// Placed above.
		case kindWorkload:
			if i := rollUpOf(rollUpRows, w); i >= 0 {
				p.RollUps[i].Workloads = append(p.RollUps[i].Workloads, number(w))
			} else {
				p.Workloads = append(p.Workloads, number(w))
			}
		case kindDME:
			p.Workloads = append(p.Workloads, number(w))
		case kindHHH:
			p.HHH = append(p.HHH, number(w))
		default:
			n := number(w)
			if n.Note == "" {
				n.Note = unplaced
			}
			p.Other = append(p.Other, n)
		}
	}
	return p
}

// rollUpOf returns the index of the first roll-up whose states include all
// of w's, or -1 when none does or w has no states.
func rollUpOf(rollUps []Workload, w Workload) int {
	if len(w.States) == 0 {
		return -1
	}
	for i, r := range rollUps {
		if includesAll(r.States, w.States) {
			return i
		}
	}
	return -1
}

// includesAll reports whether every state in want is in have.
func includesAll(have, want []string) bool {
	for _, s := range want {
		if !slices.Contains(have, s) {
			return false
		}
	}
	return true
}

// number copies a row into the tree with its states spelled out and its
// note.
func number(w Workload) Number {
	states := make([]State, 0, len(w.States))
	for _, code := range w.States {
		states = append(states, State{Code: code, Name: stateName(code)})
	}
	return Number{ID: w.ID, States: states, Note: notes[w.ID]}
}
