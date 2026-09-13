---
name: cg-extract-types
description: Autonomously audits, extracts, and centralizes domain structs, raw generic instantiations, and Result wrappers into dedicated types.go files as single reusable named types across the codebase.
---

# Skill: Extracting Generic Types & Models to `types.go` (`cg-extract-types`)

This skill governs autonomous auditing, extraction, and centralization of scattered domain structs, raw generic Result envelopes (`ResultSlice[T]`, `ResultMap[K, V]`, `Result[T]`, `Wrap[T]`), and repeated generic type parameters into dedicated, package-level `types.go` files (or language-equivalent leaf `types/` modules) as single reusable named types.

## Mandatory Architectural Invariants

1. **Dedicated `types.go` Single Reusable Type Definition:**
   - Every package managing domain models, payloads, or Result envelopes MUST define them inside a dedicated `types.go` file within the package directory as a single reusable named type.
   - Never declare ad-hoc unexported structs or raw generic Result envelopes inline in implementation files (`importer.go`, `store.go`, `service.go`).
   - Sibling implementation files in the package import or consume the single named type directly without re-specifying generic parameters (`ResultSlice[...]`).

2. **Single Reusable Type Alias for Repeated Generics (Spec 27, Rule 2):**
   - When a generic type is used 3 or more times with the same type argument, or wraps a domain payload struct, create a canonical type alias in `types.go`.
   - Pattern:
     ```go
     type (
         ScheduleExportBundle struct { ... }
         ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]
         ScheduleExportBundleSingleResult = result.Wrap[ScheduleExportBundle]
     )
     ```

3. **Exported Models & Affirmative Field Naming:**
   - All domain payload structs used across function boundaries MUST be exported with PascalCase names.
   - All boolean fields in models, structs, and envelopes MUST carry affirmative prefixes (`is*` or `has*`):
     - `isDefined bool` (TOTAL BAN on bare `defined bool`)
     - `hasRecords bool` / `hasRecord bool`
     - `isStopped bool` (TOTAL BAN on bare `stop bool`)
     - `isForced bool` / `isPaused bool` / `isEnabled bool`

4. **Pointer-Attached Null Safety & Fluent Predicates:**
   - Result containers returned by functions (`Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`) must have inspection methods attached to pointer receivers.
   - Code consuming `types.go` aliases must utilize the 4 fluent pointer-safe predicates:
     - `res.IsCountOtherThan(expectedCount int) bool`: Returns `true` if operation failed or count != expected.
     - `res.IsEmpty() bool`: Returns `true` if empty collection, nil/zero payload, or nil receiver.
     - `res.HasRecord() bool` (or `res.HasRecords()`): Returns `true` if succeeded and count > 0.
     - `res.IsDefined() bool`: Returns `true` if succeeded and payload is non-null/non-empty with count > 0.
   - **Mandatory Replacement for `!isEmpty`:** NEVER use inverted negative empty checks (`!isEmpty`, `!res.IsEmpty()`). Always use affirmative `isDefined` / `res.IsDefined()` when asserting that data or records are present.

5. **Clean Package Boundaries & Zero Dependency Cycles:**
   - `types.go` files are leaf definitions within their package.
   - `types.go` MUST only import foundation utility packages (e.g. `pkg/result`, `pkg/appfault`) and standard library packages.
   - `types.go` MUST NEVER import heavy business logic, services, or circular sibling packages.

6. **Targeted Verification (No Full CI/CD Runner):**
   - Verify changes using targeted file-level linters (`python linter-scripts/check-enum-and-boolean.py <files>`, `python linter-scripts/check-function-lengths.py <files>`).
   - DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline turns.

7. **Zero Releases:**
   - Strictly forbidden from bumping version numbers, running release scripts, or cutting releases at the conclusion of this refactor.

---

## Code Patterns: Before vs After

### Pattern 1: Transitional Anti-Pattern (Unexported Inline Struct & Scattered Generics)

```go
// ❌ ANTI-PATTERN: unexported inline struct in importer.go and repetitive raw generic returns
// schedule/importer.go
package schedule

type scheduleExportBundle struct {
    ScheduleId   string `json:"scheduleId"`
    WorkflowName string `json:"workflowName"`
    Payload      []byte `json:"payload"`
}

func parseImportSQLite(filePath string) result.ResultSlice[scheduleExportBundle] {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return result.FailSlice[scheduleExportBundle](appfault.Wrap(err, "parseImportSQLite", appfault.ErrCodeDatabaseQueryFailure))
    }
    // ...
    return result.OkSlice(bundles)
}
```

### Pattern 2: Canonical `types.go` Centralization (Single Reusable Type)

```go
// ✅ REQUIRED: Centralized types.go with exported model and single reusable type alias
// schedule/types.go
package schedule

import (
    "coding-guidelines/common/pkg/appfault"
    "coding-guidelines/common/pkg/result"
)

type (
    // ScheduleExportBundle defines the exported schedule archive payload.
    ScheduleExportBundle struct {
        ScheduleId   string `json:"scheduleId"`
        WorkflowName string `json:"workflowName"`
        Payload      []byte `json:"payload"`
    }

    // ScheduleExportBundleResult is the canonical single reusable result envelope for bundle slices.
    ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]

    // ScheduleExportBundleSingleResult is the canonical single reusable result envelope for a single bundle.
    ScheduleExportBundleSingleResult = result.Wrap[ScheduleExportBundle]
)
```

```go
// ✅ REQUIRED: Clean implementation file consuming the single reusable type from types.go
// schedule/importer.go
package schedule

import (
    "os"
    "coding-guidelines/common/pkg/appfault"
    "coding-guidelines/common/pkg/result"
)

func parseImportSQLite(filePath string) ScheduleExportBundleResult {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return result.FailSlice[ScheduleExportBundle](appfault.Wrap(err, "parseImportSQLite", appfault.ErrCodeDatabaseQueryFailure))
    }
    // ...
    return result.OkSlice(bundles)
}
```

---

## Polyglot Equivalents

| Language | Leaf Type File Location | Model Definition | Single Reusable Result Alias |
|:---|:---|:---|:---|
| **Go** | `<pkg>/types.go` | `type ScheduleExportBundle struct { ... }` | `type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]` |
| **TypeScript** | `src/<module>/types.ts` | `export interface ScheduleExportBundle { ... }` | `export type ScheduleExportBundleResult = ResultSlice<ScheduleExportBundle>;` |
| **C#** | `<Module>/Types.cs` | `public record ScheduleExportBundle(...);` | `public sealed record ScheduleExportBundleResult : ResultSlice<ScheduleExportBundle>;` |
| **Rust** | `<module>/types.rs` | `pub struct ScheduleExportBundle { ... }` | `pub type ScheduleExportBundleResult = ResultSlice<ScheduleExportBundle>;` |

---

## Autonomous Execution Workflow

### Phase 1: Discovery & Ledger

1. Run ripgrep to inventory unexported domain structs declared in implementation files and raw generic returns:
   ```bash
   rg "type [a-z][A-Za-z0-9]+ struct \{" -g "*.go" -g "!*test*.go"
   rg "result\.ResultSlice\[" -g "*.go" -g "!*types.go"
   rg "result\.Wrap\[" -g "*.go" -g "!*types.go"
   rg "appfault\.ResultMap\[" -g "*.go" -g "!*types.go"
   ```
2. Build the violation ledger in `.lovable/plans/pending/XX-types-go-extraction-audit.md` grouping violations by package.
3. Decompose into granular subtasks in `.lovable/plans/subtasks/XX-types-go-extraction/`.

### Phase 2: Extraction & Verification

1. For each target package:
   - Create or open `<pkg>/types.go`.
   - Extract and export domain models.
   - Declare single reusable type aliases for Result envelopes.
   - Refactor implementation files to use the single type alias.
   - Update call sites and test assertions.
2. Verify affirmative boolean field naming (`isDefined bool`).
3. Run targeted file-level linters:
   ```bash
   python linter-scripts/check-enum-and-boolean.py <modified_files>
   python linter-scripts/check-relative-paths.py
   ```
4. Record modified files in `.lovable/temp/recent-file-changes.json`.

---

## Routine Execution Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes time. Verify code strictly using targeted file-level linters on the specific modified files.
