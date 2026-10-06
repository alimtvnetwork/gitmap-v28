# Feature Catalog 02: Git Operations and Cloning

- **Domain:** Semantic Flat Commits, Pull-All Status (PAS), and Multi-Repo Cloning
- **Status:** Authoritative Capability Catalog

## 1. Semantic Flat Commit Suite (`gitmap commit` / `cm`)
- Single-command atomic commit creation with automatic dirty check.
- High-speed alias `cm` for fast execution cycles.
- Auto-stage flag (`--all` / `-a`) stages tracked and untracked changes before committing.
- Non-git directory guard preventing execution outside valid repositories.

## 2. Pull-All Status (PAS) & CPAR
- Worker concurrency scales dynamically across available CPU cores.
- Fast-forward auto-rebase eliminates merge noise across clean branches.
- Commit-Pull-Array-Runner (CPAR) sequences multi-repo pipelines with fail-fast options.

## 3. Remote Cloning & LFS Fallback
- Clone-From-Remote (CFR) and Clone-From-Remote-Parallel (CFRP) distribute repository clones across fleet nodes.
- LFS smudge fallback intercepts missing LFS server objects without halting clone operations.
