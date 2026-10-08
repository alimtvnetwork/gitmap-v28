# Architecture Specification: GitMap Universal Runner, Audit DB, Supabase Vault, PE-AI & Muse Skills

> **Specification ID:** 239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills  
> **Status:** Ratified  
> **Version:** 1.0.0  
> **Target Subsystems:** `cli/cmdrun/`, `cli/cmdsupabase/`, `cli/cmdpipeline/`, `cli/cmd/`, `01-prompts/27-muse-prompts/`, `.agents/skills/`  

---

## 1. Executive Summary

This architecture specification defines the comprehensive design and system contracts for five interconnected GitMap subsystems:
1. **Universal Polyglot File Runner (`gitmap run <file>`):** An intelligent execution router that automatically detects file extensions (`.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go`), resolves verified interpreters from the local installation cache, enqueues execution lifecycles into the SQLite task audit database, records structured error telemetry into an isolated run errors database upon failure, and exposes inspection commands (`gitmap run errors`, `gitmap run history`).
2. **Multi-Supabase Database Vault (`gitmap supabase`):** A zero-cleartext secrets management subsystem allowing multi-project Supabase registration via single-line commands (`gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]`). Sensitive keys and database connection strings are protected with authenticated AES-GCM / RSA vault encryption using machine-derived keys before being stored in SQLite.
3. **Chore Semantic Commits (`cpc`) & Expanded Help Menus:** Introduction of `gitmap cpc "<module> - <summary>"` (alias for `commit-push-chore`) adhering to strict hyphen formatting, alongside full-form expansion across all commit and fleet help menus (`cpf`: `commit-push-feature`, `cpb`: `commit-push-bug`, `cpc`: `commit-push-chore`, `cpr`: `commit-push-release`, `pcp`: `pull-commit-push`, `pas`: `pull-all-ssh`).
4. **Pipeline PE-T Diagnostics, `--ai` Flag & Source Attribution:**
   - `gitmap pe -t`: Polling loop checking every 2 minutes; if errors or stack traces are detected in logs, it immediately displays diagnostic telemetry via `gitmap pe` and terminates execution.
   - `gitmap pe -ud` (`--until-done`): Continuous execution monitoring until full CI/CD run terminates, returning final pipeline verdict.
   - `gitmap pe --ai`: Suppression of system clipboard modifications and clipboard terminal notices, keeping AI agent context clean.
   - Enhanced Source Attribution: Enriches error reports with logger name, workflow name, run URL, and sanitized relative log paths.
   - Immediate Git Hash Retrieval: Instant cache resolution for completed commit SHAs without blocking network requests.
5. **Muse Master Prompt & Skills Modernization:** Updating `01-prompts/27-muse-prompts/01-muse-master-prompt.md`, `.agents/skills/muse-master-prompt/`, and `.agents/skills/gitmap/` to mandate GitMap native command primacy (`gitmap run`, `gitmap task`, `gitmap pe --ai`, `gitmap sync`, `gitmap aum search`), updating coding guidelines to require `--ai` on all pipeline error inspections.

---

## 2. Architectural Invariants

1. **Deterministic Extension Resolution Invariant:**
   - When executing `gitmap run <target>`, if `<target>` exists directly on disk, its extension dictates the execution runtime.
   - If `<target>` does not exist as specified, the runner MUST probe candidate extensions in strict deterministic order: `.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go`.
   - If both a saved macro and a file match `<target>`, explicit file targets (containing path separators or matching disk files) take precedence, with fallback to macro execution if no file matches.
2. **Dual-Audit Invariant (Task Audit & Run Errors DB):**
   - Every file execution MUST register a lifecycle entry in the SQLite tasks database (`TaskQueue` status `pending` -> `running` -> `completed` / `failed`).
   - If the process exits with a non-zero code or encounters an unhandled runtime exception, a structured error record MUST be inserted into the `run_errors` table in the SQLite database, capturing the target file, interpreter, exit code, duration, stdout snippet, stderr snippet, and timestamp.
3. **Zero Cleartext Secrets Invariant:**
   - Cleartext database passwords, service role keys, and API tokens MUST NEVER be stored in SQLite, flat files, or emitted in logs.
   - All sensitive Supabase credentials MUST be encrypted via authenticated AES-256-GCM (`cli/crypto/Encrypt`) using a 32-byte master key derived from the host machine identifier and vault salt before persistence.
   - Display commands (`gitmap supabase list`, `gitmap supabase show`) MUST mask keys (`eyJhbG...****`).
4. **AI Clipboard Suppression Invariant:**
   - Under `--ai`, the pipeline error inspection suite MUST NEVER invoke `writeClipboard` or emit the `📋 Copied pipeline logs to system clipboard` banner.
   - All error summaries, stack traces, and job breakdowns MUST be rendered purely through stdout for AI consumption without altering the operating system clipboard buffer.
5. **Strict Relative Git Path Hygiene:**
   - All logged file paths, error locations, and database records MUST strictly use relative Git workspace paths (`cli/...`, `data/...`), stripping OS absolute prefixes (`C:\...`, `/home/...`).

---

## 3. Subsystem 1: Universal Polyglot File Runner (`cli/cmdrun/`)

### 3.1 Architecture Overview

The universal file runner provides a single unified interface for executing scripts across polyglot ecosystems without requiring developers or AI agents to manually remember specific interpreter flags, environment paths, or error handling.

```
                    gitmap run <target> [args...]
                                 │
                     ┌───────────▼───────────┐
                     │  Target File Resolver │
                     │  (Disk check + Probe) │
                     └───────────┬───────────┘
                                 │ Target Resolved
                     ┌───────────▼───────────┐
                     │ Task DB Audit Enqueue │
                     │   (Status: running)   │
                     └───────────┬───────────┘
                                 │
                     ┌───────────▼───────────┐
                     │ Interpreter Execution │
                     │  (.py, .ps1, .sh,     │
                     │   .js, .ts, .go)      │
                     └───────────┬───────────┘
                                 │
                   ┌─────────────┴─────────────┐
                   │                           │
          Exit Code == 0              Exit Code != 0
                   │                           │
         ┌─────────▼─────────┐       ┌─────────▼─────────┐
         │ Task DB Complete  │       │ Task DB Failed    │
         │ (Status: success) │       │ Run Errors Logged │
         └───────────────────┘       └───────────────────┘
```

### 3.2 File Extension Resolution Table

When `<target>` is provided without an extension, candidate extensions are evaluated in deterministic priority:

| Extension | Target Interpreter | Resolution Strategy | Fallback / Environment |
| :--- | :--- | :--- | :--- |
| `.py` | Python 3 | `installation.db` cached python path | `python3`, `python`, `py -3` |
| `.ps1` | PowerShell Core / Desktop | Cross-platform `pwsh` | Windows PowerShell (`powershell -NoProfile`) |
| `.sh` | POSIX Shell / Bash | `bash` | `sh`, Git Bash (`/bin/bash`) |
| `.js` | Node.js / Bun | `node` | `bun`, `deno` |
| `.ts` | TypeScript Runtime | `bun` -> `tsx` -> `ts-node` | `node --loader ts-node/esm` |
| `.go` | Go Compiler Runner | `go run` | Pre-compiled binary runner |

### 3.3 Task Queue & Audit Logging

Before initiating execution, `cli/cmdrun` interacts with the SQLite tasks database (`store.OpenTasksRootSplitDB()`):
- Creates a `TaskQueue` entry with `Action: "file-run"`, `Target: <resolved_path>`, `Status: "running"`.
- Records timing metrics (`StartedAt`).
- Upon clean exit (`0`), updates status to `"completed"` in `TaskHistory`.

### 3.4 Run Errors Database (`run_errors`)

On execution failure (non-zero exit code or execution launch failure), `cli/cmdrun` persists an error record:

```sql
CREATE TABLE IF NOT EXISTS run_errors (
    ErrorId        TEXT PRIMARY KEY,
    FilePath       TEXT NOT NULL,
    FileExtension  TEXT NOT NULL,
    Interpreter    TEXT NOT NULL,
    ExitCode       INTEGER NOT NULL,
    DurationMs     INTEGER NOT NULL,
    ErrorMessage   TEXT NOT NULL,
    StdoutSnippet  TEXT NOT NULL DEFAULT '',
    StderrSnippet  TEXT NOT NULL DEFAULT '',
    ExecutedArgs   TEXT NOT NULL DEFAULT '',
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxRunErrors_FilePath ON run_errors(FilePath);
CREATE INDEX IF NOT EXISTS IdxRunErrors_CreatedAt ON run_errors(CreatedAt DESC);
```

### 3.5 CLI Subcommands for Error Remediation

- `gitmap run <file> [args...]`: Execute target script with auto-extension resolution and audit logging.
- `gitmap run errors` (aliases: `gitmap run-errors`, `gitmap run err`): Displays tabular list of recent execution failures, showing Error ID, File, Exit Code, Duration, and Error Snippet.
- `gitmap run history`: Displays historical execution audits from `TaskHistory`.
- `gitmap run clear-errors`: Purges error history from `run_errors`.

### 3.6 Disambiguation in `macro_root_dispatch.go` & `roottooling.go`

In `cli/cmd/macro_root_dispatch.go`, `runMacroRootRun`:
1. Check if `args[0]` matches known subcommands (`errors`, `history`, `clear-errors`). If yes, route to `cmdrun`.
2. Check if `args[0]` exists on disk as a file or if appending candidate extensions resolves to an existing file. If yes, route to `cmdrun.RunFile()`.
3. If no file matches, attempt to load as dynamic macro via `macro.LoadMacroPolymorphic()`.
4. If neither matches, output clear usage message describing both file runner and macro execution.

---

## 4. Subsystem 2: Multi-Supabase Database Vault & Encryption (`cli/cmdsupabase/`)

### 4.1 Architecture Overview

GitMap enables managing multiple Supabase projects across development, staging, and production environments, supporting fleet coordination without storing plaintext secrets.

```
 gitmap supabase add <alias> <url> <anon> <service> [db_url]
                             │
                ┌────────────▼────────────┐
                │ Machine Master Key Gen  │
                │ (Device ID + Salt KDF)  │
                └────────────┬────────────┘
                             │
                ┌────────────▼────────────┐
                │   AES-256-GCM Encrypt   │
                │  - anon_key -> enc     │
                │  - service_key -> enc  │
                │  - db_url -> enc       │
                └────────────┬────────────┘
                             │
                ┌────────────▼────────────┐
                │ SQLite Split DB Vault   │
                │ (supabase/sql.db)       │
                └─────────────────────────┘
```

### 4.2 Cryptographic Key Derivation & AES-GCM Vault

1. **Master Key Derivation:**
   - Derived using PBKDF2 or SHA-256 HMAC combining a local machine hardware identifier (CPU/motherboard UUID) and an installation-scoped salt stored in `.gitmap/data/installation/vault.salt`.
   - Generates a deterministic 32-byte key for AES-256.
2. **Authenticated Encryption:**
   - Utilizes `cli/crypto.Encrypt(plaintext, key)`: AES-GCM with standard 12-byte random nonces and base64 encoding.
   - Integrity and authenticity are guaranteed; tampered ciphertext returns decryption failure.

### 4.3 Database Schema (`supabase/sql.db`)

Stored under `.gitmap/data/installation/supabase/sql.db`:

```sql
CREATE TABLE IF NOT EXISTS supabase_databases (
    Alias          TEXT PRIMARY KEY,
    ProjectRef     TEXT NOT NULL DEFAULT '',
    ApiUrl         TEXT NOT NULL,
    AnonKeyEnc     TEXT NOT NULL,
    ServiceKeyEnc  TEXT NOT NULL,
    DbUrlEnc       TEXT NOT NULL DEFAULT '',
    Status         TEXT NOT NULL DEFAULT 'active',
    Description    TEXT NOT NULL DEFAULT '',
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxSupabase_Status ON supabase_databases(Status);
```

### 4.4 CLI Command Interface

- `gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]`: Encrypts and saves credentials.
- `gitmap supabase list` (alias: `gitmap supabase ls`): Lists registered projects with masked keys.
- `gitmap supabase test <alias>` (alias: `gitmap supabase ping`): Tests connectivity via REST endpoint (`<url>/rest/v1/`).
- `gitmap supabase get <alias> [--decrypt]`: Retrieves project configuration (only decrypts if explicitly requested).
- `gitmap supabase remove <alias>` (alias: `gitmap supabase rm`): Deletes project registration.
- `gitmap supabase export-env <alias>`: Emits environment variables (`SUPABASE_URL`, `SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_ROLE_KEY`) for process scripting.

---

## 5. Subsystem 3: Chore Commits (`cpc`) & Help Menus

### 5.1 Semantic Commit Taxonomy Expansion

GitMap standardizes commit commands to automate staging, prefixing, and pushing:

| Command | Alias | Full Form | Prefix Added | Target Scope |
| :--- | :--- | :--- | :--- | :--- |
| `gitmap commit-push-feature` | `cpf` | Commit & Push Feature | `Feature: ` | New features, capabilities |
| `gitmap commit-push-bug` | `cpb` | Commit & Push Bugfix | `Bug: ` | Defect fixes, regressions |
| `gitmap commit-push-chore` | `cpc` | Commit & Push Chore | `Chore: ` | Maintenance, refactoring, dependencies, specs |
| `gitmap commit-push-release` | `cpr` | Commit & Push Release | `Release: ` | Version bumps, releases |
| `gitmap pull-commit-push` | `pcp` | Pull Commit Push | Contextual | Staging pull, commit, and push |
| `gitmap pull-all-ssh` | `pas` | Pull All SSH | N/A | Multi-node fleet git pull |

### 5.2 Commit Message Rules

- All commit messages passed to `cpc` MUST use hyphen-separated formatting:
  `gitmap cpc "<module> - <summary>"`
- GitMap automatically prepends `Chore: `, resulting in `Chore: <module> - <summary>`.
- Colons inside `<summary>` are forbidden to prevent duplicate prefix artifacts.

### 5.3 Help Text Modernization

The help menus (`cli/cmd/commit_help_menu.go`, `cli/helptext/`) MUST explicitly present full forms and aliases:
- Document `cpf` as `commit-push-feature`.
- Document `cpb` as `commit-push-bug`.
- Document `cpc` as `commit-push-chore`.
- Document `cpr` as `commit-push-release`.
- Document `pas` as `pull-all-ssh`.

---

## 6. Subsystem 4: Pipeline PE-T Diagnostics, `--ai` Flag & Source Details

### 6.1 Diagnostic Polling with `-t` and `-ud`

In `cli/cmdpipeline/`:
- **`gitmap pe -t` (2-Minute Dynamic Check):**
  - Polls GitHub Actions workflow status at 2-minute intervals (`120s`).
  - Evaluates log extracts and stack traces on each cycle.
  - If errors are detected in the active or failed run, it immediately outputs `gitmap pe` diagnostics and exits with code 1.
  - If no error is present and runs are in progress, continues waiting until the next interval or timeout.
- **`gitmap pe -ud` / `--until-done`:**
  - Runs in a continuous monitoring loop until all pipeline jobs reach a terminal state (`completed`, `failure`, `success`, `cancelled`).
  - Emits the final status and error report upon completion.

### 6.2 The `--ai` Flag & Clipboard Suppression

- **Problem:** Conventional `gitmap pe` execution triggers OS clipboard write operations (`clipboard.WriteAll`), emitting notices like `📋 Copied pipeline error logs to clipboard`. In autonomous AI environments, clipboard modifications cause race conditions, OS security popups, and terminal pollution.
- **Solution:** Passing `--ai` sets `PipelineErrorFlags.IsAI = true`. When active:
  - `writeClipboard()` is bypassed.
  - `printClipboardNotice()` is suppressed.
  - Output is streamed strictly to stdout in cleanly structured Markdown.

### 6.3 Enhanced Source Attribution

Every error section in `gitmap pe` reports exact provenance details:
- **Logger Name / Origin:** Originating job runner or tool (e.g., `golangci-lint`, `pytest`, `cargo-test`).
- **Workflow Name:** GitHub Actions workflow identifier (e.g., `CI / Build and Test`).
- **Run URL:** Direct HTML link to the GitHub run.
- **Relative Log Path:** Sanitized relative path to cached log file (e.g., `data/pipeline/<slug>/<run-id>.log`), removing `%LOCALAPPDATA%` and absolute machine paths.

### 6.4 Immediate Git Hash Retrieval

- When inspecting errors for a specific commit SHA or offset (`gitmap pe <sha>`):
  - If the commit SHA exists in the local SQLite pipeline database and all its runs are marked completed, `cli/cmdpipeline` immediately returns the cached verdict without initiating external HTTP requests.
  - If runs are active or incomplete, cache hit is rejected to query live status.

---

## 7. Subsystem 5: Muse Master Prompt & Skills Modernization

### 7.1 Master Prompt Updates (`01-prompts/27-muse-prompts/01-muse-master-prompt.md`)

- Incorporate GitMap command substitution matrix into the onboarding protocol.
- Instruct Muse agents to use `gitmap run <file>` for test scripts and verification tools.
- Enforce `gitmap pe --ai` across all CI/CD monitoring tasks.
- Include `gitmap cpc "<module> - <summary>"` in commit standards.

### 7.2 Skill Updates (`.agents/skills/`)

- **`muse-master-prompt/skill.md`:** Add rules for `gitmap run`, `gitmap cpc`, `gitmap supabase`, and `gitmap pe --ai`.
- **`gitmap/SKILL.md`:** Expand the command cheat sheet with universal runner, Supabase vault, and chore commit documentation.
- **`02-spec/02-coding-guidelines/06-ai-optimization/`:** Document the mandate that all AI agents MUST use `gitmap pe --ai` to prevent clipboard interference.

---

## 8. Implementation Roadmap & Subtask Allocation

```
┌────────────────────────────────────────────────────────────────────────┐
│ Task 239: GitMap Runner, Audit DB, Supabase Vault, PE-AI & Muse Skills │
└────────────────────────────────────┬───────────────────────────────────┘
                                     │
         ┌───────────────────────────┴───────────────────────────┐
         │                                                       │
┌────────▼────────────────────────┐             ┌────────────────▼────────────────┐
│ Subtask 01: Universal Runner    │             │ Subtask 02: Supabase Vault      │
│ & Task Audit DB (`cmdrun/`)     │             │ & AES-GCM Enc (`cmdsupabase/`)  │
└─────────────────────────────────┘             └─────────────────────────────────┘
         │                                                       │
         ├───────────────────────────┬───────────────────────────┤
         │                                                       │
┌────────▼────────────────────────┐             ┌────────────────▼────────────────┐
│ Subtask 03: Chore Commits       │             │ Subtask 04: PE-T AI & Source    │
│ (`cpc`) & Help Text Full Forms  │             │ Telemetry (`cmdpipeline/`)      │
└─────────────────────────────────┘             └─────────────────────────────────┘
                                     │
                        ┌────────────▼────────────┐
                        │ Subtask 05: Muse Master │
                        │ Prompt & Skills Sync    │
                        └─────────────────────────┘
```

| Subtask ID | Focus Area | Key Components |
| :--- | :--- | :--- |
| **Subtask-01** | Universal Polyglot Runner & Audit DB | `cli/cmdrun/`, `macro_root_dispatch.go`, `roottooling.go`, `run_errors` table |
| **Subtask-02** | Multi-Supabase Database Vault | `cli/cmdsupabase/`, `supabase_databases` table, AES-GCM encryption |
| **Subtask-03** | Chore Commits & Help Expansion | `cli/cmd/commit_push.go`, `commit_help_menu.go`, `cpc` / `commit-push-chore` |
| **Subtask-04** | Pipeline PE-T, `--ai` & Attribution | `cli/cmdpipeline/`, `-t`, `-ud`, `--ai`, source telemetry |
| **Subtask-05** | Muse Prompt & Skills Modernization | `01-prompts/27-muse-prompts/`, `.agents/skills/`, coding guidelines |
