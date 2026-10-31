# Architecture Spec 360: Mandatory AI Name Enforcement, On-The-Fly Teaching & Settings UI Configuration

- **Spec ID:** `360-mandatory-ai-name-and-settings-ui-configuration`
- **Component:** GitMap CLI, Settings Subsystem, Web UI & Skills
- **Status:** Active
- **Created At:** 2026-10-31
- **Author:** Antigravity AI Assistant

---

## 1. Executive Summary

This architecture specification formalizes mandatory AI agent name attribution across GitMap autonomous commands (`gitmap login set ai-auto <agent-name>` and `gitmap cpf/cpb/cpc/cfr/cpr --ai-agent <name>`), on-the-fly teaching diagnostics when an agent name is missing, and end-to-end configurability of default `ai-auto` parameters through both the CLI settings engine (`gitmap settings set ...`) and the Web Settings UI (`ui_settings.json`).

---

## 2. Core Architecture Invariants

1. **Mandatory AI Agent Name:**
   - In autonomous mode, every AI agent MUST provide its identifier (e.g. `muse`, `gemini`, `cursor-model`, `antigravity`, `claude`).
   - If an AI agent attempts to execute `gitmap login set ai-auto` without an agent name (and without `--global` / `global`), GitMap rejects the invocation and emits actionable on-the-fly teaching examples.
   - If a commit is executed while an AI bot identity is active and `--ai-agent` is omitted, GitMap enforces attribution (falling back to configured settings or providing on-the-fly teaching guidance).

2. **On-The-Fly Teaching System:**
   - When required parameters are missing, error messages MUST NOT be opaque stack traces or bare exit codes.
   - Output must clearly demonstrate the correct command syntax with diverse, realistic examples (`muse`, `gemini`, `cursor-model`, `antigravity`) so autonomous agents immediately self-correct without human intervention.

3. **Multi-Tier Settings Resolution:**
   - Default `ai-auto` configuration resolves through a deterministic hierarchy:
     1. Explicit CLI arguments/flags (`--ai-agent <name>`, `gitmap login set ai-auto <name>`)
     2. Web UI settings (`ui_settings.json` -> `attributes["ai_auto.default_*"]`)
     3. Persistent SQLite database settings (`store.DB` / `config.SetVariable("global", ...)`)
     4. Built-in system defaults (`gitmap automation bot - runs autonomously`, `info@rgitmap.com`).

4. **Web Settings UI Integration:**
   - The Web Settings UI (`ui_assets_markup_settings_a.go`) provides responsive form controls for:
     - `ai_auto.default_name` (Default Bot Name)
     - `ai_auto.default_email` (Default Bot Email)
     - `ai_auto.default_agent` (Default AI Agent Name)
   - The backend API (`/api/settings`) exposes and saves these attributes atomically.

5. **Cross-Repository Synchronization:**
   - Canonical skills (`gitmap`, `gitmap-automation-bot`, `execute-parent-task-with-n-steps-v6`, `execute-parent-task-with-n-steps-v7`) and execution prompts are synchronized across `gitmap` and `coding-guidelines`, enforcing mandatory AI name attribution in the prompt contract.
