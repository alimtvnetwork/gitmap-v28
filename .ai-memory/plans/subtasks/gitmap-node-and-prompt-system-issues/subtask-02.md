# Subtask 02: Interactive Prompt Tree Suite & Telegram Prompt Sorting

**Parent Plan:** `.ai-memory/plans/gitmap-node-and-prompt-system-issues.md`  
**Subtask ID:** `subtask-02`  
**Worker:** Worker 02  
**Status:** In Progress  

## Scope of Ownership
- `cli/cmdagy/agy_ui.go`
- `cli/cmdagy/agy_ui_html.go`
- `cli/cmdagy/agy_telegram_subcommands.go`
- `cli/cmdui/ui_assets_js_prompts_help.go`
- `cli/cmdui/ui_assets_markup_ops.go`

## Implementation Checklist
- [ ] In `cli/cmdagy/agy_ui.go`, implement backend handlers for:
  - `/api/prompts/backup`: captures active/queued prompts into backup batch
  - `/api/prompts/restore`: restores batch into active queue
  - `/api/prompts/reset`: resets queue and prompt tree state
  - `/api/prompts/refresh`: performs immediate status re-scan
- [ ] In `cli/cmdagy/agy_ui_html.go`, replace naive `alert(msg)` tree view with interactive visual tree modal/panel including Send, Backup, Restore, Reset, and Refresh buttons.
- [ ] In `cli/cmdagy/agy_telegram_subcommands.go`, support machine targeting (`/status <machine>`, `/rp <machine>`, `/rwi <machine> ...`) and add Telegram send verification.
- [ ] In `cli/cmdui/ui_assets_js_prompts_help.go` and `cli/cmdagy/agy_ui_html.go`, sort prompt queues with:
  - `running` prompts strictly on TOP.
  - `queued` / `enqueued` prompts immediately below running prompts.
  - Completed / idle prompts at the bottom.
- [ ] Ensure explicit badges: `RUNNING` (yellow/green) and `QUEUED` / `ENQUEUED` (cyan).
- [ ] Verify with targeted test: `go test ./cmdagy -run TestAgyUI`.
