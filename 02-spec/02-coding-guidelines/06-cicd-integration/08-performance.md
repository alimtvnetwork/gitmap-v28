# Performance — Probe Order, Parallelism, Timeouts

> **Version:** 1.0.0
> **Updated:** 2026-04-19
> **Status:** Active

The linter pack must finish in **< 30 seconds** on a typical 50 kLOC
repository running on a free GitHub Actions runner (2 vCPU, 7 GB RAM).
This doc captures the three techniques that make that budget achievable
and mirrors the same tricks already documented for the `install.sh`
self-update probe.

---

## 1. Middle-out probe ordering

When the orchestrator scans a repo, it does **not** walk files in
alphabetical order. Alphabetical walks bias toward `apps/`, `assets/`,
and `cmd/` — directories that are usually thin — and leave the heavy
`internal/`, `pkg/`, `src/`, and `vendor-adjacent` trees for the end,
delaying the first finding until late in the run.

**Middle-out** sorts candidate roots by **estimated byte weight** and
probes from the median outward, alternating heavier-then-lighter:

```
sorted by size:   [A, B, C, D, E, F, G]   (A smallest, G largest)
median index:     3 → D
probe order:      D, E, C, F, B, G, A
```

This surfaces violations from the densest code first, so:

- Early-exit modes (`--fail-fast`) get the worst offender within the
  first ~20% of wall-clock time.
- CI logs show meaningful findings before the runner times out on a
  pathological repo.
- Cache warmup in tree-sitter parsers benefits the largest files first,
  amortizing parser-init cost across the longest-running checks.

The same algorithm is used in `install.sh` to probe release-asset
mirrors median-first instead of in declaration order.

### Implementation hook

```python

# linters-cicd/checks/_lib/walker.py

def walk_files_middle_out(root, extensions):
    files = walk_files(root, extensions)
    files.sort(key=lambda p: p.stat().st_size)
    return _middle_out(files)
```

Plugins opt in by calling `walk_files_middle_out` instead of
`walk_files`. Default remains alphabetical for deterministic test
fixtures.

---

## 2. Parallel check execution

`run-all.sh` currently iterates checks **sequentially** for log
readability. Performance mode runs them in parallel with a worker pool
sized to `nproc - 1` (minimum 1).

| Mode | Trigger | Behavior |
|------|---------|----------|
| Sequential | default | One check at a time, ordered log output |
| Parallel | `--jobs N` or env `LINTERS_JOBS=N` | N workers, results merged at end |
| Auto | `--jobs auto` | `N = max(1, nproc - 1)` |

Parallelism is at the **(rule, language)** level — each registry entry
runs in its own process. SARIF output stays per-tool (one `runs[]`
entry per worker), so the merge step in `scripts/merge-sarif.py` is
unchanged.

### Why per-tuple, not per-file

Per-file parallelism wastes Python interpreter startup (~60 ms × file
count). Per-tuple parallelism amortizes startup across the full file
walk and lets each tree-sitter parser stay warm.

---

## 3. Per-check timeout budgets

A runaway check (regex catastrophic backtracking, infinite AST loop)
must not block the entire pipeline. Every plugin invocation is wrapped
in a hard timeout.

| Scope | Default budget | Override |
|-------|----------------|----------|
| Single check (rule × language) | **20 s** | `--check-timeout SECONDS` |
| Whole orchestrator run | **120 s** | `--total-timeout SECONDS` |
| Per-file parse inside a check | **2 s** | hard-coded in plugin lib |

On timeout the orchestrator emits a SARIF result with
`level: "error"` and `ruleId: "TOOL-TIMEOUT"`, attaches the offending
file path when known, and continues with the remaining checks. Exit
code becomes `2` (tool error) so CI surfaces it distinctly from real
findings (`1`).

### Why these numbers

- **20 s per check** — measured worst case on the 50 kLOC corpus is
  ~7 s for `magic-strings` on TS. 20 s gives 3× headroom.
- **120 s total** — fits inside GitHub Actions' default 6-hour job
  budget with massive safety margin while still failing fast on
  pathological inputs.
- **2 s per file** — anything slower indicates a parser pathology, not
  a real file. Skip and emit `TOOL-TIMEOUT` for that file only.

---

## Acceptance criteria

1. `time bash linters-cicd/run-all.sh --path . --jobs auto` on this
   repo completes in < 30 s wall-clock.
2. Synthetic test repo with 1 × 10 MB minified JS file does **not**
   hang the orchestrator — `TOOL-TIMEOUT` is emitted within 25 s.
3. Middle-out probe order is unit-tested with a fixed file-size
   manifest in `linters-cicd/checks/_lib/walker_test.py`.
4. Parallel mode produces byte-identical merged SARIF (after sorting
   `results[]` by `ruleId, uri, startLine`) compared to sequential
   mode.

---

## Cross-References

- [Orchestrator](./01-index.md#layer-1--portable-check-scripts-linters-cicdchecks)
- [Plugin Model](./03-plugin-model.md)
- [Acceptance Criteria](./97-acceptance-criteria.md)
- `install.sh` middle-out mirror probe — same pattern, different domain

---

*Part of [CI/CD Integration](./01-index.md)*
