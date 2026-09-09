# Subtask 87.04: Parallel Nested Ifs & Enum Checkers Upgrades

## Objective
Upgrade `linter-scripts/check-nested-ifs.py` and `linter-scripts/check-enum-and-boolean.py` to use chunked parallel execution (8 files/chunk) instead of individual single-file futures, add 5-second snapshot heartbeats, and support `--changed-only`.

## Proposed Changes
1. `linter-scripts/check-nested-ifs.py`:
   - Replace single-future dispatch with 8-file chunking (`chunk_items`).
   - Add 5-second snapshot heartbeat printing: `[Snapshot 5s] Scanned X/Y files (Z%) | Active workers: W | Throughput: N files/sec`.
   - Add `--changed-only` support reading from `.lovable/temp/git-changed-files.json`.
2. `linter-scripts/check-enum-and-boolean.py`:
   - Apply 8-file chunking to reduce thread scheduling overhead.
   - Add 5-second snapshot heartbeat.
   - Add `--changed-only` support reading from `.lovable/temp/git-changed-files.json`.

## Acceptance Criteria
- [ ] Both linters execute successfully across the full repository.
- [ ] Both linters support `--changed-only` and complete in < 1s.
- [ ] 0 violations found.
- [ ] Coding guidelines met.
