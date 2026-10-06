// Command codegen updates the entries tables in backend/pkg/carc/carc.go
// and backend/pkg/rarc/rarc.go from the JSON lists beside them,
// backend/pkg/carc.json and backend/pkg/rarc.json. Both Go files are written
// by hand except for their tables' rows, which sit between two marker
// comments:
//
//	var entries = []Entry{
//		// codegen: entries begin
//		{Code: "1", Short: "Deductible Amount", Description: "Deductible Amount", Status: "current", Start: "01/01/1995"},
//		// codegen: entries end
//	}
//
// codegen rewrites exactly the lines between the markers: one row per code,
// in the JSON's order and in the same format every time, so a refreshed
// list changes only the rows of the codes that changed, and a rerun on
// unchanged lists writes nothing. Everything outside the markers is left as
// it is. A file whose markers are missing, repeated, or out of order stops
// the run before anything is written. Run it after refreshing either JSON
// list from the x12.org page its _meta names:
//
//	codegen -root .            # at the repository root (make codes)
//	codegen -root . -dry-run   # say what would change; write nothing
//
// The paths under the root are fixed; only the root is a flag. codegen
// alone, or codegen help, prints the help (help.go) before anything is
// parsed.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// The two comments that bound each file's entries table. Each must be a
// line of its own, once in the file, the begin marker first.
const (
	beginMarker = "// codegen: entries begin"
	endMarker   = "// codegen: entries end"
)

// codeEntry is one row of backend/pkg/carc.json or rarc.json.
type codeEntry struct {
	Code         string `json:"code"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Start        string `json:"start"`
	Stop         string `json:"stop"`
	LastModified string `json:"lastModified"`
	Notes        string `json:"notes"`
}

// codeFile is the shape of both backend/pkg/carc.json and rarc.json; of the
// _meta object at the top, codegen reads only the count.
type codeFile struct {
	Meta struct {
		Count int `json:"count"`
	} `json:"_meta"`
	Codes []codeEntry `json:"codes"`
}

// list is one code list: its JSON and the Go file whose table it fills.
type list struct {
	jsonPath string
	goPath   string
}

// plan is what a run does to one Go file: the file as it is, the file with
// the new rows, and which codes the new rows add, change, or remove.
type plan struct {
	list    list
	codes   int
	current []byte
	next    []byte
	added   []string
	changed []string
	removed []string
}

func main() {
	// Help comes before the flags are parsed, so typing codegen alone or
	// codegen help never writes.
	if len(os.Args) == 1 || isHelp(os.Args[1]) {
		fmt.Print(helpText())
		return
	}
	flags := flag.NewFlagSet("codegen", flag.ContinueOnError)
	flags.Usage = func() { fmt.Print(helpText()) }
	root := flags.String("root", ".", "the repository root, where backend/pkg lives")
	dryRun := flags.Bool("dry-run", false, "say what would change and write nothing")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(2)
	}
	if flags.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "codegen: unknown command %q\n\n", flags.Arg(0))
		fmt.Print(helpText())
		os.Exit(2)
	}

	lists := []list{
		{
			jsonPath: filepath.Join(*root, "backend", "pkg", "carc.json"),
			goPath:   filepath.Join(*root, "backend", "pkg", "carc", "carc.go"),
		},
		{
			jsonPath: filepath.Join(*root, "backend", "pkg", "rarc.json"),
			goPath:   filepath.Join(*root, "backend", "pkg", "rarc", "rarc.go"),
		},
	}

	// Both files are checked and their new contents built before either is
	// written, so a problem in one stops the run with nothing changed.
	var plans []plan
	var errs []error
	for _, l := range lists {
		p, err := planList(l)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		plans = append(plans, p)
	}
	if err := errors.Join(errs...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "codegen: nothing written")
		os.Exit(1)
	}

	for _, p := range plans {
		if bytes.Equal(p.current, p.next) {
			fmt.Printf("codegen: %s: %d codes from %s, already up to date; nothing written\n", p.list.goPath, p.codes, p.list.jsonPath)
			continue
		}
		what := fmt.Sprintf("%d added, %d changed, %d removed", len(p.added), len(p.changed), len(p.removed))
		if len(p.added)+len(p.changed)+len(p.removed) == 0 {
			what = "no code added, changed, or removed, but rows reordered or reformatted"
		}
		if *dryRun {
			fmt.Printf("codegen: %s: %d codes from %s: %s; not written (-dry-run)\n", p.list.goPath, p.codes, p.list.jsonPath, what)
		} else {
			if err := writeFile(p.list.goPath, p.next); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("codegen: %s: %d codes from %s: %s; written\n", p.list.goPath, p.codes, p.list.jsonPath, what)
		}
		printCodes("added", p.added)
		printCodes("changed", p.changed)
		printCodes("removed", p.removed)
	}
}

// planList reads one JSON list and its Go file and builds the file with the
// list's rows between the markers. The error names the file when the list
// fails a check or the markers are missing, repeated, or out of order.
func planList(l list) (plan, error) {
	codes, err := readList(l.jsonPath)
	if err != nil {
		return plan{}, err
	}
	current, err := os.ReadFile(l.goPath)
	if err != nil {
		return plan{}, fmt.Errorf("codegen: reading %s: %w", l.goPath, err)
	}
	lines := strings.Split(string(current), "\n")
	begin, end, err := findMarkers(l.goPath, lines)
	if err != nil {
		return plan{}, err
	}

	// Rows take the begin marker's indentation, which gofmt gives the
	// elements of the table it sits in.
	indent := lines[begin][:len(lines[begin])-len(strings.TrimLeft(lines[begin], " \t"))]
	rows := make([]string, len(codes))
	for i, e := range codes {
		rows[i] = indent + row(e)
	}

	next := make([]string, 0, begin+1+len(rows)+len(lines)-end)
	next = append(next, lines[:begin+1]...)
	next = append(next, rows...)
	next = append(next, lines[end:]...)

	p := plan{
		list:    l,
		codes:   len(codes),
		current: current,
		next:    []byte(strings.Join(next, "\n")),
	}
	p.added, p.changed, p.removed = compareRows(lines[begin+1:end], codes, rows)
	return p, nil
}

var codePattern = regexp.MustCompile(`^[A-Za-z0-9]+$`)

// readList reads one JSON list and checks it: every code is letters and
// digits, no code appears twice, and the list holds as many codes as its
// _meta count says.
func readList(path string) ([]codeEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("codegen: reading %s: %w", path, err)
	}
	var file codeFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return nil, fmt.Errorf("codegen: parsing %s: %w", path, err)
	}
	seen := make(map[string]bool, len(file.Codes))
	for _, e := range file.Codes {
		if !codePattern.MatchString(e.Code) {
			return nil, fmt.Errorf("codegen: checking %s: code %q has characters other than letters and digits", path, e.Code)
		}
		if seen[e.Code] {
			return nil, fmt.Errorf("codegen: checking %s: code %q appears more than once", path, e.Code)
		}
		seen[e.Code] = true
	}
	if file.Meta.Count != len(file.Codes) {
		return nil, fmt.Errorf("codegen: checking %s: _meta.count is %d but the list holds %d codes", path, file.Meta.Count, len(file.Codes))
	}
	return file.Codes, nil
}

// findMarkers returns the indexes of the begin and end marker lines. Each
// marker must be on exactly one line, and the begin marker must come first.
func findMarkers(path string, lines []string) (begin, end int, err error) {
	var begins, ends []int
	for i, line := range lines {
		switch strings.TrimSpace(line) {
		case beginMarker:
			begins = append(begins, i)
		case endMarker:
			ends = append(ends, i)
		}
	}
	if len(begins) != 1 || len(ends) != 1 {
		return 0, 0, fmt.Errorf("codegen: finding the entries markers in %s: %q is on %s and %q on %s; each must be on exactly one line",
			path, beginMarker, lineNumbers(begins), endMarker, lineNumbers(ends))
	}
	if ends[0] < begins[0] {
		return 0, 0, fmt.Errorf("codegen: finding the entries markers in %s: %q (line %d) comes before %q (line %d)",
			path, endMarker, ends[0]+1, beginMarker, begins[0]+1)
	}
	return begins[0], ends[0], nil
}

// lineNumbers says where a marker was found, numbering lines from 1:
// "no line", "line 12", or "lines 12, 40".
func lineNumbers(at []int) string {
	switch len(at) {
	case 0:
		return "no line"
	case 1:
		return fmt.Sprintf("line %d", at[0]+1)
	}
	numbers := make([]string, len(at))
	for i, n := range at {
		numbers[i] = strconv.Itoa(n + 1)
	}
	return "lines " + strings.Join(numbers, ", ")
}

// row is one table row for e: Code, Short, Description, and Status always,
// then Start, Stop, LastModified, and Notes when the list gives them, each
// value a quoted Go string. The format never varies, so an unchanged code
// always gets the same line.
func row(e codeEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "{Code: %q, Short: %q, Description: %q, Status: %q",
		e.Code, shortDescription(e.Description), e.Description, e.Status)
	if e.Start != "" {
		fmt.Fprintf(&b, ", Start: %q", e.Start)
	}
	if e.Stop != "" {
		fmt.Fprintf(&b, ", Stop: %q", e.Stop)
	}
	if e.LastModified != "" {
		fmt.Fprintf(&b, ", LastModified: %q", e.LastModified)
	}
	if e.Notes != "" {
		fmt.Fprintf(&b, ", Notes: %q", e.Notes)
	}
	b.WriteString("},")
	return b.String()
}

// shortDescription cuts desc at the first ". " or ";", whichever comes
// first, verbatim and with no paraphrase. A description with neither is
// kept whole.
func shortDescription(desc string) string {
	cut := len(desc)
	if i := strings.Index(desc, ". "); i != -1 && i < cut {
		cut = i
	}
	if i := strings.Index(desc, ";"); i != -1 && i < cut {
		cut = i
	}
	return strings.TrimSpace(desc[:cut])
}

// compareRows compares the rows now between the markers with the new rows,
// code by code: a code only in the new rows is added, a code only in the
// old rows is removed, and a code in both whose row differs is changed. An
// old line that does not begin with a quoted code is not counted.
func compareRows(old []string, codes []codeEntry, rows []string) (added, changed, removed []string) {
	oldRows := make(map[string]string, len(old))
	var oldCodes []string
	for _, line := range old {
		code, ok := rowCode(line)
		if !ok {
			continue
		}
		oldRows[code] = strings.TrimSpace(line)
		oldCodes = append(oldCodes, code)
	}
	newCodes := make(map[string]bool, len(codes))
	for i, e := range codes {
		newCodes[e.Code] = true
		was, ok := oldRows[e.Code]
		switch {
		case !ok:
			added = append(added, e.Code)
		case was != strings.TrimSpace(rows[i]):
			changed = append(changed, e.Code)
		}
	}
	for _, code := range oldCodes {
		if !newCodes[code] {
			removed = append(removed, code)
		}
	}
	return added, changed, removed
}

// rowCode reads the code from a row that row wrote: the quoted string after
// "{Code: ".
func rowCode(line string) (string, bool) {
	rest, ok := strings.CutPrefix(strings.TrimSpace(line), "{Code: ")
	if !ok {
		return "", false
	}
	quoted, err := strconv.QuotedPrefix(rest)
	if err != nil {
		return "", false
	}
	code, err := strconv.Unquote(quoted)
	if err != nil {
		return "", false
	}
	return code, true
}

// printCodes prints the codes of one kind of change on one line: up to
// twenty, then how many more, so filling an empty table stays readable.
func printCodes(kind string, codes []string) {
	const most = 20
	switch {
	case len(codes) == 0:
	case len(codes) <= most:
		fmt.Printf("  %s: %s\n", kind, strings.Join(codes, " "))
	default:
		fmt.Printf("  %s: %s and %d more\n", kind, strings.Join(codes[:most], " "), len(codes)-most)
	}
}

// writeFile replaces path with data through a temporary file in the same
// directory and a rename, keeping the file's permissions, so an interrupted
// run leaves the old file or the new one and never part of either.
func writeFile(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("codegen: writing %s: %w", path, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*")
	if err != nil {
		return fmt.Errorf("codegen: writing %s: %w", path, err)
	}
	defer os.Remove(tmp.Name()) // fails quietly once the rename has moved it
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("codegen: writing %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("codegen: writing %s: %w", path, err)
	}
	if err := os.Chmod(tmp.Name(), info.Mode().Perm()); err != nil {
		return fmt.Errorf("codegen: writing %s: %w", path, err)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("codegen: writing %s: %w", path, err)
	}
	return nil
}
