# Subtask 02: Settings CLI & Web Settings UI Configuration

- **Parent Plan:** `360-mandatory-ai-name-and-settings-ui-configuration`
- **Slug:** `mandatory-ai-name-and-settings-ui-configuration - Task 02`
- **Status:** `completed`
- **Assigned Subagent:** Worker 02

---

## Deliverables
1. Update `cli/cmdui/ui_assets_markup_settings_a.go`:
   - Add responsive form group under `settings-pane-speed` (Speed & Automation):
     - `ai_auto.default_name` (Default AI Auto Bot Name)
     - `ai_auto.default_email` (Default AI Auto Email)
     - `ai_auto.default_agent` (Default AI Agent Name)
2. Update `cli/cmdui/ui_server_settings.go`:
   - Expose default values in `loadSettings()`.
   - Ensure `mergeSettings()` and `handleAPISettings()` persist `ai_auto.*` attributes to `ui_settings.json` and sync with `store.DB` / `config`.
3. Update `cli/cmdagy/agy_telegram_email_settings.go`:
   - Include `AIAutoDefaultName`, `AIAutoDefaultEmail`, `AIAutoDefaultAgent` in `SpeedSettingsStatus`.
   - Support `gitmap settings set ai_auto.default_name <val>`, `ai_auto.default_email <val>`, `ai_auto.default_agent <val>`.
   - Display AI Auto section in `runSettingsShow`.
