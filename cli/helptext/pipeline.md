# pipeline

Query live CI/CD pipeline runs, calculate remaining wait times (ETA), fetch step failure logs, export telemetry, clear local databases, and feed errors directly into Google Antigravity.

## Synopsis

```bash
gitmap pipeline <subcommand> [flags]
gitmap pl <subcommand> [flags]
```

## Aliases & Shorthands

- `gitmap pipeline`
- `gitmap pl`
- `gitmap pipe`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `status` | Live CI/CD execution state, active workflow, ETA, and pending PRs |
| `waittime, eta` | Remaining estimated wait time in seconds (machine-friendly integer) |
| `errors, err` | Aggregated failure logs for current or past commit by offset (`-1`, `-2`, `-3`) |
| `errors clear` | Clear error reports, logs, and database for current repository (or target repo) |
| `clear, clean` | Clear logs, reports, and database for current repository (or target repo) |
| `fix errors agy` | Feed errors, git log & RCA prompt to Antigravity, inject prompt & queue (alias: `aef`) |
| `clear-db, db clear` | Reset/purge local pipeline split SQLite database and telemetry error logs (`-y`) |
| `history, hist` | Visual pipeline run tree and commit status summary for recent commits (`-n 5`) |
| `logs` | Full step logs for target commit, offset, or latest workflow run |
| `pipeline-ai status` | Auto-delay (-t <sec>), stream live errors, and switch to fix |
| `pipeline-ai errors` | Extract failing workflow errors with AI remediation guidance |
| `help` | Show this pipeline command suite documentation |

---

## Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--force` | `-f` | boolean | `false` | Bypass duplicate check and resend errors previously dispatched |
| `--all` | | boolean | `false` | Scan all projects for failures (or retain verbose logs) |
| `--projects` | | integer | `3` | Number of failing projects to batch-fix with Antigravity |
| `--limit` | | integer | `3` | Limit number of failing projects per batch run |
| `--reset-batch` | | boolean | `false` | Reset multi-project batch cursor to the first project |
| `--yes` | `-y` | boolean | `false` | Skip confirmation prompt when clearing pipeline database |
| `--no-inject` | | boolean | `false` | Skip direct Antigravity CLI/IDE injection |
| `--json` | `-j` | boolean | `false` | Output structured JSON payload for scripting and AI agents |
| `--file` | | string | `""` | Write error logs or telemetry to the specified file path |
| `--tempfile` | | string | `""` | Write error logs to `.ai-memory/temp/<filename>` |
| `--split` | | boolean | `false` | Separate failure logs across multiple jobs into distinct files |
| `--clip` | | boolean | `false` | Copy logs directly to system clipboard |
| `--count` | `-n` | integer | `5` | Number of recent commits to display in pipeline history |
| `--help` | `-h` | boolean | `false` | Show help for pipeline commands |

---

## Top-Level Shortcuts

| Shortcut | Equivalent Command |
|----------|-------------------|
| `gitmap aef` | `gitmap pipeline fix errors agy` |
| `gitmap pipeline-fix` | `gitmap pipeline fix errors agy` |
| `gitmap error-logs` | `gitmap pipeline error-logs` |
| `gitmap errors` | `gitmap pipeline errors` |
| `gitmap logs` | `gitmap pipeline logs` |
| `gitmap history` | `gitmap pipeline history` |
| `gitmap waittime` | `gitmap pipeline waittime` |
| `gitmap eta` | `gitmap pipeline eta` |
| `gitmap pipeline clear-db` | `gitmap pipeline db clear` |

---

## Examples

```bash
# Check live CI/CD pipeline status and ETA
gitmap pipeline status

# Machine-readable ETA in seconds
gitmap pipeline waittime

# View failure logs for the latest commit or specific offset
gitmap pipeline errors
gitmap pipeline errors -1

# Clear pipeline logs and database for current repository
gitmap pipeline clear -y

# Clear pipeline data for a specific repository
gitmap pipeline clear alimtvnetwork/gitmap-v28 -y

# Batch fix failing pipelines with Google Antigravity
gitmap pipeline errors agy fix
gitmap aef

# View pipeline history tree
gitmap pipeline history -n 5
```

---

## Accurate Failure Detection Architecture

GitMap pipeline telemetry prevents false-positive "clean status" masking:

- **Branch and Commit Scoped Telemetry**: Pipeline error evaluation is strictly scoped to the active commit SHA and branch name. Workflow runs are never deduplicated globally by workflow title alone.
- **Failures on Feature Branches**: Even if a parent or default branch (`main`) is passing, failure runs on the current commit (e.g. `cat-my-v12` commit `384f5d2`) will accurately be detected, evaluated, and presented with full step logs.

---

## 1. Clear Pipeline Repository Data (`clear`, `errors clear`, `clear-db`)

Purge recorded pipeline run history, cached step outputs, error reports, and the repository SQLite database.

### Usage

```bash
# Inside repository without arguments (auto-detects current repo)
gitmap pipeline clear [-y]
gitmap pipeline errors clear [-y]

# From anywhere specifying target repository
gitmap pipeline clear <repo> [-y]
gitmap pipeline errors clear <repo> [-y]

# Equivalent legacy subcommands
gitmap pipeline clear-db [-y]
gitmap pipeline db clear [-y]
```

### Examples

```bash
# Interactively clear pipeline cache & DB for current repository
gitmap pipeline clear

# Non-interactively clear pipeline cache & DB without confirmation prompt
gitmap pipeline clear -y

# Clear pipeline errors and logs for current repository
gitmap pipeline errors clear

# Clear pipeline data for a specific remote repository
gitmap pipeline clear alimtvnetwork/gitmap-v28 -y
gitmap pipeline errors clear alimtvnetwork/gitmap-v28 -y
```

---

## 2. Inspect Live Status & Telemetry

### Check Live Status in Human-Readable Format

```bash
$ gitmap pipeline status
  Repository: alimtvnetwork/gitmap-v28
  Workflow:   GoReleaser
  Status:     RUNNING (in_progress)
  ETA:        84s
  Last Tag:   v6.154.0
```

### Machine-Readable Status for AI Agents

```bash
$ gitmap pipeline status --json
{
  "isRunning": true,
  "etaSeconds": 84,
  "lastTagRelease": "v6.154.0",
  "pendingPipelines": 4,
  "pendingTasks": 0,
  "pendingPRs": 0,
  "repo": "alimtvnetwork/gitmap-v28",
  "activeWorkflow": "GoReleaser",
  "lastStatus": "in_progress",
  "lastRunUrl": "https://github.com/alimtvnetwork/gitmap-v28/actions/runs/33329109649",
  "updatedAt": "2026-08-30T18:48:19Z"
}
```

---

## 3. Negative Commit Offsets (`-1`, `-2`, `-3`)

Inspect failure logs from earlier commits when diagnosing regressions across recent pushes:

```bash
# Aggregate all workflow errors for the previous commit (-1)
gitmap pipeline errors -1

# Inspect errors for 2 commits prior (-2)
gitmap pipeline errors -2

# Inspect errors for 3 commits prior (-3)
gitmap pipeline errors -3
```

---

## 4. Pipeline History Tree

Display a visual ASCII tree of recent commit statuses, run durations, and failure states:

```bash
# View pipeline history for the last 5 commits
gitmap pipeline history

# View pipeline history for the last 10 commits
gitmap pipeline history -n 10
```

---

## 5. Feed Pipeline Errors to Antigravity (`agy fix`)

Extract failure step logs and git history, generate 4-part Root Cause Analysis (RCA) prompt into `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`, inject prompt directly into Antigravity IDE/CLI, and copy to clipboard:

```bash
# Fix pipeline errors for current repository
gitmap pipeline errors agy fix

# Using top-level shortcut
gitmap aef

# Batch fix failing pipelines across all tracked projects (batch of 3)
gitmap pipeline errors agy fix --all

# Run again to automatically advance to the next batch of failing projects
gitmap pipeline errors agy fix --all

# Custom batch limit of 5 projects
gitmap pipeline errors agy fix --all --limit 5

# Reset batch cursor back to the start
gitmap pipeline errors agy fix --reset-batch
```

---

## 6. Clear Pipeline Logs and Database (`clear`, `errors clear`)

Clean and purge pipeline logs, run artifacts, error reports, and database records:

```bash
# Clear pipeline logs & database for current repo (run inside repo)
gitmap pipeline clear -y

# Clear pipeline logs & database for a specific repo (run from anywhere)
gitmap pipeline clear alimtvnetwork/gitmap-v28 -y

# Clear error reports & logs for current repo
gitmap pipeline errors clear
```

---

## 7. Pipeline-AI Live Error Streaming & Fast-Forward Remediation

Automatically delay (`-t <seconds>`), query pipeline status, and stream live errors without waiting for the full suite to finish:

```bash
# Check status with dynamic auto-delay (delays then queries)
gitmap pipeline-ai status -t 45

# Structured JSON output with live error diagnostics and nextAiCommand
gitmap pipeline-ai status --json

# When errors occur during run, NextAiCommand automatically switches:
# Next Action: gitmap pipeline fix agy
# Stop waiting for remaining jobs. Fix current failure now so CI/CD can rerun in parallel.

# Directly inspect errors with AI remediation guidance
gitmap pipeline-ai errors
```

---

## See Also

- [agy](agy.md) — Comprehensive guide for Google Antigravity prompts, reruns, and workspaces
- [storage](storage.md) — Inspect disk space, pipeline logs, and SQLite database storage
- [ui](ui.md) — Browser dashboard to view pipeline telemetry and use the web terminal
- [llm](llm.md) — Guidance for autonomous AI agents monitoring CI/CD pipelines
