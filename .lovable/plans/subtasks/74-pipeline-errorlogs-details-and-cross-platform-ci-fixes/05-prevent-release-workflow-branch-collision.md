# Subtask 05: Prevent Release Workflow Race Condition Collision

## Scope
- In `.github/workflows/release.yml`:
  - Restrict the `push` trigger to tags matching `v*` only (remove branches `release/*` from trigger), or configure `softprops/action-gh-release@v2` with concurrency groups to prevent duplicate parallel executions from attempting to upload identical assets at the same time.
- In `03-ai-scripts/29-release-orchestrator.py`:
  - When pushing release artifacts, push the release branch and tag with appropriate sequence or push only the release tag.
- Keep all functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.

## Files Touched
- `.github/workflows/release.yml`
- `03-ai-scripts/29-release-orchestrator.py`
