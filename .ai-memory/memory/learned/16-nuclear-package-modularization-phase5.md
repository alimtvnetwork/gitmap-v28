# Learned: Nuclear Package Modularization & Heavy Test Isolation (Phase 5)

## 1. Context & Motivation
Large Go packages (such as `cli/cmd` with >1,000 files) severely impair developer ergonomics, compilation velocity, and code analysis. Concurrently, routine unit test execution can degrade from milliseconds to minutes if slow tests (those executing subprocess `git` commands, disk repos, or network calls) are mixed into package test suites.

## 2. Core Architectural Principles

### Strict Unidirectional DAG Architecture
- Packages must form a strict Directed Acyclic Graph (DAG) with zero import cycles.
- **Leaf Packages**: `cli/constants`, `cli/store`, `cli/apperror`, `cli/cliexit`, `cli/pipelinedb`, `cli/repodb`, `cli/model`. Leaf packages never import upper domain or orchestration layers.
- **Domain Subpackages**: `cli/cmdfixrepo`, `cli/cmddb`, `cli/cmdpipeline`, `cli/cmdmacro`, `cli/cmdzip`, etc. Domain packages import ONLY leaf packages (and sibling packages lower in the DAG). Domain packages never import `cli/cmd`.
- **Root Orchestrator**: `cli/cmd` imports domain subpackages and routes CLI subcommands downwards via thin bridges in `cli/cmd/clihelpers.go`.
- **Entrypoint**: `cli/main.go` imports `cli/cmd`.

### Heavy Test Isolation to `cli/tests/heavy_test`
- Any test that executes `exec.Command("git", ...)` or initializes real on-disk git repositories belongs exclusively in `cli/tests/heavy_test/` under `package heavy_test`.
- Routine package tests (`*_test.go` within domain packages) must remain 100% in-memory unit tests that execute in <0.01s.
- `03-ai-scripts/33-test-inventory-generator.py` parses individual test function bodies and ignores documentation comments, accurately categorizing tests into:
  - `tier: "slow"` (duration >= 4.0s) -> isolated in `cli/tests/heavy_test`
  - `tier: "fast"` (duration = 0.005s) -> standard unit tests

### Bridge Pattern in `cli/cmd/clihelpers.go`
- Subpackage public APIs are exposed via `exports.go`.
- `cli/cmd/clihelpers.go` provides backward-compatible function and type aliases (e.g., `func runFixRepo(args []string) error { return cmdfixrepo.RunFixRepo(args) }`), ensuring existing command tables and internal test dispatchers continue functioning without churn.
