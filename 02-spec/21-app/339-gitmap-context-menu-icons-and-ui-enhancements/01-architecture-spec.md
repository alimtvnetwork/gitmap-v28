# Master Architecture Spec: GitMap Context Menu Icons, Show-Log Non-Git Guard, Default Browser Decoupling & UI Enhancements

- **Sequence**: 339
- **Title**: Context Menu Icons, Non-Git Show-Log UX, Default Browser Decoupling, Profile/Macro UI
- **Status**: IN_PROGRESS

---

## 1. Problem Statement & Root Cause Analysis

1. **Context Menu Icons**:
   - The root GitMap context menu item (`Directory\shell\gitmap` and `Directory\Background\shell\gitmap`) had `{exe},0` set for its icon, but child entries also wrote `{exe},0`.
   - On Windows 10/11 Explorer, when sub-items have duplicate or generic icon references, or when legacy submenus (`70_git`) remain, child entries render redundant branch/log icons or clash with the root icon.
   - Clean separation is needed: the root menu item gets the official GitMap application icon, while child leaves omit redundant icon paths unless explicitly distinct. Legacy orphaned registry keys (like `70_git`) must be purged upon install/reinstall.

2. **Show-Log on Non-Git Folders**:
   - Right-clicking on workspace container folders (e.g., `D:\work` which contains 69 repositories but is not itself a git repository) ran `gitmap show-log`.
   - `cmdlogviewer.RunShowLog` printed an error to stderr and exited with code 0 without launching the UI:
     `✗ '.' is not a Git repository.`
   - Graceful resolution: If `show-log` is invoked in a non-git directory that contains child repositories, launch the GitMap Web UI dashboard or automatically resolve to the first active/target repository so the user gets an interactive UI rather than an abrupt command-line error.

3. **Chrome Profile Decoupling**:
   - Web UI launcher (`LaunchBrowserAppWindow`) attempted `--app` mode via Google Chrome. Users reported disrupted sessions or isolated profile behavior.
   - Decouple browser launching: Provide default system browser opening via native OS opener (`rundll32 url.dll,FileProtocolHandler` on Windows, `open` on macOS, `xdg-open` on Linux) with an explicit opt-in for app mode or default browser preference without forcing separate profiles.

4. **Installer, Profile, Macro & Settings Context Menu & UI Overhaul**:
   - Provide direct context menu entries for:
     - Installer Manager (`ui installer`)
     - Profile Management (`install profile`, `install list`, tree inspection)
     - Macro Builder UI (`ui macro` and `macro list`)
     - Settings UI (`ui settings`)
   - Expand the Fleet Web UI (`cmdui`) with rich Profile Tree views, single-item install/uninstall actions, profile installer/uninstaller execution, and macro builder/replay UI.
