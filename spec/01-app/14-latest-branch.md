# Latest Branch Command

## Overview

`gitmap latest-branch` (alias: `lb`) finds the most recently updated
remote branch by commit date and displays its name, SHA, date, and
subject line.

This is useful after a `git fetch` to quickly identify which branch
received the most recent push — especially in multi-branch workflows
where the "active" branch changes frequently.

## Command Signature

```
gitmap latest-branch [flags]
gitmap latest-branch <n>
gitmap lb [flags]
gitmap lb <n>
```

A bare integer argument (e.g. `gitmap lb 3`) is shorthand for
`--top <n>` — it shows the N most recently updated remote branches.

## Behavior

1. **Validate** — confirm the current directory is inside a Git repo
   (`git rev-parse --is-inside-work-tree`).
2. **Fetch** — unless `--no-fetch` is set, run `git fetch --all --prune` to update remote refs.
3. **List remote branches** — run `git branch -r`, trim whitespace,
   exclude `HEAD ->` pointer lines.
4. **Filter by remote** — if `--all-remotes` is not set, keep only
   branches matching the `--remote` value (default: `origin`).
5. **Read tip commits** — for each remote branch, run
   `git log -1 --format="%cI|%H|%s" <ref>` to get ISO-8601 commit
   date, full SHA, and subject.
6. **Sort** — order by commit date descending.
7. **Pick latest** — select the first (most recent) entry.
8. **Resolve branch name** — run
   `git for-each-ref --points-at=<sha> refs/remotes/<remote> --format="%(refname:short)"`
   to find which branch(es) point exactly at the SHA. Strip the
   `<remote>/` prefix.
9. **Contains fallback** — if `--contains-fallback` is set and
   `--points-at` returned nothing, fall back to
   `git branch -r --contains <sha>` (filtered to the selected remote).
10. **Display** — print branch name(s), remote, SHA (short), commit
    date, subject, and the original remote ref used.
11. **Top N** — if `--top <n>` is set (n > 0), also display the top N
    most recently updated remote branches in a table.

## Flags

| Flag                  | Type   | Default    | Description                                      |
|-----------------------|--------|------------|--------------------------------------------------|
| `--remote <name>`     | string | `origin`   | Remote to filter branches against                |
| `--all-remotes`       | bool   | `false`    | Include branches from all remotes                |
| `--contains-fallback` | bool   | `false`    | Fall back to `--contains` if `--points-at` empty |
| `--top <n>`           | int    | `0`        | Show top N most recently updated branches        |
| `--format <fmt>`      | string | `terminal` | Output format: `terminal`, `json`, `csv`         |
| `--json`              | bool   | `false`    | Shorthand for `--format json`                    |
| `--no-fetch`          | bool   | `false`    | Skip `git fetch` (use existing remote refs)      |
| `--sort <order>`      | string | `date`     | Sort order: `date` (descending) or `name` (A-Z)  |
| `--filter <pattern>`  | string | `""`       | Filter branches by glob or substring pattern     |

## Output Format

### Default (single latest branch)

```
  Latest branch: feature/v1.5.1
  Remote:        origin
  SHA:           a1b2c3d
  Commit date:   06-Mar-2025 03:22 PM
  Subject:       Fix auth token refresh
  Ref:           origin/feature/v1.5.1
```

### With `--top 3`

Appends a table after the main output:

```
  Top 3 most recently updated remote branches (origin):
  DATE                           BRANCH                SHA      SUBJECT
  06-Mar-2025 03:22 PM           feature/v1.5.1        a1b2c3d  Fix auth token refresh
  05-Mar-2025 09:10 AM           main                  d4e5f6a  Merge PR #42
  04-Mar-2025 05:45 PM           release/v2.3.0        b7c8d9e  Bump version
```

### With `--format csv`

Outputs CSV with a header row to stdout. When combined with `--top`,
all N rows are included:

```
branch,remote,sha,commitDate,subject,ref
feature/v1.5.1,origin,a1b2c3d,06-Mar-2025 03:22 PM,Fix auth token refresh,origin/feature/v1.5.1
```

### With `--format csv --top 3`

```
branch,remote,sha,commitDate,subject,ref
feature/v1.5.1,origin,a1b2c3d,06-Mar-2025 03:22 PM,Fix auth token refresh,origin/feature/v1.5.1
main,origin,d4e5f6a,05-Mar-2025 09:10 AM,Merge PR #42,origin/main
release/v2.3.0,origin,b7c8d9e,04-Mar-2025 05:45 PM,Bump version,origin/release/v2.3.0
```
### With `--json` (or `--format json`)

```json
{
  "branch": ["feature/v1.5.1"],
  "remote": "origin",
  "sha": "a1b2c3d",
  "commitDate": "06-Mar-2025 03:22 PM",
  "subject": "Fix auth token refresh",
  "ref": "origin/feature/v1.5.1"
}
```

### With `--json --top 3`

```json
{
  "branch": ["feature/v1.5.1"],
  "remote": "origin",
  "sha": "a1b2c3d",
  "commitDate": "06-Mar-2025 03:22 PM",
  "subject": "Fix auth token refresh",
  "ref": "origin/feature/v1.5.1",
  "top": [
    {
      "branch": "feature/v1.5.1",
      "sha": "a1b2c3d",
      "commitDate": "06-Mar-2025 03:22 PM",
      "subject": "Fix auth token refresh"
    },
    {
      "branch": "main",
      "sha": "d4e5f6a",
      "commitDate": "05-Mar-2025 09:10 AM",
      "subject": "Merge PR #42"
    },
    {
      "branch": "release/v2.3.0",
      "sha": "b7c8d9e",
      "commitDate": "04-Mar-2025 05:45 PM",
      "subject": "Bump version"
    }
  ]
}
```


| Condition                        | Message                                                              |
|----------------------------------|----------------------------------------------------------------------|
| Not inside a Git repo            | `Error: not inside a Git repository.`                                |
| No remote branches found         | `Error: no remote-tracking branches found for remote '<name>'.`      |
| No remote branches (all remotes) | `Error: no remote-tracking branches found on any remote.`            |
| Cannot read commit info          | `Error: could not read commit info for remote branches.`             |

## Implementation Notes

- All git commands use `os/exec` directly (no shell wrappers).
- The `|` delimiter in the `--format` string is safe because commit
  subjects are the last field (split with limit 3).
- Branch name resolution strips `<remote>/` prefix using
  `strings.TrimPrefix`.
- SHA display is truncated to 7 characters for readability.
- The command does **not** require `.gitmap/output/` or a previous scan.

## File Layout

| File                              | Purpose                                      |
|-----------------------------------|----------------------------------------------|
| `cmd/latestbranch.go`             | Config struct, flag parsing, orchestrator    |
| `cmd/latestbranchresolve.go`      | Result type, branch name resolve helpers     |
| `cmd/latestbranchoutput.go`       | JSON/CSV/terminal output formatters          |
| `gitutil/latestbranch.go`         | Core git operations: list, filter, sort      |
| `gitutil/latestbranchresolve.go`  | ReadBranchTips, ResolvePointsAt, ResolveContains |
| `constants/constants.go`          | Command name, alias, messages, flags         |
