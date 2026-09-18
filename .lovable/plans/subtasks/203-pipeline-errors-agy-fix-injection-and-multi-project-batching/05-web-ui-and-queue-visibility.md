# Subtask 05: Web UI Documentation, Commands Dataset & Queue Visibility

## Objective
Update the Web UI dataset and documentation pages to document direct AGY injection, multi-project batch limits, and queue visibility.

## Target Files
- `src/data/commands.ts`
- `src/pages/Pipeline.tsx`
- `src/pages/AGYPrompts.tsx`

## Status
Completed: 2026-09-18

## Verification
- `src/data/commands.ts`: Registered `pipeline errors agy fix` (aliases: `pipeline fix errors agy`, `aef`, `agy fix-pipeline`) with flags (`--all`, `--projects`, `--limit`, `--reset-batch`, `--no-inject`), examples, and cross-links.
- `src/pages/Pipeline.tsx`: Added CLI command cards for `gitmap pipeline errors agy fix` and `gitmap pipeline errors agy fix --all`.
- `src/pages/AGYPrompts.tsx`: Added documentation section for Direct Antigravity Injection & Multi-Project Batching.
- TypeScript AST / syntax validated. Zero absolute paths.
