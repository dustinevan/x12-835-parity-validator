// Package compare is vdb compare: one vendor submission against the MAC
// submission.
package compare

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/help"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
	submissionctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/submission"
)

// Run compares one vendor submission with the MAC submission.
func Run(ctx context.Context, args []string, s settings.Settings) error {
	flags := flag.NewFlagSet("compare", flag.ContinueOnError)
	flags.Usage = func() { fmt.Print(Help()) }
	local := flags.Bool("local", false, help.Local)
	prefix := flags.String("submission", "", "the vendor submission prefix, such as vendor-1/2026-10-02")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return help.ErrUsage
	}
	if flags.NArg() > 0 {
		return help.UsageError(fmt.Sprintf("compare: unexpected argument %q", flags.Arg(0)), Help())
	}
	if *prefix == "" {
		return help.UsageError("compare needs -submission", Help())
	}
	if *local {
		s = s.Local()
	}

	database, err := settings.OpenDatabase(ctx, s)
	if err != nil {
		return err
	}
	defer database.Close()
	fmt.Println(settings.CompareLine(s))
	controller := submissionctrl.NewController(database, s.Set, s.Workloads)
	return submissionctrl.NewSeeder(controller, os.Stdout).Compare(ctx, *prefix)
}

// Help is the help for vdb compare.
func Help() string {
	return help.Render("vdb compare: compare one vendor submission with the MAC submission", []help.Section{
		{Heading: "DESCRIPTION", Body: `  Pairs the vendor submission's claims with the MAC's by match key, one MAC
  workload at a time (APP_SERVER_COMPARE_WORKLOADS of them at once), compares
  each pair under the compared set APP_SERVER_COMPARE_SET names, and records
  every comparison and difference, replacing an earlier run's. The submission
  is hidden from the views while it runs and until it is summarized; a failed
  workload sets its status to failed, and rerunning compare is safe.
`},
		{Heading: "DEFAULTS", Body: `  The database comes from the APP_SERVER_DB_* settings; -local uses the local
  development database (localhost:5432/validator, TLS off). Compared set 1, one
  workload at a time.
`},
		{Heading: "OPTIONS", Body: `  -submission PREFIX   the vendor submission, such as vendor-1/2026-10-02
  -local               the local development database (localhost:5432/validator, TLS off)
`},
		{Heading: "USAGE", Body: `  vdb compare -submission PREFIX [-local]
`},
		{Heading: "EXAMPLES", Body: `  vdb compare -submission vendor-1/2026-10-02 -local                                 # one submission; summarize it next
  APP_SERVER_COMPARE_WORKLOADS=4 vdb compare -submission vendor-1/2026-10-02 -local  # four workloads at once
`},
	})
}
