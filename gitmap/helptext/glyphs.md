# gitmap --glyphs

Global switch that selects between **rich** (emoji) and **safe**
(ASCII) glyph rendering. Use it to fix mojibake (`Γ£ô`, `≡ƒôª`) on
terminals whose font or code page can't render the rich set.

## Usage

    gitmap --glyphs <auto|rich|safe> <command> [...]
    gitmap --glyphs=safe scan
    GITMAP_GLYPHS=safe gitmap release v1.0.0

Resolution order (first wins):

1. `--glyphs <mode>` on the command line.
2. `GITMAP_GLYPHS` environment variable.
3. **auto** — picks `safe` on legacy Windows ConsoleHost
   (powershell.exe 5.1 / cmd.exe), `rich` everywhere else
   (Windows Terminal, VS Code, iTerm2, all *nix TTYs).

## Modes

| Mode | When to use |
|------|-------------|
| `auto` *(default)* | Best for most users — detects the host and picks the right set. |
| `rich` | Force emoji output even on detected legacy hosts (you've installed a Nerd Font). |
| `safe` | Force ASCII — for log scrapers, CI logs, screen-share recordings, or any host where emoji render as boxes / mojibake. |

## Example translations (safe mode)

| Rich   | Safe     |
|--------|----------|
| `✓`    | `v`      |
| `✗`    | `x`      |
| `→`    | `->`     |
| `⚠`    | `!`      |
| `✅`   | `[OK]`   |
| `❌`   | `[X]`    |
| `📦`   | `[pkg]`  |
| `🎉`   | `[done]` |
| `📁`   | `[dir]`  |
| `🏷`   | `[tag]`  |

## See also

- **`gitmap --theme`** — controls ANSI color palette. Orthogonal to glyphs.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter glyphs
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap glyphs
```
