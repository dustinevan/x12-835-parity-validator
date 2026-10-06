// Package seed is vdb seed: the chain over the generated test data, the
// MAC set, then load, compare, and summarize for every vendor submission.
package seed

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/help"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
)

// Run loads, compares, and summarizes the generated test data.
func Run(ctx context.Context, args []string, s settings.Settings) error {
	flags := flag.NewFlagSet("seed", flag.ContinueOnError)
	flags.Usage = func() { fmt.Print(Help()) }
	local := flags.Bool("local", false, help.Local)
	dir := flags.String("dir", "data/generated", "the generated data: DIR/mac and DIR/<vendor>/<date>")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return help.ErrUsage
	}
	if flags.NArg() > 0 {
		return help.UsageError(fmt.Sprintf("seed: unexpected argument %q", flags.Arg(0)), Help())
	}
	if *local {
		s = s.Local()
	}
	if s.S3Enabled {
		return fmt.Errorf("seed reads local generated data and does not run with S3 enabled")
	}

	database, err := settings.OpenDatabase(ctx, s)
	if err != nil {
		return err
	}
	defer database.Close()
	fmt.Println(settings.CompareLine(s))
	seeder, err := settings.OpenSeeder(ctx, database, s)
	if err != nil {
		return err
	}
	return seeder.Seed(ctx, *dir)
}

// Help is the help for vdb seed.
func Help() string {
	return help.Render("vdb seed: load, compare, and summarize the generated test data", []help.Section{
		{Heading: "DESCRIPTION", Body: `  Loads the MAC submission from DIR/mac and every vendor submission from
  DIR/<vendor>/<YYYY-MM-DD>, then compares and summarizes each vendor
  submission in prefix order. Labels come from the vendor directory name
  (vendor-1 becomes Vendor 1). It reads local files only and does not
  run with S3 enabled. Each stage replaces its own earlier results, so a rerun
  rebuilds everything.
`},
		{Heading: "DEFAULTS", Body: `  -dir data/generated, what make gen writes. The database comes from the
  APP_SERVER_DB_* settings; -local uses the local development database
  (localhost:5432/validator, TLS off).
`},
		{Heading: "OPTIONS", Body: `  -dir DIR   the generated data: DIR/mac and DIR/<vendor>/<date>
  -local     the local development database (localhost:5432/validator, TLS off)
`},
		{Heading: "USAGE", Body: `  vdb seed [-dir DIR] [-local]
`},
		{Heading: "EXAMPLES", Body: `  vdb seed -local                   # what make db runs after migrating
  vdb seed -dir /tmp/small -local   # another generated set
`},
	})
}
