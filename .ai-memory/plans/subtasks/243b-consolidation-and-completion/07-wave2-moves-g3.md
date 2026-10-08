# Subtask 07 — Wave 2 moves, group G3 (cmdopen→cmdrepo, NEW packages)

## Objective
Move implementation files out of `cli/cmd/` into NEW owning `cmdX` packages per the move map
(`.ai-memory/plans/subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md`,
Wave 3). MECHANICAL move wave: copy with package-clause swap, verify, delete source.
The lead rewires dispatch callers centrally afterwards — you do NOT touch dispatch.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section.

## Owned files
ONLY the source files listed below (in `cli/cmd/`) and their destinations. Do NOT touch any other file.

- cli/cmd/open.go → cli/cmdopen/open.go
- cli/cmd/open_opener.go → cli/cmdopen/open_opener.go
- cli/cmd/open_target.go → cli/cmdopen/open_target.go
- cli/cmd/orphans.go → cli/cmdorphans/orphans.go
- cli/cmd/dopending.go → cli/cmdpending/dopending.go
- cli/cmd/dopendingretry.go → cli/cmdpending/dopendingretry.go
- cli/cmd/pending.go → cli/cmdpending/pending.go
- cli/cmd/pending_commits_cmd.go → cli/cmdpending/pending_commits_cmd.go
- cli/cmd/pending_commits_types.go → cli/cmdpending/pending_commits_types.go
- cli/cmd/pendingclear.go → cli/cmdpending/pendingclear.go
- cli/cmd/pendingtaskhelper.go → cli/cmdpending/pendingtaskhelper.go
- cli/cmd/pin_cmd.go → cli/cmdpin/pin_cmd.go
- cli/cmd/probe.go → cli/cmdprobe/probe.go
- cli/cmd/probeflags.go → cli/cmdprobe/probeflags.go
- cli/cmd/proberender.go → cli/cmdprobe/proberender.go
- cli/cmd/probereport.go → cli/cmdprobe/probereport.go
- cli/cmd/profile.go → cli/cmdprofiles/profile.go
- cli/cmd/profileops.go → cli/cmdprofiles/profileops.go
- cli/cmd/profiles_cmd.go → cli/cmdprofiles/profiles_cmd.go
- cli/cmd/profiles_ops.go → cli/cmdprofiles/profiles_ops.go
- cli/cmd/profileutil.go → cli/cmdprofiles/profileutil.go
- cli/cmd/projectrepos.go → cli/cmdprojectrepos/projectrepos.go
- cli/cmd/projectreposoutput.go → cli/cmdprojectrepos/projectreposoutput.go
- cli/cmd/projectreposrender.go → cli/cmdprojectrepos/projectreposrender.go
- cli/cmd/prune.go → cli/cmdprune/prune.go
- cli/cmd/pruneops.go → cli/cmdprune/pruneops.go
- cli/cmd/reconcile.go → cli/cmdreconcile/reconcile.go
- cli/cmd/reconcile_cmd.go → cli/cmdreconcile/reconcile_cmd.go
- cli/cmd/reconcile_db.go → cli/cmdreconcile/reconcile_db.go
- cli/cmd/reconcile_prompt.go → cli/cmdreconcile/reconcile_prompt.go
- cli/cmd/reconcile_types.go → cli/cmdreconcile/reconcile_types.go
- cli/cmd/recreate_repo.go → cli/cmdrecreate/recreate_repo.go
- cli/cmd/regoldens.go → cli/cmdregoldens/regoldens.go
- cli/cmd/regoldens_diff.go → cli/cmdregoldens/regoldens_diff.go
- cli/cmd/regoldens_diff_print.go → cli/cmdregoldens/regoldens_diff_print.go
- cli/cmd/regoldens_diff_scope.go → cli/cmdregoldens/regoldens_diff_scope.go
- cli/cmd/regoldens_dryrun.go → cli/cmdregoldens/regoldens_dryrun.go
- cli/cmd/regoldens_exec.go → cli/cmdregoldens/regoldens_exec.go
- cli/cmd/regoldens_precheck.go → cli/cmdregoldens/regoldens_precheck.go
- cli/cmd/regoldens_validate.go → cli/cmdregoldens/regoldens_validate.go
- cli/cmd/clearreleasejson.go → cli/cmdrelease/clearreleasejson.go
- cli/cmd/release.go → cli/cmdrelease/release.go
- cli/cmd/release_notes_opts.go → cli/cmdrelease/release_notes_opts.go
- cli/cmd/release_scan_commits.go → cli/cmdrelease/release_scan_commits.go
- cli/cmd/release_tools.go → cli/cmdrelease/release_tools.go
- cli/cmd/releasealias.go → cli/cmdrelease/releasealias.go
- cli/cmd/releasealias_git.go → cli/cmdrelease/releasealias_git.go
- cli/cmd/releaseargs.go → cli/cmdrelease/releaseargs.go
- cli/cmd/releaseautobump.go → cli/cmdrelease/releaseautobump.go
- cli/cmd/releaseautoregister.go → cli/cmdrelease/releaseautoregister.go
- cli/cmd/releasebranch.go → cli/cmdrelease/releasebranch.go
- cli/cmd/releasepending.go → cli/cmdrelease/releasepending.go
- cli/cmd/releasepersist.go → cli/cmdrelease/releasepersist.go
- cli/cmd/releaserebase.go → cli/cmdrelease/releaserebase.go
- cli/cmd/releasepull.go → cli/cmdrelease/releasepull.go
- cli/cmd/releaserecentclone.go → cli/cmdrelease/releaserecentclone.go
- cli/cmd/releasescan.go → cli/cmdrelease/releasescan.go
- cli/cmd/releaseself.go → cli/cmdrelease/releaseself.go
- cli/cmd/releaseundo.go → cli/cmdrelease/releaseundo.go
- cli/cmd/releaseundorange.go → cli/cmdrelease/releaseundorange.go
- cli/cmd/remediation_box.go → cli/cmdremediation/remediation_box.go
- cli/cmd/remediation_local.go → cli/cmdremediation/remediation_local.go
- cli/cmd/remediation_matcher.go → cli/cmdremediation/remediation_matcher.go
- cli/cmd/remediation_state_ops.go → cli/cmdremediation/remediation_state_ops.go
- cli/cmd/remediation_suggest.go → cli/cmdremediation/remediation_suggest.go
- cli/cmd/replace.go → cli/cmdreplace/replace.go
- cli/cmd/replace_classify.go → cli/cmdreplace/replace_classify.go
- cli/cmd/replace_dashn.go → cli/cmdreplace/replace_dashn.go
- cli/cmd/replaceapply.go → cli/cmdreplace/replaceapply.go
- cli/cmd/replaceaudit.go → cli/cmdreplace/replaceaudit.go
- cli/cmd/replaceflags.go → cli/cmdreplace/replaceflags.go
- cli/cmd/replaceversion.go → cli/cmdreplace/replaceversion.go
- cli/cmd/replaceversionrun.go → cli/cmdreplace/replaceversionrun.go
- cli/cmd/replacewalk.go → cli/cmdreplace/replacewalk.go
- cli/cmd/repo_cmd_dispatch.go → cli/cmdrepo/repo_cmd_dispatch.go
- cli/cmd/repo_create_init.go → cli/cmdrepo/repo_create_init.go
- cli/cmd/repo_create_params.go → cli/cmdrepo/repo_create_params.go
- cli/cmd/repo_create_profile.go → cli/cmdrepo/repo_create_profile.go
- cli/cmd/repo_create_remote.go → cli/cmdrepo/repo_create_remote.go
- cli/cmd/repo_create_report.go → cli/cmdrepo/repo_create_report.go
- cli/cmd/repo_create_slug.go → cli/cmdrepo/repo_create_slug.go
- cli/cmd/repo_create_smart.go → cli/cmdrepo/repo_create_smart.go
- cli/cmd/repo_create_tokens.go → cli/cmdrepo/repo_create_tokens.go
- cli/cmd/repo_db_dispatch.go → cli/cmdrepo/repo_db_dispatch.go
- cli/cmd/repo_db_ops.go → cli/cmdrepo/repo_db_ops.go
- cli/cmd/repo_recreate.go → cli/cmdrepo/repo_recreate.go

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
