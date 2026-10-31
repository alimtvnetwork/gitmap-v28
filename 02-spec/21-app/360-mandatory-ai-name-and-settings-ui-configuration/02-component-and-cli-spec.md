# Component & CLI Spec 360: Mandatory AI Name Enforcement, On-The-Fly Teaching & Settings UI Configuration

- **Spec ID:** `360-mandatory-ai-name-and-settings-ui-configuration`
- **Component:** `cli/cmdlogin`, `cli/cmdcommit`, `cli/cmdui`, `cli/cmdagy`, `cli/config`
- **Status:** Active
- **Created At:** 2026-10-31

---

## 1. CLI Behavioral Specification

### 1.1 `gitmap login set ai-auto`

- **Syntax:**
  ```bash
  gitmap login set ai-auto <agent-name> [--global]
  gitmap login set global ai-auto [<agent-name>]
  ```
- **Validation Rule:**
  - If `args` contains `set ai-auto` without an `<agent-name>` and without `--global`/`global`:
    - Check if `ai_auto.default_agent` is configured in Settings. If configured, apply that default.
    - If NOT configured, return an actionable error with on-the-fly teaching:
      ```
      Error: AI agent name is mandatory for 'gitmap login set ai-auto <agent-name>'.

      💡 Examples to teach on the fly:
        • gitmap login set ai-auto muse
        • gitmap login set ai-auto gemini
        • gitmap login set ai-auto cursor-model
        • gitmap login set ai-auto antigravity
        • gitmap login set ai-auto claude

      To configure a generic machine-wide bot without a specific agent name, use:
        • gitmap login set global ai-auto
      ```

### 1.2 `gitmap commit` Semantic Verbs (`cpf`, `cpb`, `cpc`, `cfr`, `cpr`)

- **Syntax:**
  ```bash
  gitmap cpf "<message>" --ai-agent <name>
  ```
- **Behavior:**
  - When committing with an AI bot identity active, `--ai-agent <name>` provides cognitive heatmap attribution.
  - If `--ai-agent` is omitted:
    - Check configured `ai_auto.default_agent` in Settings.
    - If set, use it for telemetry and author naming.
    - If empty, display a teaching guidance message demonstrating `--ai-agent <name>`.

### 1.3 `gitmap settings` CLI & Web UI Settings

- **CLI Keys:**
  - `ai_auto.default_name`: Bot name template or default string.
  - `ai_auto.default_email`: Bot email address.
  - `ai_auto.default_agent`: Fallback agent name when not specified.
- **Commands:**
  ```bash
  gitmap settings set ai_auto.default_name "gitmap automation bot - runs autonomously"
  gitmap settings set ai_auto.default_email "info@rgitmap.com"
  gitmap settings set ai_auto.default_agent "muse"
  gitmap settings
  ```
- **Web UI:**
  - Located under `Speed & Automation` pane in Web Settings UI (`/api/settings`).
  - Inputs with `data-key="ai_auto.default_name"`, `data-key="ai_auto.default_email"`, `data-key="ai_auto.default_agent"`.
