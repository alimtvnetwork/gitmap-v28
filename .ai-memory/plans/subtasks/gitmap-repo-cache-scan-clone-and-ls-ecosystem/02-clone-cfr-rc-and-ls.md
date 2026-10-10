# Subtask Plan 02: Clone & LS Engine: clone/cfr/cfrp --rc/rc, Split-DB Cache, Short TreeView & ls rc

- **Parent Task Slug:** `gitmap-repo-cache-scan-clone-and-ls-ecosystem`
- **Subtask ID:** `Task-02`
- **Owner:** Worker 02
- **Status:** PENDING
- **Owned Files:**
  - `cli/cmdclone/clone_rc.go`
  - `cli/cmdclone/clone_rc_db.go`
  - `cli/cmdclone/clone_rc_tree.go`
  - `cli/cmdclone/clone.go`
  - `cli/cmdclone/clonefixrepo.go`
  - `cli/cmdlist/list_rc.go`
  - `cli/cmdlist/list.go`
  - `cli/cmdclone/clone_rc_test.go`

---

## 1. Technical Objective
Implement `gitmap clone/cfr/cfrp --rc/rc` and `gitmap ls rc`:
1. In `cli/cmdclone/clone.go` & `cli/cmdclone/clonefixrepo.go`:
   - Intercept `--rc`, `-rc`, `--repo-cache`, `rc`, `repo-cache`.
   - Dispatch to `RunCloneRC(args, mode)`.
2. In `cli/cmdclone/clone_rc_db.go`:
   - SQLite Split-DB engine in `.gitmap/data/repocache/sql.db` (`RepoCacheManifest`, `RepoCacheEntry`).
   - Invalidation based on mtime and SHA256 file hash (<5ms).
3. In `cli/cmdclone/clone_rc_tree.go`:
   - Format raw URLs to short `owner/repo` (e.g. `alimtvnetwork/gitmap-v28`).
   - Classify protocol: `[SSH]` vs `[Public HTTPS]`.
   - Render interactive TreeView menu with option to clone all, specific manifest, or specific repo.
4. In `cli/cmdclone/clone_rc.go`:
   - Main orchestrator handling `-y` / `--yes` flag, numbered arguments, and dispatching to `cloner.RunMultiClone` or CFR engine.
5. In `cli/cmdlist/list.go` & `cli/cmdlist/list_rc.go`:
   - Intercept `gitmap ls rc` / `gitmap list rc`.
   - Query all manifests, display breakdown of `[SSH]` vs `[Public HTTPS]`, and output actionable import command hints.
6. In `cli/cmdclone/clone_rc_test.go`:
   - Unit tests verifying discovery priority (`01-gitmap/gitmap.json` first), Split-DB caching, URL sanitization, and flag parsing.
