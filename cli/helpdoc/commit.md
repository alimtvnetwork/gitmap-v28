# commit

Stages all changes automatically and creates a flat git commit with the given message. No manual `git add` needed.

## Why use this instead of git commit?

Instead of typing two separate commands every time:

```bash
git add -A
git commit -m "your message"
```

You can run:

```bash
gitmap commit "your message"
```

All modified, new, and deleted files are automatically staged and committed flat without forced semantic prefixes.

## Aliases

- `gitmap commit "<message>"`
- `gitmap cm "<message>"`
- `gitmap commit-all "<message>"`
- `gitmap ca "<message>"`

## Semantic Commit Shortcuts & Full Forms

- `cpf <msg>`: `commit-push-feature` (Feature: <msg>)
- `cpb <msg>`: `commit-push-bug` (Bug: <msg>)
- `cpc <msg>`: `commit-push-chore` (Chore: <msg>)
- `cpr <msg>`: `commit-push-release` (Release: <msg>)
- `pcp <msg>`: `pull-commit-push` (pull rebase, stage, commit, and push)
- `pas`: `pull-all-ssh` (pull all repositories using SSH transport)

## Flags

- `-m "<message>"`: Commit message (optional, positional arguments are also accepted)
- `-p, --push`: Push to remote branch after committing
- `-n, --dry-run`: Preview staged changes without committing
- `-tz, --timezone <tz>`: Canonical IANA target timezone for schedule evaluation (e.g., `America/New_York`, `UTC`, `Asia/Tokyo`, or `Local`)
- `-sc, --schedule`: Enforce commit scheduling rules:
  - **Weekday policy (Mon–Fri)**: Commits must be completed before **08:00 AM** in target timezone. Attempting commits after 08:00 AM on weekdays is restricted.
  - **Weekend policy (Sat–Sun)**: Commits are completely unconstrained and permitted at any hour.
- `-h, --help`: Show help documentation

## Timezone & Commit Scheduling Policy

When `--schedule` (`-sc`) is active (or enabled in UI settings via `commit.schedule_active`), GitMap evaluates the current local time in the specified `--timezone` (`-tz`, default `commit.timezone` or `Local`):

- **Weekdays (Mon–Fri)**: Commits must be completed strictly before **08:00 AM**. Commits attempted after 08:00 AM on weekdays will trigger a closed-window restriction.
- **Weekends (Sat–Sun)**: Commits are unconstrained (24 hours a day, anytime).

## Examples

```bash
# Flat commit with quoted message:
gitmap commit "refactor: simplify split-db path resolution"
gitmap cm "fix: handle empty slice in scanner"

# Flat commit with unquoted message words:
gitmap cm update readme with new install steps

# Auto-stage and open git editor:
gitmap commit

# Auto-stage, commit flat, and push to remote:
gitmap commit "feat: add user telemetry" --push
gitmap cm "feat: add user telemetry" -p

# Enforce commit schedule under target timezone (weekday < 08:00 AM cutoff):
gitmap commit "chore: morning maintenance" --schedule --timezone America/New_York
gitmap cm "fix: quick morning fix" -sc -tz UTC
```
