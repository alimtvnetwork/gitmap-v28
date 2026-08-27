# find-regex

Finds files within the repository by evaluating a regex pattern against the relative file paths.
Note: Regex searches only evaluate filenames, not file contents.

## Limits
Use --limit <n> or -l <n> to cap results.

## Examples
```bash
gitmap find-regex "^cmd/.*\.go$"
gitmap find-regex "v[0-9]+.*\.json$" -l 50
```
