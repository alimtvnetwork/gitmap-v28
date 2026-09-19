# RCA 60: Pipeline History Import, Typecheck Fields, Nested If, and Error Management Linters

## 1. Symptom
GitHub Actions run `#35428297196` failed across 23 pipeline sections:
1. `cmdpipeline/pipeline_history.go:516:11: undefined: filepath` causing cascade package compilation failures across `cli`, `cmd`, `cmdagy`, `cmdssh`, etc.
2. `cmdssh/ssh_exec_command_test.go:174:38: undefined: strings`
3. `cmdssh/ssh_exec_ui_test.go:32:4: unknown field ID in struct literal of type db.SSHConnection`
4. `pipelinedb/pipeline_split_db_test.go:296:13: undefined: store`
5. `cmd/dbmigrate.go:93: undefined: store`
6. `cmd/env.go:23: undefined: helptext`
7. `cmd/auto_alias_test.go:128: undefined: store.OpenInMemory`
8. `examples_golden_test.go:42: help files missing '## Examples' section (agy.md, pipeline.md)`
9. Boolean linter: `shouldSkipFallback` banned prefix in `cli/cmdpipeline/pipeline_logs.go:131`, inverted boolean `!recSuccess.IsSuccess` in `cli/cmdpipeline/pipeline_recorder_test.go:34`.
10. Nested-if linter: 13 depth-2 nested if violations across 10 files (`pipeline_split_db_fs.go`, `install_aliases.go`, `pipeline.go`, `pipeline_db_ops.go`, `rm.go`, `reconcile_prompt.go`, `auto_alias.go`, `sshexec.go`, `ssh_transfer.go`, `ssh_exec_ui.go`).
11. Error management linter: 3 swallowed database errors (`_, _ = db.Exec(...)`) in `cli/cmd/storage_reset.go`.

## 2. Root Cause
- Missing package imports (`path/filepath`, `strings`, `store`, `helptext`) after recent helper extractions.
- Struct literal drift in unit test (`db.SSHConnection` has no `ID` field; primary key is `Alias`).
- Missing `OpenInMemory` in `cli/store`.
- Help text files used `### Examples` instead of required top-level `## Examples`.
- Unflattened conditional blocks creating depth-2 nested `if` statements.
- `_, _ = db.Exec(...)` calls violating zero-swallow error policy in error management guidelines.

## 3. Resolution
1. Added `"path/filepath"` to `cli/cmdpipeline/pipeline_history.go`.
2. Added `"strings"` to `cli/cmdssh/ssh_exec_command_test.go`.
3. Added `"github.com/alimtvnetwork/gitmap-v28/cli/store"` to `cli/cmd/dbmigrate.go` and `cli/pipelinedb/pipeline_split_db_test.go`.
4. Added `"github.com/alimtvnetwork/gitmap-v28/cli/helptext"` to `cli/cmd/env.go`.
5. Implemented `OpenInMemory` in `cli/store/store.go` for isolated memory DB test instances.
6. Fixed `db.SSHConnection` literal in `cli/cmdssh/ssh_exec_ui_test.go` removing nonexistent `ID` field.
7. Renamed `shouldSkipFallback` to `isSkipFallback`, replaced `!recSuccess.IsSuccess` with `recSuccess.IsFail()`.
8. Flattened all 13 nested-if statements using guard clauses, early returns, and dedicated helper functions.
9. Checked error returns on `db.Exec` in `cli/cmd/storage_reset.go`.
10. Added `## Examples` sections to `cli/helptext/pipeline.md` and `cli/helptext/agy.md`.

## 4. Prevention & Learnings
- Always verify all imported packages are explicitly declared when calling helpers across package boundaries.
- Run local policy and type checkers before release triggers to ensure zero linter regressions.
