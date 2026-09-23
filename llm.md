# Gitmap LLM Specification & AI Agent Guidelines

> Public Instruction URL: https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md

Gitmap is a high-performance developer CLI and autonomous automation engine designed specifically for AI coding agents, LLMs, and engineers to explore, search, refactor, execute macros, and self-heal CI/CD pipelines with maximum reliability.

---

## 1. AI Agent Standard Operating Procedure (Order of Commands)

When an autonomous AI Agent is assigned a coding, refactoring, or debugging task, it MUST follow this structured 5-phase lifecycle in exact sequential order:

```text
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 1: Discovery & Context Gathering                     │
  │  ➔ gitmap find-files, find-files-any, search, list-files   │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 2: Code Modification & Refactoring                   │
  │  ➔ gitmap replace, replace-regex, targeted code edits       │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 3: Local Verification, Linting & Autofix             │
  │  ➔ gitmap ai list, gitmap ai run, gitmap ai fix, linting   │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 4: Semantic Commit & Push                            │
  │  ➔ gitmap commit-push-feature (cpf), commit-push-bug (cpb)  │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 5: Non-Blocking CI Telemetry & Self-Healing Loop     │
  │  ➔ gitmap pipeline-ai status --json, gitmap error-logs      │
  └─────────────────────────────────────────────────────────────┘
```

---

## 2. Interactive Shell Macros & Workflow Automation

Gitmap provides an interactive shell macro recording and replay engine with dynamic directory tracking, environment variable expansion, and JSON/YAML structured reporting with line-by-line step logs.

### Macro Recording:

- **Start Interactive Recording**:
  ```bash
  gitmap macro record "deploy-workflow"
  ```
- **In-Session Commands**:
  - `stop` / `exit` / `quit`: Save macro and finish.
  - `cancel` / `abort`: Abort recording without saving.
  - `undo` / `undo-steps <N> [-y]`: Undo the last recorded step or last N steps.
  - `redo` / `redo-steps <N>`: Restore previously undone step(s).
  - `list` / `steps`: View currently recorded steps.
  - `help` / `?`: Show in-session help.
- **Dynamic Path Expansion & Directory Tracking**:
  - Automatically expands `%TEMP%`, `%USERPROFILE%`, `$HOME`, `~`, and environment variables.
  - Dynamically updates working directory when running `cd <dir>`, `cd ..`, `cd -`, or `gitmap cd <repo>`.

### Macro Replay & Structured Report Exporting:

- **Standard Replay**:
  ```bash
  gitmap macro run deploy-workflow
  ```
- **Structured JSON Output (with step logs array)**:
  ```bash
  gitmap macro run deploy-workflow --json
  ```
- **Structured YAML Output**:
  ```bash
  gitmap macro run deploy-workflow --yaml
  ```
- **Export Report to File (JSON or YAML with Absolute Path Confirmation)**:
  ```bash
  # Accepts --file, --filepath, --out, --output, -o, -f
  gitmap macro run deploy-workflow --json --file "reports/deploy.json"
  gitmap macro run deploy-workflow --yaml --out "reports/deploy.yaml"
  ```
  *Structured Payload Structure:*
  ```json
  {
    "macro": "deploy-workflow",
    "status": "success",
    "totalSteps": 2,
    "executedSteps": 2,
    "failedSteps": 0,
    "elapsedSeconds": 0.45,
    "startedAt": "2026-09-01T01:30:00Z",
    "completedAt": "2026-09-01T01:30:00.45Z",
    "outputFile": "D:\\work\\gitmap\\reports\\deploy.json",
    "steps": [
      {
        "stepNum": 1,
        "commandLine": "git status",
        "workingDir": "D:\\work\\gitmap",
        "status": "success",
        "exitCode": 0,
        "elapsedSeconds": 0.15,
        "logs": [
          "On branch main",
          "Your branch is up to date with 'origin/main'."
        ]
      }
    ]
  }
  ```

- **Inspect & List Macros in JSON/YAML**:
  ```bash
  gitmap macro list --json
  gitmap macro list --yaml
  gitmap macro show deploy-workflow --json
  gitmap macro show deploy-workflow --yaml --file "macro-spec.yaml"
  ```

---

## 3. Autonomous CI/CD Telemetry & Self-Healing Loop

AI Agents must monitor remote CI/CD workflows non-blockingly without spamming GitHub API endpoints.

### 1. Check Live Pipeline Status with Auto-Delay:

```bash
gitmap pipeline-ai status --json
```
*Output:*
```json
{
  "isRunning": true,
  "activeWorkflow": "Release",
  "etaSeconds": 180,
  "sleepSeconds": 20,
  "nextAiCommand": "gitmap pipeline-ai status -t 180",
  "pendingPipelines": 1,
  "lastStatus": "in_progress",
  "lastRunUrl": "https://github.com/alimtvnetwork/gitmap-v28/actions/runs/33413393831"
}
```

### 2. Auto-Delay Protection & Dynamic Waiting:

- **Default Delay**: `pipeline-ai status` automatically pauses for 20 seconds before checking to prevent rate-limit bans.
- **Dynamic Wait (`-t <seconds>`)**: When `etaSeconds` is returned, the AI executes `gitmap pipeline-ai status -t <etaSeconds>`, ensuring the agent sleeps until the pipeline is likely finished.

### 3. Extract Failing CI Error Logs to File for RCA:

```bash
gitmap pipeline error-logs --json --tempfile "ci-failure.json"
```
*Behavior:* Extracts the exact failing step name, exit code, and failure logs into `ci-failure.json` for 4-part Root Cause Analysis without polluting context memory.

---

## 4. High-Performance Code Search & File Inspection Engine

AI Agents should use Gitmap's dedicated search tools rather than scanning huge directories or using heavy shell find loops.

### Finding Files by Name:

- **Exact Match**: `gitmap find-files "constants.go" -ext "go"` (alias: `gitmap ff "constants.go"`)
- **Substring Match**: `gitmap find-files-any "record" -ext "go"` (alias: `gitmap ffa "record"`)
- **Prefix Match**: `gitmap find-files-startswith "01-" -ext "md"` (alias: `gitmap ffs "01-"`)
- **Suffix Match**: `gitmap find-files-endswith "_test.go"` (alias: `gitmap ffe "_test.go"`)
- **Wildcard / Glob**: `gitmap find "*macro*.go"` (alias: `gitmap f "*macro*.go"`)
- **List Files in Subtree**: `gitmap list-files "macro/*"`

### Content Search & Targeted Reading:

- **Indexed Keyword Search**: `gitmap search "DirTracker"`
- **Structured JSON Multi-Repo Regex**: `gitmap repo-search-json "ProcessCd"` (alias: `gitmap rsj "ProcessCd"`)
- **Find and Read Sections**:
  ```bash
  gitmap find-read "record_dir.go"
  gitmap find-regex-read "func parseDirectoryChange" -ext "go"
  ```
- **Repo Navigation**: `gitmap cd <repo_name>`

---

## 5. Native AI Scripts Catalog & Autofix Engine (`gitmap ai` / `gitmap scripts`)

Gitmap embeds native discovery, live streaming execution, and repository autofix suites for the 44 autonomous Python scripts in `03-ai-scripts/`:

### Catalog & Script Discovery:

- **List All Scripts**: `gitmap ai list` (alias: `gitmap ai ls`, `gitmap scripts ls`)
- **Filter by Category**: `gitmap ai list --category guidelines` (or `-c audit`, `-c ci`, `-c fix`, `-c spec`)
- **Search by Keyword**: `gitmap ai list --search "format"` (or `-s "naming"`)
- **Show Fix Mode Only**: `gitmap ai list --fix` (alias: `-f`)
- **Structured JSON Export**: `gitmap ai list --json`
- **Verbose Metadata**: `gitmap ai list --verbose` (alias: `-v`)

### Direct & Numeric Execution:

- **Run by Number**: `gitmap ai run 01` (executes `01-ai-instruction-writer.py`)
- **Run by Full Name**: `gitmap ai run 06-cicd-local-runner.py`
- **Direct Dispatch Shortcut**: `gitmap ai 01` or `gitmap scripts 14`
- **Pass Arguments**: `gitmap ai run 06 --dry-run`

### Standardized Repository Autofix Suites:

- **Fix All Targets**: `gitmap ai fix all` (runs guidelines, newlines, paths, naming, encoding, and spelling)
- **Fix Coding Guidelines**: `gitmap ai fix guidelines`
- **Fix Trailing Newlines**: `gitmap ai fix newlines`
- **Fix Documentation Paths**: `gitmap ai fix paths`
- **Fix Naming Conventions**: `gitmap ai fix naming`
- **Fix File Encodings**: `gitmap ai fix encoding`
- **Fix Common Misspellings**: `gitmap ai fix spelling`

---

## 6. Alternative Commands & Aliases Cheat Sheet

| Action / Goal | Standard CLI Command | Fast AI Shortcut / Alias | Notes |
| :--- | :--- | :--- | :--- |
| **AI Scripts Catalog** | `gitmap ai list` | `gitmap ai ls`, `gitmap scripts ls` | Catalog registered AI scripts with category filter |
| **Run AI Script** | `gitmap ai run <token>` | `gitmap ai <token>` | Live streaming execution by number, slug, or alias |
| **Repository Autofix** | `gitmap ai fix [target]` | `gitmap ai fix all` | Standardized autofix targets (guidelines, paths, etc.) |
| **Pipeline Status (AI)** | `gitmap pipeline-ai status --json` | `gitmap pl-ai status --json` | Auto-delays 20s and suggests `nextAiCommand` |
| **Dynamic Wait Time** | `gitmap pipeline-ai status -t <sec>`| `gitmap pl-ai status -t <sec>` | Sleeps `<sec>` before querying |
| **Get Bare ETA** | `gitmap pipeline waittime` | `gitmap eta` | Returns integer seconds |
| **Extract Error Logs** | `gitmap pipeline error-logs` | `gitmap error-logs` | Dumps failing step logs to workspace file |
| **Full Workflow Logs** | `gitmap pipeline logs` | `gitmap logs` | Complete workflow execution trace |
| **Record Macro** | `gitmap macro record <name>` | `gitmap m rec <name>` | Interactive live recording with undo/redo |
| **Run Macro (JSON)** | `gitmap macro run <name> --json` | `gitmap m run <name> --json` | Captures step `logs: []` line-by-line |
| **Run Macro (YAML)** | `gitmap macro run <name> --yaml` | `gitmap m run <name> -y` | Structured YAML execution report |
| **Export Macro File** | `gitmap macro run <name> --file <p>` | `gitmap m run <name> -o <p>` | Saves report to file & prints path |
| **Find File (Exact)** | `gitmap find-files <name>` | `gitmap ff <name>` | Match exact filename with optional `-ext` |
| **Find File (Any)** | `gitmap find-files-any <str>` | `gitmap ffa <str>` | Match substring in filename |
| **Find File (Prefix)** | `gitmap find-files-startswith` | `gitmap ffs <prefix>` | Match prefix in filename |
| **Find File (Suffix)** | `gitmap find-files-endswith` | `gitmap ffe <suffix>` | Match suffix in filename (e.g. `_test.go`) |
| **Commit Feature** | `gitmap commit-push-feature` | `gitmap cpf "<msg>"` | Stage, commit, and push feature branch |
| **Commit Bug Fix** | `gitmap commit-push-bug` | `gitmap cpb "<msg>"` | Stage, commit, and push bugfix branch |
| **Commit Release** | `gitmap commit-push-release` | `gitmap cpr "<msg>"` | Stage, commit, and push release chore |
| **AI Scripts** | `gitmap ai [run/list/fix]` | `gitmap scripts` | Native AI scripts runner, catalog, and fix suite |
| **AI Scaffolder** | `gitmap ai create <name>` | `gitmap scripts new <name>` | Scaffold new AI linter/fixer/auditor script skeleton |
| **Automation Engine** | `gitmap automation <cmd>` | `gitmap auto`, `gitmap py-auto` | Native Go search (lazy regex), polyglot newlines, in-memory cache, and Go vs Py benchmarks |
| **Antigravity** | `gitmap antigravity` | `gitmap agy`, `gitmap ag` | AI agent workspaces and config sync |
| **VS Code PM** | `gitmap vscode` | `gitmap vsc` | VS Code Project Manager integrations |
| **Scheduler** | `gitmap schedule` | `gitmap sc` | Background cron and interval scheduler |
| **Repo Status** | `gitmap status` | `gitmap st` | Dirty, ahead, behind across all repos |

---

## Examples

```bash

# Workflow 1: Record and Replay Macro with JSON Output and Logs Array

gitmap macro record test-build

# (In session: go test ./... -> stop)

gitmap macro run test-build --json --file "reports/test-build.json"

# Workflow 2: Non-blocking CI Status Check and Auto-Delay

gitmap pipeline-ai status --json

# Workflow 3: Fast File Discovery

gitmap find-files "record_dir.go" -ext "go"
gitmap find-regex-read "func ProcessCd" -ext "go"

# Workflow 4: AI Scripts Discovery, Execution, Scaffolding and Standard Autofix

gitmap ai list --category guidelines
gitmap ai create custom-linter --type linter --parallel
gitmap ai new schema-validator --type auditor --dry-run
gitmap ai run 01
gitmap ai fix all

# Workflow 5: Native Go Automation, Lazy Regex Search, Polyglot Newlines, and Benchmarks

gitmap automation search "func Run" cli --ext .go
gitmap automation newlines --fix
gitmap automation cache warm cli
gitmap automation cache status
gitmap automation benchmark all
```

---

## 7. Lazy Regex Pattern Matching & Test Diagnostics

When authoring or refactoring Go code and unit tests that match text patterns or parse outputs, autonomous AI agents MUST follow the `cli/lazyregex` pattern matching architecture:

- **Total Ban on Raw `regexp.MustCompile`**: Use `lazyregex.New(pattern)` for thread-safe lazy compilation and deduplication.
- **Total Ban on Blind Test Assertions**: Never write `if !re.MatchString(content) { t.Error("expected match") }`. Opaque checks conceal failure context from AI agents and reviewers.
- **Mandatory `MatchResult` Wrapped Envelopes**:
  ```go
  var rsLazyRegex = lazyRegex.MatchResult(comparing)
  if rsLazyRegex.IsFailed() {
      t.Error(rsLazyRegex.AppError())
  }
  ```
- **Fluent MatchGroup Accessors**:
  - `rs.Items()`: Returns `[]string` of all captured submatches (`[0]` = full match, `[1..N]` = groups).
  - `rs.Map()`: Returns `GroupMap` containing named capture groups (`(?P<name>...)`).
  - `rs.First()`: Returns full match or empty string.
  - `rs.Last()`: Returns last captured submatch.
  - `rs.FirstOrDefault("default")`: Returns first match or fallback default if empty.
- **Diagnostic Output**: On failure, `rs.AppError()` reports:
  ```text
  [E1000:VALIDATION] validation: regex pattern does not match content
    Pattern:   "^(?P<action>add|join|nodes)$"
    Comparing: "gitmap cluster unknown-command"
    Length:    29 bytes (at=lazyregex/match_error.go:32)
  ```
