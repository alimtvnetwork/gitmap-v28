# Pipeline Historical Success ETA & Error Diagnostics

## Metadata

- **Module**: `spec/21-app/09-pipeline-historical-eta/`
- **Version**: 1.0.0
- **Status**: Stable
- **AI Confidence**: Production-Ready
- **Ambiguity**: None
- **Keywords**: `pipeline`, `eta`, `historical-baseline`, `dynamic-timeout`, `-t`, `ci-cd-errors`, `diagnostics`

## Purpose

This specification defines the architecture, data structures, and operational contracts for GitMap's enhanced CI/CD pipeline model. Key features include:

1. **Historical Success Baseline Algorithm**: Derives ETA approximations strictly from past successful workflow executions (`conclusion == "success"`), computing average run durations per pipeline segment while discarding distorted cancelled or fast-failing runs.
2. **Dynamic Timeline & Timeout (`-t` / `--timeout` / `--timeline`)**: Eliminates manual interval guessing by automatically computing adaptive retry cadences and deadline budgets from the historical ETA.
3. **Targeted CI/CD Error Extraction**: Automatically extracts and displays only relevant error lines and failure diagnostics during pipeline status inspection and wait-time tracking, stripping away passing step noise.

## File Inventory

| File | Purpose |
|---|---|
| [`00-overview.md`](./00-overview.md) | Module entry point, scoring, and file inventory |
| [`01-historical-success-eta-algorithm.md`](./01-historical-success-eta-algorithm.md) | Success-only filtering, segment-level averages, and remaining ETA math |
| [`02-dynamic-timeline-timeout-t-flag.md`](./02-dynamic-timeline-timeout-t-flag.md) | Adaptive polling timeline, timeout heuristics, and `-t` CLI flag |
| [`03-ci-cd-error-extraction-and-diagnostics.md`](./03-ci-cd-error-extraction-and-diagnostics.md) | High-precision failure regex filtering and diagnostic terminal formatting |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Given/When/Then verification contracts and automated test gates |
| [`99-consistency-report.md`](./99-consistency-report.md) | Structural and architectural compliance audit |
