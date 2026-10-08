# Subtask 06 — Wave 2 moves, group G2 (cmdcompletion→cmdmultigroup, NEW packages)

## Objective
Move implementation files out of `cli/cmd/` into NEW owning `cmdX` packages per the move map
(`.ai-memory/plans/subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md`,
Wave 3). MECHANICAL move wave: copy with package-clause swap, verify, delete source.
The lead rewires dispatch callers centrally afterwards — you do NOT touch dispatch.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section.

## Owned files
ONLY the source files listed below (in `cli/cmd/`) and their destinations. Do NOT touch any other file.

- cli/cmd/completion.go → cli/cmdcompletion/completion.go
- cli/cmd/completion_keys_printers.go → cli/cmdcompletion/completion_keys_printers.go
- cli/cmd/completion_printers.go → cli/cmdcompletion/completion_printers.go
- cli/cmd/console_other.go → cli/cmdconsole/console_other.go
- cli/cmd/console_windows.go → cli/cmdconsole/console_windows.go
- cli/cmd/copy_paste.go → cli/cmdcopy/copy_paste.go
- cli/cmd/create_cmd.go → cli/cmdcreate/create_cmd.go
- cli/cmd/create_ops.go → cli/cmdcreate/create_ops.go
- cli/cmd/dashboard.go → cli/cmddashboard/dashboard.go
- cli/cmd/diff.go → cli/cmddiff/diff.go
- cli/cmd/diff_help_menu.go → cli/cmddiff/diff_help_menu.go
- cli/cmd/dispatchcompare.go → cli/cmddiff/dispatchcompare.go
- cli/cmd/dispatchdiff.go → cli/cmddiff/dispatchdiff.go
- cli/cmd/diffprofiles.go → cli/cmddiffprofiles/diffprofiles.go
- cli/cmd/diffprofilesops.go → cli/cmddiffprofiles/diffprofilesops.go
- cli/cmd/diffprofilesrender.go → cli/cmddiffprofiles/diffprofilesrender.go
- cli/cmd/child_path.go → cli/cmddir/child_path.go
- cli/cmd/dir_scope_matcher.go → cli/cmddir/dir_scope_matcher.go
- cli/cmd/docs.go → cli/cmddocs/docs.go
- cli/cmd/env.go → cli/cmdenv/env.go
- cli/cmd/envops.go → cli/cmdenv/envops.go
- cli/cmd/envplatform_unix.go → cli/cmdenv/envplatform_unix.go
- cli/cmd/envplatform_windows.go → cli/cmdenv/envplatform_windows.go
- cli/cmd/envregistry.go → cli/cmdenv/envregistry.go
- cli/cmd/envvalidate.go → cli/cmdenv/envvalidate.go
- cli/cmd/exec.go → cli/cmdexec/exec.go
- cli/cmd/execprint.go → cli/cmdexec/execprint.go
- cli/cmd/explorer.go → cli/cmdexplorer/explorer.go
- cli/cmd/export.go → cli/cmdexport/export.go
- cli/cmd/exportrender.go → cli/cmdexport/exportrender.go
- cli/cmd/find_duplicates.go → cli/cmdfind/find_duplicates.go
- cli/cmd/find_duplicates_git.go → cli/cmdfind/find_duplicates_git.go
- cli/cmd/find_entry.go → cli/cmdfind/find_entry.go
- cli/cmd/find_files.go → cli/cmdfind/find_files.go
- cli/cmd/find_files_match.go → cli/cmdfind/find_files_match.go
- cli/cmd/find_help_menu.go → cli/cmdfind/find_help_menu.go
- cli/cmd/findnext.go → cli/cmdfindnext/findnext.go
- cli/cmd/findnextflags.go → cli/cmdfindnext/findnextflags.go
- cli/cmd/findnextrender.go → cli/cmdfindnext/findnextrender.go
- cli/cmd/fix_cmd.go → cli/cmdfix/fix_cmd.go
- cli/cmd/fix_diagnostics.go → cli/cmdfix/fix_diagnostics.go
- cli/cmd/fix_execute.go → cli/cmdfix/fix_execute.go
- cli/cmd/fix_ls.go → cli/cmdfix/fix_ls.go
- cli/cmd/fix_ls_table.go → cli/cmdfix/fix_ls_table.go
- cli/cmd/fixauth.go → cli/cmdfixauth/fixauth.go
- cli/cmd/fixcredential.go → cli/cmdfixauth/fixcredential.go
- cli/cmd/gitrm.go → cli/cmdgitrm/gitrm.go
- cli/cmd/gitrm_entry.go → cli/cmdgitrm/gitrm_entry.go
- cli/cmd/gomod.go → cli/cmdgomod/gomod.go
- cli/cmd/gomodbranch.go → cli/cmdgomod/gomodbranch.go
- cli/cmd/gomodreplace.go → cli/cmdgomod/gomodreplace.go
- cli/cmd/group.go → cli/cmdgroup/group.go
- cli/cmd/groupadd.go → cli/cmdgroup/groupadd.go
- cli/cmd/groupcreate.go → cli/cmdgroup/groupcreate.go
- cli/cmd/groupdelete.go → cli/cmdgroup/groupdelete.go
- cli/cmd/grouplist.go → cli/cmdgroup/grouplist.go
- cli/cmd/groupremove.go → cli/cmdgroup/groupremove.go
- cli/cmd/groupscoped.go → cli/cmdgroup/groupscoped.go
- cli/cmd/groupshow.go → cli/cmdgroup/groupshow.go
- cli/cmd/hasanyupdates.go → cli/cmdhaschange/hasanyupdates.go
- cli/cmd/haschange.go → cli/cmdhaschange/haschange.go
- cli/cmd/headtail.go → cli/cmdheadtail/headtail.go
- cli/cmd/command_history_cli.go → cli/cmdhistory/command_history_cli.go
- cli/cmd/history.go → cli/cmdhistory/history.go
- cli/cmd/history_purge_cmd.go → cli/cmdhistory/history_purge_cmd.go
- cli/cmd/historyrender.go → cli/cmdhistory/historyrender.go
- cli/cmd/historyreset.go → cli/cmdhistory/historyreset.go
- cli/cmd/historyrewrite.go → cli/cmdhistory/historyrewrite.go
- cli/cmd/historyrewrite_flags.go → cli/cmdhistory/historyrewrite_flags.go
- cli/cmd/historyrewrite_helpers.go → cli/cmdhistory/historyrewrite_helpers.go
- cli/cmd/historyrewrite_paths.go → cli/cmdhistory/historyrewrite_paths.go
- cli/cmd/historyrewrite_pin.go → cli/cmdhistory/historyrewrite_pin.go
- cli/cmd/historyrewrite_push.go → cli/cmdhistory/historyrewrite_push.go
- cli/cmd/historyrewrite_sandbox.go → cli/cmdhistory/historyrewrite_sandbox.go
- cli/cmd/historyrewrite_verify.go → cli/cmdhistory/historyrewrite_verify.go
- cli/cmd/hygiene_format.go → cli/cmdhygiene/hygiene_format.go
- cli/cmd/hygiene_parallel.go → cli/cmdhygiene/hygiene_parallel.go
- cli/cmd/hygiene_parallel_map.go → cli/cmdhygiene/hygiene_parallel_map.go
- cli/cmd/hygiene_parallel_workers.go → cli/cmdhygiene/hygiene_parallel_workers.go
- cli/cmd/import_all_json_cmd.go → cli/cmdimport/import_all_json_cmd.go
- cli/cmd/importcmd.go → cli/cmdimport/importcmd.go
- cli/cmd/index_cmd.go → cli/cmdindex/index_cmd.go
- cli/cmd/inject.go → cli/cmdinject/inject.go
- cli/cmd/inject_idempotency.go → cli/cmdinject/inject_idempotency.go
- cli/cmd/interactive.go → cli/cmdinteractive/interactive.go
- cli/cmd/ip_cmd.go → cli/cmdip/ip_cmd.go
- cli/cmd/ip_resolver.go → cli/cmdip/ip_resolver.go
- cli/cmd/ip_subcmds.go → cli/cmdip/ip_subcmds.go
- cli/cmd/ipchange_cmd.go → cli/cmdip/ipchange_cmd.go
- cli/cmd/join.go → cli/cmdjoin/join.go
- cli/cmd/jsonextract.go → cli/cmdjsonextract/jsonextract.go
- cli/cmd/latestbranch.go → cli/cmdlatestbranch/latestbranch.go
- cli/cmd/latestbranchcsv.go → cli/cmdlatestbranch/latestbranchcsv.go
- cli/cmd/latestbranchoutput.go → cli/cmdlatestbranch/latestbranchoutput.go
- cli/cmd/latestbranchrender.go → cli/cmdlatestbranch/latestbranchrender.go
- cli/cmd/latestbranchresolve.go → cli/cmdlatestbranch/latestbranchresolve.go
- cli/cmd/latestbranchswitch.go → cli/cmdlatestbranch/latestbranchswitch.go
- cli/cmd/addlfsinstall.go → cli/cmdlfscommon/addlfsinstall.go
- cli/cmd/lfscommon.go → cli/cmdlfscommon/lfscommon.go
- cli/cmd/list.go → cli/cmdlist/list.go
- cli/cmd/list_preview.go → cli/cmdlist/list_preview.go
- cli/cmd/list_tree.go → cli/cmdlist/list_tree.go
- cli/cmd/listreleases.go → cli/cmdlist/listreleases.go
- cli/cmd/listreleasesallrepos.go → cli/cmdlist/listreleasesallrepos.go
- cli/cmd/listreleasesload.go → cli/cmdlist/listreleasesload.go
- cli/cmd/listreleasesrender.go → cli/cmdlist/listreleasesrender.go
- cli/cmd/listversions.go → cli/cmdlist/listversions.go
- cli/cmd/listversionsrender.go → cli/cmdlist/listversionsrender.go
- cli/cmd/listversionsutil.go → cli/cmdlist/listversionsutil.go
- cli/cmd/llm_entry.go → cli/cmdllm/llm_entry.go
- cli/cmd/llm_help_menu.go → cli/cmdllm/llm_help_menu.go
- cli/cmd/llmdocs.go → cli/cmdllm/llmdocs.go
- cli/cmd/llmdocscommands.go → cli/cmdllm/llmdocscommands.go
- cli/cmd/llmdocsgroups.go → cli/cmdllm/llmdocsgroups.go
- cli/cmd/llmdocsheader.go → cli/cmdllm/llmdocsheader.go
- cli/cmd/llmdocsrender.go → cli/cmdllm/llmdocsrender.go
- cli/cmd/llmdocssections.go → cli/cmdllm/llmdocssections.go
- cli/cmd/locate_entry.go → cli/cmdlocate/locate_entry.go
- cli/cmd/dispatchmovemerge.go → cli/cmdmerge/dispatchmovemerge.go
- cli/cmd/merge.go → cli/cmdmerge/merge.go
- cli/cmd/merge_json_cmd.go → cli/cmdmerge/merge_json_cmd.go
- cli/cmd/move.go → cli/cmdmerge/move.go
- cli/cmd/movemergeflags.go → cli/cmdmerge/movemergeflags.go
- cli/cmd/migrate.go → cli/cmdmigrate/migrate.go
- cli/cmd/migrate_cmd.go → cli/cmdmigrate/migrate_cmd.go
- cli/cmd/migrate_wizard.go → cli/cmdmigrate/migrate_wizard.go
- cli/cmd/mkdir.go → cli/cmdmkdir/mkdir.go
- cli/cmd/multigroup.go → cli/cmdmultigroup/multigroup.go
- cli/cmd/multigroupops.go → cli/cmdmultigroup/multigroupops.go

## Method (per file, in order)
1. Verify the source exists at `cli/cmd/<file>`. If missing, record SKIP (do not guess, do not create).
2. Verify the destination does NOT already exist (`cli/<newpkg>/<file>`). If it exists, record BLOCKED (do not overwrite).
3. Read the source fully (via `gitmap cat`). Write the destination with EXACTLY one change: the `package cmd` line becomes `package <newpkg>`. Everything else byte-identical (keep build tags, comments, imports, blank lines).
4. Read back the destination's first 3 lines to confirm the package clause; confirm the file is non-empty.
5. Delete the source file. (Source removal is the authorized move — nothing beyond the listed files.)
6. Do NOT edit any `cli/cmd/` file. Do NOT edit dispatch tables. Do NOT create doc.go files.

## Done criteria
- [ ] Every listed file: moved (dest exists with new package clause) + source deleted, or recorded SKIP/BLOCKED with reason.
- [ ] No file outside the list touched.
- [ ] Report: moved count, skipped list, blocked list.

## Hard rules
- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on grep, rg, git grep, Select-String.
- NEVER run ANY git command. NEVER run `go build` / `go test` / `go vet` (lead verifies centrally).
- Relative paths only. All filenames already lowercase.
- Do NOT prune aliases. Do NOT touch versioning.
