# Subtask 18: Modernize Pending Commits Table and Tree Remediation with Tokyo Night

## Status: Pending
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdpending/pending_commits_render.go`

---

## 1. Objective

Modernize the pending commits table banner, header, and footer in `cli/cmdpending/pending_commits_render.go` by replacing hardcoded ASCII / sharp box-drawing characters with `termout.BorderRounded` and `termout.TokyoBoxStyle`. Enhance the remediation option subtree branches with Tokyo Night color styling.

---

## 2. Scope & Implementation Details

1. **Table Banner & Header Modernization**:
   - In `renderTableBanner`:
     - Replace `┌───` with curved corner characters (`╭───`, `╰───`) or `termout.TokyoHeader("Gitmap Pending Commits Summary")`.
     - Align with standard width using `termout.VisualStringWidth`.
   - In `renderTableSummary`:
     - Format metric pills using `termout.TokyoMetricPill` for `Scanned`, `Dirty`, `Uncommitted`, and `Unpushed`.

2. **Remediation Subtree Branch Styling**:
   - In `renderRemediationTitleBranch` and `renderRemediationChildBranch`:
     - Style tree branch glyphs (`├──`, `└──`, `│`) using `termout.SafeTokyoStyle(termout.TokyoMuted)`.
     - Highlight option titles with `termout.SafeTokyoStyle(termout.TokyoPrimary)`.
     - Highlight suggested commands with `termout.SafeTokyoStyle(termout.TokyoAccent)`.

3. **Table Footer Modernization**:
   - In `renderTableFooter`:
     - Replace sharp bottom corners `└────...────┘` with rounded corner box drawing `╰────...────╯`.
     - Format the remediation guidance tip using subtle Tokyo colors.

4. **Decomposition Strategy**:
   - Decompose styling functions into small helpers strictly $\le 15$ lines.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Ifs**: Invert conditions to eliminate nested conditional branching.
- **Whitespace Hygiene**: Blank line after closing brace `}` and blank line before every `return`.
- **Zero Builds & Tests**: Strict prohibition on running `go test` or `go build`.

---

## 4. Verification Checklist

- [ ] Banner and table borders use rounded box characters matching `termout.BorderRounded`.
- [ ] Summary metrics formatted with `termout.TokyoMetricPill`.
- [ ] Tree remediation branches (`├──`, `└──`) styled with Tokyo Night tokens.
- [ ] Table footer uses rounded bottom framing.
- [ ] All functions decomposed to $\le 15$ lines with zero nested `if` statements.
