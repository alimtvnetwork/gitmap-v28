# Avoid 02: Macro Idempotency and Binary Hygiene

- **Category:** Anti-Pattern & Strict Constraint
- **Status:** Mandatory Enforcement

## 1. Core Rule
NEVER author non-idempotent macro actions or leave orphan child processes in developer profiles.

## 2. Rationale
- Non-idempotent scripts crash on repeated runs when target assets already exist or are absent.
- Lingering background processes hold open ports and file locks, corrupting subsequent runs.

## 3. Enforcement
- Always specify preconditions (`exists`, `absent`, `matches`, `process_running`).
- Terminate child process groups on macro task finalization.
