# Subtask 03: Function Size Decomposition (<8 lines preferred, 15 lines cap)

Slug: 03-function-size-decomposition
Parent Plan: 01-coding-guideline-fixes
Status: pending

## Objective
Decompose monolithic functions exceeding 15 lines into concise, focused helper routines ($\le 8$ lines preferred, $\le 15$ lines hard cap).

## Concrete Execution Steps (30 Steps)

1. `gitmap/archive/archive_test.go:30`: Refactor `TestFormatFromPath_SingleExtensions` (16 lines) by extracting extension test table.
2. `gitmap/archive/extract.go:64`: Refactor `completeCompactExtract` (16 lines) into `validateCompactExtract` + `finalizeCompactExtract`.
3. `gitmap/archive/extract.go:174`: Refactor `safeJoin` (18 lines) into `resolveCleanPath` + `checkPathBoundary`.
4. `gitmap/archive/extract.go:216`: Refactor `findDeepestRoot` (20 lines) into `collectSubDirs` + `determineDeepestRoot`.
5. `gitmap/archive/source.go:92`: Refactor `ResolveSource` (16 lines) into `resolveHTTPSource` + `resolveLocalSource`.
6. `gitmap/archive/source.go:151`: Refactor `downloadWithAria2c` (21 lines) into `buildAria2cCmd` + `executeAria2cDownload`.
7. `gitmap/archive/source.go:177`: Refactor `downloadWithHTTP` (20 lines) into `initHTTPRequest` + `streamHTTPResponse`.
8. `gitmap/archive/source.go:241`: Refactor `AutoDetectSingleArchive` (25 lines) into `matchArchiveFiles` + `selectSingleArchive`.
9. `gitmap/cliexit/kind_test.go:17`: Refactor `TestKindCode_Table` (19 lines) by extracting table runner helper.
10. `gitmap/cliexit/kind_test.go:53`: Refactor `TestWithKindExtra_TagsContext` (17 lines) into helper sub-assertions.
11. `gitmap/cliexit/report_test.go:17`: Refactor `TestWriteStructured_HumanMode` (18 lines) into helper test verification.
12. `gitmap/cliexit/report_test.go:57`: Refactor `TestWriteStructured_JSONMode` (16 lines) into assertion sub-routine.
13. `gitmap/cliexit/report_test.go:77`: Refactor `assertJSONOutput` (16 lines) into `unmarshalReportJSON` + `verifyReportFields`.
14. `gitmap/cliexit/report_test.go:98`: Refactor `TestWriteStructured_OmitsEmpty` (18 lines) into table runner.
15. `gitmap/cliexit/report.go:168`: Refactor `writeJSON` (19 lines) into `marshalReportData` + `writeJSONStream`.
16. `gitmap/clonefrom/execute_checkout_test.go:49`: Refactor `TestBuildGitArgs_NoCheckoutOnlyForSkipMode` (23 lines) into sub-tests.
17. `gitmap/clonefrom/execute_checkout_test.go:79`: Refactor `TestExecute_SkipCheckout_NoWorkingTree` (19 lines) into helper.
18. `gitmap/clonefrom/execute_checkout_test.go:114`: Refactor `TestExecute_ForceCheckout_BranchMissingFails` (25 lines) into setup + verify.
19. `gitmap/clonefrom/execute_concurrent_test.go:77`: Refactor `TestExecuteWithHooksConcurrent_HookOrderAndCount` (23 lines).
20. `gitmap/clonefrom/execute_dest_test.go:49`: Refactor `TestExecute_MkdirParentFailureIsFailedRow` (16 lines).
21. `gitmap/clonefrom/execute_lfs_fix.go:39`: Refactor `executeLFSFix` (26 lines) into `detectLFSFiles` + `runUnTrackFix`.
22. `gitmap/clonefrom/execute_test.go:47`: Refactor `TestExecute_SkipsNonEmptyDest` (17 lines).
23. `gitmap/clonefrom/execute_test.go:96`: Refactor `TestRenderSummary_TalliesAllStatuses` (17 lines).
24. `gitmap/clonefrom/parse_test.go:17`: Refactor `TestParseFile_JSON` (21 lines).
25. `gitmap/clonefrom/parse_test.go:44`: Refactor `TestParseFile_CSV` (21 lines).
26. `gitmap/clonefrom/parse.go:33`: Refactor `ParseFile` (16 lines) into `openSourceFile` + `dispatchByExt`.
27. `gitmap/clonefrom/parse.go:101`: Refactor `jsonRow` (16 lines) into `validateRow` + `convertToScanRecord`.
28. `gitmap/clonefrom/parsecsv.go:42`: Refactor `readCSVRows` (18 lines) into `parseCSVLine` + `appendCSVRecord`.
29. `gitmap/clonefrom/parsecsv.go:85`: Refactor `indexCSVHeader` (16 lines) into `normalizeHeader` + `mapHeaderColumns`.
30. `gitmap/cloner/cloner.go:103`: Refactor `parseTextFile` (17 lines) into `scanCloneLines` + `handleScannerErr`.

## Target Verification Files
- `gitmap/archive/*`
- `gitmap/cliexit/*`
- `gitmap/clonefrom/*`
- `gitmap/cloner/*`
