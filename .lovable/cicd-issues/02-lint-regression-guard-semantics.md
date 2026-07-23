# CI/CD Issue 02 — Lint guard semantics unified to baseline-diff

## Pipeline
- **Workflow:** `.github/workflows/ci.yml`
- **Jobs:** `lint-baseline-guard` (formerly `lint-hard-floor`, originally `lint-regression-guard`), `lint-baseline-diff`
- **Scripts:** `.github/scripts/check-single-linter-diff.sh`, `.github/scripts/lint-diff.py`

## Symptom
User asked to "verify the lint regression guard ignores the baseline and fails CI only on newly introduced golangci-lint issues." Verification revealed the contract was **not** uniform — the job mixed two enforcement models (hard-floor for `unused`+`G115`, baseline-diff for the rest).

## Root Cause (historical)
Two distinct enforcement models lived under the "regression guard" umbrella:

| Linter / Check | Original Model | Original Implementation |
|---|---|---|
| `unused` | Hard floor (no baseline) | `check-lint-regressions.sh` |
| `gosec G115` (integer-overflow) | Hard floor (no baseline) | `check-lint-regressions.sh` |
| `misspell` | Baseline diff (new only) | `check-single-linter-diff.sh` |
| `gocritic` | Baseline diff (new only) | `check-single-linter-diff.sh` |
| `exhaustive` | Baseline diff (new only) | `check-single-linter-diff.sh` |
| Full report | Baseline diff | `lint-diff.py` (in `lint-baseline-diff` job) |

Baseline is cached per-linter as `golangci-<linter>-baseline-main-<sha>`, refreshed only on successful pushes to `main`.

## Status
✅ Resolved (2026-04-29, two passes).

**Pass 1** — User initially chose Option (b) (rename only). Job renamed `lint-regression-guard` → `lint-hard-floor`; comment block updated to clarify the split enforcement model.

**Pass 2** — User then chose Option (a) (convert `unused` + `G115` to baseline-diff to match the rest). Final state:

| Change | Detail |
|---|---|
| Job renamed | `lint-hard-floor` → `lint-baseline-guard` |
| Job label | `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)` |
| Script removed | `.github/scripts/check-lint-regressions.sh` (deleted; the hard-floor mechanism no longer exists) |
| New sub-steps | `unused` and `gosec G115` now use `check-single-linter-diff.sh` with rolling per-linter caches keyed by SHA, restored via `golangci-<name>-baseline-main-` prefix |
| Script extended | `check-single-linter-diff.sh` now accepts `TEXT_FILTER` (regex on `.Text`, used to scope `gosec` to G115) and `LABEL` (display name override). Filter applied uniformly to current AND baseline so the diff stays apples-to-apples |
| Aggregator updated | `test-summary` job's `needs:` now includes `lint-baseline-guard` so the rename surfaces in the required-check chain |

## Resolution Path (historical)
Two options were presented:
- **(a)** Convert `unused` + `G115` to baseline-diff semantics (any new issue fails; existing baseline tolerated). ← **Selected (Pass 2).**
- **(b)** Rename `lint-regression-guard` → `lint-hard-floor` so the job name reflects zero-tolerance enforcement and stops implying baseline-diff. ← Selected (Pass 1), then superseded.

## Prevention
- All five linters in `lint-baseline-guard` now share one model (baseline-diff) and one script (`check-single-linter-diff.sh`). New linters added to this job MUST use the same contract — no more hard-floor carve-outs.
- If a hard-floor is ever genuinely needed, give it its own dedicated job (do NOT mix models inside one job again).
- When scoping a multi-rule analyzer (like `gosec`) to a single rule, use `TEXT_FILTER` env var rather than custom shell post-processing.

## Branch-protection follow-up
The required-check name in branch protection rules (if any) needs updating from `Lint Regression Guard (unused + G115)` (or the intermediate `Lint Hard Floor …`) to **`Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)`**. The old check name will not appear after this change merges.

## Related
- Session memory: `.lovable/memory/workflow/04-ci-hardening-session.md`
- Architecture log: `.lovable/memory/tech/ci-pipeline-architecture.md`
