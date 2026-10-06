// Command gen writes synthetic X12 835 test data: MAC reference files
// under real MAC workload numbers, vendor files for several submissions
// whose differences come from a catalog of difference patterns (vendor.go),
// and a truth file listing every difference it put between them and the
// pattern behind each. See docs/testdata/REALISM.md.
//
// The same flags and seed always produce byte-identical output. gen alone
// prints its help (help.go); it runs only when given an option.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	Out         string
	Claims      int
	Seed        int64
	Submissions []string
	Vendors     []string
	DryRun      bool // list what would be written, and where, without writing
}

func main() {
	// Help comes before the flags are parsed, so typing gen never writes.
	if len(os.Args) == 1 || isHelp(os.Args[1]) {
		fmt.Print(helpText())
		return
	}
	var cfg config
	var submissions, vendors string
	flags := flag.NewFlagSet("gen", flag.ContinueOnError)
	flags.Usage = func() { fmt.Print(helpText()) }
	flags.StringVar(&cfg.Out, "out", "data/generated", "output directory")
	flags.IntVar(&cfg.Claims, "claims", 20000, "MAC claims to generate")
	flags.Int64Var(&cfg.Seed, "seed", 42, "random seed")
	flags.StringVar(&submissions, "submissions", "2026-10-02,2026-10-16,2026-10-30", "vendor submission dates, comma-separated")
	flags.StringVar(&vendors, "vendors", "vendor-1,vendor-2", "vendor names, comma-separated")
	flags.BoolVar(&cfg.DryRun, "dry-run", false, "list the files gen would write, and where, without writing")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(2)
	}
	if flags.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "gen: unknown command %q\n\n", flags.Arg(0))
		fmt.Print(helpText())
		os.Exit(2)
	}
	cfg.Submissions = strings.Split(submissions, ",")
	cfg.Vendors = strings.Split(vendors, ",")

	summary, err := run(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(summary)
}

// outputFile is one file gen writes: its path under the output directory
// and its bytes.
type outputFile struct {
	Path string
	Data []byte
}

// run generates everything, writes it under cfg.Out (or, with DryRun, lists
// what it would write there), and returns a summary.
func run(cfg config) (string, error) {
	if cfg.Claims < 60 {
		return "", fmt.Errorf("gen: checking flags: -claims must be at least 60")
	}
	for _, d := range cfg.Submissions {
		if _, err := time.Parse("2006-01-02", d); err != nil {
			return "", fmt.Errorf("gen: parsing submission date %q: %w", d, err)
		}
	}
	for _, v := range cfg.Vendors {
		if v == "" || v == "mac" || strings.ContainsAny(v, `/\.`) {
			return "", fmt.Errorf("gen: checking flags: vendor name %q is not a plain directory name", v)
		}
	}

	w := generateWorld(cfg)

	var outputs []outputFile
	for _, f := range macFiles(w) {
		outputs = append(outputs, outputFile{Path: f.Path, Data: render(w.Workloads, f)})
	}
	macCount := len(outputs)

	var summary strings.Builder
	fmt.Fprintf(&summary, "MAC: %d claims in %d payments, %d files, %d workloads\n", len(w.Claims), len(w.Payments), macCount, len(w.Workloads))
	fmt.Fprintf(&summary, "%-26s %6s %6s %7s %7s %6s %9s %10s  %s\n", "submission", "files", "equal", "matched", "missing", "extra", "set-aside", "duplicated", "claims per pattern")

	var truth bytes.Buffer
	enc := json.NewEncoder(&truth)
	enc.SetEscapeHTML(false)
	for v := range cfg.Vendors {
		for d := range cfg.Submissions {
			sub := buildSubmission(cfg, w, v, d)
			for _, f := range sub.Files {
				outputs = append(outputs, outputFile{Path: f.Path, Data: render(w.Workloads, f)})
			}
			for _, r := range sub.Truth {
				if err := enc.Encode(r); err != nil {
					return "", fmt.Errorf("gen: encoding truth: %w", err)
				}
			}
			summary.WriteString(summarize(sub))
		}
	}
	outputs = append(outputs, outputFile{Path: "truth.jsonl", Data: truth.Bytes()})

	// Remove only what a previous run wrote, so stale files never mix in.
	stale := []string{"mac", "truth.jsonl"}
	stale = append(stale, cfg.Vendors...)
	if cfg.DryRun {
		fmt.Fprintf(&summary, "dry run: would remove %s under %s, then write these %d files there:\n", strings.Join(stale, ", "), cfg.Out, len(outputs))
		for _, f := range outputs {
			fmt.Fprintf(&summary, "  %s  (%d bytes)\n", f.Path, len(f.Data))
		}
		summary.WriteString("dry run: nothing written\n")
		return summary.String(), nil
	}
	for _, s := range stale {
		if err := os.RemoveAll(filepath.Join(cfg.Out, s)); err != nil {
			return "", fmt.Errorf("gen: removing old output: %w", err)
		}
	}
	for _, f := range outputs {
		if err := writeFile(cfg.Out, f.Path, f.Data); err != nil {
			return "", err
		}
	}
	fmt.Fprintf(&summary, "wrote %d files under %s\n", len(outputs), cfg.Out)
	return summary.String(), nil
}

func writeFile(root, path string, data []byte) error {
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("gen: creating directory for %s: %w", path, err)
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		return fmt.Errorf("gen: writing %s: %w", path, err)
	}
	return nil
}

// summarize reports one submission. The equal share is matched claims with no
// difference, over reference claims that are not set aside (so missing
// claims count against it). Patterns are listed by their catalog number
// with the claims they touched.
func summarize(sub submission) string {
	var equal, matched, missing, extra, setAside, duplicated int
	byPattern := map[string]int{}
	for _, r := range sub.Truth {
		switch r.Outcome {
		case "matched":
			matched++
			if len(r.Differences) == 0 {
				equal++
			}
		case "missing":
			missing++
		case "extra":
			extra++
		case "set_aside":
			setAside++
		case "duplicated":
			duplicated++
		}
		for _, s := range r.Patterns {
			byPattern[s]++
		}
	}
	var counts []string
	for n, s := range catalog {
		if c := byPattern[s.Name]; c > 0 {
			counts = append(counts, fmt.Sprintf("%d:%d", n, c))
		}
	}
	rate := 100 * float64(equal) / float64(max(1, matched+missing))
	return fmt.Sprintf("%-26s %6d %5.1f%% %7d %7d %6d %9d %10d  %s\n", sub.Name, len(sub.Files), rate, matched, missing, extra, setAside, duplicated, strings.Join(counts, " "))
}
