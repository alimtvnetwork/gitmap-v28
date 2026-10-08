# Subtask 05 — Wave 2 moves, group G1 (existing packages + cmdadd→cmdcommittransfer)

## Objective
Move implementation files out of `cli/cmd/` into their owning `cmdX` packages per the move map
(`.ai-memory/plans/subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md`,
Waves 2–3). This is a MECHANICAL move wave: copy with package-clause swap, verify, delete source.
The lead rewires dispatch callers centrally afterwards — you do NOT touch dispatch.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section.

## Owned files
ONLY the source files listed below (in `cli/cmd/`) and their destinations. Do NOT touch any other file —
especially NOT any `cli/cmd/` file outside this list (dispatch tables, `root*.go`, keep-list files are lead-owned).

### G1a — Wave-2 list (into EXISTING packages; keep filenames)
- cli/cmd/ai_memory_server.go → cli/cmdai/ai_memory_server.go
- cli/cmd/agm_cmd.go → cli/cmdagy/agm_cmd.go
- cli/cmd/apps.go → cli/cmdapps/apps.go
- cli/cmd/remotetransport.go → cli/cmdclone/remotetransport.go
- cli/cmd/shellhandoff.go → cli/cmdclone/shellhandoff.go
- cli/cmd/what_configs_cmd.go → cli/cmdconfig/what_configs_cmd.go
- cli/cmd/which_format_cmd.go → cli/cmdconfig/which_format_cmd.go
- cli/cmd/cmd_db.go → cli/cmddb/cmd_db.go
- cli/cmd/dbmigrate.go → cli/cmddb/dbmigrate.go
- cli/cmd/dbreset.go → cli/cmddb/dbreset.go
- cli/cmd/downloaderconfig.go → cli/cmddownload/downloaderconfig.go
- cli/cmd/hd_download.go → cli/cmddownload/hd_download.go
- cli/cmd/error_cmd.go → cli/cmderrors/error_cmd.go
- cli/cmd/error_export.go → cli/cmderrors/error_export.go
- cli/cmd/error_ls.go → cli/cmderrors/error_ls.go
- cli/cmd/error_schema.go → cli/cmderrors/error_schema.go
- cli/cmd/error_tracker.go → cli/cmderrors/error_tracker.go
- cli/cmd/error_warnings.go → cli/cmderrors/error_warnings.go
- cli/cmd/dispatch_folder_tree.go → cli/cmdfoldertree/dispatch_folder_tree.go
- cli/cmd/folder_entry.go → cli/cmdfoldertree/folder_entry.go
- cli/cmd/addignoreattrs.go → cli/cmdignore/addignoreattrs.go
- cli/cmd/ignore_entry.go → cli/cmdignore/ignore_entry.go
- cli/cmd/binarylocations.go → cli/cmdinstall/binarylocations.go
- cli/cmd/cat.go → cli/cmdmacro/cat.go
- cli/cmd/macro_root_dispatch.go → cli/cmdmacro/macro_root_dispatch.go
- cli/cmd/nodes_cmd.go → cli/cmdnodes/nodes_cmd.go
- cli/cmd/nodes_ping_cmd.go → cli/cmdnodes/nodes_ping_cmd.go
- cli/cmd/power.go → cli/cmdos/power.go
- cli/cmd/power_format.go → cli/cmdos/power_format.go
- cli/cmd/power_ops.go → cli/cmdos/power_ops.go
- cli/cmd/processattr_other.go → cli/cmdos/processattr_other.go
- cli/cmd/processattr_windows.go → cli/cmdos/processattr_windows.go
- cli/cmd/rm.go → cli/cmdrm/rm.go
- cli/cmd/rm_resilience.go → cli/cmdrm/rm_resilience.go
- cli/cmd/dedupe.go → cli/cmdscan/dedupe.go
- cli/cmd/sc_help_menu.go → cli/cmdschedule/sc_help_menu.go
- cli/cmd/space.go → cli/cmdspace/space.go
- cli/cmd/sync.go → cli/cmdsync/sync.go
- cli/cmd/sync_help_menu.go → cli/cmdsync/sync_help_menu.go
- cli/cmd/task.go → cli/cmdtask/task.go
- cli/cmd/taskfilter.go → cli/cmdtask/taskfilter.go
- cli/cmd/taskops.go → cli/cmdtask/taskops.go
- cli/cmd/tasks.go → cli/cmdtask/tasks.go
- cli/cmd/tasks_list.go → cli/cmdtask/tasks_list.go
- cli/cmd/tasks_ops.go → cli/cmdtask/tasks_ops.go
- cli/cmd/tasksync.go → cli/cmdtask/tasksync.go
- cli/cmd/ui_cmd.go → cli/cmdui/ui_cmd.go
- cli/cmd/ui_help_menu.go → cli/cmdui/ui_help_menu.go
- cli/cmd/selfupdate.go → cli/cmdupdate/selfupdate.go
- cli/cmd/nginx.go → cli/cmdvhost/nginx.go
- cli/cmd/nginx_ini.go → cli/cmdvhost/nginx_ini.go
- cli/cmd/nginx_list.go → cli/cmdvhost/nginx_list.go
- cli/cmd/githubdesktop.go → cli/cmdvscode/githubdesktop.go
- cli/cmd/githubdesktop_help_menu.go → cli/cmdvscode/githubdesktop_help_menu.go
- cli/cmd/githubdesktop_optimize.go → cli/cmdvscode/githubdesktop_optimize.go

### G1b — Wave-3 list, cmdadd→cmdcommittransfer (into NEW packages; create the package dir)
- cli/cmd/add_entry.go → cli/cmdadd/add_entry.go
- cli/cmd/amend.go → cli/cmdamend/amend.go
- cli/cmd/amendaudit.go → cli/cmdamend/amendaudit.go
- cli/cmd/amendauditrender.go → cli/cmdamend/amendauditrender.go
- cli/cmd/amendexec.go → cli/cmdamend/amendexec.go
- cli/cmd/amendexecprint.go → cli/cmdamend/amendexecprint.go
- cli/cmd/amendlist.go → cli/cmdamend/amendlist.go
- cli/cmd/amendlistrender.go → cli/cmdamend/amendlistrender.go
- cli/cmd/audit.go → cli/cmdaudit/audit.go
- cli/cmd/audit_db.go → cli/cmdaudit/audit_db.go
- cli/cmd/audit_finish.go → cli/cmdaudit/audit_finish.go
- cli/cmd/auditlegacy.go → cli/cmdauditlegacy/auditlegacy.go
- cli/cmd/auditlegacy_diff_render.go → cli/cmdauditlegacy/auditlegacy_diff_render.go
- cli/cmd/auditlegacy_diffs.go → cli/cmdauditlegacy/auditlegacy_diffs.go
- cli/cmd/auditlegacy_emit.go → cli/cmdauditlegacy/auditlegacy_emit.go
- cli/cmd/auditlegacy_parse.go → cli/cmdauditlegacy/auditlegacy_parse.go
- cli/cmd/auditlegacy_report.go → cli/cmdauditlegacy/auditlegacy_report.go
- cli/cmd/aum_help_menu.go → cli/cmdaum/aum_help_menu.go
- cli/cmd/backup.go → cli/cmdbackup/backup.go
- cli/cmd/backup_cloud.go → cli/cmdbackup/backup_cloud.go
- cli/cmd/backup_cloud_ops.go → cli/cmdbackup/backup_cloud_ops.go
- cli/cmd/bash_runner.go → cli/cmdbash/bash_runner.go
- cli/cmd/bookmark.go → cli/cmdbookmark/bookmark.go
- cli/cmd/bookmarklist.go → cli/cmdbookmark/bookmarklist.go
- cli/cmd/bookmarklistrender.go → cli/cmdbookmark/bookmarklistrender.go
- cli/cmd/bookmarkrun.go → cli/cmdbookmark/bookmarkrun.go
- cli/cmd/bookmarksave.go → cli/cmdbookmark/bookmarksave.go
- cli/cmd/branch.go → cli/cmdbranch/branch.go
- cli/cmd/branch_help_menu.go → cli/cmdbranch/branch_help_menu.go
- cli/cmd/browse.go → cli/cmdbrowse/browse.go
- cli/cmd/cd.go → cli/cmdcd/cd.go
- cli/cmd/cd_suggest.go → cli/cmdcd/cd_suggest.go
- cli/cmd/cd_workdir_resolver.go → cli/cmdcd/cd_workdir_resolver.go
- cli/cmd/cddefault.go → cli/cmdcd/cddefault.go
- cli/cmd/cdops.go → cli/cmdcd/cdops.go
- cli/cmd/cfrppriorversion.go → cli/cmdcfrppriorversion/cfrppriorversion.go
- cli/cmd/changelog.go → cli/cmdchangelog/changelog.go
- cli/cmd/changelog_regen.go → cli/cmdchangelog/changelog_regen.go
- cli/cmd/changeloggen.go → cli/cmdchangelog/changeloggen.go
- cli/cmd/changelogprint.go → cli/cmdchangelog/changelogprint.go
- cli/cmd/changelogwrap.go → cli/cmdchangelog/changelogwrap.go
- cli/cmd/clean_dev_entry.go → cli/cmdclean/clean_dev_entry.go
- cli/cmd/cluster.go → cli/cmdcluster/cluster.go
- cli/cmd/cluster_h_stubs.go → cli/cmdcluster/cluster_h_stubs.go
- cli/cmd/cluster_help_menu.go → cli/cmdcluster/cluster_help_menu.go
- cli/cmd/cluster_ops.go → cli/cmdcluster/cluster_ops.go
- cli/cmd/clustercommand.go → cli/cmdcluster/clustercommand.go
- cli/cmd/clusterflags.go → cli/cmdcluster/clusterflags.go
- cli/cmd/clusterstatus.go → cli/cmdcluster/clusterstatus.go
- cli/cmd/clustersubcmd.go → cli/cmdcluster/clustersubcmd.go
- cli/cmd/code.go → cli/cmdcode/code.go
- cli/cmd/codingguidelines.go → cli/cmdcodingguidelines/codingguidelines.go
- cli/cmd/codingguidelines_commit.go → cli/cmdcodingguidelines/codingguidelines_commit.go
- cli/cmd/codingguidelines_compat.go → cli/cmdcodingguidelines/codingguidelines_compat.go
- cli/cmd/commit_batch.go → cli/cmdcommit/commit_batch.go
- cli/cmd/commit_cmd.go → cli/cmdcommit/commit_cmd.go
- cli/cmd/commit_help_menu.go → cli/cmdcommit/commit_help_menu.go
- cli/cmd/commit_push.go → cli/cmdcommit/commit_push.go
- cli/cmd/commit_ui.go → cli/cmdcommit/commit_ui.go
- cli/cmd/commitin.go → cli/cmdcommitin/commitin.go
- cli/cmd/commitin_ui_server.go → cli/cmdcommitin/commitin_ui_server.go
- cli/cmd/commitpull.go → cli/cmdcommitpull/commitpull.go
- cli/cmd/commitpull_bootstrap.go → cli/cmdcommitpull/commitpull_bootstrap.go
- cli/cmd/committransfer.go → cli/cmdcommittransfer/committransfer.go
- cli/cmd/dispatchcommittransfer.go → cli/cmdcommittransfer/dispatchcommittransfer.go

## Method (per file, in order)
1. Verify the source exists at `cli/cmd/<file>`. If missing, record SKIP (do not guess, do not create).
2. Verify the destination does NOT already exist (`cli/<newpkg>/<file>`). If it exists, record BLOCKED (do not overwrite).
3. Read the source fully (via `gitmap cat`). Write the destination with EXACTLY one change: the `package cmd` line becomes `package <newpkg>`. Everything else byte-identical (keep build tags, comments, imports, blank lines).
4. Read back the destination's first 3 lines to confirm the package clause; confirm the file is non-empty.
5. Delete the source file. (Source removal is the authorized move — nothing beyond the listed files.)
6. Do NOT edit any `cli/cmd/` file. Do NOT edit dispatch tables. Do NOT create doc.go files.

## Done criteria
- [ ] Every listed file: moved (dest exists with new package clause) + source deleted, or recorded SKIP/BLOCKED with reason.
- [ ] No file outside the lists touched.
- [ ] Report: moved count, skipped list, blocked list.

## Hard rules
- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on grep, rg, git grep, Select-String.
- NEVER run ANY git command. NEVER run `go build` / `go test` / `go vet` (lead verifies centrally).
- Relative paths only. All filenames already lowercase.
- Do NOT prune aliases. Do NOT touch versioning.
