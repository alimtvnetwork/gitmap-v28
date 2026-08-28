# tail

Reads the last N lines of a file efficiently.
Defaults to 10 lines if not specified.

## Why use this instead of shell commands?
Instead of running:
`cat readme.md | Select -Last 20`
or
`tail -n 20 readme.md`

Use `gitmap tail` to instantly read just the end of a file:
`gitmap tail "readme.md" 20`

## Examples
```bash
gitmap tail "server.log"
gitmap tail "readme.md" 50
```
