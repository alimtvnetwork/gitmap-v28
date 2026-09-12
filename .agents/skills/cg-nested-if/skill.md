---
name: cg-nested-if
description: Autonomously scan, audit, refactor, and verify repository-wide control flow flattening, eliminating nested if statements (nesting depth > 1) and enforcing guard clauses, early returns, and mandatory braces against spec/02-coding-guidelines/.
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
