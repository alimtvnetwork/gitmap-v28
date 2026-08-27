# repo-regex

Searches file contents using regex.

## Why use this instead of shell commands?
Instead of running expensive recursive regex searches:
`Get-ChildItem -Path gitmap/cmd -Filter *.go -Recurse | Select-String -Pattern "func run[A-Z]"`
or
`cat gitmap/store/db.go | Select-String "func GetDB"`

You should use the optimized lazy regex engine:
`gitmap repo-regex "func run[A-Z]"`

## Examples
```bash
gitmap repo-regex "func run[A-Z]"
gitmap repo-regex "^func Get[A-Za-z]+\("
```
