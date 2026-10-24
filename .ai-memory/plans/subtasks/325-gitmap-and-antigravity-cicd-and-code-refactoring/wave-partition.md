# Wave partition — plan 325 enum/constants refactor (re-audit)

Re-audit at tip fd096c26 (2026-10-25). Method: GitMap search only (`gitmap aum search`); stray enums are decl-verified (`type XType` string/int/uint8 outside the package's own `enums/` subpackage). Size = non-test `.go` files in the package directory. All paths are repo-relative. Do NOT touch `cli/cmdcommit` (dry-run worker territory) or the coordinator-owned pending plan file.

## Verdict: Wave 2 is NOT done — do it first

No wave-2 commit exists in the log since the 03:05 audit. Two Wave 2 items did land, unlabeled: `cli/cmdpull` (chain to slice lookup, see below) and the `cli/cmdcargo` exemplar. Everything else in Wave 2 is open, above all 101 plain-error `Parse` functions in done packages' `enums/` dirs. Order: Wave 1 remainder, Wave 2 remainder, Wave 3, Wave 4, then sweeps (sweep list B lands after its owning wave; sweep list A may run in parallel with the waves).

## Refreshed totals (03:05 audit -> now)

- Own `enums/` dirs: 27 -> 27, plus shared `cli/enums` = 28 dirs; files 49 -> 78 (51 non-test, 27 test); 100 `*Type` decls inside.
- Stray-enum packages (decl-verified): 48 -> 64 packages, 109 files, 160 types. Rise is partly method: this audit counts uint8/int enums the 03:05 count missed (`cli/cmd/commitin` alone holds 8 uint8 types in one file, plus 1 int type).
- String `+` on enum values: 9 -> 4 true-enum sites (all wave-owned, none left for sweeps); plus 9 `Version.String()` sites, see Wave 2.
- `arg ==` lines: 517 tree-wide -> 511 in `cli/` (505 non-test).
- `regexp.MustCompile`: 166 -> 166 (162 non-test, 4 test hits in 2 files); classified: 156 package-level static, 6 func-local.
- `regexp.Compile(`: 18 -> 17 (15 non-test, all dynamic; 2 test).
- Enum `Parse` returning plain `error`: "every done enum" -> 107 (101 done-enums + 1 shared + 5 stray-enum). Only `cli/cmdcargo` (4/4) returns `*apperror.AppError` among done packages.
- `lazyregex.New` users (non-test): cmdnodes 4, cmdpipeline 3, cmdcommit 3 (territory), cmdssh 1, committransfer 2, searcher 2. Pattern to copy: `cli/lazyregex` (`lazyregex.New`, pattern from constants, lazy compile).

## Wave 1 remainder (store/dbengine core leftovers, single worker)

- `cli/store`: `backup_prompts_split_types.go` PromptStatusType, `heatmap_schema.go` FileTouchType, `repo_transport.go:56` `string(strategy)` site (RepoTransportQueryType), `downloader_seed.go:87,91` chains (2 lines).
- `cli/dbengine`: `dialect.go` DbType, `operators.go` SqlOperatorType, `query.go` whereType.
- Already done, do not redo: `OpenSQLiteDB` flat conversion, dbengine pragma dedupe (`applySQLitePragmas` now 0 hits), the 3 agy error-discard sites (now `IsFailure()` + warning).
- Open question for the Wave 1 worker: batch 3a renamed DbSQLite to DBSQLite but left the type name DbType — confirm whether DbType to DBType is intended before renaming (blast radius).

## Wave 2 remainder (done-package cleanups — do first; 24 units: 19 enums dirs + cmdclone, cmdcargo, cmdrelease, release, shared ctxmodetype)

Parse conversions (101, one enums file per package; convert return from `(Type, error)` to `(Type, *apperror.AppError)`; most call sites compile unchanged because AppError satisfies error — fix the rest via build errors):
- `cli/cmdadd/enums` add_enums.go (3), `cli/cmdbranch/enums` branch_enums.go (4), `cli/cmdconfig/enums` config_enums.go (7), `cli/cmdconflict/enums` conflict_enums.go (4), `cli/cmdcopy/enums` copy_enums.go (4), `cli/cmddiff/enums` diff_enums.go (4), `cli/cmddocs/enums` docs_enums.go (6), `cli/cmdhealth/enums` health_enums.go (8), `cli/cmdhooks/enums` hooks_enums.go (4), `cli/cmdlogin/enums` login_enums.go (8), `cli/cmdplans/enums` plan_enums.go (4), `cli/cmdprune/enums` prune_enums.go (4), `cli/cmdrevert/enums` revert_enums.go (3), `cli/cmdsearch/enums` search_enums.go (10), `cli/cmdstale/enums` stale_enums.go (3), `cli/cmdstats/enums` stats_enums.go (10), `cli/cmdstatus/enums` status_enums.go (8), `cli/cmdsync/enums` sync_enums.go (4), `cli/cmdundo/enums` undo_enums.go (3). Several define short alias Parses too (for example ParseDocsCategory and ParseDocsCategoryType) — convert both.
- Shared: `cli/enums/ctxmodetype/variant.go:105` Parse (1).
- Chains riding with the Parse packages (convert to set lookup): `cli/cmdconfig/which_format_cmd.go` (1), `cli/cmdlogin/login_ssh.go` (1), `cli/cmdplans/plans.go` (5), `cli/cmdsearch/helpers_shim.go` + `search.go` (10), `cli/cmdsync/sync_dispatch.go` (1).
- `cli/cmdclone` (partially done): `clone_rc.go` CloneDispatchModeType plus `multiclone_options.go` chain (7 lines, still literal flags — constantize first, then set lookup).
- `cli/cmdcargo` residual: `cargo_flags.go:23` (already constant-based 2-way chain — slice contains).
- `cmdrelease`/`release` `+` sites (Version is a struct, not an enum, but the 03:05 audit counted these — convert via strings.Join/prefix helper or record an explicit exemption, do not silently leave): `cli/cmdrelease/releaseundo.go:47`, `cli/release/workflow.go:99,207`, `cli/release/workflowgithub.go:36`, `cli/release/workflowpending.go:97,105`, `cli/release/workflowvalidate.go:100`, `cli/release/metadata.go:84` (struct field, same pattern).
- DONE, no Wave 2 item: `cli/cmdpull` — chain is now `slices.Contains(pullValuedFlagTokens, arg)` in `pull_cli_input.go`, and the FlagWorkerShort suspect is resolved by comment (registered as `fs.Int(FlagWorkerShort)` at line 170; Go flag accepts `--w`). Landed post-audit, unlabeled. `cli/cmdcargo` exemplar also DONE (`cargo_exec.go` uses `enums.ToolchainEnvPath.String()`, `.WithValue()`, `.AssignmentPrefix()`; zero `+` sites left).
- DONE, no remainder: cloner, pipelinedb, repodb (enums, no Parse, no stray, no chain). Their regex files are in sweep list A/B.

## Wave 3 — remaining SMALL packages (size at most 20; 76 packages, 614 non-test files: 42 stray-enum + 34 chain-only)

Stray-enum packages (decl file: enums; chains noted where present):
- `cli/archive` (6) types.go: CompressionModeType, FormatType, SourceKindType
- `cli/cliexit` (5) exitcodes.go: ExitCodeType; kind.go: KindType; report.go: OutputModeType
- `cli/cloneconcurrency` (3) pull_concurrency.go: CPUProfileType
- `cli/cmd/commitin/workspace` (6) source.go: SourceKindType
- `cli/cmd/folder` (7) folder.go: OutputFormatType; chains 12
- `cli/cmdapps` (6) types.go: AppScopeType, PackageManagerType
- `cli/cmdcpar` (3) types.go: CPARReviewActionType
- `cli/cmddesktopsync` (4) desktopsync_ops.go: SyncResultType; chains 1
- `cli/cmddownload` (10) types.go: DownloadEngineType; chains 9
- `cli/cmdfind` (7) find_files.go: MatchKindType; chains 5
- `cli/cmdfixreleasetags` (10) types.go: AuditReasonType, ProtectionStatusType, StepStatusType; chains 1
- `cli/cmdhistory` (15) historyrewrite.go: HistoryModeType; chains 4
- `cli/cmdhygiene` (6) hygiene_format.go: HygieneFormatType
- `cli/cmdide` (9) ide_types.go: IDETargetType; chains 7
- `cli/cmdlemon` (18) command_whitelist.go: WhitelistAuthType; lemon_types.go: EmailFormatType; service_manager.go: ServiceActionType, ServiceStatusType; chains 1
- `cli/cmdmail` (10) folder_watcher.go: FileEventType; chains 1
- `cli/cmdpurge` (11) preflight.go: TargetType
- `cli/cmdpushfix` (6) pushfix_ui.go: PushFixBadgeType
- `cli/cmdreplace` (10) replace_classify.go: ReplaceModeType; chains 2
- `cli/cmdresolver` (9) repo_feature.go: RepoFeatureTargetType
- `cli/cmdsends` (1) sends_grouping.go: FileCategoryType
- `cli/cmdsetup` (16) setup_perms.go: PermsAppType (no-error Parse — verify default-on-unknown, see below)
- `cli/cmdvhost` (16) vhost_types.go: VHostSiteType; chains 1; plus `vhost_templates.go:123` string(siteType) site
- `cli/cmdzsh` (19) types.go: PackageManagerType, ZshThemeType
- `cli/committransfer` (15) types.go: DirectionType, PreferPolicyType
- `cli/config` (6) validate_shape.go: JSONKindType
- `cli/db` (9) enums.go: CommandKindType, ResultStatusType; plus plain-error Parse in the same file (2)
- `cli/diag` (6) errreport.go: PhaseType; types.go: ActionType, RenderFormatType, SuggestionCategoryType
- `cli/diff` (3) tree.go: EntryKindType
- `cli/downloaderconfig` (3) types.go: ConfigKeyType
- `cli/glyphs` (5) glyphs.go: ModeType (no-error Parse — verify)
- `cli/logging` (1) jsonlog.go: LevelType
- `cli/movemerge` (14) types.go: ChoiceType, DiffKindType, DirectionType, EndpointKindType, PreferPolicyType
- `cli/render` (8) preflag.go: PrettyModeType; pretty_parse.go: BlockKindType
- `cli/secrets` (8) ghtoken.go: SourceType
- `cli/startup` (18) add.go: AddStatusType; remove.go: RemoveStatusType; winbackend.go: BackendType; plus plain-error ParseBackend in winbackend.go (1)
- `cli/templates` (10) diff.go: DiffStatusType; merge.go: MergeOutcomeType; resolver.go: SourceType
- `cli/termgeneric` (4) header.go: AlignmentType
- `cli/termout` (15) termtable_types.go: AlignmentType, BorderStyleType, WrapModeType; themes.go: RoleType; theme_theme.go: ModeType (no-error Parse — verify)
- `cli/types` (1) agent_types.go: AgentActionType
- `cli/usercontext` (7) types.go: AuthStatusType
- `cli/vscodepm` (15) mergemode.go: MergeModeType; plus plain-error ParseMergeMode in the same file (1)
- Chain-only packages (no enum move; convert `arg ==` chains; size, chain lines): `cli/cmd/llm` (13, 2), `cli/cmdamend` (9, 2), `cli/cmdautofix` (17, 1), `cli/cmdbackup` (5, 1), `cli/cmdchrome` (19, 1), `cli/cmdcluster` (9, 3), `cli/cmdcode` (3, 1), `cli/cmdcommittransfer` (4, 2), `cli/cmdcursor` (7, 12), `cli/cmddir` (4, 7), `cli/cmderrors` (11, 1), `cli/cmdfix` (15, 4), `cli/cmdfixauth` (2, 2), `cli/cmdfixgit` (12, 7), `cli/cmdflatten` (2, 4), `cli/cmdgit` (6, 4), `cli/cmdgroup` (10, 1), `cli/cmdignore` (12, 3), `cli/cmdlfs` (2, 1), `cli/cmdlist` (13, 5 — also convert the Version `+` site in `listreleasesload.go:139`), `cli/cmdports` (4, 4), `cli/cmdprofiles` (7, 1), `cli/cmdpy` (3, 1), `cli/cmdrecreate` (2, 1), `cli/cmdremediation` (6, 4), `cli/cmdrepo` (16, 1), `cli/cmdservice` (8, 1), `cli/cmdsplit` (3, 1), `cli/cmdsummary` (7, 1), `cli/cmdsupabase` (5, 2), `cli/cmdtoolchain` (5, 3), `cli/cmdvariable` (2, 1), `cli/cmdwinutil` (4, 1), `cli/cmdworkgroup` (6, 7). Chain files per package: re-derive with one search (`gitmap aum search "arg ==" --path cli/<pkg>`).

## Wave 4 — remaining LARGE packages (size over 20; 24 packages, 1496 non-test files: 16 stray-enum + 8 chain-only)

Stray-enum packages:
- `cli/cluster` (29) broadcast.go: CommandType; registry.go: NodeStateType; node_resolver.go: TargetSelectorType
- `cli/cmd/commitin` (21) conflict.go: ConflictDecisionType; enums.go (in-package file, not an `enums/` subdir — move as a unit, keep uint8): ConflictModeType, InputKindType, RunStatusType, CommitOutcomeType, SkipReasonType, ExclusionKindType, MessageRuleKindType, FunctionIntelLanguageType
- `cli/cmdagent` (33) agent_log.go: ActionType; chains 2
- `cli/cmdagy` (238) agy_fix_pipeline_types.go: AgyInjectionModeType, AgyConvStatusType; agy_running_prompts_types.go: PromptExportFormatType (no-error Parse — verify); chains 47
- `cli/cmdai` (33) ai_analysis_types.go: AiTaskStatusType, AiOperationType; ai_train_types.go: LlmTrainingFormatType, ExportFormatType; types.go: ScriptCategoryType, FixTargetType, ScriptTemplateType; chains 2
- `cli/cmdfoldertree` (30) treesearch_match.go: SearchModeType; chains 8
- `cli/cmdinstall` (112) install_archive_types.go: ArchiveInstallStrategyType; install_muse.go: MusePlatformType; installer_types.go: InstallerOSType; osdetect.go: OSType; chains 3; plus `install_muse.go:97` string(platform) site
- `cli/cmdmacro` (28) macro_add_interactive.go: InteractiveLoopActionType; chains 19
- `cli/cmdnodes` (48) nodes_clone_types.go: NodesCloneKindType; nodes_commits_types.go: CommitActionKindType; chains 49
- `cli/cmdos` (126) os_dm_types.go: DMType; os_dock_types.go: DockPositionType; os_theme_types.go: ThemeModeType; os_tui_types.go: OSTUIViewModeType; os_tweak_types.go: TweakTargetType; chains 10
- `cli/cmdschedule` (22) schedule_os.go: OSActionType; chains 15; plus `schedule_os.go:67` string(action) site
- `cli/cmdspace` (23) email_tasktype.go: EmailTaskType; chains 12; plus plain-error ParseEmailTaskType in the same file (1)
- `cli/cmdssh` (194) ssh_deploy_sync.go: DeploySyncModeType, DeployActionType; ssh_troubleshoot.go: SSHDiagnosisType; chains 81
- `cli/cmdupdate` (43) update_ui_types.go: UpdateStageType; chains 15
- `cli/cmdvisibility` (25) visibilitydriftguard.go: DriftActionType
- `cli/macro` (21) recurse.go: ContextKeyType
- Chain-only packages (size, chain lines): `cli/cmd` (124, 11), `cli/cmdautomation` (85, 7), `cli/cmdcg` (22, 1), `cli/cmdchromeprofile` (43, 10), `cli/cmddb` (65, 23), `cli/cmdpending` (22, 14), `cli/cmdpipeline` (82, 2), `cli/cmdvscode` (27, 2).
- The four largest (cmdagy, cmdssh, cmdos, cmdinstall) each get a dedicated worker; do not batch them with other packages.

## Sweep: lazy regex (83 files, 156 package-level static hits)

Convert package-level `regexp.MustCompile` to `lazyregex.New` (pattern from constants + lazy compile). Lists are file-disjoint from wave-touched files (enum decl files, enum Parse files, `arg ==` files). `regexp.Compile(` needs no sweep: all 15 non-test hits are dynamic (pattern variables), plus the lazyregex implementation itself — keep.
- List A — parallel-safe packages (23 files, 37 hits; run any time): `cli/clonefrom` execute_lfs_fix.go; `cli/cloner` safe_pull.go; `cli/cmd/commitin/funcintel` lang_go.go, lang_java.go, lang_javascript.go, lang_php.go, lang_python.go, lang_rust.go, lang_ts.go; `cli/cmd/commitin/message` title_replace.go; `cli/cmdasset` asset_prnt_fetcher.go; `cli/cmdcodingguidelines` codingguidelines_compat.go; `cli/cmdenv` envvalidate.go; `cli/cmdfixrepo` fixrepo_identity.go; `cli/cmdpull` pull_stream.go; `cli/cmdspec` spec.go, spec_progress.go; `cli/constants` constants_commitin_tagreplay.go; `cli/helpdoc` bump.go, emitter.go, fixtureversion.go, hash.go; `cli/visibility` version.go.
- List B — wave packages, file-disjoint (60 files, 119 hits; land after the owning wave): `cli/cmd/commitin` config_json.go; `cli/cmd/commitin/workspace` expand.go; `cli/cmdagent` agent_subtask_resolve.go, agent_task_sqlite_paths.go; `cli/cmdagy` agy_clean_cache_procs.go, agy_conv_inspect_ops.go, agy_enhance_prompts.go, agy_image_assets.go, agy_rerun_queue.go, agy_running_prompts_backup.go; `cli/cmdautofix` markdown.go, naming.go, paths.go (static hits only — the dynamic hit in paths.go stays); `cli/cmdautomation` doc_links.go, help_audit.go, milestones.go, naming.go, release_bump.go, relpaths.go, rules_audit.go, schema_audit.go, sequence_auditor.go, version_sync.go; `cli/cmdchromeprofile` chromeprofile_snapshot_files.go; `cli/cmdclone` helpers.go, multiclone_parse.go; `cli/cmdinstall` agm_uninstall.go, agm_update_ui.go, install_urls.go, installprobe.go; `cli/cmdnodes` nodes_clone_probe.go; `cli/cmdos` os_fix_mirrors.go; `cli/cmdplans` plans_parser_meta.go; `cli/cmdrelease` release_winscan.go, releasepending.go, releaserebase.go; `cli/cmdreplace` replaceversion.go; `cli/cmdschedule` duration_parser.go; `cli/cmdsetup` helpers.go, setup_wp.go (static hits only — dynamic hit stays), setup_wp_salts.go; `cli/cmdssh` ssh_known_hosts_prune.go, ssh_path_expand.go; `cli/cmdsync` boundary.go, mirror.go, worker.go; `cli/cmdupdate` update_ui.go, updateprobe.go; `cli/cmdvisibility` visibilityexceptlatest.go; `cli/cmdvscode` vscode_repair_bin.go; `cli/committransfer` message.go; `cli/config` variables.go; `cli/macro` expand.go; `cli/release` scan_parser.go; `cli/render` pretty_emit.go, prettypost.go, prettypost_spans.go; `cli/store` templates_array_vars.go, templates_split_import_export.go; `cli/cmddb` cmd_view_helpers.go (static hits only — the per-alias dynamic hit in the loop stays).
- NOT sweep targets: `cli/cmd/commitin/runlog/tagreplay.go` (already lazy via sync.Once); `cli/cmdautofix/misspell.go` (init-time, pattern built from map keys — sweep worker decides keep vs sync.Once); `cli/templates/merge.go` (dynamic per-tag hit; file is wave-owned anyway); 4 test hits in `cli/cmd/outputfilenames_contract_test.go` and `cli/store/erd_parity_test.go` (test-only, convert only if the sweep owner wants).

## Sweep: shared constants (tier 1 first)

No existing constant covers these (0 hits in `cli/constants`). Create the constants in `cli/constants`, then replace the literal uses. Files below are non-test and not wave-touched; `helpdoc` schedule.md is documentation, skip it.
- `application/json` and/or `Content-Type` (44 code files, 18 pkgs; the 45th hit is the skipped doc file above): `cli/clonenext` github.go; `cli/cmd` hd_command_exec.go, hd_terminal.go, helpdashboard.go; `cli/cmdagent` agent_ui_assets.go, agent_ui_logs.go; `cli/cmdagy` agy_account_switch_notify.go, agy_sug_ui.go, agy_ui_script.go, agy_ui_tree_handlers.go, agy_ui_tree_html.go; `cli/cmdai` ai_memory_server.go; `cli/cmdcommitin` commitin_ui_server.go; `cli/cmddaemon` daemon_client.go, daemon_server.go; `cli/cmddb` rest_envelope.go, rest_swagger.go, rest_swagger_asset.go, rest_ui.go, rest_visualizer_js2.go; `cli/cmdmail` mail_sender.go, mail_ui_html.go, mail_ui_server.go; `cli/cmdpurge` release_prune.go; `cli/cmdschedule` schedule_ui_html.go, schedule_ui_server.go; `cli/cmdspace` schedule_trigger_sender.go, send_email_smtp.go; `cli/cmdtemplates` templates_ui_server.go; `cli/cmdui` ui_assets_js_fleet.go, ui_assets_js_nodes.go, ui_assets_js_settings.go, ui_assets_js_terminal.go, ui_server.go, ui_server_help.go, ui_server_instances.go, ui_server_settings.go, ui_server_ssh.go, ui_server_terminal.go; `cli/cmdupdate` fleet_remote.go; `cli/cmdurl` cmd_url.go; `cli/cmdvhost` vhost_templates_content.go; `cli/release` assetsupload.go, githubapi.go.
- `Bearer ` (9 files, 7 pkgs): `cli/cmd/workflow_open_pr.go`, `cli/cmddoctor/doctor_extra.go`, `cli/cmdlogin/constants/constants.go`, `cli/cmdlogin/login_token.go`, `cli/cmdpurge/release_prune.go`, `cli/cmdsupabase/supabase_ping.go`, `cli/fspath/release-version.ps1`, `cli/fspath/release-version.sh`, `cli/fspath/release_version.py` (scripts — convert only if the sweep owner extends scope beyond Go).
- Tier 2 (verify each hit is code, not help text, fixture, or doc, before converting; re-derive with `gitmap aum search "<literal>" --path cli`): `already exists` (67 files, 31 pkgs), `.gitignore` (49, 23), `localhost` (38, 14), `README.md` (32, 14), `go.mod` (17, 11). Counts are non-test, non-wave-touched files.

## Excluded, and name collisions never to merge

- LEAVE in place: `cli/apperror` ErrorType/SeverityType (error infra, coordinator decision b); `cli/pipelinedb/enums` and `cli/repodb/enums` record structs (decision c; no Parse anyway).
- Sweep-owned, not wave work: `cli/constants` (36 own `*Type` enums — the shared-constants home) and `cli/lazyregex` (pattern library).
- Dry-run worker territory, excluded from every list above: `cli/cmdcommit` (stray OriginType in commit_info_types.go, 7 `arg ==` lines, 3 `lazyregex.New` uses already).
- Test files and `cli/cmd/commitin/e2e` test support: excluded.
- Never merge (same name, different meaning): ModeType (glyphs, termout; render PrettyModeType is distinct), DirectionType and PreferPolicyType (committransfer, movemerge), AlignmentType (termgeneric, termout), SourceKindType (archive, cmd/commitin/workspace), SourceType (secrets, templates), PackageManagerType (cmdapps, cmdzsh), ActionType (cmdagent, diag).
- No-error enum Parses (return the type, no error — wave worker verifies the default-on-unknown behavior, no error conversion applies): glyphs Parse, termout Parse, cmdsetup ParsePermsAppType, cmdagy ParsePromptExportFormatType. cmdpull ParseGitPullOutput (PullStepType, no error) is likewise not a conversion item.

## Coordination rules

- Package lists are disjoint: Wave 1, Wave 2, Wave 3, Wave 4 share no package (Wave 2 units are `enums/` subdirs except cmdclone, cmdcargo, cmdrelease, release, and shared ctxmodetype).
- Sweep lists are file-disjoint from wave-touched files by construction; list B shares packages with waves but never files. A wave's completion note lists every file it touched; sweep workers re-check their list against it before editing.
- Enum moves go to the package's own `enums/` subpackage (new or existing), following the cmdcargo exemplar: enum type, String(), WithValue/AssignmentPrefix-style helpers where the exemplar has them, and Parse returning `*apperror.AppError`.
- Counts method note: `arg ==` and literal searches cover `cli/` only; non-`*Type` string types without a String() method would not be caught (spot checks found none beyond the listed flag-value and struct String() methods, which are not enums).
