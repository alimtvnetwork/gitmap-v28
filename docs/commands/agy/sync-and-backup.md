# Antigravity Synchronization & Backup

Synchronize repositories with Antigravity and backup project workspace configurations.

## Commands

### `gitmap agy sync [path]`
Synchronize all tracked GitMap repositories into Antigravity workspace JSON records.
```bash
gitmap agy sync
```

### `gitmap agy export-projects [file]`
Exports all Antigravity project JSON files into a portable ZIP backup archive.
```bash
gitmap agy export-projects agy_backup.zip
```

### `gitmap agy import-projects <file>`
Restores Antigravity project workspace records from a previously exported ZIP archive.
```bash
gitmap agy import-projects agy_backup.zip
```
