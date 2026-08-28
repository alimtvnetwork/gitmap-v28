# write

Writes text to a file, overwriting existing content.
Creates the file and parent directories if they don't exist.

## Why use this instead of shell commands?
Instead of running:
`Set-Content config.txt "value"`
or:
`echo "content" > file.txt`

Use `gitmap write` for cross-platform consistency:
`gitmap write "config.txt" "value"`

## Examples
```bash
gitmap write "hello.txt" "Hello World"
gitmap write "config/settings.json" "{}"
```
