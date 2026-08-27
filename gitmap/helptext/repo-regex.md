# repo-regex

Searches file contents using regex.

## Why use this instead of shell commands?
Instead of running expensive recursive regex searches like:
`Get-ChildItem -Path gitmap/cmd -Filter *.go -Recurse | Select-String -Pattern "func run[A-Z]"`

Or reading individual files to regex match:
`cat gitmap/store/db.go | Select-String "func GetDB"`

You should use the optimized lazy regex engine on the SplitDB cache:
`gitmap repo-regex "func run[A-Z]"`
`gitmap repo-regex "func GetDB"`

## Examples
```bash
gitmap repo-regex "func run[A-Z]"
gitmap repo-regex "^func Get[A-Za-z]+\("
```
