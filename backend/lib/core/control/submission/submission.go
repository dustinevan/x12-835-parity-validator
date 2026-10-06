// Package submission runs one submission through the pipeline stages, each
// its own step and each safe to rerun: load its claims; for a vendor
// submission, compare it with the MAC submission; then summarize it.
// Compare works by MAC workload, a set number of them at once (NewController's
// workloads; vdb's APP_SERVER_COMPARE_WORKLOADS, default 1), reading claim
// pairs as the database streams them, so no submission is ever held in
// memory; the caller (vdb) runs submissions in sequence. Seeder (seed.go)
// loads, compares, and summarizes a directory of submissions through a
// Controller.
package submission

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/query"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/cell"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claimrepeat"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/comparison"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/difference"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/maccell"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/measure"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/payment"
	submissiontable "github.com/dustinevan/x12-835-parity-validator/backend/lib/data/submission"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/tagcell"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/tagtotal"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// Controller loads, compares, and summarizes submissions.
type Controller struct {
	db db.DB

	// set is the compared set Compare runs under (APP_SERVER_COMPARE_SET):
	// the elements the comparison leaves out of a claim. Compare records
	// its name on the submission.
	set service.CompareSet

	// workloads is how many MAC workloads Compare runs at once (vdb's
	// APP_SERVER_COMPARE_WORKLOADS, at least 1). Each running workload holds
	// two connections from the pool, which the caller sizes for it.
	workloads int
}

// NewController returns a Controller. set is the compared set Compare runs
// under and workloads how many MAC workloads it compares at once (below 1
// is treated as 1).
func NewController(database db.DB, set service.CompareSet, workloads int) Controller {
	return Controller{db: database, set: set, workloads: max(workloads, 1)}
}

// InputFile is one 835 file of a submission. Name is relative to the
// submission's prefix. Input files are only ever read.
type InputFile struct {
	Name string
	Data []byte
}

// ListDirectory lists the names of the regular files directly under dir
// (not dot files, not subdirectories), in name order: the files
// ReadDirectory reads. vdb's load -dry-run shows them without reading them.
func ListDirectory(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("submission: reading %s: %w", dir, err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.Type().IsRegular() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		names = append(names, entry.Name())
	}
	return names, nil
}

// ReadDirectory reads every file ListDirectory names as one submission's
// input files.
func ReadDirectory(dir string) ([]InputFile, error) {
	names, err := ListDirectory(dir)
	if err != nil {
		return nil, err
	}
	var files []InputFile
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("submission: reading %s: %w", name, err)
		}
		files = append(files, InputFile{Name: name, Data: data})
	}
	return files, nil
}

// Loaded is what Load wrote: its row counts.
type Loaded struct {
	Claims       int64 // claim rows
	Payments     int64 // payment rows, one per 835 transaction
	RepeatedKeys int64 // claim_repeat rows: match keys on more than one claim
}

// Load replaces a submission: in one transaction it deletes everything
// derived from the submission's earlier files (claims, payments, the keys
// the claims repeat, for the MAC its legacy cube, and for a vendor its
// comparisons, differences, and summaries), inserts the submission, its
// claims, and its payments (one row per 835 transaction, the header its
// claims sit under), and records the match keys that occur on more than
// one claim (claim_repeat), found here once so that compare never counts
// claims per key. For the MAC submission it also writes the legacy cube
// (mac_cell) from the claims, the Overview's legacy figures. It returns
// the row counts it wrote.
func (c Controller) Load(ctx context.Context, ns submissiontable.NewSubmission, files []InputFile) (Loaded, error) {
	if ns.Side != submissiontable.SideMAC && ns.Side != submissiontable.SideVendor {
		return Loaded{}, fmt.Errorf("submission: loading %s: side must be %q or %q", ns.Prefix, submissiontable.SideMAC, submissiontable.SideVendor)
	}
	if ns.Side == submissiontable.SideMAC {
		existing, err := submissiontable.NewTableAccess(c.db).MAC(ctx)
		if err == nil && existing.Prefix != ns.Prefix {
			return Loaded{}, fmt.Errorf("submission: loading %s: the MAC submission is already %s", ns.Prefix, existing.Prefix)
		}
	}

	var claims []claim.Claim
	var payments []payment.Payment
	for _, f := range files {
		fileClaims, filePayments, err := service.Split(ns.Prefix, f.Name, f.Data)
		if err != nil {
			return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
		}
		claims = append(claims, fileClaims...)
		payments = append(payments, filePayments...)
	}
	ns.ClaimCount = int64(len(claims))

	tx, err := c.db.Begin(ctx)
	if err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: beginning: %w", ns.Prefix, err)
	}
	defer tx.Rollback(ctx)

	// Remove every row derived from the submission's earlier files.
	if err := tagtotal.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := tagcell.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := cell.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := measure.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := difference.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := comparison.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := claimrepeat.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := maccell.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := claim.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := payment.NewTableAccess(tx).DeleteBySubmission(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if err := submissiontable.NewTableAccess(tx).Delete(ctx, ns.Prefix); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}

	if err := submissiontable.NewTableAccess(tx).Create(ctx, ns); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if _, err := claim.NewTableAccess(tx).CreateMany(ctx, claims); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if _, err := payment.NewTableAccess(tx).CreateMany(ctx, payments); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	repeated, err := query.NewLoadQuery(tx).CreateClaimRepeats(ctx, ns.Prefix)
	if err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
	}
	if ns.Side == submissiontable.SideMAC {
		if err := query.NewLoadQuery(tx).CreateMACCells(ctx, ns.Prefix); err != nil {
			return Loaded{}, fmt.Errorf("submission: loading %s: %w", ns.Prefix, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return Loaded{}, fmt.Errorf("submission: loading %s: committing: %w", ns.Prefix, err)
	}
	return Loaded{Claims: ns.ClaimCount, Payments: int64(len(payments)), RepeatedKeys: repeated}, nil
}

// batchRows is how many comparison rows compareWorkload collects before it
// writes them, and their difference rows, with COPY. Measured 2026-10-01:
// 10,000-row COPYs ran 18–28% faster than 1,000-row ones.
const batchRows = 10000

// Compare compares a vendor submission with the MAC submission by MAC
// workload, the Controller's workloads of them at once, under the
// Controller's compared set, and records the result, replacing any earlier
// run.
// The match key contains the workload, so no pair, set-aside, duplicated,
// missing, or extra key crosses one, and the workloads are independent of
// each other: each writes its own rows, whatever order they finish in. The
// database does the matching: for each workload the pairs (a key on
// exactly one claim per side, read from the claim_repeat rows load wrote)
// stream out of a join and Go compares them one at a time
// (service.ComparePair), writing the rows in batches of batchRows with
// COPY; the set-aside, duplicated, missing, and extra rows are written by
// four INSERT ... SELECT statements. Nothing holds a submission, or a
// workload, in memory. Every statement stands on its own: no transaction
// is held open while Go reads and compares, since the status gate hides
// the submission until the end and the submission's old rows are deleted
// before the first workload is written. With workloads of 1 the workloads
// run one after another in workload order, as a single loop would.
//
// The submission is first set back to 'loaded', so the views hide it and
// Summarize refuses it before its rows change; a run that dies leaves it
// there. After the last workload the status is set to 'compared' with the
// set's name, and the submission stays hidden until Summarize runs. The
// first workload to fail stops the others (their context is canceled),
// sets the status to 'failed', and is the error returned; the failed
// workloads are left partly written. Rerunning Compare is safe either way.
// Compare does not touch the summary tables.
func (c Controller) Compare(ctx context.Context, vendorPrefix string) (service.Summary, error) {
	submissions := submissiontable.NewTableAccess(c.db)
	vendor, err := submissions.Get(ctx, vendorPrefix)
	if err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}
	if vendor.Side != submissiontable.SideVendor {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: not a vendor submission", vendorPrefix)
	}
	mac, err := submissions.MAC(ctx)
	if err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}
	if err := submissions.ClearCompared(ctx, vendor.Prefix); err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}
	if err := difference.NewTableAccess(c.db).DeleteBySubmission(ctx, vendor.Prefix); err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}
	if err := comparison.NewTableAccess(c.db).DeleteBySubmission(ctx, vendor.Prefix); err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}

	workloads, err := query.NewCompareQuery(c.db).ListWorkloads(ctx, mac.Prefix, vendor.Prefix)
	if err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}
	// Run the workloads, c.workloads at a time. The loop takes a slot from
	// running before it starts each one and the workload's goroutine gives
	// the slot back when it is done, so the loop blocks while every slot is
	// taken. Each goroutine writes its own entry of summaries, which are
	// added up after every goroutine has finished, so nothing is shared
	// while they run except the first failure, kept under a mutex. A
	// failure cancels runCtx, which ends the other workloads' streams and
	// writes; the loop then starts no more. ctx itself stays usable for
	// the status update at the end.
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	running := make(chan struct{}, c.workloads)
	summaries := make([]service.Summary, len(workloads))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failedWorkload string // the first workload to fail, with its error
	var failure error
	for i, workload := range workloads {
		running <- struct{}{}
		if runCtx.Err() != nil {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-running }()
			s, err := c.compareWorkload(runCtx, mac.Prefix, vendor.Prefix, workload)
			if err != nil {
				mu.Lock()
				if failure == nil {
					failedWorkload, failure = workload, err
				}
				mu.Unlock()
				cancel()
				return
			}
			summaries[i] = s
		}()
	}
	wg.Wait()
	if failure == nil && ctx.Err() != nil {
		// The caller gave up before every workload had started.
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, ctx.Err())
	}
	if failure != nil {
		if failErr := submissions.SetFailed(ctx, vendor.Prefix); failErr != nil {
			return service.Summary{}, fmt.Errorf("submission: comparing %s: workload %s: %w (and %v)", vendorPrefix, failedWorkload, failure, failErr)
		}
		return service.Summary{}, fmt.Errorf("submission: comparing %s: workload %s: %w", vendorPrefix, failedWorkload, failure)
	}
	var summary service.Summary
	for _, s := range summaries {
		summary.Equal += s.Equal
		summary.Different += s.Different
		summary.Missing += s.Missing
		summary.Duplicated += s.Duplicated
		summary.Extra += s.Extra
		summary.SetAside += s.SetAside
		summary.Differences += s.Differences
	}

	if err := submissions.SetCompared(ctx, vendor.Prefix, c.set.Name); err != nil {
		return service.Summary{}, fmt.Errorf("submission: comparing %s: %w", vendorPrefix, err)
	}
	return summary, nil
}

// compareWorkload compares one MAC workload of a vendor submission and
// returns its counts: stream the pairs and compare each as it arrives,
// writing every batchRows pairs with COPY; then write the set-aside,
// duplicated, missing, and extra rows. The stream reads on its own
// connection while the writes go through the pool, so a running workload
// holds two connections at most; while a COPY runs, the stream's
// connection waits. The caller has already removed the submission's
// earlier rows.
func (c Controller) compareWorkload(ctx context.Context, macPrefix, vendorPrefix, workload string) (service.Summary, error) {
	compare := query.NewCompareQuery(c.db)
	comparisons := comparison.NewTableAccess(c.db)
	differences := difference.NewTableAccess(c.db)

	pairs, err := compare.Pairs(ctx, macPrefix, vendorPrefix, workload)
	if err != nil {
		return service.Summary{}, err
	}
	defer pairs.Close()

	var summary service.Summary
	comparisonRows := make([]comparison.Comparison, 0, batchRows)
	var differenceRows []difference.Difference
	for pairs.Next() {
		mac, vendor := pairs.Pair()
		row, diffs, err := service.ComparePair(c.set, vendorPrefix, mac, vendor)
		if err != nil {
			return service.Summary{}, err
		}
		if row.Outcome == comparison.OutcomeEqual {
			summary.Equal++
		} else {
			summary.Different++
		}
		summary.Differences += len(diffs)
		comparisonRows = append(comparisonRows, row)
		differenceRows = append(differenceRows, diffs...)
		if len(comparisonRows) < batchRows {
			continue
		}
		if _, err := comparisons.CreateMany(ctx, comparisonRows); err != nil {
			return service.Summary{}, err
		}
		if _, err := differences.CreateMany(ctx, differenceRows); err != nil {
			return service.Summary{}, err
		}
		comparisonRows = comparisonRows[:0]
		differenceRows = differenceRows[:0]
	}
	if err := pairs.Err(); err != nil {
		return service.Summary{}, err
	}
	if _, err := comparisons.CreateMany(ctx, comparisonRows); err != nil {
		return service.Summary{}, err
	}
	if _, err := differences.CreateMany(ctx, differenceRows); err != nil {
		return service.Summary{}, err
	}

	setAside, err := compare.CreateSetAside(ctx, vendorPrefix, macPrefix, workload)
	if err != nil {
		return service.Summary{}, err
	}
	duplicated, err := compare.CreateDuplicated(ctx, vendorPrefix, macPrefix, workload)
	if err != nil {
		return service.Summary{}, err
	}
	missing, err := compare.CreateMissing(ctx, vendorPrefix, macPrefix, workload)
	if err != nil {
		return service.Summary{}, err
	}
	extra, err := compare.CreateExtra(ctx, vendorPrefix, macPrefix, workload)
	if err != nil {
		return service.Summary{}, err
	}
	summary.SetAside = int(setAside)
	summary.Duplicated = int(duplicated)
	summary.Missing = int(missing)
	summary.Extra = int(extra)
	return summary, nil
}

// Summarize writes a compared vendor submission's summary tables from its
// comparison rows: measure, tag_total, cell, and tag_cell. Each table is
// rebuilt in its own short transaction: delete the submission's rows, then
// one INSERT ... SELECT per row set. A summarized
// submission is first set back to 'compared', which hides it from the views;
// the last step sets 'summarized', which shows it. If a step fails, the
// submission stays hidden and rerunning Summarize is safe.
func (c Controller) Summarize(ctx context.Context, vendorPrefix string) error {
	submissions := submissiontable.NewTableAccess(c.db)
	vendor, err := submissions.Get(ctx, vendorPrefix)
	if err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	if vendor.Side != submissiontable.SideVendor {
		return fmt.Errorf("submission: summarizing %s: not a vendor submission", vendorPrefix)
	}
	if vendor.Status != submissiontable.StatusCompared && vendor.Status != submissiontable.StatusSummarized {
		return fmt.Errorf("submission: summarizing %s: status is %q; compare it first", vendorPrefix, vendor.Status)
	}
	mac, err := submissions.MAC(ctx)
	if err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	if err := submissions.ClearSummarized(ctx, vendor.Prefix); err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}

	if err := c.summarizeMeasures(ctx, vendor.Prefix, mac.Prefix); err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	if err := c.summarizeTagTotals(ctx, vendor.Prefix); err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	if err := c.summarizeCells(ctx, vendor.Prefix); err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	if err := c.summarizeTagCells(ctx, vendor.Prefix); err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	if err := submissions.SetSummarized(ctx, vendor.Prefix); err != nil {
		return fmt.Errorf("submission: summarizing %s: %w", vendorPrefix, err)
	}
	return nil
}

// summarizeMeasures replaces a submission's measure rows in one transaction.
func (c Controller) summarizeMeasures(ctx context.Context, vendorPrefix, macPrefix string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("measures: beginning: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := measure.NewTableAccess(tx).DeleteBySubmission(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := query.NewSummaryQuery(tx).CreateMeasures(ctx, vendorPrefix, macPrefix); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("measures: committing: %w", err)
	}
	return nil
}

// summarizeTagTotals replaces a submission's tag_total rows in one
// transaction.
func (c Controller) summarizeTagTotals(ctx context.Context, vendorPrefix string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tag totals: beginning: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := tagtotal.NewTableAccess(tx).DeleteBySubmission(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := query.NewSummaryQuery(tx).CreateTagTotals(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tag totals: committing: %w", err)
	}
	return nil
}

// summarizeCells replaces a submission's cell rows, the summary cube, in
// one transaction.
func (c Controller) summarizeCells(ctx context.Context, vendorPrefix string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cells: beginning: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := cell.NewTableAccess(tx).DeleteBySubmission(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := query.NewSummaryQuery(tx).CreateCells(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("cells: committing: %w", err)
	}
	return nil
}

// summarizeTagCells replaces a submission's tag_cell rows, the tag cube, in
// one transaction.
func (c Controller) summarizeTagCells(ctx context.Context, vendorPrefix string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tag cells: beginning: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := tagcell.NewTableAccess(tx).DeleteBySubmission(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := query.NewSummaryQuery(tx).CreateTagCells(ctx, vendorPrefix); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tag cells: committing: %w", err)
	}
	return nil
}
