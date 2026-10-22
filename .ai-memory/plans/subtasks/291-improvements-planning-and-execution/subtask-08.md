# Subtask 08: Strongly-Typed Git Constants & Command Enums

> **Subtask ID:** Subtask-08  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-08.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate untyped magic string literals across Git execution wrappers by expanding `cli/constants/constants_git.go`.
2. Define a strongly-typed enum `GitCommandType` ending in the mandatory `Type` suffix.
3. Introduce typed constants for common Git subcommands and flags:
   - `GitRevList = "rev-list"` / `GitCmdRevList GitCommandType = "rev-list"`
   - `GitCountFlag = "--count"` / `GitFlagCount GitCommandType = "--count"`
   - `GitSymbolicRef = "symbolic-ref"` / `GitCmdSymbolicRef GitCommandType = "symbolic-ref"`
   - Core subcommands: `clone`, `pull`, `rebase`, `status`, `fetch`.
4. Preserve backward compatibility for existing callers while migrating command invocations to strongly-typed constants.

---

## 2. Detailed Technical Plan

### 2.1 Define `GitCommandType` Enum in `cli/constants/constants_git.go`
- Add the strongly-typed enum definition:
  ```go
  // GitCommandType represents strongly-typed Git subcommands and execution flags.
  type GitCommandType string

  const (
      GitCmdClone       GitCommandType = "clone"
      GitCmdPull        GitCommandType = "pull"
      GitCmdRebase      GitCommandType = "rebase"
      GitCmdStatus      GitCommandType = "status"
      GitCmdRevList     GitCommandType = "rev-list"
      GitCmdRevParse    GitCommandType = "rev-parse"
      GitCmdSymbolicRef GitCommandType = "symbolic-ref"
      GitCmdForEachRef  GitCommandType = "for-each-ref"
      GitCmdCatFile     GitCommandType = "cat-file"
      GitCmdLsRemote    GitCommandType = "ls-remote"
  )
  ```

### 2.2 Add Matching Command & Flag Constants
- Define string-compatible constants for Git flags and commands:
  ```go
  const (
      GitRevList       = "rev-list"
      GitCountFlag     = "--count"
      GitSymbolicRef   = "symbolic-ref"
      GitFlagPorcelain = "--porcelain"
      GitFlagAbbrevRef = "--abbrev-ref"
  )
  ```

### 2.3 Guideline Alignment
- Verify all new type identifiers end with `Type` (`GitCommandType`, `GitRefCommandType`, `GitRemoteNameType`).
- Enforce blank line spacing between type definitions and const blocks.
- Ensure no functions or constants violate the 500-line file cap ratchet for `constants_git.go`.

---

## 3. Verification Criteria

- [ ] `GitCommandType` is defined as a strong type alias for `string` in `cli/constants/constants_git.go`.
- [ ] `GitRevList`, `GitCountFlag`, and `GitSymbolicRef` are defined as constants.
- [ ] Type names end with `Type` suffix per repository coding standards.
- [ ] All constants are grouped logically with descriptive documentation comments.
- [ ] `cli/constants/constants_git.go` remains under the 500-line file limit.
