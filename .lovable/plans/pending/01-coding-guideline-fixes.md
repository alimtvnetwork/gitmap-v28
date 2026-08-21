# Plan: Coding Guideline Audit & Enforcement (v4)

Slug: 01-coding-guideline-fixes
Status: pending
Created: 2026-08-22

## Context & Objectives
A comprehensive, multi-pass deep-scan audit was executed across the entire repository (Go backend, React/TypeScript frontend, and PowerShell scripts) against v1.4.5 coding guidelines and top-notch anti-hallucination rules:
- **Function Size**: Strictly $\le 8$ lines preferred, 15 lines hard cap.
- **Return Line Gap**: Blank line strictly mandatory before every `return` statement (except single-line `if`).
- **Boolean Principles**: Positively named (`is*`, `has*`), inverse naming (`isHonest`/`isDishonest`, `isFound`/`isMissing`), zero inverted `!is*` checks.
- **Go Wrapped Results**: Single return parameter wrapping data and mutually exclusive status (`IsSuccess`, `IsFailed`).
- **File Size Caps**: Files $\le 300$ lines, React components $\le 100$ lines, structs $\le 120$ lines.
- **Zero Magic Literals & Zero Swallowed Errors**: Centralized constants and contextual error logs.

---

## Root Cause & Fallout Analysis

### 1. Inverted Booleans & Negative Logic
- **Root Cause**: Developers naturally write `if !ok` or `if (!open)` inline instead of capturing the state in an explicit positive or inverse boolean variable.
- **Blast Radius & Fallout**: Refactoring to `isMissing := !ok; if isMissing == true` has zero external fallout while improving readability and static verification.

### 2. Return Line-Gaps
- **Root Cause**: Inline calculations immediately preceding `return res` without a separating newline.
- **Blast Radius & Fallout**: Purely stylistic; automated formatting and AST tools verify zero semantic drift.

### 3. Monolithic Functions & File Sizes
- **Root Cause**: Gradual feature additions caused orchestration functions (e.g. `scanner.Scan`, `cloner.loadRecords`, React pages) to exceed 15 lines and files to exceed 300 lines.
- **Blast Radius & Fallout**: Splitting into sub-components and helper functions requires exporting or keeping package-private helpers within the same directory. All unit tests must be re-run to confirm zero regression.

### 4. Wrapped Booleans & Error Envelopes
- **Root Cause**: Go functions returning `(bool, error)` or bare `bool` rather than a unified `Result` struct.
- **Blast Radius & Fallout**: Refactoring internal helpers to return a result struct encapsulates status without breaking public API callers when wrapped appropriately.

---

## Execution Subtasks (150+ Steps across 6 Subtasks)

| Subtask File | Scope | Steps | Target Files |
|---|---|---|---|
| [`01-inverted-booleans.md`](../subtasks/01-coding-guideline-fixes/01-inverted-booleans.md) | Inverted Booleans & Inverse Naming | 30 steps | `gitmap/clonefrom/*`, `gitmap/clonenext/*`, `gitmap/clonenow/*`, `src/components/*` |
| [`02-return-line-gaps.md`](../subtasks/01-coding-guideline-fixes/02-return-line-gaps.md) | Missing Blank Lines Before Returns | 30 steps | `gitmap/archive/*`, `gitmap/clonefrom/*`, `gitmap/cloner/*`, `gitmap/store/*` |
| [`03-function-size-decomposition.md`](../subtasks/01-coding-guideline-fixes/03-function-size-decomposition.md) | Function Decomposition (<8 lines) | 30 steps | `gitmap/archive/*`, `gitmap/cliexit/*`, `gitmap/clonefrom/*`, `gitmap/cloner/*` |
| [`04-file-size-caps-and-react-splits.md`](../subtasks/01-coding-guideline-fixes/04-file-size-caps-and-react-splits.md) | File Caps (<300 lines, React <100 lines) | 25 steps | `scanner/scanner.go`, `GenericCLI.tsx`, `Release.tsx`, `Install.tsx`, `TabOrderMap.tsx` |
| [`05-wrapped-booleans-and-single-return.md`](../subtasks/01-coding-guideline-fixes/05-wrapped-booleans-and-single-return.md) | Go Wrapped Result Structs | 20 steps | `gitmap/cloner/*`, `gitmap/scanner/*`, `gitmap/store/*`, `gitmap/cmd/*` |
| [`06-magic-constants-and-nested-ifs.md`](../subtasks/01-coding-guideline-fixes/06-magic-constants-and-nested-ifs.md) | Enums, Constants & Flattening | 25 steps | `gitmap/constants/*`, `gitmap/cmd/*`, `src/pages/*` |

---

## Verification Plan

### Automated Verification
1. `go test ./... -short` across all Go packages to verify complete functional parity.
2. AST verification script checking all Go functions for $\le 15$ lines and return line gaps.
3. TypeScript / React compilation check (`npx tsc --noEmit` / `npm run build`).
4. Parity and helptext test suites (`go test ./gitmap/constants/... ./gitmap/helptext/...`).

### Manual Inspection
1. Verify all `readme.md` files remain strictly lowercase.
2. Confirm zero `.lovable/memories/` directories created.
3. Verify strictly positive boolean naming (`is*`, `has*`, inverse naming).
