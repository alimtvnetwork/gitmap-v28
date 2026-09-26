# gitmap recreate-repo

Recreates a remote repository on GitHub from an existing local git repository,
creating an automated pre-flight backup branch, re-anchoring `origin` remote,
and synchronizing workspaces.

## Aliases

recreate

## Usage

    gitmap recreate-repo [name] [--public]

## Description

1. Verifies that the current directory is a valid git repository.
2. Creates an automatic pre-flight backup branch (`backup/recreate-<timestamp>`).
3. Uses the GitHub CLI (`gh repo create`) to initialize and publish the repository.
4. Re-anchors remote `origin` and tracks `main`.
5. Executes full workspace synchronization.

## Examples

```bash
# Recreate current repo on parent account as private
gitmap recreate-repo my-new-repo

# Recreate as public
gitmap recreate my-new-repo --public
```

