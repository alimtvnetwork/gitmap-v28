# Execution Plan: Clean Temp Scripts & Fix Git Tracking

## Task Summary

The previous executions left temporary python scripts (`fix_*.py`) in the root directory, which were incorrectly committed and tracked by git. These need to be moved to `.lovable/temp-scripts/`, removed from git tracking, and `.gitignore` needs to be updated to ensure this doesn't happen again.

## Actionable Items & Execution Steps

1. **Identify**: Find all `fix_*.py` scripts in the root directory (e.g., `fix_args.py`, `fix_clean.py`, `fix_workflows.py`, etc.).
2. **Move**: Move these files into the `.lovable/temp-scripts/` directory.
3. **Git Hygiene**:
   - Remove the files from git tracking using `git rm`.
   - Verify that `.lovable/temp-scripts/` (or `.lovable/temp-scripts`) is explicitly ignored in `.gitignore`. If not, add it.
4. **Release Bump**:
   - Bump `version.json` minor version.
   - Update `changelog.md` to reflect the removal of untracked AI scripts.
5. **Commit & Push**:
   - Commit the deletions and `.gitignore` updates.
   - Push to `main`.

## Coding Guidelines Checklist (To Enforce)

- [x] Temp Script Sandboxing: strictly enforce `.lovable/temp-scripts/` gitignore rules.
- [x] No generic garbage variable names.
- [x] Boolean conventions used (is/has prefixes, no negatives).
- [x] Git Working Tree is completely clean after execution.
