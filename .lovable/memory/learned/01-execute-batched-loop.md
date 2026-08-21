# Learned Patterns & Conventions: Batched Refactoring Execution

**Slug:** execute-batched-loop
**Date:** 2026-08-22

## 1. Boolean Extraction Scope Safety
- When extracting negative booleans (e.g. !isPRCURL(token) to isNotPRCURL := !isPRCURL(token)), ensure the variable is not redeclared in the same scope using := multiple times. Use distinct domain-specific names (e.g. isFlagUpdateSkipped, isEnvUpdateSkipped).

## 2. Multi-Agent Work Locking
- Using .lovable/temp/active-locks.json prevents merge collisions by ensuring disjoint file sets across concurrent agents.
- Agent 1: gitmap/archive/
- Agent 2: gitmap/cmd/move.go, gitmap/cmd/rm.go, gitmap/completion/powershell.go, src/pages/Troubleshooting.tsx
- Agent 3: gitmap/cliexit/kind.go, gitmap/diff/tree.go

## 3. Strict Coding Guidelines Compliance
- Functions exceeding 15 lines refactored into $\le 15$ line helpers.
- Nested if blocks flattened into early guard clauses.
- Enum type aliases suffixed with Type (e.g. KindType, EntryKindType).
