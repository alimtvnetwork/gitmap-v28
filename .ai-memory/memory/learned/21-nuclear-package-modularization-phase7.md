# Learned Memory 21: Nuclear Package Modularization (Phase 7) & DAG Decoupling

## Key Insights & Architectural Patterns

1. **Acyclic Package Modularization (Zero Cycles)**:
   - When extracting cohesive command families out of a monolithic package (`cli/cmd`), enforce a strict Directed Acyclic Graph (DAG):
     `Leaf (constants, model, store, apperror, fsutil)` -> `Domain Subpackage (cmdinstaller, cmdchrome, cmdsetup, cmdinstall)` -> `Top-level Command Bridge (cli/cmd)` -> `main.go`.
   - Never import `cli/cmd` from any domain subpackage.
   - If a domain subpackage needs functionality or callbacks from upper layers, use the delegate hook pattern:
     Define `var HookFn func(...)` in the subpackage `exports.go` and inject the implementation inside `cli/cmd/clihelpers.go:init()`.

2. **Cross-Subpackage Ingestion Hierarchy**:
   - Leaf subpackages (e.g. `cmdinstaller`) must have zero dependencies on sibling domain packages.
   - Higher-order domain subpackages (e.g. `cmdinstall`) can cleanly import leaf subpackages (e.g. `cmdinstaller`) without cycles.
   - Forwarders in `cli/cmd/clihelpers.go` maintain exact AST and signature backward compatibility for commands remaining in `cmd`.

3. **Heavy Test Isolation**:
   - Heavy tests that invoke subprocesses (`exec.Command`), git commands, sockets, or sleeps belong in `cli/tests/heavy_test/` (`package heavy_test`).
   - Domain subpackage unit tests remain 100% fast, pure in-memory tests executing in milliseconds.

4. **Boolean Guidelines for Helpers**:
   - When flattening nested `if` statements in helpers, avoid `should*` prefixes (which violate boolean guidelines).
   - Use affirmative `has*`, `is*`, or `supports*` prefixes (e.g., `hasTrailingFlagValue`).
