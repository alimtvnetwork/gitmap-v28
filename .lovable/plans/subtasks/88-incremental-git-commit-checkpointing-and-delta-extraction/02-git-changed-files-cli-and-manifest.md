# Subtask 02: Git Changed Files CLI & Multi-Format Checkpoint Manifest

## Objective
Update `03-ai-scripts/27-git-changed-files.py` to persist `last_commit_hash` and checkpoint metadata in JSON, YAML, and TXT manifests, and provide CLI controls for incremental vs full window extraction.

## Requirements
1. Update `serialize_json(...)`:
   - Include top-level fields:
     - `last_commit_hash`: current HEAD SHA.
     - `previous_checkpoint_hash`: previous SHA or None.
     - `commit_range`: string (e.g. `a1b2c3d4..e5f6g7h8`).
     - `is_incremental`: boolean.
     - `commit_window`: int.
     - `generated_at`: float epoch.
     - `total_files`: int.
     - `files`: list of records with `path`, `status`, `extension`.
2. Update `serialize_yaml(...)`:
   - Include metadata header:
     - `last_commit_hash`: ...
     - `commit_range`: ...
     - `is_incremental`: ...
     - `total_files`: ...
     - `changed_files`: list of entries.
3. Update `serialize_txt(...)`:
   - Include `# last_commit_hash: ...` and `# commit_range: ...` header comments followed by file paths.
4. CLI Enhancements in `parse_cli_args()`:
   - `--no-incremental` / `--full-window` flag to bypass previous checkpoint.
   - `--checkpoint <path>` custom checkpoint path.
   - `--since-commit <sha>` explicitly diff against specific commit.
5. Summary Banner:
   - Display `📌 Last Processed Commit : <prev_sha>`
   - Display `🎯 Current HEAD Commit   : <current_sha>`
   - Display `🔄 Extraction Mode        : INCREMENTAL / FULL WINDOW`
   - Display `📐 Evaluated Range        : <range>`

## Coding Guidelines
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested `if` statements.
- Strict relative paths in logs and outputs.
