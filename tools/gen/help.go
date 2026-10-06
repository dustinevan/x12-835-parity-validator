// Help for gen, printed straight to stdout, never through a pager. Prose
// first (DESCRIPTION, HOW IT WORKS, DEFAULTS), then OPTIONS, USAGE, and
// EXAMPLES last, so the bottom of the help, nearest the prompt, is what the
// reader types from.
package main

import (
	"os"
	"strings"
)

// The colors for section headings: a few ANSI sequences, used only when
// stdout is a terminal.
const (
	ansiBold  = "\033[1m"
	ansiCyan  = "\033[36m"
	ansiReset = "\033[0m"
)

// isHelp reports whether arg asks for help.
func isHelp(arg string) bool {
	switch arg {
	case "help", "-help", "--help", "-h":
		return true
	}
	return false
}

// stdoutIsTerminal reports whether stdout is a terminal (a character
// device), so headings are colored only where the colors will be seen.
func stdoutIsTerminal() bool {
	info, err := os.Stdout.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// section is one block of help: a heading and its body, already indented
// and ending in a newline.
type section struct {
	heading string
	body    string
}

// renderHelp joins a title line and its sections, coloring the headings
// when stdout is a terminal.
func renderHelp(title string, sections []section) string {
	color := stdoutIsTerminal()
	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n")
	for _, s := range sections {
		b.WriteString("\n")
		if color {
			b.WriteString(ansiBold + ansiCyan + s.heading + ansiReset)
		} else {
			b.WriteString(s.heading)
		}
		b.WriteString("\n")
		b.WriteString(s.body)
	}
	return b.String()
}

func helpText() string {
	return renderHelp("gen: synthetic X12 835 test data for the validator", []section{
		{"DESCRIPTION", `  Writes a MAC reference set, vendor submissions for several dates, and
  truth.jsonl, one line per vendor claim listing its outcome, every
  difference the comparison should find between it and the MAC's claim, and
  the pattern behind each. Everything is synthetic: made-up payees,
  patients, identifiers, and amounts, nothing from a real remittance. The
  one real thing is the payer identifier (N1*PR N104): the MAC files are
  written under CMS's published contractor numbers for one state's Part A
  and Part B workload in each of the twelve A/B jurisdictions and the four
  DME MACs, so jurisdiction and Part read correctly downstream. The same
  flags and seed always produce byte-identical files. A run first removes
  what an earlier run wrote under the output directory (mac, each vendor's
  directory, truth.jsonl), so stale files never mix in.
`},
		{"HOW IT WORKS", `  1. A world of -claims MAC claims is drawn from -seed: professional, DME,
     and institutional claims (inpatient with MIA and DRG, outpatient with
     MOA, and the other bill types) in Medicare's proportions, spread over
     the sixteen jurisdictions by their published workload shares, with
     dates of service in October and November 2025, paid on weekly cycles,
     and written as one 835 file per workload per cycle under <out>/mac.
     The statuses include processed as primary and forwarded, processed as
     secondary, denied, reversals with their replacement and void chains
     (CLP09 1, 7, and 8 on institutional claims), and a few claims sent to
     the wrong contractor. A handful of keys repeat on the MAC side.
  2. For each vendor and submission date, the MAC claims are copied with
     differences planted from a catalog of difference patterns: ways a
     replacement adjudication system's output can differ from the MAC's,
     each applied to the claims where it would show (vendor.go lists the
     catalog with its taxonomy entries). The vendor-1 profile plants
     mostly pricing patterns, vendor-2 mostly edit and history patterns,
     both convention, timing, and front-end patterns; each later
     submission drops some. A few claims are missing, extra, or sent
     twice. The copies are written under <out>/<vendor>/<date>.
  3. truth.jsonl records each vendor claim's outcome (matched, missing,
     extra, set_aside, duplicated), its differences, and the pattern and
     taxonomy entry behind each difference. Two dollar amounts within five
     cents of each other are the same amount to the comparison, so a
     planted cent change (sequestration rounding) is in the files but not
     among the differences.
`},
		{"DEFAULTS", `  -out data/generated   -claims 20000   -seed 42
  -vendors vendor-1,vendor-2   -submissions 2026-10-02,2026-10-16,2026-10-30
`},
		{"OPTIONS", `  -out DIR            the output directory
  -claims N           MAC claims to generate (at least 60; the default set takes a few seconds)
  -seed N             the random seed
  -vendors V,V        vendor directory names, comma-separated
  -submissions D,D    vendor submission dates, YYYY-MM-DD, comma-separated
  -dry-run            list the files gen would write, and where, without writing
`},
		{"USAGE", `  gen -out DIR [-claims N] [-seed N] [-vendors V,V] [-submissions D,D] [-dry-run]
  gen alone prints this help; it runs only when given an option.
`},
		{"EXAMPLES", `  gen -out data/generated                    # the standard test data (what make gen writes)
  gen -out data/generated -dry-run           # what it would write, without writing
  gen -out /tmp/small -claims 1000 -seed 7   # a small set with another seed
`},
	})
}
