package migrate

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
)

// runDown is vdb migrate down: roll back the most recently applied
// migration only, after a typed yes unless -yes is given. -dry-run names
// that migration and stops.
func runDown(ctx context.Context, args []string, s settings.Settings) error {
	opts, done, err := parse("down", args)
	if err != nil || done {
		return err
	}
	m, database, err := open(ctx, s, opts)
	if err != nil {
		return err
	}
	defer database.Close()
	return m.down(ctx, opts.yes, opts.dryRun)
}

const removeMigrationSQL = `DELETE FROM public.migrations WHERE version = @version`

// down rolls back the most recently applied migration only. With dryRun it
// names that migration and stops.
func (m migrator) down(ctx context.Context, yes, dryRun bool) error {
	var last migration
	for _, mg := range m.migrations {
		if _, ok := m.applied[mg.version]; ok && mg.version > last.version {
			last = mg
		}
	}
	if last.version == 0 {
		fmt.Println("No migrations are applied.")
		return nil
	}
	if last.down == "" {
		return fmt.Errorf("migration %d has no .down.sql file", last.version)
	}
	if dryRun {
		fmt.Printf("dry run: would roll back migration %d (%s) on %s, which deletes its tables; nothing changed\n", last.version, last.description, m.target)
		return nil
	}
	if !yes {
		ok, err := confirm(ctx, m.target, fmt.Sprintf("ROLL BACK migration %d (%s); this can delete data", last.version, last.description))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("nothing changed")
			return nil
		}
	}

	err := pgx.BeginFunc(ctx, m.db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, last.down); err != nil {
			return fmt.Errorf("rolling back migration %d: %w", last.version, err)
		}
		if _, err := tx.Exec(ctx, removeMigrationSQL, pgx.NamedArgs{"version": last.version}); err != nil {
			return fmt.Errorf("removing migration record %d: %w", last.version, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("rolled back %3d  %s\n", last.version, last.description)
	return nil
}
