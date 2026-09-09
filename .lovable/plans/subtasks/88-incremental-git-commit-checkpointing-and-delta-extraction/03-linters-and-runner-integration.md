# Subtask 03: Linters & Local Runner Integration

## Objective
Update linters (`check-relative-paths.py`, `check-nested-ifs.py`, `check-enum-and-boolean.py`) and CI runner (`06-cicd-local-runner.py`) to leverage incremental commit checkpoints safely and backward-compatibly.

## Requirements
1. `linter-scripts/check-relative-paths.py`:
   - Safely read `files` from manifest whether entries are dicts (`item["path"]`) or strings (`item`).
   - Display checkpoint metadata when `--changed-only` is active: `[Checkpoint: <last_hash> | Range: <commit_range>]`.
2. `linter-scripts/check-nested-ifs.py`:
   - Safely read `files` from manifest using backward-compatible extractor.
3. `linter-scripts/check-enum-and-boolean.py`:
   - Safely read `files` from manifest using backward-compatible extractor.
4. `03-ai-scripts/06-cicd-local-runner.py`:
   - Support `--no-incremental` flag forwarding to `27-git-changed-files.py`.
   - Prevent memory oversubscription / Windows commit limit exhaustion:
     - Run `Smart Unit Tests & Coverage` sequentially or limit test parallelism to prevent running simultaneously with `Web App Build`.

## Coding Guidelines
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested `if` statements.
