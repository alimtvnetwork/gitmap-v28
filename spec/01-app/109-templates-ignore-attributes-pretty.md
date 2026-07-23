# Templates: .gitignore, .gitattributes, and Pretty Markdown Renderer

Status: Spec (Phase 0 scaffolding in progress)
Owner: gitmap CLI
Related memory: `mem://features/templates-ignore-attributes`
Related plan: `.lovable/memory/plans/04-templates-ignore-attributes-plan.md`

## 1. Goals

1. Ship curated, language-aware `.gitignore` and `.gitattributes` templates
   with `gitmap`. No more hand-copied snippets.
2. Idempotent merge: running `gitmap add ignore go node` twice must NOT
   duplicate any line. Comments are deduped too.
3. Preserve user-authored content. Managed lines live inside marker blocks.
   User content sits outside, untouched.
4. Default to safe `.gitattributes` (text eol normalization for source files,
   LFS only where it actually helps — `*.svg` is text, not LFS).
5. Render command output (markdown help, changelog) with a small "pretty"
   layer that collapses redundant fenced code blocks and highlights
   `"quoted strings"` in cyan.

## 2. CLI Surface

| Command                          | Alias  | Action                                                    |
|----------------------------------|--------|-----------------------------------------------------------|
| `gitmap add ignore [langs...]`   | `ai`   | Merge `<lang>.gitignore` templates into `./.gitignore`    |
| `gitmap add attributes [langs...]` | `aa` | Merge `<lang>.gitattributes` templates into `./.gitattributes` |
| `gitmap add lfs-install`         | `alfs` | Run `git lfs install` and ensure LFS attrs are present    |
| `gitmap templates list`          | `tl`   | List available languages and template kinds               |
| `gitmap templates show <kind> <lang>` | `ts` | Print one template to stdout                          |

`add` is a real subcommand router (see §6), not a flat command. This keeps
room for `add hooks`, `add editorconfig`, etc. without polluting the top
level.

### 2.1 Resolved open questions

| Question                                  | Decision                                      |
|-------------------------------------------|-----------------------------------------------|
| Alias names                               | `ai` / `aa` (2-letter, matches existing CLI)  |
| Real `add` subcommand router              | Yes. New `dispatchAdd` in `rootcore.go`       |
| Audit trail for curated common template   | Per-template `# source:` header + spec table  |
| Read-only install-dir fallback            | `~/.gitmap/templates/` overlay (user > embed) |
| Pretty renderer normalization fixtures    | `gitmap/render/testdata/pretty/*.in.md` pairs |

## 3. Templates Folder Layout

```
gitmap/templates/
  embed.go                  # //go:embed assets/**
  resolver.go               # overlay: ~/.gitmap/templates > embedded
  materialize.go            # first-run extraction to ~/.gitmap/templates
  assets/
    ignore/
      common.gitignore      # OS junk, IDE, gitmap artifacts
      go.gitignore
      node.gitignore
      python.gitignore
      rust.gitignore
      csharp.gitignore
    attributes/
      common.gitattributes  # text eol=lf for known source extensions
      go.gitattributes
      node.gitattributes
      python.gitattributes
      rust.gitattributes
      csharp.gitattributes
    lfs/
      common.gitattributes  # binary patterns: *.png, *.jpg, *.zip, ...
                            # NOTE: *.svg is text, not LFS
```

Every template file starts with a header:

```
# source: github.com/github/gitignore@<sha>/Go.gitignore (curated)
# kind: ignore
# lang: go
# version: 1
```

The `source:` line is the audit trail. When we re-curate from upstream we
bump `version:` and update the SHA.

## 4. Merge Algorithm

### 4.1 Marker Block

```
# >>> gitmap-ignore (do not edit between markers) >>>
<merged + deduped template lines>
# <<< gitmap-ignore <<<

# user entries (preserved across re-runs)
<existing user-defined rules>
```

### 4.2 Steps

1. Read existing target file (if any).
2. Split into `managed` (between markers) and `user` (everything else).
3. For each requested language: load template from overlay then embedded,
   prepend `common.<kind>` exactly once.
4. Concatenate, then single-pass dedupe with a `map[string]struct{}` keyed
   on the trimmed line. Comments are deduped too.
5. Re-emit: marker block on top, `# user entries` separator, user content
   below.
6. Write atomically (temp file + rename).

### 4.3 Idempotence Test

```
gitmap add ignore go node
gitmap add ignore go node     # MUST be a no-op (byte-identical file)
gitmap add ignore python      # adds python lines, leaves go/node intact
```

## 5. Read-Only Install Fallback

`gitmap` may live in `C:\Program Files\gitmap\` (read-only). On first
template command we materialize the embedded `assets/` to:

- Windows: `%USERPROFILE%\.gitmap\templates\`
- Unix:    `~/.gitmap/templates/`

Resolution order per request:

1. `~/.gitmap/templates/<kind>/<lang>.<ext>`  (user override)
2. embedded `assets/<kind>/<lang>.<ext>`       (fallback)

This lets users edit templates locally without touching the binary. The
overlay is also where the `templates list` command discovers user-added
languages.

## 6. `add` Subcommand Router

New file `gitmap/cmd/rootadd.go`:

```go
func dispatchAdd(command string) bool {
    if command != constants.CmdAdd && command != constants.CmdAddAlias {
        return false
    }
    if len(os.Args) < 3 {
        printAddUsage()
        os.Exit(1)
    }
    sub, rest := os.Args[2], os.Args[3:]
    switch sub {
    case constants.AddSubIgnore, constants.AddSubIgnoreAlias:
        runAddIgnore(rest)
    case constants.AddSubAttributes, constants.AddSubAttributesAlias:
        runAddAttributes(rest)
    case constants.AddSubLFSInstall, constants.AddSubLFSInstallAlias:
        runAddLFSInstall(rest)
    default:
        fmt.Fprintf(os.Stderr, constants.ErrUnknownAddSubcommand, sub)
        os.Exit(1)
    }
    return true
}
```

Wired into `dispatch()` in `root.go`.

## 7. Pretty Markdown Renderer

Lives in `gitmap/render/pretty.go`. Used by help output and changelog
display.

### 7.1 Rules

1. If a fenced code block's content matches the trimmed text of the
   immediately preceding paragraph, collapse to a single yellow line:
   `→ <content>` and drop the fence.
2. Highlight `"double-quoted strings"` in cyan. Single quotes untouched
   (apostrophes).
3. Indent body content under headings by 2 spaces (already implemented in
   changelog UI; mirror in CLI).
4. Subtitles (italic line directly under a heading) render in muted color.

### 7.2 Fixture Corpus

`gitmap/render/testdata/pretty/`:

```
case-001-collapse-redundant-fence.in.md
case-001-collapse-redundant-fence.want.txt
case-002-quoted-strings.in.md
case-002-quoted-strings.want.txt
case-003-subtitle-and-indent.in.md
case-003-subtitle-and-indent.want.txt
case-004-no-collapse-when-different.in.md
case-004-no-collapse-when-different.want.txt
case-005-nested-fences.in.md
case-005-nested-fences.want.txt
```

Test loops over `*.in.md`, renders, compares to `*.want.txt`. New edge
cases get added as paired files, no test code changes needed.

## 8. Phased Delivery

| Phase | Scope                                                      |
|-------|------------------------------------------------------------|
| 0     | Scaffold `templates/` package, embed.FS, materialize, overlay resolver, constants, plan/spec/memory |
| 1     | Seed common + 5 language templates (go, node, python, rust, csharp) |
| 2     | `add` subcommand router + `add ignore` with merge/dedupe   |
| 3     | `add attributes` + `add lfs-install`                       |
| 4     | `templates list` / `templates show`                        |
| 5     | Pretty renderer + fixture corpus                           |
| 6     | Wire pretty renderer into help and changelog CLI output    |

Phase 0 ships first and is independently useful (other commands can read
templates immediately).
