# Plan 191 (Consolidated): SSH Config Sanitize, Clone --force Reclone, Cluster/SC/SSH-Join Help Parity & Architecture

## Execution Context & Lifecycle
- **Task Start**: Initiated to resolve OpenSSH client crash triggered by server daemon directive `authorizedkeysfile` in `~/.ssh/config`, implement `--force` / `-f` / `--reclone` on `gitmap clone`, establish parity for `ssh-join` in root cluster help and route `cluster join`, deeply enrich `servers-clients` (`sc`) help text with command breakdowns, joining mechanics, and install/gitmap examples, and document cluster architecture while providing an automated `gitmap cluster init` JSON configuration generator.
- **Workflow**: Continuous N-Step Self-Loop across 2-Agent Concurrency and Parallel Micro-Batches.
- **Total Steps/Loops**: 5 subtasks executed in parallel micro-batches.
- **Status**: 100% Completed & Verified.

---

## Consolidated Subtasks Ledger

### Subtask 01: SSH Config Sanitizer & Bad Directive Auto-Healing
- **Problem**: `authorizedkeysfile` is an `sshd_config` server daemon directive, invalid in client `~/.ssh/config`. When present, OpenSSH client crashes with `Bad configuration option: authorizedkeysfile` (exit 255), breaking `git clone` and `ssh` commands.
- **Implementation**:
  - `cli/cmdssh/sshconfig.go` (modified): Implemented `SanitizeSSHConfigContent(content string) (string, int)` and `SanitizeSSHConfigFile(path string) (bool, *apperror.AppError)`. Scans lines for 20+ daemon-only directives (`authorizedkeysfile`, `permitrootlogin`, `subsystem`, `strictmodes`, `clientaliveinterval`, etc.) and safely comments them out with `# [gitmap removed server-only directive]: ...` without breaking valid client directives or host configurations.
  - Auto-healing integrated into `updateSSHConfig` so any SSH config write auto-sanitizes the file, and exposed standalone via `gitmap ssh config --sanitize` (`-s`).
  - `cli/cmd/fixauth.go` (modified): Integrated `cmdssh.SanitizeSSHConfigFile("")` into `gitmap fix-auth`.
  - `cli/cmdssh/sshconfig_test.go` (new): Hermetic unit tests covering `authorizedkeysfile` isolation, multi-directive stripping, and file modification idempotency.

### Subtask 02: Clone Force & Reclone Flag (`--force`, `-f`, `--reclone`)
- **Problem**: `gitmap clone <url>` previously skipped cloning if the local folder existed (`already exists on disk`), without providing a direct force/reclone flag for single URLs.
- **Implementation**:
  - `cli/cmdclone/flags.go` (modified): Registered `--force`, `-f`, and `--reclone` as first-class aliases on `cleanFlag`.
  - `cli/cmdclone/clone.go` (modified): Refactored direct URL cloning to use `DirectCloneParams`. Implemented `cleanDirectCloneTarget(absPath, params.IsClean)` and `maybeCleanMultiFolder` which forcefully removes existing target directories with `os.RemoveAll` before cloning when `--force` / `-f` / `--clean` / `--reclone` is specified.
  - `cli/cmdclone/clonefixrepo.go` (modified): Updated caller to pass `DirectCloneParams`.
  - `cli/helptext/clone.md` (modified): Documented `--force`, `-f`, and `--reclone` in flag reference and examples.
  - `cli/cmdclone/clone_force_test.go` (new): Unit tests for flag parsing and hermetic directory removal.

### Subtask 03: Cluster SSH Join & Root Help Parity
- **Problem**: `ssh-join (sj)` was omitted from the root cluster group help, and `cluster join` routing needed verification against direct host join arguments.
- **Implementation**:
  - `cli/constants/constants_helpgroups.go` (modified): Updated `CompactCluster` to `"  servers-clients (sc), clients, cluster, ssh-join (sj)"`.
  - `cli/cmd/rootusage_groups.go` (modified): Rendered `constants.HelpSSHJoin` under `printGroupCluster()`.
  - `cli/cmd/rootusagefilter_rows.go` (modified): Added `constants.HelpSSHJoin` under `constants.HelpGroupCluster` in `allHelpRows()`.
  - `cli/cmdssh/exports.go` (modified): Confirmed `RunClusterJoinCLI` routes directly to `RunSSHJoinCLI(args)` (`gitmap cluster join <ip> [alias]`).
  - `cli/cmdssh/cluster_join_routing_test.go` (new): Unit tests verifying help group presence and CLI routing validation.

### Subtask 04: Servers-Clients (`sc`) Help Text Deep Enrichment
- **Problem**: `gitmap sc help` lacked command-by-command descriptions, explanation of how host joining functions under the hood, and examples for running GitMap operations and package installations.
- **Implementation**:
  - `cli/helptext/servers-clients.md` (modified): Extensively enriched with detailed explanations of every subcommand (`bash`, `sh`, `ps`, `cmd`, `join`, `nodes`, `rm`, `ping`, `install`, `pull`, `push`, `status`, `proj`, `restart`, `shutdown`), architectural breakdown of machine admission and key distribution ("How Join Works"), copy-pasteable examples for multi-node GitMap commands (`gitmap sc bash "gitmap status"`, `gitmap sc pull --all`), and package installations (`gitmap sc install "git,nodejs,curl"`, `gitmap sc install gitmap`).
  - `cli/helptext/sc.md` (modified): Synchronized short-form documentation with identical architectural parity.

### Subtask 05: Cluster Architecture, JSON Schema Generator & Documentation
- **Problem**: Users lacked clear conceptual explanations of `cluster` vs `exec` vs `execute`, and had no CLI command to scaffold cluster JSON topology files automatically.
- **Implementation**:
  - `cli/cmdssh/cluster_init_cmd.go` (new): Implemented `gitmap cluster init` (with aliases `template`, `init-config`, `schema`) to scaffold starter cluster topology JSON configurations (`01-config.json`) with `--control`, `--workers`, `--user`, `--password`, `--out`, and `--force` flags.
  - `cli/cmd/cluster.go` (modified): Routed `init`, `template`, `init-config`, `schema` to `cmdssh.RunClusterInitCLI`.
  - `cli/cmdssh/cluster_init_cmd_test.go` (new): Hermetic unit tests verifying JSON generation, custom node configurations, and file writing with overwrite protection.
  - `cli/helptext/cluster.md` (modified): Added upfront architectural guide explaining Kubernetes cluster orchestration vs SSH node fleets, documented `cluster init`, and added a disambiguation comparison table ("Local Exec (`exec`)" vs "Macro Replay (`execute`)" vs "Remote Multi-Node Fleet (`cluster exec` / `sc`)").
  - `cli/helptext/exec.md` (modified): Added cross-references directing users to `cluster exec` and `servers-clients` for remote multi-host execution.
  - `cli/helptext/catalog.go` (modified): Registered `cluster-init` and `cluster-template`.

---

## Verification & Cleanliness Checks
- All functions strictly adhere to <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Zero nested ifs (nesting depth <= 1).
- Universal AppError wrapping across all error paths.
- Targeted linters passed with 0 violations (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`).
