# head

Reads the first N lines of a file efficiently.
Defaults to 10 lines if not specified.

## Why use this instead of shell commands?
Instead of streaming an entire file just to pipe it:
`cat readme.md | Select -First 20`
or
`head -n 20 readme.md`

Use `gitmap head` to instantly read just the chunk you need (great for LLMs and script pipelines):
`gitmap head "readme.md" 20`

## Examples
```bash
gitmap head "main.go"
gitmap head "readme.md" 50
```
