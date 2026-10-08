# Subtask 04 — Enforcement, cmd split, enums, stale docs

## Objective
Per `02-spec/21-app/243-cli-help-displayer-and-backup-branch/04-enforcement-and-refactor.md`.

## Files (disjoint box — coordinate the move map with the coordinator first)
- The 10 files >1000 lines (split targets; one file per wave).
- `03-ai-scripts/` (file-size check script).
- `cli/cmd/` (dispatch only afterwards) + `cmdX/` packages (receivers).
- `cli/enums/` (new package).
- `.ai-memory/overview.md`, `docs/architecture/state-ownership.md` (new).

## Done criteria
- [ ] The 10 files >1000 lines split (junk-drawer clusters → own files); then the 500–1000 band, largest first.
- [ ] File-size check script (`gitmap py`) fails past 500 lines; wired into the repo's pre-commit path.
- [ ] `cli/cmd/` keeps only thin dispatch (`root*.go`, argv preprocessing, alias context); impl files moved to their `cmdX` packages; exact move map documented in the ledger BEFORE moving; `go build` verified after each wave.
- [ ] `cli/enums/` created per the Go coding-guideline enum pattern (cite the pattern); most self-contained group migrated first.
- [ ] `.ai-memory/overview.md` version corrected (from `version.json`); state-ownership doc written (`.gitmap/` vs `gitmap.json` vs SQLite).
- [ ] New-command checklist updated to the spec-03 help technique.
- [ ] `go build` + `go vet` clean after every wave; atomic commits per wave (`gitmap cpf`/`cpb`) + immediate push.
- [ ] Do NOT touch version cadence. Do NOT prune aliases.

## CMD move map

> Documented BEFORE any file move (spec 04 §3, coordinator protocol 2026-10-09).
> `cli/cmd/` holds 509 non-test, non-`root*.go` implementation files.
> Convention: moved files keep their filenames; no destination collisions
> (verified 2026-10-09: zero filename overlaps between `cli/cmd/` and any `cmdX/`).
> Collision check vs Worker D-1's 83 new files (cmdpull, cmdpipeline, dbengine,
> cmdupdate, cmdchromeprofile, cmdui, pipelinedb): Wave-1 new files land in
> cmdai/cmdide/cmdssh only — no overlap. Every future wave must re-check
> destination filenames before moving.

### Wave 1 (COMPLETED 2026-10-09 — Worker D-2; lead to verify with go build)

**A. `cli/cmd/clihelpers.go` (1129 lines) → DELETED, split by concern.**
Every identifier was verified with `gitmap aum search` before rewiring.
Pure one-line delegates were deleted and callers rewired to the existing
`cmdX` export; logic-bearing helpers moved to the owning package.

| Concern | Action |
|---|---|
| macro (`runCatCmd`, `runTouchCmd`, `runMkfileCmd`) | callers → `cmdmacro.RunCatCmd/RunTouchCmd/RunMkfileCmd` |
| vscode (`runVSCode`, `runVSCodePMSync`, `runVSCodePMPath`, `runVSCodeWorkspace`, `runFindDuplicatesVSCode`→`RunFindDuplicates`, `syncRecordsToVSCodePM`, `reportVSCodePMSoftError`, `renameVSCodePMByPath`, `runGitHubDesktopGroup`, `stripVSCodeSyncDisabledFlag`, `stripVSCodeTagFlags`) | callers → `cmdvscode.*` |
| cursor (`runCursor`) | → `cmdcursor.RunCursor` |
| ide (`syncScanRecordsToIDEs`, `applyIDEScanTargetFilter`, `applyIDEScanExcludeFilter`) | MOVED to `cli/cmdide/ide_scan_sync.go` as `SyncScanRecordsToIDEs` (+2 unexported helpers); `cmdscan.SyncRecordsToIDEsFn` hook rewired |
| ai (`stripAiFlag`) | MOVED to `cli/cmdai/ai_flag.go` as `StripAiFlag`; `root.go` argv preprocessing rewired |
| vhost (`runVHost`, `runNginxRm`→`RunVHostRm`, `runVHostEnable/Disable/Create/Test/Reload/List`) | callers → `cmdvhost.*`; dead `VHostConfig` alias, `VHostSiteType*` consts, `RenderVHostConfig` stub deleted (no callers) |
| zip (`runZip`, `runUnzipCompact`, `runZipGroup`) | → `cmdzip.*` |
| fixrepo (`runFixRepo`; `resolveFixRepoIdentity`/`fixRepoIdentity`/`copyFileForBackup`/`fixRepoBackupManifest`) | `runFixRepo` → `cmdfixrepo.RunFixRepo`; `undo.go` rewired to `cmdfixrepo.ResolveFixRepoIdentity()` (exported fields), `cmdfixrepo.CopyFileForBackup`, `cmdfixrepo.FixRepoBackupManifest`; var aliases (`rewriteFixRepoFile`, `CountUnguardedTokenHits`, `ScanUnguardedTokenHits`) → `cmdfixrepo.*` (incl. test file) |
| db (`runDB`, `runStartFresh`; `formatBytes`, `confirmOrSkip`, `isInteractiveStdin`, `hasConfirmFlag`, `parseConfirmFlag`, `truncateStr`) | callers → `cmddb.*` (12 files) |
| pipeline (`runPipeline`, `runPipelineErrors`, `runPipelineDetails`, `runPipelineAI`; `handlePipelineDB`, `resolveTempDir`, `buildErrorLogsPayload`, `FailedJobItem`, `recordRunInSplitDb`, `saveParsedFailedJobs`, `ghRunItem` + all type/var aliases) | callers → `cmdpipeline.*` (rootutility.go, cmd_invocation_test.go, mv_rm_tx_test.go, pipeline_and_repo_db_test.go); dead aliases deleted |
| fixgit (`runFixGit`, `RemediateGitIndex`, `FixGitOptions`, `FixGitIssue`) | `runFixGit` → `cmdfixgit.RunFixGit`; rest had no callers → deleted |
| cg (`runCG`, `parseCGFlags`; `CGOptions`, `CGMetadata`, `VersionInstallConfig`, `WriteCGMetadata`, `ReadCGMetadata`, `ResolveCGTarget`, `DefaultVersionInstallConfig`, `InstallVersionJSON`) | `runCG` → `cmdcg.RunCG`; `parseCGFlags` → `cmdcg.ParseCGFlags` (incl. test); rest had no callers → deleted |
| ssh (~25 ids) | delegates → `cmdssh.*`; `runRemote` MOVED to `cli/cmdssh/remote_dispatch.go` as `RunRemote`; dead `SSHJoinCmd`/`SJ*Cmd`/`SEOptions`/`SJOptions`/`SSHExecutor`/`RunSSHLogin`/`ParseMultiIPList` deleted (no callers) |
| installer (`RunInstallerCLI`; `ExportInstallerFlags`, `InstallerTreeNode`) | `RunInstallerCLI` → `cmdinstaller.RunInstallerCLI`; aliases had no callers → deleted |
| chrome (`runChrome`, `findChromeBinaryPath`, `readChromeBackup`) | → `cmdchrome.*` |
| setup (`runSetup`, `runSetupPerms`, `warnIfNoWrapper`; `resolveSetupConfigPath`, `isWrapperActive`) | → `cmdsetup.*`; last two had no callers → deleted |
| install (17 ids incl. `installOptions`/`ProfileComposition` aliases) | → `cmdinstall.*`; `isCustomStandaloneTool` → `cmdinstall.CheckCustomStandaloneTool`; `init()` DI hook closures rewired to `cmdinstall.*` |
| update (`runUpdateRunner`, `runUpdateCleanup`, `scheduleDeployedCleanupHandoff`, `initRunnerVerbose`, `expandTilde`, `createHandoffCopy`, `hasFlag`, `handleHandoffError`, `writeScriptToTemp`, `normalizeRepoPath`, `saveRepoPathToDB`) | → `cmdupdate.*`; dead `runUpdate` stub deleted (no callers) |
| pull (`runPullAll`, `runPullAllEfficient`, `runPullReleaseCD`, `runPushFix`→`cmdpushfix`, `findBySlug`, `isGitRepoCWD`, `pullOneRepo`, `ResolvePullDirectoryTargets`; `ExtractTransportFlags` dead → deleted) | → `cmdpull.*` / `cmdpushfix.RunPushFix`; `runPull`/`dispatchPullAllWithArgs`/`runPush` dispatch logic KEPT in `cmd/dispatch_pull.go` (thin dispatch) |
| clone (`runClone`, `runMultiClone`→`RunMultiCloneCommand`, `runCloneOnlyMissing`, `runCloneFixRepo`, `runCloneFixRepoPub`, `runCloneFrom`, `runCloneNext`, `runCloneNow`, `runClonePick`, `runCloneSync`, `ConvertURLToSSH/HTTPS`, `repoNameFromURL`, `extractRepoName`, `resolveCloneNextFolder`, `openInVSCode`, `registerSingleDesktop`; `ResolveCloneFixRepoName`, `RunRepoReclone` dead → deleted) | → `cmdclone.*` |
| scan (`runRescan`, `runRescanSubtree`, `expandHome`, `resolveOutFile`) | → `cmdscan.*`; `runScan` (pre-check logic) KEPT in `cmd/dispatch_scan.go` (thin dispatch) |
| doctor (`runDoctor`→`RunDoctorCmd`, `runCleanCorrupted`→`RunCleanCorrupted`, `CleanCorruptedDirs`, `CleanOptions`; `CleanResult` dead → deleted) | → `cmddoctor.*` |
| os (`runOS`, `runOSFixLink`, `RunOSCLI`, `RunOSInfoCLI`) | → `cmdos.*` |
| schedule/config/workdir/cargo (`runSchedule`, `runExportConfig`, `runImportConfig`, `runWorkDir`, `runCargo`) | → `cmdschedule` / `cmdconfig` / `cmdworkdir` / `cmdcargo` |
| argv helpers (`hasArgFlag`, `extractFlagVal`, `printJSON`, `isHelpFlag`, `isFileFlagWithArg`, `isTerminalInput`, `isExistingFile`) | KEPT in `cmd/clihelpers_argv.go` (argv preprocessing — explicitly kept per spec) |
| `init()` cross-package DI wiring | KEPT in `cmd/di_hooks.go` (dispatch glue — must stay in `cmd/`) |
| `runIP` (adapter over cmd-local `runIPCmd`) | KEPT in `cmd/dispatch_ip.go` until `ip_cmd.go` moves to `cmdip` (later wave) |

**B. Wave-1 impl tranche (most self-contained files):**

| File | Destination | Notes |
|---|---|---|
| `cli/cmd/muse_cmd.go` | `cli/cmdmuse/muse_cmd.go` (NEW pkg) | 125 lines, self-contained; 1 caller (`roottooling.go:159` → `cmdmuse.RunMuseCLI`) |
| `cli/cmd/wpr_cmd.go` | DELETED | pure delegate `RunWPR` → `cmdagy.RunWPRCLI`; rewire `roottooling.go` |
| `cli/cmd/status_target_resolver.go` | DELETED | pure delegate → `cmdpull.ResolvePullDirectoryTargets`; already rewired |

### Waves 2–4 (planned; lead validates each file before its wave)

**Wave 2 — moves into EXISTING `cmdX` packages** (mechanical; no new packages):

- `cmdai`: `ai_memory_server.go`
- `cmdagy`: `agm_cmd.go`
- `cmdapps`: `apps.go`
- `cmdaum`: (none — `aum_help_menu.go` → `cmdaum` is NEW, see wave 3)
- `cmdcargo`: (done wave 1)
- `cmdcg`: (done wave 1)
- `cmdchrome`: (done wave 1)
- `cmdchromeprofile`: (none from cmd/)
- `cmdclone`: `remotetransport.go`, `shellhandoff.go`
- `cmdconfig`: `what_configs_cmd.go`, `which_format_cmd.go`
- `cmddb`: `cmd_db.go`, `dbmigrate.go`, `dbreset.go`
- `cmddoctor`: (done wave 1)
- `cmddownload`: `downloaderconfig.go`, `hd_download.go`
- `cmderrors`: `error_cmd.go`, `error_export.go`, `error_ls.go`, `error_schema.go`, `error_tracker.go`, `error_warnings.go`
- `cmdfixgit`: (done wave 1)
- `cmdfixrepo`: (done wave 1)
- `cmdfoldertree`: `dispatch_folder_tree.go`, `folder_entry.go`
- `cmdgit`: (none — `githubdesktop*.go` → `cmdvscode`, see below)
- `cmdide`: (done wave 1)
- `cmdignore`: `addignoreattrs.go`, `ignore_entry.go`
- `cmdinstall`: `binarylocations.go`
- `cmdinstaller`: (done wave 1)
- `cmdmacro`: `cat.go`, `macro_root_dispatch.go`
- `cmdnodes`: `nodes_cmd.go`, `nodes_ping_cmd.go`
- `cmdos`: `power.go`, `power_format.go`, `power_ops.go`, `processattr_other.go`, `processattr_windows.go`
- `cmdpipeline`: (done wave 1)
- `cmdpull`: (done wave 1)
- `cmdpushfix`: (done wave 1)
- `cmdpy`: (none)
- `cmdrm`: `rm.go`, `rm_resilience.go`
- `cmdscan`: `dedupe.go`
- `cmdschedule`: `sc_help_menu.go`
- `cmdsee`: (none)
- `cmdservice`: (none)
- `cmdsetup`: (done wave 1)
- `cmdspace`: `space.go`
- `cmdssh`: (done wave 1)
- `cmdsync`: `sync.go`, `sync_help_menu.go`
- `cmdtask`: `task.go`, `taskfilter.go`, `taskops.go`, `tasks.go`, `tasks_list.go`, `tasks_ops.go`, `tasksync.go`
- `cmdtoken`: (none)
- `cmdui`: `ui_cmd.go`, `ui_help_menu.go`
- `cmdupdate`: `selfupdate.go`
- `cmdvhost`: `nginx.go`, `nginx_ini.go`, `nginx_list.go`
- `cmdvscode`: `githubdesktop.go`, `githubdesktop_help_menu.go`, `githubdesktop_optimize.go`
- `cmdworkdir`: (done wave 1)
- `cmdzip`: (done wave 1)

**Wave 3 — moves into NEW `cmdX` packages** (one package per command family):

- `cmdact`/`cmdadd`: `add_entry.go` → `cmdadd`
- `cmdamend`: `amend.go`, `amendaudit.go`, `amendauditrender.go`, `amendexec.go`, `amendexecprint.go`, `amendlist.go`, `amendlistrender.go`
- `cmdaudit`: `audit.go`, `audit_db.go`, `audit_finish.go`
- `cmdauditlegacy`: `auditlegacy.go`, `auditlegacy_diff_render.go`, `auditlegacy_diffs.go`, `auditlegacy_emit.go`, `auditlegacy_parse.go`, `auditlegacy_report.go`
- `cmdaum`: `aum_help_menu.go`
- `cmdbackup`: `backup.go`, `backup_cloud.go`, `backup_cloud_ops.go`
- `cmdbash`: `bash_runner.go`
- `cmdbookmark`: `bookmark.go`, `bookmarklist.go`, `bookmarklistrender.go`, `bookmarkrun.go`, `bookmarksave.go`
- `cmdbranch`: `branch.go`, `branch_help_menu.go`
- `cmdbrowse`: `browse.go`
- `cmdcd`: `cd.go`, `cd_suggest.go`, `cd_workdir_resolver.go`, `cddefault.go`, `cdops.go`
- `cmdcfrppriorversion`: `cfrppriorversion.go`
- `cmdchangelog`: `changelog.go`, `changelog_regen.go`, `changeloggen.go`, `changelogprint.go`, `changelogwrap.go`
- `cmdclean`: `clean_dev_entry.go`
- `cmdcluster`: `cluster.go`, `cluster_h_stubs.go`, `cluster_help_menu.go`, `cluster_ops.go`, `clustercommand.go`, `clusterflags.go`, `clusterstatus.go`, `clustersubcmd.go`
- `cmdcode`: `code.go`
- `cmdcodingguidelines`: `codingguidelines.go`, `codingguidelines_commit.go`, `codingguidelines_compat.go`
- `cmdcommit`: `commit_batch.go`, `commit_cmd.go`, `commit_help_menu.go`, `commit_push.go`, `commit_ui.go`
- `cmdcommitin`: `commitin.go`, `commitin_ui_server.go`
- `cmdcommitpull`: `commitpull.go`, `commitpull_bootstrap.go`
- `cmdcommittransfer`: `committransfer.go`, `dispatchcommittransfer.go`
- `cmdcompletion`: `completion.go`, `completion_keys_printers.go`, `completion_printers.go`
- `cmdconsole`: `console_other.go`, `console_windows.go`
- `cmdcopy`: `copy_paste.go`
- `cmdcreate`: `create_cmd.go`, `create_ops.go`
- `cmddashboard`: `dashboard.go`
- `cmddiff`: `diff.go`, `diff_help_menu.go`, `dispatchcompare.go`, `dispatchdiff.go`
- `cmddiffprofiles`: `diffprofiles.go`, `diffprofilesops.go`, `diffprofilesrender.go`
- `cmddir`: `child_path.go`, `dir_scope_matcher.go`
- `cmddocs`: `docs.go`
- `cmddownload`: (wave 2)
- `cmdenv`: `env.go`, `envops.go`, `envplatform_unix.go`, `envplatform_windows.go`, `envregistry.go`, `envvalidate.go`
- `cmdexec`: `exec.go`, `execprint.go`
- `cmdexplorer`: `explorer.go`
- `cmdexport`: `export.go`, `exportrender.go`
- `cmdfind`: `find_duplicates.go`, `find_duplicates_git.go`, `find_entry.go`, `find_files.go`, `find_files_match.go`, `find_help_menu.go`
- `cmdfindnext`: `findnext.go`, `findnextflags.go`, `findnextrender.go`
- `cmdfix`: `fix_cmd.go`, `fix_diagnostics.go`, `fix_execute.go`, `fix_ls.go`, `fix_ls_table.go`
- `cmdfixauth`: `fixauth.go`, `fixcredential.go`
- `cmdgitrm`: `gitrm.go`, `gitrm_entry.go`
- `cmdgomod`: `gomod.go`, `gomodbranch.go`, `gomodreplace.go`
- `cmdgroup`: `group.go`, `groupadd.go`, `groupcreate.go`, `groupdelete.go`, `grouplist.go`, `groupremove.go`, `groupscoped.go`, `groupshow.go`
- `cmdhaschange`: `hasanyupdates.go`, `haschange.go`
- `cmdheadtail`: `headtail.go`
- `cmdhistory`: `command_history_cli.go`, `history.go`, `history_purge_cmd.go`, `historyrender.go`, `historyreset.go`, `historyrewrite.go`, `historyrewrite_flags.go`, `historyrewrite_helpers.go`, `historyrewrite_paths.go`, `historyrewrite_pin.go`, `historyrewrite_push.go`, `historyrewrite_sandbox.go`, `historyrewrite_verify.go`
- `cmdhistoryrewrite`: (fold into `cmdhistory` — same family; lead to confirm)
- `cmdhygiene`: `hygiene_format.go`, `hygiene_parallel.go`, `hygiene_parallel_map.go`, `hygiene_parallel_workers.go`
- `cmdimport`: `import_all_json_cmd.go`, `importcmd.go`
- `cmdindex`: `index_cmd.go`
- `cmdinject`: `inject.go`, `inject_idempotency.go`
- `cmdinteractive`: `interactive.go`
- `cmdip`: `ip_cmd.go`, `ip_resolver.go`, `ip_subcmds.go`, `ipchange_cmd.go`
- `cmdjoin`: `join.go`
- `cmdjsonextract`: `jsonextract.go`
- `cmdlatestbranch`: `latestbranch.go`, `latestbranchcsv.go`, `latestbranchoutput.go`, `latestbranchrender.go`, `latestbranchresolve.go`, `latestbranchswitch.go`
- `cmdlfscommon`: `addlfsinstall.go`, `lfscommon.go`
- `cmdlist`: `list.go`, `list_preview.go`, `list_tree.go`, `listreleases.go`, `listreleasesallrepos.go`, `listreleasesload.go`, `listreleasesrender.go`, `listversions.go`, `listversionsrender.go`, `listversionsutil.go`
- `cmdllm`: `llm_entry.go`, `llm_help_menu.go`, `llmdocs.go`, `llmdocscommands.go`, `llmdocsgroups.go`, `llmdocsheader.go`, `llmdocsrender.go`, `llmdocssections.go`
- `cmdlocate`: `locate_entry.go`
- `cmdmerge`: `dispatchmovemerge.go`, `merge.go`, `merge_json_cmd.go`, `move.go`, `movemergeflags.go`
- `cmdmigrate`: `migrate.go`, `migrate_cmd.go`, `migrate_wizard.go`
- `cmdmkdir`: `mkdir.go`
- `cmdmultigroup`: `multigroup.go`, `multigroupops.go`
- `cmdmuse`: (wave 1)
- `cmdopen`: `open.go`, `open_opener.go`, `open_target.go`
- `cmdorphans`: `orphans.go`
- `cmdpending`: `dopending.go`, `dopendingretry.go`, `pending.go`, `pending_commits_cmd.go`, `pending_commits_types.go`, `pendingclear.go`, `pendingtaskhelper.go`
- `cmdpin`: `pin_cmd.go`
- `cmdpower`: (→ `cmdos`, wave 2)
- `cmdprobe`: `probe.go`, `probeflags.go`, `proberender.go`, `probereport.go`
- `cmdprofiles`: `profile.go`, `profileops.go`, `profiles_cmd.go`, `profiles_ops.go`, `profileutil.go`
- `cmdprojectrepos`: `projectrepos.go`, `projectreposoutput.go`, `projectreposrender.go`
- `cmdprune`: `prune.go`, `pruneops.go`
- `cmdreconcile`: `reconcile.go`, `reconcile_cmd.go`, `reconcile_db.go`, `reconcile_prompt.go`, `reconcile_types.go`
- `cmdrecreate`: `recreate_repo.go`
- `cmdregoldens`: `regoldens.go`, `regoldens_diff.go`, `regoldens_diff_print.go`, `regoldens_diff_scope.go`, `regoldens_dryrun.go`, `regoldens_exec.go`, `regoldens_precheck.go`, `regoldens_validate.go`
- `cmdrelease`: `clearreleasejson.go`, `release.go`, `release_notes_opts.go`, `release_scan_commits.go`, `release_tools.go`, `releasealias.go`, `releasealias_git.go`, `releaseargs.go`, `releaseautobump.go`, `releaseautoregister.go`, `releasebranch.go`, `releasepending.go`, `releasepersist.go`, `releaserebase.go`, `releasepull.go`, `releaserecentclone.go`, `releasescan.go`, `releaseself.go`, `releaseundo.go`, `releaseundorange.go`
- `cmdremediation`: `remediation_box.go`, `remediation_local.go`, `remediation_matcher.go`, `remediation_state_ops.go`, `remediation_suggest.go`
- `cmdreplace`: `replace.go`, `replace_classify.go`, `replace_dashn.go`, `replaceapply.go`, `replaceaudit.go`, `replaceflags.go`, `replaceversion.go`, `replaceversionrun.go`, `replacewalk.go`
- `cmdrepo`: `repo_cmd_dispatch.go`, `repo_create_init.go`, `repo_create_params.go`, `repo_create_profile.go`, `repo_create_remote.go`, `repo_create_report.go`, `repo_create_slug.go`, `repo_create_smart.go`, `repo_create_tokens.go`, `repo_db_dispatch.go`, `repo_db_ops.go`, `repo_recreate.go`
- `cmdresolver`: `resolver.go`, `resolver_alias.go`, `resolver_glob.go`, `resolver_path.go`, `resolver_pwd.go`, `resolver_types.go`
- `cmdrest`: `rest_enable.go`
- `cmdrevert`: `revert.go`, `revertscript.go`, `reverttxn.go`, `reverttxn_lastn.go`
- `cmdsafe`: `safe_rm_cmd.go`
- `cmdsafety`: `safety_snapshot.go`
- `cmdsearch`: `search.go`, `search_entry.go`, `search_help_menu.go`
- `cmdselfinstall`: `selfinstall.go`, `selfuninstall.go`, `selfuninstallhandoff.go`, `selfuninstallparts.go`
- `cmdsends`: `sends_cmd.go`
- `cmdsequence`: `sequence_cmd.go`
- `cmdseowrite`: `seowrite.go`, `seowritecreate.go`, `seowritecsv.go`, `seowritegit.go`, `seowriteloop.go`, `seowritetemplate.go`
- `cmdservercmd`: `serve.go`, `servercmd.go`, `servercmd_run.go`
- `cmdsetsourcerepo`: `setsourcerepo.go`
- `cmdsf`: `sf.go`
- `cmdsize`: `size.go`
- `cmdstale`: `stale.go`
- `cmdstash`: `stash_help_menu.go`
- `cmdstartup`: `startup.go`, `startup_cmd.go`, `startup_run.go`, `startupadd.go`, `startuplistfilter.go`, `startuplistrender.go`, `startupstatusjson.go`
- `cmdstats`: `stats.go`, `statsrender.go`
- `cmdstatus`: `status.go`, `status_branch.go`, `status_help_menu.go`, `status_id.go`, `status_md.go`, `status_pr.go`
- `cmdstorage`: `storage_clean.go`, `storage_cmd.go`, `storage_display.go`, `storage_drive_other.go`, `storage_drive_windows.go`, `storage_help_menu.go`, `storage_ls.go`, `storage_reset.go`, `storage_restore.go`, `storage_types.go`
- `cmdtemplates`: `templates_help_menu.go`, `templates_state_cli.go`, `templates_ui_server.go`, `templatescli.go`, `templatesdiff.go`, `templatesinit.go`
- `cmdtemprelease`: `temprelease.go`, `tempreleaseexport_test_helpers.go`, `tempreleaselist.go`, `tempreleaselistrender.go`, `tempreleaseops.go`, `tempreleaseremove.go`
- `cmdundo`: `undo.go`
- `cmduser`: `user_cmd.go`, `user_help_menu.go`, `user_info.go`, `user_profiles_cmd.go`, `user_project_cmd.go`
- `cmdvariable`: `variable_cmd.go`
- `cmdversion`: `version_tags.go`
- `cmdversionhistory`: `versionhistory.go`, `versionhistoryrender.go`
- `cmdvisibility`: all 23 `visibility*.go` files
- `cmdwatch`: `watch.go`, `watchformat.go`, `watchops.go`, `watchrender.go`
- `cmdwhoami`: `whoami.go`
- `cmdwpr`: (deleted wave 1 — was pure delegate)
- `cmdbash`: `bash_runner.go`
- `cmdpowershell`: `powershell_runner.go`
- `cmdauthor`: `author_sponsor.go`
- `cmdas`: `as.go`, `asops.go`
- `cmdappend`: `append.go`
- `cmddesktopsync`: `desktopsync.go`, `desktopsync_ops.go`
- `cmdct`: `ct.go`

**Wave 4 — STAY in `cmd/` (thin dispatch / alias context / argv preprocessing only):**
`alias.go`, `aliasops.go`, `aliasresolve.go`, `aliassuggest.go`, `auto_alias.go`,
`help.go`, `help_dynamic.go`, `help_category_groups.go`, `helpcheck.go`,
`helpdashboard.go`, `helpdashboard_download.go`, `helpdashboard_extract.go`,
`helpdashboard_fallback.go`, `helpdashboard_target.go`, `themeflag.go`,
`themeflag_parse.go`, `prettyflag.go`, `glyphsflag.go`, `fallback.go`, `hints.go`,
`commons.go`, `rich_help_dispatcher.go`, plus `root*.go` (dispatch tables),
`clihelpers_argv.go`, `di_hooks.go`, `dispatch_pull.go`, `dispatch_scan.go`,
`dispatch_ip.go` (new dispatch keepers from wave 1).

### Collision ledger
- 2026-10-09: zero filename collisions between `cli/cmd/*.go` and any `cli/cmdX/*.go`.
- 2026-10-09: Worker D-1 new files (83, in cmdpull/cmdpipeline/dbengine/cmdupdate/cmdchromeprofile/cmdui/pipelinedb) do not collide with Wave-1 D-2 files (cmdai/ai_flag.go, cmdide/ide_scan_sync.go, cmdssh/remote_dispatch.go, cmdmuse/muse_cmd.go — all new names).
- Rule: every wave re-checks destination filenames before moving; on collision, pick a non-colliding name and record it here.

## D4 workstreams (COMPLETED 2026-10-09 — Worker D-2; lead to verify with go build)

### D4a — `cli/enums/` (new)
- Read the Go enum pattern first: `~/workspace/repos/coding-guidelines/spec/02-coding-guidelines/03-coding-guidelines-spec/03-golang/01-enum-specification/01-enum-pattern.md` (+ `03-folder-structure.md`). Cited in `cli/enums/doc.go`.
- Created `cli/enums/doc.go` (package doc + pattern citation) and
  `cli/enums/ctxmodetype/variant.go` per the pattern (byte, Invalid-first,
  iota, PascalCase, single variantLabels table, String/Label/IsValid/
  IsInvalid/Is{Value}/All/ByIndex/Parse/Values/MarshalJSON/UnmarshalJSON;
  package name carries the `type` suffix; file is `variant.go`).
- Migrated the most self-contained USED enum-like group:
  `CtxMode` (`terminal`/`silent`/`prefill`) from
  `cli/constants/constants_installctx.go` → `enums/ctxmodetype.Variant`.
  Reference counts verified via `gitmap aum search`: CtxMode ~40 sites,
  all inside `cli/cmdinstall/` (vs TaskType's 56 cross-package sites).
  All references updated (impl + tests); serialized values preserved
  verbatim ("terminal"/"silent"/"prefill") per the pattern's
  protocol-driven exception (baked into Windows registry entries).
  Note: `SettingType` (3 values) was even smaller but had ZERO functional
  references — migrating an unused type would have been theater.

### D4b — `.ai-memory/overview.md`
- Stale `v3.1.0` claim fixed → `v6.515.0`, read from `version.json`
  (read-only; versioning untouched).

### D4c — `docs/architecture/state-ownership.md` (new)
- One page: `.gitmap/` working state vs `gitmap.json` derived scan artifact
  vs SQLite (`gitmap.db` owned by `store`; split DBs
  `<base>/data/<section>/<slug>/sql.db` owned per-feature via
  `store.ResolveSplitDbPath`; satellites `gitmap-errors.db`,
  `installation.db`) vs `config.json` discovery order. Researched from
  `cli/store/split_db_path.go`, `cli/pipelinedb/`, `cli/config/`,
  `cli/constants/`.

### D4d — new-command checklist
- No existing checklist file found after repo-wide search. Created
  `docs/architecture/new-command-checklist.md` mandating the spec-03
  `HelpDisplay` technique: struct-first, generated `helptext/*.md`,
  single-touch registration (+ the mechanical steps: `Cmd*` constant,
  owning `cmdX/` package, one `dispatchEntry`, CI coverage test).
