# Avoid 03: Git State and Branch Pollution

- **Category:** Anti-Pattern & Strict Constraint
- **Status:** Mandatory Enforcement

## 1. Core Rule
Subagents MUST NEVER run git commands (`git add`, `git commit`, `git push`, `git rm`, `git status`). All Git state manipulations are reserved exclusively for the lead orchestrator.

## 2. Rationale
- Concurrent subagents executing git commands trigger index lock races (`.git/index.lock`).
- Fragmented micro-commits pollute Git history.

## 3. Enforcement
- Subagents execute file modifications using direct filesystem tools only.
- Strict isolation protecting `.ai-memory/plans/pending/` and active open plans.
