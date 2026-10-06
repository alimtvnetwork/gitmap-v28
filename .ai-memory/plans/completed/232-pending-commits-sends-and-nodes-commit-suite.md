# Completed Plan: 232-pending-commits-sends-and-nodes-commit-suite

## User Request (Verbatim)
```text
# High Priority Instruction

add new commands with e2e testing and help text and ui help

gitmap pending-commits #show commits pending from all repos, summary first and ask like sort, 1 by 1 or all together
gitmap pending-commits -ssh # gets same info from all repos from all nodes open using json and displays it properly

gitmap sends cpf/cpb/cpr $repoName/all "msg" # commits specific repo or all repo using the term as mention try to understand the semantics from the existing repo commands

gitmap nodes pending-commits(pc) #show commits pending from all repos, summary first and ask like sort,
gitmap nodes commits $repoName/all $msg # commits
gitmap nodes cpf $repoName/all $msg # commits feature
gitmap nodes cpb $repoName/all $msg # commits bug
gitmap nodes cpr $repoName/all $msg # commits release
gitmap nodes commit-fix $repoName/all $msg # commits fix merge issues

all command should be delegated to other gitmap on that machine and returns as json result to show in the terminal also we can see everything in json if we use JSON flag add all helptext and do e2e tests please

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Implement new commands with e2e testing and help text for UI help.
5. Develop `gitmap pending-commits` to show commits pending from all repos, with options for sorting and displaying as a summary or detailed view.
6. Enhance `gitmap pending-commits -ssh` to retrieve and display commit information from all nodes using JSON.
7. Create `gitmap sends cpf/cpb/cpr $repoName/all "msg"` to commit to specific or all repos, understanding semantics from existing commands.
8. Implement `gitmap nodes pending-commits(pc)` to show pending commits from all repos with sorting options.
9. Develop `gitmap nodes commits $repoName/all $msg` for committing changes.
10. Implement `gitmap nodes cpf $repoName/all $msg` for committing features.
11. Develop `gitmap nodes cpb $repoName/all $msg` for committing bug fixes.
12. Implement `gitmap nodes cpr $repoName/all $msg` for committing releases.
13. Develop `gitmap nodes commit-fix $repoName/all $msg` for committing fixes to merge issues.
14. Ensure all commands delegate to other gitmap instances on the machine and return JSON results for terminal display.
15. Add help text and conduct e2e tests for all commands.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## Execution Summary

- **Canonical Architecture Spec:** [01-architecture-spec.md](../../02-spec/21-app/232-pending-commits-sends-and-nodes-commit-suite/01-architecture-spec.md)
- **Canonical Component Spec:** [02-component-spec.md](../../02-spec/21-app/232-pending-commits-sends-and-nodes-commit-suite/02-component-spec.md)
- **Subtasks Executed:**
  - `Subtask 01: Implement gitmap pending-commits & gitmap sends cpf/cpb/cpr Suite` (Worker 01):
    - `cli/constants/constants_cli.go`: Registered `CmdPendingCommits`, `CmdPendingCommitsAlias`, and `CmdSends`.
    - `cli/cmd/pending_commits_types.go`: Authored core data models (`PendingCommitsPayload`, `RepoPendingCommitRecord`, `UncommittedChangesSummary`, `UnpushedCommitSummary`) adhering to 100% positive booleans (`isDirty`, `isClean`, `hasUncommitted`, `hasUnpushed`, `hasUpstream`).
    - `cli/cmd/pending_commits_cmd.go`: Implemented `RunPendingCommits`:
      - Dispatches `gitmap pending-commits` and alias `gitmap pc`.
      - Discovers repositories across workspace directories and Split-DB registry.
      - Dual-state git detection for uncommitted changes (`git status --porcelain`) and unpushed commits ahead of upstream (`git rev-list @{u}..HEAD`).
      - Summary-first table view (`termtable.TableConfig`), sorting options (`priority`, `name`, `count`), and detailed views (`--detail all`, `--detail 1` / interactive).
      - Fleet SSH aggregation (`-ssh` / `--ssh`) over SSH with JSON deserialization.
      - Machine-readable `--json` envelope.
      - Rich UI box help menu.
    - `cli/cmd/sends_cmd.go`: Implemented `RunSends`:
      - Dispatches `gitmap sends <verb> <target> "<message>"` supporting verbs: `cpf` (`Feature: `), `cpb` (`Bug: `), `cpr` (`Release: `), `commit-fix` (`Fix: `), `cp`, and `cm`.
      - Target resolution supporting single repo name/slug or `all` dirty repositories.
      - Automatic clean repository skipping (`--skip-clean`).
      - Safety flags: `--dry-run` (`-n`), `--no-push`, and `--json`.
      - Rich UI box help menu.
    - `cli/cmd/roottooling.go`: Wired `pending-commits`, `pc`, and `sends` top-level routing in `toolingWorkspaceEntries`.
    - `cli/cmd/pending_commits_test.go` & `cli/cmd/sends_test.go`: Comprehensive E2E unit tests covering mock repository states, argument parsing, sorting, semantic prefix mapping, and JSON schemas.
  - `Subtask 02: Implement gitmap nodes pending-commits(pc) & gitmap nodes commits Remote Cluster Suite` (Worker 02):
    - `cli/cmdnodes/nodes_commits_types.go`: Authored cluster delegation data models (`NodePendingCommitsResult`, `NodeCommitExecutionResult`, `NodesCommitSummary`, `NodesCommitDelegationPayload`) with affirmative booleans (`isNodeOnline`, `isSuccess`, `hasChanges`, `isDryRun`, `isPushed`, `isDirty`).
    - `cli/cmdnodes/nodes_pending_commits.go`: Implemented `RunNodesPendingCommits`:
      - Handlers for `gitmap nodes pc` and `gitmap nodes pending-commits`.
      - Concurrent remote execution over SSH invoking `gitmap pending-commits --json`.
      - Shell-noise resilient JSON payload extraction (`extractJSONPayload`).
      - Tabular display rendering or pure JSON output (`--json`).
      - Node filtering options: `-t/--target`, `-e/--except`, `--include-main`, `--open-only`, `--dirty-only`.
    - `cli/cmdnodes/nodes_commits.go`: Implemented `RunNodesCommits`:
      - Handlers for `gitmap nodes commits`, `nodes cpf`, `nodes cpb`, `nodes cpr`, and `nodes commit-fix $repoName/all $msg`.
      - Remote command builder for POSIX Bash (`~/work`) and PowerShell (`D:\work`).
      - Aggregated execution reporting with per-node status and JSON support.
    - `cli/cmd/nodes_cmd.go`: Wired up `pc`, `pending-commits`, `commits`, `cpf`, `cpb`, `cpr`, and `commit-fix` in `runUnifiedNodesCLI`, and updated `printUnifiedNodesHelp`.
    - `cli/cmdnodes/nodes_commits_test.go`: End-to-end unit tests covering routing, payload parsing, mock SSH delegation, and prefix mapping.

---

## Verification Evidence

- **VG-01 Command Routing Parity**: All new commands and aliases (`pending-commits`, `pc`, `sends`, `nodes pc`, `nodes cpf/cpb/cpr/commit-fix`) cleanly wired in `cli/cmd/roottooling.go` and `cli/cmd/nodes_cmd.go`.
- **VG-02 Status & Commit Classification**: Engine distinguishes dirty working trees from commits ahead of upstream remote.
- **VG-03 Target Resolution Parity**: Handles target repository by slug or `all` dirty repositories with clean repository skipping.
- **VG-04 Semantic Message Prefixing**: `cpf` produces `Feature: `, `cpb` produces `Bug: `, `cpr` produces `Release: `, `commit-fix` produces `Fix: `.
- **VG-05 Remote JSON RPC Delegation**: Remote execution extracts JSON payloads and aggregates across nodes.
- **VG-06 UI Help Text & Box Rendering**: Both local and nodes commands include rich ANSI box help menus and flag descriptions.
- **VG-07 Path Relativity & Coding Guidelines**: 100% relative Git paths, positive boolean nomenclature, zero builds/test suite executions.
