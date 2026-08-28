# search

Searches repository file contents on-the-fly. This performs an immediate filesystem walk instead of querying the cache.

## Why use this instead of shell commands?
Instead of running complex pipelines like:
`grep -rn "func run" gitmap/cmd/`
`Get-ChildItem -Path gitmap/cmd -Filter *.go | Select-String "func dispatch[A-Z]"`
or:
`cat .github/scripts/smoke-installer.ps1 | Select-String "constants.go" -Context 3,3`

Use `gitmap` to instantly find references across your codebase:
`gitmap search "constants.go"`
`gitmap search "func run"`

## Examples
```bash
gitmap search "constants.go"
gitmap search "func run"
```
