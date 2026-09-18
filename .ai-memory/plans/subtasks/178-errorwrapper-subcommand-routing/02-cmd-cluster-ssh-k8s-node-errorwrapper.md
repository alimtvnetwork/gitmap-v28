# Subtask 02: Cluster, SSH, and Node Command Routers ErrorWrapper Refactor

## 1. Objectives
Refactor cluster and SSH router/dispatch signatures across `cli/cmd/cluster.go`, `cli/cmdssh/ssh.go`, `cli/cmdssh/cluster_node_cmd.go`, and `cli/cmdssh/cluster_k8s_cmd.go` to return `result.ErrorWrapper`, and centralize domain types in `cli/cmdssh/types.go`.

## 2. Disjoint Target Files
- `cli/cmd/cluster.go`
- `cli/cmdssh/ssh.go`
- `cli/cmdssh/cluster_node_cmd.go`
- `cli/cmdssh/cluster_k8s_cmd.go`
- `cli/cmdssh/types.go`

## 3. Implementation Details

### A. `cli/cmdssh/types.go`
Centralize domain models, options, and Result aliases for `cmdssh`:
- `SSHTarget`
- `SSHJoinOptions`
- `BootstrapResult`
- `ClusterRunResult`
- `SSHHealthResult`
- Canonical type aliases:
  `type ClusterRunResultSlice = result.ResultSlice[ClusterRunResult]`
  `type BootstrapResultSlice = result.ResultSlice[BootstrapResult]`

### B. `cli/cmdssh/cluster_node_cmd.go`
1. Export `RouteClusterNodeCLI(args []string) result.ErrorWrapper`:
   - If help requested, returns `result.MatchWrapper(showClusterNodeHelp())`.
   - Dispatches `routeClusterNodeCommand(ctx, args[0], args[1:])`. If matched, returns `res` directly (NO `.AsError()`).
   - Default returns `result.FailureWrapper(apperror.NewValidationError(fmt.Sprintf("unknown cluster node subcommand: %s", args[0])))`.
2. `RunClusterNodeCLI(args []string) error`:
   - Returns `RouteClusterNodeCLI(args).AsError()`.

### C. `cli/cmdssh/cluster_k8s_cmd.go`
1. Export `RouteClusterK8sCLI(args []string) result.ErrorWrapper`:
   - If help requested, returns `result.MatchWrapper(showClusterK8sHelp())`.
   - Dispatches `routeClusterK8sCommand(ctx, args[0], args[1:])`. If matched, returns `res` directly (NO `.AsError()`).
   - Default returns `result.FailureWrapper(apperror.NewValidationError(fmt.Sprintf("unknown cluster k8s subcommand: %s", args[0])))`.
2. `RunClusterK8sCLI(args []string) error`:
   - Returns `RouteClusterK8sCLI(args).AsError()`.

### D. `cli/cmdssh/ssh.go`
1. `dispatchSSH(ctx context.Context, args []string, parent *cobra.Command) result.ErrorWrapper`:
   - If empty args, returns `result.MatchWrapper(handleEmptySSHArgs())`.
   - Calls `resPrimary := dispatchPrimarySSH(ctx, sub, args[1:], parent)`. If matched, returns `resPrimary` directly.
   - If fallback matched, returns `result.SuccessWrapper()`.
   - Returns `result.MatchWrapper(runSSHLogin(parent, args, ctx))`.
2. `runSSH(args []string) error`:
   - Calls `dispatchSSH(...)`.
   - Returns `res.AsError()`.

### E. `cli/cmd/cluster.go`
1. `routeCluster(args []string) result.ErrorWrapper`:
   - If root help, prints help and returns `result.SuccessWrapper()`.
   - Dispatches `resHelp := dispatchInvertedClusterHelp(args)`. If matched, returns `resHelp` directly.
   - Dispatches `resSub := dispatchClusterSubcommand(args[0], args[1:])`. If matched, returns `resSub` directly.
   - Default returns `result.FailureWrapper(apperror.NewSimple("unknown command", "E9000"))`.
2. In `routeClusterLegacyOps`:
   - `case "node": return cmdssh.RouteClusterNodeCLI(rest)` directly.
3. In `routeClusterK8s`:
   - `case "k8s", "kube", "kubernetes": return cmdssh.RouteClusterK8sCLI(rest)` directly.
4. `runCluster(args []string) error`:
   - Returns `routeCluster(args).AsError()`.

## 4. Verification
- Functions strictly <= 15 lines.
- Affirmative booleans only.
- Strict Unix LF line endings.
- Targeted linters only.
