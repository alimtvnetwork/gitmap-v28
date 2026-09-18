# 4-Part Root Cause Analysis (RCA): Pipeline CI Build Failures on Commit 5cb748b

## 1. Why It Happened (High-Level Business & Architectural Context)
In commit `5cb748b` (and preceding commit `d72decf5`), we implemented direct command routing for `gitmap fix agy`, compound `gitmap agy errors fix` normalization, `gitmap cluster init` generator, and enriched documentation for `gitmap servers-clients`. While local linters for nesting, booleans, and error management passed with zero violations, Go compiler checks and golden fixture tests failed in CI due to missing per-file package imports in `cli/cmd`, a slice value type mismatch in `runFix` parameter extraction, an unused standard library import in a test file, a non-standard heading name in markdown help documentation, and formatting whitespace drift.

---

## 2. How It Happened (Technical Execution Flow)
1. **Per-File Import Scope in Go**: In Go, import declarations are file-scoped, not package-scoped. Although `github.com/alimtvnetwork/gitmap-v28/cli/cmdagy` was imported in `cli/cmd/root.go` and `cli/cmd/find_duplicates.go`, it was omitted from the import blocks of `cli/cmd/clihelpers.go` and `cli/cmd/rootutility.go`. When the compiler parsed `clihelpers.go:723` (`cmdagy.RunPipelineFixAgyCLI`) and `rootutility.go:115` (`cmdagy.RunPipelineFixAgyCLI`), it failed with `undefined: cmdagy`.
2. **Slice Value Type Mismatch in `executeFixTarget`**: In `cli/cmd/fix_cmd.go`, `LoadRemediationState()` returns `[]RemediationItem` (a slice of structs). When refactoring `runFix` to extract `executeFixTarget`, the parameter was declared as `items []*RemediationItem` (slice of pointers), causing a type mismatch error when calling `executeFixTarget(args, aliasOverride, items)` and subsequently calling `resolveFixTarget(args, aliasOverride, items)`.
3. **Unused Import in Test**: In `cli/cmdssh/cluster_init_cmd_test.go`, `"os"` was imported on line 4 but never referenced, violating Go's strict unused import rule during `go vet` and compiler runs.
4. **Golden Fixture Help Text Regex Gate**: In `cli/helptext/servers-clients.md`, the examples section was authored under `## Concrete Examples`. The golden fixture test `TestEveryHelpFileHasExamples` in `cli/helptext/examples_golden_test.go` specifically searches for `\n## Examples` followed by fenced code blocks. Because of the word "Concrete", `strings.Index(md, "\n## Examples")` returned -1.
5. **Gofmt Whitespace Drift**: Minor spacing drift in `cli/cmdagy/agy_cmd.go` and `cli/cmdagy/agy_pipeline_fix_errors_test.go` caused `test_gofmt_check_clean_repo` in `test_ci_scripts.py` to fail.

---

## 3. Root Cause (Exact File, Line, and Mechanism)
- `cli/cmd/clihelpers.go:723:37`: Referenced identifier `cmdagy` without importing `"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"`.
- `cli/cmd/rootutility.go:115:104`: Referenced identifier `cmdagy` without importing `"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"`.
- `cli/cmd/fix_cmd.go:46`: Declared `items []*RemediationItem` instead of `items []RemediationItem`.
- `cli/cmdssh/cluster_init_cmd_test.go:4:2`: Imported `"os"` without using it.
- `cli/helptext/servers-clients.md:111`: Titled heading `## Concrete Examples` instead of `## Examples`.
- `cli/cmdagy/agy_cmd.go` & `cli/cmdagy/agy_pipeline_fix_errors_test.go`: Unformatted whitespace detected by `gofmt -l`.

---

## 4. Code Fix (Exact Changes)

### Fix 1: Add Missing `cmdagy` Import in `cli/cmd/clihelpers.go`
```go
import (
    ...
    "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
    "github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
    "github.com/alimtvnetwork/gitmap-v28/cli/cmdcg"
```

### Fix 2: Add Missing `cmdagy` Import in `cli/cmd/rootutility.go`
```go
import (
    ...
    "github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
    "github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
    "github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh"
```

### Fix 3: Fix `items []RemediationItem` in `cli/cmd/fix_cmd.go`
```go
func executeFixTarget(args []string, aliasOverride string, items []RemediationItem) error {
    item, action, err := resolveFixTarget(args, aliasOverride, items)
    if err != nil {
        return err
    }
    if item == nil {
        return nil
    }

    return applyFixRecipe(item, action)
}
```

### Fix 4: Remove Unused `"os"` Import in `cli/cmdssh/cluster_init_cmd_test.go`
```go
package cmdssh

import (
    "path/filepath"
    "testing"
)
```

### Fix 5: Rename Section Heading in `cli/helptext/servers-clients.md`
```markdown
---

## Examples

### 1. Bash & POSIX Shell Fan-Out
```

### Fix 6: Run `gofmt -w` on Unformatted Go Files
```powershell
gofmt -w cli/cmdagy/agy_cmd.go cli/cmdagy/agy_pipeline_fix_errors_test.go
```
