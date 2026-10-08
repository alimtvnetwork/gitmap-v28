# Subtask 08 — Wave 2 moves, group G4 (cmdresolver→cmdct, NEW packages)

## Objective
Move implementation files out of `cli/cmd/` into NEW owning `cmdX` packages per the move map
(`.ai-memory/plans/subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md`,
Wave 3). MECHANICAL move wave: copy with package-clause swap, verify, delete source.
The lead rewires dispatch callers centrally afterwards — you do NOT touch dispatch.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section.

## Owned files
ONLY the source files listed below (in `cli/cmd/`) and their destinations. Do NOT touch any other file.

- cli/cmd/resolver.go → cli/cmdresolver/resolver.go
- cli/cmd/resolver_alias.go → cli/cmdresolver/resolver_alias.go
- cli/cmd/resolver_glob.go → cli/cmdresolver/resolver_glob.go
- cli/cmd/resolver_path.go → cli/cmdresolver/resolver_path.go
- cli/cmd/resolver_pwd.go → cli/cmdresolver/resolver_pwd.go
- cli/cmd/resolver_types.go → cli/cmdresolver/resolver_types.go
- cli/cmd/rest_enable.go → cli/cmdrest/rest_enable.go
- cli/cmd/revert.go → cli/cmdrevert/revert.go
- cli/cmd/revertscript.go → cli/cmdrevert/revertscript.go
- cli/cmd/reverttxn.go → cli/cmdrevert/reverttxn.go
- cli/cmd/reverttxn_lastn.go → cli/cmdrevert/reverttxn_lastn.go
- cli/cmd/safe_rm_cmd.go → cli/cmdsafe/safe_rm_cmd.go
- cli/cmd/safety_snapshot.go → cli/cmdsafety/safety_snapshot.go
- cli/cmd/search.go → cli/cmdsearch/search.go
- cli/cmd/search_entry.go → cli/cmdsearch/search_entry.go
- cli/cmd/search_help_menu.go → cli/cmdsearch/search_help_menu.go
- cli/cmd/selfinstall.go → cli/cmdselfinstall/selfinstall.go
- cli/cmd/selfuninstall.go → cli/cmdselfinstall/selfuninstall.go
- cli/cmd/selfuninstallhandoff.go → cli/cmdselfinstall/selfuninstallhandoff.go
- cli/cmd/selfuninstallparts.go → cli/cmdselfinstall/selfuninstallparts.go
- cli/cmd/sends_cmd.go → cli/cmdsends/sends_cmd.go
- cli/cmd/sequence_cmd.go → cli/cmdsequence/sequence_cmd.go
- cli/cmd/seowrite.go → cli/cmdseowrite/seowrite.go
- cli/cmd/seowritecreate.go → cli/cmdseowrite/seowritecreate.go
- cli/cmd/seowritecsv.go → cli/cmdseowrite/seowritecsv.go
- cli/cmd/seowritegit.go → cli/cmdseowrite/seowritegit.go
- cli/cmd/seowriteloop.go → cli/cmdseowrite/seowriteloop.go
- cli/cmd/seowritetemplate.go → cli/cmdseowrite/seowritetemplate.go
- cli/cmd/serve.go → cli/cmdservercmd/serve.go
- cli/cmd/servercmd.go → cli/cmdservercmd/servercmd.go
- cli/cmd/servercmd_run.go → cli/cmdservercmd/servercmd_run.go
- cli/cmd/setsourcerepo.go → cli/cmdsetsourcerepo/setsourcerepo.go
- cli/cmd/sf.go → cli/cmdsf/sf.go
- cli/cmd/size.go → cli/cmdsize/size.go
- cli/cmd/stale.go → cli/cmdstale/stale.go
- cli/cmd/stash_help_menu.go → cli/cmdstash/stash_help_menu.go
- cli/cmd/startup.go → cli/cmdstartup/startup.go
- cli/cmd/startup_cmd.go → cli/cmdstartup/startup_cmd.go
- cli/cmd/startup_run.go → cli/cmdstartup/startup_run.go
- cli/cmd/startupadd.go → cli/cmdstartup/startupadd.go
- cli/cmd/startuplistfilter.go → cli/cmdstartup/startuplistfilter.go
- cli/cmd/startuplistrender.go → cli/cmdstartup/startuplistrender.go
- cli/cmd/startupstatusjson.go → cli/cmdstartup/startupstatusjson.go
- cli/cmd/stats.go → cli/cmdstats/stats.go
- cli/cmd/statsrender.go → cli/cmdstats/statsrender.go
- cli/cmd/status.go → cli/cmdstatus/status.go
- cli/cmd/status_branch.go → cli/cmdstatus/status_branch.go
- cli/cmd/status_help_menu.go → cli/cmdstatus/status_help_menu.go
- cli/cmd/status_id.go → cli/cmdstatus/status_id.go
- cli/cmd/status_md.go → cli/cmdstatus/status_md.go
- cli/cmd/status_pr.go → cli/cmdstatus/status_pr.go
- cli/cmd/storage_clean.go → cli/cmdstorage/storage_clean.go
- cli/cmd/storage_cmd.go → cli/cmdstorage/storage_cmd.go
- cli/cmd/storage_display.go → cli/cmdstorage/storage_display.go
- cli/cmd/storage_drive_other.go → cli/cmdstorage/storage_drive_other.go
- cli/cmd/storage_drive_windows.go → cli/cmdstorage/storage_drive_windows.go
- cli/cmd/storage_help_menu.go → cli/cmdstorage/storage_help_menu.go
- cli/cmd/storage_ls.go → cli/cmdstorage/storage_ls.go
- cli/cmd/storage_reset.go → cli/cmdstorage/storage_reset.go
- cli/cmd/storage_restore.go → cli/cmdstorage/storage_restore.go
- cli/cmd/storage_types.go → cli/cmdstorage/storage_types.go
- cli/cmd/templates_help_menu.go → cli/cmdtemplates/templates_help_menu.go
- cli/cmd/templates_state_cli.go → cli/cmdtemplates/templates_state_cli.go
- cli/cmd/templates_ui_server.go → cli/cmdtemplates/templates_ui_server.go
- cli/cmd/templatescli.go → cli/cmdtemplates/templatescli.go
- cli/cmd/templatesdiff.go → cli/cmdtemplates/templatesdiff.go
- cli/cmd/templatesinit.go → cli/cmdtemplates/templatesinit.go
- cli/cmd/temprelease.go → cli/cmdtemprelease/temprelease.go
- cli/cmd/tempreleaseexport_test_helpers.go → cli/cmdtemprelease/tempreleaseexport_test_helpers.go
- cli/cmd/tempreleaselist.go → cli/cmdtemprelease/tempreleaselist.go
- cli/cmd/tempreleaselistrender.go → cli/cmdtemprelease/tempreleaselistrender.go
- cli/cmd/tempreleaseops.go → cli/cmdtemprelease/tempreleaseops.go
- cli/cmd/tempreleaseremove.go → cli/cmdtemprelease/tempreleaseremove.go
- cli/cmd/undo.go → cli/cmdundo/undo.go
- cli/cmd/user_cmd.go → cli/cmduser/user_cmd.go
- cli/cmd/user_help_menu.go → cli/cmduser/user_help_menu.go
- cli/cmd/user_info.go → cli/cmduser/user_info.go
- cli/cmd/user_profiles_cmd.go → cli/cmduser/user_profiles_cmd.go
- cli/cmd/user_project_cmd.go → cli/cmduser/user_project_cmd.go
- cli/cmd/variable_cmd.go → cli/cmdvariable/variable_cmd.go
- cli/cmd/version_tags.go → cli/cmdversion/version_tags.go
- cli/cmd/versionhistory.go → cli/cmdversionhistory/versionhistory.go
- cli/cmd/versionhistoryrender.go → cli/cmdversionhistory/versionhistoryrender.go
- cli/cmd/visibility.go → cli/cmdvisibility/visibility.go
- cli/cmd/visibility_alias.go → cli/cmdvisibility/visibility_alias.go
- cli/cmd/visibility_branch.go → cli/cmdvisibility/visibility_branch.go
- cli/cmd/visibility_cmd.go → cli/cmdvisibility/visibility_cmd.go
- cli/cmd/visibility_config.go → cli/cmdvisibility/visibility_config.go
- cli/cmd/visibility_csv.go → cli/cmdvisibility/visibility_csv.go
- cli/cmd/visibility_db.go → cli/cmdvisibility/visibility_db.go
- cli/cmd/visibility_delete.go → cli/cmdvisibility/visibility_delete.go
- cli/cmd/visibility_diff.go → cli/cmdvisibility/visibility_diff.go
- cli/cmd/visibility_export.go → cli/cmdvisibility/visibility_export.go
- cli/cmd/visibility_flags.go → cli/cmdvisibility/visibility_flags.go
- cli/cmd/visibility_help_menu.go → cli/cmdvisibility/visibility_help_menu.go
- cli/cmd/visibility_list.go → cli/cmdvisibility/visibility_list.go
- cli/cmd/visibility_pr.go → cli/cmdvisibility/visibility_pr.go
- cli/cmd/visibility_repos.go → cli/cmdvisibility/visibility_repos.go
- cli/cmd/visibility_resolve.go → cli/cmdvisibility/visibility_resolve.go
- cli/cmd/visibility_run.go → cli/cmdvisibility/visibility_run.go
- cli/cmd/visibility_scan.go → cli/cmdvisibility/visibility_scan.go
- cli/cmd/visibility_status.go → cli/cmdvisibility/visibility_status.go
- cli/cmd/visibility_sync.go → cli/cmdvisibility/visibility_sync.go
- cli/cmd/visibility_table.go → cli/cmdvisibility/visibility_table.go
- cli/cmd/visibility_types.go → cli/cmdvisibility/visibility_types.go
- cli/cmd/visibility_undo.go → cli/cmdvisibility/visibility_undo.go
- cli/cmd/watch.go → cli/cmdwatch/watch.go
- cli/cmd/watchformat.go → cli/cmdwatch/watchformat.go
- cli/cmd/watchops.go → cli/cmdwatch/watchops.go
- cli/cmd/watchrender.go → cli/cmdwatch/watchrender.go
- cli/cmd/whoami.go → cli/cmdwhoami/whoami.go
- cli/cmd/bash_runner.go → cli/cmdbash/bash_runner.go
- cli/cmd/powershell_runner.go → cli/cmdpowershell/powershell_runner.go
- cli/cmd/author_sponsor.go → cli/cmdauthor/author_sponsor.go
- cli/cmd/as.go → cli/cmdas/as.go
- cli/cmd/asops.go → cli/cmdas/asops.go
- cli/cmd/append.go → cli/cmdappend/append.go
- cli/cmd/desktopsync.go → cli/cmddesktopsync/desktopsync.go
- cli/cmd/desktopsync_ops.go → cli/cmddesktopsync/desktopsync_ops.go
- cli/cmd/ct.go → cli/cmdct/ct.go

NOTE: `cli/cmd/visibility*.go` — the move map says "all 23 visibility*.go files". The 23 names above are the
concrete inventory from the map's era; if any are missing or extra visibility*.go files exist in cli/cmd/,
move the ones that exist, record SKIP for missing ones, and NEVER move a file not matching visibility*.go.
NOTE: `cli/cmd/bash_runner.go → cli/cmdbash/` appears twice in the source move map (dedupe: move it once; it is
listed here in G4 — do not also expect it elsewhere).

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
