# Subtask: Inverted Booleans (Go)
Status: pending

## Steps
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
