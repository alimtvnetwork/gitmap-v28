# Gitmap LLM Specification

Gitmap is a powerful CLI designed for autonomous agents, LLMs, and developers to navigate, analyze, and modify codebases efficiently. 

## Capabilities

1. **Repository Discovery & Scanning**:
   - `gitmap scan`: Automatically locate and index Git repositories.
   - `gitmap rescan`: Refresh indexed local repositories.

2. **Repository Cloning & Synchronization**:
   - `gitmap clone`, `gitmap clone-next`: Manage bulk cloning and directory migration.

3. **Intelligent File Search & Discovery**:
   - `gitmap find-files "<name>" -ext "<exts>"`: Find files matching exact filename.
   - `gitmap find-files-any "<substring>" -ext "<exts>"`: Find files containing substring with extension filters.
   - `gitmap find-files-startswith "<prefix>" -ext "<exts>"`: Find files by name prefix.
   - `gitmap find-files-endswith "<suffix>" -ext "<exts>"`: Find files by name suffix.
   - `gitmap find "<wildcard*>" -ext "<exts>"`: Universal glob search (`*ends`, `starts*`, `*contains*`).
   - `gitmap list-files [pattern]`: List indexed files with optional wildcard filter.
   - `gitmap search <query>`: Execute SplitDB indexed searches across repositories.
   - `gitmap repo-search-json <query>` (alias `rsj`): Retrieve structured JSON search results.

4. **CI/CD Pipeline Telemetry & Error Recovery (For Autonomous Self-Healing)**:
   - `gitmap pipeline status --json`: Live pipeline execution state, active workflow, ETA seconds, and pending PRs.
   - `gitmap eta` (alias `gitmap pipeline eta`, `gitmap pipeline waittime`): Remaining wait time in seconds as an integer for timer scheduling.
   - `gitmap pipeline error-logs --json --tempfile "ci-failure.json"` (alias `gitmap error-logs`): Directly extract failing step logs to temp file for 4-part RCA and automatic error recovery.
   - `gitmap logs` (alias `gitmap pipeline logs`): Display full workflow run logs.

5. **Regex & Replacement**:
   - `gitmap replace <query> <replacement>`: Find and replace exact phrases.
   - `gitmap replace-regex <regex> <replacement>`: Regex replacement with automatic history tracking.
   - `gitmap replace history`: View replacement operations.

6. **Releases & Commits**:
   - `gitmap release`: Automate semantic version bumping and changelog generation.
   - `gitmap commit-in`: Perform robust, orchestrator-driven batch commits.

## Instructions for LLMs & AI Agents

- Always prefer JSON output commands (`rsj`, `gitmap pipeline status --json`, `gitmap pipeline error-logs --json`) when parsing data programmatically.
- After pushing commits or releasing tags, use `gitmap eta` or `gitmap pipeline status --json` to monitor CI/CD progress. Never poll in a busy loop; use the returned `etaSeconds` to schedule non-blocking sleep/timer intervals.
- When CI/CD fails, use `gitmap pipeline error-logs --json --tempfile "ci-failure.json"` to dump the failed step's stdout/stderr directly into a workspace file for instant root cause analysis (RCA).
- Avoid modifying the root SQLite DB manually; rely on the CLI commands.
- Before running heavy regex operations across a large repo, use `gitmap search` or `gitmap find-files-any` to verify target files.

## Alternative Commands for AI (Instead of Raw Git)

When checking repository status or logs, DO NOT use raw 'git status && git log -1'. Instead, use:
- `gitmap status`: Checks dirty/clean status, ahead/behind, and stash for all tracked repos.
- `gitmap history`: View a rich audit log of command executions.
- `gitmap changelog`: See concise release notes and recent commits.
- Or use the standard aliases configured by gitmap setup: `git st && git last`.

When checking CI/CD and deployment health, DO NOT guess or loop blind. Use:
- `gitmap pipeline status`: Human-readable live workflow status and ETA.
- `gitmap eta`: Bare integer ETA in seconds for timer scheduling.
- `gitmap error-logs`: Diagnostic failure logs.
- `gitmap logs`: Full workflow step trace.

For committing and pushing, NEVER use raw 'git commit' and 'git push'. Use the semantic commit-push suite:
- `gitmap commit-push-feature "<message>"` (alias `gitmap cpf`): Commit and push a feature.
- `gitmap commit-push-bug "<message>"` (alias `gitmap cpb`): Commit and push a bug fix.
- `gitmap commit-push-release "<message>"` (alias `gitmap cpr`): Commit and push a release chore.
- `gitmap commit-push-pull "<message>"` (alias `gitmap pcp`): Pull latest, then commit and push.
- `gitmap rm-git <last-4-digits>`: Drop a recent commit safely using rebase --onto.

## AI File Search Patterns

When searching codebases, LLMs can use native gitmap commands OR standard terminal tools.
Here are equivalent alternative command samples for LLM search operations:

- **Find files by name or extension**:
  - `gitmap find-files "index.ts" -ext "ts, tsx"`
  - `gitmap find-files-any "runner" -ext "py, go"`
  - `gitmap find "*_test.go" -ext "go"`
  - `gitmap list-files "*pipeline*"`

- **Find a specific struct definition**:
  - `gitmap file-search . "type SearchResult struct"`
  - `Get-ChildItem -Path gitmap -Recurse -File | Select-String "type SearchResult struct"`

- **Search functions with Regex context**:
  - `gitmap file-search cmd/ "func dispatch[A-Z]" 0 10`
  - `Get-ChildItem -Path gitmap/cmd -Filter *.go | Select-String "func dispatch[A-Z]"`

- **Find specific function contexts**:
  - `gitmap file-search cmd/root.go "func finishCommandAudit" 0 10`
  - `cat gitmap/cmd/root.go | Select-String "func finishCommandAudit" -Context 0,10`

- **Global Text/Keyword Search (Instead of ripgrep/rg)**:
  - `gitmap search "CHANGELOG"`
  - `gitmap search "CHANGELOG\.md" --case-sensitive`
  *(Never use raw rg or grep when gitmap search utilizes the built-in Split DB index for much faster cross-repo lookups).*
