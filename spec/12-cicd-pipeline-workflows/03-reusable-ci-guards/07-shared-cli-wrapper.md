# Pattern 07 — Shared CLI Wrapper

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

The six guards in this folder (Patterns 01–06) each ship as a
stand-alone Bash or Python script with its own argument style,
environment-variable conventions, and exit-code semantics. Wiring
them into a workflow means remembering six different invocations:

```yaml
- run: bash check-cmd-naming.sh gitmap/cmd
- run: CONST_DIR=gitmap/constants bash check-naming.sh
- run: python3 check-collisions.py
- run: python3 lint-diff.py --current cur.json --baseline base.json
- run: python3 lint-suggest.py --current cur.json --baseline base.json --out $GITHUB_STEP_SUMMARY
- run: bash test-summary.sh ./test-artifacts
```

Six different surfaces means six things to learn, six things to
document, six places where flag names drift over time. The shared
CLI wrapper unifies them behind one entry point with one flag
vocabulary so contributors and AI assistants can reason about CI
at the **phase** level (`check`, `lint`, `test`) instead of at the
**script** level.

---

## Algorithm

```
INPUT:  --phase <check|lint|test|all> [--guard <name>] [options]
OUTPUT: aggregated exit code; optional JSON summary

1. parse flags; validate --phase is one of the four known values
2. if --guard is set: dispatch that single guard, exit with its code
3. else: enumerate guards bound to the phase
     for each guard in declared order:
         dispatch guard with normalized env / args
         capture exit code
         aggregate -> max(current, incoming)
4. emit JSON summary if --json given
5. exit aggregated code
```

The aggregator uses `max()` so that a tool error (`2`) on guard A
followed by a clean run (`0`) on guard B still surfaces as `2` —
highest severity wins.

---

## Phase → Guard Map

| Phase | Guards (declared order) | Patterns |
|-------|-------------------------|----------|
| `check` | `forbidden-names`, `naming-baseline`, `collisions` | 01, 02, 03 |
| `lint` | `lint-diff`, `lint-suggest` | 04, 05 |
| `test` | `test-summary` | 06 |
| `all` | all of the above | 01–06 |

`--guard <name>` overrides the phase grouping and runs exactly one
guard. Useful in local development (`--guard collisions`) and when
composing custom workflows.

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `bash scripts/ci-runner.sh --phase <name> [options]` |
| Exit `0` | Every selected guard returned `0` |
| Exit `1` | At least one guard reported a violation; none were tool errors |
| Exit `2` | At least one guard hit a tool error (missing script, bad input) |
| Exit `64` | Usage error (unknown flag, missing `--phase`, unknown guard) |
| stderr | `::error::[ci-runner] …` annotations + `=== guard: name ===` section banners |
| stdout (JSON) | Written to `--json <path>` if provided; shape below |

### JSON Summary Shape

```json
{
  "phase": "check",
  "overall": 1,
  "guards": [
    {"guard": "forbidden-names", "exit": 0},
    {"guard": "naming-baseline", "exit": 1},
    {"guard": "collisions",      "exit": 0}
  ]
}
```

Single-guard mode emits the simpler form:

```json
{"guard": "collisions", "exit": 0}
```

---

## Configuration Surface

| Flag | Purpose | Example |
|------|---------|---------|
| `--phase` | Required. Phase to run. | `check`, `lint`, `test`, `all` |
| `--guard` | Run a single guard by name. | `collisions` |
| `--source-dir` | Source directory for check-phase guards. | `cmd/`, `src/lib/` |
| `--baseline` | Baseline file path. | `.github/scripts/constants-baseline.txt` |
| `--results-dir` | Matrix-test results directory. | `./test-artifacts` |
| `--json` | Emit JSON summary to path. | `/tmp/ci.json` |
| `--fix` | Regenerate baseline (naming-baseline only). | — |
| `--scripts-dir` | Where guard scripts live. | `.github/scripts` (default) |
| `--verbose` | Diagnostic messages to stderr. | — |
| `--help` | Print usage and exit `0`. | — |

The wrapper does **not** invent new flags for guards; it normalizes
the way they're passed. Each dispatcher function translates the
wrapper flags into the guard's native argument style:

- `forbidden-names` receives `--source-dir` as positional `$1`
- `naming-baseline` receives `--source-dir` and `--baseline` as
  `CONST_DIR` and `BASELINE_FILE` env vars
- `lint-diff` and `lint-suggest` receive `--baseline` as
  `--baseline <path>` arg
- `test-summary` receives `--results-dir` as positional `$1`

---

## Adaptations Per Language

The wrapper is Bash-only by intent — every CI runner ships Bash.
Adapting it to a different host repo is purely about pointing at
the right scripts:

### Go repo

```yaml
- run: bash scripts/ci-runner.sh --phase all --source-dir gitmap/constants --baseline .github/scripts/constants-baseline.txt --results-dir ./test-artifacts --json $RUNNER_TEMP/ci.json
```

### Node / TypeScript repo

Substitute the underlying scripts with TS-aware variants but keep the
same wrapper:

```yaml
- run: bash scripts/ci-runner.sh --phase check --source-dir src/lib/constants --baseline .baselines/constants.txt
```

### Python repo

```yaml
- run: bash scripts/ci-runner.sh --phase lint --baseline .baselines/ruff.json
```

### Rust repo

```yaml
- run: bash scripts/ci-runner.sh --phase all --source-dir src --baseline .baselines/clippy.json --results-dir ./cargo-test-output
```

---

## Local Developer UX

```bash

# Run all check-phase guards

./scripts/ci-runner.sh --phase check --source-dir src/constants

# Iterate on one guard

./scripts/ci-runner.sh --guard collisions --verbose

# Regenerate the naming baseline after an approved rename

./scripts/ci-runner.sh --guard naming-baseline --fix

# Full pre-push gate, machine-readable summary

./scripts/ci-runner.sh --phase all --json /tmp/ci.json
echo "exit=$?"; cat /tmp/ci.json | jq .
```

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| `set -e` aborts on first guard failure → later guards never run | Wrapper uses `set -uo pipefail` (no `-e`); each dispatcher returns its code, aggregator collects all |
| Missing guard script silently skipped | Dispatchers verify file exists; absence returns `EXIT_TOOL_ERROR` |
| Single-guard exit `0` masks a tool error elsewhere | Use `--phase all` in CI; reserve `--guard` for local debugging |
| JSON summary unwriteable (read-only path) | Wrapper logs the failure but does not change exit code — JSON is best-effort |
| Bash 3 (macOS default) lacks `mapfile`/`readarray` | Wrapper avoids both; uses `while read` from process substitution |

---

## Reference Implementation

See `scripts/ci-runner.sh` in this repo for the working wrapper that
dispatches to the six guards under `.github/scripts/`.

---

## Cross-References

- [00-overview.md](./00-overview.md) — Pattern inventory
- [01-forbidden-name-guard.md](./01-forbidden-name-guard.md) through [06-matrix-test-aggregator.md](./06-matrix-test-aggregator.md) — The six dispatched guards
- [99-ai-implementation-guide.md](./99-ai-implementation-guide.md) — Workflow skeleton wiring the wrapper

---

*Shared CLI wrapper — v1.0.0 — 2026-04-21*