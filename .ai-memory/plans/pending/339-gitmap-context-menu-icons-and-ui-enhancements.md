# Plan: GitMap Context Menu Icons, Show-Log Non-Git Guard, Browser Decoupling & UI Enhancements

- **ID**: 339-gitmap-context-menu-icons-and-ui-enhancements
- **Status**: IN_PROGRESS

## Subtasks
1. `01-context-menu-icons-and-cleanup.md`: Clean orphaned registry entries, ensure crisp root icon on `gitmap`, remove confusing child icons.
2. `02-show-log-non-git-guard.md`: Handle non-git folders gracefully in `cmdlogviewer.RunShowLog`, auto-fallback to child repo or Web UI dashboard.
3. `03-browser-decouple-chrome-profile.md`: Launch system default browser cleanly without forcing Chrome profiles or isolated user-data.
4. `04-installer-profiles-macro-settings-ui.md`: Wire Macro, Profile Tree, Installer and Settings into context menu and Web UI tabs.
5. `05-verification-and-release.md`: Run test suites, gofmt, update changelogs, minor version bump to `v6.627.0`.
