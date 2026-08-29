# find

Finds files within the repository using exact matches or wildcards.
Queries the local repository SplitDB for high-speed file discovery.

## Why use this instead of shell commands?

Instead of running:
`Get-ChildItem -Path gitmap/cmd -Filter *.go -Recurse`
or
`find . -name "*.go"`

You should use:
`gitmap find "*.go"`

## Wildcard Support

- Exact: `"filename.txt"`
- Starts-with: `"*name.txt"`
- Ends-with: `"file*"`
- Contains: `"*name*"`

## Limits

Use `--limit <n>` or `-l <n>` to cap results.

## Examples

```bash
gitmap find "main.go"
gitmap find "*.go" -l 10
gitmap find "*config*"
```
