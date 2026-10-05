# Subtask 02: Staging, Packaging & Path Normalization

## Status: DONE
## Owner: Worker Subagent 1 & Lead
## Dependencies: Subtask 01

---

## Deliverables
1. Execution of packaging on local Windows host: 2,572 items bundled into 42.85 MB tarball.
2. Sanitization of workspace JSON paths and project manager paths from Windows (`d:\work\...`) to POSIX (`/home/a/git-work/...`).
3. Binary packing of `state.vscdb`, `conversation-search.db`, and `workspaceStorage`.

---

## Acceptance Criteria
- [x] Tarball archive created in local staging area (`/tmp/` or designated path).
- [x] Verification of tarball contents and size.
- [x] Exclusions of lock files, socket files, GPU caches, and statsig bloat.
