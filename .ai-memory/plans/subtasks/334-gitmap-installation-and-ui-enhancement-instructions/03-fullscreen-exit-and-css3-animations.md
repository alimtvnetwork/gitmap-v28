# Subtask 03: Fullscreen, Exit Endpoints & CSS3 Animations

- **Plan**: `334-gitmap-installation-and-ui-enhancement-instructions`
- **Subtask**: 03
- **Status**: PENDING

---

## Action Items
1. In `cli/cmdui/ui_server.go`, add `/api/ui/exit` endpoint shutting down the HTTP server and terminating the background process.
2. In `cli/cmdui/ui_assets_markup_misc.go`, add Fullscreen (`⛶`) and Exit (`✕ Close`) header buttons.
3. In `cli/cmdui/ui_assets_js_settings.go`, wire the fullscreen toggle and exit API call with `window.close()`.
4. In `cli/cmdlogviewer/ui_assets_html.go` & `ui_assets_script.go`, add `#btn-fullscreen` and wire Fullscreen API.
5. In `cli/cmdui/ui_assets_theme.go` and `cli/cmdlogviewer/ui_assets_html.go`, elevate CSS with rich micro-interactions, cubic-bezier transitions, glassmorphic cards, and keyframe animations (`@keyframes appEntrance`, `@keyframes pulseGlow`).
