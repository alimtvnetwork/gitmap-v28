# Specification 288: Codebase Improvements Prompts & Skills Architecture

- **Sequence:** 288
- **Slug:** `288-create-prompts-for-codebase-improvements`
- **File:** `02-spec/21-app/288-create-prompts-for-codebase-improvements/01-architecture-spec.md`
- **Status:** APPROVED & SPECIFIED

---

## 1. Architectural Mission & Systems Objectives

The objective of this specification is to architect, standardize, and integrate three canonical AI improvement prompts and their corresponding multi-agent skills across the Prompt Architect library (`01-prompts/30-improvements/`), Google Antigravity Agent Skills (`.agents/skills/`), and Cursor Skills (`.cursor/skills/`):

1. **Prompt 1: Find 100 Improvements in the UI Level** (`01-prompts/30-improvements/01-find-100-improvements-in-ui-level.md`):
   - Directs autonomous AI agents to discover, audit, and systematically catalog 100 UI/UX improvements.
   - Enforces detailed sub-points: target component, what it fixes, how it improves visual ergonomics/contrast, design tokens, color samples (hex, ANSI SGR, Lipgloss AdaptiveColor), and Bubble Tea/TUI layout patterns.
2. **Prompt 2: Find 100 Improvements in the Codebase** (`01-prompts/30-improvements/02-find-100-improvements-in-codebase-coding-guidelines.md`):
   - Directs autonomous AI agents to audit source code and uncover 100 concrete code improvements strictly based on grounded coding guidelines and error management specifications.
   - Enforces a rigorous reporting schema: relative file path, line numbers, function/symbol name, defect context, coding guideline citation, and concrete code refactoring diffs.
   - Core tenets: leaf function error sovereignty (never swallow errors), 15-line function cap, zero nested ifs, positive booleans (`is*`, `has*`), and Split-DB repository encapsulation.
3. **Prompt 3: Find 100 Improvements in the Whole System** (`01-prompts/30-improvements/03-find-100-improvements-in-whole-system-ux.md`):
   - Directs autonomous AI agents to identify 100 whole-system improvements that enhance user experience, developer velocity, and operational stability.
   - Enforces structured analysis: cause, effect, and non-destructive implementation blueprints ensuring existing codebases, working trees, and databases are protected without regressions.

---

## 2. Directory Hierarchy & Artifact Topology

```text
gitmap/
├── 01-prompts/
│   ├── readme.md (Updated catalog index)
│   └── 30-improvements/
│       ├── readme.md (Subsystem catalog & guide)
│       ├── 01-find-100-improvements-in-ui-level.md
│       ├── 02-find-100-improvements-in-codebase-coding-guidelines.md
│       └── 03-find-100-improvements-in-whole-system-ux.md
├── .agents/skills/
│   ├── find-100-ui-improvements/
│   │   └── skill.md
│   ├── find-100-codebase-improvements/
│   │   └── skill.md
│   └── find-100-system-ux-improvements/
│       └── skill.md
├── .cursor/skills/
│   ├── find-100-ui-improvements/
│   │   └── skill.md
│   ├── find-100-codebase-improvements/
│   │   └── skill.md
│   └── find-100-system-ux-improvements/
│       └── skill.md
└── 02-spec/21-app/288-create-prompts-for-codebase-improvements/
    ├── 01-architecture-spec.md
    ├── 02-prompts-and-skills-design-spec.md
    └── 03-sync-and-compatibility-spec.md
```

---

## 3. Design Principles & Non-Negotiable Invariants

1. **Strict Error Return Sovereignty:** Prompts and generated code must require that fallible methods never return void and never swallow errors with blank discards (`_ = ...`). Leaf functions must return `error` or `*appfault.AppError`.
2. **Non-Destructive Safe Evolution:** Whole-system improvements must include explicit rollback blueprints, safe defaults, and autostash protections to ensure zero user data loss.
3. **Multi-Platform ANSI & Color Fidelity:** UI prompts must provide grounded color palettes mapping directly to GitMap's semantic roles (`RoleFG`, `RolePrimary`, `RoleAccent`, `RoleWarn`, `RoleError`, `RoleOK`, `RoleMuted`) with 24-bit TrueColor and 256-color SGR fallbacks.
4. **Synchronizability & Parity:** All prompts and skills must adhere to the sync boundary rules of `03-ai-scripts/38-sync-prompts-scripts.py`, ensuring automatic propagation to all 43 connected repositories.
5. **Path Hygiene:** All file paths must be strictly relative from repository root. No absolute paths or file URIs.
