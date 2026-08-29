# append

Appends text to a file.
Creates the file and parent directories if they don't exist.

## Why use this instead of shell commands?

Instead of running:
`Add-Content .lovable/cicd-issues/index.md "- [18-ci.md](...)"`
or:
`echo "content" >> file.txt`

Use `gitmap append` for cross-platform consistency:
`gitmap append ".lovable/cicd-issues/index.md" "- [18-ci.md](...)"`

## Examples

```bash
gitmap append "hello.txt" "Hello World"
gitmap append "config/settings.json" "{}"
```
