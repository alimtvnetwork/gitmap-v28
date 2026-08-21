# Subtask 02: Return Line-Gap Enforcement

Slug: 02-return-line-gaps
Parent Plan: 01-coding-guideline-fixes
Status: pending

## Objective
Enforce the mandatory blank line before every `return` statement across Go and TypeScript source files (exempting single-line `if` guards).

## Concrete Execution Steps (30 Steps)

1. `gitmap/archive/archive.go:168`: Add blank line before `return FormatUnknown`
2. `gitmap/archive/archive.go:182`: Add blank line before `return FormatUnknown`
3. `gitmap/archive/extract.go:61`: Add blank line before `return format, nil`
4. `gitmap/archive/extract.go:82`: Add blank line before `return res, nil`
5. `gitmap/archive/extract.go:119`: Add blank line before `return written, nil`
6. `gitmap/archive/extract.go:137`: Add blank line before `return writeArchiveFile(entry, clean, written)`
7. `gitmap/archive/extract.go:146`: Add blank line before `return nil`
8. `gitmap/archive/extract.go:169`: Add blank line before `return nil`
9. `gitmap/archive/extract.go:213`: Add blank line before `return flattened, moveEntries(root, finalDir, entries)`
10. `gitmap/archive/extract.go:236`: Add blank line before `return root, flattened, nil`
11. `gitmap/archive/extract.go:247`: Add blank line before `return nil`
12. `gitmap/archive/extract.go:275`: Add blank line before `return copyDirEntry(src, dst, path, d)`
13. `gitmap/archive/extract.go:308`: Add blank line before `return streamToFile(in, dst, mode)`
14. `gitmap/archive/extract.go:319`: Add blank line before `return err`
15. `gitmap/archive/list.go:56`: Add blank line before `return nil`
16. `gitmap/archive/list.go:62`: Add blank line before `return out, mholtToFormat(format), err`
17. `gitmap/clonefrom/execute_checkout.go:59`: Add blank line before `return "", true`
18. `gitmap/clonefrom/execute_dest.go:44`: Add blank line before `return fmt.Sprintf(...), false`
19. `gitmap/clonefrom/execute.go:75`: Add blank line before `return runRowLifecycle(...)`
20. `gitmap/clonefrom/execute.go:89`: Add blank line before `return Result{...}`
21. `gitmap/clonefrom/execute.go:99`: Add blank line before `return runGitClone(r, dest, cwd)`
22. `gitmap/clonefrom/execute.go:121`: Add blank line before `return handleGitCloneError(...)`
23. `gitmap/clonefrom/execute.go:128`: Add blank line before `return cmd.CombinedOutput()`
24. `gitmap/clonefrom/execute.go:136`: Add blank line before `return trimGitError(outputStr, err), false`
25. `gitmap/clonefrom/execute.go:146`: Add blank line before `return applyLfsFix(...)`
26. `gitmap/cloner/cloner.go:63`: Add blank line before `return cloneAll(...)`
27. `gitmap/cloner/cloner.go:85`: Add blank line before `return records, nil`
28. `gitmap/cloner/cloner.go:118`: Add blank line before `return records, nil`
29. `gitmap/cloner/cloner.go:133`: Add blank line before `return rec`
30. `gitmap/store/bookmark.go:71`: Add blank line before `return results, nil`

## Target Verification Files
- `gitmap/archive/*_test.go`
- `gitmap/clonefrom/*_test.go`
- `gitmap/cloner/*_test.go`
- `gitmap/store/*_test.go`
