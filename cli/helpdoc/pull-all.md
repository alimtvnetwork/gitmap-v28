# gitmap pull-all

> **Generated — do not hand-edit.** This topic is rendered from the
> `cmdpull.PullAllHelpDisplay` HelpDisplay struct. Regenerate with:
>
> `gitmap py 03-ai-scripts/51-helptext-generator.py --topics pull-all`

```text
Batch-pull every tracked repository in the catalog (gitmap pull-all)

  Usage:
    gitmap pull-all [flags]
    gitmap pa [--probe] [--status | --table | --json]
    gitmap pa --probe [-y]
    gitmap pat [flags]
    gitmap pull all table [flags]

Aliases:
  pa  Batch pull every tracked repository (fast concise mode)
  pat  Batch pull with the full post-pull status table
  pull-all-table  Batch pull with the full post-pull status table
Flags:
  --probe  Probe GitHub account & workspace for companion repositories (repo-secrets & repo-cache) before pulling
  --status, --table  Render the full post-pull repository status table
  --json  Suppress interactive progress bar and emit JSON summary to stdout
  --ssh, -s  Dispatch pull-all across SSH cluster fleet and Local VM concurrently
  --verbose  Enable verbose logging
  -p, --parallel <N>  Explicitly set concurrency worker count (default: 0 = auto; overrides presets)
  --auto-scale  Dynamically adjust worker concurrency based on CPU core availability (default)
  --high-perf, --turbo  Maximum concurrency preset when CPU pressure is low (up to 12 parallel workers)
  --low-cpu, --conservative  Throttled concurrency preset when system is under high CPU pressure (1-2 workers)
  --only-available  Skip repos whose latest probe reports no new tag
  --stop-on-fail  Halt the batch after the first failure
Examples:
  gitmap pa  Default fast batch pull (shows only active/updated repos)
  gitmap pa --status  Full post-pull status table
    e.g. gitmap pat
  gitmap pa --json  Machine-readable JSON summary
  gitmap pa --ssh  SSH fleet batch pull (fleet + local VM concurrently)
  gitmap pull-all --parallel 8 --only-available --stop-on-fail  8-way parallel, only repos with new commits, stop on first failure
  gitmap pa --probe  Probe companion repositories (repo-secrets & repo-cache) before pulling
    e.g. gitmap pa --probe -y
  gitmap pa --probe --status  Probe companion repos and render the full status table
    e.g. gitmap pa --probe --json
  gitmap pa --auto-scale  CPU auto-scaling and concurrency presets
    e.g. gitmap pa --high-perf

  All pull flags are forwarded verbatim; --all is injected automatically and is idempotent (passing it again is a no-op).

* Run gitmap scan first to populate the database (see scan help)
```
