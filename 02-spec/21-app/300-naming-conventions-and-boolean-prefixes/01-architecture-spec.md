# Plan 300 Architecture Specification: Naming Conventions & Positive Boolean Standardization

## User Request (Verbatim)
```text
@[\1](file;.agents/skills/\1)
```

## 1. Executive Summary & Architectural Intent
This specification governs repository-wide compliance with Gitmap Strict Coding Guidelines:
- **Positive Boolean Rule:** All booleans MUST begin with `is` or `has` ONLY (`can`, `should`, `was`, `will` are banned). NO negative booleans (`hasNone`, `isNotFound`, `hasNotExist`, `isExists`).
- **Implicit Boolean Rule:** Eliminate explicit boolean comparisons (`== true`, `=== true`, `== false`, `=== false`).
- **Acronym & Identifier Hygiene:** Replace awkward or compound negative logic with clean, readable predicates (`isDefined`, `isMissing`, `isEmpty`).

## 2. Scope Boundaries
- **In-Scope Subsystems & Files:**
  - Cluster 1 (Explicit & Negative Booleans):
    - `cli/cmddb/rest_visualizer_js1.go`
    - `cli/cmdui/ui_assets_js_fleet.go`
    - `cli/cmdagent/agent_ui_tasks.go`
    - `cli/cmdssh/sshjoin_rm_cmd.go`
    - `cli/cmdagy/agy_running_prompts_copy.go`
    - `cli/cmdstatus/status_resolve.go`
    - `cli/cmdautomation/changed_files_cmd.go`
  - Cluster 2 (Banned Prefixes: `can`, `should`, `was` in AI, Nodes, Install, Scan, VSCode):
    - `cli/cmdai/ai_remove_revert.go`
    - `cli/cmdai/ai_train.go`
    - `cli/cmdnodes/nodes_probe.go`
    - `cli/cmdnodes/nodes_collect.go`
    - `cli/cmdchromeprofile/chromeprofile_zip_export.go`
    - `cli/cmdinstall/install_muse.go`
    - `cli/cmdscan/scan.go`
    - `cli/cmdvscode/vscodepmsync_flags.go`
  - Cluster 3 (Banned Prefixes: `can`, `should`, `was` in OSClean, Release, Special Repos, Sends):
    - `cli/osclean/terminal_clean.go`
    - `cli/release/gitstderrfilter.go`
    - `cli/cmd/special_repos_cli.go`
    - `cli/cmd/special_repos_ops.go`
    - `cli/cmd/cmd_sends.go`
- **Excluded / Deferred:**
  - `go build` and `go test` (strictly banned under Rule 1).
  - Version bumping or changelog edits (strictly banned under Rule 10).
  - Raw git commands inside subagents (strictly banned under Rule 7).

## 3. Core Technical Invariants
1. **Positive Boolean Naming:** All booleans start with `is` or `has`. Zero `can`, `should`, `was`, `will`, `hasNo*`, `isNot*`.
2. **Implicit Comparison:** Replace `val === true` with `Boolean(val)` or `val`. Replace `isExists == false` with `!isPresent`.
3. **AST & Linter Compliance:** Zero violations reported by `03-ai-scripts/08-naming-autofixer.py` and `python 03-ai-scripts/50-fastgate.py`.

## 4. Multi-Agent Partitioning & Ownership Matrix
- **Wave 1 (Tree Swarm Concurrency $A = 3$, $H = 2$):**
  - Worker 01 (`self`): Owns Cluster 1 (Explicit & Negative Booleans).
  - Worker 02 (`self`): Owns Cluster 2 (Banned Prefixes in AI, Nodes, Install, Scan, VSCode).
  - Worker 03 (`self`): Owns Cluster 3 (Banned Prefixes in OSClean, Release, Special Repos, Sends).
- **Wave 2 (Verification & Push):**
  - Root Orchestrator: Validates AST via FastGate and Guideline Autofixer, verifies secrets gate, updates registry, and executes atomic GitMap commit.
