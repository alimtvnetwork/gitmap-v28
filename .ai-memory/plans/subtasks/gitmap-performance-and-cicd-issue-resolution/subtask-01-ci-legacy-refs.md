# Subtask 01: Fix Legacy Refs in Readme Conventions Spec and Pass CI Script Tests

- **Subtask ID**: `subtask-01-ci-legacy-refs`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `02-spec/01-spec-authoring-guide/13-root-readme-conventions.md`
- **Assigned Worker**: Worker 01
- **Status**: PENDING

---

## 1. Objectives
1. Remove legacy reference violations matching `/gitmap-v[567]\b/` in `02-spec/01-spec-authoring-guide/13-root-readme-conventions.md`:
   - Line 20: Replace `gitmap-v6` with `gitmap-v28`.
   - Line 257: Replace `gitmap-v6` with `gitmap-v28`.
2. Verify `python .github/scripts/check-legacy-refs.py` exits with code 0.
3. Verify `test_legacy_refs_runs_clean_on_repo` in `.github/scripts/tests/test_ci_scripts.py` passes cleanly.

---

## 2. Verification
- `python .github/scripts/check-legacy-refs.py`
