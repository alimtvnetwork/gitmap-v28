# search

Searches repository file contents on-the-fly. This performs an immediate filesystem walk instead of querying the cache.

## Why use this instead of shell commands?
Instead of running:
`grep -rn "func run" gitmap/cmd/`

Use:
`gitmap search "func run"`

## Examples
```bash
gitmap search "func run"
```
