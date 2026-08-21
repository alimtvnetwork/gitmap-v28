# Plan: Coding Guideline Audit & Enforcement (v4)

Slug: 01-coding-guideline-fixes
Steps: 152
Status: pending
Created: 2026-08-22

## Context
Massive audit of the codebase to enforce v1.4.5 coding guidelines. Target violations:
- Inverted booleans (!isSuccess)
- Magic strings and numbers
- Swallowed errors (empty catch {})
- Monolithic functions (> 15 lines)
- Nested if statements

## Strategy
3 Concurrent Sub-Agents will be spawned to tackle the subtasks.

## 1. Subtask: Inverted Booleans (Go)
1. gitmap/cmd/installctx_linux_e2e_test.go: line 213, extract !hasZenity to isZenityMissing := !hasZenity
2. gitmap/cmd/prettyflag.go: line 75, extract !hasPrettyPrefix(arg) to isMissingPrefix := !hasPrettyPrefix(arg)
3. gitmap/cmd/prettyflag.go: line 101, extract !hasValue to isValueMissing := !hasValue
4. gitmap/cmd/pull.go: line 155, extract !isGitRepoCWD() to isNotGitRepo := !isGitRepoCWD()
5. gitmap/cmd/pullreleasecd.go: line 135, extract !isPRCURL(token) to isNonPRCURL := !isPRCURL(token)
6. gitmap/cmd/pullreleasecd.go: line 138, extract !isPRCURL(token) to isNonPRCURL
7. gitmap/cmd/pullreleasecd.go: line 141, extract !isPRCURL(token) to isNonPRCURL
8. gitmap/cmd/push.go: line 135, extract !isGitRepoCWD() to isNotGitRepo := !isGitRepoCWD()
9. gitmap/cmd/push.go: line 219, extract !isNonFastForwardRejection(stderr) to isOtherRejection := !isNonFastForwardRejection(stderr)
10. gitmap/cmd/reclone_confirm.go: line 49, extract !isStdinInteractive() to isNonInteractive := !isStdinInteractive()
11. gitmap/cmd/reclone_validate.go: line 114, extract !isPlausibleGitURL(picked) to isImplausibleURL := !isPlausibleGitURL(picked)
12. gitmap/cmd/reclone_validate.go: line 174, extract !isSchemeChar(ch) to isNonSchemeChar := !isSchemeChar(ch)
13. gitmap/cmd/regoldens_diff.go: line 42, extract !isGitWorkingTree() to isNotWorkingTree := !isGitWorkingTree()
14. gitmap/cmd/regoldens_diff.go: line 100, extract !isGoldenFixturePath(path) to isNonGoldenPath := !isGoldenFixturePath(path)
15. gitmap/cmd/regoldens_diff.go: line 157, extract !isGoldenFixturePath(fields[2]) to isNonGoldenPath := !isGoldenFixturePath(fields[2])
16. gitmap/cmd/release.go: line 59, extract !shouldAutoBumpMinor(...) to mustSkipAutoBump := !shouldAutoBumpMinor(...)
17. gitmap/cmd/releasealias_git.go: line 32, extract !isWorkingTreeDirty(target) to isWorkingTreeClean := !isWorkingTreeDirty(target)
18. gitmap/cmd/replacewalk_test.go: line 15, extract !isExcludedDir(name) to isIncludedDir := !isExcludedDir(name)
19. gitmap/cmd/reporeclone.go: line 74, extract !isGitRepoDir(cwd) to isNotGitRepo := !isGitRepoDir(cwd)
20. gitmap/cmd/reporeclone.go: line 90, extract !isGitRepoDir(abs) to isNotGitRepo := !isGitRepoDir(abs)
21. gitmap/cmd/reporeclone.go: line 149, extract !isStdinTTY() to isNonTTY := !isStdinTTY()
22. gitmap/cmd/schemaregistry_contract_test.go: line 105, extract !isSchemaAccepted("demo", 3) to isSchemaRejected := !isSchemaAccepted("demo", 3)
23. gitmap/cmd/schemaregistry_contract_test.go: line 126, extract !shouldUpdateSchema("from-flag") to mustRetainSchema := !shouldUpdateSchema("from-flag")
24. gitmap/cmd/selfinstall.go: line 207, extract !isConcreteShellFamily(tok) to isUnknownShell := !isConcreteShellFamily(tok)
25. gitmap/cmd/selfuninstallhandoff.go: line 50, extract !isWindows() to isNonWindows := !isWindows()
26. gitmap/cmd/selfuninstallparts.go: line 52, extract !isGitmapArtifact(e.Name()) to isForeignArtifact := !isGitmapArtifact(e.Name())
27. gitmap/cmd/sshbind.go: line 19, extract !isGitRepoCWD() to isNotGitRepo := !isGitRepoCWD()
28. gitmap/cmd/task_unit_test.go: line 97, extract !isIgnored(...) to isIncluded := !isIgnored(...)
29. gitmap/cmd/templatescli.go: line 118, extract !isValidKindFilter(kindFilter) to isInvalidKind := !isValidKindFilter(kindFilter)
30. gitmap/cmd/updatecleanup_extra.go: line 42, extract !isRemovableDriveRootShim(...) to isNonRemovableShim := !isRemovableDriveRootShim(...)
31. gitmap/cmd/updatedebugwindows_json.go: line 51, extract !isDebugWindowsJSONRequested() to isReleaseMode := !isDebugWindowsJSONRequested()
32. gitmap/cmd/updaterepo.go: line 123, extract !canPromptForRepoPath() to mustSkipPrompt := !canPromptForRepoPath()
33. gitmap/cmd/whoami.go: line 21, extract !isGitRepoCWD() to isNotGitRepo := !isGitRepoCWD()
34. gitmap/committransfer/replay.go: line 75, extract !hasStagedChanges(...) to isTreeClean := !hasStagedChanges(...)
35. gitmap/detector/detector.go: line 61, extract !isInterestingFile(name) to isUninteresting := !isInterestingFile(name)

## 2. Subtask: Inverted Booleans (TS)
36. src/components/docs/DocsTooltip.tsx: line 65, extract !isValidElement(child) to isInvalid := !isValidElement(child)
37. src/components/docs/TerminalDemo.tsx: line 39, extract !isPlaying to isPaused := !isPlaying
38. src/components/spec/SpecSectionCard.tsx: line 60, extract !isCollapsed to isExpanded := !isCollapsed
39. src/components/ui/carousel.tsx: line 189, extract !canScrollPrev to cannotScrollPrev := !canScrollPrev
40. src/components/ui/carousel.tsx: line 217, extract !canScrollNext to cannotScrollNext := !canScrollNext
41. src/hooks/use-mobile.tsx: line 18, extract !!isMobile to return isMobile directly or cast appropriately
42. src/pages/ReleaseVersion.tsx: line 29, extract !isValid to isInvalid := !isValid

## 3. Subtask: Swallowed Errors (PowerShell)
43. fix-repo.ps1: line 139, add explicit Write-Error log to empty catch block
44. fix-repo.ps1: line 260, add explicit Write-Error log to empty catch block
45. gitmap/scripts/Get-LastRelease.ps1: lines 47, 69, 93, 112, add log to empty catch blocks
46. gitmap/scripts/install.ps1: lines 93, 99, 130, 175, 196, 359, 370, 417, 530, 651, 684, 756, 773, 1162, 1191, 1255, 1279, 1421, 1453 add log to empty catch blocks
47. gitmap/scripts/release-version.ps1: lines 203, 300, 314, 355, 389, 438, 453, 462, 513 add log to empty catch blocks
48. install-quick.ps1: add log to empty catch blocks (approx 8 blocks)
49. install.ps1: add log to empty catch blocks (approx 15 blocks)
50. run.ps1: add log to empty catch blocks (approx 20 blocks)
51. uninstall-quick.ps1: add log to empty catch blocks (approx 6 blocks)

## 4. Subtask: Long Functions & Nested Ifs (Go)
52-100. Audit all func declarations in gitmap/cmd/*.go exceeding 15 lines. Extract guard clauses and internal logic into separate helpers to meet the 8-15 line target. (Granular mapping delegated to Agent 3).

## 5. Subtask: Long Functions & Nested Ifs (TS/React)
101-152. Audit all TSX components in src/components/ and ensure all inline effects and functions comply with the 15-line hard cap. (Granular mapping delegated to Agent 3).
