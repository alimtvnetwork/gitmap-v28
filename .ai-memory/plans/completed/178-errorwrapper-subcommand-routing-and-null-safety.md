# Plan 178: Subcommand Routing ErrorWrapper Architecture & Universal AppError Returns

## 1. Overview & Problem Statement
In previous iterations, helper sub-routers (such as `routeEnvVariableSub`, `routeClusterCore`, `dispatchGroupCRUD`, `routeGitProfileSub`) were refactored to return `result.ErrorWrapper`. However, their parent routers (such as `routeEnvSub`, `dispatchGroup`, `routeProfileSub`, `dispatchSSH`, `routeClusterNodeCommand`, `routeClusterK8sCommand`, `handlePipelineDB`) still declared a return type of standard library `error`. As a consequence, these parent routers were forced to unpack the `result.ErrorWrapper` using `.AsError()` or `.AppError()`:

```go
// ❌ ANTI-PATTERN: Intermediate router returns error, unpacking ErrorWrapper via .AsError()
func routeEnvSub(sub string, args []string) error {
    resVar := routeEnvVariableSub(sub, args)
    if resVar.IsMatched() {
        return resVar.AsError() // or resVar.AppError()
    }
    if sub == constants.CmdEnvPathAdd {
        return routeEnvPath(args)
    }
    return apperror.NewSimple(constants.ErrEnvSubcommand, "E9000")
}
```

This anti-pattern has several critical deficiencies:
1. **Dual-Handling & Premature Conversion**: It forces intermediate routing layers to convert structured `*apperror.AppError` and match state into standard `error` prematurely, rather than preserving the strongly typed `result.ErrorWrapper` throughout the dispatch pipeline.
2. **Go Typed-Nil Interface Trap**: In Go, returning `*apperror.AppError` directly as a standard library `error` interface causes `err != nil` to evaluate to `true` even when the pointer is `nil`. Unpacking via `.AsError()` was a workaround for this signature flaw, when the true solution is for the router to return `result.ErrorWrapper` directly.
3. **Loss of Inspection Predicates**: Callers cannot inspect `.IsSuccess()`, `.IsFailed()`, `.IsMatched()`, or `.IsInvalid()` across routing boundaries.

This plan systematically refactors all parent routers, sub-dispatchers, and subcommands across `cli/cmd/`, `cli/cmdssh/`, `cli/cmdpipeline/`, `cli/cmddb/`, `cli/cmdvscode/`, `cli/cmdzsh/`, and `cli/cmdmacro/` to return `result.ErrorWrapper`, converting to standard `error` via `.AsError()` ONLY at the final CLI/Cobra entry point boundary.

---

## 2. Violation Ledger

| Package / File | Function Name | Current Flawed Signature | Canonical Refactored Signature | Rationale |
|---|---|---|---|---|
| `cli/cmd/env.go` | `routeEnvSub` | `func routeEnvSub(...) error` | `func routeEnvSub(...) result.ErrorWrapper` | Parent router must return `result.ErrorWrapper` and forward `resVar` without `.AsError()`. |
| `cli/cmd/env.go` | `routeEnvPath` | `func routeEnvPath(...) error` | `func routeEnvPath(...) result.ErrorWrapper` | Propagate `result.ErrorWrapper` from `dispatchEnvPath`. |
| `cli/cmd/env.go` | `dispatchEnvPath` | `func dispatchEnvPath(...) error` | `func dispatchEnvPath(...) result.ErrorWrapper` | Return `result.MatchWrapper` / `result.FailureWrapper`. |
| `cli/cmd/group.go` | `dispatchGroup` | `func dispatchGroup(...) error` | `func dispatchGroup(...) result.ErrorWrapper` | Propagate `resCRUD` and `resScoped` directly without `.AsError()`. |
| `cli/cmd/group.go` | `showActiveGroup` | `func showActiveGroup() *apperror.AppError` | `func showActiveGroup() result.ErrorWrapper` | Eliminate typed nil hazard when calling from `runGroup`. |
| `cli/cmd/profile.go` | `routeProfileSub` | `func routeProfileSub(...) error` | `func routeProfileSub(...) result.ErrorWrapper` | Propagate `resGit`, `resDB`, `resChrome`, `resInstall` directly without `.AsError()`. |
| `cli/cmd/cluster.go` | `routeCluster` | N/A (inline in `runCluster`) | `func routeCluster(...) result.ErrorWrapper` | Encapsulate cluster routing cleanly, returning `result.ErrorWrapper`. |
| `cli/cmdssh/cluster_node_cmd.go` | `RouteClusterNodeCLI` | `RunClusterNodeCLI(...) error` | `func RouteClusterNodeCLI(...) result.ErrorWrapper` | Export `RouteClusterNodeCLI` returning `result.ErrorWrapper`; `RunClusterNodeCLI` wraps `.AsError()`. |
| `cli/cmdssh/cluster_k8s_cmd.go` | `RouteClusterK8sCLI` | `RunClusterK8sCLI(...) error` | `func RouteClusterK8sCLI(...) result.ErrorWrapper` | Export `RouteClusterK8sCLI` returning `result.ErrorWrapper`; `RunClusterK8sCLI` wraps `.AsError()`. |
| `cli/cmdssh/ssh.go` | `dispatchSSH` | `func dispatchSSH(...) error` | `func dispatchSSH(...) result.ErrorWrapper` | Propagate `resPrimary` directly without `.AsError()`. |
| `cli/cmdpipeline/pipeline_db_dispatch.go` | `routePipelineDB` | `handlePipelineDB(...) error` | `func routePipelineDB(...) result.ErrorWrapper` | Intermediate router returns `result.ErrorWrapper`; `handlePipelineDB` wraps `.AsError()`. |
| `cli/cmd/bookmark.go` | `routeBookmarkSub` | `func routeBookmarkSub(...) *apperror.AppError` | `func routeBookmarkSub(...) result.ErrorWrapper` | Wrap sub-actions into `MatchWrapperAppErr` / `MatchWrapper`. |
| `cli/cmd/profiles_cmd.go` | `routeProfilesSub` | `func routeProfilesSub(...) error` | `func routeProfilesSub(...) result.ErrorWrapper` | Wrap subcommand execution in `MatchWrapper`. |
| `cli/cmd/startup_cmd.go` | `routeStartupSubcommand` | `func routeStartupSubcommand(...) error` | `func routeStartupSubcommand(...) result.ErrorWrapper` | Return `result.ErrorWrapper` for startup actions. |
| `cli/cmd/storage_cmd.go` | `routeStorageSubcommand` | `func routeStorageSubcommand(...) error` | `func routeStorageSubcommand(...) result.ErrorWrapper` | Wrap drive report & clean in `MatchWrapper`. |
| `cli/cmd/cd.go` | `routeCDSub` | `func routeCDSub(...) error` | `func routeCDSub(...) result.ErrorWrapper` | Wrap cd sub-actions in `MatchWrapper`. |
| `cli/cmddb/cmddb_dispatch.go` | `routeDbSubcommand` | `func routeDbSubcommand(...) error` | `func routeDbSubcommand(...) result.ErrorWrapper` | Wrap db actions in `MatchWrapper`. |
| `cli/cmdvscode/vscode_cmd.go` | `routeVSCodeProjectAction` | `func routeVSCodeProjectAction(...) error` | `func routeVSCodeProjectAction(...) result.ErrorWrapper` | Return `result.ErrorWrapper` for vscode projects. |
| `cli/cmdvscode/vscode_cmd.go` | `routeVSCodeMaintenanceAction` | `func routeVSCodeMaintenanceAction(...) error` | `func routeVSCodeMaintenanceAction(...) result.ErrorWrapper` | Return `result.ErrorWrapper` for vscode maintenance. |
| `cli/cmdzsh/zsh_cmd.go` | `routeZshSubcommand` | `func routeZshSubcommand(...) error` | `func routeZshSubcommand(...) result.ErrorWrapper` | Wrap zsh dispatch in `MatchWrapper`. |
| `cli/cmdmacro/macro_cmd.go` | `routeMacroSubcommand` | `func routeMacroSubcommand(...) error` | `func routeMacroSubcommand(...) result.ErrorWrapper` | Return `result.ErrorWrapper` for macro routing. |
| `cli/cmdssh/types.go` | Domain Types & Aliases | Scattered across implementation files | Centralized in `types.go` | Centralize `SSHTarget`, `SSHJoinOptions`, `BootstrapResult`, `ClusterRunResult`. |

---

## 3. Disjoint Subtasks Decomposition

- **Subtask 01**: `01-cmd-env-group-profile-bookmark-errorwrapper.md`
  - Targets: `cli/cmd/env.go`, `cli/cmd/group.go`, `cli/cmd/profile.go`, `cli/cmd/bookmark.go`
  - Eliminates `resVar.AsError()`, `resCRUD.AsError()`, `resScoped.AsError()`, `resGit.AsError()`, etc.
  - Refactors `routeEnvSub`, `dispatchGroup`, `routeProfileSub`, `routeBookmarkSub` to return `result.ErrorWrapper`.

- **Subtask 02**: `02-cmd-cluster-ssh-k8s-node-errorwrapper.md`
  - Targets: `cli/cmd/cluster.go`, `cli/cmdssh/ssh.go`, `cli/cmdssh/cluster_node_cmd.go`, `cli/cmdssh/cluster_k8s_cmd.go`, `cli/cmdssh/types.go`
  - Creates `cli/cmdssh/types.go` with centralized types.
  - Exports `RouteClusterNodeCLI` and `RouteClusterK8sCLI` returning `result.ErrorWrapper`.
  - Refactors `dispatchSSH` and `routeCluster` to return `result.ErrorWrapper`.

- **Subtask 03**: `03-cmd-auxiliary-routers-errorwrapper.md`
  - Targets: `cli/cmdpipeline/pipeline_db_dispatch.go`, `cli/cmd/profiles_cmd.go`, `cli/cmd/startup_cmd.go`, `cli/cmd/storage_cmd.go`, `cli/cmd/cd.go`, `cli/cmddb/cmddb_dispatch.go`, `cli/cmdvscode/vscode_cmd.go`, `cli/cmdzsh/zsh_cmd.go`, `cli/cmdmacro/macro_cmd.go`
  - Refactors all remaining auxiliary command routers to return `result.ErrorWrapper`.
  - Ensures clean delegation and eliminates all intermediate `.AsError()` conversions.

---

## 4. Quality Gates & Constraints
1. Functions strictly <= 15 lines (target <= 8 lines).
2. Affirmative booleans only (`is*`, `has*`).
3. Zero nested ifs (nesting depth <= 1).
4. Strict Unix LF line endings.
5. TOTAL BAN on running `go test`, `go build`, or runner scripts (`06-cicd-local-runner.py`) during routine execution turns.
6. NO per-file committing: all files accumulated in working tree, committed together in a single grouped atomic commit at the final step before pushing.
