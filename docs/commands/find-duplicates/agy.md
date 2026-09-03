# `gitmap agy find-duplicates`

Audit duplicate Antigravity project JSON records and obtain instant remediation commands.

## Usage

```bash
gitmap agy find-duplicates
gitmap find-duplicates agy
```

## Remediation Output

Beneath the duplicate findings, the CLI outputs ready-to-run copyable commands:

```text
Remediation & Fix Commands for Antigravity:
──────────────────────────────────────────────────────────────────────────
● Fix Single (Delete specific duplicate project ID):
  gitmap agy rm <project-id>

● Fix All Together (Deduplicate & keep newest per path):
  gitmap agy optimize-projects
```
