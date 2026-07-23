# gitmap help-dashboard

Serve the interactive documentation site locally in your browser.

**Alias:** `hd`

## Usage

```
gitmap help-dashboard [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--port <number>` | Port to serve on (default: 5173) |

## How It Works

1. Locates the `docs-site/` directory relative to the gitmap binary
2. If `docs-site/` is missing but `docs-site.zip` exists, extracts it automatically
3. If a pre-built `dist/` folder exists, serves it with a built-in HTTP server
4. If no `dist/` found, falls back to `npm install && npm run dev`
5. Opens the dashboard in your default browser automatically

## Prerequisites

- **Static mode**: No dependencies — serves pre-built files directly
- **Auto-extract mode**: `docs-site.zip` is downloaded by the installer and extracted on first run
- **Dev mode (fallback)**: Requires Node.js and npm on PATH

## Examples

    $ gitmap help-dashboard

    Serving docs from /usr/local/bin/docs-site/dist on http://localhost:5173
    Opening http://localhost:5173 in browser...

    $ gitmap hd --port 8080

    No pre-built dist/ found, falling back to npm run dev
    Running npm install...
    Starting dev server from /usr/local/bin/docs-site...
    Opening http://localhost:8080 in browser...

Press Ctrl+C to stop the server.

## Browseable Pages

The dashboard exposes a help page for every CLI command. The pages
listed below were added in v2.96.0 — v2.98.0 and may not yet be
linked from older bookmarks:

| Page | Command | Spec |
|------|---------|------|
| `/diff` | `gitmap diff LEFT RIGHT` (alias `df`) | spec/01-app/97-move-and-merge.md (companion) |
| `/mv` | `gitmap mv LEFT RIGHT` (alias `move`) | spec/01-app/97-move-and-merge.md |
| `/merge-both` | `gitmap merge-both LEFT RIGHT` (alias `mb`) | spec/01-app/97-move-and-merge.md |
| `/merge-left` | `gitmap merge-left LEFT RIGHT` (alias `ml`) | spec/01-app/97-move-and-merge.md |
| `/merge-right` | `gitmap merge-right LEFT RIGHT` (alias `mr`) | spec/01-app/97-move-and-merge.md |
| `/commit-right` | `gitmap commit-right LEFT RIGHT` (alias `cmr`) **LIVE** | spec/01-app/106-commit-left-right-both.md |
| `/commit-left` | `gitmap commit-left LEFT RIGHT` (alias `cml`) *(scaffold)* | spec/01-app/106-commit-left-right-both.md |
| `/commit-both` | `gitmap commit-both LEFT RIGHT` (alias `cmb`) *(scaffold)* | spec/01-app/106-commit-left-right-both.md |
| `/as` | `gitmap as [name]` (alias `s-alias`) | helptext/as.md |
| `/release-alias` | `gitmap release-alias <name> <ver>` (alias `ra`) | helptext/release-alias.md |
| `/release-alias-pull` | `gitmap release-alias-pull <name> <ver>` (alias `rap`) | helptext/release-alias-pull.md |

The same pages render in the CLI via `gitmap help <command>`
(e.g. `gitmap help merge-both`, `gitmap help release-alias`).

## See Also

- docs — Open the hosted documentation website
- dashboard — Generate an HTML analytics dashboard for a repo
- mv / merge-both / merge-left / merge-right — Move & merge command family
- commit-right / commit-left / commit-both (cmr / cml / cmb) — Commit-replay family (spec 106)
- diff — Read-only preview before merge-both
- as / release-alias / release-alias-pull — Register-and-release-anywhere workflow


## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter help-dashboard
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
