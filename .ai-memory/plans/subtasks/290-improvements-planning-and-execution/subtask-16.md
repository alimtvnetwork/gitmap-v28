# Subtask 16: Purge Dracula Colors & Bind Tokyo Night Tokens in CG Audit

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdcg/cg_audit.go`

---

## 1. Objective

Purge all hardcoded Dracula hex color literals (`#8be9fd`, `#50fa7b`, `#ffb86c`, `#ff5555`, `#6272a4`) from `cli/cmdcg/cg_audit.go` and bind audit diagnostics to centralized `termout.SafeTokyoStyle` tokens and Tokyo status pill badges.

---

## 2. Scope & Implementation Details

1. **Dracula Hex Code Purge**:
   - Replace legacy style variables (lines 13–19):
     ```go
     // REMOVE
     cgAuditTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
     cgPassStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
     cgWarnStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
     cgFailStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
     cgMutedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))

     // REPLACE WITH TOKYO TOKENS
     cgAuditTitleStyle = termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true)
     cgPassStyle       = termout.SafeTokyoStyle(termout.TokyoOK).Bold(true)
     cgWarnStyle       = termout.SafeTokyoStyle(termout.TokyoWarn).Bold(true)
     cgFailStyle       = termout.SafeTokyoStyle(termout.TokyoError).Bold(true)
     cgMutedStyle      = termout.SafeTokyoStyle(termout.TokyoMuted)
     ```

2. **Audit Check Rendering Refactor**:
   - In `renderAuditCheckItem` or check renderers, use `termout.TokyoOKPill("PASS")`, `termout.TokyoWarnPill("WARN")`, and `termout.TokyoErrorPill("FAIL")`.
   - In score rendering, use `termout.RenderTokyoProgress` for compliance score visualization.

3. **TTY & Safe Color Fallback**:
   - Because `termout.SafeTokyoStyle` automatically verifies `colorAllowed()`, color output degrades safely in piped stdout without ANSI escape pollution.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines.
- **Zero Nested Ifs**: Flatten all branches with guard clauses and early returns.
- **Boolean Naming**: Ensure all booleans use `is` or `has` prefixes (`isPass`, `hasDirectory`).
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] Zero Dracula hex color strings (`#8be9fd`, etc.) remain in `cli/cmdcg/cg_audit.go`.
- [x] Variables bound to `termout.SafeTokyoStyle(...)`.
- [x] Status formatting leverages `termout.TokyoOKPill` / `TokyoWarnPill` / `TokyoErrorPill`.
- [x] Functions conform to $\le 15$ lines.
