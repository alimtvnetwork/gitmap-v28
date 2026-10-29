# Context Menu and Retirement Specification

**Version:** 6.630.0
**Updated:** 2026-10-29

## 1. Scope

This document specifies the file-class context menu and the later catalog retirement for sequence 354 (`gitmap-ui-and-functionality-improvements`). It is the behavior contract. It does not change Go sources. Another agent owns the context-menu Go files, and those files stay untouched until that owner releases them.

Implementation of the menu and diff behavior is blocked. See the pending plan `.ai-memory/plans/pending/354-gitmap-ui-and-functionality-improvements.md`, subtask ST-06.

## 2. Current file menu

`ctxFileMenu` in `cli/cmdinstall/installctxentries.go` registers the file-class cascade. The items, in order, are:

1. Diff
2. Diff later
3. Diff with previous version
4. Show log
5. Blame
6. Commit to main

The first five items are the diff group (inspection). Commit to main is the commit group.

`leafCommands` in `cli/cmdinstall/installctx.go` writes the verb default, the `\command` subkey, and the optional Icon and Extended values. It writes no separator. The file cascade therefore has no Tortoise-style gap between the diff group and the commit group.

## 3. Current working-directory command

`commandTemplate` in `cli/cmdinstall/installctx.go` builds the PowerShell command stored on each verb. Every branch calls `Set-Location` on `%V`.

On the file class, `%V` is the clicked file, not the folder that contains it. `Set-Location` accepts only a directory. When `%V` is a file, PowerShell rejects the location change and prints that the path does not exist. That error does not stop the script. The script continues and launches `gitmap diff`.

The terminal and silent branches guard with `Test-Path` before `Set-Location`. `Test-Path` succeeds for a file, so the guard still enters `Set-Location` and the same rejection occurs. The prefill branch calls `Set-Location` with no path-type check.

## 4. Current single-file diff

`runSinglePathDiff` in `cli/cmddiff/single_diff.go` launches the external compare tool whenever a tool is installed. `launchExternalToolDiff` always passes one path. It does not ask whether that file differs from `HEAD`.

`LaunchToolWithPaths` and `LaunchToolWithPathsAsync` in `cli/cmdgit/compare_launcher.go` print a line that starts with `Launching` and names the tool. When the selected tool is Beyond Compare, that name is `bcomp`. The async path used by single-file diff prints the launching line and starts `bcomp` even when the file has no changes against `HEAD`.

## 5. Required menu behavior

When the context-menu Go files are free to edit, the file cascade must show a registry separator between the diff group and the commit group.

The separator is the Windows shell `CommandFlags` value on the Commit to main verb (`60_commit`): `REG_DWORD` `0x00000020` (`ECF_SEPARATORBEFORE`). `leafCommands` must write that value for the commit verb and must leave the five diff-group verbs without a leading separator. Explorer then draws the Tortoise-style gap immediately before Commit to main, after Blame.

Reinstalling the context menu must rewrite the file-class verbs so an existing install gains the same gap.

## 6. Required working-directory behavior

`commandTemplate` must call `Set-Location` only when `%V` is a directory. Apply that rule on the prefill branch and on the terminal and silent branches.

A positive directory test is required before `Set-Location`. `Test-Path -LiteralPath '%V' -PathType Container` is the check. When `%V` is a file, skip `Set-Location` and continue with the gitmap invocation. PowerShell must not print that the path does not exist for a file that is present.

## 7. Required no-changes behavior

`runSinglePathDiff` must compare the file with `HEAD` before any external launch. When the file has no diff against `HEAD`, print one clear no-changes line that names the path and states that there are no changes against `HEAD`, then return. Do not call the Beyond Compare launcher. `cli/cmdgit/compare_launcher.go` must not print `Launching bcomp` for that file, because the launch must not start.

When the file does differ from `HEAD`, keep the current one-path Beyond Compare launch and the existing launching line.

A missing git repository keeps the current non-repository message and must not launch Beyond Compare.

## 8. Ownership block

ST-06 stays blocked while these files are owned elsewhere:

- `cli/cmdinstall/installctxentries.go`
- `cli/cmdinstall/installctx.go`
- `cli/cmddiff/single_diff.go`
- `cli/cmdgit/compare_launcher.go`

This specification does not edit them.

## 9. Catalog retirement

Generic UI catalogs are scheduled for removal from the working tree later. Removal means deleting the directories in a normal commit. It does not mean a history rewrite.

Inclusive range to remove later:

- `02-spec/21-app/340-next-100-ui-improvements`
- `02-spec/21-app/341-next-100-ui-improvements`
- `02-spec/21-app/342-next-100-ui-improvements`
- `02-spec/21-app/343-next-100-ui-improvements`
- `02-spec/21-app/344-next-100-ui-improvements`
- `02-spec/21-app/345-next-100-ui-improvements`
- `02-spec/21-app/346-next-100-ui-improvements`

That deletion is plan subtask ST-07. It is blocked until the heatmap no longer lists those catalogs as active work.

## 10. Time-restriction record (blocked)

Time-restriction content lives in `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling`. The heatmap lists that task as active.

Record 349 as blocked. Keep the directory. Do not delete it in the catalog retirement. Do not rewrite git history to drop it. ST-07 must not treat 349 as part of the 340–346 range.
