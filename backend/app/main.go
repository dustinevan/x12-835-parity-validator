// The validator service: configuration, startup, and assembly. Handlers
// live in app/api, route registration and middleware in app/web, and the
// UI handler in app/ui.
package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"

	"github.com/dustinevan/x12-835-parity-validator/backend/app/ui"
	appweb "github.com/dustinevan/x12-835-parity-validator/backend/app/web"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/control/views"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/logger"
)

var build = "local"

// The frontend build writes site/dist before the release binary is compiled.
// Keeping the assets in the executable gives the container one server process.
//
//go:embed site
var bundle embed.FS

func main() {
	log := logger.New(os.Stdout, "validator", build)
	if err := run(log); err != nil {
		if errors.Is(err, errUsage) {
			// The argument and the help have been printed.
			os.Exit(2)
		}
		log.Errorw("shutdown", "error", err)
		_ = log.Sync()
		os.Exit(1)
	}
	_ = log.Sync()
}

func run(log *zap.SugaredLogger) error {
	cfg := struct {
		conf.Version
		Web struct {
			APIHost           string        `conf:"default:127.0.0.1:8089"`
			UIProxy           string        `conf:"help:Development UI server URL; empty serves the bundled site"`
			ReadHeaderTimeout time.Duration `conf:"default:5s"`
			ReadTimeout       time.Duration `conf:"default:5s"`
			WriteTimeout      time.Duration `conf:"default:0,help:off; the request deadline is RequestTimeout"`
			RequestTimeout    time.Duration `conf:"default:5m,help:deadline for one API request (the Timeout middleware)"`
			IdleTimeout       time.Duration `conf:"default:120s"`
			ShutdownTimeout   time.Duration `conf:"default:30s"`
		}
		Store struct {
			Root   string `conf:"default:data,help:local directory for inputs (while S3 is disabled) and for output"`
			Output string `conf:"default:diffs,help:output directory under the root"`
		}
		// S3 holds real Medicare 835s (protected health information) inside the
		// CMS environment. It stays disabled everywhere else; inputs then come
		// from the local filesystem.
		S3 struct {
			Enabled         bool   `conf:"default:false,env:S3_ENABLED,help:read inputs from S3 (read-only); only inside the CMS environment"`
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
			StatementTimeout time.Duration `conf:"default:30s"`
		}
		// Compare names the compared set (service.CompareSet) the diff
		// endpoint compares under; vdb reads the same setting for its runs.
		Compare struct {
			Set string `conf:"default:1,help:the compared set by name: which elements the comparison leaves out of a claim"`
		}
	}{Version: conf.Version{Build: build, Desc: "X12 835 Parity Validator"}}

	// validator's only arguments are help. The check comes before the
	// settings are parsed, so help never needs a valid environment, and any
	// other argument is a mistake: settings come from APP_SERVER_* variables,
	// never from flags.
	if len(os.Args) > 1 {
		usage, err := conf.UsageInfo("APP_SERVER", &cfg)
		if err != nil {
			return fmt.Errorf("app: formatting settings: %w", err)
		}
		if !isHelp(os.Args[1]) {
			fmt.Fprintf(os.Stderr, "validator: unknown argument %q\n\n", os.Args[1])
			fmt.Print(helpText(usage))
			return errUsage
		}
		fmt.Print(helpText(usage))
		return nil
	}
	if _, err := conf.Parse("APP_SERVER", &cfg); err != nil {
		return fmt.Errorf("app: parsing configuration: %w", err)
	}
	// Every read of the environment happens here, in the first lines of
	// run, so what the program needs is visible in one place and a missing
	// value fails at startup, never after the program is serving. A value
	// the platform provides under its own name is read here under that name
	// and handed to the SDK explicitly, because the AWS libraries read the
	// environment themselves when handed nothing.
	// ECS sets these four in a task that has a task role; they are not
	// APP_SERVER_ settings, and sys.OpenS3 uses them only for task-role.
	containerCredentialsRelativeURI := os.Getenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI")
	containerCredentialsFullURI := os.Getenv("AWS_CONTAINER_CREDENTIALS_FULL_URI")
	containerAuthorizationToken := os.Getenv("AWS_CONTAINER_AUTHORIZATION_TOKEN")
	containerAuthorizationTokenFile := os.Getenv("AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE")
	defer log.Infow("shutdown complete")
	if !filepath.IsLocal(cfg.Store.Output) {
		return fmt.Errorf("app: configuring output: directory must be relative to the store")
	}
	set, err := service.LookupCompareSet(cfg.Compare.Set)
	if err != nil {
		return fmt.Errorf("app: configuring compare: %w", err)
	}
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("app: formatting configuration: %w", err)
	}
	log.Infow("startup", "config", out)

	var site http.Handler
	if cfg.Web.UIProxy != "" {
		site, err = ui.Proxy(log, cfg.Web.UIProxy)
	} else {
		dist, subErr := fs.Sub(bundle, "site/dist")
		if subErr != nil {
			return fmt.Errorf("app: opening bundled UI: %w", subErr)
		}
		site, err = ui.Bundled(dist)
	}
	if err != nil {
		return fmt.Errorf("app: configuring UI: %w", err)
	}

	// Output (difference files) always goes to the local filesystem.
	if err := os.MkdirAll(cfg.Store.Root, 0700); err != nil {
		return fmt.Errorf("app: creating store root: %w", err)
	}
	root, err := os.OpenRoot(cfg.Store.Root)
	if err != nil {
		return fmt.Errorf("app: opening store root: %w", err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			log.Errorw("closing store root", "error", err)
		}
	}()
	outputs := sys.FileSystem{Root: root}

	// Inputs come from the local filesystem unless S3 is explicitly enabled.
	// S3 holds protected health information; it is read-only (sys.S3 has no
	// write method) and enabled only inside the CMS environment.
	var inputs sys.Reader = outputs
	if !cfg.S3.Enabled {
		log.Infow("startup", "inputs", "filesystem", "root", cfg.Store.Root, "s3", "disabled")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// The client is built only from these settings: no environment
		// variables, shared config files, or default credential chain.
		s3Inputs, err := sys.OpenS3(ctx, sys.S3Config{
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
		})
		if err != nil {
			return fmt.Errorf("app: opening S3: %w", err)
		}
		inputs = s3Inputs
		// The token is never logged; containerEndpoint says only whether ECS
		// set one of the two endpoint variables.
		log.Infow("startup", "inputs", "s3", "bucket", cfg.S3.Bucket, "prefix", cfg.S3.Prefix, "access", "read-only",
			"credentials", cfg.S3.Credentials, "containerEndpoint", containerCredentialsRelativeURI != "" || containerCredentialsFullURI != "")
	}

	log.Infow("startup", "status", "opening database", "host", cfg.DB.Host, "database", cfg.DB.Database, "tls", cfg.DB.TLSMode)
	dbConfig := db.Config{
		Host:             cfg.DB.Host,
		DBName:           cfg.DB.Database,
		User:             cfg.DB.User,
		Password:         cfg.DB.Password,
		TLSMode:          cfg.DB.TLSMode,
		TLSRootCert:      cfg.DB.TLSRootCert,
		StatementTimeout: cfg.DB.StatementTimeout,
	}
	dbConn, err := db.Open(context.Background(), dbConfig)
	if err != nil {
		return fmt.Errorf("app: opening database: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "closing database")
		dbConn.Close()
	}()
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelStartup()
	if err := db.StatusCheck(startupCtx, log, dbConn); err != nil {
		return fmt.Errorf("app: checking database: %w", err)
	}
	if err := db.MigrationCheck(startupCtx, dbConn); err != nil {
		return fmt.Errorf("app: checking database: %w", err)
	}
	log.Infow("startup", "status", "database ready")

	// The diff service is shared by pointer: it holds the lock that allows
	// one diff at a time.
	diffService := &service.Service{Inputs: inputs, Outputs: outputs, OutputDirectory: cfg.Store.Output, Set: set}
	log.Infow("startup", "comparedSet", set.Name)
	router := appweb.RegisterRoutes(appweb.RouteConfig{
		Log:             log,
		RequestTimeout:  cfg.Web.RequestTimeout,
		Diff:            diffService,
		Outputs:         outputs,
		OutputDirectory: cfg.Store.Output,
		Views:           views.NewController(dbConn),
	})

	// The router's own mux answers 404 for an unknown API path and 405 for
	// a known path with the wrong method. Only /api/ and the root paths
	// /health, /diff, and /diffs/ reach it; every other path is the UI.
	api := router.Handler()
	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/api", http.NotFoundHandler())
	mux.Handle("/health", api)
	mux.Handle("/diff", api)
	mux.Handle("/diffs/", api)
	mux.Handle("/", site)

	requestContext, cancelRequests := context.WithCancel(context.Background())
	defer cancelRequests()
	server := http.Server{
		Addr: cfg.Web.APIHost, Handler: mux,
		ReadTimeout: cfg.Web.ReadTimeout, WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout: cfg.Web.IdleTimeout, ReadHeaderTimeout: cfg.Web.ReadHeaderTimeout,
		BaseContext: func(_ net.Listener) context.Context { return requestContext },
		ErrorLog:    zap.NewStdLog(log.Desugar()),
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdown)
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	log.Infow("listening", "address", server.Addr)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("app: running HTTP server: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "signal", sig.String())
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
	defer cancel()
	// Cancel requests partway through shutdown, leaving time to release
	// their resources before the server's shutdown deadline.
	cancelTimer := time.AfterFunc(cfg.Web.ShutdownTimeout*2/3, cancelRequests)
	defer cancelTimer.Stop()
	if err := server.Shutdown(ctx); err != nil {
		cancelRequests()
		closeErr := server.Close()
		return fmt.Errorf("app: stopping HTTP server: %w", errors.Join(err, closeErr))
	}
	return nil
}
