// Package migrate is vdb migrate: status, up, and down on the schema
// migrations embedded in lib/sys/db, one action per file (status.go,
// up.go, down.go) and the migrator they share in migrator.go.
package migrate

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/help"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// Run runs migrate status, up, or down. args[0] is the action; the flags
// follow it. The command's name alone, or help after it, is handled by
// main before Run is called.
func Run(ctx context.Context, args []string, s settings.Settings) error {
	action := args[0]
	switch action {
	case "status":
		return runStatus(ctx, args[1:], s)
	case "up":
		return runUp(ctx, args[1:], s)
	case "down":
		return runDown(ctx, args[1:], s)
	}
	return help.UsageError(fmt.Sprintf("migrate: unknown command %q", action), Help())
}

// options are the flags every action takes. status accepts -yes and
// -dry-run too and ignores them, so the three actions read alike.
type options struct {
	local  bool
	yes    bool
	dryRun bool
}

// parse reads an action's flags. A help flag prints the help and returns
// done; an unexpected argument is a usage error.
func parse(action string, args []string) (opts options, done bool, err error) {
	flags := flag.NewFlagSet("migrate "+action, flag.ContinueOnError)
	flags.Usage = func() { fmt.Print(Help()) }
	flags.BoolVar(&opts.local, "local", false, help.Local)
	flags.BoolVar(&opts.yes, "yes", false, "skip the typed confirmation (up, down)")
	flags.BoolVar(&opts.dryRun, "dry-run", false, "say what would be applied or rolled back, and stop (up, down)")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return opts, true, nil
		}
		return opts, false, help.ErrUsage
	}
	if flags.NArg() > 0 {
		return opts, false, help.UsageError(fmt.Sprintf("migrate %s: unexpected argument %q", action, flags.Arg(0)), Help())
	}
	return opts, false, nil
}

// open reads the embedded migrations, opens the database named by s (or
// the local one), names it, and reads what is applied. The caller closes
// the database.
func open(ctx context.Context, s settings.Settings, opts options) (migrator, db.DB, error) {
	if opts.local {
		s = s.Local()
	}
	migrations, err := readMigrations()
	if err != nil {
		return migrator{}, db.DB{}, err
	}
	database, err := db.Open(ctx, s.DB)
	if err != nil {
		return migrator{}, db.DB{}, err
	}
	fmt.Println("database", settings.Describe(s.DB))

	if _, err := database.Exec(ctx, createMigrationsTableSQL); err != nil {
		database.Close()
		return migrator{}, db.DB{}, fmt.Errorf("creating migrations table: %w", err)
	}
	m := migrator{db: database, target: settings.Describe(s.DB), migrations: migrations}
	m.applied, err = m.appliedChecksums(ctx)
	if err != nil {
		database.Close()
		return migrator{}, db.DB{}, err
	}
	return m, database, nil
}

// Help is the help for vdb migrate and each of its actions.
func Help() string {
	return help.Render("vdb migrate: apply, roll back, and report the schema migrations", []help.Section{
		{Heading: "DESCRIPTION", Body: `  The migrations are the files embedded in vdb (lib/sys/db/migrations), so vdb
  and the validator service always agree on the schema. The migrations table
  records each applied file's version and checksum: up refuses an applied file
  whose checksum has changed, and status reports such a file as changed. down
  rolls back only the latest migration, which deletes its tables.
`},
		{Heading: "DEFAULTS", Body: `  The database comes from the APP_SERVER_DB_* settings; -local uses the local
  development database (localhost:5432/validator, TLS off). up and down ask for
  a typed yes unless -yes is given.
`},
		{Heading: "OPTIONS", Body: `  -local     the local development database (localhost:5432/validator, TLS off)
  -yes       skip the typed confirmation (up, down)
  -dry-run   say what would be applied or rolled back, and stop (up, down)
`},
		{Heading: "USAGE", Body: `  vdb migrate status [-local]
  vdb migrate up [-local] [-yes] [-dry-run]
  vdb migrate down [-local] [-yes] [-dry-run]
`},
		{Heading: "EXAMPLES", Body: `  vdb migrate status -local           # each migration: applied, pending, or applied but changed since
  vdb migrate up -local -yes          # apply the pending migrations without the prompt
  vdb migrate down -local -dry-run    # which migration a down would roll back, on which database
  vdb migrate down -local -yes        # roll back the latest migration; this deletes its tables
`},
	})
}
