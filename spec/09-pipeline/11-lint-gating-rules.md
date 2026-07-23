# 11 — Lint Gating Rules

How CI decides which `golangci-lint` findings block a merge and which are tolerated.

---

## TL;DR

| Linter | Gating Model | What fails CI |
|---|---|---|
| `unused` | Baseline-diff | Only NEW dead code introduced by this PR |
| `gosec G115` (integer overflow) | Baseline-diff | Only NEW G115 findings introduced by this PR |
| `misspell` | Baseline-diff | Only NEW misspellings introduced by this PR |
| `gocritic` | Baseline-diff | Only NEW gocritic advisories introduced by this PR |
| `exhaustive` | Baseline-diff | Only NEW missing-case warnings introduced by this PR |
| Everything else (full report) | Baseline-diff | Only NEW findings vs. cached main baseline |
| Pre-existing findings (any linter) | **Tolerated** | Never block CI |

**There are currently NO hard-gated (zero-tolerance) linters.** The hard-floor mechanism was removed on 2026-04-29 in favor of a uniform baseline-diff contract. See [History](#history) for why.

---

## What "baseline-diff" means

For each guarded linter:

1. CI runs the linter against the current tree → `current.json`.
2. CI restores the most recent main-branch baseline from cache → `baseline.json`.
3. CI computes `current` minus `baseline` (keyed by `file|line|text`).
4. **Only findings present in current but not in baseline fail the build.**
5. On successful pushes to `main`, the new `current.json` is saved as the next baseline.

Pre-existing findings show up in the report as context but never block merges. This lets contributors land incremental improvements on a single area without first cleaning up unrelated legacy issues elsewhere.

---

## How each linter is wired

Five "single-linter diff" sub-steps live in the `lint-baseline-guard` job (`.github/workflows/ci.yml`). All use the same script with different inputs:

```yaml
env:
  LINTER: <golangci-lint analyzer name>
  LABEL: <optional display name; defaults to LINTER>
  TEXT_FILTER: <optional regex on .Text to scope a multi-rule analyzer>
  BASELINE: /tmp/lint-<name>-baseline/report.json
  CURRENT_OUT: /tmp/lint-<name>-current/report.json
run: bash .github/scripts/check-single-linter-diff.sh gitmap-v27
```

Each linter has its own rolling per-SHA cache slot (`golangci-<name>-baseline-main-<sha>`) restored via prefix match.

### Special case: `gosec G115`

`gosec` emits dozens of rules (G101 through G601). To gate only G115 (integer overflow conversions), the sub-step uses `TEXT_FILTER=G115` — a regex applied to each finding's `.Text` field. The filter is applied to **both** current and baseline so the diff stays apples-to-apples (a pre-existing G115 in the baseline will not be flagged as "new" in current).

---

## Full-report diff (separate job)

In addition to the per-linter sub-steps, the `lint-baseline-diff` job runs the **entire** golangci-lint configuration (`.golangci.yml`) and diffs the full JSON report via `.github/scripts/lint-diff.py`. This catches NEW findings from any analyzer enabled in the repo config, even ones without a dedicated sub-step.

The per-linter sub-steps exist on top of the full diff because:

- They guarantee each finding has a `path:line:col` annotation (path-less typecheck errors are dropped), so the GitHub PR-files view always underlines the exact location.
- They isolate each class so a flood of unrelated path-less errors can't drown out a real regression in the PR UI.

---

## When NOT to use baseline-diff

The current state is intentional: every linter goes through baseline-diff. But if a future rule class genuinely cannot ship even once (e.g. a new security check where any occurrence is a CVE), give it its own dedicated **hard-floor** job. Do **not** mix models inside one job — that ambiguity is what motivated the 2026-04-29 cleanup.

A hard-floor job would:
- Run `golangci-lint --no-config --disable-all --enable=<analyzer>` so repo excludes can't mask the finding.
- Fail on **any** count > 0, no baseline, no diff.
- Live in its own GitHub Actions job (separate `name:`, separate required check).

---

## Adding a new baseline-diff linter

Copy an existing block in `lint-baseline-guard` (e.g. the `unused` one) and adjust four things:

1. **Cache slot path** — `/tmp/lint-<name>-baseline` and `/tmp/lint-<name>-current`.
2. **Cache key** — `golangci-<name>-baseline-main-${{ github.sha }}` with matching `restore-keys:` prefix.
3. **`LINTER:` env** — the golangci-lint analyzer name.
4. **(Optional) `TEXT_FILTER:` + `LABEL:`** — when scoping a multi-rule analyzer to a single rule.

That's it. The script handles JSON parsing, normalization, diff, and annotation emission.

---

## Tests

The diff/guard contract is locked in by `.github/scripts/tests/run-tests.sh` (run by the `lint-script-tests` CI job). The 7-case suite covers:

- Empty current report
- Pre-existing findings (unchanged) → tolerated
- New finding → fails with file:line annotation
- Missing baseline → seeding mode (warn, don't fail)
- `TEXT_FILTER` scoping (G115 only)
- `TEXT_FILTER` symmetry (applied to both sides)
- Path-less findings dropped

If you change the script, run the tests locally first:

```sh
bash .github/scripts/tests/run-tests.sh
```

---

## History

| Date | Change |
|---|---|
| pre-2026-04-29 | Job named `lint-regression-guard`. Mixed contract: `unused` + `gosec G115` were hard-floor (zero tolerance, ignored baseline) while `misspell`, `gocritic`, `exhaustive` were baseline-diff. Misleading umbrella name. |
| 2026-04-29 (pass 1) | Renamed → `lint-hard-floor` to make the split semantics honest. Comments updated per-step. |
| 2026-04-29 (pass 2) | Unified all five linters under baseline-diff. Hard-floor script (`check-lint-regressions.sh`) deleted. `check-single-linter-diff.sh` extended with `TEXT_FILTER` + `LABEL` to handle the gosec-G115 scoping. Job renamed → `lint-baseline-guard`. |
| 2026-04-29 (pass 3) | Added `lint-script-tests` CI job + 7-case unit-test suite to lock in the contract. |

Related issue file: `.lovable/cicd-issues/02-lint-regression-guard-semantics.md`.

---

## Files referenced

| Path | Role |
|---|---|
| `.github/workflows/ci.yml` | Defines `lint-baseline-guard`, `lint-baseline-diff`, `lint-script-tests` jobs |
| `.github/scripts/check-single-linter-diff.sh` | Per-linter baseline-diff engine (the only script needed) |
| `.github/scripts/lint-diff.py` | Full-report baseline diff (separate job) |
| `.github/scripts/tests/run-tests.sh` | Unit tests for the diff engine |
