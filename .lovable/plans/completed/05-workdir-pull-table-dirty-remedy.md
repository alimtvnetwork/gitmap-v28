# Plan: WorkDir Management, Recursive Top-Level Pull, Rich Pull Table & Dirty Remediation

## Context

Full architectural expansion providing:
1. WorkDir management (`gitmap workdir` / `wd`) with default work path override.
2. Recursive top-level Git repository auto-discovery with nested repo pruning.
3. Rich pull table with last commit SHA, active branch, and PR indicator.
4. Detailed dirty status diagnosis and copy-pasteable remediation command recipes.
5. Release version bump and local CI/CD sanity checks.

Inputs:
- .lovable/spec/commands/04-workdir-pull-table-dirty-remedy.md

## Execution Model

50 micro-tasks executed sequentially in a continuous self-loop with full test validation after every step.
