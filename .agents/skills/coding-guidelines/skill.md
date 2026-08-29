---
name: coding-guidelines
description: >-
  Use this skill to audit, review, and enforce coding guidelines across all languages (Go, TS, Python, PHP, C#).
---

# Coding Guidelines Audit & Enforcement Skill

Autonomously audits, reviews, and enforces repository coding guidelines against `.lovable/coding-guidelines/coding-guidelines.md`, `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `spec/17-consolidated-guidelines/`.

## Core Audit Checkpoints

1. **Size Limits & Anti-Compression:**
   - Functions: <= 15 lines max (target <= 8 lines).
   - Files: <= 100 lines standard cap (hard ceiling 300 lines).
   - React components (.tsx): <= 100 lines max per file.
   - Zero line compression cheating (never merge multiple statements into one line).

2. **Nesting & Braces:**
   - Zero nested `if` statements (nesting depth <= 1).
   - Invert conditions into guard clauses and early returns.
   - Return New Line (R13-R16): blank line before `return`/`throw` and after closing `}`.

3. **Booleans & Logic:**
   - Implicit boolean evaluation only (`if isReady`, never `== true` or `== false`).
   - Affirmative naming (`is*`, `has*`, `can*`, `should*`).
   - No inverted success checks (`!isSuccess` -> `isFail` or `isMissing`).
   - No mixed polarity chains in single condition.

4. **Error Management:**
   - Zero swallowed errors (`catch {}`, `_ = err`).
   - Wrap errors in `AppError` with operational context.
   - Universal response envelopes.
   - Zero dual-handling (no log + return or panic + return).

5. **Constants, Enums & Casing:**
   - Zero magic strings or numbers.
   - Enum type names MUST end with `Type` suffix.
   - PascalCase acronym abbreviations: `Id`, `Url`, `Api` (never `ID`, `URL`, `API`).

6. **0-100 Scoring Formula:**
   - `Score = 100 - (Critical_Count * 10) - (Major_Count * 5) - (Minor_Count * 2)`
   - Bounded: `final_score = max(0, min(100, Score))`
