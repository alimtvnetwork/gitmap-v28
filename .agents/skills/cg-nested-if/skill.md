---
name: cg-nested-if
description: Autonomously scan, audit, refactor, and verify repository-wide control flow flattening, eliminating nested if statements (nesting depth > 1) and enforcing guard clauses, early returns, and mandatory braces against 02-spec/02-coding-guidelines/.
---

# Skill: Nested If Elimination & Guard Clauses (`cg-nested-if`)

This skill governs autonomous scanning, auditing, refactoring, and verification of conditional structures across Go, TypeScript, and Python codebases, eliminating nested `if` statements (nesting depth > 1) and single-line compressed `if`s.

## Core Architectural Directives

1. **Zero Nested `if` Mandate (Absolute Ban):**
   - Nested `if` blocks (an `if` statement directly or indirectly inside another `if` body or `else` body) are strictly prohibited across all languages (Go, TypeScript, Python, PHP).
   - Maximum conditional nesting depth within any function or closure is exactly 1.
   - Code reviews finding a nested `if` are an automatic rejection.

2. **Flattening Mechanisms:**
   - **Inverted Guard Clauses:** Convert outer conditions into early return/exit checks (`if (!condition) { return ...; }`), leaving primary logic unindented at root block depth.
   - **Helper Function Decomposition:** Extract nested processing logic into focused, single-purpose helper functions (target $\le 8$ lines, hard cap 15 lines).
   - **Combined Affirmative Conditions:** Combine sequential prerequisite checks using clear, named boolean variables or safe boolean operators (`&&`) where appropriate.

3. **Mandatory Braces & Anti-Compression (Rule 1):**
   - Single-line `if` statements without curly braces `{}` are strictly forbidden in all languages:
     - ❌ `if (x) return y;`
     - ❌ `if err != nil { return err }` on one line
     - ✅ Multi-line with explicit braces:
       ```typescript
       if (x) {
         return y;
       }
       ```
   - Never collapse statements onto one line to bypass function size caps.

4. **Formatting & Spacing Standards:**
   - Exactly one blank line before every `if` statement (unless at the very start of a block).
   - Exactly one blank line before every `return` statement (unless it is the only statement in a block).
   - Exactly one blank line after every closing brace `}` of a control flow block.

5. **Exemptions (Spec Grounded):**
   - Enum `switch` parsers (`FromString` / `fromString`) with a `default` error return fallback.
   - Optional-field validation where outer condition checks existence (`isDefined()`) and inner condition validates the same field.

6. **Targeted Verification Protocol:**
   - Run `python linter-scripts/check-nested-ifs.py` to audit both Go and TypeScript files.
   - Run targeted type checks and linters (`npx tsc --noEmit`, `go vet ./...`).
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
