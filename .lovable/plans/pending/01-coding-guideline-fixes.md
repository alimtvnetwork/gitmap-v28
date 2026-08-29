# Plan: Coding Guideline Audit & Enforcement (v4)

- slug: 01-coding-guideline-fixes
- status: pending
- steps_count: 150
- created: 2026-08-29

## Context
Comprehensive, deep multi-stage audit of the entire codebase for coding guideline violations, boolean anti-patterns, missing enum suffixes, cyclomatic nesting, and error-handling flaws across Go backend packages. Structured into exactly 150 granular steps.

## Fallout & Blast Radius Analysis
- **Enum Suffix Renames**: Modifying type definitions requires updating all call sites across CLI commands, tests, and struct field definitions. Blast radius is contained within Go internal packages.
- **Boolean Inversions**: Changing `!isX` to explicit `isX == false` or `isMissing` preserves runtime semantics without breaking API contracts or downstream scripts.
- **Nested If Flattening & Function Extraction**: Uses guard clauses and single-responsibility helper extractions to strictly comply with the <= 15-line function cap.
- **CI/CD Guard**: No CI/CD workflows, GitHub Actions, or validation rules will be bypassed.

## Enqueued Granular Tasks (Steps 1 to 150)

1. **enum-suffix**: `gitmap/archive/create.go:36` - Type alias `CompressionMode` acting as enum lacks `Type` suffix. **Fix**: Rename `CompressionMode` to `CompressionModeType`.
2. **enum-suffix**: `gitmap/cliexit/report.go:50` - Type alias `OutputMode` acting as enum lacks `Type` suffix. **Fix**: Rename `OutputMode` to `OutputModeType`.
3. **enum-suffix**: `gitmap/cmd/commitin/enums.go:15` - Type alias `ConflictMode` acting as enum lacks `Type` suffix. **Fix**: Rename `ConflictMode` to `ConflictModeType`.
4. **enum-suffix**: `gitmap/cmd/commitin/enums.go:37` - Type alias `InputKind` acting as enum lacks `Type` suffix. **Fix**: Rename `InputKind` to `InputKindType`.
5. **enum-suffix**: `gitmap/cmd/commitin/enums.go:62` - Type alias `RunStatus` acting as enum lacks `Type` suffix. **Fix**: Rename `RunStatus` to `RunStatusType`.
6. **enum-suffix**: `gitmap/cmd/commitin/enums.go:93` - Type alias `CommitOutcome` acting as enum lacks `Type` suffix. **Fix**: Rename `CommitOutcome` to `CommitOutcomeType`.
7. **enum-suffix**: `gitmap/cmd/commitin/enums.go:118` - Type alias `SkipReason` acting as enum lacks `Type` suffix. **Fix**: Rename `SkipReason` to `SkipReasonType`.
8. **enum-suffix**: `gitmap/cmd/commitin/enums.go:146` - Type alias `ExclusionKind` acting as enum lacks `Type` suffix. **Fix**: Rename `ExclusionKind` to `ExclusionKindType`.
9. **enum-suffix**: `gitmap/cmd/commitin/enums.go:168` - Type alias `MessageRuleKind` acting as enum lacks `Type` suffix. **Fix**: Rename `MessageRuleKind` to `MessageRuleKindType`.
10. **enum-suffix**: `gitmap/cmd/commitin/enums.go:193` - Type alias `FunctionIntelLanguage` acting as enum lacks `Type` suffix. **Fix**: Rename `FunctionIntelLanguage` to `FunctionIntelLanguageType`.
11. **enum-suffix**: `gitmap/cmd/commitin/workspace/source.go:25` - Type alias `SourceKind` acting as enum lacks `Type` suffix. **Fix**: Rename `SourceKind` to `SourceKindType`.
12. **enum-suffix**: `gitmap/vscodepm/mergemode.go:38` - Type alias `MergeMode` acting as enum lacks `Type` suffix. **Fix**: Rename `MergeMode` to `MergeModeType`.
13. **swallowed-error**: `gitmap/cmd/chromeprofile_copy.go:157` - Swallowed error variable: `_ = err`. **Fix**: Handle or wrap error using `apperror.Wrap`.
14. **inverted-bool**: `gitmap/archive/extract.go:105` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
15. **inverted-bool**: `gitmap/archive/extract.go:232` - Inverted boolean logic: `!isDir`. **Fix**: Extract into positive boolean check or use explicit `== false`.
16. **inverted-bool**: `gitmap/archive/list.go:42` - Inverted boolean logic: `!isExtractor`. **Fix**: Extract into positive boolean check or use explicit `== false`.
17. **inverted-bool**: `gitmap/cliexit/report_test.go:85` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
18. **inverted-bool**: `gitmap/cloneconcurrency/resolve_test.go:28` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
19. **inverted-bool**: `gitmap/cloneconcurrency/resolve_test.go:45` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
20. **inverted-bool**: `gitmap/clonefrom/execute.go:80` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
21. **inverted-bool**: `gitmap/clonefrom/execute.go:95` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
22. **inverted-bool**: `gitmap/clonefrom/execute_checkout_test.go:166` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
23. **inverted-bool**: `gitmap/clonefrom/execute_checkout_test.go:183` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
24. **inverted-bool**: `gitmap/clonefrom/execute_lfs_fix_test.go:34` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
25. **inverted-bool**: `gitmap/clonefrom/execute_lfs_fix_test.go:46` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
26. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:40` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
27. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:51` - Inverted boolean logic: `!hasKey`. **Fix**: Extract into positive boolean check or use explicit `== false`.
28. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:73` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
29. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:78` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
30. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:88` - Inverted boolean logic: `!hasKey`. **Fix**: Extract into positive boolean check or use explicit `== false`.
31. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:134` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
32. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:152` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
33. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:162` - Inverted boolean logic: `!hasConst`. **Fix**: Extract into positive boolean check or use explicit `== false`.
34. **inverted-bool**: `gitmap/clonefrom/summary_provenance_test.go:71` - Inverted boolean logic: `!exists`. **Fix**: Extract into positive boolean check or use explicit `== false`.
35. **inverted-bool**: `gitmap/clonefrom/validate.go:47` - Inverted boolean logic: `!isValidBranchName`. **Fix**: Extract into positive boolean check or use explicit `== false`.
36. **inverted-bool**: `gitmap/clonefrom/validate.go:50` - Inverted boolean logic: `!isValidCheckout`. **Fix**: Extract into positive boolean check or use explicit `== false`.
37. **inverted-bool**: `gitmap/clonenext/remoteupdate.go:86` - Inverted boolean logic: `!exists`. **Fix**: Extract into positive boolean check or use explicit `== false`.
38. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:65` - Inverted boolean logic: `!isGitWorkTree`. **Fix**: Extract into positive boolean check or use explicit `== false`.
39. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:139` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
40. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:236` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
41. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:247` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
42. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:273` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
43. **inverted-bool**: `gitmap/clonenow/parse.go:211` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
44. **inverted-bool**: `gitmap/clonenow/parsetext.go:35` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
45. **inverted-bool**: `gitmap/clonenow/parsetext.go:56` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
46. **inverted-bool**: `gitmap/clonenow/parse_schema.go:152` - Inverted boolean logic: `!hasURL`. **Fix**: Extract into positive boolean check or use explicit `== false`.
47. **inverted-bool**: `gitmap/clonenow/parse_schema_json.go:62` - Inverted boolean logic: `!hasJSONURL`. **Fix**: Extract into positive boolean check or use explicit `== false`.
48. **inverted-bool**: `gitmap/cloner/cache.go:127` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
49. **inverted-bool**: `gitmap/cloner/lfs_retry_test.go:12` - Inverted boolean logic: `!isLFSSmudgeFailure`. **Fix**: Extract into positive boolean check or use explicit `== false`.
50. **inverted-bool**: `gitmap/cloner/runners.go:134` - Inverted boolean logic: `!isGitRepo`. **Fix**: Extract into positive boolean check or use explicit `== false`.
51. **inverted-bool**: `gitmap/cloner/safe_pull.go:48` - Inverted boolean logic: `!isGitRepo`. **Fix**: Extract into positive boolean check or use explicit `== false`.
52. **inverted-bool**: `gitmap/cluster/exec_proj.go:94` - Inverted boolean logic: `!hasPs1`. **Fix**: Extract into positive boolean check or use explicit `== false`.
53. **inverted-bool**: `gitmap/cluster/exec_proj.go:98` - Inverted boolean logic: `!hasPs1`. **Fix**: Extract into positive boolean check or use explicit `== false`.
54. **inverted-bool**: `gitmap/cluster/exec_proj.go:98` - Inverted boolean logic: `!hasSh`. **Fix**: Extract into positive boolean check or use explicit `== false`.
55. **inverted-bool**: `gitmap/cluster/exec_ps.go:36` - Inverted boolean logic: `!isWin`. **Fix**: Extract into positive boolean check or use explicit `== false`.
56. **inverted-bool**: `gitmap/cluster/exec_ps.go:39` - Inverted boolean logic: `!isWin`. **Fix**: Extract into positive boolean check or use explicit `== false`.
57. **inverted-bool**: `gitmap/cluster/exec_ps.go:46` - Inverted boolean logic: `!isWin`. **Fix**: Extract into positive boolean check or use explicit `== false`.
58. **inverted-bool**: `gitmap/cluster/node_resolver.go:109` - Inverted boolean logic: `!hasSeparator`. **Fix**: Extract into positive boolean check or use explicit `== false`.
59. **inverted-bool**: `gitmap/cluster/node_resolver.go:127` - Inverted boolean logic: `!isValidRange`. **Fix**: Extract into positive boolean check or use explicit `== false`.
60. **inverted-bool**: `gitmap/cluster/node_resolver.go:136` - Inverted boolean logic: `!hasTrailing`. **Fix**: Extract into positive boolean check or use explicit `== false`.
61. **inverted-bool**: `gitmap/cluster/node_resolver.go:142` - Inverted boolean logic: `!isValidTrailing`. **Fix**: Extract into positive boolean check or use explicit `== false`.
62. **inverted-bool**: `gitmap/cluster/node_resolver.go:162` - Inverted boolean logic: `!isValidInt`. **Fix**: Extract into positive boolean check or use explicit `== false`.
63. **inverted-bool**: `gitmap/cluster/node_resolver.go:165` - Inverted boolean logic: `!isValidInt`. **Fix**: Extract into positive boolean check or use explicit `== false`.
64. **inverted-bool**: `gitmap/cluster/node_resolver.go:174` - Inverted boolean logic: `!hasTrailing`. **Fix**: Extract into positive boolean check or use explicit `== false`.
65. **inverted-bool**: `gitmap/cluster/node_resolver.go:180` - Inverted boolean logic: `!isValidTrailing`. **Fix**: Extract into positive boolean check or use explicit `== false`.
66. **inverted-bool**: `gitmap/cmd/agy_cmd.go:59` - Inverted boolean logic: `!hasEnoughArgs`. **Fix**: Extract into positive boolean check or use explicit `== false`.
67. **inverted-bool**: `gitmap/cmd/agy_cmd.go:88` - Inverted boolean logic: `!isCreated`. **Fix**: Extract into positive boolean check or use explicit `== false`.
68. **inverted-bool**: `gitmap/cmd/agy_cmd.go:119` - Inverted boolean logic: `!hasEnoughArgs`. **Fix**: Extract into positive boolean check or use explicit `== false`.
69. **inverted-bool**: `gitmap/cmd/agy_cmd.go:226` - Inverted boolean logic: `!hasEnoughArgs`. **Fix**: Extract into positive boolean check or use explicit `== false`.
70. **inverted-bool**: `gitmap/cmd/amendaudit_jsonschema_contract_test.go:52` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
71. **inverted-bool**: `gitmap/cmd/amendlist_jsonschema_contract_test.go:44` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
72. **inverted-bool**: `gitmap/cmd/amendlist_jsonschema_contract_test.go:64` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
73. **inverted-bool**: `gitmap/cmd/audit.go:15` - Inverted boolean logic: `!shouldAuditCommand`. **Fix**: Extract into positive boolean check or use explicit `== false`.
74. **inverted-bool**: `gitmap/cmd/audit.go:24` - Inverted boolean logic: `!shouldAudit`. **Fix**: Extract into positive boolean check or use explicit `== false`.
75. **inverted-bool**: `gitmap/cmd/auditlegacy.go:83` - Inverted boolean logic: `!isAuditScannable`. **Fix**: Extract into positive boolean check or use explicit `== false`.
76. **inverted-bool**: `gitmap/cmd/auditlegacy_report.go:119` - Inverted boolean logic: `!hasDiffs`. **Fix**: Extract into positive boolean check or use explicit `== false`.
77. **inverted-bool**: `gitmap/cmd/bookmarklist_jsonschema_contract_test.go:34` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
78. **inverted-bool**: `gitmap/cmd/bookmarklist_jsonschema_contract_test.go:54` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
79. **inverted-bool**: `gitmap/cmd/cddefault.go:40` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
80. **inverted-bool**: `gitmap/cmd/cfrppriorversion.go:29` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
81. **inverted-bool**: `gitmap/cmd/cg_resolver_test.go:11` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
82. **inverted-bool**: `gitmap/cmd/cg_worker.go:111` - Inverted boolean logic: `!hasFiles`. **Fix**: Extract into positive boolean check or use explicit `== false`.
83. **inverted-bool**: `gitmap/cmd/changelog.go:125` - Inverted boolean logic: `!found`. **Fix**: Extract into positive boolean check or use explicit `== false`.
84. **inverted-bool**: `gitmap/cmd/chromeprofile.go:39` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
85. **inverted-bool**: `gitmap/cmd/chromeprofile.go:174` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
86. **inverted-bool**: `gitmap/cmd/chromeprofile_csv.go:85` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
87. **inverted-bool**: `gitmap/cmd/chromeprofile_csv_test.go:123` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
88. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:58` - Inverted boolean logic: `!isKnownMergeWhat`. **Fix**: Extract into positive boolean check or use explicit `== false`.
89. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:64` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
90. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:70` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
91. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:318` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
92. **inverted-bool**: `gitmap/cmd/chromeprofile_preferences.go:55` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
93. **inverted-bool**: `gitmap/cmd/chromeprofile_register.go:94` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
94. **inverted-bool**: `gitmap/cmd/chromeprofile_register_test.go:63` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
95. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:39` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
96. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:50` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
97. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:58` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
98. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:67` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
99. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:74` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
100. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:81` - Inverted boolean logic: `!ok`. **Fix**: Extract into positive boolean check or use explicit `== false`.
101. **nested-if**: `gitmap/archive/extract.go:192` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
102. **nested-if**: `gitmap/archive/extract.go:233` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
103. **nested-if**: `gitmap/archive/extract.go:265` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
104. **nested-if**: `gitmap/archive/extract.go:293` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
105. **nested-if**: `gitmap/archive/source.go:175` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
106. **nested-if**: `gitmap/archive/source.go:195` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
107. **nested-if**: `gitmap/archive/source.go:199` - Nested `if` statement detected at indentation depth 4. **Fix**: Flatten with early returns or guard clauses.
108. **nested-if**: `gitmap/cliexit/report.go:183` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
109. **nested-if**: `gitmap/clonefrom/parsecsv.go:55` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
110. **nested-if**: `gitmap/clonefrom/render.go:113` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
111. **nested-if**: `gitmap/clonefrom/summary.go:43` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
112. **nested-if**: `gitmap/clonefrom/summary.go:47` - Nested `if` statement detected at indentation depth 4. **Fix**: Flatten with early returns or guard clauses.
113. **nested-if**: `gitmap/clonefrom/summary.go:117` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
114. **nested-if**: `gitmap/clonefrom/summary_terminal.go:51` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
115. **nested-if**: `gitmap/clonefrom/summary_terminal.go:73` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
116. **nested-if**: `gitmap/clonefrom/summary_terminal.go:146` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
117. **nested-if**: `gitmap/clonefrom/validate.go:44` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
118. **nested-if**: `gitmap/clonefrom/validate.go:47` - Nested `if` statement detected at indentation depth 4. **Fix**: Flatten with early returns or guard clauses.
119. **nested-if**: `gitmap/clonefrom/validate.go:50` - Nested `if` statement detected at indentation depth 5. **Fix**: Flatten with early returns or guard clauses.
120. **nested-if**: `gitmap/clonefrom/validate.go:181` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
121. **nested-if**: `gitmap/clonenext/localstate.go:123` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
122. **nested-if**: `gitmap/clonenext/repodetect.go:46` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
123. **nested-if**: `gitmap/clonenow/clonenow.go:102` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
124. **nested-if**: `gitmap/clonenow/execute.go:97` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
125. **nested-if**: `gitmap/clonenow/execute.go:141` - Nested `if` statement detected at indentation depth 3. **Fix**: Flatten with early returns or guard clauses.
126. **long-func**: `gitmap/archive/create.go:108` - Function `CreateArchive` exceeds 15 lines (20 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
127. **long-func**: `gitmap/archive/extract.go:38` - Function `CompactExtract` exceeds 15 lines (16 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
128. **long-func**: `gitmap/archive/extract.go:66` - Function `completeCompactExtract` exceeds 15 lines (20 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
129. **long-func**: `gitmap/archive/extract.go:92` - Function `extractAllIntoDir` exceeds 15 lines (20 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
130. **long-func**: `gitmap/archive/extract.go:125` - Function `extractArchiveEntry` exceeds 15 lines (17 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
131. **long-func**: `gitmap/archive/extract.go:177` - Function `safeJoin` exceeds 15 lines (21 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
132. **long-func**: `gitmap/archive/extract.go:202` - Function `promoteRealRoot` exceeds 15 lines (16 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
133. **long-func**: `gitmap/archive/extract.go:219` - Function `findDeepestRoot` exceeds 15 lines (22 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
134. **long-func**: `gitmap/archive/extract.go:282` - Function `copyDirEntry` exceeds 15 lines (17 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
135. **long-func**: `gitmap/archive/extract.go:327` - Function `archiveBaseName` exceeds 15 lines (16 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
136. **long-func**: `gitmap/archive/list.go:29` - Function `ListEntries` exceeds 15 lines (20 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
137. **long-func**: `gitmap/archive/list.go:50` - Function `extractListEntries` exceeds 15 lines (16 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
138. **long-func**: `gitmap/archive/source.go:96` - Function `ResolveSource` exceeds 15 lines (22 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
139. **long-func**: `gitmap/archive/source.go:131` - Function `resolveHTTP` exceeds 15 lines (21 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
140. **long-func**: `gitmap/archive/source.go:157` - Function `downloadWithAria2c` exceeds 15 lines (24 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
141. **long-func**: `gitmap/archive/source.go:184` - Function `downloadWithHTTP` exceeds 15 lines (23 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
142. **long-func**: `gitmap/archive/source.go:227` - Function `resolveGit` exceeds 15 lines (17 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
143. **long-func**: `gitmap/archive/source.go:249` - Function `AutoDetectSingleArchive` exceeds 15 lines (28 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
144. **long-func**: `gitmap/cliexit/report.go:149` - Function `sortedExtraLines` exceeds 15 lines (16 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
145. **long-func**: `gitmap/cliexit/report.go:168` - Function `writeJSON` exceeds 15 lines (24 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
146. **long-func**: `gitmap/clonefrom/execute_hooks.go:34` - Function `ExecuteWithHooks` exceeds 15 lines (16 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
147. **long-func**: `gitmap/clonefrom/execute_lfs_fix.go:39` - Function `executeLFSFix` exceeds 15 lines (38 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
148. **long-func**: `gitmap/clonefrom/jsonschema_helpers.go:41` - Function `rootSchema` exceeds 15 lines (17 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
149. **long-func**: `gitmap/clonefrom/parse.go:33` - Function `ParseFile` exceeds 15 lines (20 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
150. **long-func**: `gitmap/clonefrom/parse.go:81` - Function `parseJSON` exceeds 15 lines (17 lines). **Fix**: Extract helper functions, table-driven dispatch, or guard clauses.
