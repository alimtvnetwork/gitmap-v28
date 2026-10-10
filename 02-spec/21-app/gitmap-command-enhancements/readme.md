# GitMap Command Enhancements Specification Suite (269 / gitmap-command-enhancements)

## Document Catalog

1. [`01-architecture-spec.md`](./01-architecture-spec.md): Complete architecture specification for `summary`, `full summary` (`fs`), `fs+pe`, `pe all`, and `merge-ai` (`ma`).
2. [`repo-feature.md`](./repo-feature.md): Universal Destination Resolver specification ("Repo Feature") handling remote URLs, git folders, non-git folders, and bare repo slugs.

## Command Reference Summary

- `gitmap summary [$repo] [N]`: Summarizes last $N$ releases (default $N=8$, $\le 200$ words gist, top 5 heated churn files, Split-DB `summary.db` cache).
- `gitmap full summary [N]` / `gitmap full status [N]` / `gitmap fs [N]`: Workspace heatmap tree view (default $N=3$, 48h active or dirty filter, dirty file list, master sanitize command).
- `gitmap full summary+pe` / `gitmap fs+pe [--json]` / `gitmap fspe`: CI/CD error telemetry fusion into the workspace tree.
- `gitmap pe all [--json]` / `gitmap pe all --force-all`: High-speed CI/CD pipeline diagnosis with 48h activity filter and concise green status.
- `gitmap nodes fs`, `nodes fspe`, `nodes summary $repo`, `nodes pe all`: Fleet SSH delegation with local-machine precedence deduplication.
- `gitmap merge-ai <dest> <sources>` / `gitmap ma config.json`: Multi-repo consolidation with single-commit staging, chronological inspection, `01_`/`02_` collision sequencing, `merge-ai-manifest.json`, and root `instruction.md` consolidation checklist.
