# Pipeline

Query live CI/CD pipeline runs, calculate remaining wait times (ETA), fetch step failure logs, and export telemetry.

## Aliases

pl

## Usage

    gitmap pipeline <subcommand> [flags]
    gitmap pl <subcommand> [flags]

## Subcommands

| Subcommand     | Description                                                                           |
|----------------|---------------------------------------------------------------------------------------|
| status         | Live CI/CD execution state, active workflow, ETA, and pending PRs                     |
| waittime, eta  | Remaining estimated wait time in seconds (machine-friendly integer)                   |
| errors, err    | Aggregated failure logs for current or past commit by offset (`-1`, `-2`, `-3`)       |
| fix errors agy | Feed errors & git log with RCA prompt to Antigravity and queue follow-up (alias: `aef`)|
| history, hist  | Visual pipeline run tree and commit status summary for recent commits (`-n 5`)       |
| logs           | Full step logs for target commit, offset, or latest workflow run                      |
| help           | Show this pipeline command suite documentation                                        |

## Flags

| Flag         | Type    | Default | Description                                                     |
|--------------|---------|---------|-----------------------------------------------------------------|
| -f, --force  | boolean | false   | Bypass duplicate check and resend errors previously dispatched  |
| --json       | boolean | false   | Output structured JSON payload for scripting and AI agents      |
| --file       | string  | ""      | Write error logs or telemetry to the specified file path        |
| --tempfile   | string  | ""      | Write error logs to `.lovable/temp/<filename>` (configurable)   |
| --split      | boolean | false   | Separate failure logs across multiple jobs into distinct files  |
| --clip       | boolean | false   | Copy logs directly to system clipboard                         |
| --all        | boolean | false   | Retain verbose and passing lines without noise filtering        |
| -n <count>   | integer | 5       | Number of recent commits to display in pipeline history         |

## Top-Level Shortcuts

| Shortcut           | Equivalent Command                  |
|--------------------|-------------------------------------|
| gitmap aef         | gitmap pipeline fix errors agy      |
| gitmap pipeline-fix| gitmap pipeline fix errors agy      |
| gitmap error-logs  | gitmap pipeline error-logs          |
| gitmap errors      | gitmap pipeline errors              |
| gitmap logs        | gitmap pipeline logs                |
| gitmap history     | gitmap pipeline history             |
| gitmap waittime    | gitmap pipeline waittime            |
| gitmap eta         | gitmap pipeline eta                 |

## Examples

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

### Extract Failure Logs Directly to Temp File

```bash
$ gitmap pipeline error-logs --json --tempfile "ci-failure.json"
  ✓ Output written to .lovable\temp\ci-failure.json
```

### Inspect Errors for Previous Commit (-1, -2, -3)

```bash
# Aggregate all workflow errors for the previous distinct commit
$ gitmap pipeline errors -1

# Inspect errors for 3 commits prior
$ gitmap pipeline errors -3
```

### View Pipeline History Tree Across Recent Commits

```bash
$ gitmap pipeline history
```

### Export Full Pipeline Logs to Clipboard or Temp File

```bash
# Copy latest workflow run logs directly to system clipboard
$ gitmap pipeline logs --clip

# Save logs for a specific commit to temp file
$ gitmap pipeline logs 066e04b --tempfile "pipeline-run.log"
```

## See Also

- [storage](storage.md) — Inspect disk space, pipeline logs, and SQLite database storage
- [ui](ui.md) — Launch the browser dashboard to view pipeline telemetry and use the web terminal
- [llm](llm.md) — Guidance for autonomous AI agents monitoring CI/CD pipelines
