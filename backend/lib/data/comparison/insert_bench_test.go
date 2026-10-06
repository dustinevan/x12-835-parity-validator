package comparison

// Raw insert rate of CreateMany: 10,000 synthetic comparison rows written to
// a scratch database in CreateMany calls of 100, 1,000, and 10,000 rows,
// each with one transaction around all the calls and with one transaction
// per call. The pipeline makes one CreateMany call per 10,000 pairs, each
// standing on its own (Controller.compareWorkload), so the per-call column
// is the pipeline's shape.
//
// The benchmark opens a database, so it runs only when VALIDATOR_BENCH_DB=1.
// It uses a scratch database named validator_bench. Create it first and
// apply the schema with vdb:
//
//	createdb -h localhost -U postgres validator_bench
//	APP_SERVER_DB_DATABASE=validator_bench bin/vdb migrate up -yes
//	VALIDATOR_BENCH_DB=1 go test -run '^$' -bench CreateMany -benchtime 10x ./backend/lib/data/comparison/
//
// It writes only rows whose submission_prefix is benchPrefix, and removes
// them (outside the timer) before each run.

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

const (
	benchPrefix = "bench/comparison"
	benchRows   = 10000
)

// benchConfig names the scratch database only.
var benchConfig = db.Config{
	Host:     "localhost:5432",
	DBName:   "validator_bench",
	User:     "postgres",
	Password: "postgres",
	TLSMode:  db.TLSDisable,
}

// benchWriter holds the database and the rows one benchmark run writes.
type benchWriter struct {
	database db.DB
	rows     []Comparison
}

const vacuumSQL = `VACUUM public.comparison`

// reset removes the benchmark's rows and vacuums the table, so every timed
// run writes into the same state, with no dead rows left by the run before.
func (w benchWriter) reset(ctx context.Context) error {
	if err := NewTableAccess(w.database).DeleteBySubmission(ctx, benchPrefix); err != nil {
		return err
	}
	_, err := w.database.Exec(ctx, vacuumSQL)
	return err
}

// oneTransaction writes every row in CreateMany calls of batchSize rows
// inside a single transaction.
func (w benchWriter) oneTransaction(ctx context.Context, batchSize int) error {
	tx, err := w.database.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for start := 0; start < len(w.rows); start += batchSize {
		end := min(start+batchSize, len(w.rows))
		if _, err := NewTableAccess(tx).CreateMany(ctx, w.rows[start:end]); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// transactionPerBatch writes every row in CreateMany calls of batchSize
// rows, each call in its own transaction.
func (w benchWriter) transactionPerBatch(ctx context.Context, batchSize int) error {
	for start := 0; start < len(w.rows); start += batchSize {
		end := min(start+batchSize, len(w.rows))
		tx, err := w.database.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := NewTableAccess(tx).CreateMany(ctx, w.rows[start:end]); err != nil {
			tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

// benchTags are twelve of the element tags the comparison stage writes.
var benchTags = []string{"CLP02", "CLP04", "CLP05", "SVC02", "SVC03", "SVC04", "CAS*CO*45", "CAS*PR*1", "CAS*PR*2", "MOA", "AMT*B6", "PLB"}

// syntheticComparisons returns n distinct rows: a quarter are "different"
// with three tags, the rest "equal" with none, close to the generated data.
func syntheticComparisons(n int) []Comparison {
	rows := make([]Comparison, n)
	for i := range rows {
		key := sha256.Sum256([]byte(fmt.Sprintf("bench-comparison-%d", i)))
		c := Comparison{
			SubmissionPrefix:              benchPrefix,
			MatchKey:                      key[:],
			Outcome:                       OutcomeEqual,
			MACWorkloadID:                 fmt.Sprintf("%05d", 10000+i%12),
			ClaimType:                     "professional",
			MedicareBeneficiaryIdentifier: fmt.Sprintf("BENCH%07d", i),
			DateOfService:                 time.Date(2026, time.Month(1+i%12), 1+i%28, 0, 0, 0, 0, time.UTC),
			TotalClaimChargeAmount:        fmt.Sprintf("%d.%02d", 100+i%900, i%100),
			MACClaimStatusCode:            "1",
			VendorClaimStatusCode:         "1",
			MACClaimPaymentAmount:         "200.00",
			VendorClaimPaymentAmount:      "200.00",
		}
		if i%4 == 0 {
			c.Outcome = OutcomeDifferent
			c.VendorClaimPaymentAmount = "190.50"
			j := i / 4
			c.DifferenceCount = 3
			c.Tags = []string{benchTags[j%12], benchTags[(j+5)%12], benchTags[(j+7)%12]}
		}
		rows[i] = c
	}
	return rows
}

func BenchmarkCreateMany(b *testing.B) {
	if os.Getenv("VALIDATOR_BENCH_DB") != "1" {
		b.Skip("set VALIDATOR_BENCH_DB=1 to run against a scratch database named validator_bench (see the file header)")
	}
	ctx := context.Background()
	database, err := db.Open(ctx, benchConfig)
	if err != nil {
		b.Fatal(err)
	}
	defer database.Close()
	writer := benchWriter{database: database, rows: syntheticComparisons(benchRows)}

	for _, batchSize := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("batch=%d/one-transaction", batchSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				if err := writer.reset(ctx); err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
				if err := writer.oneTransaction(ctx, batchSize); err != nil {
					b.Fatal(err)
				}
			}
			b.ReportMetric(float64(benchRows*b.N)/b.Elapsed().Seconds(), "rows/s")
		})
		b.Run(fmt.Sprintf("batch=%d/transaction-per-batch", batchSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				if err := writer.reset(ctx); err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
				if err := writer.transactionPerBatch(ctx, batchSize); err != nil {
					b.Fatal(err)
				}
			}
			b.ReportMetric(float64(benchRows*b.N)/b.Elapsed().Seconds(), "rows/s")
		})
	}
	if err := writer.reset(ctx); err != nil {
		b.Fatal(err)
	}
}
