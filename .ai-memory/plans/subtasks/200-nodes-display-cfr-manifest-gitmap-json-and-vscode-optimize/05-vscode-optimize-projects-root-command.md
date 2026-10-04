# Subtask 05: VS Code Optimize Projects Root Command

> **Parent Plan:** `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md`
> **Status:** `PENDING`
> **Target Files:**
> - `cli/vscodepm/optimize.go`
> - `cli/cmdvscode/vscode_optimize.go`
> - `cli/cmd/rootcore.go`
> - `cli/cmd/rootutility.go`
> - `cli/cmdvscode/vscode_cmd.go`

---

## Technical Specification

1. **Root Command Registration**:
   - Register `gitmap vscode optimize-projects` and aliases:
     - `gitmap vsc optimize-projects`
     - `gitmap vpm optimize`
     - `gitmap vscode-optimize`
     - `gitmap vsc-optimize`
     - Direct routing in `cli/cmd/rootcore.go` / `cli/cmd/rootutility.go`.

2. **Optimization Logic in `vscodepm`**:
   - Parse `projects.json` from standard paths:
     - Windows: `%APPDATA%\Code\User\globalStorage\alefragnani.project-manager\projects.json`
     - Linux: `~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json`
     - macOS: `~/Library/Application Support/Code/User/globalStorage/alefragnani.project-manager/projects.json`
   - Detect:
     a) Exact duplicates (identical normalized `rootPath`).
     b) Multiple projects sharing the same base repo name located in different folders (duplicate checkouts).
     c) Missing project directories (where `rootPath` does not exist on disk).
   - Provide Advice:
     - For duplicate copies: print recommended relocation/consolidation target directory (e.g. `Advise: consolidate duplicate copy 'D:\other\repo' into canonical location './repo'`).
   - Clean & Prune:
     - Prune missing project paths from `projects.json`.
     - Prune duplicate entries.
     - Atomically save updated `projects.json`.
   - Output summary table with clear status icons and counts.
