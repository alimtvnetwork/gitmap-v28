# Plan 288: Create Prompts for Codebase Improvements

## Status: Completed
- **Plan ID:** 288
- **Slug:** `288-create-prompts-for-codebase-improvements`
- **Spec Reference:** [02-spec/21-app/288-create-prompts-for-codebase-improvements/01-architecture-spec.md](../../../02-spec/21-app/288-create-prompts-for-codebase-improvements/01-architecture-spec.md)
- **Scope:** Canonical Prompts Creation, 3 Improvement Prompts (UI, Codebase Guidelines, System UX), Antigravity Agent Skills, Cursor Skills, Multi-Repo Sync Architecture
- **Created At:** 2026-10-22
- **Parent Goal:** Establish canonical, production-grade prompts and skills under `01-prompts/30-improvements/`, `.agents/skills/`, and `.cursor/skills/` to empower autonomous AI agents to discover, triage, and execute 100-item improvement campaigns across UI ergonomics, codebase coding guidelines, and whole-system user experience.

---

## 1. Executive Summary

This plan executes the creation of three authoritative, self-contained AI prompts and their corresponding multi-agent skills:
1. **Prompt 1: Find 100 Improvements in the UI Level** (`01-prompts/30-improvements/01-find-100-improvements-in-ui-level.md`):
   - Instructs the AI on finding 100 concrete UI/UX improvements across terminal TUI (Bubble Tea/Lipgloss), Unicode tables, interactive modal pickers, progress indicators, and embedded web visualizers.
   - Mandates detailed sub-points: what it fixes, how it improves visual ergonomics/contrast, design system tokens, color samples (hex, ANSI SGR codes, Lipgloss AdaptiveColor), and visual hierarchy.
2. **Prompt 2: Find 100 Improvements in the Codebase** (`01-prompts/30-improvements/02-find-100-improvements-in-codebase-coding-guidelines.md`):
   - Instructs the AI on identifying 100 codebase improvements strictly based on grounded coding guidelines and error management rules.
   - Enforces detailed reporting with exact relative file path, line numbers, function/symbol name, defect context, and concrete refactoring code diffs.
   - Covers error return sovereignty (never swallow errors, return error or result wrapper), 15-line limit, zero nested ifs, positive booleans (`is*`, `has*`), and Split-DB repository encapsulation.
3. **Prompt 3: Find 100 Improvements in the Whole System** (`01-prompts/30-improvements/03-find-100-improvements-in-whole-system-ux.md`):
   - Instructs the AI on identifying 100 whole-system improvements that enhance user experience, developer ergonomics, and workflow reliability.
   - Mandates cause-and-effect analysis, non-destructive execution mechanics (safe fallbacks, zero regressions, preserving working trees and databases), and safety verification.
4. **Skills Parity & Multi-Repo Sync Integration:**
   - 3 Agent Skills in `.agents/skills/` (`find-100-ui-improvements`, `find-100-codebase-improvements`, `find-100-system-ux-improvements`).
   - 3 Cursor Skills in `.cursor/skills/` with matching frontmatter and directives.
   - Catalog index in `01-prompts/30-improvements/readme.md` and registered in `01-prompts/readme.md`.
   - Guaranteed compatibility with `03-ai-scripts/38-sync-prompts-skills-scripts.py`.

---

## 2. Subtasks Breakdown

- [x] **Subtask 288.1:** Specifications Authoring & Task Decomposition
- [x] **Subtask 288.2:** Develop Prompt 1: Find 100 Improvements in the UI Level
- [x] **Subtask 288.3:** Develop Prompt 2: Find 100 Improvements in the Codebase (Guidelines & Errors)
- [x] **Subtask 288.4:** Develop Prompt 3: Find 100 Improvements in the Whole System (UX & Safety)
- [x] **Subtask 288.5:** Create Antigravity Agent Skills & Cursor Skills
- [x] **Subtask 288.6:** Prompt Library Indexing, FastGate AST Linter & Atomic Push

---

## 3. Specifications Index

- [01-architecture-spec.md](../../../02-spec/21-app/288-create-prompts-for-codebase-improvements/01-architecture-spec.md)
- [02-prompts-and-skills-design-spec.md](../../../02-spec/21-app/288-create-prompts-for-codebase-improvements/02-prompts-and-skills-design-spec.md)
- [03-sync-and-compatibility-spec.md](../../../02-spec/21-app/288-create-prompts-for-codebase-improvements/03-sync-and-compatibility-spec.md)

---

## 4. Verification & Quality Gates

- **FastGate AST Linter:** Zero violations on relative paths, nested ifs, positive booleans, and 500-line limits.
- **Path Hygiene:** 100% strictly relative Git paths (zero absolute path violations).
- **Git Atomic Commit:** `gitmap cpf "prompts - create prompts and skills for codebase improvements"`.
