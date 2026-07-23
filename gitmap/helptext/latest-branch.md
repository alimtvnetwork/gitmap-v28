# gitmap latest-branch

Find the most recently updated remote branch in the current repository.

## Alias

lb

## Usage

    gitmap latest-branch [--top N] [--format json|csv|terminal]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --top \<N\> | 1 | Number of branches to show |
| --format json\|csv\|terminal | terminal | Output format |
| --no-fetch | false | Skip git fetch before query |
| --sort date\|name | date | Sort order |
| --switch, -s | false | Checkout the resolved latest branch after printing the report |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Show the latest branch

    gitmap lb

**Output:**

    Latest branch: feature/auth-redesign
    Last commit:   2 hours ago (2025-03-10 14:30)
    Author:        developer@example.com
    Commit:        abc1234 — Add OAuth2 provider

### Example 2: Top 5 most recent branches

    gitmap lb 5

**Output:**

     #  BRANCH                    LAST COMMIT          AUTHOR
     1  feature/auth-redesign     2 hours ago          dev@example.com
     2  bugfix/login-fix          5 hours ago          dev@example.com
     3  main                      1 day ago            dev@example.com
     4  develop                   2 days ago           team@example.com
     5  feature/payments          3 days ago           dev@example.com

### Example 3: CSV output for scripting

    gitmap lb 3 --format csv

**Output:**

    branch,last_commit,author,commit_sha
    feature/auth-redesign,2025-03-10T14:30:00Z,dev@example.com,abc1234
    bugfix/login-fix,2025-03-10T09:15:00Z,dev@example.com,def5678
    main,2025-03-09T18:00:00Z,dev@example.com,ghi9012

### Example 4: JSON output without fetch

    gitmap latest-branch --format json --no-fetch

**Output:**

    {
      "branch": "feature/auth-redesign",
      "last_commit": "2025-03-10T14:30:00Z",
      "author": "dev@example.com",
      "commit_sha": "abc1234"
    }

### Example 5: Jump to the latest branch

    gitmap lb -s

**Output:**

      ▶ Switching to feature/auth-redesign...
    Switched to branch 'feature/auth-redesign'
    Your branch is up to date with 'origin/feature/auth-redesign'.

## See Also

- [branch](branch.md) — `gitmap b def` jumps to the default branch
- [status](status.md) — View repo branch and status info
- [release-branch](release-branch.md) — Create a release branch
- [watch](watch.md) — Live-refresh status dashboard

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter latest-branch
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
