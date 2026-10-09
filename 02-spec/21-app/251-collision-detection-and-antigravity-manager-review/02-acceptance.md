# Spec 251 — Acceptance criteria

## Benchmark & cache
- **B01**: `agent collisions` first line matches `^\[\d+ms\]`. Evidence: CLI transcript.
- **B02**: `agent stats`, `agent heatmap`, `agent ps` same first-line benchmark.
  `--json` includes `elapsed_ms` first. Evidence: transcripts.
- **B03**: SQLite-as-cache documented (this spec §A1); `collisions` opens exactly
  one DB (central). Evidence: spec + code.

## agent ls / show
- **L01**: `gitmap agent ls` lists ≤15 tasks with ID, slug, status, agents,
  subtasks rollup, updated. Evidence: transcript.
- **L02**: subtask subtree shown with IDs + statuses; broken DBs show BROKEN
  (no crash). Evidence: transcript (incl. a broken-DB fixture).
- **L03**: `gitmap agent show <slug>` prints full detail + subtask subtree.
  Evidence: transcript.
- **L04**: `ls` completes <300ms for 15 tasks (benchmark line proves it).

## Antigravity
- **G01**: Attribution answered (Jeik = upstream or fork?) with evidence.
- **G02**: PR #6 reviewed; ideas extracted; considered commits recorded in memory.
- **G03**: Scope tier map documented; no deletions; guardrails added.
- **G04**: `unwrap()`/`expect()` reduced in proxy handlers (count before/after).

## Release
- **R01**: gitmap 6.519.0 tagged+pushed; `gitmap version` = 6.519.0.
- **R02**: honest report on `gitmap update` (sandbox-blocked → rebuilt from source).
