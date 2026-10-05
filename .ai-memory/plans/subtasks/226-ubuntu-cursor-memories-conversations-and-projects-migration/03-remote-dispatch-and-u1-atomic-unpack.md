# Subtask 03: Remote Dispatch & u1 Atomic Unpack

## Status: DONE
## Owner: Worker Subagent 2 & Lead
## Dependencies: Subtask 02

---

## Deliverables
1. Remote transfer of migration archive and unpack script to node `u1` via GitMap SSH copy (`gitmap ssh copy`).
2. Automatic pre-flight snapshots created on `u1`:
   - `/home/a/.config/Cursor/User.bak.20261005_140914`
   - `/home/a/.cursor.bak.20261005_140914`
3. Extraction of payload with proper permissions (`0755` dirs, `0644` files, `a:a` ownership).

---

## Acceptance Criteria
- [x] Safe non-destructive deployment.
- [x] Rollback snapshot verified on `u1`.
- [x] Remote extraction command exits 0.
