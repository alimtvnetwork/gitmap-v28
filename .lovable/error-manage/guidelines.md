# Error Management Guidelines

These rules dictate how errors should be caught and logged across the codebase to reduce scattered try/catch blocks and guarantee unified logging formatting.

## 1. Centralized Query Wrapper
- **MANDATORY USAGE**: All queries or operations that can fail in Go/TypeScript/PHP/Python **MUST** be wrapped using a centralized wrapper (e.g., `gitmap/store/wrapper.go`, `queryWrapper.ts` or `query_wrapper.py`).
- **NO SCATTERED TRY/CATCH**: Do not write raw `try/catch` blocks scattered across application logic. Use the wrapper to handle the exception, explicitly logging the failure, and returning a structured state.

## 2. Structured Error State
- The wrapper should return a defined type containing explicit boolean properties `isSuccess: true/false` and `isFailure: true/false` instead of relying on consumers to check `!isSuccess` or `!data` or using `is_fail`.

## 3. Explicit Logging
- All caught errors MUST be explicitly logged. The wrapper is responsible for formatting this log.
- Do not introduce random magic strings for logging. If a specific context string is required, pass it to the wrapper and declare it clearly.

## 4. Release Skew Post-Mortem Memory
- **Root Cause Memory**: In `v6.87.1`, a user pushed a tag manually on top of the `v6.87.0` codebase. This caused a release failure because `constants.go` remained hardcoded at `6.87.0`, causing the `Installer Smoke Test` to fail since it verified the exact binary output.
- **Prevention**: Never trigger a release via tag unless the canonical version file (`constants.go`) has been bumped in source control first. If bumping, `git push origin main` must be sent explicitly before `git push origin <tag>`.
