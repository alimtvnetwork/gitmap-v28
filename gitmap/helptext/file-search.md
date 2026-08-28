# file-search

Searches for a regex pattern within a specific file and provides contextual lines.
Queries the local repository SplitDB for high-speed file discovery and caches results.

## Why use this instead of shell commands?
Instead of running:
`cat gitmap/cmd/search_entry.go | Select-String "cmd.Flags" -Context 0,5`
or:
`grep -A 5 -B 0 "cmd.Flags" gitmap/cmd/search_entry.go`

Use `gitmap file-search` to leverage the SQLite repo database for instant, cached results:
`gitmap file-search "gitmap/cmd/search_entry.go" "cmd.Flags" 0 5`

## Usage
`gitmap file-search <file> <regex> [contextBefore] [contextAfter]`

## Examples
```bash
gitmap file-search "gitmap/cmd/search_entry.go" "cmd.Flags" 0 5
gitmap file-search "readme.md" "Install" 2 2
```
