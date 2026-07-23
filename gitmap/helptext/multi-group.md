# gitmap multi-group

Select multiple groups for batch operations (pull, status, exec).

## Alias

mg

## Usage

    gitmap multi-group <group1,group2,...|clear|pull|status|exec>

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)
- Create groups with `gitmap group create` (see group.md)

## Examples

### Example 1: Select multiple groups

    gitmap mg backend,frontend

**Output:**

    Multi-group set: backend, frontend
    Total repos: 8 (5 backend + 3 frontend)

### Example 2: Pull all repos across selected groups

    gitmap mg pull

**Output:**

    Pulling 8 repos (groups: backend, frontend)...
    [backend]
      [1/5] billing-svc (main)... updated (2 commits)
      [2/5] auth-gateway (main)... Already up to date.
      [3/5] payments-api (main)... Already up to date.
      [4/5] user-svc (develop)... updated (1 commit)
      [5/5] notification-svc (main)... Already up to date.
    [frontend]
      [1/3] web-app (develop)... updated (5 commits)
      [2/3] admin-panel (main)... Already up to date.
      [3/3] landing-page (main)... Already up to date.
    ✓ 8 repos pulled (3 updated, 5 up to date)

### Example 3: Status across multiple groups

    gitmap mg status

**Output:**

    [backend]
    REPO             BRANCH   STATUS  AHEAD/BEHIND
    billing-svc      main     clean   0/0
    auth-gateway     main     dirty   1/0
    [frontend]
    web-app          develop  dirty   3/0
    admin-panel      main     clean   0/0
    ✓ 4 repos (2 dirty, 2 clean)

### Example 4: Clear multi-group selection

    gitmap mg clear

**Output:**

    ✓ Multi-group selection cleared
    (No active group — commands will target all repos)

## See Also

- [group](group.md) — Manage and activate single groups
- [pull](pull.md) — Pull a specific repo
- [status](status.md) — View repo statuses
- [exec](exec.md) — Run git across repos

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter multi-group
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
