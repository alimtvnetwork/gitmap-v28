# Subtask 04: Live Verification on u1 & Evidence Ledger

## Status: QUEUED
## Owner: Lead Orchestrator
## Dependencies: Subtask 03

---

## Deliverables
1. Run automated remote verification commands on node `u1`:
   - `test -f ~/.config/Cursor/User/settings.json`
   - `test -f ~/.config/Cursor/User/globalStorage/state.vscdb` (check size > 100 MB, SQLite integrity)
   - `test -f ~/.config/Cursor/User/globalStorage/conversation-search.db`
   - `test -d ~/.config/Cursor/User/workspaceStorage` (check workspace count)
   - `test -d ~/.cursor/projects` (check project count)
   - `test -d ~/.cursor/agents` & `test -d ~/.cursor/skills-cursor`
2. Update `00-master-audit-ledger.md` with complete evidence.

---

## Acceptance Criteria
- [ ] 100% of verification commands return true and exit 0.
- [ ] Zero corruption of SQLite databases.
- [ ] Evidence recorded in audit ledger.
