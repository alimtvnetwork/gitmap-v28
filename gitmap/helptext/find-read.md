# find-read

Finds files within the repository using exact matches or wildcards, and then reads their contents to the terminal.
Skips files larger than 300KB if they were not cached by the indexer.

## Examples
```bash
gitmap find-read "constants.go"
gitmap find-read "*.md" -l 2
```
