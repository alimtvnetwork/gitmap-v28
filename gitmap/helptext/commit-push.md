# commit-push

Stages all changes, commits with the given message, and pushes to the remote in one command.

## Why use this instead of shell commands?

Instead of running three separate commands:

```bash
git add -A
git commit -m "your message"
git push
```

You should use:

```bash
gitmap commit-push "your message"
```

## Aliases

- `gitmap cp "message"`

## Examples

```bash
gitmap commit-push "fix: resolve null pointer in scanner"
gitmap cp "docs: update readme with new install steps"
gitmap commit-push "refactor: extract shared DB resolver to cmd_db.go"
```
