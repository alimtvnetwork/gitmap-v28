# `gitmap db sizes list`

Displays disk space consumed across all GitMap database and cache files.

## Usage

```bash
gitmap db sizes list
gitmap db sizes ls
```

## Description

Outputs a consolidated breakdown of:
1. Master SQLite file (`gitmap.db`)
2. Write-Ahead Log (`gitmap.db-wal`) and Shared Memory (`gitmap.db-shm`)
3. Total combined size of all split search databases in `repo_search/`
4. Release metadata cache files in `.gitmap/release/`
