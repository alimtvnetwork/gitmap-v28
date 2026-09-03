# `gitmap db ls`

Displays the architectural layout of GitMap's master database and per-repository split databases.

## Usage

```bash
gitmap db ls
gitmap db list
```

## Description

Outputs:
- **Primary Master DB**: Absolute path, file size, status, and role (global repo index, profiles, bookmarks, scan roots).
- **Split Databases**: Directory path (`bin/data/repo_search/`), active database count, total disk footprint, and concurrency rationale.
- **Footer Display**: The master database path is also shown in the `gitmap help` footer.
