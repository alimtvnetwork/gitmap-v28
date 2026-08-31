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
  │  Phase 3: Local Verification & Linting                      │
  │  ➔ python .github/scripts/go-format-check.py, go test       │
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

## 5. Alternative Commands & Aliases Cheat Sheet

| Action / Goal | Standard CLI Command | Fast AI Shortcut / Alias | Notes |
| :--- | :--- | :--- | :--- |
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
| **Pull & Push** | `gitmap pull-commit-push` | `gitmap pcp "<msg>"` | Pull latest, commit, and push |
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
```
