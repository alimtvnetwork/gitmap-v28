# Plan 187: Naming Conventions, Bare Ok Elimination & Boolean Prefixes

> **Task Origin & Objective**:
> - Autonomous execution under Coding Guidelines Version 2.2.0 (`cg-naming`, `cg-execute naming`).
> - Eliminate bare `ok` identifiers across type assertions, map lookups, channel receives, and comma-ok idioms.
> - Replace negative boolean variables (`hasNo*`, `hasNone`, `isNot*`, `isNotesMissing`, `isMissingExtractor`) with affirmative names and inverted guard conditions.
> - Eliminate awkward `isExists` identifiers (e.g. `isExistsError` -> `isAlreadyExistsError`).
> - Enforce mandatory affirmative boolean prefixes (`is*`, `has*` only) and discrete assertions.
> - Zero line-compression cheating, functions <= 8–15 lines, defer build verification strictly to final step.

---

## 1. Problem Analysis & Violation Ledger

| # | File Path | Violation Type | Description | Remediation Plan |
|:---:|---|---|---|---|
| 1 | `cli/cmd/remediation_state_ops.go:60` | `NEGATIVE_BOOLEAN` | `hasNone := len(remaining) == 0` | Replaced with affirmative `hasRemaining := len(remaining) > 0` and inverted guard `if !hasRemaining` |
| 2 | `cli/cmd/githubdesktop_test.go:13` | `NEGATIVE_BOOLEAN` | `hasNone := hasGHDesktopInstallFlag(...)` | Renamed to `hasOther` in test assertion |
| 3 | `cli/cmd/commitin/orchestrator/save_profile.go:62` | `AWKWARD_EXISTS` | `isExistsError(saveErr)` | Renamed to `isAlreadyExistsError(saveErr)` |
| 4 | `cli/store/purge_history_migrate.go:72` | `NEGATIVE_BOOLEAN` | `isNotesMissing := !db.columnExists(...)` | Replaced with affirmative `hasNotes := db.columnExists(...)` and `!hasNotes` |
| 5 | `cli/store/purge_history_migrate.go:77` | `NEGATIVE_BOOLEAN` | `isCommentsMissing := !db.columnExists(...)` | Replaced with affirmative `hasComments := db.columnExists(...)` and `!hasComments` |
| 6 | `cli/apperror/apperror_test.go:43` | `BARE_OK` | `val, ok := appErr.Ctx["file"]` | Replaced with `val, isFound := appErr.Ctx["file"]` |
| 7 | `cli/archive/extract.go:110` | `BARE_OK` | `extractor, ok := format.(archives.Extractor)` and `isMissingExtractor` | Replaced with `extractor, isExtractor := format.(archives.Extractor)` and guard `if !isExtractor` |
| 8 | `cli/clonefrom/execute.go:125` | `BARE_OK` | `detail, ok := prepareAndClone(cloneParams)` | Replaced with `detail, isCloned := prepareAndClone(...)` |
| 9 | `cli/clonefrom/execute.go:131` | `BARE_OK` | `coDetail, coOK := runPostCloneCheckout(...)` | Replaced with `coDetail, isCheckedOut := runPostCloneCheckout(...)` |
| 10 | `cli/clonefrom/execute.go:142` | `BARE_OK` | `detail, ok := prepareDestParent(params.AbsDest)` | Replaced with `detail, isParentPrepared := prepareDestParent(...)` |
| 11 | `cli/clonefrom/parse.go:112` | `BARE_OK` | `if d, ok := obj[constants.CSVColumnDepth].(float64); ok` | Replaced with `if d, isFloat := obj[constants.CSVColumnDepth].(float64); isFloat` |
| 12 | `cli/clonefrom/summary_scheme.go:54` | `BARE_OK` | `if hit, ok := matchKnownScheme(url); ok` | Replaced with `if hit, isKnown := matchKnownScheme(url); isKnown` |
| 13 | `cli/clonefrom/validate.go:162` | `BARE_OK` | `if i, ok := seen[key]; ok` | Replaced with `if i, isSeen := seen[key]; isSeen` |
| 14 | `cli/clonenow/parsetext.go:35` | `BARE_OK` | `row, ok := textRowFromLine(line)` | Replaced with `row, isRowValid := textRowFromLine(line)` |
| 15 | `cli/clonenow/parsetext.go:58` | `BARE_OK` | `url, dest, ok := extractCloneArgs(fields)` | Replaced with `url, dest, isExtracted := extractCloneArgs(fields)` |
| 16 | `cli/clonenow/parse.go:224` | `BARE_OK` | `i, ok := idx[name]` | Replaced with `i, isFound := idx[name]` |
| 17 | `cli/clonenow/parse.go:298` | `BARE_OK` | `if idx, ok := seen[key]; ok` | Replaced with `if idx, isSeen := seen[key]; isSeen` |

---

## 2. Verification & Quality Gates

- `go vet ./cmd/... ./store/... ./apperror/... ./archive/... ./clonefrom/... ./clonenow/...`: Clean exit code 0.
- `03-ai-scripts/26-go-code-formatter.py`: All 12 files formatted with gofmt.
- `03-ai-scripts/33-test-inventory-generator.py`: 3,536 tests indexed.
- `03-ai-scripts/08-naming-autofixer.py`: 2,807 files checked and 100% compliant with implicit boolean conventions.
