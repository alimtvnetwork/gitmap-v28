# Subtask 12: Unified Key-Value Output Engine

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/render.go`

---

## 1. Objective

Implement `PrintKeyValPadded` and `FormatKeyValPadded` in `cli/termout/render.go` to unify terminal key-value output across CLI subcommands, eliminating scattered `fmt.Printf` padding calls and integrating `termout.PadRightMax` with semantic role colors.

---

## 2. Scope & Implementation Details

1. **`FormatKeyValPadded`**:
   - Signature: `FormatKeyValPadded(key string, val any, keyWidth int) string`
   - Key formatting: Appends colon suffix `key + ":"`.
   - Padding: Uses `PadRightMax(key+":", keyWidth)` to achieve uniform alignment backed by the zero-allocation space pool.
   - Styling:
     - When colors are enabled (`colorAllowed()`):
       - Key formatted with `RoleSGRBold(RolePrimary)`.
       - Separator/Colon reset: `Reset()`.
       - Value formatted with `RoleSGR(RoleFG)` or string conversion.
     - When colors are disabled:
       - Renders plain string: `fmt.Sprintf("  %-*s %v", keyWidth, key+":", val)`.
   - Line constraint: $\le 12$ lines.

2. **`PrintKeyValPadded`**:
   - Signature: `PrintKeyValPadded(key string, val any, keyWidth int)`
   - Directly writes `FormatKeyValPadded(key, val, keyWidth)` followed by a newline to `output.UI()`.
   - Line constraint: $\le 5$ lines.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Max 15 lines per function.
- **Zero Nested Ifs**: Guard clauses and early returns only.
- **Theme Consistency**: Strictly utilize `RoleSGR` semantic roles instead of raw ANSI escape sequences.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `FormatKeyValPadded` correctly pads the key with colon using `PadRightMax`.
- [x] `PrintKeyValPadded` writes output directly to `output.UI()`.
- [x] Plain-text fallback works when stdout is not a TTY or `NO_COLOR` is active.
- [x] No functions exceed 15 lines.
