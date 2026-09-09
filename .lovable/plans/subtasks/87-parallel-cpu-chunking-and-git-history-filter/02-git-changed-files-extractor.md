# Subtask 87.02: Git Changed Files Extractor CLI

## Objective
Create standalone CLI script `03-ai-scripts/27-git-changed-files.py` to extract changed files from the last $N$ commits into `.lovable/temp/` in JSON, YAML, and TXT formats with dictionary deduplication.

## Proposed Changes
1. CLI Arguments:
   - `--commits`, `-n`: Number of commits to inspect (default: 20).
   - `--output-dir`, `-o`: Target directory (default: `.lovable/temp/`).
   - `--format`: Output format filter (`all`, `json`, `yaml`, `txt`).
   - `--staged-only`: Restrict to only staged/working-tree changes.
   - `--quiet`, `-q`: Suppress stdout summary.
2. File Extraction & Deduplication:
   - Call `git diff --name-only HEAD~N..HEAD`.
   - Call `git status --porcelain` to capture staged and untracked changes.
   - Use dict `{canonical_path: record}` to ensure zero duplicates.
3. Multi-Format Output:
   - `.lovable/temp/git-changed-files.json`
   - `.lovable/temp/git-changed-files.yaml` (pure python yaml serializer without external deps)
   - `.lovable/temp/git-changed-files.txt` (one path per line for easy shell iteration)
4. Speed:
   - Must complete in < 0.2 seconds.

## Acceptance Criteria
- [ ] Script runs via `python 03-ai-scripts/27-git-changed-files.py` with exit code 0.
- [ ] All 3 files generated in `.lovable/temp/`.
- [ ] Zero duplicate paths in any of the output files.
- [ ] Functions <= 15 lines, affirmative booleans, blank lines before returns.
