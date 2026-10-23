# Subtask 311-01: CI/CD Diagnosis

## Objective

Produce a root-cause entry for every failing check at the current HEAD of
main, with evidence, before any fix is attempted.

## Steps

1. Fetch GitHub check-runs for HEAD via the API; list every failing check
   (the failure set is moving because the spec-310 stream is active — do not
   rely on the f723012c5 baseline alone).
2. For each failing check, fetch its annotations. For checks whose
   annotations are generic (Test (unit), go build + test, go test -race,
   Lint Script Unit Tests), narrow the failing package/test from the job
   definition in .github/workflows/ and the repo layout.
3. Record each root cause in
   02-spec/21-app/311-cicd-issues-and-end-to-end-testing/02-component-spec.md
   (failure map table) and mark this subtask's table below.

## Root-Cause Log

Baseline re-verified at HEAD abb35f01d (GitHub check-runs: 51 runs, same 11
failures as f723012c5; E2E x3, Test integration/store/tui, Heavy Tests,
Installer Smoke x2, Lint, File Size Gate all green).

| Check | Root cause | Evidence | Fix owner |
|---|---|---|---|
| Lint Baseline Guard / Diff | NEW unused funcs vs baseline: recordVacuumTelemetry, handleVacuumError, executeVacuumOnDB, vacuumSingleTarget (cli/cmddb/db_vacuum.go), requireOnlineCheck (cli/cmd/helpcheck.go:182), hasCredentialPersistenceFailure (cli/cloner/safe_pull.go:135), formatFolderFileURI (cli/cmdnodes/nodes_deploy_ide.go:217), hasSkipIDEsFlag (cli/cmdnodes/nodes_clone_types.go:78), parseGitHubSlugFromURL (cli/cmdpe/pe_gh.go:96); misspell DOUB at cli/cmddb/rest_swagger.go:231 (a strings.Contains type-prefix probe, sibling code in cli/cmdautomation/db_generate.go uses the full "DOUBLE"). RESOLVED in 311-02: all 9 funcs verified callerless via gitmap aum search and deleted (incl. the apperror import that only the dead helpcheck func used); DOUB -> DOUBLE (same behavior for real SQL type names) | Check-run annotations at abb35f01d + aum caller verification (15 hits = definitions/internal cluster calls only) | 311-02 DONE |
| Legacy Refs Check | `gitmap-v6` references at 02-spec/01-spec-authoring-guide/13-root-readme-conventions.md lines 20, 257 matching /gitmap-v[567]\b/. RESOLVED in 311-02: both now reference gitmap-v28 | Check-run annotations | 311-02 DONE |
| Full Suite Guard | Aggregator of the other failures (its annotations list the same lint findings) | Check-run annotations | clears with constituents |
| Test (unit) | Runtime test failures in unit scope (./clonenext/... ./cloner/... ./cmd/... ./config/... ./formatter/... ./scanpipe/... ./release/...); compile gate green, integration/store/tui legs green. CI logs are auth-walled and annotations generic, so failing tests are being named by a scoped local diagnostic run (diagnosis only; CI remains the verification gate) | ci.yml test matrix; diagnostic run in flight | 311-02 |
| go build + test x3 | Same unit failures (job runs go test -short ./... on all OSes) | cross-platform.yml | 311-02 |
| go test -race | Hot packages ./cmd/... ./cloneconcurrency/... ./visibility/... ./store/... ./uipref/... — expected to share the cmd/... unit failures; race-specific findings TBD from the diagnostic | race-detector.yml | 311-02 |
| Lint Script Unit Tests | Runs .github/scripts/tests/test_ci_scripts.py; failing case being identified by a local run of that suite | ci.yml lint-script-tests job | 311-02 |
| Relative Path Check | Checker run against a clean HEAD worktree: exactly 3 violations, all in 02-spec/21-app/310-relative-paths/01-architecture-spec.md lines 16/18/20 — the spec-310 spec quotes the absolute path it exists to eliminate (`d:\work\gitmap\...` inside @[...] references). RESOLVED in 311-02: quoted prefix replaced with `<repo-root>\`; checker re-run on the full working tree: PASS (9,836 files, 0 violations) | linter-scripts/check-relative-paths.py run on HEAD worktree + working tree | 311-02 DONE |
| Test (unit) — working-tree signal | Full unit scope PASSES against the working tree (17 packages ok, incl. cli/cmd in 22.5s): the committed-tree failures are either fixed by the spec-310 stream's uncommitted edits or are CI-environment-specific. Discriminator: same scope running against a clean HEAD worktree | Local diagnostic run (diagnosis only; CI remains the gate) | 311-01 |
| Test (unit) / go build+test x3 / race — ROOT CAUSE | Clean HEAD-worktree run of the unit scope: exactly 4 failing test functions. (1-3) TestPipelineCommandInvocations, TestTopLevelErrorLogsAndLogs, TestTopLevelPipelineErrorsShortcut (cli/cmd/cmd_invocation_test.go:63-127): the tests asserted `gitmap pe` error-logs succeeds unconditionally; since the pe no-data fix (db665c939) the invocation correctly fails with typed E9000 "pipeline status unknown" whenever gh is unauthenticated (CI runners included) — the tests still encode the pre-fix fabricated-success contract. (4) TestHandleOrphanedMetaDecline (cli/release/workflow_test.go:80): asserts err.Error() equals the bare message, but handleOrphanedMeta now returns the typed apperror envelope ("[E9000:EXECUTION] execution: release aborted by user (at=...)") from the repo-wide fmt.Errorf migration. The same pe tests sit in the race scope (./cmd/...) and in the -short ./... cross-platform scope, so one root cause explains Test (unit), all three go build+test legs, and the race job (pending race-job confirmation; no DATA RACE lines seen in any local run). RESOLUTION: test-contract updates already authored in-tree by the parallel repair stream (assertPipelineLogsOutcome helper in cmd_invocation_test.go + pipeline_errorlogs_test.go; typed-envelope assertions in workflow_test.go) — the working tree passes the full unit scope with them. They are committed with task 311 (disclosed in the release notes) because a green release commit is impossible without them | HEAD-worktree diagnostic run (full output retained) + in-tree diff review | 311-02 (commit) |
| Lint Script Unit Tests | HEAD-worktree run in progress (with go on PATH, unlike the earlier local run whose 4 failures included the local-only 'go not found' artifact in the installer-smoke test). Known at HEAD: the legacy-refs script test must fail on the committed gitmap-v6 references — fixed by the 311-02 doc edit | test_ci_scripts.py inventory + local run | 311-01/02 |

## Exit Criteria

- Every failing check at HEAD has a named root cause and a fix route.
- No check is written off as "flake" without a re-run proving it.
