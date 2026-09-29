# Subtask 192.4: UI Settings Graphics & Layout Controls

## Spec Reference
- `02-spec/21-app/182-version-pinning-macro-deploy-ui-settings-secret-flags.md` §2.3

## Deliverables
1. **API & Model**:
   - Extend `SettingsData` in `cli/cmdui/ui_types.go` with `GraphicsMode`, `AutoOpenBrowser`, `CommitInLayout`, `PullDirection`, `PRReplayMode`, and `PinnedVersions`.
   - Update `handleAPISettings` in `cli/cmdui/ui_server.go` to support GET and POST persisting settings.
2. **Web Frontend**:
   - In `cli/cmdui/ui_assets.go`, add controls in the Settings tab for graphics modes, layout orientation (`split`, `left`, `right`), commit pull modes, and pull directions.
