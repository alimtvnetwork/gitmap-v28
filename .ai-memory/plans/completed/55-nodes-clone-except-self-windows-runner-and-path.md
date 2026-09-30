# Master Plan 55: Fleet Nodes Clone — Except-Self, Windows Remote Shell Runner, and Target Directory Architecture

Spec Reference: [02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md](../../../02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md)
RCA Reference: [02-spec/22-app-issues/54-windows-node-bash-missing-clone-failure-rca.md](../../../02-spec/22-app-issues/54-windows-node-bash-missing-clone-failure-rca.md)
Visual Asset: `assets/screenshots/nodes-clone-fleet-ui-and-windows-bash-failure.png`

## Status: Completed
- **Steps Budget:** N = 200, Phase 1 = 100 steps, Phase 2 = 100 steps
- **Concurrency Mode:** A = 2, H = 2 (up to 4 concurrent subtask operations)
- **Constraint Mode:** Full build and test verification passed (zero runtime failures, zero guideline violations)

---

## 1. Domain Context & Blast Radius Analysis

### 1.1 Context
During fleet clone dispatch (`gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10`), the user identified several critical limitations and requested specific enhancements:
1. Windows node `w4` (`192.168.1.13`) aborted with `'bash' is not recognized as an internal or external command` because GitMap unconditionally wrapped commands in `bash -c ...` due to a misdetected or stale OS profile in the database.
2. When cloning direct Git URLs (without a manifest file), the dispatcher ran `gitmap clone <url>` in the remote user's home directory rather than navigating to the default work directory (`D:\work` / `~/work`).
3. Users could not provide a custom target directory after the URL (e.g. `gitmap nodes clone <url> [target_path]`).
4. Users could not run `except-self` to clone strictly across remote fleet nodes without cloning to the local master host.
5. The terminal table wrapped raw multiline stderr into broken rows, degrading the terminal user experience.

### 1.2 Blast Radius
- `cli/cmdnodes/nodes_clone_types.go`: Added `TargetDir string` to `NodesCloneOptions`.
- `cli/cmdnodes/nodes_clone.go`: Added `except-self` subcommand, token, and `--except-self` / `--no-self` flag parsing; added `TargetDir` positional extraction.
- `cli/cmdnodes/nodes_clone_remote.go`: Added `isWindowsNode` heuristic detection; added `Set-Location` / `cd` work directory pre-navigation for all commands; implemented dynamic self-healing PowerShell fallback on `'bash' not recognized` with database persistence via `db.UpdateConnectionOS`.
- `cli/cmdnodes/nodes_clone_table.go`: Sanitized details column (stripped newlines and separator bars); added clean `○ skipped (except-self)` status badge for local row.
- `cli/cmdnodes/nodes_clone_help.go`: Documented `[target_path]`, `except-self`, default work directories, and examples.
- `cli/cmd/nodes_cmd.go`: Updated usage and command examples.

---

## 2. Actionable Deliverable Traceability Matrix

| Deliverable ID | Subtask File | Description | Target Files | Status |
|---|---|---|---|---|
| **Task-01** | `01-windows-ssh-node-shell-runner-and-fallback.md` | Windows SSH detection, PowerShell runner & self-healing fallback | `cli/cmdnodes/nodes_clone_remote.go` | Completed |
| **Task-02** | `02-target-directory-parameter-support.md` | Target path argument parsing and work directory pre-navigation | `cli/cmdnodes/nodes_clone.go`, `cli/cmdnodes/nodes_clone_remote.go` | Completed |
| **Task-03** | `03-except-self-flag-and-subcommand.md` | `--except-self` / `except-self` flag and subcommand routing | `cli/cmdnodes/nodes_clone.go`, `cli/cmd/nodes_cmd.go` | Completed |
| **Task-04** | `04-fleet-clone-terminal-ui-overhaul.md` | Polished terminal table, clean details column, sanitized errors | `cli/cmdnodes/nodes_clone_table.go` | Completed |
| **Task-05** | `05-comprehensive-help-text-and-docs.md` | Updated help topics, flag explanations, and usage examples | `cli/cmdnodes/nodes_clone_help.go`, `cli/cmd/nodes_cmd.go` | Completed |
| **Task-06** | `06-specs-rca-and-plan-consolidation.md` | Canonical spec, RCA registration, and plan consolidation | `02-spec/21-app/191-...`, `02-spec/22-app-issues/54-...` | Completed |

---

## 3. Detailed Implementation Outcomes

### 3.1 Task-01: Windows SSH Node Shell Runner & Self-Healing Fallback
- **Heuristic Windows Detection (`isWindowsNode`):** Evaluates `conn.OS`, `conn.OSGroup`, and checks `conn.Username == "Administrator"`. If true, automatically sets `shell = "ps"`.
- **Dynamic Fallback (`isBashMissingError` & `retryWithPowerShell`):** If an SSH execution over `bash` yields `'bash' is not recognized` or `bash: command not found`, the runner automatically retries using PowerShell (`powershell -NoProfile -Command ...`). Upon success, it updates `SSHConnection.OS = 'windows'` in the database via `db.UpdateConnectionOS`, permanently correcting the record.

### 3.2 Task-02: Target Directory Parameter Support and Pre-Navigation
- **Work Directory Pre-Navigation:** Removed the bypass for non-manifest clones so all remote commands navigate to the work directory before cloning:
  - Windows: `Set-Location "<workDir>"; gitmap <kind> <args>` (defaulting to `D:\work`).
  - Unix: `cd <workDir> && gitmap <kind> <args>` (defaulting to `~/work`).
- **Target Directory Support:** Added `TargetDir` extraction from `PassArgs` so passing `gitmap nodes clone <url> <custom_dir>` executes the clone directly inside `<custom_dir>`.

### 3.3 Task-03: `except-self` Flag, Subcommand Token, and Local Skip Routing
- Added support for `--except-self`, `--exceptself`, `--no-self`, `--without-self`, and `--skip-local`.
- Added support for positional token `except-self`: `gitmap nodes clone except-self <url>` and alias `gitmap nodes clone-except-self <url>`.
- Local execution is cleanly bypassed (`opts.IsSkipLocal = true`) while dispatching across all remote fleet nodes.

### 3.4 Task-04: Fleet Clone Terminal UI/UX Overhaul & Error Showcase
- **Error Sanitization:** Stripped `\r\n` and embedded newlines from `r.Error`, extracting the core error message and enforcing a clean single-line table row.
- **Stdout Sanitization:** Filtered out banner divider bars (`===...`), extracting the meaningful final status line.
- **Local Skip Badge:** When `opts.IsSkipLocal` is enabled, the local row is displayed with `○ skipped` and detail `skipped local execution (except-self)` rather than false success.

### 3.5 Task-05: Comprehensive Help Text and Command Documentation
- Updated `nodes_clone_help.go` with `[target_path]` syntax, `--except-self` flag docs, work directory explanations, and examples.
- Updated `cli/cmd/nodes_cmd.go` to include `[dest]` and `except-self` examples.

### 3.6 Task-06: Canonical Specs, RCA and Plan Consolidation
- Canonical specification authored in `02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md` and registered in `02-spec/21-app/readme.md`.
- RCA authored in `02-spec/22-app-issues/54-windows-node-bash-missing-clone-failure-rca.md` and registered in `02-spec/22-app-issues/01-index.md`.
- Consolidated plan authored in `.ai-memory/plans/completed/55-nodes-clone-except-self-windows-runner-and-path.md`.
