# Subtask 265-stdout-ci-gates · 06 — Benchmarks for Hot Paths

Scope: implementer of the benchmark suite for task 265.
Spec home: `02-spec/21-app/265-stdout-ci-gates/04-deprecation-and-pe-parallel.md` (Part B context) and `03-ci-gates.md`.
Rules: files ≤300 lines, lowercase filenames, positive-prefixed booleans (`hasX`, `isX`), `*apperror.AppError` for errors, relative paths only. Benchmark funcs are named `Benchmark*` in `*_test.go` files next to the code under test. Do NOT run `go build`, `go test`, or any git commands — the lead runs the benchmarks and owns commits/releases.

## Background

Only 3 benchmarks exist today (`cli/cmdautomation/search_benchmark_test.go` x2, `cli/fsutil/path_normalize_test.go` x1). Research confirms 1,037 test files / 4,292 test funcs with zero duplicate test names — no obvious breakage from adding new `Benchmark*` funcs.

## Steps

- [ ] `BenchmarkAumSearch` — hot path: the `aum search` engine in `cli/cmdautomation/search.go`. Benchmark a representative query against a fixed fixture corpus (small, generated in-test; no network, no repo-wide walk). Mirror the style of the existing funcs in `search_benchmark_test.go`.
- [ ] `BenchmarkFixScan` — hot path: the `fix` scan in `cli/cmdautofix/fix.go`. Benchmark the scan phase on a synthetic repo tree built in `t.TempDir()`.
- [ ] `BenchmarkPipelineErrorsAll` — hot path: `collectAllPipelineSummary` in `cli/cmdpipeline/pipeline_all_errors.go`. Benchmark against a seeded SQLite fixture (no GitHub API calls; stub the fetch layer).
- [ ] `BenchmarkSpecIssue` — hot path: spec issuance in `cli/cmdspec/spec.go:108`. Benchmark the issuance path with a temp DB.
- [ ] Each benchmark: call `b.ReportAllocs()`, reset the timer after fixture setup (`b.ResetTimer()`), no sleeps, no wall-clock dependence.
- [ ] Lead runs each once and records the baselines in Evidence below. Where the baselines should live long-term (a new `05-benchmarks.md` spec vs inside the Part B spec) is an owner call — do NOT create the spec file unasked; record here first.

## Baselines

(hot path, ns/op, allocs/op — recorded by the lead after the run)
