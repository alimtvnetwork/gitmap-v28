# gitmap --theme

Global palette switch for gitmap's terminal output. Pick the look that
matches your terminal so checkmarks, banners, and prompts read clearly
without ANSI noise.

## Usage

    gitmap --theme <bright|standard|monochrome> <command> [...]
    gitmap --theme=standard scan
    GITMAP_THEME=mono gitmap release v1.4.6

Resolved in this order (first wins):

1. `--theme <mode>` on the command line.
2. `GITMAP_THEME` environment variable.
3. Default: `bright`.

The `--theme` flag is **global** — it works in front of any
subcommand and is stripped before the subcommand sees its own
flagset.

## Modes

| Mode | When to use |
|------|-------------|
| `bright` *(default)* | Modern dark-themed terminals (Windows Terminal, iTerm, VS Code). Bright-bold ANSI with accent colors. |
| `standard` | Light themes, older terminals, screen-share recordings where bright-bold is too loud. Downgrades to the pre-v5.13 plain palette. |
| `monochrome` / `mono` | Piping into `diff`, `less` without `-R`, log scrapers, CI logs. Strips every SGR escape so output is plain text. |

## Examples

```
# Brighten everything — the default; explicit form for scripts
gitmap --theme bright status

# Calmer palette for a light terminal
gitmap --theme standard scan

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

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
