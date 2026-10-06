// vdb is the validator's database tool. It applies, rolls back, and
// reports the schema migrations embedded in lib/sys/db, so the tool and the
// service always agree on the schema, and it runs the pipeline stages on
// submissions: load, compare, and summarize, each on its own and safe to
// rerun; a compared submission is hidden from the views until it is
// summarized. seed runs all three on the generated data.
//
// vdb reads the same APP_SERVER_DB_* and APP_SERVER_S3_* settings as the web
// service (vdb settings lists them), so one set of environment variables
// configures both. -local ignores them and uses the local development
// database with S3 off. APP_SERVER_COMPARE_SET (default 1) names the
// compared set compare runs under; an unknown name fails at startup.
// APP_SERVER_COMPARE_WORKLOADS (default 1) is how many MAC workloads compare
// runs at once; each holds two database connections, and the pool is sized
// from it.
//
// One package per command, each with its Run, its flags, and its help
// beside them: migrate (with status, up, and down, one file each), load,
// compare, summarize, and seed. settings is what main hands every command
// after reading the environment once, with the database and Seeder
// openings the pipeline commands share; help is what every command's help
// is built from. This file reads the environment and routes.
//
// vdb alone prints its help, and so does each command's name alone or with
// help after it; the check comes before the settings are parsed, so typing
// a name never starts work.
//
//	vdb migrate status -local
//	vdb load -dir data/generated/mac -prefix mac -side mac -label MAC -local
//	vdb compare -submission vendor-1/2026-10-02 -local
//	vdb summarize -submission vendor-1/2026-10-02 -local
//	vdb seed -local
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/compare"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/help"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/load"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/migrate"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/seed"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/settings"
	"github.com/dustinevan/x12-835-parity-validator/backend/app/vdb/summarize"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

var build = "local"

func main() {
	// Ctrl+C, or a container's stop signal, cancels the context every
	// command runs under: the statement in progress is cancelled and run
	// returns an error naming the stage that was running.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := run(ctx)
	if err == nil {
		return
	}
	if errors.Is(err, help.ErrUsage) {
		// The error and the help have been printed.
		os.Exit(2)
	}
	if errors.Is(err, context.Canceled) || ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "vdb: interrupted:", err)
		fmt.Fprintln(os.Stderr, "vdb: the stages reported above finished; the one named here did not, and rerunning it is safe")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "vdb:", err)
	os.Exit(1)
}

// command is one of vdb's commands: what runs it, and its help.
type command struct {
	run  func(context.Context, []string, settings.Settings) error
	help func() string
}

// commands is the route table, by name.
var commands = map[string]command{
	"migrate":   {migrate.Run, migrate.Help},
	"load":      {load.Run, load.Help},
	"compare":   {compare.Run, compare.Help},
	"summarize": {summarize.Run, summarize.Help},
	"seed":      {seed.Run, seed.Help},
}

func run(ctx context.Context) error {
	cfg := struct {
		conf.Version
		// S3 holds real Medicare 835s (protected health information) inside the
		// CMS environment. It stays disabled everywhere else; load then reads
		// the local filesystem.
		S3 struct {
			Enabled         bool   `conf:"default:false,env:S3_ENABLED,help:load reads inputs from S3 (read-only); only inside the CMS environment"`
			Bucket          string `conf:"env:S3_BUCKET"`
			Prefix          string `conf:"env:S3_PREFIX"`
			Region          string `conf:"env:S3_REGION"`
			Credentials     string `conf:"default:instance-role,env:S3_CREDENTIALS,help:static | instance-role (EC2) | task-role (ECS)"`
			AccessKeyID     string `conf:"env:S3_ACCESS_KEY_ID,help:static credentials only"`
			SecretAccessKey string `conf:"env:S3_SECRET_ACCESS_KEY,mask,help:static credentials only"`
			SessionToken    string `conf:"env:S3_SESSION_TOKEN,mask,help:static credentials only; optional"`
		}
		DB struct {
			Host             string        `conf:"default:localhost:5432"`
			Database         string        `conf:"default:validator"`
			User             string        `conf:"default:postgres"`
			Password         string        `conf:"default:postgres,mask"`
			TLSMode          string        `conf:"default:disable,help:disable | require | verify-full"`
			TLSRootCert      string        `conf:"help:PEM file of trusted certificate authorities; used by verify-full"`
			StatementTimeout time.Duration `conf:"default:0s,help:zero means no limit; loads and summaries run long statements"`
		}
		// Compare names the compared set (service.CompareSet) compare runs
		// under; the web service's diff endpoint reads the same setting.
		// Workloads is vdb's own: how many MAC workloads compare runs at
		// once. Each one holds two database connections, so the pool is
		// sized from it in settings.OpenDatabase.
		Compare struct {
			Set       string `conf:"default:1,help:the compared set by name: which elements the comparison leaves out of a claim"`
			Workloads int    `conf:"default:1,help:MAC workloads compared at once; set from the database's measured write rate"`
		}
	}{Version: conf.Version{Build: build, Desc: "X12 835 Parity Validator database tool"}}

	// Help comes before the settings are parsed, so typing vdb, or a
	// command's name, never starts work and never needs a valid environment.
	if len(os.Args) == 1 || help.IsHelp(os.Args[1]) {
		fmt.Print(helpText())
		return nil
	}
	name, args := os.Args[1], os.Args[2:]
	if name == "settings" {
		usage, err := conf.UsageInfo("APP_SERVER", &cfg)
		if err != nil {
			return fmt.Errorf("formatting settings: %w", err)
		}
		fmt.Print(settingsHelp(usage))
		return nil
	}
	c, ok := commands[name]
	if !ok {
		return help.UsageError(fmt.Sprintf("unknown command %q", name), helpText())
	}
	// A command's name alone prints its help, except seed, which runs with
	// its defaults.
	if (len(args) == 0 && name != "seed") || (len(args) > 0 && help.IsHelp(args[0])) {
		fmt.Print(c.help())
		return nil
	}

	// The command comes first, so conf reads no flags; settings come from
	// the environment.
	if _, err := conf.Parse("APP_SERVER", &cfg); err != nil {
		return fmt.Errorf("parsing configuration: %w", err)
	}
	// Every read of the environment happens here, in the first lines of
	// run, so what the program needs is visible in one place and a missing
	// value fails at startup, never partway through a run. A value the
	// platform provides under its own name is read here under that name and
	// handed to the SDK explicitly, because the AWS libraries read the
	// environment themselves when handed nothing.
	// ECS sets these four in a task that has a task role; they are not
	// APP_SERVER_ settings, and sys.OpenS3 uses them only for task-role.
	containerCredentialsRelativeURI := os.Getenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI")
	containerCredentialsFullURI := os.Getenv("AWS_CONTAINER_CREDENTIALS_FULL_URI")
	containerAuthorizationToken := os.Getenv("AWS_CONTAINER_AUTHORIZATION_TOKEN")
	containerAuthorizationTokenFile := os.Getenv("AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE")
	set, err := service.LookupCompareSet(cfg.Compare.Set)
	if err != nil {
		return fmt.Errorf("configuring compare: %w", err)
	}
	if cfg.Compare.Workloads < 1 {
		return fmt.Errorf("configuring compare: workloads must be at least 1, got %d", cfg.Compare.Workloads)
	}
	s := settings.Settings{
		DB: db.Config{
			Host:             cfg.DB.Host,
			DBName:           cfg.DB.Database,
			User:             cfg.DB.User,
			Password:         cfg.DB.Password,
			TLSMode:          cfg.DB.TLSMode,
			TLSRootCert:      cfg.DB.TLSRootCert,
			StatementTimeout: cfg.DB.StatementTimeout,
		},
		S3Enabled: cfg.S3.Enabled,
		S3: sys.S3Config{
			Bucket:                          cfg.S3.Bucket,
			Prefix:                          cfg.S3.Prefix,
			Region:                          cfg.S3.Region,
			Credentials:                     cfg.S3.Credentials,
			AccessKeyID:                     cfg.S3.AccessKeyID,
			SecretAccessKey:                 cfg.S3.SecretAccessKey,
			SessionToken:                    cfg.S3.SessionToken,
			ContainerCredentialsRelativeURI: containerCredentialsRelativeURI,
			ContainerCredentialsFullURI:     containerCredentialsFullURI,
			ContainerAuthorizationToken:     containerAuthorizationToken,
			ContainerAuthorizationTokenFile: containerAuthorizationTokenFile,
		},
		Set:       set,
		Workloads: cfg.Compare.Workloads,
	}
	return c.run(ctx, args, s)
}

// helpText is the top-level help: vdb alone, vdb help, or an unknown
// command.
func helpText() string {
	return help.Render("vdb: the validator's database tool, for the schema migrations and the batch pipeline", []help.Section{
		{Heading: "DESCRIPTION", Body: `  vdb applies the schema migrations embedded in it, so it and the validator
  service always agree on the schema, and runs the pipeline stages on
  submissions: load, compare, and summarize, each on its own and safe to
  rerun, since a stage replaces its own earlier results. seed runs all three
  on the generated test data. Every command names its database before it
  starts. Ctrl+C cancels the statement in progress; vdb reports what finished
  and exits non-zero, and the interrupted stage is rerun as it was.
`},
		{Heading: "HOW IT WORKS", Body: `  load       reads one submission's 835 files into the claim tables; its status becomes loaded
  compare    pairs a vendor submission's claims with the MAC's and records every difference; compared
  summarize  builds the summary tables the views read; summarized, and only then do the views show it
`},
		{Heading: "DEFAULTS", Body: `  The database and S3 come from the APP_SERVER_DB_* and APP_SERVER_S3_*
  settings (vdb settings lists them). -local is how the local development
  database is reached: it ignores those settings and uses
  localhost:5432/validator as postgres with TLS off, and S3 off. Without it,
  vdb connects to whatever the environment names, such as a bench copy or a
  deployment's database. S3 is only ever listed and read. migrate up and
  down ask for a typed yes unless -yes is given; seed reads data/generated.
`},
		{Heading: "OPTIONS", Body: `  -local     the local development database (localhost:5432/validator, TLS off) with S3 off
  -dry-run   migrate up, migrate down, load: say what would happen and stop
  -yes       migrate up, migrate down: skip the typed confirmation
             each command's own options: vdb <command> help
`},
		{Heading: "USAGE", Body: `  vdb migrate status|up|down [-local] [-yes] [-dry-run]
  vdb load -dir DIR -prefix PREFIX -side mac|vendor -label LABEL [-date YYYY-MM-DD] [-local] [-dry-run]
  vdb compare -submission PREFIX [-local]
  vdb summarize -submission PREFIX [-local]
  vdb seed [-dir data/generated] [-local]
  vdb settings
  vdb <command> help
`},
		{Heading: "EXAMPLES", Body: `  vdb migrate status -local                                    # what the local database has applied
  vdb seed -local                                              # load, compare, and summarize data/generated
  vdb migrate down -local -yes && vdb migrate up -local -yes   # rebuild the local schema after editing 001; then seed
  vdb compare -submission vendor-1/2026-10-02 -local           # rerun one stage; summarize it next
  vdb load -dir data/generated/mac -prefix mac -side mac -label MAC -local -dry-run   # the files a load would read
  APP_SERVER_DB_DATABASE=validator_bench vdb migrate status    # another database, named by the environment
`},
	})
}

// settingsHelp is the help for vdb settings. usage is what conf formats
// from the configuration struct; only its ENVIRONMENT section is shown,
// since the --flags it also lists are never read: the command comes first.
func settingsHelp(usage string) string {
	_, env, found := strings.Cut(usage, "ENVIRONMENT\n")
	if !found {
		env = usage
	}
	return help.Render("vdb settings: the APP_SERVER_* settings vdb reads, with their defaults", []help.Section{
		{Heading: "DESCRIPTION", Body: `  vdb reads these from the environment, never from flags: the same database
  and S3 settings as the validator service, plus COMPARE_WORKLOADS, which is
  vdb's own. -local replaces the database settings with the local
  development database and turns S3 off; the compare settings still apply.
`},
		{Heading: "SETTINGS", Body: env},
	})
}
