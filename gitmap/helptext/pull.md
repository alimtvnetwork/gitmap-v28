# gitmap pull

Pull a specific tracked repository by slug, group, or all at once.

## Alias

p

## Usage

    gitmap pull <repo-name> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -A, --alias \<name\> | — | Target a repo by its alias |
| --group \<name\> | — | Pull all repos in a group |
| --all | false | Pull all tracked repos |
| --verbose | false | Enable verbose logging |
| --parallel \<N\> | 1 | Run up to N pulls concurrently (worker pool) |
| --only-available | false | Skip repos whose latest probe reports no new tag |
| --stop-on-fail | false | Halt the batch after the first failure |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Pull a single repo by slug

    gitmap pull my-api

**Output:**

    Pulling my-api (main)...
    remote: Enumerating objects: 5, done.
    remote: Counting objects: 100% (5/5), done.
    Already up to date.

### Example 2: Pull all repos in a group

    gitmap p --group backend

**Output:**

    Pulling 5 repos in group 'backend'...
    [1/5] billing-svc (main)... updated (3 new commits)
    [2/5] auth-gateway (main)... Already up to date.
    [3/5] payments-api (main)... updated (1 new commit)
    [4/5] user-svc (develop)... Already up to date.
    [5/5] notification-svc (main)... Already up to date.
    ✓ 5 repos pulled (2 updated, 3 up to date)

### Example 3: Pull all tracked repos with verbose logging

    gitmap pull --all --verbose

**Output:**

    [verbose] Log file: gitmap-debug-2025-03-10T14-30.log
    Pulling 42 tracked repos...
    [1/42] my-api (main)... updated (7 commits)
    [2/42] web-app (develop)... Already up to date.
    [3/42] billing-svc (main)... updated (2 commits)
    ...
    ✓ 42 repos pulled (12 updated, 30 up to date)
    [verbose] Debug log written

### Example 4: Pull by alias

    gitmap pull -A api

**Output:**

    Pulling my-api (main)...
    Already up to date.

### Example 5: Parallel pull, only what's actually new

First refresh the probe so `--only-available` has fresh data:

    gitmap probe --all
    gitmap pull --all --only-available --parallel 4

**Output:**

    [1/6] billing-svc (main)... updated (3 new commits)
    [2/6] auth-gateway (main)... updated (1 new commit)
    [3/6] payments-api (main)... updated (5 new commits)
    [4/6] user-svc (develop)... updated (2 new commits)
    [5/6] notification-svc (main)... updated (1 new commit)
    [6/6] config-svc (main)... updated (4 new commits)
    ✓ Pull complete: 6 succeeded, 0 failed (4-way parallel)

## See Also

- [scan](scan.md) — Scan directories to populate the database
- [clone](clone.md) — Clone repos from output files
- [status](status.md) — Check repo statuses before pulling
- [group](group.md) — Manage groups for targeted pulls
- [alias](alias.md) — Manage repo aliases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pull
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
