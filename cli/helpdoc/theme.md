# gitmap --theme

Global palette switch for gitmap's terminal output. Pick the look that
matches your terminal so checkmarks, banners, and prompts read clearly
without ANSI noise.

## Usage

    gitmap --theme <slug> <command> [...]
    gitmap --theme=tokyo-night scan
    GITMAP_THEME=mono gitmap release v1.4.6

Resolved in this order (first wins):

1. `--theme <slug>` on the command line.
2. `GITMAP_THEME` environment variable.
3. Default: `navy-purple`.

The `--theme` flag is **global** — it works in front of any
subcommand and is stripped before the subcommand sees its own
flagset.

## Themes

Named palettes from the coding-guidelines design-system theme
catalogue (semantic roles: fg, fg-muted, fg-subtle, primary, accent,
warn, error, ok — never literal colors in code). Piped output never
carries ANSI: colors render on a real TTY only.

| Theme | When to use |
|-------|-------------|
| `navy-purple` *(default)* | Developer tools. Deep midnight navy + electric violet primary + ice cyan accent. |
| `tokyo-night` | Modern CLI tools. Storm blue-grey canvas with neon pastel accents. |
| `vscode-dark-plus` | Instant familiarity for engineers. Classic VS Code Dark+ palette. |
| `monokai-pro` | High-impact telemetry and status. Warm slate + fluorescent rose/gold. |

Legacy aliases still work: `bright` → `navy-purple`, `standard` →
`vscode-dark-plus`, `monochrome` / `mono` → strip every SGR escape
(for piping into `diff`, `less` without `-R`, log scrapers, CI logs).

## Examples

```

# Default palette — the dev-tools look; explicit form for scripts

gitmap --theme navy-purple status

# Neon pastel palette for a modern CLI feel

gitmap --theme tokyo-night scan

# Classic IDE palette

gitmap --theme vscode-dark-plus scan

# Strip all color for log scraping

gitmap --theme=mono release v2.0.0 > release.log

# Pin via env so every subprocess inherits it

export GITMAP_THEME=monochrome
gitmap clone json --target-dir ./projects
```

## How it works

The flag is parsed once at startup and written to `GITMAP_THEME`. For
`standard` and `monochrome` modes, gitmap wraps `os.Stdout` /
`os.Stderr` with a tiny ANSI rewrite filter so every existing `Print`
in the codebase adapts automatically — no per-command wiring. For
`bright` the filter is bypassed entirely (zero overhead).

TTY detection (`gitmap help foo` pretty-rendering, prompts) is
captured BEFORE the pipe wrap, so theming a piped run does not
falsely demote the renderer.

## See also

- **`gitmap help`** — most subcommands honor `--pretty` / `--no-pretty`
  for markdown rendering, which is orthogonal to `--theme`.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter theme
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
