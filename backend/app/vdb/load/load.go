// Package load is vdb load: one submission's 835 files into the database.
package load

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/help"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
	submissionctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/submission"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/submission"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
)

// Run reads one submission's input files and loads them.
func Run(ctx context.Context, args []string, s settings.Settings) error {
	flags := flag.NewFlagSet("load", flag.ContinueOnError)
	flags.Usage = func() { fmt.Print(Help()) }
	local := flags.Bool("local", false, help.Local)
	dryRun := flags.Bool("dry-run", false, "name the database and list the files, then stop without loading")
	dir := flags.String("dir", "", "the directory of 835 files, or with S3 enabled a key prefix under APP_SERVER_S3_PREFIX")
	prefix := flags.String("prefix", "", "the submission prefix: mac, or vendor/date such as vendor-1/2026-10-02")
	side := flags.String("side", "", "mac or vendor")
	label := flags.String("label", "", "a display label: MAC, Vendor 1")
	date := flags.String("date", "", "the submission date, YYYY-MM-DD (vendor submissions)")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return help.ErrUsage
	}
	if flags.NArg() > 0 {
		return help.UsageError(fmt.Sprintf("load: unexpected argument %q", flags.Arg(0)), Help())
	}
	if *dir == "" || *prefix == "" || *side == "" || *label == "" {
		return help.UsageError("load needs -dir, -prefix, -side, and -label", Help())
	}
	if *side != submission.SideMAC && *side != submission.SideVendor {
		return help.UsageError(fmt.Sprintf("load: -side must be %s or %s", submission.SideMAC, submission.SideVendor), Help())
	}
	ns := submission.NewSubmission{Prefix: *prefix, Side: *side, Label: *label}
	if *date != "" {
		submittedOn, err := time.Parse("2006-01-02", *date)
		if err != nil {
			return fmt.Errorf("parsing -date: %w", err)
		}
		ns.SubmittedOn = submittedOn
	}
	if *local {
		s = s.Local()
	}
	if *dryRun {
		return dryRunLoad(ctx, s, *dir, ns)
	}

	database, err := settings.OpenDatabase(ctx, s)
	if err != nil {
		return err
	}
	defer database.Close()
	seeder, err := settings.OpenSeeder(ctx, database, s)
	if err != nil {
		return err
	}
	return seeder.Load(ctx, *dir, ns)
}

// dryRunLoad names the database a load would open, without opening it, and
// lists the files it would read: a local directory, or with S3 enabled the
// objects under the key prefix, which is only listed.
func dryRunLoad(ctx context.Context, s settings.Settings, dir string, ns submission.NewSubmission) error {
	fmt.Printf("database %s (dry run: not opened)\n", settings.Describe(s.DB))
	source := "directory " + dir
	var names []string
	if s.S3Enabled {
		inputs, err := sys.OpenS3(ctx, s.S3)
		if err != nil {
			return err
		}
		source = fmt.Sprintf("s3://%s/%s (read-only)", s.S3.Bucket, path.Join(s.S3.Prefix, dir))
		names, err = inputs.List(ctx, dir)
		if err != nil {
			return err
		}
	} else {
		var err error
		names, err = submissionctrl.ListDirectory(dir)
		if err != nil {
			return err
		}
	}
	fmt.Printf("dry run: load would read %s as submission %s (side %s, label %q):\n", source, ns.Prefix, ns.Side, ns.Label)
	for _, name := range names {
		fmt.Println(" ", name)
	}
	if len(names) == 0 {
		fmt.Println("dry run: no input files, so the load would fail; nothing loaded")
		return nil
	}
	fmt.Printf("dry run: %d files; nothing loaded\n", len(names))
	return nil
}

// Help is the help for vdb load.
func Help() string {
	return help.Render("vdb load: read one submission's 835 files into the database", []help.Section{
		{Heading: "DESCRIPTION", Body: `  Reads every file directly under -dir (not dot files, not subdirectories) as
  one submission and writes its claims and payments in one transaction,
  replacing everything an earlier load of the same prefix wrote, including
  its comparisons and summaries. No files is an error. With
  APP_SERVER_S3_ENABLED=true, -dir is a key prefix under APP_SERVER_S3_PREFIX,
  and S3 is only listed and read. The database must be at vdb's migration
  version.
`},
		{Heading: "DEFAULTS", Body: `  The database comes from the APP_SERVER_DB_* settings; -local uses the local
  development database (localhost:5432/validator, TLS off) with S3 off. -date
  is empty: the MAC submission has none.
`},
		{Heading: "OPTIONS", Body: `  -dir DIR           the directory of 835 files (with S3, the key prefix)
  -prefix PREFIX     the submission prefix: mac, or vendor/date such as vendor-1/2026-10-02
  -side mac|vendor   which side of the comparison the files are
  -label LABEL       a display label: MAC, Vendor 1
  -date YYYY-MM-DD   the submission date (vendor submissions)
  -local             the local development database (localhost:5432/validator, TLS off), S3 off
  -dry-run           name the database and list the files, then stop without loading
`},
		{Heading: "USAGE", Body: `  vdb load -dir DIR -prefix PREFIX -side mac|vendor -label LABEL [-date YYYY-MM-DD] [-local] [-dry-run]
`},
		{Heading: "EXAMPLES", Body: `  vdb load -dir data/generated/mac -prefix mac -side mac -label MAC -local   # the MAC reference submission
  vdb load -dir data/generated/vendor-1/2026-10-02 -prefix vendor-1/2026-10-02 \
      -side vendor -label "Vendor 1" -date 2026-10-02 -local                # one vendor submission; compare it next
  vdb load -dir data/generated/mac -prefix mac -side mac -label MAC -local -dry-run   # the database and files, nothing loaded
`},
	})
}
