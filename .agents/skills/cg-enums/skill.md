---
name: cg-enums
description: Autonomously scan, audit, refactor, and verify repository-wide constants, enums, magic literals, and rune conversions against spec/02-coding-guidelines/, enforcing *Type suffixes, constants centralization, and eliminating rune number conversions.
---

# Skill: Constants & Enums Architecture (`cg-enums`)

This skill governs autonomous scanning, auditing, refactoring, and verification of constants, enums, magic string/number literals, and rune conversions across Go, TypeScript, and Python codebases adhering to `spec/02-coding-guidelines/`.

## Core Architectural Directives

1. **Mandatory `*Type` Suffix (Rule 1):**
   - Every enum type definition across all languages MUST end with `Type`.
   - Examples:
     - Go: `type ConflictModeType uint8`, `type InputKindType uint8`, `type SourceKindType uint8`
     - TypeScript: `export enum BuildStatusType { ... }`, `export const TabIdType = { ... } as const`
     - Python: `class SeverityType(StrEnum): ...`
   - In Go, when refactoring existing enums to `*Type`, always provide a backward-compatible type alias (`type ConflictMode = ConflictModeType`) to prevent consumer breakage.

2. **Zero Raw Rune Numeric Conversions (Rule 2):**
   - Raw integer numeric character casts like `rune(10)`, `rune(13)`, `rune(0)` are strictly banned.
   - Always use centralized character/string constants from `constants` or `core` packages (e.g., `NewLineUnix`, `LineFeedByte`).

3. **Dedicated Packages & Centralization (Rule 3):**
   - Constants and enums must reside in dedicated directories and files (`constants/`, `enums/`, `src/types/`, `src/enums/`).
   - Magic strings and numbers repeated across files must be extracted into dedicated constant files.

4. **No Raw String Unions for Enums in TypeScript (Rule 4):**
   - TypeScript status, mode, and state collections must use typed enums or `as const` object maps ending with `*Type` rather than raw string unions.

5. **Size & Hygiene Constraints:**
   - Files must strictly adhere to canonical size limits (target $\le 100$ coding lines).
   - Functions must not exceed 8–15 lines.
   - Never use line compression or semicolon flattening to bypass size caps.

6. **Targeted Verification Protocol:**
   - Run `python linter-scripts/check-enum-guidelines.py` to audit enums, constants, and rune casts.
   - Run `node linter-scripts/check-enum-and-boolean.mjs` for TS enum and boolean validation.
   - Run targeted type checks and linters (`go vet ./...`, `npx tsc --noEmit`).
   - DO NOT run `06-cicd-local-runner.py` or execute unit test suites during routine refactoring turns.
   - Record modified files under atomic lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
