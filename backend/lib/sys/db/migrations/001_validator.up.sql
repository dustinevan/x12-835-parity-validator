-- Version: 1
-- Description: validator schema: submissions, claims, comparisons, differences, summaries
--
-- Database roles. Three login roles, least privilege (docs/DELIVERY.md):
--   validator_owner  owns the tables; runs vdb migrate only
--   validator_batch  reads and writes rows; runs vdb load, compare, summarize
--   validator_read   reads rows; the web service, which never writes
-- A DBA runs this once, as the administrator (the RDS master user):
--
--   CREATE DATABASE validator;
--   \connect validator
--   CREATE ROLE validator_owner LOGIN PASSWORD '...';
--   CREATE ROLE validator_batch LOGIN PASSWORD '...';
--   CREATE ROLE validator_read LOGIN PASSWORD '...';
--   REVOKE ALL ON DATABASE validator FROM PUBLIC;
--   GRANT CONNECT ON DATABASE validator TO validator_owner, validator_batch, validator_read;
--   REVOKE ALL ON SCHEMA public FROM PUBLIC;
--   GRANT USAGE, CREATE ON SCHEMA public TO validator_owner;
--   GRANT USAGE ON SCHEMA public TO validator_batch, validator_read;
--   ALTER ROLE validator_read SET default_transaction_read_only = on;
--
-- then, after the first vdb migrate up, once more as validator_owner:
--
--   GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO validator_batch;
--   REVOKE INSERT, UPDATE, DELETE ON public.migrations FROM validator_batch;
--   GRANT SELECT ON ALL TABLES IN SCHEMA public TO validator_read;
--   ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO validator_batch;
--   ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO validator_read;
--
-- The second part waits for the migrations table, which vdb creates, so the
-- batch role can read it but not change it. The default privileges cover
-- tables that later migrations create.
--
-- Pipeline stages. The input files (S3, or the local input directory) are
-- the golden source and are only ever read. Every table is derived from
-- them by one stage, and each derived table has exactly one writer, its
-- stage, triggered by its submission reaching that stage:
--
--   stage      writes                                   runs when
--   load       submission, claim, claim_repeat, payment  a submission's files arrive
--              and mac_cell (the MAC submission only)
--   compare    comparison, difference                   a vendor submission is loaded
--   summarize  measure, tag_total, cell, tag_cell       a vendor submission is compared
--
-- Nothing edits a derived row afterward. To recompute, rerun the stage: it
-- deletes the submission's rows and writes them again. Compare and summarize
-- also advance submission.status. The one exception is load replacing a
-- submission's files: it removes every row derived from the earlier files,
-- in every stage's tables, because those files no longer exist.
--
-- submission.status: loaded -> compared -> summarized, or failed.
-- - loaded: the submission and its claims are written. The MAC submission
--   stops here; it has nothing to compare.
-- - compared: the comparison and difference rows are written by MAC
--   workload (the match key contains the workload, so no key crosses one),
--   a set number of workloads at once (vdb's APP_SERVER_COMPARE_WORKLOADS,
--   default 1), every statement on its own with no transaction held open:
--   the submission's old rows are deleted first, then per workload the
--   pairs stream out of a join and are compared one at a time, and the
--   set-aside, duplicated, missing, and extra rows are written by
--   INSERT ... SELECT. The submission is first set back to 'loaded' with
--   no compared set, so the views hide it and summarize refuses it before
--   its rows change (a run that dies leaves it there); the status is set
--   to 'compared', with the compared set, after the last workload.
-- - summarized: summarize rebuilt each summary table in its own short
--   transaction (delete the submission's rows, one INSERT ... SELECT). It
--   first sets a summarized submission back to 'compared', so a rerun hides
--   the submission until every table is rebuilt.
-- - failed: a stage could not finish. Compare sets it when a workload
--   fails, leaving that workload partly written; summarize refuses it. A
--   failed summarize leaves 'compared'. Rerunning either stage is safe:
--   compare deletes the submission's rows before writing them again.
-- The views show the MAC submission and the vendor submissions whose status is
-- 'summarized'. That status gate, not one long transaction, keeps a
-- half-finished submission off the screens.
--
-- Designed from the views (front end first):
-- - comparison is the view-shaped table: one row per (vendor submission,
--   match key). The claim lists and one-claim views read it.
-- - claim holds each claim's bytes, for the one-claim view and as the
--   comparison's input. The MAC-processed files are a submission like any other.
-- - difference holds each differing element with both values, for the
--   claim page.
-- - payment holds each 835 transaction's header, so a result can be cut by
--   the MAC's payment context (payer, payee, payment date, method) without
--   reading the files again.
-- - measure, tag_total, cell, and tag_cell hold what a page would
--   otherwise compute by scanning a whole submission on load, so the headline
--   views read a few rows even at 100 million claims.
--
-- Principles:
-- - Keys are facts anyone can recompute: a submission's prefix, a file name
--   and positions within it, and the match key (SHA-256 of the six key
--   fields, which are also stored readably). No generated IDs.
-- - No foreign keys: each submission is loaded in a fixed order by one process.
-- - Column names use the X12 005010X221A1 element names or the CMS term.

-- Written by load. Compare sets status, compared_at, and compared_set;
-- summarize sets status.
CREATE TABLE IF NOT EXISTS public.submission (
    prefix              TEXT PRIMARY KEY,   -- where the files arrived: 'mac', 'vendor-1/2026-10-02', or s3://...
    side                TEXT NOT NULL,      -- 'mac' | 'vendor'
    label               TEXT NOT NULL,      -- opaque: 'MAC', 'Vendor 1'; never a vendor's name
    submitted_on        DATE,               -- NULL for the MAC submission
    status              TEXT NOT NULL,      -- 'loaded' | 'compared' | 'summarized' | 'failed'
    claim_count         BIGINT NOT NULL DEFAULT 0,
    loaded_at           TIMESTAMPTZ NOT NULL,
    compared_at         TIMESTAMPTZ,
    compared_set        TEXT                -- service.CompareSet.Name of the compare run; NULL until compared
);

-- Written by load. One row per claim (loop 2100 with its 2110 service
-- lines) per submission.
CREATE TABLE IF NOT EXISTS public.claim (
    submission_prefix                       TEXT NOT NULL,
    match_key                               BYTEA NOT NULL,   -- SHA-256 of the six key fields below
    file_name                               TEXT NOT NULL,    -- relative to the submission prefix
    payment_number                          INT NOT NULL,     -- 1st, 2nd... transaction set (ST..SE) in the file
    claim_number                            INT NOT NULL,     -- 1st, 2nd... claim in the payment
    medicare_beneficiary_identifier         TEXT NOT NULL,    -- key: 2100 NM1*QC NM109
    date_of_service                         DATE NOT NULL,    -- key: DTM*232, else earliest 2110 DTM*472
    total_claim_charge_amount               NUMERIC NOT NULL, -- key: CLP03
    federal_taxpayer_identification_number  TEXT NOT NULL,    -- key: 1000B REF*TJ
    payee_identification_code               TEXT NOT NULL,    -- key: 1000B N1*PE N104 (NPI)
    mac_workload_id                         TEXT NOT NULL,    -- key: 1000A N1*PR N104
    patient_control_number                  TEXT NOT NULL,    -- CLP01, display only
    claim_type                              TEXT NOT NULL,    -- 'institutional' | 'professional' | 'dme'
    claim_status_code                       TEXT NOT NULL,    -- CLP02
    claim_payment_amount                    NUMERIC NOT NULL, -- CLP04
    payer_claim_control_number              TEXT NOT NULL,    -- CLP07, each system's own; display only
    claim_frequency_code                    TEXT NOT NULL DEFAULT '', -- CLP09, as sent; '' when not sent. A cube dimension
    adjudication                            TEXT NOT NULL DEFAULT '', -- 'inpatient' when the claim carries MIA, 'outpatient' when MOA, '' otherwise. A cube dimension
    jurisdiction                            TEXT NOT NULL DEFAULT '', -- the MAC jurisdiction of mac_workload_id (backend/pkg/macs): 'JE', 'JA'...; '' when the workload is not in the table
    part                                    TEXT NOT NULL DEFAULT '', -- 'A', 'B', or 'DME' from the same lookup; '' when unknown
    delimiters                              TEXT NOT NULL,    -- element, component, repetition, segment terminator
    segments                                BYTEA NOT NULL,   -- the claim's bytes exactly as received
    PRIMARY KEY (submission_prefix, match_key, file_name, payment_number, claim_number)
);

-- Written by load, after the submission's claims: one row per match key
-- that occurs more than once in the submission, with how many claims carry
-- it. Found once, at load, so compare never counts claims per key again.
-- Compare reads it to sort the keys: no row on either side, compared as a
-- pair; a row on the MAC side, set aside; a row on the vendor side only,
-- duplicated.
CREATE TABLE IF NOT EXISTS public.claim_repeat (
    submission_prefix  TEXT NOT NULL,
    match_key          BYTEA NOT NULL,
    claims             INT NOT NULL,     -- claims of the submission carrying the key; always more than 1
    PRIMARY KEY (submission_prefix, match_key)
);

-- Written by load. One row per 835 transaction (ST..SE) per submission: the
-- payment header the transaction's claims sit under, so results can be cut
-- by the MAC's payment context. A comparison groups by the MAC claim's
-- transaction (each pair has exactly one MAC claim in exactly one MAC
-- transaction; the claim's file_name and payment_number name the row). A
-- vendor submission's rows serve only the vendor's own per-transaction
-- checks and its extra claims. total_paid is stored for a balance check of
-- one side against itself (total_paid = the sum of the transaction's
-- claims' claim_payment_amount - plb_total); it is never compared between
-- sides.
CREATE TABLE IF NOT EXISTS public.payment (
    submission_prefix  TEXT NOT NULL,
    file_name          TEXT NOT NULL,     -- relative to the submission prefix
    payment_number     INT NOT NULL,      -- 1st, 2nd... transaction set (ST..SE) in the file
    payer_identifier   TEXT NOT NULL,     -- 1000A N1*PR N104 (the MAC workload)
    payer_name         TEXT NOT NULL,     -- 1000A N1*PR N102
    payee_identifier   TEXT NOT NULL,     -- 1000B N1*PE N104 (NPI)
    payee_name         TEXT NOT NULL,     -- 1000B N1*PE N102
    payee_tax_id       TEXT,              -- REF*TJ (the last of the header's, payer's, and payee's), the value the claims' federal_taxpayer_identification_number carries
    payee_state        TEXT,              -- 1000B N4 N402
    payment_date       DATE,              -- BPR16, check issue or EFT effective date
    payment_method     TEXT,              -- BPR04: ACH, CHK, NON...
    handling_code      TEXT,              -- BPR01: I, H, C...
    trace_number       TEXT,              -- TRN02, check or EFT trace number
    production_date    DATE,              -- DTM*405
    total_paid         NUMERIC,           -- BPR02, total actual provider payment; one side's own balance check only, never compared between sides
    claim_count        INT NOT NULL,      -- CLP segments in the transaction
    plb_count          INT NOT NULL,      -- PLB segments in the transaction
    plb_total          NUMERIC,           -- the sum of the PLB adjustment amounts sent (PLB04, PLB06... PLB14); 0 with no PLB
    PRIMARY KEY (submission_prefix, file_name, payment_number)
);

-- Written by load, for the MAC submission only, after its claims: the
-- legacy cube. One row per distinct combination of the nine dimensions the
-- MAC claim itself carries, with the claims and what was paid and billed,
-- so the front page shows the legacy side's own figures (/api/front's
-- legacy block) whatever a vendor did with its claims. It counts claims,
-- not match keys: a key the MAC repeats counts each of its claims. The
-- dimensions read as in cell: paid_band from claim_payment_amount ('zero',
-- 'under 100', 'under 1000', 'under 10000', '10000 and more'),
-- service_month the first day of the month of date_of_service, and every
-- dimension text and never NULL. Load deletes and rewrites the rows with
-- the submission's other rows.
CREATE TABLE IF NOT EXISTS public.mac_cell (
    submission_prefix     TEXT NOT NULL,
    jurisdiction          TEXT NOT NULL,     -- 'JE', 'JA'...; '' unknown
    part                  TEXT NOT NULL,     -- 'A' | 'B' | 'DME' | ''
    mac_workload_id       TEXT NOT NULL,
    claim_type            TEXT NOT NULL,     -- 'institutional' | 'professional' | 'dme'
    adjudication          TEXT NOT NULL,     -- 'inpatient' | 'outpatient' | ''
    claim_status_code     TEXT NOT NULL,     -- CLP02
    claim_frequency_code  TEXT NOT NULL,     -- CLP09; '' when not sent
    paid_band             TEXT NOT NULL,
    service_month         DATE NOT NULL,
    claims                BIGINT NOT NULL,   -- MAC claims in the cell
    paid                  NUMERIC NOT NULL,  -- the sum of CLP04
    billed                NUMERIC NOT NULL,  -- the sum of CLP03
    PRIMARY KEY (submission_prefix, jurisdiction, part, mac_workload_id, claim_type, adjudication, claim_status_code, claim_frequency_code, paid_band, service_month)
);

-- Written by compare. One row per (vendor submission, match key): the result
-- the claim lists and one-claim views read.
CREATE TABLE IF NOT EXISTS public.comparison (
    submission_prefix                TEXT NOT NULL,
    match_key                        BYTEA NOT NULL,
    outcome                          TEXT NOT NULL,     -- 'equal' | 'different' | 'missing' | 'duplicated' | 'extra' | 'set_aside'
    mac_workload_id                  TEXT NOT NULL,
    claim_type                       TEXT NOT NULL,
    medicare_beneficiary_identifier  TEXT NOT NULL,
    date_of_service                  DATE NOT NULL,
    total_claim_charge_amount        NUMERIC NOT NULL,
    mac_claim_status_code            TEXT,              -- NULL when extra
    vendor_claim_status_code         TEXT,              -- NULL when missing, duplicated, or set aside
    mac_claim_payment_amount         NUMERIC,
    vendor_claim_payment_amount      NUMERIC,
    claim_frequency_code             TEXT NOT NULL DEFAULT '', -- the four below are copied from the MAC claim, or from the vendor claim when extra: the attributes the summary cube groups by
    adjudication                     TEXT NOT NULL DEFAULT '', -- 'inpatient' | 'outpatient' | ''
    jurisdiction                     TEXT NOT NULL DEFAULT '',
    part                             TEXT NOT NULL DEFAULT '',
    difference_count                 INT NOT NULL DEFAULT 0,
    tags                             TEXT[] NOT NULL DEFAULT '{}', -- what differs, e.g. {CLP04,CAS*CO*253}
    PRIMARY KEY (submission_prefix, match_key)
);

-- Written by compare. One row per differing element in a compared claim
-- pair.
CREATE TABLE IF NOT EXISTS public.difference (
    submission_prefix  TEXT NOT NULL,
    match_key          BYTEA NOT NULL,
    element            TEXT NOT NULL,   -- 'CLP04', 'SVC03', 'CAS*CO*253', 'MOA', 'AMT*B6'...
    occurrence         INT NOT NULL,    -- 1st, 2nd... of the same element in the claim
    line_control_number TEXT,           -- REF*6R of the service line, for line-level elements
    mac_value          TEXT,            -- NULL when absent on the MAC side
    vendor_value       TEXT,            -- NULL when absent on the vendor side
    amount_difference  NUMERIC,         -- vendor minus MAC; numeric elements only
    PRIMARY KEY (submission_prefix, match_key, element, occurrence)
);

-- Written by summarize. A vendor submission's headline measures, by claim
-- type and workload; 'all' rows hold the totals. Read by the overview,
-- workloads, differences, and unmatched views.
CREATE TABLE IF NOT EXISTS public.measure (
    submission_prefix        TEXT NOT NULL,
    claim_type               TEXT NOT NULL,     -- 'all' | 'institutional' | 'professional' | 'dme'
    mac_workload_id          TEXT NOT NULL,     -- 'all' | a workload
    mac_claims               BIGINT NOT NULL,   -- MAC claims in scope (the expected count)
    equal_count              BIGINT NOT NULL,
    different_count          BIGINT NOT NULL,
    missing_count            BIGINT NOT NULL,
    extra_count              BIGINT NOT NULL,
    set_aside_count          BIGINT NOT NULL,   -- MAC claims whose key the MAC repeats
    duplicated_count         BIGINT NOT NULL,   -- MAC claims whose key the vendor sent more than once (one claim per key)
    outcome_match_rate       NUMERIC,           -- same CLP02 among matched pairs ÷ MAC claims not set aside
    paid_claims              BIGINT NOT NULL,   -- matched pairs where the MAC paid (> 0)
    median_deviation         NUMERIC,           -- |vendor − MAC paid| ÷ MAC paid, over paid_claims
    p95_deviation            NUMERIC,
    PRIMARY KEY (submission_prefix, claim_type, mac_workload_id)
);

-- Written by summarize. For each tag of a vendor submission's different
-- claims: how many carry it and the dollars of that element, by claim type
-- and workload; 'all' rows hold the totals, as in measure. Read by the
-- differences view (its rows and each row's previous_claims).
CREATE TABLE IF NOT EXISTS public.tag_total (
    submission_prefix  TEXT NOT NULL,
    claim_type         TEXT NOT NULL,     -- 'all' | 'institutional' | 'professional' | 'dme'
    mac_workload_id    TEXT NOT NULL,     -- 'all' | a workload
    tag                TEXT NOT NULL,     -- an element, as in comparison.tags
    claims             BIGINT NOT NULL,   -- different claims whose tags contain the tag
    dollars            NUMERIC NOT NULL,  -- sum of |amount_difference| of the element's difference rows in those claims
    PRIMARY KEY (submission_prefix, claim_type, mac_workload_id, tag)
);

-- Written by summarize. The summary cube: one row per distinct combination
-- of the twelve dimensions below among a vendor submission's comparison
-- rows, with counts and sums. The front page reads it with WHERE filters
-- and SUM, so a reader's selection never groups the comparison rows again.
-- Every dimension is text and never NULL: '' means not applicable (a
-- vendor status on a missing row, a payment band on anything but a pair).
-- paid_band is the MAC paid amount (the vendor's when extra): 'zero',
-- 'under 100', 'under 1000', 'under 10000', '10000 and more'.
-- payment_band, on pairs only, is |vendor paid - MAC paid|: 'match' within
-- five cents (the comparison's five-cent rule), 'up to 1.00', 'up to 1%'
-- and 'up to 5%' of the MAC paid amount when it is more than zero, else
-- 'over 5%'; 'match' and 'up to' include their edges.
-- service_month is the first day of the month of date_of_service.
-- The sums: mac_paid, vendor_paid, and billed (CLP03) add every row with
-- an absent amount as zero; net_difference (vendor minus MAC) and
-- gross_difference (|vendor minus MAC|) add pairs only; same_status counts
-- pairs whose CLP02 agree. No rate is stored.
CREATE TABLE IF NOT EXISTS public.cell (
    submission_prefix         TEXT NOT NULL,
    jurisdiction              TEXT NOT NULL,     -- 'JE', 'JA'...; '' unknown
    part                      TEXT NOT NULL,     -- 'A' | 'B' | 'DME' | ''
    mac_workload_id           TEXT NOT NULL,
    claim_type                TEXT NOT NULL,     -- 'institutional' | 'professional' | 'dme'
    adjudication              TEXT NOT NULL,     -- 'inpatient' | 'outpatient' | ''
    mac_claim_status_code     TEXT NOT NULL,     -- CLP02; '' when extra
    vendor_claim_status_code  TEXT NOT NULL,     -- CLP02; '' when missing, duplicated, or set aside
    claim_frequency_code      TEXT NOT NULL,     -- CLP09; '' when not sent
    paid_band                 TEXT NOT NULL,
    service_month             DATE NOT NULL,
    payment_band              TEXT NOT NULL,     -- '' when not a pair
    outcome                   TEXT NOT NULL,     -- as comparison.outcome
    claims                    BIGINT NOT NULL,   -- comparison rows (keys) in the cell
    mac_paid                  NUMERIC NOT NULL,
    vendor_paid               NUMERIC NOT NULL,
    billed                    NUMERIC NOT NULL,
    net_difference            NUMERIC NOT NULL,
    gross_difference          NUMERIC NOT NULL,
    same_status               BIGINT NOT NULL,
    PRIMARY KEY (submission_prefix, jurisdiction, part, mac_workload_id, claim_type, adjudication, mac_claim_status_code, vendor_claim_status_code, claim_frequency_code, paid_band, service_month, payment_band, outcome)
);

-- Written by summarize. The same cube for what differs: per cell and tag
-- (an element in comparison.tags), the different pairs carrying the tag
-- and the dollars of that element in them (the sum of |amount_difference|
-- of its difference rows, as tag_total computes it). Only different pairs
-- have tags, so outcome is always 'different' here.
CREATE TABLE IF NOT EXISTS public.tag_cell (
    submission_prefix         TEXT NOT NULL,
    jurisdiction              TEXT NOT NULL,
    part                      TEXT NOT NULL,
    mac_workload_id           TEXT NOT NULL,
    claim_type                TEXT NOT NULL,
    adjudication              TEXT NOT NULL,
    mac_claim_status_code     TEXT NOT NULL,
    vendor_claim_status_code  TEXT NOT NULL,
    claim_frequency_code      TEXT NOT NULL,
    paid_band                 TEXT NOT NULL,
    service_month             DATE NOT NULL,
    payment_band              TEXT NOT NULL,
    outcome                   TEXT NOT NULL,
    tag                       TEXT NOT NULL,     -- an element, as in comparison.tags
    claims                    BIGINT NOT NULL,   -- different pairs in the cell whose tags contain the tag
    dollars                   NUMERIC NOT NULL,  -- sum of |amount_difference| of the element's difference rows in those pairs
    PRIMARY KEY (submission_prefix, jurisdiction, part, mac_workload_id, claim_type, adjudication, mac_claim_status_code, vendor_claim_status_code, claim_frequency_code, paid_band, service_month, payment_band, outcome, tag)
);

CREATE INDEX IF NOT EXISTS comparison_tags_idx ON public.comparison USING GIN (tags);
CREATE INDEX IF NOT EXISTS comparison_outcome_idx ON public.comparison (submission_prefix, outcome);
CREATE INDEX IF NOT EXISTS comparison_match_key_idx ON public.comparison (match_key);
CREATE INDEX IF NOT EXISTS claim_match_key_idx ON public.claim (match_key);
-- Compare reads one submission's claims one MAC workload at a time.
CREATE INDEX IF NOT EXISTS claim_workload_idx ON public.claim (submission_prefix, mac_workload_id);
