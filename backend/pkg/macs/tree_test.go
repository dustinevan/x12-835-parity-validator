package macs

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// treeNumbers lists every number in the tree with where it sits:
// "JM A roll-up", "JM A under 11001", "JM A workload", "JM A HHH",
// "JM A other".
func treeNumbers(tree []Jurisdiction) map[string][]string {
	at := map[string][]string{}
	for _, j := range tree {
		for _, p := range j.Parts {
			where := j.Name + " " + p.Name
			for _, r := range p.RollUps {
				at[r.ID] = append(at[r.ID], where+" roll-up")
				for _, w := range r.Workloads {
					at[w.ID] = append(at[w.ID], where+" under "+r.ID)
				}
			}
			for _, w := range p.Workloads {
				at[w.ID] = append(at[w.ID], where+" workload")
			}
			for _, w := range p.HHH {
				at[w.ID] = append(at[w.ID], where+" HHH")
			}
			for _, w := range p.Other {
				at[w.ID] = append(at[w.ID], where+" other")
			}
		}
	}
	return at
}

// TestTreeHoldsEveryNumberOnce checks that each of the table's numbers
// appears exactly once in the tree, under its own jurisdiction and Part,
// and that the tree holds nothing else.
func TestTreeHoldsEveryNumberOnce(t *testing.T) {
	if len(Workloads) != 136 {
		t.Errorf("the table has %d rows, want 136 (update this test when the table changes)", len(Workloads))
	}
	at := treeNumbers(Tree())
	for _, w := range Workloads {
		places := at[w.ID]
		if len(places) != 1 {
			t.Errorf("%s appears %d times in the tree (%v), want once", w.ID, len(places), places)
			continue
		}
		if prefix := w.Jurisdiction + " " + w.Part + " "; !strings.HasPrefix(places[0], prefix) {
			t.Errorf("%s sits at %q, want under %q", w.ID, places[0], prefix)
		}
	}
	if len(at) != len(Workloads) {
		t.Errorf("the tree holds %d numbers, the table %d", len(at), len(Workloads))
	}
}

// TestTreeOrder checks the jurisdictions against CMS's order and the
// Parts against A, B, DME, and each jurisdiction's contractor against
// every one of its rows.
func TestTreeOrder(t *testing.T) {
	tree := Tree()
	var names []string
	for _, j := range tree {
		names = append(names, j.Name)
		var parts []string
		for _, p := range j.Parts {
			parts = append(parts, p.Name)
		}
		want := []string{"A", "B"}
		if slices.Contains(dme, j.Name) {
			want = []string{"DME"}
		}
		if !slices.Equal(parts, want) {
			t.Errorf("%s has Parts %v, want %v", j.Name, parts, want)
		}
		for _, w := range Workloads {
			if w.Jurisdiction == j.Name && w.Contractor != j.Contractor {
				t.Errorf("%s: contractor %q, but its jurisdiction %s shows %q", w.ID, w.Contractor, j.Name, j.Contractor)
			}
		}
	}
	want := []string{"J5", "J6", "J8", "J15", "JE", "JF", "JH", "JJ", "JK", "JL", "JM", "JN", "JA", "JB", "JC", "JD"}
	if !slices.Equal(names, want) {
		t.Errorf("Tree() jurisdictions = %v, want %v", names, want)
	}
}

// TestTreePlaces checks the kinds the rules give: HHH numbers only under
// Part A, every roll-up with at least one workload, no workload that fits
// two roll-ups, 06014 among the numbers with a note, and a few numbers
// whose places the CRs state.
func TestTreePlaces(t *testing.T) {
	at := treeNumbers(Tree())
	for id, places := range at {
		if strings.HasSuffix(places[0], " HHH") && !strings.Contains(places[0], " A HHH") {
			t.Errorf("HHH number %s sits at %q, want Part A", id, places[0])
		}
	}
	for _, j := range Tree() {
		for _, p := range j.Parts {
			var rollUps []Workload
			for _, r := range p.RollUps {
				if len(r.Workloads) == 0 {
					t.Errorf("roll-up %s has no workload under it", r.ID)
				}
				w, _ := Lookup(r.ID)
				rollUps = append(rollUps, w)
			}
			for _, w := range Workloads {
				if w.Jurisdiction != j.Name || w.Part != p.Name || kindOf(w) != kindWorkload || len(w.States) == 0 {
					continue
				}
				fits := 0
				for _, r := range rollUps {
					if includesAll(r.States, w.States) {
						fits++
					}
				}
				if fits > 1 {
					t.Errorf("workload %s fits %d roll-ups, want at most one", w.ID, fits)
				}
			}
			for _, n := range p.Other {
				if n.Note == "" {
					t.Errorf("%s sits among the numbers the numbering does not place, with no note", n.ID)
				}
			}
		}
	}
	tests := []struct {
		id    string
		place string
	}{
		{"06014", "J6 A other"},   // CR 11976's internal processing number
		{"14013", "JK A roll-up"}, // CR 8303: the FISS roll-up for NH and VT Part A
		{"14311", "JK A under 14013"},
		{"14511", "JK A under 14013"},
		{"14111", "JK A under 14011"},
		{"04911", "JH A workload"}, // CO, NM, OK, TX: 04011 lacks OK
		{"04311", "JH A workload"}, // no CR names a roll-up for Oklahoma Part A
		{"01211", "JE A under 01011"},
		{"05901", "J5 A workload"}, // no states
		{"15102", "J15 B under 15002"},
		{"10312", "JJ B under 10012"},
		{"16013", "JA DME workload"},
		{"06004", "J6 A HHH"},
	}
	for _, tt := range tests {
		if got := at[tt.id]; len(got) != 1 || got[0] != tt.place {
			t.Errorf("%s sits at %v, want %q", tt.id, got, tt.place)
		}
	}
}

// TestStateNames checks that every state the table uses has a name.
func TestStateNames(t *testing.T) {
	for _, w := range Workloads {
		for _, s := range w.States {
			if _, ok := stateNames[s]; !ok {
				t.Errorf("%s: state %q has no name", w.ID, s)
			}
		}
	}
	if len(stateNames) != 56 {
		t.Errorf("stateNames holds %d codes, want 56: 50 states, DC, AS, GU, MP, PR, VI", len(stateNames))
	}
}

// node is one line of a drawn tree and the lines under it.
type node struct {
	text     string
	children []node
}

// meaning is what each Part's claims are, as the drawing writes it.
var meaning = map[string]string{
	"A":   "institutional claims, FISS",
	"B":   "professional claims, MCS",
	"DME": "durable medical equipment",
}

// listed joins words as English does: "SC", "VA and WV", "AL, GA, and TN".
func listed(words []string) string {
	switch len(words) {
	case 0:
		return ""
	case 1:
		return words[0]
	case 2:
		return words[0] + " and " + words[1]
	}
	return strings.Join(words[:len(words)-1], ", ") + ", and " + words[len(words)-1]
}

// leaf is a number's line: its ID and its state names.
func leaf(n Number) string {
	var names []string
	for _, s := range n.States {
		names = append(names, s.Name)
	}
	return n.ID + "  " + strings.Join(names, ", ")
}

// drawingOf turns a jurisdiction into the drawing's nodes: a roll-up shows
// its states, and the HHH claims when it carries them, in quotes; a
// workload under no roll-up in a Part that has roll-ups says so.
func drawingOf(j Jurisdiction) node {
	root := node{text: j.Name + "  " + j.Contractor}
	for _, p := range j.Parts {
		pn := node{text: fmt.Sprintf("Part %s  (%s)", p.Name, meaning[p.Name])}
		for _, r := range p.RollUps {
			var scope []string
			for _, s := range r.States {
				scope = append(scope, s.Code)
			}
			if r.CarriesHHH {
				scope = append(scope, "the HHH claims")
			}
			rn := node{text: fmt.Sprintf("roll-up %s  %q", r.ID, listed(scope))}
			for _, w := range r.Workloads {
				rn.children = append(rn.children, node{text: leaf(w)})
			}
			pn.children = append(pn.children, rn)
		}
		for _, w := range p.Workloads {
			text := leaf(w)
			if len(p.RollUps) > 0 {
				text += "  (no roll-up named)"
			}
			pn.children = append(pn.children, node{text: text})
		}
		for _, w := range p.HHH {
			pn.children = append(pn.children, node{text: w.ID + "  home health and hospice, " + w.Note})
		}
		for _, w := range p.Other {
			pn.children = append(pn.children, node{text: w.ID + "  " + w.Note})
		}
		root.children = append(root.children, pn)
	}
	return root
}

// draw writes n's children under it with the drawing's branches.
func draw(b *strings.Builder, n node, indent string) {
	for i, c := range n.children {
		branch, next := "├── ", "│   "
		if i == len(n.children)-1 {
			branch, next = "└── ", "    "
		}
		b.WriteString(indent + branch + c.text + "\n")
		draw(b, c, indent+next)
	}
}

// jmDrawing is the JM subtree as the MAC map's drawing shows it, with the
// contractor's name as CMS prints it and two spaces before each note.
const jmDrawing = `JM  Palmetto GBA, LLC
├── Part A  (institutional claims, FISS)
│   ├── roll-up 11001  "SC and the HHH claims"
│   │   └── 11201  South Carolina
│   ├── roll-up 11003  "VA and WV"
│   │   ├── 11301  Virginia
│   │   └── 11401  West Virginia
│   ├── 11501  North Carolina  (no roll-up named)
│   └── 11004  home health and hospice, Region C
└── Part B  (professional claims, MCS)
    ├── 11202  South Carolina
    ├── 11302  Virginia
    ├── 11402  West Virginia
    └── 11502  North Carolina
`

// TestTreeJM checks the JM subtree against the drawing, line for line.
func TestTreeJM(t *testing.T) {
	var jm *Jurisdiction
	tree := Tree()
	for i := range tree {
		if tree[i].Name == "JM" {
			jm = &tree[i]
		}
	}
	if jm == nil {
		t.Fatal("Tree() has no JM")
	}
	root := drawingOf(*jm)
	var b strings.Builder
	b.WriteString(root.text + "\n")
	draw(&b, root, "")
	if got := b.String(); got != jmDrawing {
		t.Errorf("the JM subtree draws as\n%s\nwant\n%s", got, jmDrawing)
	}
}
