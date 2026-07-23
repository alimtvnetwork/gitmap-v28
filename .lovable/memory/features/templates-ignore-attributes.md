---
name: templates-ignore-attributes
description: Embedded .gitignore/.gitattributes templates with idempotent marker-block merge, ~/.gitmap/templates overlay, add ignore/attributes/lfs-install subcommands, pretty markdown renderer
type: feature
---

# Templates: ignore, attributes, pretty renderer

## Summary

`gitmap` ships curated `.gitignore` and `.gitattributes` templates per
language, plus a pretty markdown renderer for CLI output.

## CLI

- `gitmap add ignore [langs...]` (alias `ai`)
- `gitmap add attributes [langs...]` (alias `aa`)
- `gitmap add lfs-install` (alias `alfs`)
- `gitmap templates list` (alias `tl`)
- `gitmap templates show <kind> <lang>` (alias `ts`)

`add` is a real subcommand router (see `gitmap/cmd/rootadd.go`).

## Storage

- Embedded: `gitmap/templates/assets/{ignore,attributes,lfs}/<lang>.<ext>`
- User overlay: `~/.gitmap/templates/...` (materialized on first run)
- Resolution: overlay wins, embed is fallback. Lets users in read-only
  install paths (e.g. `C:\Program Files\gitmap`) still customize.

## Merge contract

- Managed content lives between markers:
  `# >>> gitmap-ignore ... >>>` ... `# <<< gitmap-ignore <<<`
- User content preserved below `# user entries` separator.
- Single-pass dedupe via `map[string]struct{}` keyed on trimmed line.
  Comments dedupe too. Re-running with same args is byte-identical no-op.

## Defaults that differ from upstream

- `*.svg` is `text eol=lf`, NOT LFS (it's XML, diffs cleanly).
- `common.gitattributes` sets `text eol=lf` for all known source
  extensions across languages, so cross-platform checkouts behave.

## Audit trail

Every template file starts with:
```
# source: github.com/github/gitignore@<sha>/<file> (curated)
# kind: ignore | attributes | lfs
# lang: <lang>
# version: <int>
```
Bump `version:` when re-curating from upstream.

## Pretty renderer

`gitmap/render/pretty.go`. Rules:
1. Fenced block whose content == preceding paragraph → collapse to
   `→ <content>` in yellow, drop fence.
2. `"double-quoted strings"` → cyan. Single quotes untouched.
3. Subtitle (italic under heading) → muted color.
4. Body indented 2 spaces under headings.

Fixtures: `gitmap/render/testdata/pretty/case-NNN-*.{in.md,want.txt}`.
Add new cases as paired files; test loop picks them up automatically.

## Spec

`spec/01-app/109-templates-ignore-attributes-pretty.md`

## Plan

`.lovable/memory/plans/04-templates-ignore-attributes-plan.md`
