# Subtask 04: Verification & Quality Gates

## Objective
Run verification tests for incremental commit checkpointing across all scenarios, then run quality gates to confirm green status.

## Test Scenarios
1. **Scenario 1: Fresh Run**:
   - Run `python 03-ai-scripts/27-git-changed-files.py --commits 20`.
   - Verify `.lovable/temp/git-changed-files.json` contains `last_commit_hash` equal to `git rev-parse HEAD`.
2. **Scenario 2: Immediate Re-run (Same Commit)**:
   - Re-run `python 03-ai-scripts/27-git-changed-files.py`.
   - Verify extraction mode is `INCREMENTAL`, range is `<hash>..HEAD (no new commits)`, and only uncommitted working-tree files are returned.
3. **Scenario 3: Forced Full Window**:
   - Run `python 03-ai-scripts/27-git-changed-files.py --no-incremental`.
   - Verify extraction mode is `FULL WINDOW` with full 20-commit history returned.
4. **Scenario 4: Linter Integration**:
   - Run `python linter-scripts/check-relative-paths.py --changed-only`.
   - Run `python linter-scripts/check-nested-ifs.py --changed-only`.
   - Run `python linter-scripts/check-enum-and-boolean.py --changed-only`.
   - Ensure all exit 0 with 0 violations.
5. **Scenario 5: Local CI Runner**:
   - Run `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`.
   - Ensure all gates exit 0 cleanly.
