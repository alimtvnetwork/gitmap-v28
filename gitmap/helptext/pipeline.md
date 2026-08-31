# Pipeline

Query live CI/CD pipeline runs, calculate remaining wait times (ETA), fetch step failure logs, and export telemetry.

## Aliases

pl

## Usage

    gitmap pipeline <subcommand> [flags]
    gitmap pl <subcommand> [flags]

## Subcommands

| Subcommand     | Description                                                          |
|----------------|----------------------------------------------------------------------|
| status         | Live CI/CD execution state, active workflow, ETA, and pending PRs    |
| waittime, eta  | Remaining estimated wait time in seconds (machine-friendly integer)  |
| logs           | Full step logs for the latest workflow run                           |
| error-logs     | Failure logs for the latest failed workflow step                     |
| help           | Show this pipeline command suite documentation                       |

## Flags

| Flag         | Type    | Default | Description                                                     |
|--------------|---------|---------|-----------------------------------------------------------------|
| --json       | boolean | false   | Output structured JSON payload for scripting and AI agents      |
| --file       | string  | ""      | Write error logs or telemetry to the specified file path        |
| --tempfile   | string  | ""      | Write error logs to `.lovable/temp/<filename>` (configurable)   |
| --split      | boolean | false   | Separate failure logs across multiple jobs into distinct files  |

## Top-Level Shortcuts

| Shortcut           | Equivalent Command                  |
|--------------------|-------------------------------------|
| gitmap error-logs  | gitmap pipeline error-logs          |
| gitmap logs        | gitmap pipeline logs                |
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

## See Also

- [ui](ui.md) — Launch the browser dashboard to view pipeline telemetry and use the web terminal
- [llm](llm.md) — Guidance for autonomous AI agents monitoring CI/CD pipelines
