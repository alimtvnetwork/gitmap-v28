# Transaction Log 21: Release v6.187.0 Minor Version Bump

> **Directory:** `05-changes-history/21-release-v6-187-0/`  
> **Date:** 2026-09-05  
> **Author/Agent:** Antigravity AI  
> **Module Affected:** Repository-wide (`version.json`, `package.json`, `gitmap/constants/constants.go`, `changelog.md`, `.lovable/`)  
> **Status:** Completed & Verified  

---

## 1. Context & User Directives

The user requested a minor release:
```text
release minor
```

In accordance with the `release-and-versioning` skill:
1. SSoT manifests (`version.json`, `package.json`, `gitmap/constants/constants.go`) were bumped from `6.186.0` to `6.187.0`.
2. `changelog.md` was updated with the new `## [v6.187.0]` entry, install one-liners, and itemized release bullet points.
3. Release notes were generated at `.lovable/release/release-notes-v6.187.0.md`.
4. `.lovable/user-preferences` was updated to reflect active version `v6.187.0`.
5. Version synchronization was validated via `03-ai-scripts/14-version-sync-checker.py`.
6. Complete 16-gate local CI/CD test runner verified 100% green pass.

---

## 2. Release Highlights for v6.187.0

- **Type-Safe File Utilities (`pkg/fileutil`)**:
  - `FileActionType` and `FileModeType` enums adhering to repository guidelines.
  - `FileWrapper` struct with constructors and file I/O operations.
- **Enhanced `appfault` Error Object**:
  - First-class `errorId` field on `AppError`.
  - Dedicated wrap constructors: `WrapFailure`, `WrapWriterFailure`, `WrapReaderFailure`, `NewWithId`.
- **Local CI/CD Runner Overhaul (`03-ai-scripts/06-cicd-local-runner.py`)**:
  - Parallel worker group execution based on CPU core scaling.
  - 3-batch I/O throttling protecting heavy builds and SQLite database single-writer access.
  - Fixed SQLite race condition by prioritizing `SQLPragmaBusyTimeout5s` before `SQLPragmaJournalWAL`.
  - Quiet mode on success: `✔ All passed. (16 gates in 34.61s)`.
  - Full CLI options: `--all-paths`, `--sync`, `--output`, `--json`.
- **Reusable Worker Pool Base Engine (`03-ai-scripts/02-shared-engine.py`)**:
  - Generic `run_worker_pool` and `add_worker_cli_arguments` engine.
  - Modernized `16-installer-smoke-tester.py` and `28-go-preflight-ci.py`.
- **Architectural Research Documents**:
  - `research/12-fileutil-wrapper-and-filemode-enums.md`
  - `research/13-enhanced-cicd-local-runner-cli.md`
  - `research/14-generic-worker-pool-engine-and-installer-tester.md`

---

## 3. Files Modified & Created

### Modified
1. `version.json` - Bumped `Version` from `6.186.0` to `6.187.0`.
2. `package.json` - Bumped `version` from `6.186.0` to `6.187.0`.
3. `gitmap/constants/constants.go` - Bumped `Version` from `6.186.0` to `6.187.0`.
4. `.lovable/user-preferences` - Pinned active version `v6.187.0`.
5. `changelog.md` - Added release header, quick install commands, and release notes for `v6.187.0`.
6. `05-changes-history/01-index.md` - Registered transaction log 21.

### Created
1. `.lovable/release/release-notes-v6.187.0.md` - Standalone release notes for `v6.187.0`.
2. `05-changes-history/21-release-v6-187-0/01-transaction-log.md` - This transaction log.

---

## 4. Verification Results

1. **Version Parity Check**:
   Executed `python 03-ai-scripts/14-version-sync-checker.py`:
   - Output: `✅ Version synchronization verified: v6.187.0 (87.26ms)`.

2. **Binary Build Verification**:
   Executed `go build -C gitmap -o ../bin/gitmap.exe .` and `bin/gitmap.exe version`:
   - Output: `gitmap v6.187.0`.

3. **Spell Check**:
   Executed `python .github/scripts/misspell-changed.py`:
   - Output: `✅ Spell check passed on 7 file(s)`.

4. **Full 16-Gate CI/CD Suite**:
   Executed `python 03-ai-scripts/06-cicd-local-runner.py`:
   - Output: `✔ All passed. (16 gates in 34.61s)`.
