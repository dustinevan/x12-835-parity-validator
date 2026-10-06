package migrate

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

const createMigrationsTableSQL = `
CREATE TABLE IF NOT EXISTS public.migrations (
    version      INTEGER PRIMARY KEY,
    description  TEXT NOT NULL,
    checksum     TEXT NOT NULL,
    applied_at   TIMESTAMPTZ NOT NULL
)`

type migration struct {
	version     int
	description string
	up          string
	down        string
}

// checksum identifies the exact up SQL that was applied, so an edited
// migration is noticed rather than silently skipped.
func (m migration) checksum() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.TrimSpace(m.up))))
}

// migrator applies and rolls back migrations on one database. target
// names the database (settings.Describe) in confirmations and dry runs.
type migrator struct {
	db         db.DB
	target     string
	migrations []migration
	applied    map[int]string // version -> checksum
}

const appliedChecksumsSQL = `SELECT version, checksum FROM public.migrations`

func (m migrator) appliedChecksums(ctx context.Context) (map[int]string, error) {
	rows, err := m.db.Query(ctx, appliedChecksumsSQL)
	if err != nil {
		return nil, fmt.Errorf("reading applied migrations: %w", err)
	}
	defer rows.Close()

	applied := map[int]string{}
	for rows.Next() {
		var version int
		var checksum string
		if err := rows.Scan(&version, &checksum); err != nil {
			return nil, fmt.Errorf("reading applied migrations: %w", err)
		}
		applied[version] = checksum
	}
	return applied, rows.Err()
}

// confirm asks for a typed yes on stdin. It returns false for any other
// answer, or at end of input, and ctx's error when Ctrl+C arrives while it
// waits, so the prompt does not outlive the signal.
func confirm(ctx context.Context, target, action string) (bool, error) {
	fmt.Printf("\nDatabase: %s\n%s? Type yes to continue: ", target, action)
	answers := make(chan string, 1)
	go func() {
		answer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		answers <- answer
	}()
	select {
	case <-ctx.Done():
		fmt.Println()
		return false, ctx.Err()
	case answer := <-answers:
		return strings.TrimSpace(answer) == "yes", nil
	}
}

// readMigrations pairs each NNN_name.up.sql with its NNN_name.down.sql and
// reads the version and description from the up file's header comments.
func readMigrations() ([]migration, error) {
	files, err := fs.Glob(db.Migrations, "migrations/*.up.sql")
	if err != nil {
		return nil, fmt.Errorf("reading migrations: %w", err)
	}

	var migrations []migration
	for _, name := range files {
		up, err := fs.ReadFile(db.Migrations, name)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", name, err)
		}
		m := migration{up: string(up)}
		scanner := bufio.NewScanner(strings.NewReader(m.up))
		for scanner.Scan() {
			line := scanner.Text()
			if value, ok := strings.CutPrefix(line, "-- Version:"); ok {
				fmt.Sscan(strings.TrimSpace(value), &m.version)
			}
			if value, ok := strings.CutPrefix(line, "-- Description:"); ok {
				m.description = strings.TrimSpace(value)
			}
		}
		if m.version == 0 {
			return nil, fmt.Errorf("%s: missing -- Version: header", name)
		}

		down, err := fs.ReadFile(db.Migrations, strings.TrimSuffix(name, ".up.sql")+".down.sql")
		if err == nil {
			m.down = string(down)
		}
		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	for i := 1; i < len(migrations); i++ {
		if migrations[i].version == migrations[i-1].version {
			return nil, fmt.Errorf("two migrations share version %d", migrations[i].version)
		}
	}
	return migrations, nil
}
