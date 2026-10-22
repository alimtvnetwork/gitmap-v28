# Architecture Spec: 304 - Constants & Strongly-Typed Enums

## 1. Executive Summary & Intent
This specification formalizes the remediation of Enum and Constant violations in accordance with `01-prompts/15-cg-execute/06-constants-and-enums.md` and repository guidelines (`02-spec/02-coding-guidelines/`).
All Go enum type identifiers must bear the mandatory `*Type` suffix (e.g. `AppScopeType`, `DownloadEngineType`, `CPUProfileType`, `UpdateStageType`), while providing type aliases for backwards compatibility. Raw numerical rune casts (`rune(0)`) are eliminated in favor of literal character zero (`'\x00'`).

## 2. Invariants & Rules
1. **Mandatory Type Suffix:** All enum definitions must end in `Type`.
2. **Backwards-Compatible Aliases:** Maintain `type Enum = EnumType` so existing caller code compiles seamlessly without disruption.
3. **No Numerical Rune Casts:** Replace `rune(0)` with `'\x00'`.
4. **Zero Nested Ifs:** All conditionals remain flattened with early returns.
5. **File Size Bounds:** All touched files must stay <= 500 lines.
6. **Zero Builds & Zero Tests:** FastGate AST verification and `check-enum-guidelines.py` only.

## 3. Disjoint Micro-Batch Allocation
- **Worker 01 (Apps, Concurrency, Downloader, UserContext):**
  - Files: `cli/cmdapps/types.go`, `cli/cloneconcurrency/pull_concurrency.go`, `cli/cmddownload/types.go`, `cli/usercontext/types.go`.
  - Action: Enforce `AppScopeType`, `PackageManagerType`, `CPUProfileType`, `DownloadEngineType`, `AuthStatusType`.
- **Worker 02 (Release Tags, Folder Tree, Dock, Update UI):**
  - Files: `cli/cmdfixreleasetags/types.go`, `cli/cmdfoldertree/treesearch_match.go`, `cli/cmdos/os_dock_types.go`, `cli/cmdupdate/update_ui_types.go`.
  - Action: Enforce `AuditReasonType`, `ProtectionStatusType`, `StepStatusType`, `SearchModeType`, `DockPositionType`, `UpdateStageType`.
- **Worker 03 (Diag, TermGeneric, TermOut, Heatmap Query):**
  - Files: `cli/diag/types.go`, `cli/termgeneric/header.go`, `cli/termout/themes.go`, `cli/heatmap/heatmap_query.go`.
  - Action: Enforce `SuggestionCategoryType`, `RenderFormatType`, `AlignmentType`, `RoleType`; eliminate `rune(0)` in `heatmap_query.go`.

## 4. Verification Evidence Contract
- FastGate AST inspection: `python 03-ai-scripts/50-fastgate.py` exit 0.
- Enum linter: `python linter-scripts/check-enum-guidelines.py`.
- Forbidden strings check: `python linter-scripts/check-forbidden-strings.py` exit 0.
