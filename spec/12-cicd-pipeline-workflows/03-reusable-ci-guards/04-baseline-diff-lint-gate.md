# Pattern 04 — Baseline-Diff Lint Gate

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

You enable a strict linter on a mature codebase. The first run reports
thousands of pre-existing findings. Two unsatisfactory options:

1. **Fail every PR until all findings are fixed.** Blocks all work.
2. **Disable the new linter entirely.** No quality bar at all.

The baseline-diff pattern offers option 3: **fail only on findings
introduced by the current change**. The full lint report is computed,
then diffed against a cached baseline (last successful run on the
default branch). Pre-existing ("UNCHANGED") findings are tolerated;
new ("NEW") findings fail the build; resolved ("FIXED") findings are
celebrated in the summary.

Over time the baseline shrinks and the codebase converges on zero
findings without ever blocking a merge.

---

## Algorithm

```
INPUT:  current_report.json, baseline_report.json (optional)
OUTPUT: exit 1 if NEW findings exist (and baseline present), else exit 0

1. current  = load_findings(current_report)
2. baseline_present = baseline file exists AND non-empty
3. baseline = load_findings(baseline_report) if present else empty_set

4. added     = current  - baseline   # set difference
5. fixed     = baseline - current
6. unchanged = current  & baseline

7. print human-readable summary (counts + per-file NEW/FIXED list)

8. seeding mode (no baseline yet):
     emit ::warning per finding (so reviewers see it)
     exit 0  -- never gate the very first run

9. gate mode (baseline present):
     if added is non-empty:
        emit ::error file=…,line=… per finding
        exit 1
     else:
        exit 0
```

---

## Finding Normalization

A "finding" is a 4-tuple: `(file, line, linter, message)`. Severity,
source-line snippet, and column are intentionally **excluded** because
they're noisy across linter versions and would cause spurious diffs
when the linter is upgraded.

The set difference uses Python's built-in `set` semantics — exact
tuple equality. Reordering of issues in the JSON, edits to unrelated
files, and severity changes never produce false positives.

---

## Baseline Storage

| Approach | Pros | Cons |
|----------|------|------|
| GitHub Actions cache (keyed by `default-branch + linter-version`) | Auto-evicted, no commits | Cache miss = seeding mode (acceptable) |
| Committed baseline file | Reproducible locally, version-controlled | Adds noise to git history |
| S3 / artifact storage | Cross-workflow shareable | Extra infra |

The reference implementation uses the GitHub Actions cache for
simplicity. The script tolerates cache misses by entering "seeding
mode" (warnings, exit 0).

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `python3 lint-diff.py --current cur.json --baseline base.json` |
| Exit `0` | No new findings, OR baseline missing (seeding) |
| Exit `1` | One or more new findings (and baseline present) |
| Exit `2` | `--current` argument missing or unreadable |
| stdout | Human-readable diff summary + GitHub Actions annotations |

---

## Configuration Surface

| Variable | Purpose | Example |
|----------|---------|---------|
| `--current` | Current run's lint report | `golangci.json`, `eslint.json`, `ruff.json` |
| `--baseline` | Cached baseline (may be missing/empty) | `.cache/lint-baseline.json` |
| `LOAD_FINDINGS()` | Per-linter normalizer to the 4-tuple | One function per JSON shape |

---

## Adaptations

### golangci-lint (Go)

JSON shape: `{"Issues": [{"Pos": {"Filename":..., "Line":...}, "FromLinter":..., "Text":...}]}`

Normalizer:
```python
for issue in data.get("Issues") or []:
    pos = issue.get("Pos") or {}
    yield (pos.get("Filename",""), int(pos.get("Line",0) or 0),
           issue.get("FromLinter",""), (issue.get("Text") or "").strip())
```

### ESLint (Node / TypeScript)

JSON shape (with `--format json`): `[{"filePath":..., "messages":[{"line":..., "ruleId":..., "message":...}]}]`

Normalizer:
```python
for file_report in data:
    file = file_report.get("filePath","")
    for msg in file_report.get("messages") or []:
        yield (file, int(msg.get("line",0) or 0),
               msg.get("ruleId","") or "eslint",
               (msg.get("message") or "").strip())
```

### Ruff (Python)

JSON shape (with `--output-format json`): `[{"filename":..., "location":{"row":...}, "code":..., "message":...}]`

Normalizer:
```python
for issue in data:
    yield (issue.get("filename",""),
           int((issue.get("location") or {}).get("row", 0) or 0),
           issue.get("code",""),
           (issue.get("message") or "").strip())
```

### Clippy (Rust)

JSON shape (cargo `--message-format=json`): line-delimited JSON, filter
`reason == "compiler-message"` and `message.level in {"warning",
"error"}`.

Normalizer extracts `message.spans[0].file_name`, `line_start`,
`message.code.code`, `message.message`.

---

## Output Format

```
========================================================================
  GOLANGCI-LINT DIFF vs LAST SUCCESSFUL MAIN
========================================================================

  + NEW       :    2
  - FIXED     :    5
  = UNCHANGED :  118

  NEW findings:
    pkg/scan/runner.go
      + L42 [errcheck] Error return value of `os.Setenv` is not checked
      + L91 [gocritic] paramTypeCombine: parameters can be combined

  FIXED findings:
    pkg/release/main.go
      - L17 [unused] func `unusedHelper` is unused
    ...
========================================================================
```

Per-finding annotations follow:
```
::error file=pkg/scan/runner.go,line=42::[errcheck] Error return value of `os.Setenv` is not checked (NEW)
```

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| Linter version upgrade changes message text → all findings appear "new" | Pin the linter version in CI; bump deliberately and regenerate the baseline in the same PR |
| Cache key includes branch SHA → never hits | Key on default-branch name + linter-version, NOT current SHA |
| Diff fails when JSON is empty or malformed | Treat parse error as "no findings"; emit `::warning::` and continue |
| Findings on generated code count as new | Exclude generated files via the linter config, not in the diff script |

---

## Reference Implementation

See `.github/scripts/lint-diff.py` in `gitmap-v6`.

---

## Cross-References

- [00-overview.md](./00-overview.md)
- [05-actionable-lint-suggestions.md](./05-actionable-lint-suggestions.md) — Companion script that turns NEW findings into PR comments
- [02-grandfather-baseline-naming.md](./02-grandfather-baseline-naming.md) — Same "grandfather" idea for identifier names

---

*Baseline-diff lint gate — v1.0.0 — 2026-04-21*