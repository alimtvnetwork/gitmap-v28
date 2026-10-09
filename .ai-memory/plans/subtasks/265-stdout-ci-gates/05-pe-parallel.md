# Subtask 265-stdout-ci-gates · 05 — pe-all Parallel Fetch

Scope: implementer of the `pipeline errors all` worker pool for task 265.
Spec home: `02-spec/21-app/265-stdout-ci-gates/04-deprecation-and-pe-parallel.md` Part B (read it first).
Rules: files ≤300 lines, lowercase filenames, positive-prefixed booleans (`hasX`, `isX`), `*apperror.AppError` for errors, relative paths only. Do NOT run `go build`, `go test`, or any git commands — the lead verifies the build and owns commits/releases.

## Steps

- [ ] Worker pool in `collectAllPipelineSummary()` (`cli/cmdpipeline/pipeline_all_errors.go`): semaphore-channel pattern mirroring `executeParallelFetchWorkers` in `cli/cmdpipeline/pipeline_logs_fetch.go` (currently cap 8 — use the computed default instead).
- [ ] Worker count: `workers = max(runtime.NumCPU(), 3)`; add a `--workers` flag parsed like `Limit` in `cli/cmdpipeline/pipeline_flags.go` to override it.
- [ ] NO nested parallelism inside workers: one repo per worker, strictly sequential within the worker.
- [ ] Last-commit-only fetch per repo: use the `gh run list --commit <sha> --json` pattern from `cli/cmdpipeline/pipeline_logs_target.go:152`, or direct REST with a Bearer token from `secrets.Resolve()`.
- [ ] SQLite hash cache: before fetching, compare the repo's recorded commit hash in the pipeline DB against the current HEAD hash — match means skip the fetch and mark the repo done; mismatch means fetch + record the new hash.
- [ ] Deterministic output: re-sort results by repo slug after the parallel fetch; guard shared summary counters with a mutex.
- [ ] Keep every touched file ≤300 lines; move formatting helpers to a sibling file if `pipeline_all_errors.go` grows past the limit.
- [ ] Live verification (lead runs): `gitmap pipeline errors all` on a multi-repo catalog — record sequential baseline timing vs parallel timing; confirm `--workers 1` reproduces sequential behavior; confirm a second run with no new commits performs zero fetches (hash-cache hit); confirm output order is slug-sorted and stable across runs.

## Evidence

(baseline vs parallel timings, `--workers 1` check, hash-cache second-run check, output-stability check)
