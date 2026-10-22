# Subtask 12: Implement Visual-Aligned Key-Value Pair Output

## Status: Pending
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/termout/render.go`

---

## 1. Objective

Implement `FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string` and `PrintKeyValuePair(key, val string, keyWidth int, theme Theme)` in `cli/termout/render.go`. These functions provide precision column alignment using printable terminal visual widths (`VisualStringWidth`), stripping ANSI codes and correctly handling East Asian or emoji double-width characters without column misalignment.

---

## 2. Scope & Implementation Details

1. **`FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string`**:
   - Signature: `FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string`
   - Label Assembly:
     - Combine key with colon: `keyLabel := key + ":"`
     - Calculate visual width using `VisualStringWidth(keyLabel)`.
     - Compute padding: `padLen := keyWidth - VisualStringWidth(keyLabel)`. If `padLen < 1`, set `padLen = 1`.
     - Construct padded label: `paddedKey := keyLabel + SpacePool(padLen)`.
   - Monochrome / Non-TTY Handling:
     - If `!colorAllowed()`, return `fmt.Sprintf("  %s %s", paddedKey, val)`.
   - Color / Theme Styling:
     - Resolve key color: use `theme.HeaderColor` if non-empty, otherwise fallback to `RoleSGRBold(RolePrimary)`.
     - Return `fmt.Sprintf("  %s%s%s %s", keyColor, paddedKey, Reset(), val)`.

2. **`PrintKeyValuePair(key, val string, keyWidth int, theme Theme)`**:
   - Signature: `PrintKeyValuePair(key, val string, keyWidth int, theme Theme)`
   - Implementation:
     - Calls `FormatKeyValuePair(key, val, keyWidth, theme)`.
     - Outputs line directly to `output.UI()` via `fmt.Fprintln(output.UI(), ...)`.

3. **Decomposition Strategy**:
   - Helper `resolveKeyColor(theme Theme) string` to isolate header color resolution.
   - Helper `buildPaddedKeyLabel(key string, keyWidth int) string` to handle visual width padding arithmetic.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Functions must be $\le 15$ lines of executable code.
- **Zero Nested Ifs**: Invert all checks into guard clauses.
- **Visual String Width**: Must use `VisualStringWidth` instead of byte length `len()` to avoid Unicode misalignment.
- **Return Styling**: Blank line after closing brace `}` and blank line before every `return`.
- **Zero Builds & Tests**: No `go test` or `go build` executions allowed.

---

## 4. Verification Checklist

- [ ] `FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string` exported in `cli/termout/render.go`.
- [ ] `PrintKeyValuePair(key, val string, keyWidth int, theme Theme)` exported in `cli/termout/render.go`.
- [ ] Visual column padding accurately handles multi-byte and wide runes.
- [ ] Non-TTY output safely generates unstyled padded key-value output.
- [ ] Output routed to `output.UI()` for clean test redirection and output capture.
- [ ] All functions $\le 15$ lines without nested conditional branching.
