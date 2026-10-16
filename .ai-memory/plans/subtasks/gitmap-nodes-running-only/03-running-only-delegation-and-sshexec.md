# Subtask 03 — Running-Only Delegation, SSH Exec & Fleet Filtering

## File Box
- `cli/cmdssh/sshexec.go` (modify)
- `cli/cmdnodes/nodes_filter.go` (modify)
- `cli/cmdnodes/nodes_scan.go` (modify)

Do NOT touch any other file outside this file box during the execution of this subtask.

---

## Context & Objectives
Executing batch commands across fleet machines using `gitmap se "<cmd>"` or `gitmap nodes scan` currently reports offline nodes with prominent warning blocks and includes offline placeholder stubs in JSON telemetry.

When running through `gitmap nro exec "<cmd>"` or when specifying `--running-only` / `-ro`, operators desire a streamlined execution path:
1. Target solely machines that are online and reachable.
2. Completely suppress offline warning banners (`Notice: The following machine(s) are currently OFF...`).
3. Omit offline nodes from JSON output, emitting only records for nodes where commands actually executed.
4. Accelerate candidate node filtering via liveness caching.

Subtask 03 enhances `cli/cmdssh/sshexec.go`, `cli/cmdnodes/nodes_filter.go`, and `cli/cmdnodes/nodes_scan.go` to support first-class `--running-only` behavior.

---

## Detailed Component Specifications

### 1. `cli/cmdssh/sshexec.go` (Modification)
Integrate running-only flag parsing, output suppression, and JSON filtering:

- **Options Struct Update:**
  ```go
  type seOptions struct {
  	Exclude       string
  	Except        string
  	ExceptOS      string
  	TargetOS      string
  	Target        string
  	IP            string
  	IsJSON        bool
  	IsRunningOnly bool
  	Args          []string
  	IsShowHelp    bool
  }
  ```

- **Flag Parsing (`extractSECustomFlags`):**
  Recognize running-only flags:
  ```go
  if a == "--running-only" || a == "-ro" || a == "--open-only" || a == "--online-only" {
  	opts.IsRunningOnly = true
  	continue
  }
  ```

- **Execution Flow Adaptation (`executeOnAllSSH`):**
  ```go
  func executeOnAllSSH(conns []db.SSHConnection, args []string, isJSON bool, isRunningOnly bool)
  ```
  - **Banner Suppression:**
    - When `isRunningOnly` is true:
      - Skip calling `printOfflineAtStart(offline)`.
      - Only print the active execution banner (`printInjectedBanner(online, cmdStr)`).
    - When `isRunningOnly` is false:
      - Retain standard behavior with full offline notices.
  - **Summary Output (`printExecFinishSummary`):**
    - When `isRunningOnly` is true:
      - Suppress `Summary of offline machines` block.
      - Output only completion message for active nodes:
        `✓ SSH Execution completed across %d active node(s).`

- **JSON Output Stream (`executeOnAllSSHJSON`):**
  ```go
  func executeOnAllSSHJSON(online, offline []db.SSHConnection, args []string, isRunningOnly bool) {
  	var results []NodeExecResult
  	if !isRunningOnly {
  		for _, c := range offline {
  			results = append(results, NodeExecResult{
  				Alias:  c.Alias,
  				IP:     c.IPAddress,
  				Status: "offline",
  			})
  		}
  	}
  	// Execute across online nodes and serialize results...
  }
  ```
  - When `isRunningOnly` is true: zero offline stubs are appended; the JSON array contains exclusively executing online nodes.

- **Help Documentation (`printSSHExecHelp`):**
  Document the flag:
  ```text
    -ro, --running-only     Execute exclusively across online nodes, omitting offline warnings
  ```

### 2. `cli/cmdnodes/nodes_filter.go` (Modification)
Support running-only and open-only filtering with accelerated reachability checks:

- **Filter Options:**
  ```go
  type NodeFilterOptions struct {
  	Target      string   `json:"target,omitempty"`
  	Except      []string `json:"except,omitempty"`
  	Include     []string `json:"include,omitempty"`
  	IncludeMain bool     `json:"include_main,omitempty"`
  	OpenOnly    bool     `json:"open_only,omitempty"`
  	RunningOnly bool     `json:"running_only,omitempty"`
  }
  ```

- **Flag Extraction (`parseSingleNodeFilterFlag`):**
  Add parsing for `--running-only` and `-ro`:
  ```go
  if arg == "--running-only" || arg == "-ro" || arg == "--open-only" {
  	opts.RunningOnly = true
  	opts.OpenOnly = true
  	return 1
  }
  ```

- **Reachability Filtering (`FilterFleetNodes`):**
  - Check if `opts.RunningOnly` or `opts.OpenOnly` is true.
  - Utilize `filterOnlineNodesConcurrently` with timeout bounded to `1200ms`.
  - Ensure positive boolean naming throughout: `isOnline`, `isConnectionCandidate`, `isNodeIncluded`.

### 3. `cli/cmdnodes/nodes_scan.go` (Modification)
Support running-only filtering during fleet scan and rescan:

- **Argument Extraction (`isFilterFlagStep` & `isScanIgnoredArg`):**
  - Update flag scanner to recognize `--running-only` and `-ro`.
  - Pass filter flags cleanly into `ParseNodeFilterOptions(args)`.
- **Target Filtering:**
  - `targets := FilterFleetNodes(conns, opts)` automatically receives `opts.RunningOnly = true`, excluding offline hosts before SSH client dial attempts begin.
  - Prevents scan timeout stalls for offline machines.
- **Help Text Update (`printNodesScanHelp`):**
  Add `--running-only` / `-ro` to flag list:
  ```text
    -ro, --running-only          Execute scan only across reachable/online nodes
  ```

---

## Implementation Checklist
- [ ] In `cli/cmdssh/sshexec.go`, add `IsRunningOnly` to `seOptions` and parse `--running-only` / `-ro`.
- [ ] Suppress start offline warning banners and finish offline lists when `IsRunningOnly` is active.
- [ ] In `cli/cmdssh/sshexec.go`, omit offline node stubs from JSON output when `IsRunningOnly` is active.
- [ ] Update `cli/cmdssh/sshexec.go` help text with `--running-only` examples.
- [ ] In `cli/cmdnodes/nodes_filter.go`, add `RunningOnly` to `NodeFilterOptions` and support `-ro` / `--running-only`.
- [ ] Ensure `FilterFleetNodes` filters out offline candidates when `RunningOnly` is set.
- [ ] In `cli/cmdnodes/nodes_scan.go`, support `--running-only` and `-ro` in pass args and filter extraction.
- [ ] Verify zero absolute paths, strictly relative git paths, and positive boolean names.
- [ ] Verify zero deep nesting in conditional branches.
