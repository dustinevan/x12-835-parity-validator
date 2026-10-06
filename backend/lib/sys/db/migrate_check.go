package db

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"strconv"
	"strings"
)

// Migrations holds every migration file. The service checks the applied
// version against it at startup; backend/app/vdb applies and rolls back migrations.
//
//go:embed migrations/*.sql
var Migrations embed.FS

// MigrationCheck fails when the database is behind or ahead of this binary's
// migrations. Ahead is as wrong as behind: it means an applied migration file
// was deleted or renumbered.
func MigrationCheck(ctx context.Context, db DB) error {
	expected, err := HighestMigrationVersion()
	if err != nil {
		return fmt.Errorf("db: checking migrations: %w", err)
	}

	var applied int
	if err := db.QueryRow(ctx, "SELECT COALESCE(MAX(version), 0) FROM public.migrations").Scan(&applied); err != nil {
		return fmt.Errorf("db: checking migrations: reading applied version: %w", err)
	}

	if applied < expected {
		return fmt.Errorf("db: database is %d migration(s) behind (applied %d, expected %d); run: vdb migrate up", expected-applied, applied, expected)
	}
	if applied > expected {
		return fmt.Errorf("db: database is %d migration(s) ahead of this binary (applied %d, expected %d); restore the missing migration files", applied-expected, applied, expected)
	}
	return nil
}

// HighestMigrationVersion reads the "-- Version: N" header of every embedded
// .up.sql file and returns the largest N.
func HighestMigrationVersion() (int, error) {
	entries, err := Migrations.ReadDir("migrations")
	if err != nil {
		return 0, fmt.Errorf("reading embedded migrations: %w", err)
	}

	highest := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		data, err := Migrations.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return 0, fmt.Errorf("reading migration %s: %w", entry.Name(), err)
		}
		version, ok := migrationVersion(string(data))
		if ok && version > highest {
			highest = version
		}
	}
	if highest == 0 {
		return 0, fmt.Errorf("no migration versions found in embedded files")
	}
	return highest, nil
}

func migrationVersion(sql string) (int, bool) {
	scanner := bufio.NewScanner(strings.NewReader(sql))
	for scanner.Scan() {
		line := scanner.Text()
		lower := strings.ToLower(line)
		if !strings.HasPrefix(lower, "-- version:") && !strings.HasPrefix(lower, "--version:") {
			continue
		}
		_, value, _ := strings.Cut(line, ":")
		version, err := strconv.Atoi(strings.TrimSpace(value))
		return version, err == nil
	}
	return 0, false
}
