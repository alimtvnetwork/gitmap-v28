# 116-nuclear-package-splitting-and-slow-tests.md: Nuclear Package Modularization, Heavy Test Segregation & Test Inventory Estimation

**Status: completed**

## 1. Executive Summary & Problem Diagnosis

The GitMap codebase currently suffers from severe package bloat and high test latency:
1. **Monolithic Package Bloat**:
   - `gitmap/cmd` contains **1,245 Go files** (907 implementation files + 338 test files), accounting for **49.1%** of the entire repository.
   - Any test or build targeting `cmd` forces the Go compiler to parse, type-check, compile, and link all 1,245 files together into a single monolithic test binary.
2. **Heavy vs. Light Test Pollution**:
   - Out of 700 test files across `gitmap/`, **36 test files** execute heavy external operations (spawning subprocesses via `exec.Command`, invoking `git`, listening on network sockets, or calling `time.Sleep`).
   - **19 of these heavy test files** reside directly inside `gitmap/cmd` (e.g. `fixgit_test.go`, `gomod_integration_test.go`, `hygiene_integration_test.go`, `pushpull_transport_e2e_test.go`, `reporeclone_e2e_test.go`).
   - Mixing 111 heavy tests with 2,387 fast, pure in-memory unit tests drags down development and testing velocity.
3. **Zero Test Inventory Timing**:
   - `.lovable/test-inventory.json` catalogs 2,498 tests, but all tests currently record `duration_sec: 0.0` and `last_status: "never_run"`.

## 2. Acyclic Architecture & Go Import Cycle Prevention

In Go, package dependency cycles (`A -> B -> A`) are strictly banned by the compiler (`import cycle not allowed`).
Our dependency analysis confirms that:
- `gitmap/cmd` is at the very top of the hierarchy: it imports core packages (`store`, `cloner`, `scanner`, `macro`, `constants`, etc.).
- **Zero core packages import `gitmap/cmd`** (only `main.go` and `tests/cmd_test` import it).

### Strict DAG (Directed Acyclic Graph) Rules
1. **Extracted Domain Subpackages**:
   - High-density command groups (e.g., `chromeprofile` with 30 files, `agy` with 29 files, `prompt` with 28 files, `install` with 21 files, `pipeline` with 17 files) can be extracted into dedicated leaf or domain packages.
   - Extracted packages MUST NOT import `gitmap/cmd`. All shared utilities, types, and flags are passed down via parameters or shared leaf packages (`gitmap/constants`, `gitmap/model`, `gitmap/cliexit`, `gitmap/apperror`).
2. **Heavy Test Isolation (`gitmap/tests/heavy_test`)**:
   - Heavy tests will be moved to `gitmap/tests/heavy_test/` (external blackbox test package `package heavy_test`).
   - External test packages can import `gitmap/cmd` and domain packages without introducing ANY cycles into the production build tree.

## 3. Test Inventory Estimation Strategy

1. **Static & Heuristic Duration Estimation**:
   - Analyze AST for heavy calls: `exec.Command` (+1.5s - 5.0s), `git` CLI calls (+2.0s - 8.0s), `time.Sleep` (+duration), network sockets (+0.5s - 2.0s).
   - Fast unit tests (pure Go structs, string formatting, CLI flags) estimated at `<0.01s`.
2. **Inventory Tagging**:
   - Update `.lovable/test-inventory.json` with duration estimates and separate category tags (`tier: "unit"` vs `tier: "heavy"`).
   - Relocate the 111 heavy tests into `gitmap/tests/heavy_test` so routine `go test ./...` skips heavy suites by default.

## 4. Master Task Breakdown

1. `01-heavy-test-analysis-and-inventory-estimation.md`:
   - Estimate durations for all 2,498 tests in `.lovable/test-inventory.json`.
   - Classify the 111 heavy tests and categorize fast vs slow tests.
2. `02-heavy-test-package-segregation.md`:
   - Extract the 19 heavy test files from `gitmap/cmd` into `gitmap/tests/heavy_test/`.
   - Verify `gitmap/cmd` unit tests run in isolation with zero subprocess latency.
3. `03-cmd-monolith-subpackage-decomposition.md`:
   - Extract distinct, cohesive command groups (e.g., `chromeprofile`, `agy`, `prompt`) into dedicated domain packages under `gitmap/`.
   - Verify strict acyclic dependency structure with `go vet ./...`.
4. `04-test-inventory-sync-and-verification.md`:
   - Re-scan and generate updated `.lovable/test-inventory.json` with new package paths.
   - Record modified files via `33-test-inventory-generator.py --record`.
   - Verify formatting with `26-go-code-formatter.py` and newline rules.

## 5. Execution Summary & Deliverables

1. Extracted gitmap/tests/heavy_test (11 heavy test files, 15 tests) to isolate slow subprocesses.
2. Extracted gitmap/cmdprompt (31 files) into dedicated modular domain package.
3. Extracted gitmap/cmdpurge (8 files) into dedicated modular domain package.
4. Updated .lovable/test-inventory.json: 121 heavy tests (259.53s) vs 2724 unit tests (13.62s).
5. 100% clean compilation verified via go vet ./... and go build ./... with zero import cycles.
