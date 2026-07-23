# Release: Orphaned Metadata Recovery

When `gitmap-v27 release` detects a `.gitmap/release/vX.Y.Z.json` file but neither
the Git tag nor the release branch exists, it now prompts the user:

1. Warns that metadata exists without a matching tag/branch.
2. Asks if the user wants to remove the orphaned JSON file.
3. If confirmed, deletes the file and proceeds with the normal release flow.
4. If declined, aborts the release.

This prevents the "already released" error when the JSON file is stale.

## CRITICAL: .gitmap/release/ Directory Policy

The `.gitmap/release/` directory must **NEVER** be modified by the AI/editor.
Release metadata JSON files are local build artifacts managed exclusively
by the CLI tool. The AI must not create, edit, or delete any files in `.gitmap/release/`.
