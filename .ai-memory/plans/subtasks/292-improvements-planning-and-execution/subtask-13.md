# Subtask 13: Define Table WrapModeType Enum and RowStripe in Termtable Types

## Status: Completed
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/termtable_types.go`

---

## 1. Objective

Define the strongly-typed enum `WrapModeType` (`WrapClip`, `WrapEllipsis`, `WrapWord`) and expand `ColumnConfig` and `TableConfig` to support cell wrapping and alternating row striping (`RowStripe bool`) in `cli/termout/termtable_types.go`.

---

## 2. Scope & Implementation Details

1. **`WrapModeType` Enum Definition**:
   - Declare `type WrapModeType string`.
   - Define constants ending in `Type` / adhering to repository enum naming standards:
     ```go
     const (
         WrapClip     WrapModeType = "clip"
         WrapEllipsis WrapModeType = "ellipsis"
         WrapWord     WrapModeType = "word"
     )
     ```
   - Provide docstrings detailing the behavior of each wrap mode:
     - `WrapClip`: Hard truncate at max width without ellipsis characters.
     - `WrapEllipsis`: Truncate with `"..."` ellipsis (default).
     - `WrapWord`: Soft wrap overflowing text across multiple lines at whitespace breaks.

2. **`ColumnConfig` Enhancement**:
   - Add field `WrapMode WrapModeType` to allow per-column override of table wrapping behavior.

3. **`TableConfig` Enhancement**:
   - Add field `WrapMode WrapModeType` as table-wide default wrapping mode.
   - Add positive boolean field `RowStripe bool` indicating whether alternating row background striping is active.

---

## 3. Strict Coding Guidelines & Invariants

- **Enum Suffix Rule**: Enum types must end in `Type` (`WrapModeType`).
- **Positive Booleans**: Use `RowStripe bool` (or `HasRowStripe bool`). Banned prefixes `can`, `should`, `was`, `will` and negative booleans are forbidden.
- **Backward Compatibility**: Preserve existing fields, constants, and type aliases (`Column`, `AlignType`).
- **Zero Builds & Tests**: Spec and plan authoring strictly avoids compiler runs.

---

## 4. Verification Checklist

- [x] `WrapModeType` type and constants (`WrapClip`, `WrapEllipsis`, `WrapWord`) declared in `cli/termout/termtable_types.go`.
- [x] `WrapMode WrapModeType` field added to `ColumnConfig`.
- [x] `WrapMode WrapModeType` field added to `TableConfig`.
- [x] `RowStripe bool` field added to `TableConfig`.
- [x] Comprehensive Go docstrings added to all newly exported types and fields.
