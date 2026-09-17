# Plan 189 (Consolidated): Scheduled Shutdown & Restart, Remote GitMap SSH Installation, OS AI Cache Cleaning, SSH Table Formatting, and Help Parity Triad

## Execution Context & Lifecycle
- **Task Start**: Initiated to implement flexible duration parsing for `gitmap schedule shutdown <duration>` and `gitmap schedule restart <duration>` (`1:45hr`, `1:45h`, `2h`, `120m`, `1day`, `1d`, `1s`, `now`), mockable OS power execution isolated from unit tests, remote GitMap installation via SSH (`gitmap cluster install gitmap [target]` and `gitmap sj install gitmap [target]`), transparent preflight GitMap detection and auto-bootstrapping in `cluster exec`, OS AI cache cleanup (`gitmap os ai-clean` and `scripts/os-ai-clean.py`) with preflight table and interactive confirmation, fixed-width ASCII table alignment for `sj ls` and `cluster nodes`, Cluster Triad conceptual help architecture across `cluster.md`, `servers-clients.md`, `ssh-join.md`, and a new Coding Guideline prompt `24-isolate-destructive-os-and-heavy-unit-tests.md`.
- **Workflow**: Continuous N-Step Self-Loop across 2-Agent Concurrency and Parallel Micro-Batches.
- **Total Steps/Loops**: 5 subtasks executed in parallel micro-batches.
- **Status**: 100% Completed & Verified.

---

## Consolidated Subtasks Ledger

### Subtask 01: Schedule Shutdown & Restart Duration Parsing and Mockable OS Power Executor
- **Files Created & Modified**:
  - `cli/cmdschedule/duration_parser.go` (new): Implemented `ParseScheduleDuration(raw string) DurationResult` supporting colon notation (`1:45hr`, `1:45h`, `01:30:00`), day notation (`1day`, `1d`), minute notation (`120m`), second notation (`1s`, `2s`), and instant triggers (`0`, `0s`, `now`) returning `result.Result[time.Duration]`.
  - `cli/cmdschedule/duration_parser_test.go` (new): Hermetic unit tests verifying all duration formats and invalid inputs.
  - `cli/cmdschedule/schedule_os.go` (new): Defined `OSActionType` (`"shutdown"`, `"restart"`), `OSActionParams`, and injectable `DefaultOSActionExecutor` with fallback native commands (`shutdown /s /t` on Windows, `shutdown -h now` on Unix). Implemented `RunSchedulePowerCLI`.
  - `cli/cmdschedule/schedule_os_test.go` (new): Unit tests with mock executors verifying parameter capture and zero host power execution.
  - `cli/cmdschedule/schedule_cmd.go` (modify): Routed `restart` and `shutdown` subcommands to `RunSchedulePowerCLI`.

### Subtask 02: Cluster & SSH-Join GitMap Installation and Automatic Bootstrapping
- **Files Created & Modified**:
  - `cli/cmdssh/cluster_install_cmd.go` (new): Implemented `RunClusterInstallCLI(args []string) error` with target resolution, flags (`--version`, `--sudo`, `--parallel`), and `BuildGitmapInstallOneLiner(osType, version)` supporting official PowerShell (`install.ps1`) and curl/bash (`install.sh`) one-liners.
  - `cli/cmdssh/cluster_install_cmd_test.go` (new): Unit tests verifying argument parsing, target dispatch, and one-liner generation across platforms.
  - `cli/cmdssh/cluster_runner.go` (modify): Added `RequiresGitmap(cmd)` and automated remote preflight check (`command -v gitmap`) with auto-bootstrap execution prior to running commands in `ExecuteNodeCommand`.
  - `cli/cmdssh/cluster_runner_test.go` (new): Unit tests verifying GitMap command detection heuristics.
  - `cli/cmd/cluster.go` (modify): Routed `install` in `routeClusterCore` to `cmdssh.RunClusterInstallCLI`.
  - `cli/cmdssh/sshjoin_cmd.go` (modify): Routed `install` in `isSJSubcommand` and `routeSJSubcommands` to `RunClusterInstallCLI`.

### Subtask 03: OS AI-Clean Command and Standalone Python Utility
- **Files Created & Modified**:
  - `cli/cmdos/os_ai_clean_targets.go` (new): Implemented scanners for Antigravity Brain Caches, System Generated Tasks, OS Temp AI Dumps, and GitMap Installer Caches, excluding protected state files.
  - `cli/cmdos/os_ai_clean.go` (new): Implemented `RunOSAICleanCLI` with ASCII preflight table, file counts, humanized sizes, interactive `[y/N]` confirmation, `--yes` / `-y` bypass, `--dry-run`, and `--json`.
  - `cli/cmdos/os_ai_clean_test.go` (new): Hermetic unit tests verifying directory scanning, size calculation, and protected file exclusion.
  - `cli/cmdos/os.go` (modify): Routed `ai-clean`, `aiclean`, `clean-ai` in `dispatchOSSubcommand`.
  - `cli/cmd/rootutility.go` (modify): Registered `ai-clean` as top-level command alias `gitmap ai-clean`.
  - `scripts/os-ai-clean.py` (new): Standalone, zero-dependency Python 3 utility mirroring the Go scanner and purger.

### Subtask 04: SSH Table Formatting, Column Alignment & Cluster Triad Help Parity
- **Files Created & Modified**:
  - `cli/cmdssh/sshjoin_table.go` (new): Implemented `RenderSSHHostsTable(out io.Writer, hosts []store.SSHHost) error` with fixed-width column formatting (`ALIAS`, `ROLE`, `HOST (IP:PORT)`, `USER`, `STATUS`, `ENROLLED`), 95-character divider, default fallbacks (`"-"`, `"worker"`, `"root"`, `22`), and node counter footer.
  - `cli/cmdssh/sshjoin_table_test.go` (new): Comprehensive unit tests verifying empty table hints, single-node alignment, multi-node formatting, and non-standard ports.
  - `cli/cmdssh/sshjoin_ls_cmd.go` (modify): Replaced unpadded tabwriter with `RenderSSHHostsTable(out, limitedHosts)`.
  - `cli/cmd/cluster_ops.go` (modify): Replaced tabwriter in `renderClusterHostsASCII` with `cmdssh.RenderSSHHostsTable(out, hosts)`.
  - `cli/helptext/catalog.go` (modify): Registered `servers-clients`, `sc`, `clients`, and `cluster-exec` in `topicSummaries`.
  - `cli/helptext/cluster.md` (modify): Added GitMap Cluster Triad ASCII architecture diagram, component roles, and copy-pasteable examples for remote shell execution, GitMap install, and scheduled power commands.
  - `cli/helptext/servers-clients.md` (modify): Added Cluster Triad documentation, examples, and cross-references.
  - `cli/helptext/ssh-join.md` (modify): Added Cluster Triad documentation, examples, and cross-references.

### Subtask 05: Coding Guideline Prompt for Destructive OS Unit Test Isolation
- **Files Created & Modified**:
  - `01-prompts/15-cg-execute/24-isolate-destructive-os-and-heavy-unit-tests.md` (new): Authored complete Coding Guideline 24 mandating mockable OS action executors, prohibiting real destructive OS commands in unit tests, and establishing fast duration testing contracts.
  - `d:/work/coding-guidelines/01-prompts/15-cg-execute/24-isolate-destructive-os-and-heavy-unit-tests.md` (new): Mirrored prompt to external coding-guidelines workspace.
  - `01-prompts/15-cg-execute/01-index.md` (modify): Registered Prompt 24 in the master index table.

---

## Verification & Cleanliness Checks
- All functions strictly adhere to $\le 15$ lines (target $\le 8$ lines).
- Affirmative booleans only (`is*`, `has*`).
- Universal AppError return wrapping.
- Strict Unix LF line endings across all files.
- Hermetic unit tests with mock executors preventing host side effects.
