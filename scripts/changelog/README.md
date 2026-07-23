# scripts/changelog

Regenerates `CHANGELOG.md` and `src/data/changelog.ts` from Conventional
Commits since the most-recent annotated git tag.

## Usage

```sh
# Regenerate the two changelog sources for the next release
make changelog VERSION=v3.92.0

# CI gate: exit 3 when the on-disk files drift from regenerated output
make changelog-check VERSION=v3.92.0
```

When `VERSION` is omitted the new entry is labelled `<latest-tag>+next`
(or `vNEXT` if the repository has no tags yet).

## Conventional Commit prefixes

| Prefix      | Section in changelog |
|-------------|----------------------|
| `feat:`     | Added                |
| `fix:`      | Fixed                |
| `docs:`     | Docs                 |
| `refactor:` | Refactor             |
| `perf:`     | Performance          |
| `test:`     | Tests                |
| `build:`    | Build                |
| `ci:`       | CI                   |
| `style:`    | Style                |
| `chore:`    | Chore                |
| `revert:`   | Reverted             |

Scoped (`feat(cli):`) and breaking (`feat!:`) variants are recognised.
Commits without a recognised prefix are reported on stderr and skipped
so a single `Changes` subject cannot pollute the release notes.

## Architecture

The generator lives in its own Go module so it doesn't add dependencies
to the `gitmap` binary:

```
scripts/changelog/
├── go.mod
├── main.go                     — flag parsing entry point
└── internal/
    ├── runner/                 — mode dispatch (write | check)
    ├── gitlog/                 — `git describe` + `git log` wrapper
    ├── group/                  — Conventional-Commit prefix → section
    ├── render/                 — Markdown + TypeScript fragments
    ├── writer/                 — splice fragments into existing files
    └── drift/                  — CI drift comparator
```
