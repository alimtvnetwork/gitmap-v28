# Deep Codebase Raw Error Audit & AGY Command Suite (Completed)

> **Task Completion Reference:** Started via user prompt requesting exhaustive raw error preservation (zero error swallowing) and full Antigravity (AGY) command suite enhancements (`rm`, `rm-rejoin-read`/`rrr`, `rm-rejoin-pin-read`/`rrpr`, `pins` defaulting to `ls`, and table display contract with `CONV NAME`, `ID`, `PROJECT`).
> **Execution Lifecycle:** Completed in 5 atomic subtask phases with zero build/test interruptions and 100% targeted linter gate pass rate.

## Master Summary & Architecture
- **Raw Error Preservation (Zero Swallowing):** Audited and refactored over 140 error-handling call sites across all `cli/` packages (`cli/cmd/`, `cli/cluster/`, `cli/archive/`, `cli/clonenext/`, `cli/cmdchromeprofile/`, `cli/cmdclone/`, `cli/cmdfixrepo/`, `cli/cmdscan/`, `cli/cmdsetup/`, `cli/release/`, `cli/store/`, etc.). Replaced unlinked `apperror.NewSimple` and unlinked `fmt.Errorf` format strings with `%w` and `apperror.WrapSimple(err, ...)`. Replaced `cliexit.HandleError(nil, code)` inside caught error blocks with `cliexit.HandleError(err, code)` to ensure all error origins, stacks, and underlying causes are displayed and never swallowed.
- **AGY Command Suite & LS Table Contract:**
  - `gitmap agy ls`: table prominently displays `CONV NAME`, `ID`, and `PROJECT` as the leading columns, prioritizing pinned projects at the top with `📌 pinned` status.
  - `gitmap agy pins`: defaults to `ls` with pinned projects listed first, supporting `add`, `rm`, `edit`, `help`, and direct targets (`gitmap agy pins <target>`).
  - `gitmap agy rm <target>` (aliases `remove`, `delete`, `del`): supports 3-digit sequence (`001`), project ID, slug, comma-separated tokens, and `rm folder <dir>`; displays explicit disk preservation notice.
  - `gitmap agy rm help` & `gitmap agy rm-rejoin-read help`: dedicated user guidance explaining that files on disk are strictly preserved and only Antigravity workspace metadata is cleaned.
  - `gitmap agy rm-rejoin-read` (`rrr`): purges Gemini conversations, retains project files, records task history with undo guidance, re-enrolls project, runs read prompt, and renames the initial conversation with project name or slug.
  - `gitmap agy rm-rejoin-pin-read` (`rrpr`, `rrbr`): same workflow as `rrr` with automatic project pinning.

---

## Consolidated Subtasks

### Subtask [01]: Error Wrapping in Cluster and Core Command Operations
- Replaced `apperror.NewSimple` with `apperror.WrapSimple(err, constants.ErrClusterInvalidPassword)` in `cli/cluster/exec_lifecycle.go`.
- Replaced `errors.New` with `fmt.Errorf("cluster: %w", err)` in `cli/cmd/clusterflags.go`.
- Passed raw `err`, `runsRes.AppError()`, `nodesRes.AppError()`, `delErr`, and `statsErr` into `cliexit.HandleError` across all 18 call sites in `cli/cmd/cluster_ops.go`.
- Passed raw `err` into `cliexit.HandleError` across all 15 call sites in `cli/cmd/code.go`.
- Wrapped directory and repo lookup errors with `%w` in `cli/cmd/cmd_db.go`.

### Subtask [02]: Error Wrapping in General CLI Commands
- Replaced `cliexit.HandleError(nil, ...)` with `cliexit.HandleError(err, ...)` across `cli/cmd/` (`addignoreattrs.go`, `auditlegacy.go`, `backup.go`, `commitin.go`, `committransfer.go`, `completion.go`, `dbmigrate.go`, `dedupe.go`, `diff.go`, `export.go`, `findnext.go`, `groupscoped.go`, `haschange.go`, `history.go`, `historyrewrite_*.go`, `importcmd.go`, `inject_idempotency.go`, `merge.go`, `move.go`, `orphans.go`, `regoldens.go`, `releasealias.go`, `releasepull.go`, `replace.go`, `replaceversionrun.go`, `reverttxn.go`, `reverttxn_lastn.go`, `safety_snapshot.go`, `selfinstall.go`, `serve.go`, `sf.go`, `startup.go`, `startupadd.go`, `templatesdiff.go`, `user_cmd.go`, `visibility*.go`).
- Replaced unlinked `apperror.NewSimple` with `apperror.WrapSimple(err, ...)` in `cli/cmd/backup_cloud_ops.go`, `cli/cmd/tasks_ops.go`, `cli/cmd/grouplist.go`, `cli/cmd/groupshow.go`, `cli/cmd/fixauth.go`, `cli/cmd/helpdashboard.go`, `cli/cmd/power_ops.go`, `cli/cmd/releasealias_git.go`, `cli/cmd/release_tools.go`, and `cli/cmd/export.go`.
- Wrapped file path errors with `%w` in `cli/cmd/helpdashboard_extract.go` and `cli/cmd/replaceflags.go`.

### Subtask [03]: Error Wrapping in Clone, Scan, FixRepo, and Release
- Preserved raw errors in `cli/archive/source.go` and `cli/clonenext/github.go`.
- Passed raw errors into `cliexit.HandleError` across `cli/cmdchromeprofile/`, `cli/cmdclone/`, `cli/cmdfixrepo/`, `cli/cmdpull/`, `cli/cmdscan/`, `cli/cmdsetup/`.
- Wrapped config and tag errors with `%w` in `cli/cmdfixrepo/fixrepo_config.go`, `cli/cmdscan/rescansubtree.go`, `cli/release/changeloggen.go`, and `cli/release/selfrelease_resolve.go`.

### Subtask [04]: AGY Command Suite & LS Table Display Contract
- Formatted `gitmap agy ls` table with `SEQ`, `CONV NAME`, `ID`, `PROJECT`, `BRANCH`, `STATUS`, `PATH`.
- Prioritized pinned projects at table top with `📌 pinned` status.
- Implemented `gitmap agy pins` defaulting to `ls` and routing `add`, `rm`, `edit`, `help`, and direct targets.
- Implemented `gitmap agy rm help` & `gitmap agy rm-rejoin-read help` guidance.
- Maintained all `cli/cmdagy` files $\le 100$ lines and functions $\le 15$ lines.

### Subtask [05]: Quality Linting, Test Inventory Recording, and Consolidation
- Verified 0 nested-if violations (`check-nested-ifs.py`).
- Verified 0 boolean guideline violations (`check-boolean-guidelines.py`).
- Verified file sizes $\le 100$ lines for all touched/created files.
- Recorded modified files in `.ai-memory/temp/recent-file-changes.json`.
