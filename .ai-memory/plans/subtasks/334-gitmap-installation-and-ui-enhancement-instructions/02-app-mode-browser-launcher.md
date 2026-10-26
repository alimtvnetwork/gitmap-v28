# Subtask 02: App-Mode Browser Launcher

- **Plan**: `334-gitmap-installation-and-ui-enhancement-instructions`
- **Subtask**: 02
- **Status**: PENDING

---

## Action Items
1. Create a helper `launchBrowserAppWindow(targetURL string) bool` in `cli/cmdchrome/` or `cli/cmdbrowse/`:
   - Checks Chrome (`cmdchrome.FindChromeBinaryPath()`) or Edge (`msedge.exe`).
   - If found, launches `exec.Command(bin, "--app="+targetURL, "--new-window")`.
2. Connect `openBrowserURL` in `cli/cmdui/ui_server.go` and `openBrowser` in `cli/cmdlogviewer/server.go` to use this app-launcher helper.
3. Keep default browser opener as graceful fallback if no app-capable browser is available.
4. Verify non-blocking detached process behavior.
