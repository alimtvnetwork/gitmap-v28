# Nuclear Modularization Phase 4: Domain Packages & Heavy Test Isolation

## Context & Architecture
To eliminate compiler bloat and test latency in the oversized `cli/cmd` monolith (1,110 files), 4 cohesive domains were extracted into standalone subpackages:
- `cli/cmdmacro` (28 files): Macro execution, add/edit/record, and multi-format export/import.
- `cli/cmdvscode` (28 files): VS Code workspace generation, project manager sync, and duplicate detection.
- `cli/cmdvhost` (13 files): Virtual host configuration, templates, ops, and site management.
- `cli/cmdzip` (8 files): Archive compression, decompression, and zipgroup management.

## Strict DAG Decoupling Rules
1. **Layer 0 (Leaves)**: `constants`, `model`, `store`, `apperror`, `cliexit`, `fsutil`, `ecosystemgroup`, `vscodepm`, `macro`, `archive`.
2. **Layer 1 (Domain Packages)**: `cmdmacro`, `cmdvscode`, `cmdvhost`, `cmdzip`, `cmdagy`, `cmdchromeprofile`, `cmdvmware`, `cmdpurge`, `cmdprompt`.
   - Domain packages NEVER import `cmd`.
   - Domain packages only import Layer 0 leaves and standard library.
3. **Layer 2 (CLI Orchestrator)**: `cli/cmd`.
   - Delegates subcommands to domain packages via bridge functions in `cli/cmd/clihelpers.go`.
4. **Zero Cycles**: `go vet ./...` and `go build ./...` pass with exit code 0.

## Heavy Test Isolation Pattern
- Tests performing external subprocess execution (`exec.Command`), git bare clones, or 200+ commit replay loops must be isolated in `cli/tests/heavy_test/` (`package heavy_test`).
- Core package unit tests remain 100% in-memory fast tests (<0.01s).
- Isolating the 220-commit history bury test reduced `cli/committransfer` test duration from >30s down to 0.599s.
