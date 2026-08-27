# rm-git

Removes a commit from git history by its SHA (or last 4+ digits of the SHA). Uses `git rebase --onto` to surgically drop the commit.

## Why use this instead of shell commands?

Instead of running a complex interactive rebase:

```bash
git log --oneline -20          # find the commit
git rebase -i <sha>^           # manually edit and drop
```

You should use:

```bash
gitmap rm-git <last-4-digits-of-sha>
```

## Aliases

- `gitmap rg <sha-fragment>`

## Examples

```bash
gitmap rm-git a5f2
gitmap rg 5a82
gitmap rm-git 022575f
```

## Warning

This command rewrites git history. After running it, you will need to force push:

```bash
git push --force-with-lease
```
