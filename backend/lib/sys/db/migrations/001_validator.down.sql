-- Version: 1
-- Description: validator schema: submissions, claims, comparisons, differences, summaries
DROP INDEX IF EXISTS public.claim_workload_idx;
DROP TABLE IF EXISTS public.tag_cell;
DROP TABLE IF EXISTS public.cell;
DROP TABLE IF EXISTS public.tag_total;
DROP TABLE IF EXISTS public.measure;
DROP TABLE IF EXISTS public.difference;
DROP TABLE IF EXISTS public.comparison;
DROP TABLE IF EXISTS public.mac_cell;
DROP TABLE IF EXISTS public.payment;
DROP TABLE IF EXISTS public.claim_repeat;
DROP TABLE IF EXISTS public.claim;
DROP TABLE IF EXISTS public.submission;
