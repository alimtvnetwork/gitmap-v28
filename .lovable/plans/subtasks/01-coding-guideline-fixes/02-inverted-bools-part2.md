# Subtask: Subagent 2: Inverted Booleans & Negation Refactoring (Part 2)

- parent_plan: 01-coding-guideline-fixes.md
- steps: 51 to 100
- status: pending

## Execution Instructions

- Strictly enforce the <= 15 lines rule for any modified or extracted functions.
- Wrap errors with `apperror.Wrap` and domain context.
- Never use negative booleans or `!isX` inverted logic.
- Place exactly one blank line before `return` statements and after closing `}` braces.

## Granular Steps

51. **inverted-bool**: `gitmap/cloner/safe_pull.go:48` - Inverted boolean logic: `!isGitRepo`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
52. **inverted-bool**: `gitmap/cluster/exec_proj.go:94` - Inverted boolean logic: `!hasPs1`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
53. **inverted-bool**: `gitmap/cluster/exec_proj.go:98` - Inverted boolean logic: `!hasPs1`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
54. **inverted-bool**: `gitmap/cluster/exec_proj.go:98` - Inverted boolean logic: `!hasSh`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
55. **inverted-bool**: `gitmap/cluster/exec_ps.go:36` - Inverted boolean logic: `!isWin`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
56. **inverted-bool**: `gitmap/cluster/exec_ps.go:39` - Inverted boolean logic: `!isWin`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
57. **inverted-bool**: `gitmap/cluster/exec_ps.go:46` - Inverted boolean logic: `!isWin`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
58. **inverted-bool**: `gitmap/cluster/node_resolver.go:109` - Inverted boolean logic: `!hasSeparator`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
59. **inverted-bool**: `gitmap/cluster/node_resolver.go:127` - Inverted boolean logic: `!isValidRange`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
60. **inverted-bool**: `gitmap/cluster/node_resolver.go:136` - Inverted boolean logic: `!hasTrailing`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
61. **inverted-bool**: `gitmap/cluster/node_resolver.go:142` - Inverted boolean logic: `!isValidTrailing`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
62. **inverted-bool**: `gitmap/cluster/node_resolver.go:162` - Inverted boolean logic: `!isValidInt`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
63. **inverted-bool**: `gitmap/cluster/node_resolver.go:165` - Inverted boolean logic: `!isValidInt`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
64. **inverted-bool**: `gitmap/cluster/node_resolver.go:174` - Inverted boolean logic: `!hasTrailing`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
65. **inverted-bool**: `gitmap/cluster/node_resolver.go:180` - Inverted boolean logic: `!isValidTrailing`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
66. **inverted-bool**: `gitmap/cmd/agy_cmd.go:59` - Inverted boolean logic: `!hasEnoughArgs`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
67. **inverted-bool**: `gitmap/cmd/agy_cmd.go:88` - Inverted boolean logic: `!isCreated`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
68. **inverted-bool**: `gitmap/cmd/agy_cmd.go:119` - Inverted boolean logic: `!hasEnoughArgs`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
69. **inverted-bool**: `gitmap/cmd/agy_cmd.go:226` - Inverted boolean logic: `!hasEnoughArgs`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
70. **inverted-bool**: `gitmap/cmd/amendaudit_jsonschema_contract_test.go:52` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
71. **inverted-bool**: `gitmap/cmd/amendlist_jsonschema_contract_test.go:44` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
72. **inverted-bool**: `gitmap/cmd/amendlist_jsonschema_contract_test.go:64` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
73. **inverted-bool**: `gitmap/cmd/audit.go:15` - Inverted boolean logic: `!shouldAuditCommand`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
74. **inverted-bool**: `gitmap/cmd/audit.go:24` - Inverted boolean logic: `!shouldAudit`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
75. **inverted-bool**: `gitmap/cmd/auditlegacy.go:83` - Inverted boolean logic: `!isAuditScannable`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
76. **inverted-bool**: `gitmap/cmd/auditlegacy_report.go:119` - Inverted boolean logic: `!hasDiffs`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
77. **inverted-bool**: `gitmap/cmd/bookmarklist_jsonschema_contract_test.go:34` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
78. **inverted-bool**: `gitmap/cmd/bookmarklist_jsonschema_contract_test.go:54` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
79. **inverted-bool**: `gitmap/cmd/cddefault.go:40` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
80. **inverted-bool**: `gitmap/cmd/cfrppriorversion.go:29` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
81. **inverted-bool**: `gitmap/cmd/cg_resolver_test.go:11` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
82. **inverted-bool**: `gitmap/cmd/cg_worker.go:111` - Inverted boolean logic: `!hasFiles`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
83. **inverted-bool**: `gitmap/cmd/changelog.go:125` - Inverted boolean logic: `!found`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
84. **inverted-bool**: `gitmap/cmd/chromeprofile.go:39` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
85. **inverted-bool**: `gitmap/cmd/chromeprofile.go:174` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
86. **inverted-bool**: `gitmap/cmd/chromeprofile_csv.go:85` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
87. **inverted-bool**: `gitmap/cmd/chromeprofile_csv_test.go:123` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
88. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:58` - Inverted boolean logic: `!isKnownMergeWhat`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
89. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:64` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
90. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:70` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
91. **inverted-bool**: `gitmap/cmd/chromeprofile_merge.go:318` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
92. **inverted-bool**: `gitmap/cmd/chromeprofile_preferences.go:55` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
93. **inverted-bool**: `gitmap/cmd/chromeprofile_register.go:94` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
94. **inverted-bool**: `gitmap/cmd/chromeprofile_register_test.go:63` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
95. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:39` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
96. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:50` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
97. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:58` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
98. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:67` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
99. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:74` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
100. **inverted-bool**: `gitmap/cmd/chromeprofile_resolve_test.go:81` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
