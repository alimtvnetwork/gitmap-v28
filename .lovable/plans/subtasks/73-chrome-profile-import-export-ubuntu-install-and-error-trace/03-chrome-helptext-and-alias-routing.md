# Subtask 03: Chrome CLI Helptext & Alias Routing Parity

## Scope
- Author missing help markdown files in `gitmap/helptext/`:
  - `import-all.md`: Detailed help and examples for `gitmap chrome import-all`.
  - `export-all.md`: Detailed help and examples for `gitmap chrome export-all`.
  - `copy-all.md`: Detailed help and examples for `gitmap chrome copy-all`.
  - `import-check.md`: Detailed help and examples for `gitmap chrome import-check`.
- Update `gitmap/helptext/print.go`:
  - Add alias resolution mapping for `import-all`, `export-all`, `copy-all`, `import-check`, and `import-ls`.
  - Ensure all markdown files contain `## Examples` with fenced code blocks satisfying AST tests.

## Files Touched
- `gitmap/helptext/import-all.md`
- `gitmap/helptext/export-all.md`
- `gitmap/helptext/copy-all.md`
- `gitmap/helptext/import-check.md`
- `gitmap/helptext/print.go`
