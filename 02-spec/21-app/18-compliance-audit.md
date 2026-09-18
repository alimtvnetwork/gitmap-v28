# Compliance Audit Summary

Last updated: 2026-03-10

## Rules Checked

| # | Rule | Source |
|---|------|--------|
| 1 | No negation in `if` conditions (`!`, `!=`, `== false`) | 01-overview §Code Style |
| 2 | Functions: 8–15 lines | 01-overview §Code Style |
| 3 | Files: 100–200 lines max | 01-overview §Code Style |
| 4 | One responsibility per package | 01-overview §Code Style |
| 5 | Blank line before `return` (unless sole line in `if`) | 01-overview §Code Style |
| 6 | No magic strings — all literals in `constants` | 01-overview §Code Style |
| 7 | No `switch` statements — use `if`/`else if` chains | 03-general/06 §Conditionals |

## Package: `cmd`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `root.go` | ~70 | ✅ Pass | Added revert + revert-runner dispatch |
| `rootflags.go` | ~50 | ✅ Pass | |
| `rootusage.go` | ~47 | ✅ Pass | Added list-versions + revert help |
| `scan.go` | ~113 | ✅ Pass | Split from 257 lines |
| `scanoutput.go` | ~155 | ✅ Pass | Extracted from scan.go |
| `scanimport.go` | ~60 | ✅ Pass | Release metadata import during scan |
| `clone.go` | ~140 | ✅ Pass | |
| `pull.go` | ~100 | ✅ Pass | Magic strings extracted |
| `rescan.go` | ~110 | ✅ Pass | Magic strings extracted |
| `status.go` | ~187 | ✅ Pass | |
| `statusformat.go` | ~135 | ✅ Pass | |
| `exec.go` | ~120 | ✅ Pass | |
| `list.go` | ~80 | ✅ Pass | |
| `listreleases.go` | ~80 | ✅ Pass | list-releases command |
| `listversions.go` | ~140 | ✅ Pass | list-versions command with changelog sub-points |
| `setup.go` | ~60 | ✅ Pass | |
| `update.go` | ~90 | ✅ Pass | |
| `updatescript.go` | ~120 | ✅ Pass | Magic strings extracted |
| `updatecleanup.go` | ~100 | ✅ Pass | Magic strings extracted |
| `release.go` | ~130 | ✅ Pass | |
| `releasebranch.go` | ~60 | ✅ Pass | |
| `releasepending.go` | ~40 | ✅ Pass | |
| `changelog.go` | ~80 | ✅ Pass | Magic strings extracted |
| `latestbranch.go` | ~80 | ✅ Pass | |
| `latestbranchresolve.go` | ~90 | ✅ Pass | |
| `latestbranchoutput.go` | ~100 | ✅ Pass | Magic strings extracted |
| `desktopsync.go` | ~100 | ✅ Pass | |
| `doctor.go` | ~60 | ✅ Pass | |
| `doctorchecks.go` | ~165 | ✅ Pass | Split; version logic extracted |
| `doctorversion.go` | ~120 | ✅ Pass | Extracted from doctorchecks.go |
| `doctorfixpath.go` | ~170 | ✅ Pass | Split; sync logic extracted |
| `doctorsync.go` | ~110 | ✅ Pass | Extracted from doctorfixpath.go |
| `doctorformat.go` | ~60 | ✅ Pass | Doctor output formatting |
| `group.go` | ~30 | ✅ Pass | |
| `groupcreate.go` | ~60 | ✅ Pass | |
| `groupdelete.go` | ~60 | ✅ Pass | |
| `groupadd.go` | ~60 | ✅ Pass | |
| `groupremove.go` | ~60 | ✅ Pass | |
| `grouplist.go` | ~50 | ✅ Pass | |
| `groupshow.go` | ~60 | ✅ Pass | |
| `dbreset.go` | ~40 | ✅ Pass | Database reset command |
| `revert.go` | ~90 | ✅ Pass | Revert command, validation, checkout, handoff |
| `revertscript.go` | ~85 | ✅ Pass | Revert-runner, PS1 script generation |
| `seowrite.go` | ~80 | ✅ Pass | SEO-write entry point |
| `seowritecreate.go` | ~60 | ✅ Pass | Template creation |
| `seowritecsv.go` | ~60 | ✅ Pass | CSV-based commit loop |
| `seowriteloop.go` | ~80 | ✅ Pass | Template-based commit loop |
| `seowritetemplate.go` | ~60 | ✅ Pass | Template loading and seeding |
| `amend.go` | ~80 | ✅ Pass | Author amendment command |
| `amendexec.go` | ~80 | ✅ Pass | Amendment execution logic |
| `amendaudit.go` | ~60 | ✅ Pass | Amendment audit logging |
| `amendlist.go` | ~80 | ✅ Pass | List stored amendments |
| `audit.go` | ~40 | ✅ Pass | Command audit tracking |
| `history.go` | ~80 | ✅ Pass | Command history display |
| `historyreset.go` | ~40 | ✅ Pass | Clear command history |
| `stats.go` | ~80 | ✅ Pass | Usage statistics |
| `bookmark.go` | ~40 | ✅ Pass | Bookmark routing |
| `bookmarksave.go` | ~60 | ✅ Pass | Save bookmarks |
| `bookmarklist.go` | ~50 | ✅ Pass | List bookmarks |
| `bookmarkrun.go` | ~40 | ✅ Pass | Run bookmarks |
| `export.go` | ~60 | ✅ Pass | Database export |
| `importcmd.go` | ~80 | ✅ Pass | Database import |
| `profile.go` | ~40 | ✅ Pass | Profile routing |
| `profileops.go` | ~80 | ✅ Pass | Profile CRUD operations |
| `profileutil.go` | ~60 | ✅ Pass | Profile utility helpers |
| `cd.go` | ~60 | ✅ Pass | CD command entry |
| `cddefault.go` | ~40 | ✅ Pass | CD default management |
| `cdops.go` | ~60 | ✅ Pass | CD navigation logic |
| `diffprofiles.go` | ~60 | ✅ Pass | Diff-profiles entry |
| `diffprofilesops.go` | ~80 | ✅ Pass | Diff-profiles comparison logic |
| `watch.go` | ~60 | ✅ Pass | Watch dashboard entry |
| `watchformat.go` | ~80 | ✅ Pass | Watch display formatting |
| `watchops.go` | ~60 | ✅ Pass | Watch refresh loop |
| `gomod.go` | ~125 | ✅ Pass | GoMod entry point, flag parsing, orchestration |
| `gomodreplace.go` | ~140 | ✅ Pass | File walking, go.mod parsing, path replacement |
| `gomodbranch.go` | ~140 | ✅ Pass | Branch creation, merge, slug derivation, git ops |
| `flags_test.go` | ~40 | ✅ Pass | |

## Package: `constants`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `constants.go` | ~110 | ✅ Pass | Core defaults, modes, formats, permissions |
| `constants_cli.go` | ~163 | ✅ Pass | CLI commands/aliases/help/flags; duplicates removed (v2.27.0 fix) |
| `constants_git.go` | ~61 | ✅ Pass | Git commands and argument strings |
| `constants_terminal.go` | ~180 | ✅ Pass | ANSI colors, banners, table headers |
| `constants_messages.go` | ~195 | ✅ Pass | User-facing messages and errors |
| `constants_release.go` | ~37 | ✅ Pass | Release workflow strings |
| `constants_store.go` | ~122 | ✅ Pass | Database schema and SQL queries |
| `constants_doctor.go` | ~91 | ✅ Pass | Diagnostic messages and binary lookup |
| `constants_update.go` | ~119 | ✅ Pass | Self-update and PowerShell templates |
| `constants_amend.go` | ~130 | ✅ Pass | Amend commands, flags, messages, SQL |
| `constants_seo.go` | ~153 | ✅ Pass | SEO-write commands, flags, messages, SQL |
| `constants_bookmark.go` | ~68 | ✅ Pass | Bookmark commands, messages, SQL |
| `constants_history.go` | ~91 | ✅ Pass | History commands, SQL, detail levels |
| `constants_stats.go` | ~66 | ✅ Pass | Stats SQL queries and formatting |
| `constants_export.go` | ~22 | ✅ Pass | Export command and messages |
| `constants_import.go` | ~20 | ✅ Pass | Import command and messages |
| `constants_profile.go` | ~47 | ✅ Pass | Profile commands and messages |
| `constants_cd.go` | ~43 | ✅ Pass | CD commands and messages |
| `constants_watch.go` | ~52 | ✅ Pass | Watch commands, display, flags |
| `constants_diffprofile.go` | ~28 | ✅ Pass | Diff-profiles commands and messages |
| `constants_clone.go` | ~10 | ✅ Pass | Clone progress format strings |
| `constants_gomod.go` | ~100 | ✅ Pass | GoMod commands, flags, messages, git args |

## Package: `release`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `workflow.go` | ~163 | ✅ Pass | Trimmed from 416; imports cleaned |
| `workflowfinalize.go` | ~145 | ✅ Pass | Added loadChangelogNotes helper |
| `workflowbranch.go` | ~165 | ✅ Pass | Extracted from workflow.go |
| `gitops.go` | ~100 | ✅ Pass | Rewritten; query functions extracted |
| `gitopsquery.go` | ~135 | ✅ Pass | Extracted from gitops.go |
| `changelog.go` | ~120 | ✅ Pass | Fixed `== false` → positive logic |
| `github.go` | ~66 | ✅ Pass | Fixed `IsDir() == false` → positive logic |
| `metadata.go` | ~145 | ✅ Pass | Added Changelog field to ReleaseMeta |
| `metadata_test.go` | ~40 | ✅ Pass | |
| `semver.go` | ~160 | ✅ Pass | Fixed switch → if/else chain |
| `semver_test.go` | ~80 | ✅ Pass | |

## Package: `formatter`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `terminal.go` | ~124 | ✅ Pass | Trimmed from 223; fixed `!quiet` |
| `terminaltree.go` | ~110 | ✅ Pass | Extracted from terminal.go |
| `csv.go` | ~60 | ✅ Pass | |
| `json.go` | ~30 | ✅ Pass | |
| `text.go` | ~30 | ✅ Pass | |
| `structure.go` | ~100 | ✅ Pass | |
| `clonescript.go` | ~40 | ✅ Pass | |
| `directclone.go` | ~70 | ✅ Pass | |
| `desktopscript.go` | ~50 | ✅ Pass | |
| `template.go` | ~30 | ✅ Pass | |
| `formatter_test.go` | ~60 | ✅ Pass | |

## Package: `cloner`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `cloner.go` | ~90 | ✅ Pass | |
| `progress.go` | ~60 | ✅ Pass | Clone progress display |
| `safe_pull.go` | ~110 | ✅ Pass | Trimmed from 213 |
| `pulldiag.go` | ~130 | ✅ Pass | Extracted from safe_pull.go |

## Package: `setup`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `setup.go` | ~131 | ✅ Pass | Trimmed from 206 |
| `setupapply.go` | ~100 | ✅ Pass | Extracted from setup.go |

## Package: `config`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `config.go` | ~78 | ✅ Pass | Fixed `os.IsNotExist` → `errors.Is` |
| `config_test.go` | ~30 | ✅ Pass | |

## Package: `scanner`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `scanner.go` | ~80 | ✅ Pass | |

## Package: `mapper`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `mapper.go` | ~110 | ✅ Pass | |
| `mapper_test.go` | ~50 | ✅ Pass | |

## Package: `model`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `record.go` | ~67 | ✅ Pass | |
| `group.go` | ~20 | ✅ Pass | |
| `release.go` | ~30 | ✅ Pass | Release metadata model |
| `amendment.go` | ~25 | ✅ Pass | Amendment record model |
| `history.go` | ~25 | ✅ Pass | Command history model |
| `stats.go` | ~20 | ✅ Pass | Stats model |
| `bookmark.go` | ~20 | ✅ Pass | Bookmark model |
| `export.go` | ~20 | ✅ Pass | Export data model |
| `profile.go` | ~20 | ✅ Pass | Profile model |

## Package: `store`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `store.go` | ~80 | ✅ Pass | |
| `repo.go` | ~90 | ✅ Pass | |
| `group.go` | ~130 | ✅ Pass | |
| `release.go` | ~80 | ✅ Pass | Release CRUD |
| `amendment.go` | ~60 | ✅ Pass | Amendment CRUD |
| `history.go` | ~80 | ✅ Pass | History CRUD |
| `stats.go` | ~60 | ✅ Pass | Stats queries |
| `bookmark.go` | ~80 | ✅ Pass | Bookmark CRUD |
| `export.go` | ~60 | ✅ Pass | Export operations |
| `import.go` | ~80 | ✅ Pass | Import operations |
| `profile.go` | ~60 | ✅ Pass | Profile management |
| `cddefault.go` | ~40 | ✅ Pass | CD default storage |
| `template.go` | ~60 | ✅ Pass | Commit template storage |

## Package: `desktop`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `desktop.go` | ~60 | ✅ Pass | |

## Package: `gitutil`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `gitutil.go` | ~120 | ✅ Pass | |
| `latestbranch.go` | ~110 | ✅ Pass | |
| `latestbranchresolve.go` | ~90 | ✅ Pass | |
| `dateformat.go` | ~40 | ✅ Pass | |
| `watchstatus.go` | ~60 | ✅ Pass | Watch-specific git status |

## Package: `verbose`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `verbose.go` | ~78 | ✅ Pass | Fixed `!l.enabled` → positive guard |

## Package: `tests`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `cmd_test/amend_test.go` | ~40 | ✅ Pass | |
| `cmd_test/seowrite_test.go` | ~40 | ✅ Pass | |
| `cmd_test/seowritecreate_test.go` | ~40 | ✅ Pass | |
| `cmd_test/seowritecsv_test.go` | ~40 | ✅ Pass | |
| `cmd_test/seowriteloop_test.go` | ~40 | ✅ Pass | |
| `cmd_test/seowritetemplate_test.go` | ~40 | ✅ Pass | |
| `constants_test/seo_constants_test.go` | ~30 | ✅ Pass | |
| `store_test/template_test.go` | ~30 | ✅ Pass | |

## Root

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `main.go` | ~10 | ✅ Pass | Entry point |

## Audit Totals

| Metric | Count |
|--------|-------|
| Total files audited | 169 |
| Passing | 169 |
| Pending | 0 |

## Recent Fixes

| Date | Issue | Fix |
|------|-------|-----|
| 2026-03-10 | `CmdAmend`, `CmdAmendAlias`, `CmdAmendList`, `CmdAmendListAlias` duplicated in `constants_cli.go` and `constants_amend.go` | Removed from `constants_cli.go` |
| 2026-03-10 | `CmdSEOWrite`, `CmdSEOWriteAlias` duplicated in `constants_cli.go` and `constants_seo.go` | Removed from `constants_cli.go` |

## Wave 2 Changes Applied

| Category | Files Changed | Details |
|----------|--------------|---------|
| File trims (≤200 lines) | 4 | `workflow.go` 416→163, `terminal.go` 223→124, `safe_pull.go` 213→110, `setup.go` 206→131 |
| Negation fixes | 6 | `changelog.go` (3×), `github.go` (2×), `metadata.go`, `semver.go`, `verbose.go`, `config.go` |
| Switch → if/else | 1 | `semver.go` Bump function |
| Constants added | 3 files | 12 git args, 3 bump levels, `SetupGlobalFlag`, `ReleaseTagPrefix` |
| Duplicate fix | 1 file | Removed 6 redeclared constants from `constants_cli.go` |

---

## Constants Inventory

Total: **22 files**, **~580+ constants** + **~11 vars** across 26+ categories.

### `constants.go` — Core Defaults (111 lines)

| Category | Constants |
|----------|-----------|
| Version | `Version` |
| Build-time vars | `RepoPath` (var) |
| Clone modes | `ModeHTTPS`, `ModeSSH` |
| Output formats | `OutputTerminal`, `OutputCSV`, `OutputJSON` |
| URL prefixes | `PrefixHTTPS`, `PrefixSSH` |
| File extensions | `ExtCSV`, `ExtJSON`, `ExtTXT`, `ExtGit` |
| Default file names | `DefaultCSVFile`, `DefaultJSONFile`, `DefaultTextFile`, `DefaultVerboseLogDir`, `DefaultStructureFile`, `DefaultCloneScript`, `DefaultDirectCloneScript`, `DefaultDirectCloneSSHScript`, `DefaultDesktopScript`, `DefaultScanCacheFile`, `DefaultConfigPath`, `DefaultSetupConfigPath`, `DefaultOutputDir`, `DefaultOutputFolder`, `DefaultBranch`, `DefaultDir`, `DefaultVersionFile` |
| Release dir | `DefaultReleaseDir` (var), `DefaultLatestFile` |
| JSON formatting | `JSONIndent` |
| Date display | `DateDisplayLayout`, `DateUTCSuffix` |
| Sort orders | `SortByDate`, `SortByName` |
| Bump levels | `BumpMajor`, `BumpMinor`, `BumpPatch` |
| Permissions | `DirPermission` |
| Safe-pull | `SafePullRetryAttempts`, `SafePullRetryDelayMS`, `WindowsPathWarnThreshold` |
| Verbose | `VerboseLogFileFmt` |

### `constants_git.go` — Git Commands & Arguments (61 lines)

| Category | Constants |
|----------|-----------|
| Core git commands | `GitBin`, `GitClone`, `GitPull`, `GitTag`, `GitCheckout`, `GitPush`, `GitFetch`, `GitBranch`, `GitLog`, `GitForEachRef`, `GitLsRemote`, `GitConfigCmd`, `GitRevParse`, `GitCatFile` |
| Git flags | `GitBranchFlag`, `GitDirFlag`, `GitFFOnlyFlag`, `GitGetFlag`, `GitAbbrevRef`, `GitLsRemoteTags`, `GitTagAnnotateFlag`, `GitTagMessageFlag`, `GitTagListFlag`, `GitBranchListFlag`, `GitCatFileTypeFlag`, `GitArgAll`, `GitArgPrune`, `GitArgRemote`, `GitArgContains`, `GitArgInsideWorkTree` |
| Git refs | `GitHEAD`, `GitOrigin`, `GitOriginPrefix`, `GitCommitPrefix`, `GitRemoteOrigin`, `GitCommitType`, `GitTagGlob` |
| Log format | `GitLogTipFormat`, `GitLogDelimiter`, `GitLogFieldCount`, `GitPointsAtFmt`, `GitRefsRemotesFmt`, `GitFormatRefnameShort`, `HeadPointer`, `ShaDisplayLength` |
| Clone instructions | `CloneInstructionFmt`, `HTTPSFromSSHFmt`, `SSHFromHTTPSFmt` |

### `constants_cli.go` — CLI Commands & Help (~163 lines)

| Category | Constants |
|----------|-----------|
| Command names | `CmdScan`, `CmdClone`, `CmdUpdate`, `CmdUpdateRunner`, `CmdUpdateCleanup`, `CmdVersion`, `CmdHelp`, `CmdDesktopSync`, `CmdPull`, `CmdRescan`, `CmdSetup`, `CmdStatus`, `CmdExec`, `CmdRelease`, `CmdReleaseBranch`, `CmdReleasePending`, `CmdChangelog`, `CmdDoctor`, `CmdLatestBranch`, `CmdList`, `CmdGroup`, `CmdDBReset`, `CmdListVersions`, `CmdRevert`, `CmdRevertRunner`, `CmdListReleases`, `CmdCDCmd` |
| Command aliases | `CmdScanAlias`, `CmdCloneAlias`, `CmdVersionAlias`, `CmdDesktopSyncAlias`, `CmdPullAlias`, `CmdRescanAlias`, `CmdStatusAlias`, `CmdExecAlias`, `CmdReleaseAlias`, `CmdReleaseBranchAlias`, `CmdReleasePendingAlias`, `CmdChangelogAlias`, `CmdLatestBranchAlias`, `CmdListAlias`, `CmdGroupAlias`, `CmdListVersionsAlias`, `CmdListReleasesAlias`, `CmdCDCmdAlias` |
| Group subcommands | `CmdGroupCreate`, `CmdGroupAdd`, `CmdGroupRemove`, `CmdGroupList`, `CmdGroupShow`, `CmdGroupDelete`, `CmdChangelogMD` |
| Clone shorthands | `ShorthandJSON`, `ShorthandCSV`, `ShorthandText` |
| Flag values | `FlagOpenValue`, `FlagJSON`, `FlagLimit`, `FlagSource` |
| Usage/help text | `UsageHeaderFmt`, `HelpUsage`, `HelpCommands`, `HelpScan`, `HelpClone`, `HelpUpdate`, `HelpUpdateCleanup`, `HelpVersion`, `HelpDesktopSync`, `HelpPull`, `HelpRescan`, `HelpSetup`, `HelpStatus`, `HelpExec`, `HelpRelease`, `HelpReleaseBr`, `HelpReleasePend`, `HelpChangelog`, `HelpDoctor`, `HelpLatestBr`, `HelpList`, `HelpGroup`, `HelpDBReset`, `HelpHelp`, `HelpListVersions`, `HelpListReleases`, `HelpRevert`, `HelpScanFlags`, `HelpConfig`, `HelpMode`, `HelpOutput`, `HelpOutputPath`, `HelpOutFile`, `HelpGitHubDesktop`, `HelpOpen`, `HelpQuiet`, `HelpCloneFlags`, `HelpTargetDir`, `HelpSafePull`, `HelpVerbose`, `HelpReleaseFlags`, `HelpAssets`, `HelpCommit`, `HelpRelBranch`, `HelpBump`, `HelpDraft`, `HelpDryRun` |
| Flag descriptions | `FlagDescConfig`, `FlagDescMode`, `FlagDescOutput`, `FlagDescOutFile`, `FlagDescOutputPath`, `FlagDescTargetDir`, `FlagDescSafePull`, `FlagDescGHDesktop`, `FlagDescOpen`, `FlagDescQuiet`, `FlagDescVerbose`, `FlagDescSetupConfig`, `FlagDescDryRun`, `FlagDescAssets`, `FlagDescCommit`, `FlagDescRelBranch`, `FlagDescBump`, `FlagDescDraft`, `FlagDescLatest`, `FlagDescLimit`, `FlagDescOpenChangelog`, `FlagDescLBRemote`, `FlagDescLBAllRemotes`, `FlagDescLBContains`, `FlagDescLBTop`, `FlagDescLBJSON`, `FlagDescLBFormat`, `FlagDescLBNoFetch`, `FlagDescLBSort`, `FlagDescLBFilter`, `FlagDescGroup`, `FlagDescAll`, `FlagDescListVerbose`, `FlagDescGroupDesc`, `FlagDescGroupColor`, `FlagDescConfirm`, `FlagDescSource` |

### `constants_amend.go` — Amend Command (~130 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdAmend`, `CmdAmendAlias`, `CmdAmendList`, `CmdAmendListAlias` |
| Flag names | `FlagAmendName`, `FlagAmendEmail`, `FlagAmendBranch`, `FlagAmendDryRun`, `FlagAmendForcePush` |
| Flag descriptions | `FlagDescAmendName`, `FlagDescAmendEmail`, `FlagDescAmendBranch`, `FlagDescAmendDryRun`, `FlagDescAmendForcePush` |
| Modes | `AmendModeAll`, `AmendModeRange`, `AmendModeHead` |
| Audit | `AmendAuditDir`, `AmendAuditFilePrefix` |
| Messages | `MsgAmendHeader`, `MsgAmendHeaderAll`, `MsgAmendAuthor`, `MsgAmendProgress`, `MsgAmendDone`, `MsgAmendAuditFile`, `MsgAmendAuditDB`, `MsgAmendForcePush`, `MsgAmendWarnPush`, `MsgAmendDryHeader`, `MsgAmendDryLine`, `MsgAmendDrySkip`, `MsgAmendCheckout`, `MsgAmendReturn`, `MsgAmendWarnRewrite` |
| Errors | `ErrAmendNoFlags`, `ErrAmendCheckout`, `ErrAmendListCommits`, `ErrAmendFilter`, `ErrAmendForcePush`, `ErrAmendAuditWrite`, `ErrAmendCommitAmend`, `ErrAmendNoCommits` |
| Help text | `HelpAmend`, `HelpAmendList`, `HelpAmendFlags`, `HelpAmendName`, `HelpAmendEmail`, `HelpAmendBr`, `HelpAmendDry`, `HelpAmendForce` |
| List flag | `FlagAmendListBranch` |
| List messages | `MsgAmendListEmpty`, `MsgAmendListHeader`, `MsgAmendListSeparator`, `MsgAmendListColumns`, `MsgAmendListRowFmt`, `ErrAmendListFailed` |
| SQL | `TableAmendments`, `SQLCreateAmendments`, `SQLInsertAmendment`, `SQLSelectAllAmendments`, `SQLSelectAmendmentsByBranch`, `SQLDropAmendments` |

### `constants_seo.go` — SEO-Write Command (~153 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdSEOWrite`, `CmdSEOWriteAlias`, `CmdCreateTemplate` |
| Flag names | `FlagSEOCSV`, `FlagSEOURL`, `FlagSEOService`, `FlagSEOArea`, `FlagSEOCompany`, `FlagSEOPhone`, `FlagSEOEmail`, `FlagSEOAddress`, `FlagSEOMaxCommits`, `FlagSEOInterval`, `FlagSEOFiles`, `FlagSEORotateFile`, `FlagSEODryRun`, `FlagSEOTemplate`, `FlagSEOCreateTemplate`, `FlagSEOAuthorName`, `FlagSEOAuthorEmail` |
| Flag descriptions | `FlagDescSEOCSV` through `FlagDescSEOAuthorEmail` (17 constants) |
| Defaults | `SEODefaultIntervalMin`, `SEODefaultIntervalMax`, `SEODefaultInterval`, `SEOSeedFile`, `SEOTemplateOutputFile` |
| Placeholders | `PlaceholderService`, `PlaceholderArea`, `PlaceholderURL`, `PlaceholderCompany`, `PlaceholderPhone`, `PlaceholderEmail`, `PlaceholderAddress` |
| Messages | `MsgSEOHeader`, `MsgSEOHeaderUnlimited`, `MsgSEOCommit`, `MsgSEOCommitOpen`, `MsgSEORotation`, `MsgSEORotationOpen`, `MsgSEODone`, `MsgSEODryTitle`, `MsgSEODryDesc`, `MsgSEODryAuthor`, `MsgSEOCreated`, `MsgSEOSeeded`, `MsgSEOGraceful`, `MsgSEOWaiting` |
| Errors | `ErrSEOURLRequired` through `ErrSEODBInsert` (14 constants) |
| Help text | `HelpSEOWrite` through `HelpSEOAuthorEmail` (19 constants) |
| SQL | `TableCommitTemplates`, `SQLCreateCommitTemplates`, `SQLInsertTemplate`, `SQLSelectTemplatesByKind`, `SQLCountTemplates`, `SQLDropCommitTemplates` |
| Template kinds | `TemplateKindTitle`, `TemplateKindDescription` |

### `constants_bookmark.go` — Bookmarks (~68 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdBookmark`, `CmdBookmarkAlias` |
| Subcommands | `CmdBookmarkSave`, `CmdBookmarkList`, `CmdBookmarkRun`, `CmdBookmarkDelete` |
| Help | `HelpBookmark` |
| Messages | `MsgBookmarkSaved`, `MsgBookmarkDeleted`, `MsgBookmarkEmpty`, `MsgBookmarkRunning`, `MsgBookmarkColumns`, `MsgBookmarkRowFmt` |
| Errors | `ErrBookmarkUsage`, `ErrBookmarkSaveUsage`, `ErrBookmarkRunUsage`, `ErrBookmarkDelUsage`, `ErrBookmarkNotFound`, `ErrBookmarkExists`, `ErrBookmarkQuery`, `ErrBookmarkSave`, `ErrBookmarkDelete` |
| SQL | `TableBookmarks`, `SQLCreateBookmarks`, `SQLInsertBookmark`, `SQLSelectAllBookmarks`, `SQLSelectBookmarkByName`, `SQLDeleteBookmark`, `SQLDropBookmarks` |

### `constants_history.go` — Command History (~91 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdHistory`, `CmdHistoryAlias`, `CmdHistoryReset`, `CmdHistoryResetAlias` |
| Help | `HelpHistory`, `HelpHistoryReset` |
| Flag descriptions | `FlagDescDetail`, `FlagDescCommand` |
| Detail levels | `DetailBasic`, `DetailStandard`, `DetailDetailed` |
| Columns | `MsgHistoryColumnsBasic`, `MsgHistoryColumnsStandard`, `MsgHistoryColumnsDetailed`, `MsgHistoryRowBasicFmt`, `MsgHistoryRowStdFmt`, `MsgHistoryRowDetailFmt` |
| Messages | `MsgHistoryEmpty`, `MsgHistoryResetDone`, `MsgHistoryStatusOK`, `MsgHistoryStatusFail` |
| Errors | `ErrHistoryResetFailed`, `ErrHistoryResetNoConfirm`, `ErrHistoryQuery` |
| SQL | `TableCommandHistory`, `SQLCreateCommandHistory`, `SQLInsertHistory`, `SQLUpdateHistory`, `SQLSelectAllHistory`, `SQLSelectHistoryByCommand`, `SQLDeleteAllHistory`, `SQLDropCommandHistory` |

### `constants_stats.go` — Usage Statistics (~66 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdStats`, `CmdStatsAlias` |
| Help | `HelpStats` |
| Flag descriptions | `FlagDescStatsCommand` |
| SQL | `SQLStatsPerCommand`, `SQLStatsForCommand`, `SQLStatsOverall` |
| Display | `MsgStatsHeader`, `MsgStatsSeparator`, `MsgStatsOverallFmt`, `MsgStatsColumns`, `MsgStatsRowFmt`, `MsgStatsEmpty`, `ErrStatsQuery` |

### `constants_export.go` — Export (~22 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdExport`, `CmdExportAlias` |
| Help | `HelpExport` |
| Defaults | `DefaultExportFile`, `FlagDescExportOut` |
| Messages | `MsgExportDone`, `MsgExportFailed` |

### `constants_import.go` — Import (~20 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdImport`, `CmdImportAlias` |
| Help | `HelpImport` |
| Messages | `MsgImportDone`, `MsgImportFailed`, `MsgImportReadFailed`, `MsgImportParseFailed`, `ErrImportNoConfirm`, `MsgImportSkipGroup` |

### `constants_profile.go` — Profiles (~47 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdProfile`, `CmdProfileAlias` |
| Subcommands | `CmdProfileCreate`, `CmdProfileList`, `CmdProfileSwitch`, `CmdProfileDelete`, `CmdProfileShow` |
| Help | `HelpProfile` |
| Defaults | `ProfileConfigFile`, `DefaultProfileName`, `ProfileDBPrefix` |
| Messages/errors | `MsgProfileCreated` through `ErrProfileConfig` (17 constants) |

### `constants_cd.go` — CD Navigation (~43 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdCD`, `CmdCDAlias` |
| Subcommands | `CmdCDRepos`, `CmdCDSetDefault`, `CmdCDClearDefault` |
| Help | `HelpCD` |
| Defaults | `CDDefaultsFile` |
| Messages/errors | `MsgCDMultipleHeader` through `ErrCDDefaultNotFound` (12 constants) |
| Flag descriptions | `FlagDescCDGroup`, `FlagDescCDPick` |

### `constants_watch.go` — Watch Dashboard (~52 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdWatch`, `CmdWatchAlias` |
| Help | `HelpWatch` |
| Defaults | `WatchDefaultInterval`, `WatchMinInterval` |
| ANSI control | `WatchClearScreen` |
| Display | `WatchBannerTop`, `WatchBannerTitle`, `WatchBannerBottom`, `WatchRefreshFmt`, `WatchLastUpdFmt`, `WatchHeaderFmt`, `WatchRowFmt`, `WatchErrorRowFmt`, `WatchSummaryFmt`, `WatchStoppedMsg` |
| Table columns | `WatchTableColumns` (var) |
| Flag descriptions | `FlagDescWatchInterval`, `FlagDescWatchNoFetch`, `FlagDescWatchJSON` |
| Errors | `ErrWatchNoRepos` |

### `constants_diffprofile.go` — Diff Profiles (~28 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdDiffProfiles`, `CmdDiffProfilesAlias` |
| Help | `HelpDiffProfiles` |
| Messages | `MsgDPHeader`, `MsgDPOnlyInHeader`, `MsgDPOnlyInRowFmt`, `MsgDPDiffHeader`, `MsgDPDiffNameFmt`, `MsgDPDiffDetailFmt`, `MsgDPSameFmt`, `MsgDPSameAllHeader`, `MsgDPSameRowFmt`, `MsgDPSummaryFmt`, `MsgDPEmpty` |
| Errors | `ErrDPUsage`, `ErrDPProfileMissing`, `ErrDPOpenFailed` |

### `constants_clone.go` — Clone Progress (~10 lines)

| Category | Constants |
|----------|-----------|
| Progress | `ProgressBeginFmt`, `ProgressDoneFmt`, `ProgressFailFmt`, `ProgressSummaryFmt`, `ProgressDetailFmt` |

### `constants_terminal.go` — UI & Display (~180 lines)

| Category | Constants |
|----------|-----------|
| ANSI colors | `ColorReset`, `ColorGreen`, `ColorRed`, `ColorYellow`, `ColorCyan`, `ColorWhite`, `ColorDim` |
| Status banner | `StatusBannerTop`, `StatusBannerTitle`, `StatusBannerBottom`, `StatusRepoCountFmt` |
| Status indicators | `StatusIconClean`, `StatusIconDirty`, `StatusDash`, `StatusSyncDash`, `StatusStashFmt`, `StatusSyncUpFmt`, `StatusSyncDownFmt`, `StatusSyncBothFmt`, `StatusStagedFmt`, `StatusModifiedFmt`, `StatusUntrackedFmt` |
| Status row formats | `StatusRowFmt`, `StatusMissingFmt`, `StatusHeaderFmt` |
| Summary formats | `SummaryJoinSep`, `SummaryReposFmt`, `SummaryCleanFmt`, `SummaryDirtyFmt`, `SummaryAheadFmt`, `SummaryBehindFmt`, `SummaryStashedFmt`, `SummaryMissingFmt`, `SummarySucceededFmt`, `SummaryFailedFmt`, `StatusFileCountSep`, `TruncateEllipsis` |
| Setup banner | `SetupBannerTop`, `SetupBannerTitle`, `SetupBannerBottom`, `SetupDryRunFmt`, `SetupAppliedFmt`, `SetupSkippedFmt`, `SetupFailedFmt`, `SetupErrorEntryFmt` |
| Changelog display | `ChangelogVersionFmt`, `ChangelogNoteFmt` |
| Exec banner | `ExecBannerTop`, `ExecBannerTitle`, `ExecBannerBottom`, `ExecCommandFmt`, `ExecRepoCountFmt`, `ExecSuccessFmt`, `ExecFailFmt`, `ExecMissingFmt`, `ExecOutputLineFmt`, `ExecSummaryRule` |
| Terminal banner | `TermBannerTop`, `TermBannerTitle`, `TermBannerBottom`, `TermFoundFmt`, `TermReposHeader`, `TermTreeHeader`, `TermCloneHeader`, `TermSeparator`, `TermTableRule` |
| Terminal repo entry | `TermRepoIcon`, `TermPathLine`, `TermCloneLine` |
| Clone help text | `TermCloneStep1`–`TermCloneStep6`, `TermCloneCmd1`–`TermCloneCmd6`, `TermCloneNote` (20 constants) |
| Folder structure MD | `StructureTitle`, `StructureDescription`, `StructureRepoFmt`, `TreeBranch`, `TreeCorner`, `TreePipe`, `TreeSpace` |
| CSV headers | `ScanCSVHeaders` (var), `LatestBranchCSVHeaders` (var) |
| Latest-branch display | `LBTermLatestFmt`, `LBTermRemoteFmt`, `LBTermSHAFmt`, `LBTermDateFmt`, `LBTermSubjectFmt`, `LBTermRefFmt`, `LBTermTopHdrFmt`, `LBTermRowFmt` |
| Latest-branch table | `LatestBranchTableColumns` (var), `StatusTableColumns` (var) |

### `constants_messages.go` — User Messages & Errors (~195 lines)

| Category | Constants |
|----------|-----------|
| Notes | `NoteNoRemote`, `UnknownRepoName` |
| GitHub Desktop | `GitHubDesktopBin`, `OSWindows`, `MsgDesktopNotFound`, `MsgDesktopAdded`, `MsgDesktopFailed`, `MsgDesktopSummary` |
| Latest-branch messages | `MsgLatestBranchFetching`, `MsgLatestBranchFetchWarning`, `LBUnknownBranch` |
| Generic errors | `ErrGenericFmt`, `ErrBareFmt` |
| OS platform | `OSDarwin` |
| OS commands | `CmdExplorer`, `CmdOpen`, `CmdXdgOpen`, `CmdWindowsShell`, `CmdArgSlashC`, `CmdArgStart`, `CmdArgEmpty` |
| Desktop sync errors | `ErrDesktopReadFailed`, `ErrDesktopParseFailed`, `ErrNoAbsPath` |
| Dispatch errors | `ErrUnknownCommand`, `ErrUnknownGroupSub` |
| Version display | `MsgVersionFmt` |
| CLI messages | `MsgFoundRepos`, `MsgCSVWritten`, `MsgJSONWritten`, `MsgTextWritten`, `MsgStructureWritten`, `MsgCloneScript`, `MsgDirectClone`, `MsgDirectCloneSSH`, `MsgDesktopScript`, `MsgCloneComplete`, `MsgAutoSafePull`, `MsgOpenedFolder`, `MsgVerboseLogFile`, `MsgDesktopSyncStart`, `MsgDesktopSyncSkipped`, `MsgDesktopSyncAdded`, `MsgDesktopSyncFailed`, `MsgDesktopSyncDone`, `MsgNoOutputDir`, `MsgNoJSONFile`, `MsgFailedClones`, `MsgFailedEntry`, `MsgPullStarting`, `MsgPullSuccess`, `MsgPullFailed`, `MsgPullAvailable`, `MsgPullListEntry`, `WarnVerboseLogFailed`, `MsgRescanReplay`, `MsgScanCacheSaved`, `MsgDBUpsertDone`, `MsgDBUpsertFailed`, `MsgUpdateStarting`, `MsgUpdateRepoPath`, `MsgUpdateVersion` |
| List/group messages | `MsgListHeader`, `MsgListSeparator`, `MsgListRowFmt`, `MsgListVerboseFmt`, `MsgListEmpty`, `MsgGroupCreated`, `MsgGroupDeleted`, `MsgGroupAdded`, `MsgGroupRemoved`, `MsgGroupHeader`, `MsgGroupRowFmt`, `MsgGroupShowHeader`, `MsgGroupShowRowFmt`, `MsgGroupEmpty`, `ErrGroupNameReq`, `ErrGroupUsage`, `ErrGroupSlugReq`, `ErrListDBFailed`, `ErrNoDatabase`, `MsgDBResetDone`, `ErrDBResetFailed`, `ErrDBResetNoConfirm` |
| Latest-branch errors | `ErrLatestBranchNotRepo`, `ErrLatestBranchNoRefs`, `ErrLatestBranchNoRefsAll`, `ErrLatestBranchNoCommits`, `ErrLatestBranchNoMatch` |
| CLI errors | `ErrSourceRequired`, `ErrCloneUsage`, `ErrShorthandNotFound`, `ErrConfigLoad`, `ErrScanFailed`, `ErrCloneFailed`, `ErrOutputFailed`, `ErrCreateDir`, `ErrCreateFile`, `ErrNoRepoPath`, `ErrUpdateFailed`, `ErrPullSlugRequired`, `ErrPullUsage`, `ErrPullLoadFailed`, `ErrPullNotFound`, `ErrPullNotRepo`, `ErrRescanNoCache`, `ErrSetupLoadFailed`, `ErrStatusLoadFailed`, `ErrExecUsage`, `ErrExecLoadFailed`, `ErrReleaseVersionRequired`, `ErrReleaseUsage`, `ErrReleaseBranchUsage`, `ErrReleaseAlreadyExists`, `ErrReleaseTagExists`, `ErrReleaseBranchNotFound`, `ErrReleaseCommitNotFound`, `ErrReleaseInvalidVersion`, `ErrReleaseBumpNoLatest`, `ErrReleaseBumpConflict`, `ErrReleaseCommitBranch`, `ErrReleasePushFailed`, `ErrReleaseVersionLoad`, `ErrReleaseMetaWrite`, `ErrChangelogRead`, `ErrChangelogVersionNotFound`, `ErrChangelogOpen` |
| List-versions errors | `ErrListVersionsNoTags` |
| Revert messages | `MsgRevertCheckout`, `MsgRevertStarting`, `MsgRevertDone`, `ErrRevertUsage`, `ErrRevertTagNotFound`, `ErrRevertCheckoutFailed`, `ErrRevertFailed`, `RevertScriptLogExec`, `RevertScriptLogExit` |
| Releases listing | `ReleaseGlob` |

### `constants_release.go` — Release & Setup (37 lines)

| Category | Constants |
|----------|-----------|
| Setup sections | `SetupSectionDiff`, `SetupSectionMerge`, `SetupSectionAlias`, `SetupSectionCred`, `SetupSectionCore`, `SetupGlobalFlag` |
| Release messages | `MsgReleaseStart`, `MsgReleaseBranch`, `MsgReleaseTag`, `MsgReleasePushed`, `MsgReleaseMeta`, `MsgReleaseLatest`, `MsgReleaseAttach`, `MsgReleaseChangelog`, `MsgReleaseReadme`, `MsgReleaseDryRun`, `MsgReleaseComplete`, `MsgReleaseBranchStart`, `MsgReleaseVersionRead`, `MsgReleaseBumpResult`, `MsgReleaseSwitchedBack`, `MsgReleasePendingNone`, `MsgReleasePendingFound`, `MsgReleasePendingFailed` |
| Release paths | `ReleaseBranchPrefix`, `ChangelogFile`, `ReadmeFile`, `ReleaseTagPrefix` |

### `constants_store.go` — Database & SQL (~122 lines)

| Category | Constants |
|----------|-----------|
| DB location | `DBDir`, `DBFile` |
| Table names | `TableRepos`, `TableGroups`, `TableGroupRepo` |
| Schema DDL | `SQLCreateRepos`, `SQLCreateGroups`, `SQLCreateGroupRepos`, `SQLEnableFK`, `SQLCreateAbsPathIndex` |
| Repo queries | `SQLUpsertRepo`, `SQLSelectAllRepos`, `SQLSelectRepoBySlug`, `SQLSelectRepoByPath`, `SQLUpsertRepoByPath` |
| Group queries | `SQLInsertGroup`, `SQLSelectAllGroups`, `SQLSelectGroupByName`, `SQLDeleteGroup`, `SQLInsertGroupRepo`, `SQLDeleteGroupRepo`, `SQLSelectGroupRepos`, `SQLCountGroupRepos` |
| Reset queries | `SQLDropGroupRepos`, `SQLDropGroups`, `SQLDropRepos` |
| Store errors | `ErrDBOpen`, `ErrDBMigrate`, `ErrDBUpsert`, `ErrDBQuery`, `ErrDBNoMatch`, `ErrDBCreateDir`, `ErrDBGroupCreate`, `ErrDBGroupQuery`, `ErrDBGroupAdd`, `ErrDBGroupRemove`, `ErrDBGroupDelete`, `ErrDBGroupNone`, `ErrDBGroupExists` |

### `constants_doctor.go` — Diagnostics (~91 lines)

| Category | Constants |
|----------|-----------|
| Doctor banners | `DoctorBannerFmt`, `DoctorBannerRule`, `DoctorIssuesFmt`, `DoctorFixPathTip`, `DoctorAllPassed`, `DoctorFixBannerFmt` |
| Doctor path/deploy | `DoctorActivePathFmt`, `DoctorDeployedFmt`, `DoctorSyncingFmt`, `DoctorRetryFmt`, `DoctorRenamedMsg`, `DoctorKillingMsg`, `DoctorKilledFmt` |
| Doctor sync failures | `DoctorSyncFailTitle`, `DoctorSyncFailDetail`, `DoctorSyncFailFix1`, `DoctorSyncFailFix2Fmt` |
| Doctor check results | `DoctorFixFlagDesc`, `DoctorOKPathFmt`, `DoctorWarnSyncFmt`, `DoctorNotOnPath`, `DoctorNoSync`, `DoctorAddPathFix`, `DoctorCannotResolve`, `DoctorAlreadySynced`, `DoctorVersionsMatch` |
| Doctor RepoPath | `DoctorRepoPathMissing`, `DoctorRepoPathDetail`, `DoctorRepoPathFix`, `DoctorRepoPathOKFmt` |
| Doctor PATH binary | `DoctorPathBinaryFmt`, `DoctorPathMissTitle`, `DoctorPathMissDetail`, `DoctorPathMissFix` |
| Doctor deploy binary | `DoctorDeployReadFail`, `DoctorDeployReadDet`, `DoctorNoDeployPath`, `DoctorNoDeployDet`, `DoctorDeployNotFound`, `DoctorDeployRunFix`, `DoctorDeployOKFmt` |
| Doctor git/go checks | `DoctorGitMissTitle`, `DoctorGitMissDetail`, `DoctorGitOKFmt`, `DoctorGitOKPathFmt`, `DoctorGoWarn`, `DoctorGoOKFmt`, `DoctorGoOKPathFmt` |
| Doctor changelog | `DoctorChangelogWarn`, `DoctorChangelogOK` |
| Doctor version mismatch | `DoctorVersionMismatch`, `DoctorVMismatchFmt`, `DoctorVMismatchFix`, `DoctorDeployMismatch`, `DoctorDMismatchFmt`, `DoctorDMismatchFix`, `DoctorBinariesDiffer`, `DoctorBDifferFmt`, `DoctorBDifferFix`, `DoctorSourceOKFmt` |
| Doctor resolve | `DoctorResolveNoRepo`, `DoctorResolveNoRead`, `DoctorResolveNoDeploy`, `DoctorResolveNotFound`, `DoctorDefaultBinary` |
| Doctor binary lookup | `GitMapBin`, `GoBin`, `GoVersionArg`, `PowershellConfigFile`, `JSONKeyDeployPath`, `JSONKeyBinaryName`, `BackupSuffix`, `GitMapSubdir` |
| Doctor format markers | `DoctorOKFmt`, `DoctorIssueFmt`, `DoctorFixFmt`, `DoctorWarnFmt`, `DoctorDetail` |

### `constants_update.go` — Self-Update & PowerShell (~119 lines)

| Category | Constants |
|----------|-----------|
| Update file patterns | `UpdateCopyFmt`, `UpdateCopyGlob`, `UpdateScriptGlob` |
| Update flags | `FlagVerbose` |
| Update UI messages | `MsgUpdateActive`, `MsgUpdateCleanStart`, `MsgUpdateCleanDone`, `MsgUpdateCleanNone`, `MsgUpdateTempRemoved`, `MsgUpdateOldRemoved`, `UpdateRunnerLogStart`, `UpdateScriptLogExec`, `UpdateScriptLogExit` |
| Update errors | `ErrUpdateExecFind`, `ErrUpdateCopyFail` |
| Update PS script | `UpdatePSHeader`, `UpdatePSDeployDetect`, `UpdatePSVersionBefore`, `UpdatePSRunUpdate`, `UpdatePSVersionAfter`, `UpdatePSVerify`, `UpdatePSPostActions` |
| Revert PS script | `RevertPSHeader`, `RevertPSBuild`, `RevertPSPostActions` |
| Backup glob | `OldBackupGlob` |
| PowerShell args | `PSBin`, `PSExecPolicy`, `PSBypass`, `PSNoProfile`, `PSNoLogo`, `PSFile`, `PSNonInteractive`, `PSCommand` |

### `constants_gomod.go` — GoMod Rename (~100 lines)

| Category | Constants |
|----------|-----------|
| Commands | `CmdGoMod`, `CmdGoModAlias` |
| Help text | `HelpGoMod`, `HelpGoModFlags`, `HelpGoModDry`, `HelpGoModNoMrg`, `HelpGoModNoTdy`, `HelpGoModVerb`, `HelpGoModExt` |
| Flag names | `FlagGoModDryRun`, `FlagGoModNoMerge`, `FlagGoModNoTidy`, `FlagGoModExt` |
| Flag descriptions | `FlagDescGoModDryRun`, `FlagDescGoModNoMerge`, `FlagDescGoModNoTidy`, `FlagDescGoModExt` |
| File constants | `GoModFile`, `GoModModuleLine`, `GoFileExt` |
| Excluded dirs | `GoModExcludeDirs` (var) |
| Branch prefixes | `GoModFeaturePrefix`, `GoModBackupPrefix` |
| Messages | `MsgGoModSummary`, `MsgGoModOld`, `MsgGoModNew`, `MsgGoModFiles`, `MsgGoModBackupBranch`, `MsgGoModFeatureBranch`, `MsgGoModMergedInto`, `MsgGoModLeftOn`, `MsgGoModVerboseFile`, `MsgGoModDryHeader`, `MsgGoModDryOld`, `MsgGoModDryNew`, `MsgGoModDryFiles`, `MsgGoModDryFile`, `MsgGoModNoImports`, `MsgGoModTidyWarn`, `MsgGoModNothingRename` |
| Errors | `ErrGoModUsage`, `ErrGoModNoFile`, `ErrGoModNoModule`, `ErrGoModNotRepo`, `ErrGoModDirtyTree`, `ErrGoModBranchExists`, `ErrGoModMergeConflict`, `ErrGoModReadFailed`, `ErrGoModWriteFailed`, `ErrGoModCommitFailed` |
| Git args | `GitAdd`, `GitAddAll`, `GitCommit`, `GitCommitMsg`, `GitMerge`, `GitMergeNoFF`, `GitStatusShort`, `GitStatus` |
| Commit format | `GoModCommitMsgFmt`, `GoModMergeMsgFmt` |
