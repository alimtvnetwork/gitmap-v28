# 369-heatmap-command-enhancements-and-high-speed-sync

Master index for Heatmap Auto-Read/Touch CLI enhancements, high-speed ModTime-first fleet synchronization, and skill updates.

## Specifications
- [`01-architecture-and-heatmap-reader.md`](01-architecture-and-heatmap-reader.md) — Heatmap auto-read and auto-touch without UDF flags.
- [`02-high-speed-modtime-sync-engine.md`](02-high-speed-modtime-sync-engine.md) — High-speed Stat and ModTime-first 3-way synchronization algorithm.
- [`03-skills-and-guidelines-update.md`](03-skills-and-guidelines-update.md) — Skill updates and fleet synchronization.

## Verification
- Unit and integration tests in `cli/heatmap/` and `cli/cmdsync/`.
- Live CLI verification of `gitmap heatmap <file>` and `gitmap cg sync`.
