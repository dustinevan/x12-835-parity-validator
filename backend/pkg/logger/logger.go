// Package logger builds the one logger the programs use: a zap SugaredLogger
// that writes JSON lines to a writer, one object per line, with RFC3339
// times to the nanosecond (the timezone spelled -06:00, not -0600), and the
// service name and build on every line.
//
// The web service passes os.Stdout: the container platform collects a
// task's stdout, and a Fargate task's disk is gone when the task stops, so
// there is no file output and no rotation. lib code never builds a logger
// or logs on its own; it returns wrapped errors, and the few places that
// must log (db.StatusCheck, the request middleware) take the logger as an
// argument. There is no global logger.
//
// The message is the phase or event and the details are key-value pairs:
//
//	log.Infow("startup", "status", "database ready")
//	log.Errorw("shutdown", "error", err)
//
// Nothing that could be PHI goes in a log line: no claim values, MBIs,
// match keys, file contents, URL paths, or query strings. Submission
// prefixes, counts, route patterns, and durations are fine.
package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a logger that writes JSON lines at Info level and above to w,
// each with service and build fields. Debug lines are not written.
func New(w io.Writer, service, build string) *zap.SugaredLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(w), zapcore.InfoLevel)
	return zap.New(core).Sugar().With("service", service, "build", build)
}
