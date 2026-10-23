# Plan 298 Architecture Specification: Result Wrapper Slice Collections, Types Centralization & AppError Sovereignty

## User Request (Verbatim)
```text
@[01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md]
```

## 1. Executive Summary & Architectural Intent
This specification governs the systematic elimination of multi-value error tuples (`([]T, error)`, `(T, error)`, `(map[K]V, error)`) across Go subsystems, standardizing on single-return monadic containers (`result.Wrap[T]`, `result.WrapSlice[T]`, `result.WrapMap[K, V]`) with pointer-attached null safety (`*result.Wrap[T]`, `*result.WrapSlice[T]`, `*result.WrapMap[K, V]`) and structured `*apperror.AppError` returns.

Following the successful execution of Plan 296 (Store, GitUtil, CmdPipeline) and Plan 297 (Codebase-wide Stutter Eradication), Plan 298 targets the core file processing, packaging, and cloner parser subsystems:
1. `cli/archive`: Archive file preparation, filesystem gathering, and extraction entry listing.
2. `cli/vscodepm`: VS Code Project Manager registry synchronization and file reading.
3. `cli/clonefrom`: Concurrent input parsers (`parseByFormat`, `parseJSON`, `parseCSV`, `readCSVRows`) and JSON schema emitters.
4. `cli/clonenext`: Batch loaders (`LoadBatchFromCSV`, `WalkBatchFromDir`).
5. `cli/clonenow`: Batch parsers (`parseByFormat`, `parseJSONWithSchema`, `parseCSVWithSchema`, `parseTextRows`).
6. `cli/clonepick`: Interactive picker runner (`RunPicker`, `normalisePaths`, `runLsTree`).

## 2. Scope Boundaries
- **In-Scope Subsystems & Files:**
  - `cli/archive/types.go`, `cli/archive/create.go`, `cli/archive/list.go`
  - `cli/vscodepm/types.go`, `cli/vscodepm/sync.go`, `cli/vscodepm/syncat.go`
  - `cli/clonefrom/types.go`, `cli/clonefrom/parse.go`, `cli/clonefrom/parsecsv.go`, `cli/clonefrom/jsonschema.go`
  - `cli/clonenext/types.go`, `cli/clonenext/batch.go`
  - `cli/clonenow/types.go`, `cli/clonenow/parse.go`, `cli/clonenow/parsetext.go`
  - `cli/clonepick/types.go`, `cli/clonepick/picker_run.go`, `cli/clonepick/parse.go`, `cli/clonepick/picker_tree.go`
- **Excluded / Deferred:**
  - Full test suite execution and `go build` (strictly banned under Rule 1).
  - Version bumping or changelog edits (strictly banned under Rule 10).
  - Raw git commands inside subagents (strictly banned under Rule 7).

## 3. Core Technical Invariants
1. **Single Monadic Return:** Fallible operations return `result.Wrap[T]`, `result.WrapSlice[T]`, or `result.WrapMap[K, V]`.
2. **Pointer-Attached Null Safety:** All inspection methods begin with `if r == nil` guards returning safe defaults.
3. **DRY Types Centralization:** All payload models and Result aliases are declared in each package's `types.go` without raw generic stutter across implementation files.
4. **Zero Double-Wrapping:** Never call `FailWrap` or wrap an AppError twice. Use canonical constructors `result.FailSlice[T](apperror.WrapSimple(err, op))` and `result.OkSlice(items)`.
5. **Strict Backward Compatibility:** Public functions with wide external callers provide monadic `*Wrap` variants while delegating legacy tuple signatures to preserve caller stability.
6. **AST Compliance:** Zero warnings from `python 03-ai-scripts/50-fastgate.py`.

## 4. Multi-Agent Partitioning & Ownership Matrix
- **Wave 1 (Planning & Baseline Extraction):**
  - Worker 01 (`research`): Survey `cli/archive` and `cli/vscodepm`.
  - Worker 02 (`research`): Survey `cli/clonefrom` and `cli/clonenext`.
  - Worker 03 (`research`): Survey `cli/clonenow` and `cli/clonepick`.
- **Wave 2 (Implementation & Type Centralization):**
  - Worker 01 (`self`): Owns `cli/archive/` and `cli/vscodepm/`.
  - Worker 02 (`self`): Owns `cli/clonefrom/` and `cli/clonenext/`.
  - Worker 03 (`self`): Owns `cli/clonenow/` and `cli/clonepick/`.
- **Wave 3 (Verification & Push):**
  - Root Orchestrator: Runs FastGate AST linting, checks secrets gate, updates plan registry, and executes atomic GitMap commit `gitmap cpf "archive - monadic result slice collections and types centralization"`.
