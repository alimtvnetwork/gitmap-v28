# gitmap pr

> **Status (v6.274.0):** **LIVE.** Spec: `02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md`.

Automated PR & merge commit replay engine with SQLite split-database tracking.

## Alias

- `pull-request`
- `pr-in` (alias for `pr in`)
- `pr-clean` / `pr-rm` (feature branch pruning)
- `pr-list` (audit recorded PRs)

## Usage

    gitmap pr LEFT RIGHT [flags]
    gitmap pr right LEFT RIGHT [flags]
    gitmap pr left LEFT RIGHT [flags]
    gitmap pr in <target> <inputs...> [flags]
    gitmap pr-clean [repo] [--yes]
    gitmap pr-list [repo] [--json]

## How It Works & Architecture

The `pr` command family checks the repository where commits are coming from (source) and the target destination repository:

1. **Commit Traversal & Merge Detection:**
   Walks commits chronologically. When a merge commit is detected (`git rev-parse <sha>^@` > 1), GitMap avoids collapsing the merge. Instead, it simulates a genuine Pull Request workflow:
   - Creates a dedicated feature branch (`feature/<slug>` or `pr/<id>`).
   - Replays the branch commits onto the feature branch.
   - Generates a structured Markdown PR description.
   - Merges the feature branch into mainline with `--no-ff` (no fast-forward).
   - Records the PR metadata in SQLite (`.gitmap/data/pr/<slug>/sql.db`).

2. **SQLite Split-DB Storage Engine:**
   Every release, feature branch, and pull request is preserved in `.gitmap/data/pr/<slug>/sql.db` across tables:
   - `PullRequest`: ID, Title, SourceBranch, TargetBranch, Status (`open`, `merged`, `closed`), Description, CreatedAt, MergedAt.
   - `PrRelease`: TagName, CommitSha, PrNumber, ReleasedAt.
   - `PrBranch`: BranchName, PrNumber, IsMerged, IsDeleted.

3. **Terminal UI & Visual Execution Graph:**
   During replay, GitMap renders styled terminal status badges (`[REPLAYING]`, `[PR CREATED]`, `[PR MERGED]`, `[SNAPSHOT SYNCED]`), real-time commit counters, and concludes with an ASCII/Unicode visual execution graph:
   ```text
   Mainline:        ●────●────●─────────● (HEAD)
                         \         /
   Feature Branch:        ○───○───○ [PR #1: feat/oauth]
   ```

4. **Final Snapshot Synchronization:**
   At the end of replay, GitMap guarantees byte-for-byte fidelity with the source release:
   - Scans and prunes any stray files in target that do not exist in source.
   - Copies all source files to target.
   - If diffs remain, creates a final synchronization commit:
     `chore(sync): synchronize final repository snapshot tree to match source <shortsha>`

5. **PR Branch Pruning (`pr-clean`):**
   Inspects merged PR branches in SQLite, renders an inspection summary table in the terminal, and prompts for confirmation before safely deleting merged feature branches (`git branch -D`).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--mirror` | `false` | Delete target files not present in source commit |
| `--include-merges` | `true` | Preserve merge commits through PR simulation |
| `--pr <mode>` | `merges` | PR creation mode (`merges`, `tags`, `all`, `release`) |
| `--since <sha\|date>` | `(auto)` | Replay commits after specific commit or date |
| `--strip <regex>` | `(config)` | Strip pattern from commit messages (repeatable) |
| `--drop <regex>` | `(config)` | Drop commits matching regex (repeatable) |
| `--conventional` | `true` | Normalize commit messages to Conventional Commits |
| `--dry-run` | `false` | Inspect planned replay and branches without writing |
| `--yes / -y` | `false` | Skip confirmation prompt |
| `--no-push` | `false` | Keep commits local; do not push to remote |

## Help Group ("HG") Discoverability

Discover PR commands under the PR help group:

```bash
# Filter help by PR group
gitmap help pr
gitmap help --group pr

# List all help groups including PR
gitmap completion --list-help-groups
```

## Examples

```bash
# Replay commits from upstream into mainline in PR mode
gitmap pr ./upstream-repo ./mainline-repo

# Replay commits right-to-left
gitmap pr left ./source ./target

# Multi-source chronological PR replay into target repo
gitmap pr in ./target ./service-a ./service-b

# List all recorded pull requests with status
gitmap pr-list

# Clean up merged PR branches with interactive prompt
gitmap pr-clean --yes
```

## JSON Scripting Examples

Machine-readable execution payload:

```json
{
  "command": "pr",
  "direction": "left-to-right",
  "source": "./upstream-repo",
  "target": "./mainline-repo",
  "prMode": "merges",
  "includeMerges": true,
  "finalSnapshotSync": true,
  "options": {
    "dryRun": false,
    "conventional": true,
    "stripPatterns": ["^WIP:"],
    "dropPatterns": ["^test: drop-me"]
  }
}
```

Audit PR list JSON output:

```bash
gitmap pr-list --json
```

```json
[
  {
    "id": 1,
    "prNumber": 1,
    "title": "PR #1: Replay merge 7a3f8c",
    "sourceBranch": "feature/7a3f8c",
    "targetBranch": "main",
    "status": "merged",
    "mergedAt": 1726804800
  }
]
```
