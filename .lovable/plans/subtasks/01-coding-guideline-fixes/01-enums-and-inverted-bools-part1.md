# Subtask: Subagent 1: Enum Suffixes, Swallowed Errors & Inverted Booleans (Part 1)

- parent_plan: 01-coding-guideline-fixes.md
- steps: 1 to 50
- status: pending

## Execution Instructions
- Strictly enforce the <= 15 lines rule for any modified or extracted functions.
- Wrap errors with `apperror.Wrap` and domain context.
- Never use negative booleans or `!isX` inverted logic.
- Place exactly one blank line before `return` statements and after closing `}` braces.

## Granular Steps

1. **enum-suffix**: `gitmap/archive/create.go:36` - Type alias `CompressionMode` acting as enum lacks `Type` suffix.
   - **Action**: Rename `CompressionMode` to `CompressionModeType`.
2. **enum-suffix**: `gitmap/cliexit/report.go:50` - Type alias `OutputMode` acting as enum lacks `Type` suffix.
   - **Action**: Rename `OutputMode` to `OutputModeType`.
3. **enum-suffix**: `gitmap/cmd/commitin/enums.go:15` - Type alias `ConflictMode` acting as enum lacks `Type` suffix.
   - **Action**: Rename `ConflictMode` to `ConflictModeType`.
4. **enum-suffix**: `gitmap/cmd/commitin/enums.go:37` - Type alias `InputKind` acting as enum lacks `Type` suffix.
   - **Action**: Rename `InputKind` to `InputKindType`.
5. **enum-suffix**: `gitmap/cmd/commitin/enums.go:62` - Type alias `RunStatus` acting as enum lacks `Type` suffix.
   - **Action**: Rename `RunStatus` to `RunStatusType`.
6. **enum-suffix**: `gitmap/cmd/commitin/enums.go:93` - Type alias `CommitOutcome` acting as enum lacks `Type` suffix.
   - **Action**: Rename `CommitOutcome` to `CommitOutcomeType`.
7. **enum-suffix**: `gitmap/cmd/commitin/enums.go:118` - Type alias `SkipReason` acting as enum lacks `Type` suffix.
   - **Action**: Rename `SkipReason` to `SkipReasonType`.
8. **enum-suffix**: `gitmap/cmd/commitin/enums.go:146` - Type alias `ExclusionKind` acting as enum lacks `Type` suffix.
   - **Action**: Rename `ExclusionKind` to `ExclusionKindType`.
9. **enum-suffix**: `gitmap/cmd/commitin/enums.go:168` - Type alias `MessageRuleKind` acting as enum lacks `Type` suffix.
   - **Action**: Rename `MessageRuleKind` to `MessageRuleKindType`.
10. **enum-suffix**: `gitmap/cmd/commitin/enums.go:193` - Type alias `FunctionIntelLanguage` acting as enum lacks `Type` suffix.
   - **Action**: Rename `FunctionIntelLanguage` to `FunctionIntelLanguageType`.
11. **enum-suffix**: `gitmap/cmd/commitin/workspace/source.go:25` - Type alias `SourceKind` acting as enum lacks `Type` suffix.
   - **Action**: Rename `SourceKind` to `SourceKindType`.
12. **enum-suffix**: `gitmap/vscodepm/mergemode.go:38` - Type alias `MergeMode` acting as enum lacks `Type` suffix.
   - **Action**: Rename `MergeMode` to `MergeModeType`.
13. **swallowed-error**: `gitmap/cmd/chromeprofile_copy.go:157` - Swallowed error variable: `_ = err`.
   - **Action**: Handle or wrap error using `apperror.Wrap`.
14. **inverted-bool**: `gitmap/archive/extract.go:105` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
15. **inverted-bool**: `gitmap/archive/extract.go:232` - Inverted boolean logic: `!isDir`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
16. **inverted-bool**: `gitmap/archive/list.go:42` - Inverted boolean logic: `!isExtractor`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
17. **inverted-bool**: `gitmap/cliexit/report_test.go:85` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
18. **inverted-bool**: `gitmap/cloneconcurrency/resolve_test.go:28` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
19. **inverted-bool**: `gitmap/cloneconcurrency/resolve_test.go:45` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
20. **inverted-bool**: `gitmap/clonefrom/execute.go:80` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
21. **inverted-bool**: `gitmap/clonefrom/execute.go:95` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
22. **inverted-bool**: `gitmap/clonefrom/execute_checkout_test.go:166` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
23. **inverted-bool**: `gitmap/clonefrom/execute_checkout_test.go:183` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
24. **inverted-bool**: `gitmap/clonefrom/execute_lfs_fix_test.go:34` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
25. **inverted-bool**: `gitmap/clonefrom/execute_lfs_fix_test.go:46` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
26. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:40` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
27. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:51` - Inverted boolean logic: `!hasKey`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
28. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:73` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
29. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:78` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
30. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:88` - Inverted boolean logic: `!hasKey`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
31. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:134` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
32. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:152` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
33. **inverted-bool**: `gitmap/clonefrom/jsonschema_test.go:162` - Inverted boolean logic: `!hasConst`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
34. **inverted-bool**: `gitmap/clonefrom/summary_provenance_test.go:71` - Inverted boolean logic: `!exists`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
35. **inverted-bool**: `gitmap/clonefrom/validate.go:47` - Inverted boolean logic: `!isValidBranchName`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
36. **inverted-bool**: `gitmap/clonefrom/validate.go:50` - Inverted boolean logic: `!isValidCheckout`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
37. **inverted-bool**: `gitmap/clonenext/remoteupdate.go:86` - Inverted boolean logic: `!exists`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
38. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:65` - Inverted boolean logic: `!isGitWorkTree`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
39. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:139` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
40. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:236` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
41. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:247` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
42. **inverted-bool**: `gitmap/clonenow/execute_idempotent.go:273` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
43. **inverted-bool**: `gitmap/clonenow/parse.go:211` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
44. **inverted-bool**: `gitmap/clonenow/parsetext.go:35` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
45. **inverted-bool**: `gitmap/clonenow/parsetext.go:56` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
46. **inverted-bool**: `gitmap/clonenow/parse_schema.go:152` - Inverted boolean logic: `!hasURL`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
47. **inverted-bool**: `gitmap/clonenow/parse_schema_json.go:62` - Inverted boolean logic: `!hasJSONURL`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
48. **inverted-bool**: `gitmap/cloner/cache.go:127` - Inverted boolean logic: `!ok`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
49. **inverted-bool**: `gitmap/cloner/lfs_retry_test.go:12` - Inverted boolean logic: `!isLFSSmudgeFailure`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
50. **inverted-bool**: `gitmap/cloner/runners.go:134` - Inverted boolean logic: `!isGitRepo`.
   - **Action**: Extract into positive boolean check or use explicit `== false`.
