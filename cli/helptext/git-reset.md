# git-reset

Resets the current git branch to a target commit SHA (or SHA prefix), making that exact commit the top of your branch history and dropping all subsequent commits. Automatically synchronizes remote tracking via force-push with lease.

## Why use this instead of shell commands?

Instead of manually running multiple destructive commands:

```bash
git reset --hard <sha>
git push origin <branch> --force-with-lease
```

You should use:

```bash
gitmap git-reset <target-sha>
```

## Aliases

- `gitmap reset-git <target-sha>`

## Options

- `--no-push` / `--local`: Reset the branch locally only without force-pushing to the remote tracking branch.

## Examples

```bash
gitmap git-reset a5f2
gitmap reset-git 055422ab2b8a
gitmap gr cf82638
gitmap git-reset 055422ab --no-push
```

## Warning

This command rewrites git history. All commits after the specified SHA will be dropped from the branch.
