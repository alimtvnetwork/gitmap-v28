# Subtask Plan 09: Add Strongly Typed Constants Ending in `Type` (`cli/constants/constants_git.go`)

- **Subtask ID:** `subtask-09`
- **Subtask Name:** `09-constants-git-typed-enums`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/constants/constants_git.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/constants/constants_git.go`, Git subcommands, ref commands, and remote identifiers are declared as untyped string constants or passed as raw string literals.

Under **Rule 5 (Strongly Typed Enums & Constants)**:
1. All enum types MUST end with the `Type` suffix (e.g. `GitRefCommandType`, `GitRemoteNameType`).
2. Raw magic strings across Git operations make refactoring difficult and prevent compile-time type verification.
3. Adding strongly typed constants for `symbolic-ref`, `origin`, and related ref subcommands establishes type safety while maintaining backward-compatible aliases for existing callers.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Add `GitRefCommandType` Enum
- **File:** `cli/constants/constants_git.go`
- **Definition:**
```go
// GitRefCommandType represents Git subcommands operating on references and symbolic pointers.
type GitRefCommandType string

const (
	GitRefCmdSymbolicRef GitRefCommandType = "symbolic-ref"
	GitRefCmdRevParse    GitRefCommandType = "rev-parse"
	GitRefCmdForEachRef  GitRefCommandType = "for-each-ref"
)
```

### 2.2 Add `GitRemoteNameType` Enum
- **File:** `cli/constants/constants_git.go`
- **Definition:**
```go
// GitRemoteNameType represents standard remote identifiers across Git repositories.
type GitRemoteNameType string

const (
	GitRemoteNameOrigin   GitRemoteNameType = "origin"
	GitRemoteNameUpstream GitRemoteNameType = "upstream"
)
```

### 2.3 Preserve Backwards-Compatible Constants
- **File:** `cli/constants/constants_git.go`
- **Requirements:**
  - Ensure untyped constants `GitSymbolicRef = "symbolic-ref"` and `GitOrigin = "origin"` remain accessible without breaking existing usages.
  - Link or alias them to the typed constants where appropriate.

---

## 3. Implementation Verification Checklist

- [x] `GitRefCommandType` declared with `Type` suffix.
- [x] `GitRefCmdSymbolicRef` and `GitRefCmdRevParse` defined under `GitRefCommandType`.
- [x] `GitRemoteNameType` declared with `Type` suffix.
- [x] `GitRemoteNameOrigin` and `GitRemoteNameUpstream` defined under `GitRemoteNameType`.
- [x] Backwards-compatible untyped constants preserved.
- [x] Zero build checks or test runs during spec phase.
- [x] Formatting follows repo vertical whitespace standards.
