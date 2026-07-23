# Help Dashboard — Local Documentation Server

## Overview

The `help-dashboard` command (alias `hd`) serves the gitmap-v27 documentation
site locally in the user's browser. It uses a dual-mode resolution
strategy: static serving when a pre-built `dist/` exists, falling back
to a live Vite dev server otherwise.

---

## Command Signature

```
gitmap-v27 help-dashboard [flags]
gitmap-v27 hd [flags]
```

### Flags

| Flag              | Description                          | Default |
|-------------------|--------------------------------------|---------|
| `--port <number>` | Port to serve the documentation on   | 5173    |

---

## Dual-Mode Resolution

### 1. Static Mode (preferred)

If `<binary-dir>/docs-site/dist/` exists and is a directory:

1. Serve the folder with Go's `net/http.FileServer`.
2. Print the serving URL to stdout.
3. Open the URL in the default browser.
4. Block until Ctrl+C, then shut down gracefully.

**Prerequisites:** None — no external dependencies required.

### 2. Dev Mode (fallback)

If no `dist/` folder is found:

1. Print a fallback notice to stdout.
2. Verify `npm` is on PATH; exit with error if missing.
3. Run `npm install` in the `docs-site/` directory.
4. Run `npm run dev -- --port <port>` as a child process.
5. Open the URL in the default browser.
6. Block until Ctrl+C, then kill the child process.

**Prerequisites:** Node.js and npm on PATH.

---

## Directory Resolution

The command locates the docs site relative to the gitmap-v27 binary:

```
<binary-dir>/
  docs-site/          ← HDDocsDir
    dist/             ← HDDistDir (static assets)
    package.json      ← used by dev fallback
```

Binary path is resolved via `os.Executable()` with symlink
evaluation through `filepath.EvalSymlinks`.

---

## Browser Opening

Cross-platform browser launch:

| OS      | Command                        |
|---------|--------------------------------|
| Windows | `cmd /c start <url>`           |
| macOS   | `open <url>`                   |
| Linux   | `xdg-open <url>`               |

Browser open failures are non-fatal (fire-and-forget via `cmd.Start()`).

---

## Error Handling

| Condition                    | Behavior                                    |
|------------------------------|---------------------------------------------|
| `docs-site/` not found      | Print error to stderr, exit 1               |
| `npm` not on PATH           | Print error to stderr, exit 1               |
| `npm install` fails         | Print error to stderr, exit 1               |
| Dev server fails to start   | Print error to stderr, exit 1               |
| Static server bind fails    | Print error to stderr, exit 1               |
| Browser open fails          | Silently ignored                            |

---

## Shutdown

- **Static mode:** `signal.Notify` captures `SIGINT`/`SIGTERM`,
  calls `server.Close()` for graceful HTTP shutdown.
- **Dev mode:** `signal.Notify` captures `SIGINT`/`SIGTERM`,
  calls `dev.Process.Kill()` to terminate the child process.

Both modes print a "Server stopped." confirmation on exit.

---

## File Layout

| File                                  | Purpose                                  |
|---------------------------------------|------------------------------------------|
| `constants/constants_helpdashboard.go`| CLI names, flag descriptions, messages   |
| `cmd/helpdashboard.go`               | Command handler and server logic         |
| `helptext/help-dashboard.md`         | Embedded help text for `--help` output   |

---

## Acceptance Criteria

1. `gitmap-v27 help-dashboard` serves `dist/` when it exists and opens
   the browser automatically.
2. `gitmap-v27 hd` works as an alias.
3. `--port 8080` serves on port 8080 instead of the default.
4. When no `dist/` exists and Node.js is available, the command
   runs `npm install` then `npm run dev` successfully.
5. When no `dist/` exists and `npm` is not on PATH, the command
   prints a clear error and exits with code 1.
6. Ctrl+C cleanly stops both static and dev servers.
7. `gitmap-v27 help-dashboard --help` prints the embedded help text.
8. All terminal messages use constants from
   `constants_helpdashboard.go` — no magic strings.
