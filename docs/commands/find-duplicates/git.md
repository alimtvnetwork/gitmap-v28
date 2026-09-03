# `gitmap git find-duplicates`

Audit duplicate Git repositories tracked in the GitMap database or colliding clone folders.

## Usage

```bash
gitmap git find-duplicates
gitmap find-duplicates git
```

## Remediation Commands

```text
Remediation & Fix Commands for Git Repositories:
──────────────────────────────────────────────────────────────────────────
● Fix Repeated Clone Paths:
  gitmap clone --fix

● Untrack specific duplicate repo:
  gitmap rm <duplicate-path>
```
