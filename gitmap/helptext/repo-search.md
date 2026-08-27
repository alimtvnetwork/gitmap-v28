# repo-search

Searches file contents using the cached DB.

## Why use this instead of shell commands?
Instead of running an expensive recursive search:
`grep -rn "func run" gitmap/cmd/`
or
`Get-ChildItem -Path gitmap/cmd -Filter *.go -Recurse | Select-String -Pattern "func run"`

You should use the SplitDB cache for instant results:
`gitmap repo-search "func run"`

## Examples
```bash
gitmap repo-search "func run"
gitmap repo-search "GetDB" -l 10
```
