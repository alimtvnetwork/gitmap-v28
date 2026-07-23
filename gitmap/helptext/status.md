# gitmap status

Show a dashboard of repository statuses (branch, clean/dirty, ahead/behind).

## Alias

st

## Usage

    gitmap status [--group <name>] [--all]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -A, --alias \<name\> | — | Target a repo by its alias |
| --group \<name\> | — | Show status for repos in a group |
| --all | false | Show status for all tracked repos |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Status of all tracked repos

    gitmap status --all

**Output:**

    REPO             BRANCH     STATUS  AHEAD/BEHIND
    my-api           main       clean   0/0
    web-app          develop    dirty   2/1
    billing-svc      main       clean   0/0
    auth-gateway     feature/x  dirty   5/0
    shared-lib       main       clean   0/3
    ✓ 5 repos (2 dirty, 3 clean)

### Example 2: Status of a specific group

    gitmap st --group backend

**Output:**

    REPO             BRANCH   STATUS  AHEAD/BEHIND
    billing-svc      main     clean   0/0
    auth-gateway     main     dirty   1/0
    payments-api     main     clean   0/2
    ✓ 3 repos (group: backend) — 1 dirty, 2 clean

### Example 3: Status of a single repo by alias

    gitmap status -A api

**Output:**

    my-api  main  clean  0/0
    ✓ Up to date

## See Also

- [watch](watch.md) — Live-refresh status dashboard
- [scan](scan.md) — Scan directories to populate the database
- [group](group.md) — Manage repo groups
- [pull](pull.md) — Pull repos to sync changes
- [alias](alias.md) — Manage repo aliases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter status
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
