# Subtask 13: Border Style Enum Specification & TableConfig Extension

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/termtable_types.go`

---

## 1. Objective

Define the `BorderStyleType` enum in `cli/termout/termtable_types.go` and extend `TableConfig` to support structured border rendering modes (`BorderRounded`, `BorderSharp`, `BorderMinimal`) while maintaining full backward compatibility.

---

## 2. Scope & Implementation Details

1. **`BorderStyleType` Enum**:
   - Must adhere to the enum naming convention (ending with `Type`):
     ```go
     // BorderStyleType defines the visual framing style for rendered terminal tables.
     type BorderStyleType string

     const (
     	// BorderRounded uses curved UTF-8 box-drawing characters (╭, ╮, ╯, ╰).
     	BorderRounded BorderStyleType = "rounded"
     	// BorderSharp uses classic rectangular box-drawing characters (┌, ┐, ┘, └).
     	BorderSharp BorderStyleType = "sharp"
     	// BorderMinimal uses understated separator lines without outer vertical bars.
     	BorderMinimal BorderStyleType = "minimal"
     )
     ```

2. **`TableConfig` Struct Extension**:
   - Update `TableConfig` with field `BorderStyle BorderStyleType`:
     ```go
     type TableConfig struct {
     	Columns      []Column
     	Rows         []Row
     	HeaderColor  string
     	BorderColor  string
     	HasBorders   bool
     	BorderStyle  BorderStyleType
     	EllipsisText string
     }
     ```
   - Retain existing fields (`Columns`, `Rows`, `HeaderColor`, `BorderColor`, `HasBorders`, `EllipsisText`) without breaking existing call sites.

---

## 3. Strict Coding Guidelines & Invariants

- **Enum Standards**: Type suffix `Type` is mandatory (`BorderStyleType`). Constant names must be PascalCase (`BorderRounded`, `BorderSharp`, `BorderMinimal`).
- **Affirmative Booleans**: `HasBorders` preserved as affirmative boolean.
- **File Sizing**: File length must remain well under the 100-line cap (currently $\approx 35$ lines).
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `BorderStyleType` defined with type alias `string` and `Type` suffix.
- [x] Constants `BorderRounded`, `BorderSharp`, `BorderMinimal` exported.
- [x] `TableConfig` includes `BorderStyle BorderStyleType`.
- [x] Existing `TableConfig` fields preserved without type breaking changes.
