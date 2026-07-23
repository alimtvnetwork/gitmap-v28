# Technical Post-Mortems

Recent critical issues and their resolutions:

1. **Windows File Locks**: Resolved via a 'rename-first' strategy.
2. **Database Fragmentation**: Fixed by anchoring `data/` to the binary location.
3. **Path Double-Nesting**: Corrected `store` package path resolution.
4. **CLI Changelog Sync**: Synchronized release notes with compiled constants.
5. **Zip Group Silent Failure**: Added explicit error reporting.
6. **Auto-commit Push Rejection**: Implemented `git pull --rebase` recovery.
7. **List-Releases Resolution**: Prioritized local metadata.
8. **Legacy UUID Detection**: Provided recovery instructions.
9. **Auto Legacy Dir Migration**: Consolidated folders at startup.
10. **Legacy ID Migration**: Rebuilt Repos table.
11. **One-liner PATH Propagation**: Broadcasted system changes.
12. **Binary Extraction Failure**: Implemented flexible filename mapping in installer.
13. **Install Script 404**: Fixed relative paths.
14. **Latest Tag Disconnect**: Updated CI to explicitly use `make_latest` for stable releases.
15. **Octal Literal Style**: Switched to `0o644` in Go for `gocritic` compliance.
16. **Redundant Newline**: Switched to `fmt.Fprint` for constants with trailing newlines.
17. **Constant Redeclaration**: Resolved by centralizing all command IDs in `constants_cli.go`.
18. **Unchecked Errors**: Added lint-compliant error checking for `db.RemoveInstalledTool`, `dev.Process.Kill()`, and `cmd.Start()`.
19. **Release Pipeline Directory Error**: Resolved `cd: dist` failure by setting explicit `working-directory: gitmap/dist`.
20. **G305 Zip Path Traversal**: Fixed `installnpp.go` to validate extracted file paths stay within target directory.
21. **G110 Decompression Bomb**: Replaced `io.Copy` with `io.LimitReader` capped at 10 MB.
22. **Format Verb Mismatch**: Fixed `fmt.Fprintf` argument count at `tasksync.go:138`; audited ~140 call sites.
23. **Code Red Error Audit**: Standardized 35+ error constants with mandatory path, operation, and reason context.
24. **CI Passthrough Gate Pattern**: Replaced job-level `if` skipping with step-level conditionals.
25. **Go Flag Ordering — Silent Flag Drop**: Fixed with `reorderFlagsBeforeArgs()` helper.
26. **CI Release Branch Cancellation Protection**: Protected `release/**` branches from `cancel-in-progress`.
27. **Setup Config Path Resolution** (v2.74.0): `gitmap setup` failed with "file not found" when run from non-binary directories. Fixed by resolving `git-setup.json` relative to the binary's installation path via `filepath.EvalSymlinks`.
28. **Shell Wrapper Detection** (v2.74.0): `gitmap cd` silently failed when called as raw binary without shell wrapper. Added `GITMAP_WRAPPER=1` env var export in shell scripts and stderr warning in `gitmap cd` when wrapper not detected.
