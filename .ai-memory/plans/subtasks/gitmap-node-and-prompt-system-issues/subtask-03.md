# Subtask 03: Supabase 5m Error Backoff, UI Throttling, Header Details & CLM Cache

**Parent Plan:** `.ai-memory/plans/gitmap-node-and-prompt-system-issues.md`  
**Subtask ID:** `subtask-03`  
**Worker:** Worker 03  
**Status:** In Progress  

## Scope of Ownership
- `cli/cmdsupabase/supabase_cmd.go`
- `cli/cmdsupabase/supabase_db.go`
- `cli/cmdagy/agy_account_switch_engine.go`
- `cli/cmdagy/agy_conv_scanner.go`
- `cli/cmdprompttemplate/template_import_export.go`
- `cli/cmdprompt/prompt_export_import.go`
- `cli/cmdui/ui_assets_markup_settings_a.go`
- `cli/cmdui/ui_assets_js_settings.go`

## Implementation Checklist
- [ ] In `cli/cmdsupabase/supabase_cmd.go` and `cli/cmdsupabase/supabase_db.go`, implement a 5-minute error suppression gate: record `lastSupabaseErrorAt` on any connection/auth error, and skip retry checks if `time.Since(lastSupabaseErrorAt) < 5*time.Minute`.
- [ ] In `cli/cmdagy/agy_account_switch_engine.go`, enforce minimum 30-second interval for account credit syncing (`lastCreditSyncAt`).
- [ ] In `cli/cmdagy/agy_ui_html.go` and `cli/cmdui/ui_assets_js_settings.go`, adjust rapid auto-refresh to 5 minutes (`300000ms`) to minimize background transactions.
- [ ] In prompt preview / header rendering, collapse repeated metadata (`default-copy-814`, PIDs, timestamps, system info, user prompt) into an expandable `<details>` button.
- [ ] In `cli/cmdagy/agy_conv_scanner.go`, add in-memory caching to debounce CLM retrievals and eliminate console flickering.
- [ ] In `cli/cmdprompttemplate/template_import_export.go` and `cli/cmdprompt/prompt_export_import.go`, rebuild template search index immediately after import/export.
- [ ] Verify with targeted test: `go test ./cmdsupabase` and `go test ./cmdprompttemplate`.
