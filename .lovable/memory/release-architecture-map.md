# Release Architecture Map

## Single Source of Truth
The canonical source of truth for the repository's version is the `Version` field in `version.json` at the root of the repository.

## Release Process
1. **Version Bump**: Update the `Version` string in `version.json` to the new semantic version. Update the `updated` timestamp.
2. **README Pinning**: Find and replace instances of the old version with the new version in the root `readme.md`.
3. **Changelog**: Prepend the release notes to `changelog.md` under the `# Changelog` header.
4. **Sub-components**: `version.json` supports sub-component sections (e.g., backend, frontend, cli). These are set to `"inherit"` to use the global root `Version`.
5. **Propagation**: Running `npm run sync` (if applicable) or committing the updated `version.json` signals to the CI/CD pipeline and AI agents that a new release is cut.
6. **Testing**: All test files (e.g., `*test*`, `*.spec.*`) are strictly excluded from automated version string scanning and replacement to prevent mock data corruption.

## Associated CI/CD
When `version.json` changes, GitHub Actions (such as `release.yml` and `goreleaser.yml`) will detect the drift, compile the binaries, and cut the Git tags and GitHub Releases automatically.
