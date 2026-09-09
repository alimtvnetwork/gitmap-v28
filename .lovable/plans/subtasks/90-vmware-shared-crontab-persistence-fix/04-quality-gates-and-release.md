# Subtask 04: Quality Gate Verification, CI Local Runner & Release

**Parent Plan:** `.lovable/plans/pending/90-vmware-shared-crontab-persistence-fix.md`  
**Target Files:** `linter-scripts/`, `03-ai-scripts/`  

## Objectives
1. Run all repository linters:
   - `python linter-scripts/check-nested-ifs.py`
   - `python linter-scripts/check-enum-and-boolean.py`
   - `python linter-scripts/check-error-management.py`
   - `python .github/scripts/go-format-check.py`
   - `python .github/scripts/tests/test_ci_scripts.py`
2. Update `.lovable/plans/01-index.md` and move plan to completed.
3. Run release orchestrator:
   - `python 03-ai-scripts/29-release-orchestrator.py --tier patch --scope "Fix VMware shared crontab persistence bad minute error and add Linux E2E tests"`
4. Push `main`, `release/v6.204.8`, and tags to `origin`.
