# Component Spec: 311-cicd-issues-and-end-to-end-testing

## 1. CI/CD Failure Map (baseline at abb35f01d — CONFIRMED root causes)

Full diagnosis with evidence: .ai-memory/plans/subtasks/311-cicd-issues-and-end-to-end-testing/01-ci-diagnosis.md

| Failing check | Confirmed root cause | Fix (all applied) |
|---|---|---|
| Lint Baseline Guard / Diff | 9 callerless functions flagged NEW vs baseline; misspell DOUB at cli/cmddb/rest_swagger.go:231 (a SQL type-prefix probe) | Deleted all 9 after aum caller verification (plus the apperror import only the dead helpcheck function used); DOUB -> DOUBLE, matching the sibling code in cli/cmdautomation/db_generate.go |
| Legacy Refs Check | legacy version refs at 02-spec/01-spec-authoring-guide/13-root-readme-conventions.md lines 20, 257 | Both now reference gitmap-v28 |
| Relative Path Check | 3 violations, all inside 02-spec/21-app/310-relative-paths/01-architecture-spec.md quoting the absolute path it exists to eliminate | Quoted prefix replaced with <repo-root>\ form; checker re-run on the full tree: PASS (9,836 files) |
| Test (unit), go build + test x3, go test -race | 4 test functions: 3 cli/cmd pipeline tests asserted the pre-fix fabricated pe success (db665c939 made no-data correctly fail with E9000 "pipeline status unknown"; CI runners have no gh auth, so the old assertions cannot hold); TestHandleOrphanedMetaDecline asserted a bare error string vs the typed apperror envelope | Test-contract updates authored in-tree by the parallel repair stream (assertPipelineLogsOutcome; typed-envelope assertions), committed with this task and disclosed in the release notes. Full unit scope passes with them |
| Lint Script Unit Tests | 3 of 19 tests: legacy-refs (above), gofmt (13 committed .go files unformatted), smoke-installer E2E (fix-git --json expected exit 0 outside a git repo, where the product correctly returns E1000 validation) | Doc fix; gofmt applied to the 13 files; the E2E case now expects [0, 1], the convention the same table uses for other repo-dependent commands |
| Full Suite Guard | Aggregator over the above | Clears when constituents pass |

## 2. End-to-End Testing of All Methods

Coverage map (harness -> what it exercises -> CI job):

- .github/scripts/e2e-cli-smoke.py — 118 end-to-end cases over 74 distinct
  first-token commands (functional invocations + help surfaces + aliases),
  run against a real built binary: E2E jobs on linux/windows/macos
  (ci.yml), Installer Smoke source mode (which embeds the same suite), and
  the Lint Script Unit Tests job. Baseline state: 117/118 — the fix-git
  --json expectation fixed by this task (now [0, 1]).
- cli/cmd/cmd_invocation_test.go — in-process end-to-end invocations of
  the install / pipeline / pe command families through the real dispatcher
  (Test (unit) + race + cross-platform jobs).
- cli/cmd/commitin/e2e — the commit-in pipeline end-to-end against scratch
  repositories (Test (unit), Test (integration) jobs).
- Installer Smoke (linux + windows) — builds from source, installs, and
  verifies the installed binary reports the expected version, then runs
  the full E2E CLI suite against it.
- history-rewrite-smoke + goreleaser-smoke workflows — the history
  rewrite and release-packaging paths end-to-end.
- The remaining registered surface (403 helpdoc entries, mostly
  subcommands and aliases of the covered families) is exercised through
  the per-package unit/integration suites and the help-surface smoke
  cases; the help system itself is guarded by the helpcheck invariants
  in the unit suite (every registered command must render help).

Deliverable evidence: CI run for the release commit with every smoke/e2e
job green (section 5).

## 3. Installation Download Speed

Profiled path and what shipped (detail: subtask 04):

- Blocking correctness fix: release.yml tars the Unix binary under its
  versioned asset name, but extractTarGzBinary only accepted entries
  named gitmap/gitmap.exe — every tar-based self-update failed with
  E1148. isReleaseBinaryName now also accepts gitmap-* entries.
- Concurrency (the guaranteed win): checksums.txt is fetched concurrently
  with the archive in the Go self-update (downloadReleaseStaging) and in
  install.sh download_asset (both mirrored copies), removing one full
  serial round trip from every install/update. The download engine used
  is recorded in update telemetry (UpdateResult.DownloaderEngine).
- Segmented download: implemented in cli/cmdupdate/update_download_parallel.go
  (range probe + 4 concurrent Range segments via WriteAt + aggregated
  progress + single-stream fallback), but OPT-IN via
  GITMAP_DOWNLOAD_SEGMENTS=1. Measured on this repo's network path
  (v6.590.0 linux-amd64 archive, 18,457,710 bytes, 3 runs each):
  single stream median 3,270 ms (1,862 / 3,270 / 8,340);
  4 segments median 9,603 ms (9,066 / 9,603 / 10,461) — parallel
  connections are ~3x SLOWER through this proxied path, so default-on
  would have made updates slower here. Segmented output verified
  byte-identical (sha256 match). On per-connection-capped links it
  remains available as an opt-in; the installer's aria2c path (16
  connections) already covers that case when aria2c is installed.
- install.ps1 concurrency deferred: Start-Job process startup rivals the
  single RTT it would save for a ~1 KB checksums file (subtask 04).
- Regression gates: Installer Smoke (linux + windows), Heavy Tests.

## 4. Release

- Bump: minor increment of version.json and every synced version file
  (current version to be read from version.json at execution time;
  v6.590.0 -> v6.591.0 at spec time).
- Tag: annotated tag for the new version, pushed (tag-based release
  convention on this machine; no GitHub token available locally).
- Release notes: a release-notes file for the new version listing the
  complete change set of the release (CI fixes, e2e coverage, download-speed
  improvement, and the pe false-green fix from db665c939/f723012c5 that this
  release carries). Relative paths only — no absolute paths, no filesystem
  URLs.

## 5. Evidence

- Baseline check-runs (abb35f01d): 51 runs, 11 failures — Test (unit),
  Full Suite Guard, Legacy Refs Check, Relative Path Check, Lint
  Baseline Diff, Lint Baseline Guard, Lint Script Unit Tests,
  go build+test (ubuntu/macos/windows), go test -race (hot packages).
- Download timing before: single-stream full-archive fetch, median
  3,270 ms (3 runs, section 3); checksums fetched serially afterwards
  (one extra RTT + transfer on every install/update).
- Download timing after: checksums fetch fully overlapped with the
  archive fetch in both the Go updater and install.sh; segmented mode
  measured but left opt-in per the section 3 data.
- Method coverage list: section 2.
- Post-fix check-runs (release commit): filled at release verification.
