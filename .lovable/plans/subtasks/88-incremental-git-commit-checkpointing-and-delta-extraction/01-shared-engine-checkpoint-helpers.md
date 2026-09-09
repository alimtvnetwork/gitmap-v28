# Subtask 01: Shared Engine Checkpoint Helpers & Incremental Diff

## Objective
Enhance `03-ai-scripts/02-shared-engine.py` with git HEAD hash extraction, commit ancestry verification, checkpoint manifest parsing, and incremental diff range calculation.

## Requirements
1. `get_git_head_commit_hash(repo_root: Path | str) -> str`:
   - Runs `git rev-parse HEAD` in `repo_root`.
   - Returns 40-character SHA string, or empty string on failure.
2. `is_git_commit_ancestor(ancestor_sha: str, descendant_sha: str, repo_root: Path | str) -> bool`:
   - Runs `git merge-base --is-ancestor <ancestor_sha> <descendant_sha>`.
   - Returns `True` if exit code is 0, `False` otherwise.
3. `load_checkpoint_commit_hash(checkpoint_path: Path) -> str | None`:
   - Safely parses JSON file if present and returns `last_commit_hash` string.
4. `calculate_git_diff_range(repo_root: Path | str, commit_count: int, prev_hash: str | None, force_full: bool) -> tuple[list[str], str, bool]`:
   - If `prev_hash` is ancestor and not `force_full`:
     - If `prev_hash == current_head`: returns `([], f"{prev_hash[:8]}..HEAD (no new commits)", True)`
     - Else: returns `(diff_lines, f"{prev_hash[:8]}..{current_head[:8]}", True)`
   - Else: returns `(diff_lines, f"HEAD~{commit_count}..HEAD", False)`
5. `extract_git_changed_files(...)`:
   - Support `checkpoint_file: Path | str | None = None` and `force_full: bool = False`.
   - Returns `tuple[list[dict[str, Any]], dict[str, Any]]` containing deduplicated file records and checkpoint metadata (`last_commit_hash`, `previous_checkpoint_hash`, `commit_range`, `is_incremental`).

## Coding Guidelines
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested `if` statements.
- Affirmative booleans (`is_incremental`, `is_ancestor`).
