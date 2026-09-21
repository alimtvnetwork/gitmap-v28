# RCA: Windows Terminal Glyph Font Rendering Issue with Subarrow in Pull All

- **Date:** 2026-09-21
- **Affected Subcommands:** `gitmap pull`, `gitmap pull-all` (`gitmap pa`), `gitmap push`
- **Target Release:** `v6.294.0`
- **Scope:** UI rendering, glyph font compatibility, Windows console safe glyph fallbacks

---

## 1. Symptom

When running `gitmap pull all` (or `gitmap pull-all`, `gitmap pull`) in Windows PowerShell or Command Prompt, the sub-level progress indicator before `resolved <N> repo(s) to pull` rendered as an unprintable tofu replacement character `[?]` instead of an arrow:
```text
  → gitmap pull (cwd: C:\Users\Administrator)
    [?] resolved 52 repo(s) to pull
```
While the header arrow `→` rendered clearly, the sub-level indicator displayed a missing glyph box.

---

## 2. Root Cause

1. **Unsupported Unicode Glyph in Console Monospace Fonts**:
   In `cli/cmdpull/pull.go`, `resolveSubArrow()` was defined as:
   ```go
   func resolveSubArrow() string {
       if glyphs.Resolve() == glyphs.ModeSafe {
           return "->"
       }

       return "↳"
   }
   ```
   The downward arrow with tip rightwards `↳` (`\u21b3`) is not included in standard Windows monospace fonts (such as Consolas, Lucida Console, or Courier New).
2. **Terminal Host Detection vs Font Availability**:
   On Windows hosts running modern terminals (e.g., Windows Terminal where `WT_SESSION` is set), `glyphs.Resolve()` resolves to `ModeRich`. In rich mode, UTF-8 strings pass through directly without fallback rewriting. Consequently, `↳` was output directly to stdout. Because the user's active console font lacked the `\u21b3` glyph, Windows Terminal rendered the fallback replacement symbol `[?]`.
3. **Inconsistent Glyphs Across Subcommands**:
   The primary invocation header in `pull.go` used `→` (`\u2192`), which is universally present across console fonts. However, sub-items in `pull.go` and `push.go` used `↳`, causing visual inconsistencies and rendering errors.

---

## 3. Resolution

1. **Harmonize Subarrow with Universal Right Arrow**:
   Updated `resolveSubArrow()` in `cli/cmdpull/pull.go` to return `→` (`\u2192`) in rich mode and `->` in safe mode:
   ```go
   func resolveSubArrow() string {
       if glyphs.Resolve() == glyphs.ModeSafe {
           return "->"
       }

       return "→"
   }
   ```
2. **Wire Subarrow Helper into Push Subcommand**:
   Replaced hardcoded `↳` in `cli/cmdpull/push.go` (lines 55, 66, 134) with `resolveSubArrow()`, ensuring consistent and safe glyph rendering across `pull` and `push`.
3. **Unit Tests**:
   Added `TestResolveSubArrow` in `cli/cmdpull/pull_ui_test.go` to verify mode-specific arrow resolution (`->` in safe mode, `→` in rich mode).
4. **Verification**:
   Verified that `go test ./cmdpull -v` passes cleanly, all Go format and coding guideline linters pass, and `gitmap pull` renders without unprintable characters.

---

## 4. Prevention & Learnings

- **Unicode Console Font Portability**: Standard Windows monospace fonts support Basic Latin, Latin-1, and fundamental punctuation/arrow glyphs (such as `→` `\u2192`, `•` `\u2022`, `✔` `\u2714`, `✖` `\u2716`), but rarely include specialized directional arrows like `↳` (`\u21b3`). Avoid using specialized Unicode 3.0+ arrow symbols without verifying fallback behavior in default console fonts.
- **Centralized Glyph Resolution**: Prefer shared glyph resolvers over inline raw Unicode literals in command output.
