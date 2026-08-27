# find

Finds files within the repository using exact matches or wildcards.
Queries the local repository SplitDB for high-speed file discovery.

## Wildcard Support
- Exact: "filename.txt"
- Starts-with: "*name.txt"
- Ends-with: "file*"
- Contains: "*name*"

## Limits
Use --limit <n> or -l <n> to cap results.

## Examples
```bash
gitmap find "main.go"
gitmap find "*.go" -l 10
gitmap find "*config*"
```
