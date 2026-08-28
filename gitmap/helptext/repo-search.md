# repo-search

Searches file contents using the cached DB.

## Why use this instead of shell commands?
Instead of running an expensive recursive search or complex pipeline:
`Get-ChildItem -Path gitmap/store -Filter *.go | Select-String "OpenDefault"`
`Get-ChildItem -Path gitmap -Recurse -File | Select-String "type SearchResult struct"`
or:
`cat .github/scripts/smoke-installer.ps1 | Select-String "constants.go" -Context 3,3`

Use the high-speed cached DB to instantly find references:
`gitmap repo-search "constants.go"`
`gitmap repo-search "OpenDefault"`

## Examples
```bash
gitmap repo-search "constants.go"
gitmap repo-search "OpenDefault"
gitmap repo-search "GetDB" -l 10
```
