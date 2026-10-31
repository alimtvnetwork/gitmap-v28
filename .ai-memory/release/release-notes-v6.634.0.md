# GitMap v6.634.0

## What's Changed in v6.634.0

- **Mandatory AI Agent Name Enforcement:** `gitmap login set ai-auto <agent-name>` now strictly mandates an explicit AI agent name (e.g. `muse`, `gemini`, `cursor-model`, `antigravity`, `claude`) for autonomous AI operations.
- **On-The-Fly Teaching System:** If an AI agent attempts to run `gitmap login set ai-auto` without an agent name, GitMap emits an actionable validation error demonstrating exact syntax and practical examples so the agent self-corrects immediately.
- **Configurable Defaults in Settings & Settings UI:** Added full support for default AI auto parameters (`ai_auto.default_name`, `ai_auto.default_email`, `ai_auto.default_agent`) accessible via CLI `gitmap settings set <key> <val>` and the Web Settings UI (`Speed & Automation` pane).
- **Multi-Tier Settings Hierarchy:** Ambient commit resolution and identity setup seamlessly query Web UI settings (`ui_settings.json`), persistent database configuration, and environment overrides before falling back to system defaults.
- **Synchronized Canonical Skills & Prompts:** Updated `gitmap`, `gitmap-automation-bot`, and V6/V7 execution prompt specifications across both `gitmap` and `coding-guidelines` repositories to enforce mandatory AI name attribution.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.634.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.634.0/install.sh | sh
```
