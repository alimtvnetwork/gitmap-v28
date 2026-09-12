# Pipeline-AI

Autonomous CI/CD pipeline monitoring with automatic preflight delay and next AI command recommendation.

## Aliases

pl-ai, plai, pipeline_ai

## Usage

    gitmap pipeline-ai <subcommand> [flags]
    gitmap pl-ai <subcommand> [flags]

## Subcommands

| Subcommand | Description |
|---|---|
| status | Auto-delay (default 20s or -t <sec>) then query status and next action |
| eta, etc, wt | Auto-delay then query remaining ETA seconds |

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| -t, --time <seconds> | integer | 20 | Delay in seconds before checking pipeline (minimum: 20s) |
| -d, --delay <seconds> | integer | 20 | Alias for -t / --time |
| --json | boolean | false | Output structured JSON with nextAiCommand |

## Examples

```bash

# Auto-delay for 20s then query pipeline status

gitmap pipeline-ai status

# Auto-delay for 75s (based on previous ETA) then query

gitmap pipeline-ai status -t 75

# Machine-readable JSON output for AI automation

gitmap pipeline-ai status -t 30 --json
```
