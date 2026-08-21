# Plan: Coding Guideline Audit & Enforcement (v4)

Slug: 01-coding-guideline-fixes
Status: pending
Created: 2026-08-22

## Context
Massive 50+ step programmatic deep read of the codebase to uncover remaining v1.4.5 coding guideline violations (long functions, nested ifs, missing enum suffixes, magic values). The previous pass missed these specific checks and hallucinated the agent completions.

## Enqueued Granular Tasks
1. **enum-suffix**: `.\gitmap\archive\archive.go:32` - Type alias Format acting as enum lacks Type suffix. **Fix**: Rename to FormatType.
2. **long-func**: `.\gitmap\archive\archive.go:92` - Function Extension exceeds 15 lines (28 lines). **Fix**: Extract helper functions or early returns.
3. **enum-suffix**: `.\gitmap\archive\create.go:33` - Type alias CompressionMode acting as enum lacks Type suffix. **Fix**: Rename to CompressionModeType.
4. **long-func**: `.\gitmap\archive\create.go:166` - Function buildArchiver exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
5. **enum-suffix**: `.\gitmap\archive\source.go:30` - Type alias SourceKind acting as enum lacks Type suffix. **Fix**: Rename to SourceKindType.
6. **long-func**: `.\gitmap\cliexit\cliexit_test.go:17` - Function TestFormatLine_Shape exceeds 15 lines (26 lines). **Fix**: Extract helper functions or early returns.
7. **enum-suffix**: `.\gitmap\cliexit\kind.go:37` - Type alias Kind acting as enum lacks Type suffix. **Fix**: Rename to KindType.
8. **enum-suffix**: `.\gitmap\cliexit\report.go:50` - Type alias OutputMode acting as enum lacks Type suffix. **Fix**: Rename to OutputModeType.
9. **long-func**: `.\gitmap\cliexit\report_test.go:17` - Function TestWriteStructured_HumanMode exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
10. **magic-value**: `.\gitmap\clonefrom\execute_lfs_fix.go:26` - Magic string/number in condition: if match := rxSmudgeFatal.FindStringSubmatch(output); len(match) == 2 {. **Fix**: Extract to named constant or EnumType.
11. **magic-value**: `.\gitmap\clonefrom\execute_lfs_fix.go:30` - Magic string/number in condition: if match := rxSmudgeError.FindStringSubmatch(output); len(match) == 2 {. **Fix**: Extract to named constant or EnumType.
12. **long-func**: `.\gitmap\clonefrom\jsonschema.go:52` - Function EmitReportSchema exceeds 15 lines (45 lines). **Fix**: Extract helper functions or early returns.
13. **magic-value**: `.\gitmap\clonefrom\parse.go:70` - Magic string/number in condition: if format == "json" {. **Fix**: Extract to named constant or EnumType.
14. **long-func**: `.\gitmap\clonefrom\parsecsv_columnerr_test.go:20` - Function TestParseCSV_RowError_NamesOffendingColumn exceeds 15 lines (51 lines). **Fix**: Extract helper functions or early returns.
15. **long-func**: `.\gitmap\clonefrom\summary_csvquoting_golden_test.go:57` - Function quotingEdgeCaseResults exceeds 15 lines (31 lines). **Fix**: Extract helper functions or early returns.
16. **long-func**: `.\gitmap\clonefrom\summary_golden_test.go:47` - Function canonicalReportResults exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
17. **magic-value**: `.\gitmap\clonenext\localstate.go:105` - Magic string/number in condition: if len(parts) == 2 {. **Fix**: Extract to named constant or EnumType.
18. **magic-value**: `.\gitmap\clonenow\clonenow.go:97` - Magic string/number in condition: if mode == "ssh" {. **Fix**: Extract to named constant or EnumType.
19. **long-func**: `.\gitmap\clonenow\crossformat_golden_test.go:54` - Function crossFormatScanRecords exceeds 15 lines (30 lines). **Fix**: Extract helper functions or early returns.
20. **magic-value**: `.\gitmap\clonenow\parse_schema.go:108` - Magic string/number in condition: if name == "httpsUrl" || name == "sshUrl" {. **Fix**: Extract to named constant or EnumType.
21. **magic-value**: `.\gitmap\clonenow\parse_schema.go:148` - Magic string/number in condition: if name == "httpsUrl" || name == "sshUrl" {. **Fix**: Extract to named constant or EnumType.
22. **enum-suffix**: `.\gitmap\cloner\audit.go:30` - Type alias AuditAction acting as enum lacks Type suffix. **Fix**: Rename to AuditActionType.
23. **long-func**: `.\gitmap\cloner\batchprogress_test.go:117` - Function TestMixedOperations_Counters exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
24. **magic-value**: `.\gitmap\cloner\runners.go:156` - Magic string/number in condition: if strings.ToLower(strings.TrimSpace(response)) == "y" {. **Fix**: Extract to named constant or EnumType.
25. **long-func**: `.\gitmap\cloner\strategy_test.go:10` - Function TestPickCloneStrategy exceeds 15 lines (50 lines). **Fix**: Extract helper functions or early returns.
26. **magic-value**: `.\gitmap\cluster\exec_git_test.go:46` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
27. **nested-if**: `.\gitmap\cluster\exec_git_test.go:47` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
28. **nested-if**: `.\gitmap\cluster\exec_install_test.go:29` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
29. **magic-value**: `.\gitmap\cluster\exec_lifecycle.go:62` - Magic string/number in condition: if node.NodeRole == "server" || node.IsServer {. **Fix**: Extract to named constant or EnumType.
30. **nested-if**: `.\gitmap\cluster\exec_lifecycle.go:69` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
31. **magic-value**: `.\gitmap\cluster\exec_ps_test.go:25` - Magic string/number in condition: if file == "pwsh" {. **Fix**: Extract to named constant or EnumType.
32. **magic-value**: `.\gitmap\cluster\exec_ps_test.go:67` - Magic string/number in condition: if file == "pwsh" {. **Fix**: Extract to named constant or EnumType.
33. **magic-value**: `.\gitmap\cluster\exec_ps_test.go:70` - Magic string/number in condition: if file == "powershell" {. **Fix**: Extract to named constant or EnumType.
34. **magic-value**: `.\gitmap\cluster\exec_ps_test.go:90` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
35. **magic-value**: `.\gitmap\cluster\exec_ps_test.go:102` - Magic string/number in condition: if file == "pwsh" {. **Fix**: Extract to named constant or EnumType.
36. **magic-value**: `.\gitmap\cluster\exec_ps_test.go:122` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
37. **enum-suffix**: `.\gitmap\cluster\registry.go:9` - Type alias NodeState acting as enum lacks Type suffix. **Fix**: Rename to NodeStateType.
38. **nested-if**: `.\gitmap\cluster\runref_test.go:34` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
39. **magic-value**: `.\gitmap\cmd\amend.go:136` - Magic string/number in condition: if f.commitHash == "HEAD" {. **Fix**: Extract to named constant or EnumType.
40. **long-func**: `.\gitmap\cmd\amendauditjson_contract_test.go:21` - Function canonicalAmendmentRecord exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
41. **magic-value**: `.\gitmap\cmd\amendexec.go:106` - Magic string/number in condition: if f.commitHash == "HEAD" {. **Fix**: Extract to named constant or EnumType.
42. **long-func**: `.\gitmap\cmd\amendlistrender.go:45` - Function buildAmendListJSONItems exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
43. **magic-value**: `.\gitmap\cmd\chromeprofile_import_csv.go:65` - Magic string/number in condition: if key == "name" {. **Fix**: Extract to named constant or EnumType.
44. **long-func**: `.\gitmap\cmd\cliexit_clone_test.go:138` - Function writeCloneNowManifest exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
45. **magic-value**: `.\gitmap\cmd\cliexit_helpers_test.go:89` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
46. **long-func**: `.\gitmap\cmd\cliexit_scan_test.go:21` - Function TestScanCLI_ExitCodes exceeds 15 lines (25 lines). **Fix**: Extract helper functions or early returns.
47. **magic-value**: `.\gitmap\cmd\cliexit_scan_test.go:50` - Magic string/number in condition: if tc.name == "failure_missing_dir" {. **Fix**: Extract to named constant or EnumType.
48. **long-func**: `.\gitmap\cmd\clonefixrepo.go:254` - Function parseCloneFixRepoArgs exceeds 15 lines (39 lines). **Fix**: Extract helper functions or early returns.
49. **long-func**: `.\gitmap\cmd\clonefixrepoparallel.go:111` - Function runOneCFRJob exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
50. **long-func**: `.\gitmap\cmd\clonefrom_flags.go:61` - Function bindCloneFromFlags exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
51. **long-func**: `.\gitmap\cmd\clonenextflags.go:74` - Function parseCloneNextFlags exceeds 15 lines (33 lines). **Fix**: Extract helper functions or early returns.
52. **long-func**: `.\gitmap\cmd\clonenow.go:115` - Function parseCloneNowFlags exceeds 15 lines (45 lines). **Fix**: Extract helper functions or early returns.
53. **long-func**: `.\gitmap\cmd\clonepick_flags.go:41` - Function parseClonePickFlags exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
54. **long-func**: `.\gitmap\cmd\clonepick_flags.go:69` - Function bindClonePickCoreFlags exceeds 15 lines (25 lines). **Fix**: Extract helper functions or early returns.
55. **long-func**: `.\gitmap\cmd\clonetermblock_golden_test.go:59` - Function cloneTermGoldenCases exceeds 15 lines (163 lines). **Fix**: Extract helper functions or early returns.
56. **long-func**: `.\gitmap\cmd\cloneurlconvert_test.go:5` - Function TestConvertURLToSSH exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
57. **long-func**: `.\gitmap\cmd\clone_idempotent_test.go:23` - Function TestRunClone_Idempotent exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
58. **long-func**: `.\gitmap\cmd\clusterflags.go:93` - Function ParseClusterFlags exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
59. **long-func**: `.\gitmap\cmd\clusterflags_test.go:8` - Function TestParseClusterFlags exceeds 15 lines (72 lines). **Fix**: Extract helper functions or early returns.
60. **long-func**: `.\gitmap\cmd\clustersubcmd_test.go:11` - Function TestClusterSubCommandParser exceeds 15 lines (31 lines). **Fix**: Extract helper functions or early returns.
61. **magic-value**: `.\gitmap\cmd\cluster_ops.go:122` - Magic string/number in condition: if format == "csv" {. **Fix**: Extract to named constant or EnumType.
62. **magic-value**: `.\gitmap\cmd\cluster_ops.go:325` - Magic string/number in condition: if strings.ToLower(n.Status) == "offline" || strings.ToLower(n.Status) == "unreachable" {. **Fix**: Extract to named constant or EnumType.
63. **long-func**: `.\gitmap\cmd\code_test.go:8` - Function TestMergeStringPaths exceeds 15 lines (37 lines). **Fix**: Extract helper functions or early returns.
64. **magic-value**: `.\gitmap\cmd\codingguidelines.go:43` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
65. **magic-value**: `.\gitmap\cmd\codingguidelines_test.go:19` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
66. **magic-value**: `.\gitmap\cmd\codingguidelines_test.go:54` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
67. **magic-value**: `.\gitmap\cmd\codingguidelines_test.go:74` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
68. **long-func**: `.\gitmap\cmd\committransfer.go:162` - Function registerCommitTransferStrings exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
69. **long-func**: `.\gitmap\cmd\completion.go:43` - Function handleCompletionList exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
70. **long-func**: `.\gitmap\cmd\diffprofilesjson_contract_test.go:20` - Function canonicalDPResult exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
71. **magic-value**: `.\gitmap\cmd\doctorchecks.go:213` - Magic string/number in condition: if status == "Valid" {. **Fix**: Extract to named constant or EnumType.
72. **magic-value**: `.\gitmap\cmd\doctorchecks.go:220` - Magic string/number in condition: if status == "NotSigned" {. **Fix**: Extract to named constant or EnumType.
73. **magic-value**: `.\gitmap\cmd\doctordupbin.go:31` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
74. **magic-value**: `.\gitmap\cmd\doctordupbin.go:105` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
75. **long-func**: `.\gitmap\cmd\downloaderconfig.go:84` - Function promptDownloaderConfig exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
76. **long-func**: `.\gitmap\cmd\envops.go:13` - Function runEnvSet exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
77. **magic-value**: `.\gitmap\cmd\envplatform_unix.go:62` - Magic string/number in condition: if shell == "zsh" {. **Fix**: Extract to named constant or EnumType.
78. **magic-value**: `.\gitmap\cmd\env_unit_test.go:59` - Magic string/number in condition: if v.Name == "B" {. **Fix**: Extract to named constant or EnumType.
79. **magic-value**: `.\gitmap\cmd\escapecwd_test.go:76` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
80. **magic-value**: `.\gitmap\cmd\expandhome_test.go:56` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
81. **long-func**: `.\gitmap\cmd\findnextjson_contract_test.go:48` - Function canonicalFindNextRow exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
82. **long-func**: `.\gitmap\cmd\fixrepo_rewrite_scan_test.go:21` - Function TestScanUnguardedTokenHits exceeds 15 lines (56 lines). **Fix**: Extract helper functions or early returns.
83. **long-func**: `.\gitmap\cmd\fixrepo_rewrite_v9tov12_test.go:64` - Function TestFixRepoRewriteV9ToV12Fixture exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
84. **long-func**: `.\gitmap\cmd\fixrepo_rewrite_v9tov12_test.go:208` - Function renderFixRepoFailureDiff exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
85. **long-func**: `.\gitmap\cmd\history.go:125` - Function printHistoryRow exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
86. **long-func**: `.\gitmap\cmd\historyrender.go:46` - Function buildHistoryJSONItems exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
87. **long-func**: `.\gitmap\cmd\historyrewrite_paths_test.go:12` - Function TestParseHistoryPathsNormalizesAllInputForms exceeds 15 lines (46 lines). **Fix**: Extract helper functions or early returns.
88. **long-func**: `.\gitmap\cmd\historyrewrite_pin.go:99` - Function buildPinCallbackPython exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
89. **long-func**: `.\gitmap\cmd\install.go:12` - Function runInstall exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
90. **magic-value**: `.\gitmap\cmd\installcleancode.go:70` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
91. **long-func**: `.\gitmap\cmd\installctxflatten.go:77` - Function slugifyCtx exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
92. **long-func**: `.\gitmap\cmd\installctxmac.go:164` - Function macDocumentWflow exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
93. **magic-value**: `.\gitmap\cmd\installctx_harness_test.go:109` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
94. **magic-value**: `.\gitmap\cmd\installdetect.go:21` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
95. **magic-value**: `.\gitmap\cmd\installdetect.go:24` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
96. **magic-value**: `.\gitmap\cmd\installscripts.go:94` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
97. **long-func**: `.\gitmap\cmd\installtools.go:293` - Function resolveChocoPackage exceeds 15 lines (30 lines). **Fix**: Extract helper functions or early returns.
98. **long-func**: `.\gitmap\cmd\installtools.go:353` - Function resolveAptPackage exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
99. **long-func**: `.\gitmap\cmd\installtools.go:382` - Function resolveBrewPackage exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
100. **long-func**: `.\gitmap\cmd\installverify.go:148` - Function toolBinaryName exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
101. **long-func**: `.\gitmap\cmd\install_unit_test.go:43` - Function TestBuildUninstallCommand exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
102. **long-func**: `.\gitmap\cmd\jsonsnapshot_helpers_failuremsg_test.go:29` - Function TestExpectDelim_FailureMessages exceeds 15 lines (25 lines). **Fix**: Extract helper functions or early returns.
103. **long-func**: `.\gitmap\cmd\jsonsnapshot_helpers_failuremsg_test.go:68` - Function TestScanEveryObjectKeysPure_PointsAtBrokenObject exceeds 15 lines (26 lines). **Fix**: Extract helper functions or early returns.
104. **long-func**: `.\gitmap\cmd\latestbranch.go:138` - Function parseLatestBranchFlags exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
105. **magic-value**: `.\gitmap\cmd\list.go:57` - Magic string/number in condition: if lower == "groups" {. **Fix**: Extract to named constant or EnumType.
106. **magic-value**: `.\gitmap\cmd\list.go:68` - Magic string/number in condition: if lower == "groups" {. **Fix**: Extract to named constant or EnumType.
107. **long-func**: `.\gitmap\cmd\listreleasesrender.go:87` - Function buildListReleasesJSONItems exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
108. **long-func**: `.\gitmap\cmd\listreleasesrender.go:122` - Function buildListReleasesAllReposJSONItems exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
109. **magic-value**: `.\gitmap\cmd\llmdocs.go:55` - Magic string/number in condition: if *format == "json" {. **Fix**: Extract to named constant or EnumType.
110. **magic-value**: `.\gitmap\cmd\llmdocs.go:111` - Magic string/number in condition: if format == "json" {. **Fix**: Extract to named constant or EnumType.
111. **long-func**: `.\gitmap\cmd\llmdocscommands.go:55` - Function buildCommandGroups exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
112. **long-func**: `.\gitmap\cmd\llmdocsgroups.go:192` - Function buildUtilityGroup exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
113. **long-func**: `.\gitmap\cmd\llmdocsheader.go:23` - Function writeLLMArchitecture exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
114. **magic-value**: `.\gitmap\cmd\llmdocsjson_contract_test.go:100` - Magic string/number in condition: if len(got) == 4 && got[3] != "example" {. **Fix**: Extract to named constant or EnumType.
115. **long-func**: `.\gitmap\cmd\llmdocssections.go:40` - Function writeLLMProjectStructure exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
116. **magic-value**: `.\gitmap\cmd\move.go:18` - Magic string/number in condition: if len(positional) == 2 {. **Fix**: Extract to named constant or EnumType.
117. **nested-if**: `.\gitmap\cmd\move.go:19` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
118. **long-func**: `.\gitmap\cmd\movemergeflags.go:19` - Function bindFlags exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
119. **long-func**: `.\gitmap\cmd\pendingclear.go:78` - Function parsePendingClearArgs exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
120. **long-func**: `.\gitmap\cmd\projectreposrender.go:39` - Function buildProjectReposJSONItems exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
121. **magic-value**: `.\gitmap\cmd\prune.go:59` - Magic string/number in condition: if answer == "y" || answer == "Y" {. **Fix**: Extract to named constant or EnumType.
122. **long-func**: `.\gitmap\cmd\release.go:141` - Function parseReleaseFlags exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
123. **long-func**: `.\gitmap\cmd\releaseargs.go:12` - Function reorderFlagsBeforeArgs exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
124. **long-func**: `.\gitmap\cmd\release_notes_opts.go:37` - Function parseReleaseNotesArgs exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
125. **long-func**: `.\gitmap\cmd\release_notes_opts.go:96` - Function classifyCommit exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
126. **magic-value**: `.\gitmap\cmd\replace.go:67` - Magic string/number in condition: if len(positional) == 2 {. **Fix**: Extract to named constant or EnumType.
127. **long-func**: `.\gitmap\cmd\rescansubtree_test.go:12` - Function TestSplitRescanSubtreeArgs exceeds 15 lines (43 lines). **Fix**: Extract helper functions or early returns.
128. **long-func**: `.\gitmap\cmd\rescansubtree_test.go:118` - Function TestExtractMaxDepthForLog exceeds 15 lines (27 lines). **Fix**: Extract helper functions or early returns.
129. **nested-if**: `.\gitmap\cmd\rm.go:129` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
130. **long-func**: `.\gitmap\cmd\rootcore.go:14` - Function coreDispatchEntries exceeds 15 lines (109 lines). **Fix**: Extract helper functions or early returns.
131. **long-func**: `.\gitmap\cmd\rootdata.go:13` - Function dataDispatchEntries exceeds 15 lines (27 lines). **Fix**: Extract helper functions or early returns.
132. **long-func**: `.\gitmap\cmd\rootflags.go:37` - Function parseScanFlags exceeds 15 lines (36 lines). **Fix**: Extract helper functions or early returns.
133. **long-func**: `.\gitmap\cmd\rootflags.go:230` - Function parseCloneFlags exceeds 15 lines (62 lines). **Fix**: Extract helper functions or early returns.
134. **long-func**: `.\gitmap\cmd\rootrelease.go:13` - Function releaseDispatchEntries exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
135. **long-func**: `.\gitmap\cmd\roottooling.go:13` - Function toolingDispatchEntries exceeds 15 lines (121 lines). **Fix**: Extract helper functions or early returns.
136. **long-func**: `.\gitmap\cmd\rootusage.go:10` - Function printUsage exceeds 15 lines (40 lines). **Fix**: Extract helper functions or early returns.
137. **long-func**: `.\gitmap\cmd\rootusagecompact.go:18` - Function compactGroups exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
138. **long-func**: `.\gitmap\cmd\rootusagefilter_rows.go:16` - Function allHelpRows exceeds 15 lines (62 lines). **Fix**: Extract helper functions or early returns.
139. **long-func**: `.\gitmap\cmd\rootusageflags.go:10` - Function printGroupUtilities exceeds 15 lines (23 lines). **Fix**: Extract helper functions or early returns.
140. **long-func**: `.\gitmap\cmd\rootusageflags.go:59` - Function printUsageFixRepoFlags exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
141. **long-func**: `.\gitmap\cmd\rootusageflags.go:137` - Function printUsageSEOFlags exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
142. **long-func**: `.\gitmap\cmd\rootutility.go:23` - Function utilityDispatchEntries exceeds 15 lines (31 lines). **Fix**: Extract helper functions or early returns.
143. **long-func**: `.\gitmap\cmd\root_url_shortcut_test.go:12` - Function TestShouldRewriteToCloneCoversReportedInvocations exceeds 15 lines (56 lines). **Fix**: Extract helper functions or early returns.
144. **magic-value**: `.\gitmap\cmd\scan_export_clonefrom_integration_test.go:110` - Magic string/number in condition: if format == "json" {. **Fix**: Extract to named constant or EnumType.
145. **long-func**: `.\gitmap\cmd\selfinstall.go:122` - Function parseSelfInstallFlags exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
146. **magic-value**: `.\gitmap\cmd\selfinstall.go:248` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
147. **magic-value**: `.\gitmap\cmd\selfinstall.go:300` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
148. **magic-value**: `.\gitmap\cmd\selfinstall.go:309` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
149. **long-func**: `.\gitmap\cmd\selfuninstall.go:51` - Function parseSelfUninstallFlags exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
150. **magic-value**: `.\gitmap\cmd\selfuninstallparts.go:65` - Magic string/number in condition: if lower == "gitmap" || lower == "gitmap.exe" {. **Fix**: Extract to named constant or EnumType.
151. **long-func**: `.\gitmap\cmd\seowrite.go:106` - Function parseSEOWriteFlags exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
152. **long-func**: `.\gitmap\cmd\seowritecreate.go:37` - Function buildSampleTemplate exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
153. **long-func**: `.\gitmap\cmd\sshgen.go:74` - Function parseSSHGenFlags exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
154. **magic-value**: `.\gitmap\cmd\sshgen.go:130` - Magic string/number in condition: if input == "R" {. **Fix**: Extract to named constant or EnumType.
155. **magic-value**: `.\gitmap\cmd\sshgen.go:136` - Magic string/number in condition: if input == "N" {. **Fix**: Extract to named constant or EnumType.
156. **magic-value**: `.\gitmap\cmd\sshgen.go:218` - Magic string/number in condition: if action == "update" {. **Fix**: Extract to named constant or EnumType.
157. **magic-value**: `.\gitmap\cmd\sshgen.go:221` - Magic string/number in condition: if action == "save" {. **Fix**: Extract to named constant or EnumType.
158. **long-func**: `.\gitmap\cmd\startup.go:72` - Function parseStartupListFlags exceeds 15 lines (23 lines). **Fix**: Extract helper functions or early returns.
159. **long-func**: `.\gitmap\cmd\startupadd.go:87` - Function parseStartupAddFlags exceeds 15 lines (28 lines). **Fix**: Extract helper functions or early returns.
160. **long-func**: `.\gitmap\cmd\startuplistfilter_test.go:18` - Function fixtureStartupEntries exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
161. **long-func**: `.\gitmap\cmd\startuplistjson_determinism_test.go:45` - Function TestStartupListJSON_DeterministicAcrossRuns exceeds 15 lines (23 lines). **Fix**: Extract helper functions or early returns.
162. **long-func**: `.\gitmap\cmd\statsjson_contract_test.go:28` - Function canonicalStatsOverall exceeds 15 lines (25 lines). **Fix**: Extract helper functions or early returns.
163. **long-func**: `.\gitmap\cmd\statusprint.go:111` - Function printStatusTableWithContext exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
164. **magic-value**: `.\gitmap\cmd\templatesdiff.go:140` - Magic string/number in condition: if kind == "attributes" {. **Fix**: Extract to named constant or EnumType.
165. **long-func**: `.\gitmap\cmd\tempreleaselistjson_contract_test.go:24` - Function canonicalTempReleaseList exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
166. **magic-value**: `.\gitmap\cmd\update.go:138` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
167. **long-func**: `.\gitmap\cmd\updatecleanup_handoff_test.go:44` - Function TestBuildCleanupChildEnvForwardsDelayAndJSONPath exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
168. **long-func**: `.\gitmap\cmd\updatehandoff_phase3.go:130` - Function spawnDeployedCleanupWindows exceeds 15 lines (23 lines). **Fix**: Extract helper functions or early returns.
169. **long-func**: `.\gitmap\cmd\updatehandoff_phase3.go:168` - Function spawnDeployedCleanupUnix exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
170. **magic-value**: `.\gitmap\cmd\updateremoteinstall.go:108` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
171. **magic-value**: `.\gitmap\cmd\updateremoteinstall.go:129` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
172. **magic-value**: `.\gitmap\cmd\updateremoteinstall.go:136` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
173. **magic-value**: `.\gitmap\cmd\updateremoteinstall.go:155` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
174. **magic-value**: `.\gitmap\cmd\updateremoteinstall_test.go:26` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
175. **magic-value**: `.\gitmap\cmd\updatescript.go:17` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
176. **long-func**: `.\gitmap\cmd\updatescript.go:139` - Function buildUpdateScript exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
177. **long-func**: `.\gitmap\cmd\versionhistoryjson_contract_test.go:28` - Function canonicalVersionHistoryRecords exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
178. **magic-value**: `.\gitmap\cmd\visibilitybulkprompt.go:98` - Magic string/number in condition: if tok == "y" || tok == "yes" {. **Fix**: Extract to named constant or EnumType.
179. **magic-value**: `.\gitmap\cmd\vscodepmsync_dedupe_test.go:38` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
180. **magic-value**: `.\gitmap\cmd\vscodepmsync_dedupe_test.go:70` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
181. **magic-value**: `.\gitmap\cmd\vscodepmsync_dedupe_test.go:106` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
182. **magic-value**: `.\gitmap\cmd\vscodepmsync_dedupe_test.go:136` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
183. **long-func**: `.\gitmap\cmd\vscodepmsync_flags.go:89` - Function parseVSCodePMSyncFlags exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
184. **magic-value**: `.\gitmap\cmd\vscodepmsync_mode_test.go:34` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
185. **magic-value**: `.\gitmap\cmd\vscodepmsync_mode_test.go:60` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
186. **magic-value**: `.\gitmap\cmd\vscodepmsync_mode_test.go:91` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
187. **magic-value**: `.\gitmap\cmd\vscodepmsync_mode_test.go:123` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
188. **magic-value**: `.\gitmap\cmd\vscodepmsync_pathtag_test.go:66` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
189. **magic-value**: `.\gitmap\cmd\vscodepmsync_pathtag_test.go:92` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
190. **magic-value**: `.\gitmap\cmd\vscodepmsync_pathtag_test.go:126` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
191. **magic-value**: `.\gitmap\cmd\vscodepmsync_test.go:83` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
192. **magic-value**: `.\gitmap\cmd\vscodepmsync_test.go:111` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
193. **long-func**: `.\gitmap\cmd\vscodesyncdisabled_test.go:15` - Function TestStripVSCodeSyncDisabledFlag exceeds 15 lines (25 lines). **Fix**: Extract helper functions or early returns.
194. **magic-value**: `.\gitmap\cmd\watchformat.go:56` - Magic string/number in condition: if snap.Status == "error" {. **Fix**: Extract to named constant or EnumType.
195. **magic-value**: `.\gitmap\cmd\watchformat.go:77` - Magic string/number in condition: if status == "dirty" {. **Fix**: Extract to named constant or EnumType.
196. **magic-value**: `.\gitmap\cmd\watchops.go:95` - Magic string/number in condition: if snap.Status == "dirty" {. **Fix**: Extract to named constant or EnumType.
197. **magic-value**: `.\gitmap\cmd\whoami.go:58` - Magic string/number in condition: if strings.HasSuffix(name, ".pub") || name == "known_hosts" || name == "config" || strings.HasPrefix(name, "known_hosts") {. **Fix**: Extract to named constant or EnumType.
198. **magic-value**: `.\gitmap\cmd\whoami.go:213` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
199. **long-func**: `.\gitmap\cmd\commitin\enums.go:206` - Function String exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
200. **long-func**: `.\gitmap\cmd\commitin\enums_test.go:20` - Function TestCommitInEnumsMatchSpec exceeds 15 lines (80 lines). **Fix**: Extract helper functions or early returns.
201. **long-func**: `.\gitmap\cmd\commitin\enums_test.go:144` - Function TestCommitInFlagNamesShape exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
202. **enum-suffix**: `.\gitmap\cmd\commitin\finalize\conflict.go:12` - Type alias ConflictDecision acting as enum lacks Type suffix. **Fix**: Rename to ConflictDecisionType.
203. **magic-value**: `.\gitmap\cmd\commitin\message\message_test.go:109` - Magic string/number in condition: if out == "Refined" {. **Fix**: Extract to named constant or EnumType.
204. **long-func**: `.\gitmap\cmd\commitin\profile\build.go:17` - Function BuildFromResolved exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
205. **long-func**: `.\gitmap\cmd\commitin\replay\replay_test.go:36` - Function TestApplyCommitWiresFullPipelineThroughHooks exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
206. **long-func**: `.\gitmap\cmd\commitin\runlog\tagreplay_test.go:75` - Function TestRecordTagReplayCreatedWritesAllColumns exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
207. **long-func**: `.\gitmap\cmd\commitin\runlog\tagreplay_test.go:233` - Function TestIsAnnotatedSemverVersionTagMatrix exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
208. **magic-value**: `.\gitmap\cmd\commitin\workspace\workspace_test.go:189` - Magic string/number in condition: if sub == "clone" && len(args) == 2 {. **Fix**: Extract to named constant or EnumType.
209. **long-func**: `.\gitmap\committransfer\interleave_test.go:11` - Function TestBuildInterleavedStreamSortsByAuthorDate exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
210. **magic-value**: `.\gitmap\committransfer\replay.go:144` - Magic string/number in condition: if first == "node_modules" && !opts.IncludeNodeMod {. **Fix**: Extract to named constant or EnumType.
211. **enum-suffix**: `.\gitmap\committransfer\types.go:17` - Type alias Direction acting as enum lacks Type suffix. **Fix**: Rename to DirectionType.
212. **enum-suffix**: `.\gitmap\committransfer\types.go:63` - Type alias PreferPolicy acting as enum lacks Type suffix. **Fix**: Rename to PreferPolicyType.
213. **magic-value**: `.\gitmap\completion\bash.go:56` - Magic string/number in condition: if [[ ${COMP_CWORD} -ge 3 ]] && [[ "$sub" == "add" || "$sub" == "show" || "$sub" == "delete" || "$sub" == "remove" || "$sub" == "rename" ]]; then. **Fix**: Extract to named constant or EnumType.
214. **long-func**: `.\gitmap\completion\bash.go:4` - Function generateBash exceeds 15 lines (113 lines). **Fix**: Extract helper functions or early returns.
215. **magic-value**: `.\gitmap\completion\detect.go:11` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
216. **magic-value**: `.\gitmap\completion\dynamic.go:120` - Magic string/number in condition: if name == "Default" || strings.HasPrefix(name, "Profile ") {. **Fix**: Extract to named constant or EnumType.
217. **magic-value**: `.\gitmap\completion\install.go:124` - Magic string/number in condition: if goos == "windows" {. **Fix**: Extract to named constant or EnumType.
218. **nested-if**: `.\gitmap\completion\powershell.go:13` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
219. **nested-if**: `.\gitmap\completion\powershell.go:61` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
220. **nested-if**: `.\gitmap\completion\powershell.go:89` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
221. **nested-if**: `.\gitmap\completion\powershell.go:110` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
222. **nested-if**: `.\gitmap\completion\powershell.go:111` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
223. **nested-if**: `.\gitmap\completion\powershell.go:124` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
224. **nested-if**: `.\gitmap\completion\powershell.go:164` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
225. **nested-if**: `.\gitmap\completion\powershell.go:177` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
226. **magic-value**: `.\gitmap\completion\zsh.go:10` - Magic string/number in condition: if (( CURRENT == 2 )); then. **Fix**: Extract to named constant or EnumType.
227. **magic-value**: `.\gitmap\completion\zsh.go:73` - Magic string/number in condition: if (( CURRENT >= 4 )) && [[ "${words[3]}" == "add" || "${words[3]}" == "show" || "${words[3]}" == "delete" || "${words[3]}" == "remove" || "${words[3]}" == "rename" ]]; then. **Fix**: Extract to named constant or EnumType.
228. **long-func**: `.\gitmap\completion\zsh.go:4` - Function generateZsh exceeds 15 lines (150 lines). **Fix**: Extract helper functions or early returns.
229. **magic-value**: `.\gitmap\config\config_test.go:14` - Magic string/number in condition: if cfg.DefaultMode == "https" {. **Fix**: Extract to named constant or EnumType.
230. **magic-value**: `.\gitmap\config\config_test.go:17` - Magic string/number in condition: if cfg.DefaultOutput == "terminal" {. **Fix**: Extract to named constant or EnumType.
231. **magic-value**: `.\gitmap\config\config_test.go:31` - Magic string/number in condition: if cfg.DefaultMode == "https" {. **Fix**: Extract to named constant or EnumType.
232. **magic-value**: `.\gitmap\config\config_test.go:51` - Magic string/number in condition: if cfg.DefaultMode == "ssh" {. **Fix**: Extract to named constant or EnumType.
233. **magic-value**: `.\gitmap\config\config_test.go:64` - Magic string/number in condition: if merged.DefaultMode == "ssh" {. **Fix**: Extract to named constant or EnumType.
234. **magic-value**: `.\gitmap\config\config_test.go:67` - Magic string/number in condition: if merged.DefaultOutput == "json" {. **Fix**: Extract to named constant or EnumType.
235. **magic-value**: `.\gitmap\config\config_test.go:84` - Magic string/number in condition: if merged.DefaultMode == "ssh" {. **Fix**: Extract to named constant or EnumType.
236. **long-func**: `.\gitmap\constants\cmd_constants_test.go:13` - Function topLevelCmds exceeds 15 lines (332 lines). **Fix**: Extract helper functions or early returns.
237. **long-func**: `.\gitmap\constants\constants_clonefrom_csvaliases_test.go:9` - Function TestCanonicalCSVColumn exceeds 15 lines (26 lines). **Fix**: Extract helper functions or early returns.
238. **enum-suffix**: `.\gitmap\constants\constants_inject_idempotency.go:44` - Type alias InjectKind acting as enum lacks Type suffix. **Fix**: Rename to InjectKindType.
239. **long-func**: `.\gitmap\db\clusterexecresult.go:37` - Function InsertClusterExecResult exceeds 15 lines (25 lines). **Fix**: Extract helper functions or early returns.
240. **long-func**: `.\gitmap\db\clusterexecresult.go:73` - Function UpdateClusterExecResult exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
241. **long-func**: `.\gitmap\db\clusternode.go:24` - Function InsertOrUpdateClusterNode exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
242. **long-func**: `.\gitmap\db\clusterrun.go:25` - Function InsertClusterRun exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
243. **long-func**: `.\gitmap\db\clusterrun.go:71` - Function SelectClusterRun exceeds 15 lines (26 lines). **Fix**: Extract helper functions or early returns.
244. **long-func**: `.\gitmap\db\enums.go:25` - Function String exceeds 15 lines (28 lines). **Fix**: Extract helper functions or early returns.
245. **long-func**: `.\gitmap\db\enums.go:56` - Function ParseCommandKind exceeds 15 lines (28 lines). **Fix**: Extract helper functions or early returns.
246. **long-func**: `.\gitmap\db\enums.go:98` - Function String exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
247. **long-func**: `.\gitmap\db\enums.go:117` - Function ParseResultStatus exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
248. **magic-value**: `.\gitmap\desktop\resolve.go:48` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
249. **enum-suffix**: `.\gitmap\diff\tree.go:13` - Type alias EntryKind acting as enum lacks Type suffix. **Fix**: Rename to EntryKindType.
250. **magic-value**: `.\gitmap\diff\tree.go:92` - Magic string/number in condition: if !opts.IncludeNodeModules && base == "node_modules" {. **Fix**: Extract to named constant or EnumType.
251. **long-func**: `.\gitmap\downloaderconfig\downloaderconfig.go:61` - Function Defaults exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
252. **enum-suffix**: `.\gitmap\errreport\errreport.go:35` - Type alias Phase acting as enum lacks Type suffix. **Fix**: Rename to PhaseType.
253. **magic-value**: `.\gitmap\formatter\clonescript.go:49` - Magic string/number in condition: if r.IdentifiedTransport == "ssh" {. **Fix**: Extract to named constant or EnumType.
254. **long-func**: `.\gitmap\formatter\desktopscript.go:24` - Function writeDesktopHeader exceeds 15 lines (35 lines). **Fix**: Extract helper functions or early returns.
255. **long-func**: `.\gitmap\formatter\desktopscript.go:69` - Function writeDesktopEntry exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
256. **long-func**: `.\gitmap\formatter\formatter_test.go:12` - Function testRecords exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
257. **magic-value**: `.\gitmap\formatter\formatter_test.go:87` - Magic string/number in condition: if len(parsed) == 2 {. **Fix**: Extract to named constant or EnumType.
258. **magic-value**: `.\gitmap\formatter\formatter_test.go:107` - Magic string/number in condition: if len(parsed) == 2 {. **Fix**: Extract to named constant or EnumType.
259. **long-func**: `.\gitmap\formatter\scangolden_contract_test.go:39` - Function canonicalScanRecords exceeds 15 lines (34 lines). **Fix**: Extract helper functions or early returns.
260. **long-func**: `.\gitmap\formatter\validate_test.go:12` - Function TestValidateRecords_TableDriven exceeds 15 lines (76 lines). **Fix**: Extract helper functions or early returns.
261. **enum-suffix**: `.\gitmap\ghtoken\ghtoken.go:27` - Type alias Source acting as enum lacks Type suffix. **Fix**: Rename to SourceType.
262. **long-func**: `.\gitmap\gitutil\branchdetect_test.go:10` - Function TestParseLsRemoteSymref exceeds 15 lines (49 lines). **Fix**: Extract helper functions or early returns.
263. **magic-value**: `.\gitmap\gitutil\gitutil.go:279` - Magic string/number in condition: if len(parts) == 2 {. **Fix**: Extract to named constant or EnumType.
264. **magic-value**: `.\gitmap\glyphs\autodetect_windows.go:24` - Magic string/number in condition: if os.Getenv("TERM_PROGRAM") == "vscode" {. **Fix**: Extract to named constant or EnumType.
265. **magic-value**: `.\gitmap\glyphs\autodetect_windows.go:27` - Magic string/number in condition: if os.Getenv("ConEmuANSI") == "ON" {. **Fix**: Extract to named constant or EnumType.
266. **long-func**: `.\gitmap\glyphs\filter.go:22` - Function buildTable exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
267. **enum-suffix**: `.\gitmap\glyphs\glyphs.go:18` - Type alias Mode acting as enum lacks Type suffix. **Fix**: Rename to ModeType.
268. **enum-suffix**: `.\gitmap\logging\jsonlog.go:20` - Type alias Level acting as enum lacks Type suffix. **Fix**: Rename to LevelType.
269. **nested-if**: `.\gitmap\macro\execute.go:26` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
270. **magic-value**: `.\gitmap\macro\execute.go:46` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
271. **magic-value**: `.\gitmap\macro\record.go:43` - Magic string/number in condition: if cmdText == "stop" || cmdText == "exit" || cmdText == "quit" {. **Fix**: Extract to named constant or EnumType.
272. **magic-value**: `.\gitmap\macro\record.go:69` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
273. **long-func**: `.\gitmap\mapper\mapper.go:88` - Function buildOneRecord exceeds 15 lines (23 lines). **Fix**: Extract helper functions or early returns.
274. **magic-value**: `.\gitmap\mapper\mapper.go:166` - Magic string/number in condition: if len(parts) == 2 {. **Fix**: Extract to named constant or EnumType.
275. **magic-value**: `.\gitmap\mapper\mapper.go:177` - Magic string/number in condition: if len(parts) == 2 {. **Fix**: Extract to named constant or EnumType.
276. **magic-value**: `.\gitmap\mapper\mapper_relroot_test.go:73` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
277. **magic-value**: `.\gitmap\mapper\mapper_test.go:61` - Magic string/number in condition: if name == "unknown" {. **Fix**: Extract to named constant or EnumType.
278. **enum-suffix**: `.\gitmap\movemerge\conflict.go:11` - Type alias Choice acting as enum lacks Type suffix. **Fix**: Rename to ChoiceType.
279. **enum-suffix**: `.\gitmap\movemerge\diff.go:6` - Type alias DiffKind acting as enum lacks Type suffix. **Fix**: Rename to DiffKindType.
280. **long-func**: `.\gitmap\movemerge\integration_test.go:21` - Function TestRunMerge_PreferNewer_BothSidesByteEqual exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
281. **enum-suffix**: `.\gitmap\movemerge\types.go:12` - Type alias EndpointKind acting as enum lacks Type suffix. **Fix**: Rename to EndpointKindType.
282. **enum-suffix**: `.\gitmap\movemerge\types.go:34` - Type alias PreferPolicy acting as enum lacks Type suffix. **Fix**: Rename to PreferPolicyType.
283. **enum-suffix**: `.\gitmap\movemerge\types.go:50` - Type alias Direction acting as enum lacks Type suffix. **Fix**: Rename to DirectionType.
284. **magic-value**: `.\gitmap\movemerge\walk.go:55` - Magic string/number in condition: if base == "node_modules" {. **Fix**: Extract to named constant or EnumType.
285. **long-func**: `.\gitmap\release\assetsbuild.go:14` - Function buildSingleTarget exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
286. **magic-value**: `.\gitmap\release\assetsbuild.go:48` - Magic string/number in condition: if target.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
287. **magic-value**: `.\gitmap\release\autocommit.go:176` - Magic string/number in condition: if answer == "y" || answer == "yes" {. **Fix**: Extract to named constant or EnumType.
288. **long-func**: `.\gitmap\release\autocommit_test.go:11` - Function TestIsNonFastForwardPushError exceeds 15 lines (22 lines). **Fix**: Extract helper functions or early returns.
289. **long-func**: `.\gitmap\release\remoteorigin_test.go:68` - Function TestParseInvalidURL exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
290. **enum-suffix**: `.\gitmap\render\preflag.go:22` - Type alias PrettyMode acting as enum lacks Type suffix. **Fix**: Rename to PrettyModeType.
291. **long-func**: `.\gitmap\render\pretty_emit.go:9` - Function emitBlock exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
292. **long-func**: `.\gitmap\render\pretty_parse.go:27` - Function parse exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
293. **long-func**: `.\gitmap\setup\pathsnippetwriter.go:88` - Function rewriteSnippetBlock exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
294. **enum-suffix**: `.\gitmap\startup\add.go:65` - Type alias AddStatus acting as enum lacks Type suffix. **Fix**: Rename to AddStatusType.
295. **magic-value**: `.\gitmap\startup\add.go:115` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
296. **magic-value**: `.\gitmap\startup\add.go:141` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
297. **magic-value**: `.\gitmap\startup\add.go:189` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
298. **long-func**: `.\gitmap\startup\add_darwin_test.go:25` - Function writeRawPlist exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
299. **magic-value**: `.\gitmap\startup\add_test.go:128` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
300. **magic-value**: `.\gitmap\startup\desktop.go:33` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
301. **magic-value**: `.\gitmap\startup\lifecycle_integration_test.go:43` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
302. **magic-value**: `.\gitmap\startup\lifecycle_integration_test.go:46` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
303. **magic-value**: `.\gitmap\startup\lifecycle_integration_test.go:58` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
304. **magic-value**: `.\gitmap\startup\plist.go:143` - Magic string/number in condition: if s.pendingKey == "ProgramArguments" {. **Fix**: Extract to named constant or EnumType.
305. **magic-value**: `.\gitmap\startup\plist.go:153` - Magic string/number in condition: if s.pendingKey == "Program" {. **Fix**: Extract to named constant or EnumType.
306. **magic-value**: `.\gitmap\startup\plistxml.go:47` - Magic string/number in condition: if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "string" {. **Fix**: Extract to named constant or EnumType.
307. **magic-value**: `.\gitmap\startup\plistxml.go:51` - Magic string/number in condition: if end, ok := tok.(xml.EndElement); ok && end.Name.Local == "array" {. **Fix**: Extract to named constant or EnumType.
308. **enum-suffix**: `.\gitmap\startup\remove.go:33` - Type alias RemoveStatus acting as enum lacks Type suffix. **Fix**: Rename to RemoveStatusType.
309. **magic-value**: `.\gitmap\startup\remove.go:88` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
310. **magic-value**: `.\gitmap\startup\remove.go:175` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
311. **magic-value**: `.\gitmap\startup\remove.go:189` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
312. **magic-value**: `.\gitmap\startup\startup.go:71` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
313. **magic-value**: `.\gitmap\startup\startup.go:75` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
314. **magic-value**: `.\gitmap\startup\startup.go:111` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
315. **magic-value**: `.\gitmap\startup\startup_test.go:19` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
316. **magic-value**: `.\gitmap\startup\startup_test.go:22` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
317. **magic-value**: `.\gitmap\startup\startup_test.go:84` - Magic string/number in condition: if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
318. **enum-suffix**: `.\gitmap\startup\winbackend.go:28` - Type alias Backend acting as enum lacks Type suffix. **Fix**: Rename to BackendType.
319. **magic-value**: `.\gitmap\startup\winbackend.go:96` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
320. **long-func**: `.\gitmap\store\cloneinteractiveselection.go:18` - Function SaveClonePickSelection exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
321. **long-func**: `.\gitmap\store\cloneinteractiveselection_load.go:53` - Function scanClonePickRow exceeds 15 lines (27 lines). **Fix**: Extract helper functions or early returns.
322. **magic-value**: `.\gitmap\store\downloader_seed.go:85` - Magic string/number in condition: if arg == "version" || arg == "--version" || arg == "-v" {. **Fix**: Extract to named constant or EnumType.
323. **magic-value**: `.\gitmap\store\migrateids.go:54` - Magic string/number in condition: if name == "Id" && colType == "TEXT" {. **Fix**: Extract to named constant or EnumType.
324. **long-func**: `.\gitmap\store\migrateids.go:100` - Function rebuildReposTable exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
325. **long-func**: `.\gitmap\store\migrate_commitin.go:21` - Function commitInDDL exceeds 15 lines (39 lines). **Fix**: Extract helper functions or early returns.
326. **long-func**: `.\gitmap\store\migrate_commitin_replaymap_helpers_test.go:20` - Function seedDeps exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
327. **long-func**: `.\gitmap\store\migrate_commitin_replaymap_test.go:45` - Function TestCommitInReplayMapHasAllSpecColumns exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
328. **long-func**: `.\gitmap\store\migrate_v15phase2.go:21` - Function migrateV15Phase2 exceeds 15 lines (38 lines). **Fix**: Extract helper functions or early returns.
329. **long-func**: `.\gitmap\store\migrate_v15phase3.go:19` - Function migrateV15Phase3 exceeds 15 lines (56 lines). **Fix**: Extract helper functions or early returns.
330. **long-func**: `.\gitmap\store\migrate_v15phase4.go:54` - Function zipGroupSpecs exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
331. **long-func**: `.\gitmap\store\migrate_v15phase4.go:81` - Function projectFamilySpecs exceeds 15 lines (45 lines). **Fix**: Extract helper functions or early returns.
332. **long-func**: `.\gitmap\store\migrate_v15phase4.go:151` - Function csharpFamilySpecs exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
333. **long-func**: `.\gitmap\store\migrate_v15phase4.go:185` - Function taskFamilySpecs exceeds 15 lines (50 lines). **Fix**: Extract helper functions or early returns.
334. **long-func**: `.\gitmap\store\migrate_v15phase4.go:266` - Function historyFamilySpecs exceeds 15 lines (34 lines). **Fix**: Extract helper functions or early returns.
335. **long-func**: `.\gitmap\store\store.go:376` - Function Reset exceeds 15 lines (52 lines). **Fix**: Extract helper functions or early returns.
336. **enum-suffix**: `.\gitmap\templates\diff.go:23` - Type alias DiffStatus acting as enum lacks Type suffix. **Fix**: Rename to DiffStatusType.
337. **magic-value**: `.\gitmap\templates\list_test.go:20` - Magic string/number in condition: if e.Kind == kindIgnore && e.Lang == "go" {. **Fix**: Extract to named constant or EnumType.
338. **magic-value**: `.\gitmap\templates\list_test.go:23` - Magic string/number in condition: if e.Kind == kindIgnore && e.Lang == "go" && e.Source != SourceEmbed {. **Fix**: Extract to named constant or EnumType.
339. **enum-suffix**: `.\gitmap\templates\merge.go:28` - Type alias MergeOutcome acting as enum lacks Type suffix. **Fix**: Rename to MergeOutcomeType.
340. **enum-suffix**: `.\gitmap\templates\resolver.go:12` - Type alias Source acting as enum lacks Type suffix. **Fix**: Rename to SourceType.
341. **magic-value**: `.\gitmap\tests\cmd_test\amend_test.go:78` - Magic string/number in condition: if commitHash == "HEAD" {. **Fix**: Extract to named constant or EnumType.
342. **long-func**: `.\gitmap\tests\cmd_test\seowritecreate_test.go:111` - Function buildSampleTemplateHelper exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
343. **long-func**: `.\gitmap\tests\cmd_test\seowrite_test.go:48` - Function TestSEOWriteConstants_FlagNames exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
344. **long-func**: `.\gitmap\tests\constants_test\seo_constants_test.go:49` - Function TestSEOErrorMessages_NonEmpty exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
345. **long-func**: `.\gitmap\tests\constants_test\seo_constants_test.go:75` - Function TestSEOHelpStrings_NonEmpty exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
346. **long-func**: `.\gitmap\tests\fixrepo_test\fixture_helpers_test.go:23` - Function setupFixtureRepo exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
347. **magic-value**: `.\gitmap\tests\fixrepo_test\gofmt_e2e_test.go:100` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
348. **long-func**: `.\gitmap\tests\scanclone_test\e2e_test.go:36` - Function e2eScanRecords exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
349. **magic-value**: `.\gitmap\tests\store_test\location_test.go:37` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
350. **enum-suffix**: `.\gitmap\theme\theme.go:35` - Type alias Mode acting as enum lacks Type suffix. **Fix**: Rename to ModeType.
351. **long-func**: `.\gitmap\tui\dashboard.go:61` - Function refreshStatuses exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
352. **long-func**: `.\gitmap\tui\logsview.go:61` - Function viewDetail exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
353. **long-func**: `.\gitmap\tui\releases.go:122` - Function viewDetail exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
354. **long-func**: `.\gitmap\tui\tui.go:114` - Function updateActiveView exceeds 15 lines (38 lines). **Fix**: Extract helper functions or early returns.
355. **long-func**: `.\gitmap\tui\tuiview.go:51` - Function renderContent exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
356. **long-func**: `.\gitmap\tui\tuiview.go:76` - Function renderStatusBar exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
357. **magic-value**: `.\gitmap\vscodepm\io.go:73` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
358. **magic-value**: `.\gitmap\vscodepm\path_test.go:128` - Magic string/number in condition: if runtime.GOOS == "darwin" {. **Fix**: Extract to named constant or EnumType.
359. **magic-value**: `.\gitmap\vscodepm\path_test.go:139` - Magic string/number in condition: if runtime.GOOS == "windows" {. **Fix**: Extract to named constant or EnumType.
360. **magic-value**: `.\remotion-demo\src\Terminal.tsx:34` - Magic string/number in condition: if (ln.kind === "prompt") {. **Fix**: Extract to named constant or EnumType.
361. **magic-value**: `.\remotion-demo\src\Terminal.tsx:164` - Magic string/number in condition: if (ln.kind === "blank") {. **Fix**: Extract to named constant or EnumType.
362. **magic-value**: `.\remotion-demo\src\Terminal.tsx:168` - Magic string/number in condition: if (ln.kind === "prompt") {. **Fix**: Extract to named constant or EnumType.
363. **enum-suffix**: `.\scripts\changelog\internal\runner\args.go:11` - Type alias Mode acting as enum lacks Type suffix. **Fix**: Rename to ModeType.
364. **long-func**: `.\src\components\docs\CodeBlock.tsx:102` - Function CodeBlock exceeds 15 lines (30 lines). **Fix**: Extract helper functions or early returns.
365. **long-func**: `.\src\components\docs\CommandBubbles.tsx:30` - Function CommandBubbles exceeds 15 lines (53 lines). **Fix**: Extract helper functions or early returns.
366. **long-func**: `.\src\components\docs\CommandCard.tsx:18` - Function CommandCard exceeds 15 lines (86 lines). **Fix**: Extract helper functions or early returns.
367. **long-func**: `.\src\components\docs\CommandCategoryGroup.tsx:18` - Function CommandCategoryGroup exceeds 15 lines (44 lines). **Fix**: Extract helper functions or early returns.
368. **long-func**: `.\src\components\docs\CommitTransferPage.tsx:124` - Function CommitTransferPage exceeds 15 lines (167 lines). **Fix**: Extract helper functions or early returns.
369. **long-func**: `.\src\components\docs\DocsLayout.tsx:14` - Function DocsLayout exceeds 15 lines (81 lines). **Fix**: Extract helper functions or early returns.
370. **long-func**: `.\src\components\docs\DocsSidebar.tsx:133` - Function DocsSidebar exceeds 15 lines (59 lines). **Fix**: Extract helper functions or early returns.
371. **magic-value**: `.\src\components\docs\DocsTooltip.tsx:42` - Magic string/number in condition: if (typeof label === "string") return label;. **Fix**: Extract to named constant or EnumType.
372. **long-func**: `.\src\components\docs\DocsTooltip.tsx:94` - Function DocsTooltip exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
373. **long-func**: `.\src\components\docs\InstallBlock.tsx:15` - Function CopyLine exceeds 15 lines (21 lines). **Fix**: Extract helper functions or early returns.
374. **long-func**: `.\src\components\docs\InstallBlock.tsx:40` - Function InstallBlock exceeds 15 lines (26 lines). **Fix**: Extract helper functions or early returns.
375. **long-func**: `.\src\components\docs\InstallHelpSection.tsx:53` - Function InstallHelpSection exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
376. **long-func**: `.\src\components\docs\SpecPage.tsx:20` - Function SpecPage exceeds 15 lines (94 lines). **Fix**: Extract helper functions or early returns.
377. **magic-value**: `.\src\components\docs\TabOrderMap.tsx:51` - Magic string/number in condition: if (el.getAttribute("aria-hidden") === "true") return false;. **Fix**: Extract to named constant or EnumType.
378. **magic-value**: `.\src\components\docs\TabOrderMap.tsx:63` - Magic string/number in condition: if (s.display === "none") return false;. **Fix**: Extract to named constant or EnumType.
379. **magic-value**: `.\src\components\docs\TabOrderMap.tsx:64` - Magic string/number in condition: if (s.visibility === "hidden" || s.visibility === "collapse") return false;. **Fix**: Extract to named constant or EnumType.
380. **magic-value**: `.\src\components\docs\TabOrderMap.tsx:66` - Magic string/number in condition: if (s.pointerEvents === "none" && node === el) {. **Fix**: Extract to named constant or EnumType.
381. **long-func**: `.\src\components\docs\TabOrderMap.tsx:119` - Function labelFor exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
382. **long-func**: `.\src\components\docs\TabOrderMap.tsx:196` - Function getTabOrder exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
383. **long-func**: `.\src\components\projects\ProjectDetailDialog.tsx:46` - Function ProjectDetailDialog exceeds 15 lines (50 lines). **Fix**: Extract helper functions or early returns.
384. **long-func**: `.\src\components\projects\ProjectDetailDialog.tsx:103` - Function GoMetadataSection exceeds 15 lines (35 lines). **Fix**: Extract helper functions or early returns.
385. **long-func**: `.\src\components\projects\ProjectDetailDialog.tsx:141` - Function CSharpMetadataSection exceeds 15 lines (36 lines). **Fix**: Extract helper functions or early returns.
386. **long-func**: `.\src\components\projects\RepoGroup.tsx:15` - Function RepoGroup exceeds 15 lines (40 lines). **Fix**: Extract helper functions or early returns.
387. **long-func**: `.\src\components\ui\alert-dialog.tsx:51` - Function AlertDialogFooter exceeds 15 lines (48 lines). **Fix**: Extract helper functions or early returns.
388. **long-func**: `.\src\components\ui\breadcrumb.tsx:69` - Function BreadcrumbEllipsis exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
389. **long-func**: `.\src\components\ui\calendar.tsx:10` - Function Calendar exceeds 15 lines (41 lines). **Fix**: Extract helper functions or early returns.
390. **magic-value**: `.\src\components\ui\carousel.tsx:77` - Magic string/number in condition: if (event.key === "ArrowLeft") {. **Fix**: Extract to named constant or EnumType.
391. **magic-value**: `.\src\components\ui\chart.tsx:285` - Magic string/number in condition: if (typeof payload !== "object" || payload === null) {. **Fix**: Extract to named constant or EnumType.
392. **magic-value**: `.\src\components\ui\chart.tsx:296` - Magic string/number in condition: if (key in payload && typeof payload[key as keyof typeof payload] === "string") {. **Fix**: Extract to named constant or EnumType.
393. **long-func**: `.\src\components\ui\dialog.tsx:59` - Function DialogFooter exceeds 15 lines (33 lines). **Fix**: Extract helper functions or early returns.
394. **long-func**: `.\src\components\ui\drawer.tsx:51` - Function DrawerFooter exceeds 15 lines (33 lines). **Fix**: Extract helper functions or early returns.
395. **long-func**: `.\src\components\ui\sheet.tsx:75` - Function SheetFooter exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
396. **magic-value**: `.\src\components\ui\sidebar.tsx:484` - Magic string/number in condition: if (typeof tooltip === "string") {. **Fix**: Extract to named constant or EnumType.
397. **long-func**: `.\src\components\ui\sonner.tsx:6` - Function Toaster exceeds 15 lines (18 lines). **Fix**: Extract helper functions or early returns.
398. **long-func**: `.\src\components\ui\toaster.tsx:4` - Function Toaster exceeds 15 lines (19 lines). **Fix**: Extract helper functions or early returns.
399. **long-func**: `.\src\hooks\use-toast.ts:145` - Function dismiss exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
400. **magic-value**: `.\src\hooks\useTheme.ts:49` - Magic string/number in condition: if (next === "light" || next === "dark") setThemeState(next);. **Fix**: Extract to named constant or EnumType.
401. **magic-value**: `.\src\hooks\useTheme.ts:54` - Magic string/number in condition: if (next === "system" || next === "user") setSourceState(next);. **Fix**: Extract to named constant or EnumType.
402. **magic-value**: `.\src\hooks\useTheme.ts:59` - Magic string/number in condition: if (event.newValue === "light" || event.newValue === "dark") {. **Fix**: Extract to named constant or EnumType.
403. **magic-value**: `.\src\lib\clipboard.ts:32` - Magic string/number in condition: if (typeof document === "undefined") {. **Fix**: Extract to named constant or EnumType.
404. **magic-value**: `.\src\lib\theme.ts:13` - Magic string/number in condition: if (typeof document === "undefined") return ThemeType.Dark;. **Fix**: Extract to named constant or EnumType.
405. **magic-value**: `.\src\lib\theme.ts:18` - Magic string/number in condition: if (res.isSuccess && (res.data === "light" || res.data === "dark")) {. **Fix**: Extract to named constant or EnumType.
406. **magic-value**: `.\src\lib\theme.ts:26` - Magic string/number in condition: if (typeof document === "undefined") return;. **Fix**: Extract to named constant or EnumType.
407. **long-func**: `.\src\pages\Architecture.tsx:4` - Function ArchitecturePage exceeds 15 lines (169 lines). **Fix**: Extract helper functions or early returns.
408. **long-func**: `.\src\pages\BatchActions.tsx:52` - Function BatchActions exceeds 15 lines (268 lines). **Fix**: Extract helper functions or early returns.
409. **long-func**: `.\src\pages\Changelog.tsx:60` - Function collapseAll exceeds 15 lines (152 lines). **Fix**: Extract helper functions or early returns.
410. **long-func**: `.\src\pages\ChangelogGenerate.tsx:5` - Function ChangelogGenerate exceeds 15 lines (155 lines). **Fix**: Extract helper functions or early returns.
411. **long-func**: `.\src\pages\ChromeProfileSpec.tsx:39` - Function ChromeProfileSpec exceeds 15 lines (139 lines). **Fix**: Extract helper functions or early returns.
412. **long-func**: `.\src\pages\ClearReleaseJSON.tsx:39` - Function ClearReleaseJSONPage exceeds 15 lines (185 lines). **Fix**: Extract helper functions or early returns.
413. **long-func**: `.\src\pages\CloneCommand.tsx:11` - Function CloneCommandPage exceeds 15 lines (98 lines). **Fix**: Extract helper functions or early returns.
414. **long-func**: `.\src\pages\CloneNext.tsx:5` - Function TerminalPreview exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
415. **long-func**: `.\src\pages\CloneNext.tsx:75` - Function CloneNextPage exceeds 15 lines (276 lines). **Fix**: Extract helper functions or early returns.
416. **long-func**: `.\src\pages\CloneNextCommand.tsx:23` - Function CloneNextCommandPage exceeds 15 lines (368 lines). **Fix**: Extract helper functions or early returns.
417. **long-func**: `.\src\pages\CloneOverview.tsx:16` - Function CloneOverviewPage exceeds 15 lines (125 lines). **Fix**: Extract helper functions or early returns.
418. **long-func**: `.\src\pages\Commands.tsx:10` - Function CommandsPage exceeds 15 lines (29 lines). **Fix**: Extract helper functions or early returns.
419. **long-func**: `.\src\pages\CommitIn.tsx:13` - Function CommitInPage exceeds 15 lines (165 lines). **Fix**: Extract helper functions or early returns.
420. **long-func**: `.\src\pages\Config.tsx:4` - Function ConfigPage exceeds 15 lines (53 lines). **Fix**: Extract helper functions or early returns.
421. **long-func**: `.\src\pages\DesignSystem.tsx:108` - Function DesignSystemPage exceeds 15 lines (267 lines). **Fix**: Extract helper functions or early returns.
422. **long-func**: `.\src\pages\FlagReference.tsx:59` - Function sortIndicator exceeds 15 lines (62 lines). **Fix**: Extract helper functions or early returns.
423. **long-func**: `.\src\pages\GettingStarted.tsx:4` - Function GettingStartedPage exceeds 15 lines (90 lines). **Fix**: Extract helper functions or early returns.
424. **long-func**: `.\src\pages\GoMod.tsx:28` - Function GoModPage exceeds 15 lines (225 lines). **Fix**: Extract helper functions or early returns.
425. **long-func**: `.\src\pages\HelpIndex.tsx:143` - Function HelpIndexPage exceeds 15 lines (101 lines). **Fix**: Extract helper functions or early returns.
426. **long-func**: `.\src\pages\Index.tsx:37` - Function HomePage exceeds 15 lines (103 lines). **Fix**: Extract helper functions or early returns.
427. **long-func**: `.\src\pages\Install.tsx:141` - Function InstallPage exceeds 15 lines (363 lines). **Fix**: Extract helper functions or early returns.
428. **long-func**: `.\src\pages\InstallGitmap.tsx:76` - Function InstallGitmapPage exceeds 15 lines (153 lines). **Fix**: Extract helper functions or early returns.
429. **long-func**: `.\src\pages\InteractiveExamples.tsx:104` - Function InteractiveExamplesPage exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
430. **long-func**: `.\src\pages\InteractiveTUI.tsx:61` - Function InteractiveTUI exceeds 15 lines (138 lines). **Fix**: Extract helper functions or early returns.
431. **long-func**: `.\src\pages\Makefile.tsx:47` - Function MakefilePage exceeds 15 lines (56 lines). **Fix**: Extract helper functions or early returns.
432. **long-func**: `.\src\pages\NotFound.tsx:4` - Function NotFound exceeds 15 lines (16 lines). **Fix**: Extract helper functions or early returns.
433. **long-func**: `.\src\pages\PostMortems.tsx:21` - Function PostMortemsPage exceeds 15 lines (88 lines). **Fix**: Extract helper functions or early returns.
434. **long-func**: `.\src\pages\Release.tsx:79` - Function ReleasePage exceeds 15 lines (434 lines). **Fix**: Extract helper functions or early returns.
435. **long-func**: `.\src\pages\ReleaseSelf.tsx:51` - Function ReleaseSelfPage exceeds 15 lines (141 lines). **Fix**: Extract helper functions or early returns.
436. **long-func**: `.\src\pages\ReleaseVersion.tsx:17` - Function ReleaseVersionPage exceeds 15 lines (32 lines). **Fix**: Extract helper functions or early returns.
437. **long-func**: `.\src\pages\releaseVersionSnippets.ts:19` - Function buildReleaseSnippets exceeds 15 lines (17 lines). **Fix**: Extract helper functions or early returns.
438. **long-func**: `.\src\pages\ScanCommand.tsx:17` - Function ScanCommandPage exceeds 15 lines (84 lines). **Fix**: Extract helper functions or early returns.
439. **long-func**: `.\src\pages\Setup.tsx:147` - Function Setup exceeds 15 lines (125 lines). **Fix**: Extract helper functions or early returns.
440. **magic-value**: `.\src\pages\SpecIndex.tsx:36` - Magic string/number in condition: if (e.key === "Escape" && document.activeElement === inputRef.current) {. **Fix**: Extract to named constant or EnumType.
441. **long-func**: `.\src\pages\Troubleshooting.tsx:279` - Function Troubleshooting exceeds 15 lines (20 lines). **Fix**: Extract helper functions or early returns.
442. **nested-if**: `.\src\pages\Troubleshooting.tsx:340` - Nested if statement detected. **Fix**: Flatten using early return or guard clause.
443. **long-func**: `.\src\pages\Troubleshooting.tsx:517` - Function CopyFixButton exceeds 15 lines (35 lines). **Fix**: Extract helper functions or early returns.
444. **long-func**: `.\src\pages\Troubleshooting.tsx:558` - Function CopyLinkButton exceeds 15 lines (34 lines). **Fix**: Extract helper functions or early returns.
445. **long-func**: `.\src\pages\Troubleshooting.tsx:599` - Function DiagnosticChecklist exceeds 15 lines (75 lines). **Fix**: Extract helper functions or early returns.
446. **long-func**: `.\src\pages\Troubleshooting.tsx:679` - Function ChecklistCommand exceeds 15 lines (27 lines). **Fix**: Extract helper functions or early returns.
447. **long-func**: `.\src\pages\VersionHistory.tsx:5` - Function TerminalPreview exceeds 15 lines (24 lines). **Fix**: Extract helper functions or early returns.
448. **long-func**: `.\src\pages\VersionHistory.tsx:38` - Function VersionHistoryPage exceeds 15 lines (212 lines). **Fix**: Extract helper functions or early returns.
449. **long-func**: `.\src\pages\Watch.tsx:35` - Function TerminalPreview exceeds 15 lines (90 lines). **Fix**: Extract helper functions or early returns.
450. **long-func**: `.\src\pages\Watch.tsx:154` - Function WatchPage exceeds 15 lines (104 lines). **Fix**: Extract helper functions or early returns.
451. **magic-value**: `.\src\types\helpJson.ts:30` - Magic string/number in condition: if (!value || typeof value !== "object") return false;. **Fix**: Extract to named constant or EnumType.
452. **magic-value**: `.\src\types\helpJson.ts:32` - Magic string/number in condition: if (typeof v.version !== "string") return false;. **Fix**: Extract to named constant or EnumType.
453. **magic-value**: `.\src\types\helpJson.ts:33` - Magic string/number in condition: if (typeof v.count !== "number") return false;. **Fix**: Extract to named constant or EnumType.
