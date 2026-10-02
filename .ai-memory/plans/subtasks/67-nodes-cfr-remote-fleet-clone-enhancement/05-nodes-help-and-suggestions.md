# Subtask 05: Nodes Help Menu, Examples & Command Suggestions Alignment

> **Parent Plan:** [67-nodes-cfr-remote-fleet-clone-enhancement](../../pending/67-nodes-cfr-remote-fleet-clone-enhancement.md)
> **Component Spec:** [02-component-spec.md](../../../../02-spec/21-app/67-nodes-cfr-remote-fleet-clone-enhancement/02-component-spec.md)
> **Primary Targets:**
> - `cli/helptext/nodes.md`
> - `cli/cmd/nodes_cmd.go`
> - `cli/cmdnodes/nodes_clone_table.go`
> - `cli/cmdnodes/nodes_clone_help.go`

---

## 1. Objective

Provide comprehensive CLI documentation, interactive help, and actionable terminal footer suggestions for fleet nodes operations:
1. Create canonical embedded help manual `cli/helptext/nodes.md` matching standard help text standards across GitMap (`gitmap help nodes`).
2. Synchronize `printUnifiedNodesHelp` in `cli/cmd/nodes_cmd.go` with full command signatures (`gitmap nodes cfr [flags] [targets] [dest]`), subfolder preservation behavior, accept keys commands (`gitmap ssh trust <alias>`), and machine queries (`gitmap machine`).
3. Update interactive sub-command help screens in `cli/cmdnodes/nodes_clone_help.go`.
4. Implement context-sensitive post-run suggestions footer in `cli/cmdnodes/nodes_clone_table.go` that guides the user after executing `gitmap nodes clone/cfr/cfrp`.

---

## 2. Detailed Technical Scope

### 2.1 Embedded Help Manual (`cli/helptext/nodes.md`)
Author the complete embedded markdown manual covering:
- Command header, overview, and aggregation description (SSH, Cluster DB, Server-Clients network).
- Full syntax table:
  - `gitmap nodes`
  - `gitmap nodes ping [target]`
  - `gitmap nodes clone [flags] <targets> [dest]`
  - `gitmap nodes cfr [flags] <targets> [dest]`
  - `gitmap nodes cfrp [flags] <targets> [dest]`
  - `gitmap nodes history`
- Flag tables with short/long forms, descriptions, and defaults:
  - `-t`, `--target`
  - `--exclude`
  - `--except-self`, `--no-self`, `--skip-local`
  - `-d`, `--dir`, `--dest`, `--target-dir`
  - `--dry-run`
  - `-j`, `--json`
- Work Directory & Relative Subfolder Preservation documentation.
- Concrete real-world examples:
  - Repository by name, full HTTPS/SSH URL, comma-separated targets.
  - Manifest file transfer (`gitmap.json`).
  - Custom target directory overrides.
  - Host key acceptance (`gitmap ssh trust <alias>`).
  - Machine telemetry queries (`gitmap machine`, `gitmap nodes ping`).

### 2.2 Terminal Help Parity (`cli/cmd/nodes_cmd.go`)
Update `printUnifiedNodesHelp()` to ensure the terminal output reflects:
- Correct signatures with optional `[dest]` parameter for `clone`, `cfr`, and `cfrp`.
- Explanation of intelligent CWD subfolder preservation.
- Clear examples demonstrating cloning, host key trust, and machine queries.

### 2.3 Subcommand Help Screens (`cli/cmdnodes/nodes_clone_help.go`)
Enhance `PrintNodesCloneHelp(kind NodesCloneKind)`:
- Add documentation for `-d, --dir, --dest, --target-dir` flags.
- Document automatic remote directory pre-creation.
- Clarify CWD subfolder preservation behavior.

### 2.4 Actionable Terminal Suggestions Footer (`cli/cmdnodes/nodes_clone_table.go`)
Author `renderFleetSuggestionsFooter(out io.Writer, opts NodesCloneOptions, results []RemoteCloneNodeResult)`:
- Hook into `renderFleetResultsTable` immediately after `renderFleetSummaryFooter`.
- Analyze node execution results:
  - When failures or unreachable nodes occur: suggest `gitmap nodes ping`, `gitmap ssh check`, and `gitmap ssh trust <alias>`.
  - When successful: suggest `gitmap nodes`, `gitmap machine`, and `gitmap ssh <alias>`.
- Always suggest `gitmap nodes history` for audit trail tracking.

---

## 3. Implementation Steps

1. **Step 1:** Create `cli/helptext/nodes.md` with complete documentation, flags, and examples.
2. **Step 2:** Update `printUnifiedNodesHelp` in `cli/cmd/nodes_cmd.go` to expose `[dest]`, key trust, and machine telemetry.
3. **Step 3:** Update `printFlagsSection` and `printDescriptionSection` in `cli/cmdnodes/nodes_clone_help.go`.
4. **Step 4:** Implement `renderFleetSuggestionsFooter` in `cli/cmdnodes/nodes_clone_table.go` and invoke it from `renderFleetResultsTable`.
5. **Step 5:** Manually test help commands:
   - `gitmap nodes help`
   - `gitmap help nodes`
   - `gitmap nodes cfr --help`

---

## 4. Verification & Quality Gates

- Verify markdown rendering via `gitmap help nodes`.
- Verify terminal help rendering via `gitmap nodes help` and `gitmap nodes cfr --help`.
- Ensure all Go functions are $\le 15$ lines, use positive boolean identifiers, and wrap errors cleanly.
- Verify zero compiler warnings with `go build ./cli/...`.
