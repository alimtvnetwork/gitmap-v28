# Subtask 02: Incremental Fingerprinting & File Change Detection Engine

## Objective
Implement two-tier Git change detection and declarative gate dependency mapping to skip unchanged passing gates in O(1) time.

## Requirements
1. **Tier 1 (Fast-Path Git Delta)**:
   - Read `git rev-parse HEAD` and `git status --porcelain=v1 -uall`.
   - Calculate changed files $\Delta_{repo}$ between previous pass commit and current working tree.
2. **Tier 2 (Gate Dependency Mapping)**:
   - Define declarative `GateSpec` for all 33 gates including tool scripts, config files, relevant path globs, exclusion globs, and upstream artifacts.
   - Gates are skipped if and only if:
     - Previous run was `PASSED` (code == 0).
     - Command line args, cwd, and env hash match.
     - Tool scripts and configs are unmodified.
     - No file in $\Delta_{repo}$ matches relevant paths.
     - Required input artifacts exist and match recorded `mtime`/size.
     - Upstream producer gates were not re-executed in this session.
3. **Cache Invalidation & Controls**:
   - Add CLI options: `--force`, `--fresh`, `--clean`, `--no-cache`.
   - When `--force` is given, purge/bypass cache and execute all gates.
4. **Coding Guidelines**:
   - All functions $\le 15$ lines.
   - Blank line before every return statement.
   - Affirmative booleans (`is_*`, `has_*`).
   - Zero swallowed exceptions.

## Target Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [x] Clean run with no modified files skips previously passing gates in < 2 seconds.
- [x] Modifying a file invalidates only relevant gates and their dependents.
- [x] `--force` triggers full re-run of all 33 gates.
