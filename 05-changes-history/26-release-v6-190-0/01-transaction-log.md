# Transaction Log 26: Release v6.190.0 Minor Version Bump

> **Directory:** `05-changes-history/26-release-v6-190-0/`  
> **Date:** 2026-09-06  
> **Author/Agent:** Antigravity AI  
> **Module Affected:** Repository-wide (`version.json`, `package.json`, `gitmap/constants/constants.go`, `changelog.md`, `readme.md`, `.lovable/`, `05-changes-history/`)  
> **Status:** Completed & Verified  

---

## 1. Context & User Directives

The user explicitly commanded a minor release via the release script:
```text
release using release script please
```

In accordance with the `release-and-versioning` skill:
1. Executed canonical release bumper `03-ai-scripts/29-release-bumper.py --bump minor` with itemized release highlights.
2. SSoT manifests (`version.json`, `package.json`, `gitmap/constants/constants.go`) were bumped from `6.189.0` to `6.190.0`.
3. `changelog.md` was updated with the new `## [v6.190.0]` entry, install one-liners, and itemized release bullet points.
4. Release notes were generated at `.lovable/release/release-notes-v6.190.0.md`.
5. Root `readme.md` was pinned to active version `v6.190.0`.
6. `.lovable/user-preferences` was updated to reflect active version `v6.190.0`.
7. Version synchronization was validated via `03-ai-scripts/14-version-sync-checker.py --all-paths` (3/3 passing).

---

## 2. Release Highlights for v6.190.0

- **Strict Thread-Safe Lazy Compile Locking**:
  - Regex creation never compiles prematurely; compilation occurs lazily on demand.
  - Under internal mutex lock, checks `isCompiled`, returns cached pointer if true, or compiles once with `regexp.Compile(pattern)` and caches on success.
  - Applied across `gitmap/lazyregex` and `04-code/golang/pkg/regexnew`.
- **CompileResult Envelope**:
  - Introduced unified `CompileResult` struct replacing tuple returns in `CompileBuilder()` and `CompileAppError()`.
  - Encapsulates compiled `*regexp.Regexp`, structured `*appfault.AppError` / `*apperror.AppError`, and `*appfault.AppBuilder` with query helpers (`IsSuccess`, `IsFailed`, `HasError`, `Regexp`, `AppError`, `Builder`, `Unwrap`).
- **Dedicated `GroupMap` Data Type**:
  - Replaced raw `map[string]string` from `GroupBy(s)` and `FindGroups(s)`.
  - Fluent, nil-safe operations: `Has`, `HasKey`, `Get`, `GetOrDefault`, `Set`, `Add`, `Remove`, `Delete`, `Keys`, `AllKeys`, `Values`, `AllValues`, `Len`, `Count`, `IsEmpty`, `HasItems`, `Clone`, `Clear`, `ToMap`, `Raw`, `ForEach`, `Filter`, `MarshalJSON`.
- **Dedicated `GroupList` Data Type**:
  - Replaced raw `[]map[string]string` from `FindAllGroups(s)`.
  - Fluent, nil-safe operations: `Items`, `Count`, `Len`, `IsEmpty`, `HasItems`, `First`, `Last`, `At`, `Add`, `AllKeys`, `Keys`, `ValuesOf`, `Find`, `Filter`, `ForEach`, `Clone`, `ToMaps`, `Raw`, `MarshalJSON`.

---

## 3. Files Modified & Created

### Modified
1. `version.json` — Bumped `Version` from `6.189.0` to `6.190.0`.
2. `package.json` — Bumped `version` from `6.189.0` to `6.190.0`.
3. `gitmap/constants/constants.go` — Bumped `Version` from `6.189.0` to `6.190.0`.
4. `readme.md` — Updated pinned version to `v6.190.0`.
5. `.lovable/user-preferences` — Pinned active version `v6.190.0`.
6. `changelog.md` — Added release header, install commands, and release notes for `v6.190.0`.
7. `05-changes-history/01-index.md` — Registered transaction log 26.

### Created
1. `05-changes-history/26-release-v6-190-0/01-transaction-log.md` — This transaction log.
2. `.lovable/release/release-notes-v6.190.0.md` — Generated release notes.

---

## 4. Verification & Quality Gates

- `python 03-ai-scripts/14-version-sync-checker.py --all-paths`: 3/3 checks passed.
- `go test -v -count=1 ./pkg/regexnew/...` in `04-code/golang`: 12/12 PASS (0.191s).
- `go test -v -count=1 ./lazyregex/... ./pipelinedb/... ./constants/...` in `gitmap`: 100% PASS.
- `python linter-scripts/check-nested-ifs.py`: 0 violations across new and modified files.
- `python linter-scripts/check-enum-guidelines.py`: PASS.
