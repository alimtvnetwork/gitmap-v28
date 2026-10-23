# Subtask 01: Commit Info CLI Command, Routing Guard & Root Aliases

## Status: READY
- **Task ID:** `315-subtask-01`
- **Parent Spec:** [02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/01-architecture-spec.md](../../../../02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/01-architecture-spec.md)
- **Module:** `cli/cmdcommit`, `cli/cmd`
- **Owned Files:**
  - `cli/cmdcommit/commit_info.go`
  - `cli/cmdcommit/commit_info_types.go`
  - `cli/cmdcommit/commit_cmd.go`
  - `cli/cmd/rootcore.go`

---

## 1. Objectives

1. **Implement `RunCommitInfo(args []string) error` in `cli/cmdcommit/commit_info.go`:**
   - Provide the primary command entry point for inspecting commit details, universal SHA resolution, cross-database telemetry, and ASCII tree output.
   - Parse command-line flags cleanly: `--json`, `--repo <path|slug>`, `--verbose`, `-v`, `--help`, `-h`.
   - Resolve target repository root from current working directory or explicit `--repo` flag using `resolveCommitInfoRepoRoot(repoFlag string) (string, error)`.
   - Dispatch query execution and rendering based on requested output mode (rich terminal tree vs machine JSON).

2. **Define Shared Domain Types in `cli/cmdcommit/commit_info_types.go`:**
   - Establish typed data models: `CommitInfoReport`, `CommitProvenance`, `CommitDbCrossReferences`, `CommitFileDelta`, and `OriginType`.
   - Standardize query parameter structs: `CommitInfoOptions` and `UniversalShaQuery`.

3. **Install Subcommand Guard in `cli/cmdcommit/commit_cmd.go`:**
   - Detect `info` as the first positional argument to `gitmap commit` or `gitmap cm` (e.g., `gitmap commit info <sha>` or `gitmap cm info <sha>`).
   - Route execution cleanly to `RunCommitInfo(args[1:])` without triggering git commit staging or push actions.
   - Guard against flag permutations and display commit-info specific usage when `gitmap commit info --help` is requested.

4. **Register Top-Level Root Aliases in `cli/cmd/rootcore.go`:**
   - Bind `commit-info` and shorthand `ci` directly to `cmdcommit.RunCommitInfo(argsTail())`.
   - Ensure aliases integrate seamlessly with the existing `coreWorkflowEntries()` routing table.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Create Data Models in `cli/cmdcommit/commit_info_types.go`
- Define `OriginType` string enum: `OriginLocalAI`, `OriginForeignAI`, `OriginHumanOrExternal`.
- Declare `CommitInfoOptions` struct:
  ```go
  type CommitInfoOptions struct {
      Query        string
      RepoPath     string
      IsJsonFormat bool
      IsVerbose    bool
  }
  ```
- Declare `CommitInfoReport`, `CommitProvenance`, `CommitDbCrossReferences`, and `CommitFileDelta` structs with JSON struct tags.

### Step 2: Implement CLI Parsing & Dispatch in `cli/cmdcommit/commit_info.go`
- Implement `RunCommitInfo(args []string) error`:
  - Check for help flags (`--help`, `-h`, `help`). If present, render help and return `nil`.
  - Parse CLI flags using helper `parseCommitInfoFlags(args []string) (CommitInfoOptions, error)`.
  - Validate that a query identifier or SHA argument is provided. If missing, render usage guidance and return validation error.
  - Call `ExecuteCommitInfo(opts CommitInfoOptions) error`.
- Implement `parseCommitInfoFlags(args []string) (CommitInfoOptions, error)`:
  - Iterate tokens and populate `CommitInfoOptions` using discrete helper functions $\le 15$ lines.
  - Support `--json`, `--repo`, `-r`, `--verbose`, `-v`.
- Implement `renderCommitInfoHelp(w io.Writer)`:
  - Render concise usage banner, examples, and options list with ANSI styling.

### Step 3: Implement Routing Guard in `cli/cmdcommit/commit_cmd.go`
- Add check in `runCommit(args []string) error`:
  ```go
  if isCommitInfoSubcommand(args) {
      return RunCommitInfo(args[1:])
  }
  ```
- Implement `isCommitInfoSubcommand(args []string) bool`:
  - Check `len(args) > 0 && strings.EqualFold(args[0], "info")`.
- Verify that standard commit actions (`gitmap commit -m "msg"`, `gitmap commit -p`) remain completely untouched and unregressed.

### Step 4: Register Root Aliases in `cli/cmd/rootcore.go`
- Locate `coreWorkflowEntries()` in `cli/cmd/rootcore.go`.
- Append dispatch entry:
  ```go
  {[]string{"commit-info", "ci"}, func() error { return cmdcommit.RunCommitInfo(argsTail()) }},
  ```
- Guarantee that functions remain strictly $\le 15$ lines and all entries conform to project naming standards.

---

## 3. Strict Coding Guidelines & Constraints

- **15-Line Function Limit:** Every function must be $\le 15$ lines (preferred $\le 8$ lines).
- **Positive Booleans Only:** Use positive boolean prefixes exclusively (`isJsonFormat`, `isVerbose`, `hasRepoFlag`, `isInfoCommand`, `hasQuery`). Banned: `noJson`, `isNotInfo`, `unspecified`.
- **Zero Nested Conditionals:** Guard clauses and early returns only; zero nested `if` statements.
- **Error Wrapping:** Wrap all I/O, parsing, and execution errors with `*apperror.AppError` via `apperror.WrapSimple(err, "context")`.
- **Relative Path Hygiene:** Zero absolute drive letters (`C:`, `D:`) or file URI schemes (`file:///`).
- **No Direct Exits in Leaf Functions:** Return errors directly to caller; never invoke `os.Exit` or exit handlers.
- **Zero Builds / Zero Tests:** Do not execute `go build` or `go test` during authoring.

---

## 4. Verification Evidence & Acceptance Criteria

1. `gitmap commit info <sha>` dispatches to `RunCommitInfo` and does not attempt to create a git commit.
2. `gitmap cm info <sha>` routes identically to `gitmap commit info <sha>`.
3. `gitmap commit-info <sha>` and `gitmap ci <sha>` dispatch directly from root CLI without error.
4. `gitmap commit info --help` displays dedicated commit inspection help menu.
5. `gitmap commit info <sha> --json` sets `isJsonFormat = true` in options.
6. All modified and authored files pass the 15-line limit and positive boolean naming checks.
