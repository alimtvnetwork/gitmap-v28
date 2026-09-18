# Scanner

## Responsibility

Walk a directory tree and identify Git repositories.

## Behavior

1. Start from the given root directory.
2. Walk recursively into subdirectories.
3. Skip directories listed in `config.excludeDirs`.
4. When a `.git` directory is found, treat its parent as a repo root.
5. Do **not** descend further into a discovered repo's subdirectories.
6. Do **not** follow symlinks.

## Data Extracted Per Repo

| Field        | Source                                |
|--------------|---------------------------------------|
| Remote URL   | `git config --get remote.origin.url`  |
| Branch       | `git rev-parse --abbrev-ref HEAD`     |
| Absolute path| Filesystem path of repo root          |
| Relative path| Path relative to scan root            |

## Edge Cases

- **No remote configured:** URL fields are empty; notes = "no remote configured."
- **Detached HEAD:** Branch = "HEAD" or the commit hash.
- **Bare repos:** Skipped (only standard work-tree repos).
