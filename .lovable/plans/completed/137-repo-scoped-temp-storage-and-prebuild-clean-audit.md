# 137-repo-scoped-temp-storage-and-prebuild-clean-audit.md: Repository-Scoped Temp Storage & Mandatory Pre-Build Cleanup Audit

**Status: completed**
**Execution Loop Count: 8 steps completed in continuous N-step self-loop (N=150 budget)**
**Task Initiation: Initiated via user request to confirm build locations (distinguishing workspace bin from OS/user temp), audit and strictly enforce repository-scoped subdirectories under OS temp (<temp>/gitmap/<category>), mandate clean-before-build execution across CI/CD runners and scripts to prevent disk waste, and update all prompts, skills, and coding guidelines.**

---

## 1. Executive Summary & Build Locations Confirmation

### Where Things Are Built:
1. **Application Executable (Workspace Build)**:
   - Primary CLI binary compiles to `bin/gitmap.exe` directly inside the repository root.
   - Installed globally for CLI usage at `C:\Users\<user>\AppData\Local\gitmap-cli\gitmap.exe`.
   - Web application assets compile to `dist/` or `build/` within the repository root.
2. **OS/User Temp Artifacts (`os.TempDir()` / `tempfile.gettempdir()`)**:
   - Strictly namespaced inside a repository-named root folder: `<temp_dir>/gitmap/<category>/`.
   - Functional categorization:
     - `<temp_dir>/gitmap/build/` — Go compiler temporary scratch files (`GOTMPDIR`), test builds.
     - `<temp_dir>/gitmap/test/` — Smoke test runtimes, Python `tempfile.TemporaryDirectory()` sandboxes (`TMPDIR`, `TEMP`, `TMP`).
     - `<temp_dir>/gitmap/purge/` — Git history purge backup archives.
     - `<temp_dir>/gitmap/downloads/` — Installer downloads and staging packages.
     - `<temp_dir>/gitmap/handoff/` — Windows self-uninstall handoff binary.
     - `<temp_dir>/gitmap/sandbox/` — History rewrite mirror clone sandboxes.
3. **Repository-Internal Temp Artifacts**:
   - Stored strictly under `.lovable/temp/` (`.lovable/temp/failures/`, `.lovable/temp/runner-eta.json`). Never in a root `.tmp/`.

### Storage Respect & Mandatory Pre-Build Cleanup:
- Prior to any compilation step in `03-ai-scripts/06-cicd-local-runner.py` or `.github/scripts/smoke-installer.py`, `clear_repo_build_temp()` is executed:
  - Deletes destination binary (`bin/gitmap.exe`).
  - Wipes `<temp>/gitmap/build/`.
  - Re-creates a clean, empty build directory.
- Guarantees zero orphaned binaries, avoids accumulating stale build bloat, and enforces clean storage reuse on every single build pass.

---

## 2. Changes Made Across Codebase, Scripts, Prompts & Skills

1. **`03-ai-scripts/06-cicd-local-runner.py`**:
   - Added `get_repo_os_temp_dir(*subdirs: str) -> Path`.
   - Added `clear_repo_build_temp() -> None` with automatic wipe before `Go Compile Gate`, `Web App Build`, `GoReleaser Snapshot Build`, and runner context initialization.
   - Directed `GOTMPDIR` to `<os_temp>/gitmap/build`.
   - Directed `TMPDIR`, `TEMP`, and `TMP` to `<os_temp>/gitmap/test`.
2. **`cli/cmd/historyrewrite_sandbox.go`**:
   - Updated `mirrorClone` to create sandboxes under `tempdir.RepoTempDir("sandbox")`.
3. **`cli/cmd/selfuninstallhandoff.go`**:
   - Updated `writeHandoffCopy` to write handoff binaries under `tempdir.RepoTempDir("handoff")`.
4. **Skills (`.agents/skills/`)**:
   - Enhanced `.agents/skills/temp-storage-and-build-hygiene/SKILL.md` with Go and Python implementation patterns and storage reuse invariants.
   - Updated `.agents/skills/execute-parent-task/skill.md` with Storage Hygiene & Pre-Build Cleanup checklist item.
   - Updated `.agents/skills/execute-parent-task-with-n-steps/skill.md` with Storage Hygiene & Pre-Build Cleanup checklist item.
5. **Prompts (`01-prompts/`)**:
   - Updated `01-prompts/14-execute/02-execute-parent-task-with-n-steps.md` with OS temp scoping and pre-build cleanup rules.
   - Updated `01-prompts/14-execute/01-execute-pending-tasks.md` with artifact sanitizer and storage hygiene directives.

---

## 3. Verification & Quality Gates Passed

- `golangci-lint run ./...`: 0 issues
- `go vet ./...`: 0 issues
- `python linter-scripts/check-nested-ifs.py`: 0 violations across 2,859 files
- `python linter-scripts/check-boolean-guidelines.py`: 0 violations across 2,859 files
- `python linter-scripts/check-relative-paths.py`: 0 violations across 6,815 files
- `python -m py_compile 03-ai-scripts/06-cicd-local-runner.py`: Clean
- Binary build: `bin/gitmap.exe` compiles cleanly and is synced to global CLI path.
