# Specification 140: SSH Fleet Update, Node Inspection, Smart Command Routing, and Macro Engine

## Status: Draft / Active
## Spec ID: 140-ssh-fleet-update-nodes-macros
## Traceability: Task-01, Task-02, Task-03, Task-04, Task-05

---

## 1. User Request (Verbatim)

```text
add commands like

gitmap update ssh all-nodes
gitmap update ssh all-nodes --except id,ip,alias
gitmap ssh update  all-nodes --except id,ip,alias
gitmap ssh nodes -v # will show gitmap version on thoese node if not present then show not installed
gitmap ssh nodes gitmap -v # will show gitmap version on thoese node if not present then show not installed
gitmap ssh nodes ls # run command in the nodes with gitmap ssh ls, understood , test it with e2e please

gitmap ssh macro add/rm/ls/edit/run "macro name" # create first in local and then after verification deploy using json to reimported by other gitmap , do a local e2e testing that it works, clear????

# [V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)
...
all functions do e2e local tests in other vms

gitmap ssh exec "any commands try check if gitmap command then execute in that machine using gitmap if not then try to execute as it is"

Understoodf?? clear??

Implement please
```

---

## 2. Overview & Architecture

This specification defines 4 major feature sets extending GitMap's SSH Fleet Management Subsystem:

1. **SSH Node Inspection (`gitmap ssh nodes ...`)**:
   - `gitmap ssh nodes -v`: Runs a probe across active nodes checking if `gitmap` is available. If present, outputs the installed version (e.g. `v6.315.0`). If not found, prints `not installed`.
   - `gitmap ssh nodes gitmap -v`: Canonical alias for `nodes -v`.
   - `gitmap ssh nodes ls`: Runs listing across nodes, or executes `gitmap ssh ls` format displaying active node statuses, IPs, and connectivity.

2. **Fleet-Wide Binary Update (`gitmap update ssh all-nodes ...`)**:
   - `gitmap update ssh all-nodes [--except id,ip,alias]`: Connects to active SSH fleet nodes, checks remote OS/architecture, copies the appropriate GitMap binary to the remote target (e.g. `C:\Program Files\GitMap\gitmap.exe` or `/usr/local/bin/gitmap`), and verifies update success.
   - `gitmap ssh update all-nodes [--except id,ip,alias]`: Canonical alias allowing command discovery directly under `gitmap ssh`.
   - `--except`: Comma-separated list of machine IDs, IPs, or aliases to skip during fleet update.

3. **Intelligent Command Delegation in `gitmap ssh exec "<cmd>"`**:
   - When a user runs `gitmap ssh exec "<cmd>"`, the runner analyzes `<cmd>`.
   - If `<cmd>` matches a recognized GitMap command/subcommand (e.g. `status`, `scan`, `repo`, `branch`, `pipeline`), and `gitmap` is installed on that node, it delegates the command via `gitmap <cmd>`.
   - If not a GitMap command, or if `gitmap` is not installed on that node, it executes the command as-is using the node's native shell (PowerShell/CMD on Windows, `/bin/sh` or `/bin/bash` on Linux/macOS).

4. **SSH Fleet Macro Engine (`gitmap ssh macro ...`)**:
   - `gitmap ssh macro add <name> [flags]`: Records or defines a macro locally with target commands.
   - `gitmap ssh macro rm <name>`: Removes a defined macro.
   - `gitmap ssh macro ls`: Lists all defined macros with metadata, target command sequences, and update timestamps.
   - `gitmap ssh macro edit <name> [flags]`: Modifies an existing macro.
   - `gitmap ssh macro run <name> [flags]`: Executes the macro across fleet nodes or locally.
   - `gitmap ssh macro deploy <name> [flags]`: Exports the macro as JSON and installs/imports it onto remote fleet nodes for remote gitmap reimport.

---

## 3. Data Contracts & Models

### A. Macro Definition (`cli/cmdssh/ssh_macro_types.go`)

```go
type SSHMacro struct {
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Commands    []string `json:"commands"`
    TargetNodes []string `json:"targetNodes,omitempty"`
    CreatedAt   string   `json:"createdAt"`
    UpdatedAt   string   `json:"updatedAt"`
    Author      string   `json:"author,omitempty"`
}

type SSHMacroCatalog struct {
    Version string     `json:"version"`
    Macros  []SSHMacro `json:"macros"`
}
```

### B. Node Version Status (`cli/cmdssh/ssh_node_types.go`)

```go
type NodeVersionStatus struct {
    NodeID      string `json:"nodeId"`
    Alias       string `json:"alias"`
    Host        string `json:"host"`
    IsInstalled bool   `json:"isInstalled"`
    Version     string `json:"version"`
    RawOutput   string `json:"rawOutput,omitempty"`
    ErrorMsg    string `json:"errorMsg,omitempty"`
}
```

---

## 4. Verification & Quality Gates

1. **Zero Nested Ifs:** All branches flattened with maximum cyclomatic conditional depth = 1.
2. **Positive Booleans:** Boolean flags and variables prefixed only with `is` or `has` (e.g. `isInstalled`, `hasExceptFlag`).
3. **AppError Wrapping:** All failures returned with structured AppError wrappers (`appfault.New` or domain error wrappers).
4. **Targeted Linters:** Must pass `check-nested-ifs.py`, `check-enum-and-boolean.py`, and `check-error-management.py`.
5. **E2E VM Testing:** Local tests against available VMs (w1, w2, w3) and simulated node runners.
