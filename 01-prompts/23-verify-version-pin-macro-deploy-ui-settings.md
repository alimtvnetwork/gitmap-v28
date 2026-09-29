# 23-verify-version-pin-macro-deploy-ui-settings.md

Use this prompt in any CI pipeline, local test runner, or autonomous AI agent to rigorously verify that installer version pinning, update pinning with fleet SSH, macro deploy commands, UI settings persistence, and repo-secrets flag deduplication are operational and regression-free.

## Pre-flight Quality Gate Checklist

1. **Verify Version Pinning & Installer Table**:
   - Run `gitmap pin gitmap 6.402.0` (or `gitmap installer pin gitmap 6.402.0`) -> Confirms pin persisted in `installation.db`.
   - Run `gitmap installer ls` -> Verify `PINNED` column displays `6.402.0` and summary footer displays `Pinned: 1 installer(s)`.
   - Run `gitmap unpin gitmap` -> Removes pinned version cleanly.

2. **Verify Update Pinning & Fleet SSH**:
   - Run `gitmap update -v 6.402.0 --pin` -> Updates gitmap to target version and pins it.
   - Run `gitmap update -v 6.402.0 --pin --ssh` -> Updates and pins across fleet nodes.
   - Run `gitmap update agm -v 1.2.0 --pin` -> Updates agm/agy to target version and pins it.

3. **Verify Macro Deploy Commands**:
   - Run `gitmap deploy macro all` (or `gitmap macro deploy all`) -> Deploys all macros across all active cluster nodes.
   - Run `gitmap deploy macro <name> <node>` (or `gitmap macro deploy <name> <node>`) -> Targets a specific macro to a specific node.
   - Run `gitmap deploy macro all --except worker-beta,192.168.1.30` -> Successfully excludes specified targets.

4. **Verify UI Settings Persistence**:
   - Access `gitmap cmdui` -> Settings card displays controls for graphics mode (`high`, `fast`, `plain`), auto-open browser toggle, commit-in layout (`split`, `left`, `right`, `stacked`), pull direction (`pull-left`, `pull-right`, `bidirectional`), and PR replay mode (`merges`, `feature-per-commit`, `direct`).
   - Settings are persisted to `~/.gitmap/ui_settings.json` and restored across server restarts.

5. **Verify Repo-Secrets Deduplication**:
   - Confirm `repo-secrets` has canonical flag keys (`isApplyTree`, `isApplyFinalSync`, `isApplyCd`, `isRecreate`, `isPushImmediate`) without redundant legacy aliases.
   - Confirm `cli/cmd/commitin/config_json.go` cleanly handles both canonical and legacy JSON aliases with zero struct duplication.

6. **Verify Clean Cache & Workstation Hygiene**:
   - Run `python 03-ai-scripts/42-clean-test-and-build-caches.py` -> Cleans temp files, Go build/test cache, and temporary test databases.

7. **CI/CD Quality Gate**:
   - Verify all workflows on `main` are 100% green via `gh run list --limit 6`.
