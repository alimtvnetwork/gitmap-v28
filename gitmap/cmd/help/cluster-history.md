# gitmap cluster history

View the audit trail of past cluster executions.

## Usage

```
gitmap cluster history [RunRef]
```

## Description

Every cluster execution (e.g., `gitmap sc ps "..."`) is assigned a unique `RunRef` format ID (e.g., `RUN-20260817-001`). These runs and their results are stored locally in the SQLite `ClusterRun` and `ClusterExecResult` tables.

- Omitting `[RunRef]` lists all past cluster runs in a table.
- Providing `[RunRef]` expands the results to show execution outcome, exit code, and stdout/stderr per-node.

## Table Columns (Listing Mode)

| Column | Description |
|--------|-------------|
| RUN ID | The unique RunRef string |
| COMMAND| The raw command string executed |
| TARGET | The target selector used (e.g., servers-clients) |
| SUCCESS| Number of nodes that successfully executed |
| FAILED | Number of nodes that failed |
| TIME   | Started at / Duration |

## Examples

```
gitmap cluster history
gitmap cluster history RUN-20260817-001
```
