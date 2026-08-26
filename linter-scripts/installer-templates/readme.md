# Installer Templates

Templates consumed by the release pipeline at tag time. Placeholders of
the form `__EmbeddedXxx__` are substituted with the pinned values for
the release being built, and the resulting files are uploaded as
release assets / committed to `main` per
[spec/16-generic-release/08-version-pinned-release-installers.md](../../spec/16-generic-release/08-version-pinned-release-installers.md).

## Files

| File | Purpose | Contract |
|------|---------|----------|
| `Status.ps1.tmpl` | Windows status emitter for update-check discovery | [spec/14-update/24-update-check-mechanism/02-status-script-json.md](../../spec/14-update/24-update-check-mechanism/02-status-script-json.md) |
| `Status.sh.tmpl`  | Unix status emitter for update-check discovery     | same |

## Placeholder Reference

| Placeholder | Example | Source |
|-------------|---------|--------|
| `__EmbeddedVersion__`          | `V1.5.0` | git tag |
| `__EmbeddedRepo__`             | `https://github.com/Owner/Repo` | repo origin |
| `__EmbeddedCommit__`           | `a1b2c3d4...` (40 chars) | `git rev-parse HEAD` |
| `__EmbeddedChecksum__`         | `Sha256:Abc123...` | release-asset checksum |
| `__EmbeddedReleaseUrl__`       | `<Repo>/releases/tag/v1.5.0` | GitHub API |
| `__EmbeddedPublishedAt__`      | `2026-04-15T10:00:00Z` | release publish time, ISO 8601 UTC |
| `__EmbeddedNotes__`            | `Short summary` (or leave for null) | release notes |
| `__EmbeddedMinSupportedFrom__` | `V1.0.0` | spec / changelog |
| `__EmbeddedNewRepoUrl__`       | `null` or `https://github.com/Owner/repo-v20` | migration manifest |

When a placeholder represents a nullable field (`Notes`, `NewRepoUrl`)
and there is no value, leave the literal `__EmbeddedXxx__` token in
place — both templates detect the unsubstituted token and emit JSON
`null` for that field. The pipeline MAY also substitute the literal
string `null`; both behaviors yield the same output.

## Output Contract

Both templates emit a single JSON document to stdout, exit 0, and
produce **no other output**. The consumer
(`UpdateCheckerService.CheckForUpdate`) parses stdout directly — any
prose or progress output would break discovery.

## Local Testing

```bash

# Substitute placeholders manually for a smoke test

sed -e 's|__EmbeddedVersion__|V1.5.0|g' \
    -e 's|__EmbeddedRepo__|https://github.com/Owner/Repo|g' \
    -e 's|__EmbeddedCommit__|abcdef0123456789abcdef0123456789abcdef01|g' \
    -e 's|__EmbeddedChecksum__|Sha256:DeadBeef|g' \
    -e 's|__EmbeddedReleaseUrl__|https://github.com/Owner/Repo/releases/tag/v1.5.0|g' \
    -e 's|__EmbeddedPublishedAt__|2026-04-15T10:00:00Z|g' \
    -e 's|__EmbeddedNotes__|Bug fixes|g' \
    -e 's|__EmbeddedMinSupportedFrom__|V1.0.0|g' \
    -e 's|__EmbeddedNewRepoUrl__|null|g' \
    Status.sh.tmpl > /tmp/Status.sh
bash /tmp/Status.sh | jq .
```

The PowerShell template can be tested similarly with
`(Get-Content Status.ps1.tmpl) -replace '__EmbeddedVersion__','V1.5.0' …`
piped to `pwsh -NoProfile -Command -`.
