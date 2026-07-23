# Versioning Strategy

The project follows Semantic Versioning (v2.63.0 current). The `release` system resolves versions using a three-tier priority: 1) Explicit CLI version argument, 2) --bump flag applied to a resolved baseline, 3) Current project version. Baseline resolution for bumps checks `.gitmap/release/latest.json` first, falling back to scanning local Git tags (`v*`) if metadata is missing. Semver normalization handles `v` prefixes and zero-padding. A critical requirement is synchronization between the compiled `Version` constant, `CHANGELOG.md`, and `.gitmap/release/` metadata.

## IMPORTANT: .gitmap/release/ Directory Policy

The `.gitmap/release/` directory should **NOT** be committed to the repository. Release metadata JSON files (`.gitmap/release/vX.Y.Z.json`, `.gitmap/release/latest.json`) are local build artifacts, not source code. They must be added to `.gitignore`.

Use `gitmap clear-release-json <version>` (alias `crj`) to remove individual release files when needed.
