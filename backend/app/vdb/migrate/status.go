package migrate

import (
	"context"
	"fmt"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
)

// runStatus is vdb migrate status: each migration as applied, pending, or
// applied from a file that has changed since.
func runStatus(ctx context.Context, args []string, s settings.Settings) error {
	opts, done, err := parse("status", args)
	if err != nil || done {
		return err
	}
	m, database, err := open(ctx, s, opts)
	if err != nil {
		return err
	}
	defer database.Close()
	return m.status()
}

func (m migrator) status() error {
	for _, mg := range m.migrations {
		state := "pending"
		if checksum, ok := m.applied[mg.version]; ok {
			state = "applied"
			if checksum != mg.checksum() {
				state = "applied, but the file has changed since"
			}
		}
		fmt.Printf("%3d  %-9s %s\n", mg.version, state, mg.description)
	}
	return nil
}
