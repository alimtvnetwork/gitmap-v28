# Subtask 03: Remote Dispatch & u1 Atomic Unpack

## Status: QUEUED
## Owner: Worker Subagent 2
## Dependencies: Subtask 02

---

## Deliverables
1. Remote transfer of migration archive and unpack script to node `u1` via GitMap SSH copy (`gitmap ssh copy`).
2. Automatic pre-flight snapshots created on `u1`:
   - `~/.config/Cursor/User.bak.<timestamp>`
   - `~/.cursor.bak.<timestamp>`
3. Extraction of payload with proper permissions (`0755` dirs, `0644` files, `a:a` ownership).

---

## Acceptance Criteria
- [ ] Safe non-destructive deployment.
- [ ] Rollback snapshot verified on `u1`.
- [ ] Remote extraction command exits 0.
