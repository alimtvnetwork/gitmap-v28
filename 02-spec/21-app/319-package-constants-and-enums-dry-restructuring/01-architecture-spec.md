# Architecture Spec: 319 - Package Constants & Enums DRY Restructuring

## 1. Executive Summary & Intent
This specification formalizes the architectural blueprint for establishing modular enum types and centralized constants across GitMap's core packages. Following `01-prompts/15-cg-execute/06-constants-and-enums.md` and repository guidelines (`02-spec/02-coding-guidelines/`), all magic strings, untyped literals, magic numbers, and raw rune casts are eradicated.

The architecture introduces:
1. Strongly-typed enum definitions residing in domain-specific `enums/` packages, where every enum type identifier bears the mandatory `*Type` suffix.
2. Centralized `constants/constants.go` files (partitioned into companion files when exceeding 300 lines) organizing flags, operational keys, table names, format templates, and error identifiers.
3. Strict decoupling and acyclic dependency graphs preventing circular imports between root CLI packages, store engines, and command dispatchers.

## 2. Invariants & Coding Standards
1. **Mandatory Type Suffix:** All enum type definitions must end with `Type` (e.g., `PullModeType`, `CloneProtocolType`, `StoreTableType`, `CommitVerbType`).
2. **Standard Enum Interface Methods:**
   - `IsValid() bool`: Validates if the enum value is within defined bounds.
   - `String() string`: Returns canonical serialized string representation.
   - Positive boolean query methods: `IsRebase()`, `IsFastForward()`, `IsShallow()`, `IsSSH()`, etc. (Negative booleans such as `IsNotValid()` or `IsUnspecified()` are strictly forbidden).
3. **No Magic Number / Magic String Literals:**
   - Raw string comparisons (e.g., `if mode == "rebase"`) are replaced with enum equality (`if mode == enums.Rebase`).
   - Magic integers in flags or timeouts are declared as typed constants.
   - Raw character rune casts (`rune(0)`) are replaced with explicit literal character zero (`'\x00'`).
4. **Size and Complexity Boundaries:**
   - Max file size $\le 300$ lines.
   - Max function size $\le 15$ lines (preferred $\le 8$ lines).
   - Zero nested `if` statements (guard clauses and early returns only).
   - Blank line after `}` and blank line before `return`.
5. **Backwards Compatibility Bridge:**
   - Preserve type aliases (`type PullMode = PullModeType`) during transitional phases if required by external packages.

## 3. Package Structure & Directory Layout
To achieve maximum Code DRYness and completely eliminate circular import risks in Go:
```
cli/
  ├── cmdpull/
  │     ├── enums/
  │     │     ├── pull_mode_enum.go       # PullModeType (Rebase, FfOnly, Merge, Auto)
  │     │     └── pull_status_enum.go     # PullStatusType (UpToDate, FastForward, Conflict, Error)
  │     └── constants/
  │           └── constants.go            # Pull CLI flags, summary tokens, progress indicators
  ├── cloner/
  │     ├── enums/
  │     │     ├── clone_protocol_enum.go  # CloneProtocolType (SSH, HTTPS, Git, File)
  │     │     └── clone_depth_enum.go     # CloneDepthType (Full, Shallow, Blobless, Treeless)
  │     └── constants/
  │           └── constants.go            # Git CLI arguments, environment keys, progress stages
  ├── cmdclone/
  │     ├── enums/
  │     │     ├── clone_strategy_enum.go  # CloneStrategyType (Batch, Sequential, Interactive, Parallel)
  │     │     └── clone_profile_enum.go   # CloneProfileType (Default, Work, Personal, Fleet)
  │     └── constants/
  │           └── constants.go            # Flags, interactive prompts, telemetry event tags
  ├── store/
  │     ├── enums/
  │     │     ├── store_table_enum.go     # StoreTableType (Repositories, Pipelines, Clusters, Workgroups)
  │     │     └── db_op_enum.go           # DbOpType (Query, Exec, Begin, Commit, Rollback)
  │     └── constants/
  │           ├── constants.go            # SQLite pragma strings, connection flags, query limits
  │           └── constants_tables.go     # Table name constants and schema migration tokens
  └── cmdcommit/ (Rollout Package 5)
        ├── enums/
        │     ├── commit_verb_enum.go     # CommitVerbType (Feat, Fix, Chore, Refactor, Docs, Test, Release)
        │     └── commit_format_enum.go   # CommitFormatType (StandardBracket, HyphenSeparated, Conventional)
        └── constants/
              └── constants.go            # Commit command flags, prompt templates, validation rules
```

## 4. Circular Dependency Mitigation
1. **Subpackage Isolation:** Placing `enums/` and `constants/` inside each domain package ensures they do not import parent packages.
2. **Upstream Leaf Packaging:** Domain `enums` packages import only standard library packages (`fmt`, `strings`, `encoding/json`).
3. **Cross-Domain Shared Enums:** Any enums shared across 2 or more unrelated CLI commands reside in `cli/enums/<domain>/` so that both child packages can import them safely without mutual dependency.

## 5. Acceptance Criteria
- [ ] Specification approved and structured for 4 primary packages + 5th rollout package.
- [ ] Every candidate enum defines mandatory `*Type` suffix, `IsValid()`, and `String()`.
- [ ] All constants partitioned into dedicated `constants/constants.go` files $\le 300$ lines.
- [ ] Zero circular imports between packages.
- [ ] Subtask plans authored in `.ai-memory/plans/subtasks/319-package-constants-and-enums-dry-restructuring/`.
- [ ] Ready for multi-agent tree swarm execution ($A = 3$, $H = 2$, Capacity = 6) on user prompt trigger.
