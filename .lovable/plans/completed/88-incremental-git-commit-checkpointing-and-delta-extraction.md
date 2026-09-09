# Plan 88: Incremental Git Commit Checkpointing & Delta Extraction Engine

## Overview
Implement persistent Git commit checkpointing in `03-ai-scripts/27-git-changed-files.py` and `03-ai-scripts/02-shared-engine.py`. When extracting changed files, record the current `HEAD` commit hash into the manifest (`.lovable/temp/git-changed-files.json`). On subsequent runs, verify the saved commit hash, compute the exact incremental diff (`<last_commit_hash>..HEAD` plus working-tree changes), and only process files changed since the previous checkpoint. This dramatically reduces subsequent linter scan times and CPU/IO overhead.

---

## Key Requirements & Scope
1. **Checkpointing & Incremental Diff Engine (`03-ai-scripts/02-shared-engine.py`)**:
   - `get_git_head_commit_hash(repo_root)`: Fast resolution of current `HEAD` 40-character SHA.
   - `is_git_commit_ancestor(ancestor_sha, descendant_sha, repo_root)`: Verification that previous checkpoint is a valid ancestor of current `HEAD`.
   - `extract_git_changed_files(...)`:
     - Accept optional `checkpoint_file` and `force_full` parameters.
     - If checkpoint exists and previous commit SHA is an ancestor of `HEAD`:
       - If `prev_sha == current_head`: diff range is empty (0 new committed files), only capture working-tree uncommitted changes.
       - If `prev_sha != current_head`: diff range is `f"{prev_sha}..HEAD"` plus working-tree uncommitted changes.
     - Fall back to full window `f"HEAD~{commit_count}..HEAD"` if checkpoint is absent, invalid, or ancestor check fails.
     - Return deduplicated files list along with checkpoint metadata.
   - `load_git_changed_files_manifest(manifest_path)`: Safe loader handling both new structured schema and legacy list formats.

2. **CLI Checkpointing & Structured Manifests (`03-ai-scripts/27-git-changed-files.py`)**:
   - Save structured manifest with `last_commit_hash`, `previous_checkpoint_hash`, `commit_range`, `is_incremental`, `total_files`, and `files`.
   - Output across JSON, YAML, and TXT with checkpoint metadata preserved.
   - Add CLI options: `--no-incremental` / `--full-window`, `--checkpoint <path>`, `--since-commit <sha>`.
   - Enhanced summary banner displaying previous checkpoint SHA, current `HEAD` SHA, extraction mode (Incremental vs Full Window), and deduplicated file count.

3. **Linter & Runner Adaptations**:
   - `linter-scripts/check-relative-paths.py`: Read structured manifest backward-compatibly; display checkpoint metadata.
   - `linter-scripts/check-nested-ifs.py`: Read structured manifest backward-compatibly.
   - `linter-scripts/check-enum-and-boolean.py`: Read structured manifest backward-compatibly.
   - `03-ai-scripts/06-cicd-local-runner.py`: Forward incremental/full flags and ensure memory-safe gate concurrency.

4. **Task-Specific Constraints**:
   - Strict functions <= 15 lines.
   - Blank line before every return statement.
   - Zero nested `if` statements.
   - Affirmative booleans (`is_incremental`, `is_ancestor`, `is_fresh`).
   - Zero external library dependencies (pure standard library).

---

## Subtasks Breakdown
- `01-shared-engine-checkpoint-helpers.md`: Add HEAD hash, ancestor verification, and incremental diff logic to `02-shared-engine.py`. [COMPLETED]
- `02-git-changed-files-cli-and-manifest.md`: Update `27-git-changed-files.py` with structured schema, multi-format serializers, and CLI options. [COMPLETED]
- `03-linters-and-runner-integration.md`: Update `check-relative-paths.py`, `check-nested-ifs.py`, `check-enum-and-boolean.py`, and `06-cicd-local-runner.py` for incremental checkpoint consumption. [COMPLETED]
- `04-verification-and-ci-gates.md`: Comprehensive test suite verifying initial extraction, subsequent delta extraction, fallback on ancestor mismatch, and full CI gates green. [COMPLETED]

---

## Success Criteria
- [x] Initial extraction writes `last_commit_hash` to `.lovable/temp/git-changed-files.json`.
- [x] Second run with no commits processes 0 committed files (or only working-tree changes) in < 0.05s.
- [x] New commit triggers incremental diff `<last_commit_hash>..HEAD` with only newly changed files.
- [x] `--no-incremental` forces full $N$-commit window extraction.
- [x] All linters exit 0 with 0 violations.
- [x] Local CI runner exits code 0 (36/36 gates green).
