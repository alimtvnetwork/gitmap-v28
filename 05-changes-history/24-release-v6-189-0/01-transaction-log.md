# Transaction Log 24: Release v6.189.0 Minor Version Bump

> **Directory:** `05-changes-history/24-release-v6-189-0/`  
> **Date:** 2026-09-06  
> **Author/Agent:** Antigravity AI  
> **Module Affected:** Repository-wide (`version.json`, `package.json`, `gitmap/constants/constants.go`, `changelog.md`, `readme.md`, `.lovable/`, `05-changes-history/`)  
> **Status:** Completed & Verified  

---

## 1. Context & User Directives

The user explicitly commanded a minor release:
```text
release minor
```

In accordance with the `release-and-versioning` skill:
1. SSoT manifests (`version.json`, `package.json`, `gitmap/constants/constants.go`) were bumped from `6.188.0` to `6.189.0`.
2. `changelog.md` was updated with the new `## [v6.189.0]` entry, install one-liners, and itemized release bullet points.
3. Release notes were generated at `.lovable/release/release-notes-v6.189.0.md`.
4. Root `readme.md` was pinned to active version `v6.189.0`.
5. `.lovable/user-preferences` was updated to reflect active version `v6.189.0`.
6. Version synchronization was validated via `03-ai-scripts/14-version-sync-checker.py --all-paths` (3/3 passing).

---

## 2. Release Highlights for v6.189.0

- **Single-Map Lazy Regex Architecture (`gitmap/lazyregex` & `pkg/regexnew`)**:
  - Eliminated redundant secondary `regexMap` lookups, transitioning to single-map pattern deduplication.
  - Added `isCompiled` boolean flag and self-contained compiled `*regexp.Regexp` state protected by internal mutex.
  - Compilation occurs lazily on demand and caches directly inside the lazy regex instance, removing unnecessary secondary map overhead.
- **Pattern Query & Inspection Methods**:
  - `Count(s string) int`: Counts non-overlapping occurrences of the pattern in the target string.
  - `IsFound(s string) bool`: Fast existence check for pattern occurrences.
  - `GroupBy(s string) map[string]string`: Extracts named submatches into key-value map pairs.
  - `FindAllGroups(s string) []map[string]string`: Extracts all named submatches across matches.
- **Structured Error Handling & Builder Integration**:
  - Added `CompileAppError()` returning a structured `*appfault.AppError` on invalid regex syntax.
  - Added `CompileBuilder()` returning a fluent diagnostic `AppBuilder` containing regex pattern and compilation metadata.
- **Recursive Deadlock Resolution (`pkg/regexnew`)**:
  - Resolved deadlock during default compiler initialization by performing direct compilation within `Compile()` rather than recursive compiler calls.
  - Concurrent test suite executes with 100 parallel goroutines and 0 deadlocks.

---

## 3. Files Modified & Created

### Modified
1. `version.json` — Bumped `Version` from `6.188.0` to `6.189.0`.
2. `package.json` — Bumped `version` from `6.188.0` to `6.189.0`.
3. `gitmap/constants/constants.go` — Bumped `Version` from `6.188.0` to `6.189.0`.
4. `readme.md` — Updated pinned version to `v6.189.0`.
5. `.lovable/user-preferences` — Pinned active version `v6.189.0`.
6. `changelog.md` — Added release header, install commands, and release notes for `v6.189.0`.
7. `05-changes-history/01-index.md` — Registered transaction log 24.

### Created
1. `05-changes-history/24-release-v6-189-0/01-transaction-log.md` — This transaction log.
2. `.lovable/release/release-notes-v6.189.0.md` — Generated release notes.

---

## 4. Verification & Quality Gates

- `python 03-ai-scripts/14-version-sync-checker.py --all-paths`: 3/3 checks passed (package.json, constants.go, changelog.md).
- `go test -v -count=1 ./pkg/regexnew/...` in `04-code/golang`: 100% PASS.
- `go test -v -count=1 ./lazyregex/... ./pipelinedb/... ./constants/...` in `gitmap`: 100% PASS.
- `python linter-scripts/check-nested-ifs.py`: 0 violations across new and modified files.
