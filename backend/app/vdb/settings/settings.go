// Package settings is what main hands every vdb command: the database, S3
// when it is enabled, and the compare options, read once from the
// environment in main; and the openings the pipeline commands
// share: the database at vdb's migration version, and the Seeder that
// reads input files from the local filesystem or from S3.
package settings

import (
	"context"
	"fmt"
	"os"

	submissionctrl "github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/submission"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// Settings is everything a command needs from the configuration.
type Settings struct {
	DB        db.Config
	S3Enabled bool
	S3        sys.S3Config
	Set       service.CompareSet // the compared set compare and seed run under
	Workloads int                // MAC workloads compare and seed run at once
}

// LocalDB is the local development database -local selects: Postgres.app's
// defaults, TLS off.
var LocalDB = db.Config{Host: "localhost:5432", DBName: "validator", User: "postgres", Password: "postgres", TLSMode: db.TLSDisable}

// Local returns the settings with the local development database and S3
// off, whatever the environment said. The compare options stay.
func (s Settings) Local() Settings {
	s.DB = LocalDB
	s.S3Enabled = false
	return s
}

// Describe is how every command names its database before it starts:
// host/name as user, and the TLS mode.
func Describe(c db.Config) string {
	return fmt.Sprintf("%s/%s as %s, tls %s", c.Host, c.DBName, c.User, c.TLSMode)
}

// OpenDatabase opens the database in s, names it, and checks that it is at
// vdb's migration version. Every pipeline command starts here. The pool
// allows two connections per workload compare runs at once, plus two
// spare: a running workload holds one for its pair stream and one for the
// write in progress, and load and summarize use one each.
func OpenDatabase(ctx context.Context, s Settings) (db.DB, error) {
	s.DB.MaxConns = 2*s.Workloads + 2
	database, err := db.Open(ctx, s.DB)
	if err != nil {
		return db.DB{}, err
	}
	fmt.Println("database", Describe(s.DB))
	if err := db.MigrationCheck(ctx, database); err != nil {
		database.Close()
		return db.DB{}, err
	}
	return database, nil
}

// OpenSeeder builds the Seeder that load and seed read input files
// through: the local filesystem, or S3 when it is enabled, which is only
// listed and read. compare and summarize read no files, so they build a
// local Seeder themselves rather than open S3 for nothing.
func OpenSeeder(ctx context.Context, database db.DB, s Settings) (submissionctrl.Seeder, error) {
	controller := submissionctrl.NewController(database, s.Set, s.Workloads)
	if !s.S3Enabled {
		return submissionctrl.NewSeeder(controller, os.Stdout), nil
	}
	inputs, err := sys.OpenS3(ctx, s.S3)
	if err != nil {
		return submissionctrl.Seeder{}, err
	}
	fmt.Printf("inputs s3://%s/%s, read-only, credentials %s\n", s.S3.Bucket, s.S3.Prefix, s.S3.Credentials)
	return submissionctrl.NewS3Seeder(controller, inputs, os.Stdout), nil
}

// CompareLine is what compare (and seed) print before they start, so a
// run's options are in its output.
func CompareLine(s Settings) string {
	return fmt.Sprintf("compared set %s, %d workloads at once", s.Set.Name, s.Workloads)
}
