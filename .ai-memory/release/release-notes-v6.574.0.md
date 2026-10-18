# GitMap v6.574.0

## What's Changed in v6.574.0

### CI Pipeline Defect Resolutions & Quality Gates

1. **GitHub Auth Test Mock Parity (`cli/usercontext`)**:
   - Updated `TestDetectGhAuth_FallbackToken` in `cli/usercontext/gh_auth_test.go` to handle `git config --global user.name` arguments, fixing cross-platform build test failures across Windows, macOS, and Ubuntu.

2. **Error Management Compliance (`cli/cmddb`, `cli/heatmap`)**:
   - Resolved 6 swallowed database call errors across `cli/cmddb/cmddb_tables.go`, `cli/heatmap/heatmap_conflict.go`, and `cli/heatmap/heatmap_view.go`. All database queries and statement executions now explicitly handle errors and return safe defaults or propagate errors.

3. **Legacy Version References Purged (`02-spec/01-spec-authoring-guide`)**:
   - Replaced legacy references (`gitmap-v6`) with canonical `gitmap-v28` and `gitmap` in `02-spec/01-spec-authoring-guide/13-root-readme-conventions.md`, restoring clean passing status for `.github/scripts/check-legacy-refs.py` and `test_ci_scripts.py`.

4. **Codebase Go Formatting Cleanliness (`cli/`)**:
   - Formatted all outstanding Go files with `gofmt -w` (`cli/cmd/rootutility_update.go`, `cli/cmdclone/clone_auth_resolve.go`, `cli/cmdclone/clone_auth_types.go`, `cli/cmdupdate/update_ui.go`, `cli/heatmap/heatmap_query.go`, `cli/heatmap/heatmap_view.go`, `cli/store/heatmap_rc_store.go`), verifying 100% clean formatting across the codebase.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.574.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.574.0/install.sh | sh
```
