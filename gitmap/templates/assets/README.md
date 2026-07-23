# Embedded template assets

This directory is embedded into the binary via `//go:embed all:assets` in
`gitmap/templates/embed.go`. Phase 1 of plan 04 fills it with curated
language templates.

Layout (Phase 1):

```
assets/
  ignore/
    common.gitignore
    go.gitignore
    node.gitignore
    python.gitignore
    rust.gitignore
    csharp.gitignore
  attributes/
    common.gitattributes
    go.gitattributes
    ...
  lfs/
    common.gitattributes
```

Every template file MUST start with the audit-trail header:

```
# source: <upstream-or-curated>
# kind: ignore | attributes | lfs
# lang: <lang>
# version: 1
```

Bump `version:` when re-curating from upstream.
