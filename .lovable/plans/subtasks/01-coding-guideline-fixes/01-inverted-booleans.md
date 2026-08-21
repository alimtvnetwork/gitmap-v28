# Subtask 01: Inverted Booleans & Inverse Naming Refactor

Slug: 01-inverted-booleans
Parent Plan: 01-coding-guideline-fixes
Status: pending

## Objective
Eliminate all inverted boolean expressions (`!ok`, `!strings.Contains`, `!is*`, `!has*`) and enforce inverse naming (`isHonest`/`isDishonest`, `isFound`/`isMissing`, `isSuccess`/`isFailed`).

## Concrete Execution Steps (30 Steps)

1. `gitmap/clonefrom/execute_checkout_test.go:150`: Replace `if !strings.HasPrefix(detail, wantPrefix)` with `isPrefixMatched := strings.HasPrefix(detail, wantPrefix); if isPrefixMatched == false`
2. `gitmap/clonefrom/execute_checkout_test.go:166`: Replace `if !ok` with `isSuccess := ok; if isSuccess == false`
3. `gitmap/clonefrom/execute_checkout_test.go:183`: Replace `if !ok || len(detail) != 0` with `isFailed := !ok || len(detail) > 0; if isFailed == true`
4. `gitmap/clonefrom/execute_checkout_test.go:196`: Replace `if !strings.Contains(err.Error(), "bogus")` with `isErrorMatched := strings.Contains(err.Error(), "bogus"); if isErrorMatched == false`
5. `gitmap/clonefrom/execute_concurrent_test.go:112`: Replace `if !bytes.Contains(got, []byte(r.URL))` with `isFound := bytes.Contains(got, []byte(r.URL)); if isFound == false`
6. `gitmap/clonefrom/execute_dest_test.go:67`: Replace `if !strings.Contains(results[0].Detail, "mkdir parent")` with `isMatched := strings.Contains(results[0].Detail, "mkdir parent"); if isMatched == false`
7. `gitmap/clonefrom/execute_test.go:108`: Replace `if !strings.Contains(out, "2 ok, 1 skipped...")` with `isMatched := strings.Contains(out, ...); if isMatched == false`
8. `gitmap/clonefrom/execute_test.go:111`: Replace `if !strings.Contains(out, "report: /tmp/r.csv")` with `isReportPresent := strings.Contains(out, ...); if isReportPresent == false`
9. `gitmap/clonefrom/parse_checkout_test.go:44`: Replace `if !strings.Contains(err.Error(), "yolo")` with `isExpectedError := strings.Contains(err.Error(), "yolo"); if isExpectedError == false`
10. `gitmap/clonefrom/parse_test.go:77`: Replace `if !strings.Contains(err.Error(), "row 2")` with `isExpectedRowError := strings.Contains(err.Error(), "row 2"); if isExpectedRowError == false`
11. `gitmap/clonefrom/parse_test.go:110`: Replace `if !strings.Contains(err.Error(), "url")` with `isUrlError := strings.Contains(err.Error(), "url"); if isUrlError == false`
12. `gitmap/clonefrom/parsecsv_columnerr_test.go:81`: Replace `if !strings.Contains(msg, want)` with `isMsgMatched := strings.Contains(msg, want); if isMsgMatched == false`
13. `gitmap/clonefrom/render.go:168`: Extract URL prefix normalization to helper `isBareDomain := !strings.Contains(...)`
14. `gitmap/clonefrom/render_test.go:36`: Replace `if !strings.Contains(out, s)` with `isSnippetFound := strings.Contains(out, s); if isSnippetFound == false`
15. `gitmap/clonefrom/render_test.go:55`: Replace `if !strings.Contains(out, "branch: dev")` with `isBranchFound := strings.Contains(out, "branch: dev"); if isBranchFound == false`
16. `gitmap/clonefrom/render_test.go:58`: Replace `if !strings.Contains(out, "depth:  5")` with `isDepthFound := strings.Contains(out, "depth:  5"); if isDepthFound == false`
17. `gitmap/clonefrom/result_schema_drift_test.go:164`: Replace `if !have[s]` with `isFieldPresent := have[s]; if isFieldPresent == false`
18. `gitmap/clonefrom/summary_golden_test.go:129`: Replace `if !bytes.Equal(gotNorm, wantNorm)` with `isEqual := bytes.Equal(gotNorm, wantNorm); if isEqual == false`
19. `gitmap/clonefrom/summary_provenance_test.go:29`: Replace `if !reflect.DeepEqual(want, got)` with `isMatch := reflect.DeepEqual(want, got); if isMatch == false`
20. `gitmap/clonefrom/summary_provenance_test.go:45`: Replace `if !allowed[p.Stage]` with `isAllowed := allowed[p.Stage]; if isAllowed == false`
21. `gitmap/clonefrom/summary_terminal_test.go:61`: Replace `if !strings.Contains(out, s)` with `isMatch := strings.Contains(out, s); if isMatch == false`
22. `gitmap/clonefrom/summary_terminal_test.go:79`: Replace `if !strings.Contains(buf.String(), "(skipped")` with `isSkippedPresent := strings.Contains(...); if isSkippedPresent == false`
23. `gitmap/clonefrom/validate.go:41`: Replace `if !looksLikeGitURL(r.URL)` with `isValidURL := looksLikeGitURL(r.URL); if isValidURL == false`
24. `gitmap/clonefrom/validate.go:50`: Replace `if !isValidCheckout(r.Checkout)` with `isValid := isValidCheckout(r.Checkout); if isValid == false`
25. `gitmap/clonenext/batch.go:146`: Replace `if !looksLikeHeader(header)` with `isHeader := looksLikeHeader(header); if isHeader == false`
26. `gitmap/clonenext/batch.go:195`: Replace `if !entry.IsDir()` with `isDir := entry.IsDir(); if isDir == false`
27. `gitmap/clonenext/github.go:126`: Replace `if !strings.Contains(url, "@")` with `isSSH := strings.Contains(url, "@"); if isSSH == false`
28. `gitmap/clonenext/remoteupdate.go:55`: Replace `if !parsed.HasVersion` with `isVersioned := parsed.HasVersion; if isVersioned == false`
29. `src/components/docs/CodeBlock.tsx:125`: Replace `if (!hasLanguage)` with `const isMissingLanguage = !hasLanguage; if (isMissingLanguage) return null;`
30. `src/components/docs/CodeBlock.tsx:151`: Replace `if (!hasPinned)` with `const isUnpinned = !hasPinned; if (isUnpinned) return code;`

## Target Verification Files
- `gitmap/clonefrom/*_test.go`
- `gitmap/clonenext/*_test.go`
- `src/components/docs/CodeBlock.tsx`
