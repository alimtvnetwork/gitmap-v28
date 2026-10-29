# Master Plan: 352 — Cursor-Like IDE Setup from GitMap

> **Status:** PENDING  
> **Sequence ID:** 352  
> **Slug:** `352-cursor-like-ide-setup-from-gitmap`  
> **Canonical Spec:** `02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/readme.md`  
> **Subtask Root:** `.ai-memory/plans/subtasks/352-cursor-like-ide-setup-from-gitmap/`  
> **Category:** Developer Ergonomics & AI IDE Orchestration  
> **Target Version:** `v6.630.0`  
> **Created:** 2026-10-29  

---

## 1. Executive Summary & Problem Statement

Developers frequently choose Cursor for three distinct AI coding superpowers:
1. **Cursor Composer**: Agentic multi-file codebase refactoring with filesystem diff views and terminal execution.
2. **Cursor Inline & Tab**: Instant sub-second inline edits (`Cmd/Ctrl+K` or `Cmd/Ctrl+I`) and context-aware ghost-text tab autocomplete.
3. **High-Speed Native Editor**: Fast keystroke response with integrated AI chat and prompt transforms.

However, proprietary AI editors create severe enterprise liabilities: vendor lock-in, recurring per-seat subscription fatigue, lack of transparency regarding code privacy, and inability to seamlessly use Bring-Your-Own-Key (BYOK) high-capability models such as **Kimi K3.1 Pro (Moonshot AI)**, **DeepSeek**, or local **Ollama** models.

### The GitMap Solution
GitMap automates the configuration of open-source and native editors directly from the GitMap CLI (`gitmap ide setup`) and GitMap Fleet Web UI (`gitmap ui` -> Settings / IDE Setup):
- **VS Code + Cline** for the Cursor Composer multi-file experience, injected with GitMap MCP server tools.
- **VS Code + Continue.dev** for Cursor Tab autocomplete and inline code transformations.
- **Zed AI Editor** for ultra-fast native editing with direct BYOK model integration.
- Built-in multi-provider presets pre-configured with context windows and endpoints: **Kimi K3.1 Pro** *(Recommended default)*, **OpenRouter**, **DeepSeek**, **Anthropic Claude**, **Ollama**, and **Custom**.

---

## 2. User Request & Non-Negotiable Directives

```markdown
# Replicate Cursor-like IDE Setup using GitMap and Kimi K3.1 Pro

Provide a 1-click / 1-command Cursor-like IDE setup directly from GitMap CLI and Fleet Web UI. Support all three primary open alternatives:
1. VS Code + Cline (Cursor Composer multi-file editing replica)
2. VS Code + Continue.dev (Cursor Tab autocomplete and inline edit replica)
3. Zed AI Editor (Fast Rust-native BYOK AI editor)
Pre-configure Kimi K3.1 Pro (Moonshot AI) as the recommended primary provider preset, alongside OpenRouter, DeepSeek, Anthropic, Ollama, and Custom endpoints. Support cross-platform configuration for Windows, macOS, and Linux with safe JSON merging, pre-flight backups, CLI flags, REST API endpoints, and interactive Web UI controls.

## task: Cursor-Like IDE Setup from GitMap
## slug: 352-cursor-like-ide-setup-from-gitmap

# Actionable Items Must Follow Non-Negotiable:
1. Git pull BEFORE starting any work (first gitmap pull, fallback git pull).
2. Author specs under 02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/ and enqueue plan in .ai-memory/plans/pending/352-cursor-like-ide-setup-from-gitmap.md (subtasks in .ai-memory/plans/subtasks/352-cursor-like-ide-setup-from-gitmap/).
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String.
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cli/...); never add absolute paths.
5. Provide multi-provider presets with Kimi K3.1 Pro as the primary recommended default.
6. Support cross-platform OS path resolution (Windows, macOS, Linux) without hardcoded separators.
7. Implement non-destructive safe JSON merging with mandatory pre-flight backup (<file>.bak-<timestamp>).
8. Implement CLI command gitmap ide setup [target] with comprehensive flags (--provider, --model, --base-url, --api-key, --dry-run, --json, --test, --force, --backup).
9. Implement REST API endpoints in Fleet Web UI (/api/ide/setup/status, /api/ide/setup/presets, /api/ide/setup/test-connection, /api/ide/setup/apply) and interactive Web UI tabbed visualizer.
10. Adhere to coding guidelines: <= 15 line functions, zero nested ifs, positive booleans (is*, has* only), leaf error sovereignty, and typed enums ending in Type.
11. Execute terminal E2E verification, test suites to 100% green pass, and atomic release ceremony via gitmap cpr.
```

---

## 3. 5-Pillar Architectural Mandate

### Pillar 1: Strategic Intent & Systems Reasoning
Unify disparate AI coding configurations into a repeatable, automated GitMap command and interactive Web UI workflow. Enable developers to run a single command (`gitmap ide setup all --provider kimi`) to instantly equip their developer environment with Cursor-grade capabilities while retaining full ownership of their API keys, network routes, and source code privacy.

### Pillar 2: Explicit Scope Boundaries
- **In-Scope:**
  - `cli/cmdide/setup/`: Core setup engine, enums, presets, path resolver, connection tester, config generators, safe merger, and pre-flight backup manager.
  - `cli/cmdide/setup_cmd.go`, `cli/cmdide/setup_flags.go`, `cli/cmdide/setup_view.go`: CLI command dispatch and terminal UI rendering.
  - `cli/cmdide/enums/`: Strongly typed enums (`IDETargetType`, `IDEProviderType`, `SetupModeType`).
  - `cli/cmdui/ui_server_ide_setup.go`: REST API server endpoints mounted in `cli/cmdui/ui_server.go`.
  - `cli/cmdui/ui_assets_js_ide_setup.go` & `ui_assets_markup_ide_setup.go`: Web UI frontend tabbed visualizer and settings integration.
  - Comprehensive unit and E2E test suites with mock server validation.
  - `.agents/skills/gitmap/skill.md`: Documentation and command recipes.
- **Out-of-Scope:**
  - Writing a custom editor binary from scratch (GitMap configures existing VS Code extensions and Zed).
  - Storing user API keys in Git or remote repositories (keys stay strictly in local OS configuration files or environment variables).
  - Removing or clobbering existing user settings, keybindings, or third-party extensions.

### Pillar 3: Modular Decomposition & Subtasks Breakdown
The plan is partitioned into five discrete subtasks executed in sequential dependency order:
- **Subtask 01:** Core IDE Setup Engine, Target Enums, and Provider Presets (`01-ide-setup-engine-and-presets.md`)
- **Subtask 02:** Specialized Configuration Generators (Cline, Continue.dev, Zed AI) & Safe Merger (`02-config-generators-cline-continue-zed.md`)
- **Subtask 03:** GitMap CLI `ide setup` Command, Subcommand Dispatch, and Flag Parser (`03-cli-command-and-flags.md`)
- **Subtask 04:** Web UI Server REST Endpoints, API Route Mounting, and Settings Visualizer (`04-ui-server-api-and-web-interface.md`)
- **Subtask 05:** Live End-to-End Verification, Cross-Platform Testing, and Test Suite Parity (`05-live-testing-and-verification.md`)

### Pillar 4: Risk Mitigation & Invariant Enforcement
- **Risk 1: Accidental JSON Clobbering**: Merging JSON files can overwrite unrelated user settings.  
  *Mitigation*: Deep unmarshaling into map trees with targeted leaf updates, atomic temporary file writing, and mandatory timestamped pre-flight backup files (`<file>.bak-<timestamp>`).
- **Risk 2: Invalid API Credentials or Endpoints**: Users could misconfigure API keys leading to silent failures in the editor.  
  *Mitigation*: Built-in connection testing engine (`--test` flag and `/api/ide/setup/test-connection` endpoint) that issues lightweight test queries and validates HTTP status codes and response structures before writing.
- **Risk 3: Cross-Platform File Paths**: Hardcoded backslashes break on macOS/Linux; hardcoded Unix paths break on Windows.  
  *Mitigation*: Dynamic `filepath.Join`, environment variable resolution (`APPDATA`, `USERPROFILE`, `HOME`), and OS-specific lookup tables.
- **Risk 4: Guideline Violations**: Regressions in function length, nested conditions, or exit codes.  
  *Mitigation*: Strict adherence to <= 15 lines per function, positive booleans only, zero nested ifs, leaf error sovereignty returning `*apperror.AppError`, and strongly typed exit codes.

### Pillar 5: Verification & End-to-End Release Gates
- All unit test suites pass at 100% green with > 85% branch coverage.
- Live CLI testing across `--dry-run`, `--json`, and real/mock provider endpoints.
- Web UI REST API verification with `httptest.NewServer`.
- Final atomic minor bump release ceremony via `gitmap cpr` with release notes update.

---

## 4. Detailed Subtask Ledger & Execution Sequence

| Step | Subtask Document | Key Deliverables | Estimated Files | Status |
|:---:|:---|:---|:---:|:---:|
| **01** | `01-ide-setup-engine-and-presets.md` | `IDETargetType`, `IDEProviderType` enums, Preset catalog, Cross-platform Path Resolver, Connection Tester, Option structs | 6 files | PENDING |
| **02** | `02-config-generators-cline-continue-zed.md` | Cline generator with MCP injection, Continue.dev generator with model arrays, Zed generator with language_models, Safe JSON Merger, Backup manager | 6 files | PENDING |
| **03** | `03-cli-command-and-flags.md` | `gitmap ide setup` command dispatch, flag parsing (`--provider`, `--model`, `--base-url`, `--api-key`, `--dry-run`, `--json`, `--test`), Lipgloss terminal UI, Help text | 5 files | PENDING |
| **04** | `04-ui-server-api-and-web-interface.md` | REST API endpoints (`/status`, `/presets`, `/test-connection`, `/apply`), Route mounting, Tabbed Web UI frontend (Cline, Continue, Zed, All), Diff modal | 5 files | PENDING |
| **05** | `05-live-testing-and-verification.md` | End-to-end integration tests, cross-platform path tests, live endpoint verification, skill updates in `gitmap/skill.md`, atomic release ceremony | 4 files | PENDING |

---

## 5. Architectural Contracts & Schema Specifications

### 5.1 Strongly Typed Enums (`cli/cmdide/enums/`)

```go
package enums

// IDETargetType represents a target editor for AI configuration.
type IDETargetType string

const (
	TargetCline    IDETargetType = "cline"
	TargetContinue IDETargetType = "continue"
	TargetZed      IDETargetType = "zed"
	TargetAllAI    IDETargetType = "all"
)

// IDEProviderType represents a supported LLM provider preset.
type IDEProviderType string

const (
	ProviderKimi       IDEProviderType = "kimi"
	ProviderOpenRouter IDEProviderType = "openrouter"
	ProviderDeepSeek   IDEProviderType = "deepseek"
	ProviderAnthropic  IDEProviderType = "anthropic"
	ProviderOllama     IDEProviderType = "ollama"
	ProviderCustom     IDEProviderType = "custom"
)
```

### 5.2 Provider Presets Catalog (`cli/cmdide/setup/presets.go`)

| Provider Key | Default Base URL | Default Model | Fallback Model | Context Window | Auth Scheme |
|:---|:---|:---|:---|:---:|:---|
| `kimi` | `https://api.moonshot.cn/v1` | `kimi-k3.1-pro` | `moonshot-v1-128k` | 128,000 | `Bearer <KEY>` |
| `openrouter` | `https://openrouter.ai/api/v1` | `moonshotai/kimi-k3.1-pro` | `anthropic/claude-3.7-sonnet` | 128,000 | `Bearer <KEY>` |
| `deepseek` | `https://api.deepseek.com/v1` | `deepseek-chat` | `deepseek-reasoner` | 64,000 | `Bearer <KEY>` |
| `anthropic` | `https://api.anthropic.com/v1` | `claude-3-7-sonnet-20250219` | `claude-3-5-sonnet-latest` | 200,000 | `x-api-key: <KEY>` |
| `ollama` | `http://localhost:11434/v1` | `deepseek-r1` | `qwen2.5-coder:32b` | 32,000 | None / Optional |
| `custom` | `--base-url <url>` | `--model <id>` | Configurable | Configurable | `Bearer <KEY>` |

### 5.3 Target Path Resolution by OS

```go
// PathResolver resolves configuration targets across operating systems.
type PathResolver interface {
	GetClineSettingsPath() (string, error)
	GetClineStorageDir() (string, error)
	GetContinueConfigPath() (string, error)
	GetZedSettingsPath() (string, error)
	IsTargetInstalled(target IDETargetType) bool
}
```

- **Windows**:
  - Cline MCP Settings: `%APPDATA%\Code\User\globalStorage\saoudrizwan.claude-dev\settings\cline_mcp_settings.json`
  - Continue Config: `%USERPROFILE%\.continue\config.json`
  - Zed Settings: `%APPDATA%\Zed\settings.json`
- **macOS**:
  - Cline MCP Settings: `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`
  - Continue Config: `~/.continue/config.json`
  - Zed Settings: `~/.config/zed/settings.json`
- **Linux**:
  - Cline MCP Settings: `~/.config/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`
  - Continue Config: `~/.continue/config.json`
  - Zed Settings: `~/.config/zed/settings.json`

### 5.4 REST API Server Route Contract (`cli/cmdui/`)

| Method | Endpoint | Request Body | Response Envelope | Purpose |
|:---:|:---|:---|:---|:---|
| `GET` | `/api/ide/setup/status` | None | `Result[IDESetupStatusSummary]` | Detect installed editors, existing config presence, and active providers. |
| `GET` | `/api/ide/setup/presets` | None | `Result[[]ProviderPresetConfig]` | List all built-in LLM presets, default models, and context limits. |
| `POST` | `/api/ide/setup/test-connection` | `TestConnectionRequest` | `Result[TestConnectionResult]` | Perform lightweight test ping to validate provider credentials and measure latency. |
| `POST` | `/api/ide/setup/apply` | `ApplySetupRequest` | `Result[ApplySetupResult]` | Atomically backup and merge AI configurations into selected targets with diff preview. |

---

## 6. Coding Standards & Review Checklist

Before any pull request or commit, verify the following standards:

- [ ] **Function Length Ceiling**: Every Go function is <= 15 lines of code (preferred <= 8 lines).
- [ ] **Zero Nested If Statements**: All conditional blocks are flattened using guard clauses and early returns.
- [ ] **Strict Positive Booleans**: All boolean variables and struct fields start with `is` or `has` (`isInstalled`, `hasConfig`, `isDryRun`, `isQuiet`, `isJson`, `isForceCreate`, `isConnectionSuccessful`, `isBackupEnabled`). Zero negative names.
- [ ] **Leaf Error Sovereignty**: Functions returning `error` or `*apperror.AppError` return immediately without calling `os.Exit` or CLI exit handlers.
- [ ] **Database Audit Integration**: CLI command errors bubble up to top-level command dispatchers so they are recorded by `finishCommandAudit`.
- [ ] **No Magic Strings or Integers**: Exit codes use strongly typed enums (`constants.ExitUsageError`, `constants.ExitGeneralError`); providers and targets use typed enums ending in `Type`.
- [ ] **Relative Path Hygiene**: All documentation, imports, test cases, and commit messages refer strictly to relative Git paths.
- [ ] **Semantic Commit Verbs**: All commits use GitMap semantic verbs (`gitmap cpf`, `gitmap cfr`, `gitmap cbf`, `gitmap cpr`) formatted as `"MODULE - summary"`.

---

## 7. Definition of Done & Release Protocol

1. Subtasks 01 through 05 completed and verified.
2. Unit and integration tests pass at 100% green (`go test -v ./cli/cmdide/setup/... ./cli/cmdui/...`).
3. CLI help verification: `gitmap ide setup --help` and `gitmap ide --help` accurately describe new capabilities.
4. GitMap Fleet Web UI verified: Settings -> IDE Setup renders tabs, provider selection, connection tester, and configuration apply flow.
5. Skill synchronization: `.agents/skills/gitmap/skill.md` updated with alternative command recipes.
6. Atomic release ceremony: Staged minor bump via `gitmap cpr "v6.630.0 cursor-like ide setup with kimi k3.1 pro and open editors"`.
