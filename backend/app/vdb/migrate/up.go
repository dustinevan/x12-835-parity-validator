package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
)

// runUp is vdb migrate up: apply every pending migration, each in its own
// transaction, after a typed yes unless -yes is given. -dry-run lists them
// and stops.
func runUp(ctx context.Context, args []string, s settings.Settings) error {
	opts, done, err := parse("up", args)
	if err != nil || done {
		return err
	}
	m, database, err := open(ctx, s, opts)
	if err != nil {
		return err
	}
	defer database.Close()
	return m.up(ctx, opts.yes, opts.dryRun)
}

const recordMigrationSQL = `
INSERT INTO public.migrations (version, description, checksum, applied_at)
VALUES (@version, @description, @checksum, @appliedAt)`

// up applies every pending migration, each in its own transaction. With
// dryRun it lists them and stops.
func (m migrator) up(ctx context.Context, yes, dryRun bool) error {
	var pending []migration
	for _, mg := range m.migrations {
		checksum, ok := m.applied[mg.version]
		if ok && checksum != mg.checksum() {
			return fmt.Errorf("migration %d was applied from a different file; roll it back before changing it", mg.version)
		}
		if !ok {
			pending = append(pending, mg)
		}
	}
	if len(pending) == 0 {
		fmt.Println("All migrations are applied.")
		return nil
	}

	for _, mg := range pending {
		fmt.Printf("pending %3d  %s\n", mg.version, mg.description)
	}
	if dryRun {
		fmt.Printf("dry run: would apply %d migration(s) to %s; nothing changed\n", len(pending), m.target)
		return nil
	}
	if !yes {
		ok, err := confirm(ctx, m.target, fmt.Sprintf("Apply %d migration(s)", len(pending)))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("nothing changed")
			return nil
		}
	}

	for _, mg := range pending {
		err := pgx.BeginFunc(ctx, m.db, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, mg.up); err != nil {
				return fmt.Errorf("applying migration %d: %w", mg.version, err)
			}
			_, err := tx.Exec(ctx, recordMigrationSQL, pgx.NamedArgs{
				"version":     mg.version,
				"description": mg.description,
				"checksum":    mg.checksum(),
				"appliedAt":   time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("recording migration %d: %w", mg.version, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Printf("applied %3d  %s\n", mg.version, mg.description)
	}
	return nil
}
