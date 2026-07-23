# Post-Mortem: Release Pipeline `dist` Directory Error

## Issue

The CI release pipeline (`release.yml`) failed with:

```
cd: dist: No such file or directory
```

## Root Cause

The compress/checksum step ran inside `gitmap-updater/` (which has no `dist/` folder) instead of `gitmap-v28/dist/` where cross-compiled binaries are output. The `cd dist` command assumed the working directory was `gitmap-v28/`, but GitHub Actions defaults to the repository root.

## Fix

Extracted compress and checksum into a separate step with an explicit `working-directory` directive:

```yaml
- name: Compress and checksum
  working-directory: gitmap-v28/dist
  run: |
    for f in gitmap-*; do
      ...
    done
```

## Prevention Rules

1. **Never use `cd` in CI scripts** — use `working-directory` in the workflow step definition.
2. **Always validate directory existence** before operating on build outputs: `test -d "$DIR" || exit 1`.
3. **Use absolute or explicitly anchored paths** for all artifact operations in multi-project monorepos.
4. **Test pipeline changes on a `release/test-*` branch** before merging to `main`.

## Related

- `spec/05-coding-guidelines/17-cicd-patterns.md` — CI/CD Patterns
- `spec/05-coding-guidelines/29-ci-sha-deduplication.md` — SHA Deduplication
- CHANGELOG.md v2.54.0 — Release Pipeline Fix
