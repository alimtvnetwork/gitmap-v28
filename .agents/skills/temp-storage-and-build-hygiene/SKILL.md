---
name: temp-storage-and-build-hygiene
description: Enforce repository-scoped temporary directory isolation and mandatory pre-build cleanup across all CI/CD, Go, and Python workflows.
---

# Temp Storage & Build Hygiene

## Core Directives

1. **Repository Namespacing in OS/User Temp**:
   - Never write loose files or arbitrary directories directly into the root of the OS/user temp directory (`os.TempDir()`, `tempfile.gettempdir()`, `$env:TEMP`, `$TMPDIR`).
   - All temporary files and sandboxes must reside within a repository-scoped folder: `<temp_dir>/gitmap/<category>/`.
   - Dedicated functional subdirectories:
     - `<temp_dir>/gitmap/build/` — Compiled test binaries, CLI executables (`GOTMPDIR`).
     - `<temp_dir>/gitmap/test/` — Smoke test sandboxes, worker runtimes (`TMPDIR`, `TEMP`, `TMP`).
     - `<temp_dir>/gitmap/purge/` — History purge backup archives.
     - `<temp_dir>/gitmap/downloads/` — External archive downloads and installers.
     - `<temp_dir>/gitmap/handoff/` — Windows self-uninstall handoff binary.
     - `<temp_dir>/gitmap/sandbox/` — History rewrite mirror clone sandbox.

2. **Mandatory Pre-Build Cleanup (Storage Respect & Reuse)**:
   - Before executing any compilation or build command (`go build`, `npm run build`, `goreleaser`, `06-cicd-local-runner.py`), always purge previous build artifacts in the target build directory (`bin/gitmap.exe` and `<temp>/gitmap/build/`).
   - Never allow stale build outputs or multiple binary versions to accumulate and waste disk storage. Storage must be reused efficiently every single build run.
   - Re-use storage paths cleanly with explicit wipe-before-write semantics.

3. **In-Repository Temp Bounding**:
   - For repository-internal temporary files, isolate strictly inside `.lovable/temp/` (e.g. `.lovable/temp/failures/`, `.lovable/temp/runner-eta.json`).
   - Root `.tmp/` creation is strictly banned.

## Implementation Standard

### Go Package (`cli/tempdir/tempdir.go`)
- Use `tempdir.RepoTempDir(subdirs ...string)` for all OS temp paths.
- Use `tempdir.BuildTempDir()` for build targets.
- Use `tempdir.TestTempDir()` for test directories.
- Always call `tempdir.ClearRepoBuildTempDir()` prior to building.

### Python Scripts (`03-ai-scripts/06-cicd-local-runner.py`, `.github/scripts/`)
- Define `get_repo_os_temp_dir(*subdirs: str) -> Path` targeting `Path(tempfile.gettempdir()) / "gitmap" / ...`.
- Define `clear_repo_build_temp() -> None` wiping destination binaries and `<temp>/gitmap/build/` before compilation.
- Direct `GOTMPDIR` to `<temp>/gitmap/build` and `TMPDIR`/`TEMP`/`TMP` to `<temp>/gitmap/test`.
