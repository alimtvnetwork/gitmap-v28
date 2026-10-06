# GitMap v6.493.0

## What's Changed in v6.493.0

- **Fleet Repository Deployment, Remote Scanner Delegation & Multi-IDE Fleet Sync Engine**:
  - Added `gitmap nodes deploy repo <slug>` and `gitmap nodes deploy repos [targets]`:
    - Deploys any repository from the local machine across remote fleet nodes (`w1`, `w2`, `u1`).
    - Streams in-memory tarball archive over encrypted SSH channels or delegates origin Git clone.
    - Automatically executes remote `gitmap rescan` so the target node's `gitmap.db` immediately indexes the new repository.
    - Automatically registers the deployed repository across all installed IDEs on the target node:
      - **VS Code**: Added to Project Manager `projects.json`.
      - **Cursor**: Added to Project Manager `projects.json`.
      - **Antigravity IDE**: Added to `~/.gemini/config/projects/<uuid>.json`.
      - **GitHub Desktop**: Registered via the `github <path>` CLI shim.
    - Added `--with-pinned` flag to sync pinned project status in `pinned_projects.json`.
    - Added `--with-conversations` flag to bundle, remap, and transfer associated Antigravity conversations (`<conv-id>.db`) and brain transcript logs (`brain/<conv-id>/`).
    - Added granular targeting and filtering flags: `--target <alias>`, `--except <list>` (default `main`), `--include <list>`, `--exclude <list>`, `--open-only`, `--dest <path>`, `--dry-run`, and `--json`.
  - Added `gitmap nodes scan` and `gitmap nodes rescan`:
    - Broadcasts remote repository discovery across all reachable fleet nodes via SSH delegation.
  - Core SSH Router Enhancement:
    - Updated `isGitmapCoreCommand` in `cli/cmdssh/ssh_exec_command.go` to recognize `scan` and `rescan`, allowing direct invocation via `gitmap ssh exec <node> scan`.
  - Specification & Planning:
    - Specification: `02-spec/21-app/229-nodes-deploy-repos-and-multi-ide-fleet-sync/01-architecture-spec.md`
    - Component Spec: `02-spec/21-app/229-nodes-deploy-repos-and-multi-ide-fleet-sync/02-component-and-fleet-spec.md`
    - Execution Plan: `.ai-memory/plans/229-nodes-deploy-repos-and-multi-ide-fleet-sync.md`

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.493.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.493.0/install.sh | sh
```
