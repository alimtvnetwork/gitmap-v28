# Component Specification: App-Window Launcher, Fullscreen & Modern UI Styling

- **Sequence**: 334
- **Component**: App-Like Browser Experience, Fullscreen & Rich CSS3 Design System
- **Status**: APPROVED

---

## 1. App-Like Browser Launcher Mode
To make browser interfaces behave like a native desktop application:
1. Detect if Google Chrome (`chrome.exe`, `google-chrome`, `/Applications/Google Chrome.app`) or Microsoft Edge (`msedge.exe`, `microsoft-edge`) is installed.
2. If available, launch using `--app=<targetURL>` flag:
   ```bash
   chrome.exe --app=http://localhost:8080/ --new-window
   ```
3. If neither is found or on headless systems, fall back cleanly to `rundll32 url.dll,FileProtocolHandler <url>`, `open <url>`, or `xdg-open <url>`.

---

## 2. Server Exit & Browser Close Endpoint
Both `cmdui` (Fleet Web UI) and `cmdlogviewer` (Show Log UI) implement graceful shutdown:
1. **Server Endpoint**:
   - `POST /api/ui/exit` (Fleet UI)
   - `POST /api/log/exit` (Show Log UI)
   - Responds with `{"status":"shutting_down"}`, triggers `server.Close()` after 100ms, and terminates background process via `os.Exit(0)`.
2. **Client-Side Window Close**:
   - User clicks `✕ Close App` or `✕ Exit`.
   - Sends `fetch('/api/ui/exit', { method: 'POST' })`.
   - Displays gentle goodbye toast / modal ("Closing application...").
   - Executes `window.close()` and sets `window.location.href = 'about:blank'`.

---

## 3. Fullscreen Capability (F11 / Button)
1. Header action button:
   - `<button id="btn-fullscreen" class="btn-subtle" title="Toggle Fullscreen (F11)">⛶ Fullscreen</button>`
2. JavaScript Event Handling:
   - Toggle:
     ```js
     if (!document.fullscreenElement) {
       document.documentElement.requestFullscreen().catch(() => {});
     } else {
       document.exitFullscreen().catch(() => {});
     }
     ```
   - State listener: updates button text/icon when entering (`⤢ Restore`) vs exiting (`⛶ Fullscreen`).

---

## 4. Rich CSS3 Design System & Animations
Adhering to the Catppuccin Mocha & Macchiato design tokens:
1. **Micro-Interactions**:
   - `cubic-bezier(0.16, 1, 0.3, 1)` easing on all interactive elements.
   - Lucrative button glow on hover: `box-shadow: 0 4px 18px -2px rgba(137, 180, 250, 0.35)`.
   - Card lift on hover: `transform: translateY(-2px); box-shadow: 0 10px 30px -4px rgba(0,0,0,0.4)`.
2. **CSS3 Keyframe Animations**:
   - `@keyframes appEntrance`: smooth fadeIn & scale (0.98 -> 1) on load.
   - `@keyframes pulseGlow`: subtle breathing glow on active connection dots.
   - Custom glassmorphic blur and scrollbar styling across all views.
