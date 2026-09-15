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
   - At the final step, stage all changes (`git add -A`), commit in a single grouped atomic commit, and push to git (`git push origin <branch>`). NEVER commit each file individually.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
