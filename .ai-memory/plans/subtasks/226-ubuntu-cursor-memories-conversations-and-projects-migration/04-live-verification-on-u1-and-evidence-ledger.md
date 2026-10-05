# Subtask 04: Live Verification on u1 & Evidence Ledger

## Status: DONE
## Owner: Lead Orchestrator
## Dependencies: Subtask 03

---

## Deliverables
1. Run automated remote verification commands on node `u1`:
   - `test -f ~/.config/Cursor/User/settings.json` (Dracula Theme, JetBrains Mono, LF line endings verified).
   - `test -f ~/.config/Cursor/User/globalStorage/state.vscdb` (167 MB, `PRAGMA integrity_check` = ok, 17,667 KV rows, 82 composer sessions).
   - `test -f ~/.config/Cursor/User/globalStorage/conversation-search.db` (4.5 MB, `PRAGMA integrity_check` = ok, 21 conversations).
   - `test -d ~/.config/Cursor/User/workspaceStorage` (51 workspace folders, 48 `workspace.json` pointing to `/home/a/git-work/`).
   - `test -d ~/.cursor/projects` (47 `home-a-git-work-*` folders).
   - `test -d ~/.cursor/skills-cursor` (29 skills).
2. Update `00-master-audit-ledger.md` with complete evidence.

---

## Acceptance Criteria
- [x] 100% of verification commands return true and exit 0.
- [x] Zero corruption of SQLite databases.
- [x] Evidence recorded in audit ledger.
