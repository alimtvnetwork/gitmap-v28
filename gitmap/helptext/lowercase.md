# lowercase

Mass-rename files matching a pattern to lowercase while preserving Git history via `git mv`.

## Examples

Convert all "OLD.md" files to "old.md", skipping node_modules and .git:
```bash
gitmap lowercase "OLD.md" "old.md" . -except "node_modules/*,.git/*"
```

Same as above, using the default ignore profile (which automatically skips volatile dirs):
```bash
gitmap lowercase "OLD.md" "old.md" -ignore default
```

## OVERVIEW

This command recursively scans the target directories and renames files to the specified lowercase target. If the file is tracked in Git, `git mv` is used to maintain file history.
