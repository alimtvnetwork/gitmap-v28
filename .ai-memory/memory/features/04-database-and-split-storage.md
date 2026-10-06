# Feature Catalog 04: Database and Split Storage

- **Domain:** Three-Tier SQLite Split-DB, DevTools Cache Discovery, and Telemetry
- **Status:** Authoritative Capability Catalog

## 1. Three-Tier SQLite Split-DB Engine
- `gitmap.db`: Core repository catalog, machine settings, and persistent aliases.
- `installation.db`: Package versions, pinned releases, and installer metadata.
- `pipeline.db` / `repodb/*.db`: Transient CI/CD telemetry, build stage timings, and repository-scoped git caches.

## 2. DevTools Dynamic Cache Discovery
- Discovers and categorizes build caches (Go, npm, pnpm, Python, Cargo, OS temp).
- Boxed tree view visualization shows directory hierarchies, sizes, and safe purging options.

## 3. Pull-Error Isolated Database
- Dedicated `pull_errors.db` captures structured error events, network latency, and branches.
