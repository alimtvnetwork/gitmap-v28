# GitMap Repo-Cache Scan, Clone, and LS Ecosystem Specification Suite

- **Directory:** `02-spec/21-app/gitmap-repo-cache-scan-clone-and-ls-ecosystem/`
- **Slug:** `gitmap-repo-cache-scan-clone-and-ls-ecosystem`
- **Specification Version:** 1.0.0
- **Status:** APPROVED

## Documents
1. [01-architecture-spec.md](01-architecture-spec.md) — System architecture, `repo-cache` companion storage resolution, auto-merge deduplication engine, sequential manifest allocator (`--separate`/`--sep`/`--s`), and verification gates.
2. [02-component-and-cli-spec.md](02-component-and-cli-spec.md) — Component contracts, CLI routing for `clone`/`cfr`/`cfrp` with `--rc`/`rc`, manifest search hierarchy, SQLite Split-DB cache (`RepoCacheManifest`, `RepoCacheEntry`), short `owner/repo` TreeView, `-y` auto-confirm, `ls rc` protocol inspector (`[SSH]` vs `[Public HTTPS]`), and copy-pasteable command hints.
