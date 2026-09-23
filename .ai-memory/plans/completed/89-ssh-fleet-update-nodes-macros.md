# Plan 89: SSH Fleet Update, Node Inspection, Smart Command Routing, and Macro Engine

Spec Reference: [02-spec/21-app/140-ssh-fleet-update-nodes-macros.md](../../../02-spec/21-app/140-ssh-fleet-update-nodes-macros.md)

## Architectural Context & Blast Radius
GitMap provides distributed SSH orchestration across local clusters and VMs (e.g. `w1`..`w4`, `u1`). Previously, inspecting remote node installations required manual per-machine probes, binary upgrades across the fleet required individual execution, command delegation didn't distinguish between GitMap commands and arbitrary shell executables, and macros were limited to local interactive recording without remote deployment.

This plan introduced end-to-end fleet inspection, automated binary distribution with node exclusions, intelligent command routing with auto-installation fallback, and a full declarative macro management and remote deployment engine.

## Deliverables & Outcomes
- **Spec 140 Authoring & Planning:**
  - Authored canonical Spec 140 in `02-spec/21-app/140-ssh-fleet-update-nodes-macros.md`.
  - Registered Spec 140 in `02-spec/21-app/01-index.md`.
  - Decomposed execution into 5 structured subtasks under `.ai-memory/plans/subtasks/89-ssh-fleet-update-nodes-macros/`.
- **SSH Node Inspection (`gitmap ssh nodes -v` / `gitmap ssh nodes ls`):**
  - Implemented `runSSHNodesVersion` in `cli/cmdssh/ssh_node_version.go` with multi-node concurrent probing.
  - Implemented cross-platform remote version resolution (`gitmap --version` with PATH inspection and binary fallback).
  - Supported `gitmap ssh nodes -v`, `gitmap ssh nodes gitmap -v`, and `gitmap ssh nodes ls` table and JSON views (`--json`).
  - Correctly outputs installed version or `not installed` alongside node reachability indicators.
- **Fleet Binary Update (`gitmap update ssh all-nodes` / `gitmap ssh update all-nodes`):**
  - Implemented `parseUpdateOptions` supporting `all-nodes` / `all` targeting and `--except <id,ip,alias>` filtering.
  - Supported token absorption for comma-separated or space-separated exclusion lists across shell environments.
  - Wired `gitmap update ssh ...` in `cli/cmd/rootutility.go` to forward to `cmdssh.RunSSHUpdateCLI`.
  - Verified remote update downloading, deployment, and live version verification on Windows and Linux nodes.
- **Smart Command Routing (`gitmap ssh exec "<cmd>"`):**
  - Enhanced `cli/cmdssh/sshexec.go` to analyze command tokens and detect GitMap subcommands.
  - For GitMap commands, verifies remote GitMap presence; if missing, triggers SSH auto-installation; if installation is unavailable, cleanly falls back to executing the raw command in the remote shell.
  - For standard shell commands, executes directly as-is without overhead.
- **Fleet Macro Engine (`gitmap ssh macro add/rm/ls/edit/run/deploy`):**
  - Implemented `cli/cmdssh/ssh_macro_ops.go` with local CRUD verification and remote JSON payload transfer.
  - Supported `add`, `rm`, `ls`, `edit`, `run`, and `deploy` verbs in `cli/cmdssh/ssh_macro_cmd.go`.
  - Fixed Windows PowerShell `$true`/`$false` literal encoding and eliminated double-wrapping in `cli/cmdssh/ssh_macro_sync.go`.
  - Enables local authoring and verification, followed by deployment to remote nodes for re-import and execution by remote GitMap instances.
- **Verification & Quality Gates:**
  - Added comprehensive unit tests in `cli/cmdssh/ssh_nodes_macro_update_test.go` covering version request parsing, output filtering, and flag processing.
  - Verified 0 nested if violations via `python linter-scripts/check-nested-ifs.py`.
  - Verified boolean and enum compliance via `python linter-scripts/check-enum-and-boolean.py`.
  - Verified error management via `python linter-scripts/check-error-management.py`.
  - Formatted Go code via `python 03-ai-scripts/26-go-code-formatter.py`.
  - Executed live end-to-end verification against real VM cluster nodes (`w1`, `w2`, `w3`).
