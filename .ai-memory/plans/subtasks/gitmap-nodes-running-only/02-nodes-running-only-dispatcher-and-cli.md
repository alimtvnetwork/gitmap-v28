# Subtask 02 — Running-Only Dispatcher, CLI Command & Presentation

## File Box
- `cli/cmdnodes/nodes_running_only.go` (new file; create it)
- `cli/cmdnodes/nodes_cmd.go` (modify)
- `cli/cmd/rootcore.go` (modify)

Do NOT touch any other file outside this file box during the execution of this subtask.

---

## Context & Objectives
Operators require a single, high-speed CLI command to inspect reachable fleet machines without being delayed by offline machine timeouts or distracted by offline warning banners.

Subtask 02 introduces `cli/cmdnodes/nodes_running_only.go` implementing `RunNodesRunningOnlyCLI`, wires the top-level commands in `cli/cmd/rootcore.go` (`nodes-running-only`, `nro`, `nodes-running`, `nodes-online`, `nodes-open`, `nodes-active`), registers subcommands in `cli/cmdnodes/nodes_cmd.go`, and renders dedicated ANSI tables and clean JSON streams containing exclusively online nodes.

---

## Detailed Component Specifications

### 1. `cli/cmdnodes/nodes_running_only.go` (New File)
Implement the running-only execution pipeline:

- **Entry Point Signature:**
  ```go
  package cmdnodes

  import (
  	"context"
  	"fmt"
  	"io"
  	"os"
  	"strings"

  	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
  	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
  )

  // RunNodesRunningOnlyCLI handles `gitmap nodes-running-only` and `gitmap nro`.
  func RunNodesRunningOnlyCLI(args []string) error
  ```

- **Delegation Interception:**
  Check if the first positional token represents a delegation command before running list operations:
  - If `isNodesRunningOnlyDelegation(args)` is true:
    - Forward to `dispatchRunningOnlyDelegation(args)` (handles `exec`, `run`, `tasks`, `resources`, `ping`, `pull`, `cp`).
    - Handled in detail in Subtask 03.

- **Option Parsing & Filtering:**
  - Parse CLI flags:
    - `-j`, `--json`: Output raw JSON array.
    - `-r`, `--refresh`: Force refresh cache.
    - `--no-cache`: Bypass liveness cache.
    - `-t`, `--target <alias|ip>`: Target filter.
    - Subsystem filters: `--ssh`, `--cluster`, `--sc`.
  - Collect nodes via `collectUnifiedNodes(ctx, false, canUseCache)`.
  - Filter to strictly running nodes:
    ```go
    func filterOnlineUnifiedNodes(nodes []UnifiedFleetNode) []UnifiedFleetNode {
    	var online []UnifiedFleetNode
    	for _, n := range nodes {
    		if !n.IsOnline {
    			continue
    		}
    		online = append(online, n)
    	}
    	return online
    }
    ```
  - Apply target and subsystem filters via `applyNodesFilter`.

- **Output Formatting:**
  - **JSON Mode:** Call `printNodesJSON(onlineNodes)`. Guaranteed to omit all offline nodes.
  - **Table Mode:** Call `renderRunningOnlyNodesTable(os.Stdout, onlineNodes, totalCollectedCount)`:
    - If `len(onlineNodes) == 0`:
      ```text
      ● No running nodes currently detected across fleet.
        (All registered nodes are currently offline or unreachable)
      ```
    - If online nodes exist:
      - Render matrix and primary hosts table containing only online nodes.
      - Render custom running-only footer:
        ```text
        ● Running Nodes: 6 active node(s) (4 offline node(s) omitted)
        ```
      - Never print offline warning banners or offline host lists.

### 2. `cli/cmdnodes/nodes_cmd.go` (Modification)
Integrate running-only subcommand detection into `RunUnifiedNodesCLI`:

- **Subcommand Recognition:**
  ```go
  func isNodesRunningOnlyRequest(args []string) bool {
  	if len(args) == 0 {
  		return false
  	}
  	tok := strings.ToLower(args[0])
  	return tok == "running-only" || tok == "running" || tok == "online" ||
  		tok == "open" || tok == "ro" || tok == "nro" || tok == "active"
  }
  ```
- **Dispatch Hook in `RunUnifiedNodesCLI`:**
  ```go
  if isNodesRunningOnlyRequest(args) {
  	return RunNodesRunningOnlyCLI(args[1:])
  }
  ```
- **Help Text Update (`printUnifiedNodesHelp`):**
  Add usage lines and documentation:
  ```text
    gitmap nodes running-only [flags]   List only currently online/active fleet nodes (alias: nro)
    gitmap nro [flags]                  Fast inspection of running nodes with 10s liveness cache
  ```

### 3. `cli/cmd/rootcore.go` (Modification)
Register top-level root command entries in `coreClusterEntries()`:

- Add dispatch mapping:
  ```go
  {
  	[]string{
  		"nodes-running-only",
  		"nro",
  		"nodes-running",
  		"nodes-online",
  		"nodes-open",
  		"nodes-active",
  	},
  	func() error { return cmdnodes.RunNodesRunningOnlyCLI(argsTail()) },
  },
  ```
- Adhere strictly to Go formatting standards: single blank lines between entry declarations, no orphaned imports.

---

## Implementation Checklist
- [ ] Create `cli/cmdnodes/nodes_running_only.go` with `RunNodesRunningOnlyCLI`.
- [ ] Implement `filterOnlineUnifiedNodes` using positive boolean logic (`!n.IsOnline`).
- [ ] Implement `renderRunningOnlyNodesTable` with active node count and omitted offline count.
- [ ] Implement raw JSON array serialization containing strictly `IsOnline == true` nodes.
- [ ] Modify `cli/cmdnodes/nodes_cmd.go` to add `isNodesRunningOnlyRequest` and route to `RunNodesRunningOnlyCLI`.
- [ ] Update `printUnifiedNodesHelp()` with `nodes running-only` and `nro` documentation.
- [ ] Modify `cli/cmd/rootcore.go` to register root-level aliases (`nodes-running-only`, `nro`, `nodes-running`, `nodes-online`, `nodes-open`, `nodes-active`).
- [ ] Verify zero absolute paths, relative git paths only, and positive boolean names.
