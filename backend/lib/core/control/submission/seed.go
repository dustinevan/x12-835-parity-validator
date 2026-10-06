package submission

import (
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	submissiontable "github.com/dustinevan/x12-835-parity-validator/backend/lib/data/submission"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
)

// Seeder runs the pipeline stages through a Controller for vdb, reporting
// each step to out. Load reads a submission's input files from a local
// directory, or, for a Seeder built with NewS3Seeder, from a key prefix in
// S3, which it only lists and reads. Seed loads, compares, and summarizes
// the generated test data.
// It lives in this package, not its own, so that one control package does
// not import another.
type Seeder struct {
	submissions Controller
	fromS3      bool
	s3          sys.S3
	out         io.Writer
}

// NewSeeder returns a Seeder that reads input files from local directories.
func NewSeeder(submissions Controller, out io.Writer) Seeder {
	return Seeder{submissions: submissions, out: out}
}

// NewS3Seeder returns a Seeder that reads input files from S3.
func NewS3Seeder(submissions Controller, inputs sys.S3, out io.Writer) Seeder {
	return Seeder{submissions: submissions, fromS3: true, s3: inputs, out: out}
}

// Seed loads the MAC submission from <dir>/mac and every vendor submission
// from <dir>/<vendor>/<YYYY-MM-DD>, then compares and summarizes each
// vendor submission in prefix order (so by vendor, then date).
// Labels come from the vendor directory name: vendor-1 -> "Vendor 1".
// It reads the generated test data from the local filesystem only.
func (s Seeder) Seed(ctx context.Context, dir string) error {
	if s.fromS3 {
		return fmt.Errorf("submission: seeding: seed reads local generated data and does not run with S3 enabled")
	}
	mac := submissiontable.NewSubmission{Prefix: "mac", Side: submissiontable.SideMAC, Label: "MAC"}
	if err := s.Load(ctx, filepath.Join(dir, "mac"), mac); err != nil {
		return err
	}

	vendorDirs, err := filepath.Glob(filepath.Join(dir, "*", "????-??-??"))
	if err != nil {
		return fmt.Errorf("submission: seeding: finding vendor submissions: %w", err)
	}
	sort.Strings(vendorDirs)
	var prefixes []string
	for _, vd := range vendorDirs {
		dateText := filepath.Base(vd)
		vendor := filepath.Base(filepath.Dir(vd))
		submittedOn, err := time.Parse("2006-01-02", dateText)
		if err != nil {
			return fmt.Errorf("submission: seeding: reading submission date %s: %w", vd, err)
		}
		ns := submissiontable.NewSubmission{
			Prefix:      path.Join(vendor, dateText),
			Side:        submissiontable.SideVendor,
			Label:       labelFor(vendor),
			SubmittedOn: submittedOn,
		}
		if err := s.Load(ctx, vd, ns); err != nil {
			return err
		}
		prefixes = append(prefixes, ns.Prefix)
	}
	for _, p := range prefixes {
		if err := s.Compare(ctx, p); err != nil {
			return err
		}
		if err := s.Summarize(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

// Load reads every input file directly under dir and loads them as one
// submission. dir is a local directory, or for an S3 Seeder a key prefix
// relative to the configured S3 prefix. Finding no files is an error:
// S3 lists a mistyped prefix as empty, and loading nothing would replace the
// submission's claims with none.
func (s Seeder) Load(ctx context.Context, dir string, ns submissiontable.NewSubmission) error {
	files, err := s.read(ctx, dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("submission: loading %s: no input files in %s", ns.Prefix, dir)
	}
	loaded, err := s.submissions.Load(ctx, ns, files)
	if err != nil {
		return err
	}
	fmt.Fprintf(s.out, "loaded   %-24s %6d claims in %d payments from %d files, %d repeated keys\n", ns.Prefix, loaded.Claims, loaded.Payments, len(files), loaded.RepeatedKeys)
	return nil
}

// read returns the input files directly under dir, in name order: a local
// directory (ReadDirectory), or the objects under a key prefix in S3.
func (s Seeder) read(ctx context.Context, dir string) ([]InputFile, error) {
	if !s.fromS3 {
		return ReadDirectory(dir)
	}
	names, err := s.s3.List(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("submission: reading %s: %w", dir, err)
	}
	files := make([]InputFile, 0, len(names))
	for _, name := range names {
		data, err := s.s3.Read(ctx, path.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("submission: reading %s: %w", name, err)
		}
		files = append(files, InputFile{Name: name, Data: data})
	}
	return files, nil
}

// Compare compares one vendor submission and reports its counts.
func (s Seeder) Compare(ctx context.Context, prefix string) error {
	sum, err := s.submissions.Compare(ctx, prefix)
	if err != nil {
		return err
	}
	fmt.Fprintf(s.out, "compared %-24s equal %d, different %d, missing %d, duplicated %d, extra %d, set aside %d keys; %d differences\n",
		prefix, sum.Equal, sum.Different, sum.Missing, sum.Duplicated, sum.Extra, sum.SetAside, sum.Differences)
	return nil
}

// Summarize summarizes one compared vendor submission, which shows it in the
// views.
func (s Seeder) Summarize(ctx context.Context, prefix string) error {
	if err := s.submissions.Summarize(ctx, prefix); err != nil {
		return err
	}
	fmt.Fprintf(s.out, "summarized %s\n", prefix)
	return nil
}

func labelFor(dirName string) string {
	words := strings.Split(dirName, "-")
	for i, w := range words {
		if w != "" {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
