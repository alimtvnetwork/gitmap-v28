# Subtask 08: Git Inspection Commands & Flag Type Enum Expansion

> **Subtask ID:** Subtask-08  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-08.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Modernize `cli/constants/constants_git.go` by introducing strongly-typed enum types for Git inspection subcommands and standard flags.
2. Define `GitInspectCommandType` with members `GitCmdCherry`, `GitCmdDiff`, and `GitCmdMergeBase` to eliminate magic string literals across git history, diff, and divergence inspection routines.
3. Define `GitFlagType` with member `GitFlagStat` (and associated string constants) for diff summary flags.
4. Maintain backward compatibility with untyped string constants (`GitCherry`, `GitDiff`, `GitMergeBase`, `GitStatFlag`).
5. Enforce strict coding guidelines:
   - Type naming: All enum types must end with the `Type` suffix.
   - Constant formatting: Strict PascalCase formatting.
   - File size ratchet: Keep `cli/constants/constants_git.go` comfortably below the 500-line ceiling.

---

## 2. Detailed Technical Plan

### 2.1 Add `GitInspectCommandType` Enum to `cli/constants/constants_git.go`
- Define the new type alias and associated constants:
  ```go
  // GitInspectCommandType represents Git subcommands inspecting histories, diffs, and relationships.
  type GitInspectCommandType string

  const (
      GitCmdCherry    GitInspectCommandType = "cherry"
      GitCmdDiff      GitInspectCommandType = "diff"
      GitCmdMergeBase GitInspectCommandType = "merge-base"
  )
  ```
- Define untyped string constants for general CLI formatting:
  ```go
  const (
      GitCherry    = "cherry"
      GitDiff      = "diff"
      GitMergeBase = "merge-base"
  )
  ```

### 2.2 Add `GitFlagType` Enum to `cli/constants/constants_git.go`
- Define the new type alias and associated flag constants:
  ```go
  // GitFlagType represents command-line flags supplied to Git subprocess executions.
  type GitFlagType string

  const (
      GitFlagStat GitFlagType = "--stat"
  )

  const (
      GitStatFlag = "--stat"
  )
  ```

### 2.3 File Sizing & Hygiene Verification
- Verify that `cli/constants/constants_git.go` retains valid Go package declarations and comments.
- Verify file length remains well under 300 lines (current length 255 lines), well within the 500-line cap ratchet.
- Verify positive naming conventions across any helper predicates or constants.

---

## 3. Verification Criteria

- [ ] `GitInspectCommandType` is declared with `GitCmdCherry`, `GitCmdDiff`, and `GitCmdMergeBase`.
- [ ] `GitFlagType` is declared with `GitFlagStat`.
- [ ] String constants `GitCherry`, `GitDiff`, `GitMergeBase`, and `GitStatFlag` are present.
- [ ] All new types end with `Type` suffix in full compliance with coding guidelines.
- [ ] File length of `cli/constants/constants_git.go` remains strictly under 500 lines.
- [ ] Strictly relative Git paths used throughout all references.
