# Subtask 01: Migration Engine Enhancement & Transition Guide

## Status: DONE
## Owner: Lead Orchestrator & Worker 1
## Dependencies: None

---

## Deliverables
1. `repo-secrets/04-ubuntu-migration/step-by-step-transition-guide.md`:
   - Detailed conceptual explanation of every stage: source discovery, path conversion rationale, SQLite preservation, tarball compression, remote transport, backup snapshotting, extraction, and post-flight verification.
2. `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py`:
   - Python 3 automation script supporting memories, conversations, and projects migration.
3. `repo-secrets/05-scripts/migrate-cursor-memories-conversations.sh`:
   - Shell executable wrapper.

---

## Acceptance Criteria
- [x] Transition guide clearly explains the step-by-step transition, thinking, and path transformation logic.
- [x] Script handles both `~/.cursor` assets and `~/.config/Cursor/User` assets (including `globalStorage` and `workspaceStorage`).
- [x] Script includes dry-run simulation mode.
