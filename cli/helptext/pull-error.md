# gitmap pull-error / gitmap pulle

Inspect, query, filter, and remediate persistent Git pull errors recorded during single-repo or batch `pull-all` operations.

## Usage

```bash
gitmap pull-error [repo-slug|all] [flags]
gitmap pulle [repo-slug|all] [flags]
```

## Aliases

- `pull-errors`, `pulle`, `pull-e`

## Description

Whenever a `gitmap pull` or batch `gitmap pull-all` (`pa`) fails—due to diverged branches, merge conflicts, network drops, SSH authentication issues, or Windows Credential Manager (`wincredman`) failures—the exact error details, stack traces, node context, and actionable remediation commands are automatically persisted into a local SQLite Split-DB at `.gitmap/data/pull/errors/sql.db`.

`gitmap pull-error` lets developers and AI agents immediately query these failure diagnostics without rerunning expensive batch operations or manually digging through logs.

## Target Resolution

- **Explicit Repo Slug:** `gitmap pull-error <slug>` queries the latest recorded errors for the specified repository.
- **Current Directory:** If no target is specified and the current working directory is inside a Git repository, the target automatically resolves to the current repository's basename directory.
- **All Repositories:** Specifying `all` (`gitmap pull-error all`) or running outside a Git repository without arguments lists the latest failures across all tracked repositories.

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--json` | `-j` | Output diagnostic records in clean structured JSON for automated tooling or AI agents |
| `--ssh` | | Query pull errors originating from remote SSH fleet nodes |
| `--help` | `-h` | Display usage instructions and examples |

## Examples

### 1. View Latest Error for Current Repository

```bash
gitmap pull-error
```

Output:

```text
  ┌─ GitMap Pull Failure Diagnostic Card ──────────────────────────
  │ Repository:    cat-my-v12 (Node: localhost | v6.451.0)
  │ Timestamp:     2026-10-01 17:42:15 UTC (2m ago)
  │ Error Type:    Diverged Branches (Non-Fast-Forward)
  │
  │ Details:
  │   fatal: Not possible to fast-forward, aborting.
  │
  │ 💡 Remediation:
  │   gitmap pull cat-my-v12 --autostash
  │   (or: cd d:\work\cat-my-v12 && git merge --abort)
  └──────────────────────────────────────────────────────────────
```

### 2. View All Recorded Pull Failures Across Fleet

```bash
gitmap pulle all
```

### 3. Machine-Readable JSON for AI Agents

```bash
gitmap pull-error --json
```

Output:

```json
[
  {
    "errorId": "err-1727804535-cat-my-v12",
    "pullRunId": 104,
    "repoSlug": "cat-my-v12",
    "repoPath": "d:\\work\\cat-my-v12",
    "nodeId": "local-node",
    "nodeVersion": "v6.451.0",
    "errorType": "Diverged Branches",
    "errorText": "fatal: Not possible to fast-forward, aborting.",
    "remediationCmd": "gitmap pull cat-my-v12",
    "createdAt": "2026-10-01T17:42:15Z"
  }
]
```
