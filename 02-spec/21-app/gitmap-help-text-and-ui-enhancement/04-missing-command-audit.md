# 04 — Missing-Command Audit Approach

## 1. Goal

Find every gap between what the dispatcher CAN do and what help SHOWS, and
every gap between what help promises and what the dispatcher implements.
Output: a ranked gap list with a decision per gap (implement / add help /
alias / remove promise). The known headline gap is `merge write` / `mr`.

## 2. Procedure

1. **Enumerate dispatcher verbs.** Collect every registered verb from the
   dispatch tables (`cli/cmd/root*.go` dispatch entries) — verbs AND all
   aliases. Normalize to lowercase; keep the verb→alias mapping.
2. **Enumerate help entries.** Collect every key in the help registry
   (`cli/cmd/help.go` entries) plus every generated `cli/helpdoc/*.md`
   file, plus every key in the new `examples.go` registry (spec 03).
3. **Diff A: verbs without help.** Any dispatcher verb (or alias) with no
   help entry is a defect — every alias is kept, so every alias gets help.
4. **Diff B: help without verbs.** Any help entry with no dispatcher verb is
   a broken promise — remove the help or implement the command.
5. **Diff C: example coverage.** Registry keys from spec 03 vs dispatcher
   verbs — any verb without an example is flagged.
6. **Lovable-migration scan.** Specifically check for migration-critical
   verbs: `merge write`/`mr` (PR creation), plus repo export/import surface.
   These are ranked P0 regardless of frequency.
7. **Decision pass.** For each gap: implement (P0/P1 only), add missing
   help/example, alias to existing, or remove stale help. Record the
   decision in the audit report.

## 3. Audit artifact

`04-conversation`-style dated markdown in the workspace root of the program
(`.ai-memory/plans/subtasks/gitmap-help-text-and-ui-enhancement/`) plus the
machine-readable verb list embedded in the audit subtask. The audit MUST run
with `gitmap aum search` / `gitmap find` only (TOTAL BAN on grep/rg).

## 4. Scope boundary

The audit PRODUCES the gap list; implementing non-`mr` gaps is NOT in program
273 scope (they become proposals for WS4 or later programs). Only
`merge write` / `mr` is implemented in WS3 — it is the one gap the user
explicitly ordered.
