# gitmap merge-in

Chronological multi-source commit replay and merge engine (alias for `gitmap commit-in` / `cin`).

## Usage

```bash
gitmap merge-in <target> <input1,input2,...> [flags]
gitmap merge-in <target> all                 [flags]
gitmap merge-in <target> -5                  [flags]
gitmap cin      <target> all                 [flags]
```

## Description

Walks one or more SOURCE git repositories in author-date order and APPENDS each commit (preserving both `AuthorDate` and `CommitterDate`) into a TARGET git repository.

- `all`: Walk all versioned sibling repositories (`<target>-v1`, `<target>-v2`, …)
- `-N`: Walk the latest `N` versioned siblings only
- Replay is fully idempotent via ShaMap stored in SQLite

## Flags

| Flag | Default | Description |
|---|---|---|
| `--default`, `-d` | `off` | Load the default profile bound to `<target>` |
| `--profile <name>`| `—` | Load `.gitmap/commit-in/profiles/<name>.json` |
| `--save-profile <name>` | `—` | Persist current options as named profile |
| `--set-default` | `off` | Mark saved profile as default for `<target>` |
| `--dry-run` | `off` | Plan and preview actions; never execute git commit |
| `--conflict <mode>` | `ForceMerge` | Conflict handling mode (`ForceMerge` or `Prompt`) |
| `--pr all\|tags\|release` | `—` | Create GitHub PR for replayed commits |
| `-h`, `--help` | `off` | Show help menu |

## Examples

```bash
# Replay all versioned siblings into canonical repo and save profile
gitmap merge-in ./canonical all --save-profile Default --set-default

# Replay specific remote forks oldest to newest
gitmap merge-in ./canonical https://github.com/org/v1.git,https://github.com/org/v2.git

# Preview commit replay without modifying disk
gitmap merge-in ./canonical -3 --dry-run
```
