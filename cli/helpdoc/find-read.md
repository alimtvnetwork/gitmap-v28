# find-read

Finds files within the repository using exact matches or wildcards, and then reads their contents to the terminal.
Skips files larger than 300KB if they were not cached by the indexer.

## Why use this instead of shell commands?

Instead of combining file finding and cat:
`Get-ChildItem -Path . -Filter "constants.go" -Recurse | Get-Content`

You should use:
`gitmap find-read "constants.go"`

## Examples

```bash
gitmap find-read "constants.go"
gitmap find-read "*.md" -l 2
```
