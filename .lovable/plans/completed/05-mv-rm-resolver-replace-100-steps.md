# Master 100-Step Plan: GitMap MV, RM Resolver, and Replace Engine

## Phase 1: Windows Long-Path & Path Normalization Utility (Steps 1–15)

1. Create `gitmap/fsutil/longpath.go` defining cross-platform long-path interfaces.
2. Implement `gitmap/fsutil/longpath_windows.go` detecting paths > 240 chars and prepending `\\?\` or `\\?\UNC\`.
3. Implement `gitmap/fsutil/longpath_other.go` as a clean pass-through for Linux and macOS.
4. Implement `StripLongPathPrefix(path string) string` to remove `\\?\` before storing in SQLite.
5. Implement `NormalizeSlashes(path string) string` converting Windows backslashes to portable clean paths.
6. Implement `TrimTrailingSlashes(path string) string` removing all trailing `/` and `\`.
7. Implement `CanonicalPath(path string) (string, error)` resolving symlinks, casing, and relative `.` / `..` segments.
8. Implement `IsSubdirectory(parent, child string) bool` preventing accidental recursive moves into self.
9. Implement `SafeRemoveAll(path string) error` wrapping Windows long-path aware directory deletion.
10. Implement `SafeRename(src, dst string) error` attempting atomic rename and falling back to cross-volume copy.
11. Implement `CopyDirectory(src, dst string) error` for cross-device relocation with file metadata preservation.
12. Add unit tests in `gitmap/fsutil/longpath_test.go` verifying 260+ character path handling on Windows.
13. Add unit tests in `gitmap/fsutil/path_normalize_test.go` verifying slash normalization and trailing slash stripping.
14. Add unit tests in `gitmap/fsutil/is_subdir_test.go` verifying recursive nesting guard detection.
15. Verify `go test ./fsutil/...` passes with zero failures.

## Phase 2: Unified Project Resolver Engine (Steps 16–30)

16. Create `gitmap/cmd/resolver_types.go` defining `ResolveTargetOptions` and `ResolvedRepo`.
17. Create `gitmap/cmd/resolver.go` with main entrypoint `ResolveRepo(db *store.DB, target string) (*model.ScanRecord, error)`.
18. Implement PWD fallback in `gitmap/cmd/resolver_pwd.go` when target is empty or `.`.
19. Implement exact absolute path matching with case-insensitive Windows comparison in `resolver.go`.
20. Implement relative path normalization in `resolver.go` handling `.\folder`, `./folder/`, `../folder`.
21. Implement alias lookup in `gitmap/cmd/resolver_alias.go` querying SQLite `Alias` table.
22. Implement exact slug matching against `Repo` table in `resolver.go`.
23. Implement folder basename fallback matching when a user supplies only the folder name.
24. Implement glob expansion in `gitmap/cmd/resolver_glob.go` supporting `*`, `?`, and `[`.
25. Implement multi-target deduplication by `Repo.ID` in `ResolveMultiRepos(db *store.DB, targets []string)`.
26. Implement clear error formatting when zero matches or ambiguous matches are found.
27. Add unit tests in `gitmap/cmd/resolver_test.go` covering `.\prompt-architect` resolution.
28. Add unit tests in `gitmap/cmd/resolver_alias_test.go` covering short alias mapping.
29. Add unit tests in `gitmap/cmd/resolver_pwd_test.go` covering current directory detection.
30. Verify `go test ./cmd/ -run TestResolve` passes cleanly.

## Phase 3: External Integrations Lifecycle Sync (Steps 31–45)

31. Create `gitmap/vscodepm/update_path.go` with `UpdateRootPath(oldPath, newPath, newName string) error`.
32. Create `gitmap/vscodepm/remove.go` with `RemoveEntry(targetPath string) error`.
33. Implement case-insensitive normalized path comparison in `vscodepm` for Windows paths.
34. Ensure `projects.json` atomic file rewrite with backup protection (`projects.json.bak`).
35. Create `gitmap/desktop/update_path.go` registering relocated repository with GitHub Desktop CLI.
36. Create `gitmap/desktop/remove.go` implementing repo removal hook for GitHub Desktop.
37. Implement interactive prompt helper `promptExternalSync(action, repoName string, yes bool) bool`.
38. Ensure `-y` / `--yes` flag automatically bypasses prompts and executes external sync.
39. Add `--no-vscode` flag support to skip VS Code Project Manager updates.
40. Add `--no-desktop` flag support to skip GitHub Desktop updates.
41. Add unit tests in `gitmap/vscodepm/update_path_test.go` verifying rootPath updates in `projects.json`.
42. Add unit tests in `gitmap/vscodepm/remove_test.go` verifying entry deletion from `projects.json`.
43. Add unit tests in `gitmap/desktop/update_path_test.go` verifying desktop registration on move.
44. Add unit tests in `gitmap/desktop/remove_test.go` verifying desktop removal handling.
45. Verify `go test ./vscodepm/... ./desktop/...` passes with zero errors.

## Phase 4: `gitmap mv` (Move) Command Engine (Steps 46–60)

46. Create `gitmap/cmd/mv_flags.go` parsing `--yes`, `-y`, `--dry-run`, `--no-vscode`, `--no-desktop`.
47. Create `gitmap/cmd/mv.go` with `runMove(args []string)`.
48. Implement destination path resolution supporting `..` (parent directory), relative paths, and absolute paths.
49. Implement safety preflight in `gitmap/cmd/mv_preflight.go`: check source existence, destination collision, disk permissions.
50. Implement dry-run preview in `gitmap/cmd/mv_dryrun.go` printing exact from/to paths and DB modifications.
51. Implement interactive confirmation prompt in `gitmap/cmd/mv_prompt.go`.
52. Implement physical filesystem relocation in `gitmap/cmd/mv_relocate.go` using long-path safe renamer.
53. Implement SQLite database transaction in `gitmap/cmd/mv_db.go` updating `Repo.AbsolutePath` and `Repo.RepoName`.
54. Update `ScanFolderId` and `Alias` table records pointing to relocated path.
55. Dispatch VS Code Project Manager `UpdateRootPath` on move success.
56. Dispatch GitHub Desktop `UpdateRepoPath` on move success.
57. Register `mv` and `move` commands in `gitmap/cmd/rootdata.go` and `roottooling.go`.
58. Add unit tests in `gitmap/cmd/mv_test.go` covering `gitmap mv repo ..`.
59. Add unit tests in `gitmap/cmd/mv_relpath_test.go` covering relative destination moves.
60. Verify `go test ./cmd/ -run TestMove` passes.

## Phase 5: Enhanced `gitmap rm` (Remove & Delete) (Steps 61–75)

61. Refactor `gitmap/cmd/rm.go` to use the Unified Project Resolver (`ResolveMultiRepos`).
62. Ensure `gitmap rm .\prompt-architect` resolves immediately without requiring exact slug match.
63. Ensure `gitmap rm ./prompt-architect/` and trailing slash variants resolve cleanly.
64. Ensure `gitmap rm .` removes current repository from PWD.
65. Implement `--db-only` flag to untrack from SQLite database without deleting on-disk files.
66. Implement cascading database deletion across `Repo`, `GroupRepo`, `Alias`, `Bookmark`, `VersionProbe`.
67. Connect VS Code Project Manager `RemoveEntry` upon repo deletion.
68. Connect GitHub Desktop `RemoveRepo` upon repo deletion.
69. Format clean colorful deletion confirmation box with repo slug, absolute path, and sync status.
70. Honor `-y` / `--yes` flag to bypass all confirmation prompts.
71. Add unit tests in `gitmap/cmd/rm_resolver_test.go` verifying relative path inputs (`.\foo`).
72. Add unit tests in `gitmap/cmd/rm_dbonly_test.go` verifying `--db-only` preserves on-disk files.
73. Add unit tests in `gitmap/cmd/rm_cascade_test.go` verifying SQLite + VSCode cascade.
74. Add unit tests in `gitmap/cmd/rm_pwd_test.go` verifying `gitmap rm .` in repo directory.
75. Verify `go test ./cmd/ -run TestRm` passes.

## Phase 6: `gitmap replace` Diagnostics & Token Engine (Steps 76–85)

76. Inspect `gitmap/cmd/replacewalk.go` and fix forward slash vs backslash path comparisons on Windows.
77. Ensure `isExcludedPrefix` handles both `D:/path` and `D:\path` uniformly with `filepath.ToSlash`.
78. Refine `isBinaryFile` sniffer in `replacewalk.go` to avoid false positives on UTF-8 / UTF-16 code files.
79. Add fallback token matching in `replaceapply.go` for variations (`errorwrapper-v3`, `errorwrapper/v3`, `errorwrapper_v3`).
80. Implement colorized diff output in `gitmap/cmd/replaceprint.go` showing matched lines.
81. Implement atomic write verification ensuring replaced files are not truncated.
82. Add unit test in `gitmap/cmd/replace_winpath_test.go` testing replacement on Windows-style paths.
83. Add unit test in `gitmap/cmd/replace_tokens_test.go` reproducing `errorwrapper-v3` -> `errorwrapper-v4`.
84. Add unit test in `gitmap/cmd/replace_bom_test.go` testing files with UTF-8/UTF-16 encoding.
85. Verify `go test ./cmd/ -run TestReplace` passes.

## Phase 7: Help System, Web Dashboard & Terminal UI (Steps 86–95)

86. Add `HelpMV = "  mv (move) <src> <dest>     Relocate repo directory with VSCode & GitHub Desktop sync"` in `constants_helpgroups.go`.
87. Add `mv` to `CompactData` in `gitmap/constants/constants_helpgroups.go`.
88. Update `gitmap/cmd/rootusage_groups.go` and `rootusagefilter_rows.go` with `mv`.
89. Add `gitmap mv` interactive documentation in `src/data/commands.ts` with usage and examples.
90. Update `gitmap rm` documentation in `src/data/commands.ts` with path examples (`.\folder`, `..`, `.`).
91. Update `gitmap replace` documentation in `src/data/commands.ts` with literal replacement examples.
92. Add colorful terminal UI headers and summary boxes for `gitmap mv`.
93. Ensure all glyphs in `mv` and `rm` output use high-contrast UTF-8 symbols (`✔`, `✖`, `📂`, `▲`).
94. Build frontend assets with `bun run build` and verify 0 bundle errors.
95. Run `npx vitest run` and ensure all frontend documentation tests pass.

## Phase 8: Full Automated Verification & Smoke Tests (Steps 96–100)

96. Run complete Go test suite: `go test ./...` in `gitmap/`.
97. Compile `gitmap.exe` binary into `bin/`.
98. Execute end-to-end CLI tests for `gitmap mv`, `gitmap rm .\path`, and `gitmap replace`.
99. Verify zero lines > 200 and zero functions > 15 lines across all new Go code.
100. Update `walkthrough.md` with full execution logs, screenshots, and verification matrices.
