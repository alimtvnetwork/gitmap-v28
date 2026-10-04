# Subtask 192.5: Repo Secrets GitMap Configuration Deduplication & Go Struct Cleanup

## Spec Reference
- `02-spec/21-app/182-version-pinning-macro-deploy-ui-settings-secret-flags.md` §2.4

## Deliverables
1. **Repo Secrets Synchronization**:
   - `git pull` on `./repo-secrets` (already verified and up to date).
   - In `./repo-secrets\01-gitmap\commit-pull-config.json` and `00-commit-pull-config.json`, ensure single clean canonical flags (e.g., `isApplyCd`, `isApplyTree`, `isApplyFinalSync`, `isRecreate`, `isPushImmediate`) without redundant duplicate keys.
   - Commit and push changes to `repo-secrets`.
2. **Go Struct Normalization**:
   - In `cli/cmd/commitin/config_json.go`, clean up `CommitInConfigJSON` struct, eliminating confusing redundant fields while safely preserving compatibility with existing configs.
