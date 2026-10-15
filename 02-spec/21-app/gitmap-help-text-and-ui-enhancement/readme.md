# gitmap-help-text-and-ui-enhancement (Program 273)

**Version:** 1.0.0
**Updated:** 2026-10-16
**Status:** Draft
**AI Confidence:** High
**Ambiguity:** Low

---

## Overview

Program 273 delivers the user-ordered terminal help text + UI help overhaul for
GitMap v6.554+. Themes and colors come from the coding-guidelines-v24 design
system — never invented. Every command gets a `help example` section (what it
does + concrete invocation + expected output) driven from a single data
registry. The missing `merge write` / `mr` command (needed for the Lovable
migration) is audited and specified here, and the Windows
VirTool:MSIL/Subti.C heuristic RCA is recorded from the ground-truth commit
385aba38e so another AI can continue it.

---

## Keywords

`help-text` · `themes` · `design-system` · `help-example` · `merge-write` · `mr` · `virus-rca` · `program-273`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | High |
| Ambiguity | Low |
| Health Score | 100/100 |

---

## Files

| # | File | Description |
|---|------|-------------|
| 01 | [01-overview.md](./01-overview.md) | Scope, waves, hard rules, and file map |
| 02 | [02-theme-design.md](./02-theme-design.md) | Theme registry sourced from coding-guidelines-v24 |
| 03 | [03-help-example-format.md](./03-help-example-format.md) | Canonical `help example` format (data-driven, single source of truth) |
| 04 | [04-missing-command-audit.md](./04-missing-command-audit.md) | Missing-command audit approach (dispatcher vs help) |
| 05 | [05-merge-write-command.md](./05-merge-write-command.md) | `merge write` / `mr` design (Lovable migration) |
| 06 | [06-twenty-new-command-proposals.md](./06-twenty-new-command-proposals.md) | Proposal template + process (proposals only, no implementation) |
| 07 | [07-windows-virus-rca.md](./07-windows-virus-rca.md) | Windows VirTool:MSIL/Subti.C RCA (ground truth 385aba38e) |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | Testable acceptance criteria for WS2–WS6 |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Structural health report |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Program ledger (273) | `../../../.ai-memory/plans/273-gitmap-help-text-and-ui-enhancement.md` |
| Execution plan | `../../../.ai-memory/plans/gitmap-help-text-and-ui-enhancement.md` |
| Subtasks | `../../../.ai-memory/plans/subtasks/gitmap-help-text-and-ui-enhancement/` |
| CLI architecture cluster | `../01-cli-architecture/01-architecture-spec.md` |
| Design system source | coding-guidelines-v24 `02-spec/07-design-system/16-theme-catalogue-and-palettes.md` |
| Design tokens source | coding-guidelines-v24 `02-spec/07-design-system/17-theme-tokens.json` |
| Spec authoring guide | `../../01-spec-authoring-guide/04-required-files.md` |
