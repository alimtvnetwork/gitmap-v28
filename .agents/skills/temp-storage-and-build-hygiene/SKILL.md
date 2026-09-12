---
name: temp-storage-and-build-hygiene
description: Enforce repository-scoped temporary directory isolation and mandatory pre-build cleanup across all CI/CD, Go, and Python workflows.
---

# Temp Storage & Build Hygiene

## Core Directives

1. **Repository Namespacing in OS/User Temp**:
   - Never write loose files or arbitrary directories directly into the root of the OS/user temp directory (`os.TempDir()`, `tempfile.gettempdir()`, `$env:TEMP`, `$TMPDIR`).
   - All temporary files and sandboxes must reside within a repository-scoped folder: `<temp_dir>/gitmap/<subfolder>/`.
   - Use dedicated functional subdirectories:
     - `<temp_dir>/gitmap/build/` — Compiled test binaries, CLI executables.
     - `<temp_dir>/gitmap/test/` — Smoke test sandboxes, worker runtimes.
     - `<temp_dir>/gitmap/purge/` — History purge backup archives.
     - `<temp_dir>/gitmap/downloads/` — External archive downloads and installers.

2. **Mandatory Pre-Build Cleanup (Storage Respect & Reuse)**:
   - Before executing any compilation or build command (`go build`, `npm run build`, `goreleaser`), always purge previous build artifacts in the target build directory.
   - Never allow stale build outputs or multiple versions to accumulate and waste disk storage.
   - Re-use storage paths cleanly with explicit wipe-before-write semantics.

3. **In-Repository Temp Bounding**:
   - For repository-internal temporary files, isolate strictly inside `.lovable/temp/` (e.g. `.lovable/temp/failures/`, `.lovable/temp/cicd/`).
   - Root `.tmp/` creation is strictly banned.
