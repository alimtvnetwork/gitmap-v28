# Avoid 01: Test and Build Execution Guards

- **Category:** Anti-Pattern & Strict Constraint
- **Status:** Mandatory Enforcement

## 1. Core Rule
NEVER run unit tests, integration test suites, compiler builds, or heavy CI/CD runners during routine coding tasks, guideline refactors, or consolidation workflows unless explicitly commanded by the user.

## 2. Rationale
- Routine test runs waste execution quota and trigger resource starvation.
- Subprocesses can create lock contention on database files and generate unwanted build artifacts.

## 3. Enforcement
- Skip test flags must be active by default in all automation scripts.
- Only run verification linters (`check-relative-paths.py`, `check-forbidden-strings.py`).
