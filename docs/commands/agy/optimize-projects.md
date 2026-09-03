# `gitmap agy optimize-projects`

Deduplicate Antigravity projects sharing identical workspace directories.

## Usage

```bash
gitmap agy optimize-projects [flags]
```

## Aliases

`gitmap agy --repeat-fix`, `gitmap agy -r`

## Description

When the same repository is repeatedly scanned or imported, multiple `<id>.json` project files can accumulate pointing to the exact same path. `optimize-projects` identifies these collision groups, keeps the newest/most recently updated project file, and deletes the redundant duplicates.

## Examples

```bash
gitmap agy optimize-projects
gitmap agy optimize-projects --except "56f4c903-1658-492f-b31d-a9c25cce4c0d"
```
