// Package summarize is vdb summarize: a compared submission's summary
// tables, which is what shows it in the views.
package summarize

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

// Run builds one compared vendor submission's summary tables.
func Run(ctx context.Context, args []string, s settings.Settings) error {
	flags := flag.NewFlagSet("summarize", flag.ContinueOnError)
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
		return help.UsageError(fmt.Sprintf("summarize: unexpected argument %q", flags.Arg(0)), Help())
	}
	if *prefix == "" {
		return help.UsageError("summarize needs -submission", Help())
	}
	if *local {
		s = s.Local()
	}

	database, err := settings.OpenDatabase(ctx, s)
	if err != nil {
		return err
	}
	defer database.Close()
	controller := submissionctrl.NewController(database, s.Set, s.Workloads)
	return submissionctrl.NewSeeder(controller, os.Stdout).Summarize(ctx, *prefix)
}

// Help is the help for vdb summarize.
func Help() string {
	return help.Render("vdb summarize: build a compared submission's summary tables, which shows it in the views", []help.Section{
		{Heading: "DESCRIPTION", Body: `  Rebuilds the measure, tag_total, cell, and tag_cell rows of a compared
  vendor submission from its comparison rows, each table in its own short
  transaction, and sets the status to summarized, which is when the views
  show it. A submission whose status is not compared or summarized is
  refused: compare it first.
`},
		{Heading: "DEFAULTS", Body: `  The database comes from the APP_SERVER_DB_* settings; -local uses the local
  development database (localhost:5432/validator, TLS off).
`},
		{Heading: "OPTIONS", Body: `  -submission PREFIX   the vendor submission, such as vendor-1/2026-10-02
  -local               the local development database (localhost:5432/validator, TLS off)
`},
		{Heading: "USAGE", Body: `  vdb summarize -submission PREFIX [-local]
`},
		{Heading: "EXAMPLES", Body: `  vdb summarize -submission vendor-1/2026-10-02 -local   # after compare; the views then show it
`},
	})
}
