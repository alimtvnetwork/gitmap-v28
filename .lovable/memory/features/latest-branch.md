# Memory: features/latest-branch

The 'latest-branch' (alias 'lb') command identifies the most recently updated remote-tracking branches.

## File Layout (code-style compliant)
- `cmd/latestbranch.go` — config struct, flag parsing, orchestrator (under 200 lines)
- `cmd/latestbranchresolve.go` — result type, resolve helpers (under 200 lines)
- `cmd/latestbranchoutput.go` — JSON/CSV/terminal output formatters (under 200 lines)
- `gitutil/latestbranch.go` — core git operations: list, filter, sort, helpers (under 200 lines)
- `gitutil/latestbranchresolve.go` — ReadBranchTips, ResolvePointsAt, ResolveContains (under 200 lines)
- `gitutil/dateformat.go` — centralized date formatting
- `constants/constants.go` — command name, alias, messages, flags, sort orders, git args

## Flags
- `--remote <name>` — filter by remote (default: origin)
- `--all-remotes` — include all remotes
- `--contains-fallback` — use `git branch -r --contains` if exact SHA resolution fails
- `--top <n>` — show top N most recently updated branches
- `--format <fmt>` — output format: `terminal` (default), `json`, `csv`
- `--json` — shorthand for `--format json`
- `--no-fetch` — skip `git fetch --all --prune`
- `--sort <order>` — `date` (default, descending) or `name` (alphabetical)
- `--filter <pattern>` — filter branches by glob or substring

## Code Style
- All functions 8-15 lines, positive logic (no `!`), blank before return
- Config struct replaces multi-return parsing; positive booleans (shouldFetch, filterByRemote)
- Chained if+return replaces switch for format dispatch
