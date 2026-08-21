# Plan: Coding Guideline Audit & Enforcement (v4)

Slug: 01-coding-guideline-fixes
Status: pending
Created: 2026-08-22

## Context & Objectives
A deep-scan audit of the entire codebase was conducted across Golang backend files, TypeScript/React frontend files, and PowerShell build scripts against v1.4.5 coding guidelines. This plan details 120+ granular, atomic tasks spanning 6 dedicated subtasks to achieve 100% compliance with zero regressions.

---

## Root Cause & Fallout Analysis

### 1. Inverted Booleans (`!isSuccess`, `!ok`, `!open`, etc.)
- **Root Cause**: Standard Go/TS idioms frequently use `if !ok` or `if (!open) return` inline without extracting dedicated positive state variables.
- **Blast Radius & Fallout**: Renaming/extracting booleans locally has zero external blast radius, provided the condition semantics are preserved verbatim (`isInvalid := !isValid` or `if isInvalid == true`).

### 2. Magic Strings & Numbers
- **Root Cause**: Format flags (`"json"`, `"csv"`), CLI modes (`"ssh"`, `"https"`), and default limits were hardcoded at call sites instead of being routed through `constants/` or owning domain enums.
- **Blast Radius & Fallout**: Extracting to constants prevents drift across CLI commands, golden tests, and JSON contracts.

### 3. Swallowed Errors in PowerShell & Shell Scripts
- **Root Cause**: PowerShell `try { ... } catch {}` blocks were used defensively to ignore expected optional cleanups without logging diagnostics to stderr.
- **Blast Radius & Fallout**: Logging `Write-Warning "[$op] $_"` ensures visibility in CI logs without failing non-critical optional cleanup flows.

### 4. Monolithic Functions Exceeding 15 Lines
- **Root Cause**: Functions accumulated input parsing, validation, execution, and error formatting into single blocks.
- **Blast Radius & Fallout**: Decomposing into single-responsibility helpers ($\le 15$ lines, $\le 8$ lines preferred) improves readability and isolated unit testing.

### 5. Nested `if` Statements
- **Root Cause**: Success branches nested inside multiple preconditions rather than utilizing early returns and guard clauses.
- **Blast Radius & Fallout**: Flattening with guard clauses simplifies cyclomatic complexity without changing behavior.

### 6. Missing `Type` Suffixes on Enums
- **Root Cause**: Legacy enum type aliases omitted the mandatory `Type` suffix.
- **Blast Radius & Fallout**: Renaming requires atomic updates to all type signatures and call sites across the package.

---

## Granular Execution Steps (120+ Steps)

### Subtask 1: Inverted Booleans (Go) - `01-inverted-booleans-go.md`
1. `gitmap/archive/extract.go:106`: Extract `!ok` to `isMissingEntry := !ok; if isMissingEntry == true`
2. `gitmap/archive/extract.go:177`: Extract `!strings.HasPrefix(...)` to `isOutsideDest := !strings.HasPrefix(...)`
3. `gitmap/archive/list.go:40`: Extract `!isExtractor` to `isNonExtractor := !isExtractor`
4. `gitmap/cliexit/cliexit_test.go:72`: Extract `!strings.Contains(out, "BUG")` to `isBugMissing := !strings.Contains(out, "BUG")`
5. `gitmap/cliexit/cliexit_test.go:75`: Extract `!strings.Contains(...)` to `isWalkMissing := !strings.Contains(...)`
6. `gitmap/cliexit/cliexit_test.go:87`: Extract `!strings.HasSuffix(...)` to `isMissingNewline := !strings.HasSuffix(...)`
7. `gitmap/cliexit/kind_test.go:81`: Extract `!strings.Contains(...)` to `isKindMissing := !strings.Contains(...)`
8. `gitmap/cliexit/report_test.go:32`: Extract `!strings.HasPrefix(...)` to `isLeadMissing := !strings.HasPrefix(...)`
9. `gitmap/cloneconcurrency/resolve_test.go:28`: Extract `!ok` to `isResolveFailed := !ok`
10. `gitmap/cloneconcurrency/resolve_test.go:44`: Extract `!ok` to `isResolveFailed := !ok`
11. `gitmap/clonefrom/depthflag_format_test.go:51`: Extract `!containsTok(...)` to `isTokMissing := !containsTok(...)`
12. `gitmap/clonefrom/depthflag_format_test.go:74`: Extract `!strings.Contains(...)` to `isFlagMissing := !strings.Contains(...)`
13. `gitmap/clonefrom/execute.go:82`: Extract `!ok` to `isExecutionFailed := !ok`
14. `gitmap/clonefrom/execute.go:123`: Extract `!confirmed` to `isDeclined := !confirmed`
15. `gitmap/clonefrom/execute_dest.go:26`: Extract `!filepath.IsAbs(absDest)` to `isRelativeDest := !filepath.IsAbs(absDest)`
16. `gitmap/clonefrom/execute_lfs_fix_test.go:27`: Extract `!ok` to `isFixFailed := !ok`
17. `gitmap/clonefrom/execute_lfs_fix_test.go:37`: Extract `!ok2` to `isSecondFixFailed := !ok2`
18. `gitmap/clonefrom/jsonschema_test.go:35`: Extract `!ok` to `isSchemaMismatch := !ok`
19. `gitmap/clonefrom/jsonschema_test.go:56`: Extract `!ok` to `isSchemaMismatch := !ok`
20. `gitmap/clonefrom/jsonschema_test.go:60`: Extract `!ok` to `isSchemaMismatch := !ok`

### Subtask 2: Inverted Booleans (TS) - `02-inverted-booleans-ts.md`
21. `src/components/docs/CloneNextCommandBuilder.tsx:107`: Extract `!s.flatten` to `const isNoFlatten = !s.flatten; if (isNoFlatten) parts.push("--no-flatten");`
22. `src/components/docs/TabOrderMap.tsx:88`: Extract `!inViewport` to `const isOutsideViewport = !inViewport; if (isOutsideViewport) ...`
23. `src/components/docs/TabOrderMap.tsx:104`: Extract `!ids` to `const isIdsMissing = !ids; if (isIdsMissing) return "";`
24. `src/components/docs/TabOrderMap.tsx:211`: Extract `!aPositive` to `const isANotPositive = !aPositive; if (isANotPositive && bPositive) return 1;`
25. `src/components/docs/TabOrderMap.tsx:258`: Extract `!open` to `const isClosed = !open; if (isClosed) return;`
26. `src/components/docs/TabOrderMap.tsx:285`: Extract `!open` to `const isClosed = !open; if (isClosed) return;`
27. `src/components/docs/TabOrderMap.tsx:288`: Extract `!target` to `const isMissingTarget = !target; if (isMissingTarget) return;`
28. `src/components/docs/TabOrderMap.tsx:299`: Extract `!active` to `const isInactive = !active; if (isInactive || active === document.body) return;`
29. `src/components/projects/ProjectDetailDialog.tsx:47`: Extract `!project` to `const isMissingProject = !project; if (isMissingProject) return null;`
30. `src/components/ui/carousel.tsx:39`: Extract `!context` to `const isMissingContext = !context; if (isMissingContext) throw ...`
31. `src/components/ui/carousel.tsx:59`: Extract `!api` to `const isMissingApi = !api; if (isMissingApi) return;`
32. `src/components/ui/carousel.tsx:89`: Extract `!api || !setApi` to `const isUninitialized = !api || !setApi; if (isUninitialized) return;`
33. `src/components/ui/chart.tsx:25`: Extract `!context` to `const isMissingContext = !context; if (isMissingContext) throw ...`
34. `src/components/ui/chart.tsx:64`: Extract `!colorConfig.length` to `const isConfigEmpty = !colorConfig.length; if (isConfigEmpty) return null;`
35. `src/components/ui/chart.tsx:146`: Extract `!value` to `const isMissingValue = !value; if (isMissingValue) return null;`
36. `src/components/ui/chart.tsx:153`: Extract `!active || !payload?.length` to `const isInactive = !active || !payload?.length; if (isInactive) return null;`
37. `src/components/ui/form.tsx:40`: Extract `!fieldContext` to `const isMissingContext = !fieldContext; if (isMissingContext) throw ...`
38. `src/components/ui/form.tsx:116`: Extract `!body` to `const isMissingBody = !body; if (isMissingBody) return null;`
39. `src/components/ui/sidebar.tsx:41`: Extract `!context` to `const isMissingContext = !context; if (isMissingContext) throw ...`
40. `src/components/ui/sidebar.tsx:480`: Extract `!tooltip` to `const isMissingTooltip = !tooltip; if (isMissingTooltip) return ...`

### Subtask 3: Swallowed Errors in PowerShell - `03-swallowed-errors-ps.md`
41. `fix-repo.ps1:139`: Replace silent catch with `catch { Write-Warning "fix-repo: failed to process branch: $_" }`
42. `fix-repo.ps1:261`: Replace silent catch with `catch { Write-Warning "fix-repo: operation failed: $_" }`
43. `gitmap/scripts/Get-LastRelease.ps1:47`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: API error: $_" }`
44. `gitmap/scripts/Get-LastRelease.ps1:69`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: JSON parse error: $_" }`
45. `gitmap/scripts/Get-LastRelease.ps1:93`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: fallback query error: $_" }`
46. `gitmap/scripts/Get-LastRelease.ps1:112`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: unexpected error: $_" }`
47. `gitmap/scripts/install.ps1:130`: Replace silent catch with `catch { Write-Warning "install: elevate failed: $_" }`
48. `gitmap/scripts/install.ps1:175`: Replace silent catch with `catch { Write-Warning "install: temp cleanup failed: $_" }`
49. `gitmap/scripts/install.ps1:196`: Replace silent catch with `catch { Write-Warning "install: asset download failed: $_" }`
50. `gitmap/scripts/install.ps1:359`: Replace silent catch with `catch { Write-Warning "install: checksum failed: $_" }`
51. `gitmap/scripts/install.ps1:371`: Replace silent catch with `catch { Write-Warning "install: signature check failed: $_" }`
52. `gitmap/scripts/install.ps1:419`: Replace silent catch with `catch { Write-Warning "install: unpack failed: $_" }`
53. `gitmap/scripts/install.ps1:532`: Replace silent catch with `catch { Write-Warning "install: link creation failed: $_" }`
54. `gitmap/scripts/install.ps1:654`: Replace silent catch with `catch { Write-Warning "install: profile registration failed: $_" }`
55. `gitmap/scripts/install.ps1:687`: Replace silent catch with `catch { Write-Warning "install: PATH update failed: $_" }`
56. `gitmap/scripts/release-version.ps1:203`: Replace silent catch with `catch { Write-Warning "release-version: tag query failed: $_" }`
57. `gitmap/scripts/release-version.ps1:300`: Replace silent catch with `catch { Write-Warning "release-version: git status failed: $_" }`
58. `gitmap/scripts/release-version.ps1:314`: Replace silent catch with `catch { Write-Warning "release-version: branch validation failed: $_" }`
59. `gitmap/scripts/release-version.ps1:355`: Replace silent catch with `catch { Write-Warning "release-version: changelog sync failed: $_" }`
60. `gitmap/scripts/release-version.ps1:389`: Replace silent catch with `catch { Write-Warning "release-version: tag creation failed: $_" }`

### Subtask 4: Monolithic Functions (>15 lines) - `04-long-functions.md`
61. `src/components/docs/TabOrderMap.tsx:34`: Refactor `findTabbableElements` into `queryTabbableNodes` + `filterVisibleNodes` ($\le 15$ lines)
62. `src/components/docs/TabOrderMap.tsx:124`: Refactor `getAccessibleLabel` into `getLabelFromAria` + `getLabelFromElement` ($\le 15$ lines)
63. `src/components/docs/TabOrderMap.tsx:200`: Refactor `compareTabOrder` into `comparePositiveTabIndex` + `compareDOMOrder` ($\le 15$ lines)
64. `src/components/docs/CodeBlock.tsx:160`: Refactor `renderHighlightedCode` into `highlightLines` + `wrapLineNumbers` ($\le 15$ lines)
65. `src/components/docs/CommandPalette.tsx:35`: Refactor `handleSelect` into `executeCommand` + `closePalette` ($\le 15$ lines)
66. `src/components/docs/TerminalDemo.tsx:75`: Refactor `TerminalDemo` render into `renderHeader` + `renderOutputBody` ($\le 15$ lines)
67. `src/components/projects/ProjectDetailDialog.tsx:60`: Refactor `ProjectDetailDialog` render into `renderProjectMeta` + `renderActionButtons` ($\le 15$ lines)
68. `src/pages/Commands.tsx:70`: Refactor `filterCommands` into `matchCategory` + `matchSearchQuery` ($\le 15$ lines)
69. `src/pages/FlagReference.tsx:50`: Refactor `sortFlags` into `compareFlagNames` + `compareCommandNames` ($\le 15$ lines)
70. `src/pages/Troubleshooting.tsx:90`: Refactor `Troubleshooting` render into `renderCategoryFilter` + `renderIssueCards` ($\le 15$ lines)
71. `gitmap/store/migrateids.go:40`: Refactor `MigrateIDs` into `prepareMigrationTx` + `applyIDRemap` ($\le 15$ lines)
72. `gitmap/formatter/formatter.go:55`: Refactor `FormatOutput` into `formatJSONStream` + `formatCSVStream` ($\le 15$ lines)
73. `gitmap/scanner/scanner.go:80`: Refactor `ScanDirectory` into `probeWorktreeMarker` + `parseGitConfig` ($\le 15$ lines)
74. `gitmap/cluster/cluster.go:65`: Refactor `ExecuteClusterCommand` into `dispatchToNode` + `aggregateResults` ($\le 15$ lines)
75. `gitmap/cmd/root.go:90`: Refactor `ExecuteRoot` into `initRootConfig` + `dispatchSubcommand` ($\le 15$ lines)

### Subtask 5: Nested If Statements - `05-nested-ifs.md`
76. `src/components/docs/CodeBlock.tsx:126`: Flatten `if (isRangeSelect) { if (next.has(lineIndex)) ... }` with early guard
77. `src/components/docs/CodeBlock.tsx:177`: Flatten `if (hljs.getLanguage(lang)) { if (res.isSuccess) ... }`
78. `src/components/docs/TabOrderMap.tsx:66`: Flatten `if (s.pointerEvents === "none" && node === el) { if (...) ... }`
79. `src/components/docs/TabOrderMap.tsx:132`: Flatten `if (el.labels && el.labels.length > 0) { if (text) ... }`
80. `src/components/docs/TabOrderMap.tsx:179`: Flatten `if (labelId) { if (ref?.textContent) ... }`
81. `src/components/docs/TabOrderMap.tsx:288`: Flatten `if (!target) { if (!active) ... }`
82. `src/components/docs/TerminalDemo.tsx:41`: Flatten `if (isPaused || isFinished) { if (isFinished) ... }`
83. `src/components/ui/carousel.tsx:59`: Flatten `if (!api) { if (event.key === "ArrowLeft") ... }`
84. `src/components/ui/chart.tsx:142`: Flatten `if (labelFormatter) { if (!value) ... }`
85. `src/components/ui/sidebar.tsx:66`: Flatten `if (setOpenProp) { if (event.key === ...) ... }`
86. `src/hooks/use-toast.ts:90`: Flatten `if (toastId) { if (action.toastId === undefined) ... }`
87. `gitmap/cluster/exec_git_test.go:47`: Flatten nested if statement with early skip return
88. `gitmap/cmd/installctx_harness_test.go:110`: Flatten nested if with guard clause
89. `gitmap/cmd/codingguidelines_test.go:20`: Flatten nested if with early return
90. `gitmap/cmd/codingguidelines_test.go:55`: Flatten nested if with early return

### Subtask 6: Enums & Magic Values - `06-enums-magic-values.md`
91. `src/components/docs/SearchBar.tsx:9`: Extract `"Search commands..."` to `DOCS_SEARCH_PLACEHOLDER` constant
92. `src/components/docs/CloneNextCommandBuilder.tsx:107`: Replace `"--no-flatten"` with `FLAG_NO_FLATTEN` constant
93. `src/components/docs/CommitTransferPage.tsx:7`: Verify `DirectionType` enum completeness
94. `src/components/docs/CodeBlock.tsx:88`: Extract `"220 10% 50%"` to `DEFAULT_ACCENT_COLOR` constant
95. `src/components/docs/CodeBlock.tsx:201`: Extract `"</span>"` to `HTML_SPAN_CLOSE` constant
96. `src/components/docs/CommandBubbles.tsx:28`: Ensure `COMMANDS_PATH` is exported from centralized routes
97. `gitmap/archive/archive.go:32`: Rename `Format` to `FormatType`
98. `gitmap/cluster/exec_git_test.go:46`: Replace `"windows"` with `constants.PlatformWindows`
99. `gitmap/cluster/exec_lifecycle.go:62`: Replace `"server"` with `constants.NodeRoleServer`
100. `gitmap/cmd/amend.go:136`: Replace `"HEAD"` with `constants.GitHead`
101. `gitmap/cmd/amendexec.go:106`: Replace `"HEAD"` with `constants.GitHead`
102. `gitmap/cmd/chromeprofile_import_csv.go:65`: Replace `"name"` with `constants.ProfileKeyName`
103. `gitmap/cmd/cluster_ops.go:122`: Replace `"csv"` with `constants.FormatCSV`
104. `gitmap/cmd/cluster_ops.go:325`: Replace `"offline"`, `"unreachable"` with `constants.NodeStatusOffline`, `constants.NodeStatusUnreachable`
105. `gitmap/cmd/doctorchecks.go:213`: Replace `"Valid"` with `constants.DoctorStatusValid`
106. `gitmap/cmd/doctorchecks.go:220`: Replace `"NotSigned"` with `constants.DoctorStatusNotSigned`
107. `gitmap/cmd/doctordupbin.go:31`: Replace `"windows"` with `constants.PlatformWindows`
108. `gitmap/cmd/doctordupbin.go:105`: Replace `"windows"` with `constants.PlatformWindows`
109. `gitmap/cmd/envplatform_unix.go:62`: Replace `"zsh"` with `constants.ShellZsh`
110. `gitmap/cmd/installdetect.go:21`: Replace `"windows"` with `constants.PlatformWindows`
111. `gitmap/cmd/installdetect.go:24`: Replace `"darwin"` with `constants.PlatformDarwin`
112. `gitmap/cmd/installscripts.go:94`: Replace `"windows"` with `constants.PlatformWindows`
113. `gitmap/cmd/list.go:57`: Replace `"groups"` with `constants.ListModeGroups`
114. `gitmap/cmd/list.go:68`: Replace `"groups"` with `constants.ListModeGroups`
115. `gitmap/cmd/llmdocs.go:55`: Replace `"json"` with `constants.FormatJSON`
116. `gitmap/cmd/llmdocs.go:111`: Replace `"json"` with `constants.FormatJSON`
117. `gitmap/cmd/move.go:18`: Replace magic `2` with `constants.MoveRequiredArgsCount`
118. `gitmap/cmd/prune.go:59`: Replace `"y"`, `"Y"` with `constants.ConfirmYesLower`, `constants.ConfirmYesUpper`
119. `gitmap/cmd/replace.go:67`: Replace magic `2` with `constants.ReplaceRequiredArgsCount`
120. `gitmap/cmd/selfinstall.go:248`: Replace `"windows"` with `constants.PlatformWindows`

---

## Verification Plan
1. `go build ./...` across all Go packages.
2. `go test ./... -short` across the Go test suite.
3. `npm run build` or Vite build check for TypeScript/React.
4. Verify all subtasks maintain $\le 15$ lines per function and zero nested `if` statements.
