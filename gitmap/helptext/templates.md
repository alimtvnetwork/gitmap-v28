# gitmap templates

Discover and inspect the embedded `.gitignore` / `.gitattributes` template
corpus that powers `gitmap add ignore`, `gitmap add attributes`, and
`gitmap add lfs-install`. Two read-only subcommands.

## Alias

tpl

## Subcommands

| Subcommand | Alias | Purpose |
|------------|-------|---------|
| `templates list [--kind] [--lang]` | `tl` | Print every available template with KIND, LANG, SOURCE, PATH |
| `templates show <kind> <lang>` | `ts` | Write a single resolved template (overlay > embed) to stdout |
| `templates init <lang>...` | `ti` | Scaffold `.gitignore` / `.gitattributes` for one or more languages |

## Flags (list)

| Flag | Description |
|------|-------------|
| `--kind <ignore\|attributes\|lfs>` | Narrow output to one kind. Unknown values exit 1. |
| `--lang <name>` | Narrow output to one language across every kind. Case-insensitive. |

Flags AND together — `--kind ignore --lang go` returns at most one row
(`ignore/go`), not the union of every ignore row and every go row.

## Kinds

| Kind | File extension | Used by |
|------|----------------|---------|
| `ignore` | `.gitignore` | `gitmap add ignore` |
| `attributes` | `.gitattributes` | `gitmap add attributes` |
| `lfs` | `.gitattributes` (LFS lines) | `gitmap add lfs-install` |

## Source resolution

Every template lookup checks two locations, in order:

1. **User overlay** — `~/.gitmap/templates/<kind>/<lang>.<ext>`
   Source label: **`user`**.
2. **Embedded corpus** — bundled into the gitmap binary via `go:embed`.
   Source label: **`embed`**.

The first hit wins. `templates list` shows which one each entry resolves
to so you can tell at a glance which templates you've forked.

## Examples

### Example 1: List every template

    gitmap templates list

**Output:**

    KIND        LANG            SOURCE  PATH
    ignore      common          embed   assets/ignore/common.gitignore
    ignore      csharp          embed   assets/ignore/csharp.gitignore
    ignore      go              embed   assets/ignore/go.gitignore
    ignore      java            embed   assets/ignore/java.gitignore
    ignore      kotlin          embed   assets/ignore/kotlin.gitignore
    ignore      node            user    /home/me/.gitmap/templates/ignore/node.gitignore
    ignore      php             embed   assets/ignore/php.gitignore
    ignore      python          embed   assets/ignore/python.gitignore
    ignore      ruby            embed   assets/ignore/ruby.gitignore
    ignore      rust            embed   assets/ignore/rust.gitignore
    ignore      swift           embed   assets/ignore/swift.gitignore
    attributes  common          embed   assets/attributes/common.gitattributes
    attributes  csharp          embed   assets/attributes/csharp.gitattributes
    attributes  go              embed   assets/attributes/go.gitattributes
    attributes  java            embed   assets/attributes/java.gitattributes
    attributes  kotlin          embed   assets/attributes/kotlin.gitattributes
    attributes  node            embed   assets/attributes/node.gitattributes
    attributes  php             embed   assets/attributes/php.gitattributes
    attributes  python          embed   assets/attributes/python.gitattributes
    attributes  ruby            embed   assets/attributes/ruby.gitattributes
    attributes  rust            embed   assets/attributes/rust.gitattributes
    attributes  swift           embed   assets/attributes/swift.gitattributes
    lfs         common          embed   assets/lfs/common.gitattributes

The `node` row above shows what a forked template looks like: SOURCE flips
from `embed` to `user` and PATH points at the absolute overlay file.

### Example 2: Print a single template to stdout

    gitmap templates show ignore go

**Output:** the raw bytes of `ignore/go.gitignore` (overlay if present,
otherwise embed), audit-trail header included:

    # source: github/gitignore
    # kind: ignore
    # lang: go
    # version: 1
    *.exe
    *.test
    *.out
    ...

### Example 3: Diff your overlay against the curated embed

    gitmap templates show ignore node > /tmp/curated-node.gitignore
    diff ~/.gitmap/templates/ignore/node.gitignore /tmp/curated-node.gitignore

`templates show` always resolves overlay-first — but **once your overlay
file exists**, the embed copy is the only way to recover the curated
bytes. Pipe `templates show` through `diff` to audit your fork before
re-syncing.

### Example 4: Use the short aliases

    gitmap tpl tl
    gitmap tpl ts attributes common

Both `tpl` (the umbrella alias) and `tl` / `ts` (the per-subcommand
aliases) round-trip identically with their long forms.

### Example 5: Filter the listing

    gitmap templates list --kind ignore
    gitmap tpl tl --lang go
    gitmap templates list --kind attributes --lang node

Filters are case-insensitive and AND together. An unknown `--kind` value
exits 1 with `templates list: unknown --kind "foo" (want ignore | attributes | lfs)`
so typos surface immediately instead of silently emptying the table.
A valid filter that simply matches no rows prints
`(no templates match the requested filter)`.

## How forking works

To customize a template, copy the embedded version to the overlay path
and edit it:

    mkdir -p ~/.gitmap/templates/ignore
    gitmap templates show ignore python > ~/.gitmap/templates/ignore/python.gitignore
    $EDITOR ~/.gitmap/templates/ignore/python.gitignore

Subsequent `gitmap add ignore python` calls (and any future `add` flow
that resolves `ignore/python`) will pick up your overlay automatically.
`gitmap templates list` will report SOURCE=`user` for that row.

To revert a fork, just delete the overlay file — the next resolve falls
back to `embed`.

## Pretty rendering

`templates show` writes raw bytes by default — perfect for the diff and
redirect workflows above. When the resolved template is **markdown**
(`.md` / `.markdown`) **and** stdout is a real TTY, the output is routed
through the same pretty markdown renderer used by `gitmap help`:

- Cyan `"double quotes"` for emphasized terms.
- Yellow `→ collapsed` lines when a fenced block restates the
  preceding paragraph.
- Muted subtitles under headings, indented bodies.

Today the embedded corpus is `.gitignore` / `.gitattributes` only, so
this kicks in for **markdown overlays** you drop into
`~/.gitmap/templates/<kind>/<lang>.md` and for any future markdown
templates added to the embed.

Three opt-outs / opt-ins, all honored across the CLI:

- `--pretty` — force ANSI rendering even when stdout is not a TTY
  (handy for `gitmap templates show notes intro.md --pretty | less -R`).
- `--no-pretty` (alias: legacy `--raw`) — strip ANSI even on a TTY.
- `GITMAP_NO_PRETTY=1` — environment opt-out, shared with `gitmap help`
  and `gitmap changelog`. Set it once in your shell profile to disable
  pretty rendering across the whole CLI.

Pipes and redirects automatically bypass the renderer (stdout is no
longer a TTY), so `templates show foo bar > out.md` always writes the
unmodified bytes.

## Notes

- `templates list` and `templates show` are pure reads. They never
  write to disk and never invoke git.
- Unknown `<kind>` or `<lang>` arguments to `templates show` exit 1
  with `template not found: kind=… lang=…`.
- The embedded corpus is versioned via the `# version: N` header on
  each file. When that integer bumps, gitmap re-resolves cleanly — but
  user overlays are **never** auto-upgraded; you decide when to refresh.

## See Also

- [add lfs-install](add-lfs-install.md) — Install Git LFS hooks + merge `lfs/common.gitattributes`
- [lfs-common](lfs-common.md) — Per-pattern `git lfs track` flow (no template)
- [setup](setup.md) — Configure Git global settings
- [doctor](doctor.md) — Diagnose binary, PATH, and config issues

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter templates
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
