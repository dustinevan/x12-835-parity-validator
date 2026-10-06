// Package db opens Postgres connections and checks the schema version: a
// pgx pool behind a small value type, named arguments passed as a map, and
// SQL written next to the method that runs it.
package db

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("db: not found")

// TLS modes for Config.TLSMode, named as Postgres names its sslmode values.
const (
	TLSDisable    = "disable"     // no encryption: local development only
	TLSRequire    = "require"     // encrypted; the server's certificate is not checked
	TLSVerifyFull = "verify-full" // encrypted; the certificate must chain to TLSRootCert and name Host
)

// Config is the required properties to use the database.
type Config struct {
	User     string
	Password string
	Host     string
	DBName   string
	TLSMode  string
	// TLSRootCert is a PEM file of the certificate authorities verify-full
	// trusts, such as the Amazon RDS bundle at /etc/ssl/rds-global-bundle.pem
	// in the image. The other modes do not use it.
	TLSRootCert string
	// StatementTimeout bounds each statement. Bulk loads and the diff read
	// long-running cursors, so this is configurable. Zero means no limit.
	StatementTimeout time.Duration
	// MaxConns is the most connections the pool opens at once. Zero keeps
	// pgxpool's default, the larger of 4 and the CPU count. vdb sets it from
	// how many workloads compare runs at once, since each holds two.
	MaxConns int
}

func connString(c Config) string {
	// sslrootcert is always written, empty unless verify-full, so a
	// PGSSLROOTCERT environment variable cannot change which certificate
	// authorities are trusted.
	rootCert := ""
	if c.TLSMode == TLSVerifyFull {
		rootCert = c.TLSRootCert
	}

	q := make(url.Values)
	q.Set("sslmode", c.TLSMode)
	q.Set("sslrootcert", rootCert)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     c.Host,
		Path:     c.DBName,
		RawQuery: q.Encode(),
	}
	return u.String()
}

// ConnStringWithPasswordRedacted is safe to log.
func ConnStringWithPasswordRedacted(c Config) string {
	redacted := c
	redacted.Password = "xxxxxx"
	return connString(redacted)
}

// Open checks the TLS settings and creates the connection pool. verify-full
// needs a root certificate file; pgx reads it here, so a missing or
// unreadable file fails Open.
func Open(ctx context.Context, cfg Config) (DB, error) {
	switch cfg.TLSMode {
	case TLSDisable, TLSRequire:
	case TLSVerifyFull:
		if cfg.TLSRootCert == "" {
			return DB{}, fmt.Errorf("db: configuring TLS: verify-full needs a root certificate file")
		}
	default:
		return DB{}, fmt.Errorf("db: configuring TLS: mode must be %s, %s, or %s", TLSDisable, TLSRequire, TLSVerifyFull)
	}
	poolConfig, err := pgxpool.ParseConfig(connString(cfg))
	if err != nil {
		return DB{}, fmt.Errorf("db: parsing config: %w", err)
	}
	if cfg.StatementTimeout > 0 {
		poolConfig.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(cfg.StatementTimeout.Milliseconds())
	}
	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = int32(cfg.MaxConns)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return DB{}, fmt.Errorf("db: opening: %w", err)
	}
	return DB{pgp: pool}, nil
}

// StatusCheck returns nil if it can talk to the database, retrying briefly
// while the database starts.
func StatusCheck(ctx context.Context, log *zap.SugaredLogger, db DB) error {
	var pingError error
	for attempts := 1; ; attempts++ {
		if attempts > 5 {
			return fmt.Errorf("db: checking status: %d ping attempts failed: %w", attempts-1, pingError)
		}
		pingError = db.pgp.Ping(ctx)
		if pingError == nil {
			return nil
		}
		log.Warnw("db: ping failed", "attempt", attempts, "error", pingError)
		select {
		case <-ctx.Done():
			return fmt.Errorf("db: checking status: %w", ctx.Err())
		case <-time.After(time.Duration(attempts) * 100 * time.Millisecond):
		}
	}
}

type DB struct {
	pgp *pgxpool.Pool
}

// Query, QueryRow, and Exec accept either positional arguments or a single
// map[string]any of named arguments.
func (db DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if named, ok := namedArgs(args); ok {
		return db.pgp.Query(ctx, sql, named)
	}
	return db.pgp.Query(ctx, sql, args...)
}

func (db DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if named, ok := namedArgs(args); ok {
		return db.pgp.QueryRow(ctx, sql, named)
	}
	return db.pgp.QueryRow(ctx, sql, args...)
}

func (db DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if named, ok := namedArgs(args); ok {
		return db.pgp.Exec(ctx, sql, named)
	}
	return db.pgp.Exec(ctx, sql, args...)
}

// CopyFrom bulk-loads rows with Postgres COPY. Claim and difference loads
// use it instead of row-by-row inserts.
func (db DB) CopyFrom(ctx context.Context, table pgx.Identifier, columns []string, rows pgx.CopyFromSource) (int64, error) {
	return db.pgp.CopyFrom(ctx, table, columns, rows)
}

func (db DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pgp.Begin(ctx)
}

func (db DB) Close() {
	db.pgp.Close()
}

func namedArgs(args []any) (pgx.NamedArgs, bool) {
	if len(args) != 1 {
		return nil, false
	}
	switch named := args[0].(type) {
	case map[string]any:
		return pgx.NamedArgs(named), true
	case pgx.NamedArgs:
		return named, true
	}
	return nil, false
}
