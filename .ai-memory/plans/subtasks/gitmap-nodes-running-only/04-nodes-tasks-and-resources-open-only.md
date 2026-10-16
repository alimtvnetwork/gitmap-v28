# Subtask 04: Node Tasks & Resources Open-Only Filtering

- **Subtask ID:** `SUBTASK-04`
- **Parent Plan:** `.ai-memory/plans/gitmap-nodes-running-only.md`
- **Component Spec:** `02-spec/21-app/gitmap-nodes-running-only/02-component-spec.md`
- **Status:** `PENDING`
- **Priority:** High
- **Dependencies:** Subtask 01 (`01-fast-fleet-node-liveness-cache-and-probe-optimization.md`)

---

## 1. Objective & Context

When running fleet interrogation commands such as `gitmap nodes tasks` and `gitmap nodes resources`, the CLI currently attempts to connect to all enrolled cluster nodes, even those that are powered off or unreachable. This causes:
1. Significant latency overhead as workers block waiting for TCP connection timeouts (typically 1200ms to 5000ms per offline host).
2. Cluttered output reports containing unreachable hosts with timeout error messages.

This subtask implements strict open-only filtering across both task queue discovery and resource telemetry interrogation:
- Adding `isOpenOnly bool` to `nodesTasksOptions` in `cli/cmdnodes/nodes_tasks_types.go`.
- Expanding CLI argument parsers in both tasks and resources to accept `--running-only`, `-ro`, `--open-only`, and `--online-only`.
- Skipping offline nodes *prior* to launching remote SSH interrogation goroutines in `executeFleetTasksInterrogation` and `executeFleetResourcesProbeUnified`.
- Eliminating connection latency penalty on offline nodes and producing clean, active-node-only telemetry.

---

## 2. Owned Files & Target Modifications

| File Path | Description of Changes |
|---|---|
| `cli/cmdnodes/nodes_tasks_types.go` | Add `isOpenOnly bool` field to `nodesTasksOptions` struct; update `parseNodesTasksOptions` to parse `--running-only`, `-ro`, `--open-only`, `--online-only`. |
| `cli/cmdnodes/nodes_tasks.go` | Update `executeFleetTasksInterrogation` to skip offline nodes when `opts.isOpenOnly` is true; update `printNodesTasksHelp` with flag documentation. |
| `cli/cmdnodes/nodes_resources.go` | Update `parseNodesResourcesOptions` to parse `--running-only` and `-ro`; update `executeFleetResourcesProbeUnified` to skip offline nodes before dispatch; update help text. |

---

## 3. Detailed Technical Design & Code Specifications

### 3.1 Task Queue Options Extension (`cli/cmdnodes/nodes_tasks_types.go`)

Extend `nodesTasksOptions` to carry `isOpenOnly`:

```go
type nodesTasksOptions struct {
	isJSON       bool
	isOpenOnly   bool
	limit        int
	targetFilter string
}
```

Update `parseNodesTasksOptions(args []string)` to recognize running-only flags:

```go
func parseNodesTasksOptions(args []string) nodesTasksOptions {
	opts := nodesTasksOptions{
		limit: 10,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--json" || a == "-j":
			opts.isJSON = true
		case a == "--open-only" || a == "--online-only" || a == "--running-only" || a == "-ro":
			opts.isOpenOnly = true
		case (a == "--limit" || a == "-l") && i+1 < len(args):
			i++
			opts.limit = parsePositiveIntLimit(args[i], opts.limit)
		case isTaskFilterCandidate(a) && opts.targetFilter == "":
			opts.targetFilter = a
		}
	}
	return opts
}
```

### 3.2 Pre-Dispatch Offline Filtering in Tasks (`cli/cmdnodes/nodes_tasks.go`)

In `executeFleetTasksInterrogation`:

```go
	for _, n := range unifiedNodes {
		isLocalNode := strings.EqualFold(n.Host, "127.0.0.1") || strings.EqualFold(n.Host, "localhost") || strings.EqualFold(n.Alias, "local")
		if isLocalNode {
			continue
		}
		if opts.isOpenOnly && !n.IsOnline {
			continue
		}
		if isTargetMatching(n.Alias, n.Host, opts.targetFilter) {
			wg.Add(1)
			sem <- struct{}{}
			go func(node UnifiedFleetNode) {
				defer func() {
					<-sem
					wg.Done()
				}()

				conn := resolveNodeSSHConnection(node, connMap)
				rep := queryRemoteNodeTasks(ctx, node, conn, opts.limit)
				mu.Lock()
				reports = append(reports, rep)
				mu.Unlock()
			}(n)
		}
	}
```

### 3.3 Resource Probing Flag & Pre-Dispatch Filtering (`cli/cmdnodes/nodes_resources.go`)

In `parseNodesResourcesOptions(args []string)`:

```go
		if a == "--open-only" || a == "--online-only" || a == "--running-only" || a == "-ro" {
			opts.isOpenOnly = true
			continue
		}
```

In `executeFleetResourcesProbeUnified(ctx context.Context, nodes []UnifiedFleetNode, opts nodesResourcesOptions)`:

```go
	for _, n := range nodes {
		if opts.isOpenOnly && !n.IsOnline {
			continue
		}
		hasTargetMatch := isNodeTargetIncluded(n, opts.targetFilter)
		if hasTargetMatch {
			wg.Add(1)
			go func(node UnifiedFleetNode) {
				defer wg.Done()
				alias, port, conn := resolveFleetNodeEndpoint(node, connMap)
				isWin := isWindowsNode(conn) || strings.Contains(strings.ToLower(node.Role), "win")
				rec := probeNodeExec(alias, fmt.Sprintf("%s:%d", node.Host, port), isWin, false, opts.benchmarkLimit, &conn)
				mu.Lock()
				records = append(records, rec)
				mu.Unlock()
			}(n)
		}
	}
```

---

## 4. Step-by-Step Implementation Plan

1. **Step 1: Update `cli/cmdnodes/nodes_tasks_types.go`**:
   - Add `isOpenOnly bool` to `nodesTasksOptions`.
   - Update `parseNodesTasksOptions` switch statement to handle `--running-only`, `-ro`, `--open-only`, `--online-only`.
2. **Step 2: Update `cli/cmdnodes/nodes_tasks.go`**:
   - In `executeFleetTasksInterrogation`, insert the pre-check `if opts.isOpenOnly && !n.IsOnline { continue }`.
   - Update `printNodesTasksHelp` to list `--running-only, -ro, --open-only: Query tasks exclusively on online nodes`.
3. **Step 3: Update `cli/cmdnodes/nodes_resources.go`**:
   - In `parseNodesResourcesOptions`, add `--running-only` and `-ro` to the condition setting `opts.isOpenOnly = true`.
   - In `executeFleetResourcesProbeUnified`, skip nodes where `opts.isOpenOnly && !n.IsOnline` before incrementing `wg.Add(1)` and spawning the probing goroutine.
   - Update `printNodesResourcesHelp` to include `--running-only` and `-ro`.

---

## 5. Invariants & Verification Checklist

- [ ] **Positive Booleans**: Use `opts.isOpenOnly` and `n.IsOnline` (no negative boolean naming or `== true`).
- [ ] **Zero Nested Ifs**: Nesting depth $\le 1$ strictly maintained; use guard clauses (`continue`).
- [ ] **Zero Probe Penalty**: Offline nodes are never dispatched to SSH worker goroutines when running in open-only mode.
- [ ] **Relative Paths**: All references use relative forward-slash paths (`cli/cmdnodes/...`).
- [ ] **Zero Build & Zero Test Rules**: Author specification strictly without intermediate build/test runs.
