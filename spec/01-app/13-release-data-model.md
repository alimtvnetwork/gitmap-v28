# Release Data Model

## Version Source: `version.json`

Located at the project root. Used as fallback when no version argument
is provided to `gitmap release`.

```json
{
  "version": "1.2.3"
}
```

| Field   | Type   | Required | Notes                    |
|---------|--------|----------|--------------------------|
| version | string | yes      | Semver without `v` prefix |

### Behavior

- **Read-only**: the tool never modifies `version.json`.
- Users update it manually or via CI scripts between releases.
- If the file is missing or unreadable and no other version source
  exists, the release command exits with an error.

---

## Release Metadata: `.gitmap/release/vX.Y.Z.json`

One file per release, stored in the `.gitmap/release/` directory at the
project root. Created after a successful release.

```json
{
  "version": "1.2.3",
  "branch": "release/v1.2.3",
  "sourceBranch": "main",
  "commit": "abc123def456789",
  "tag": "v1.2.3",
  "assets": ["./dist"],
  "zipGroups": ["docs-bundle"],
  "draft": false,
  "preRelease": false,
  "createdAt": "2026-03-05T12:00:00Z",
  "isLatest": true
}
```

| Field        | Type     | Required | Default | Notes                              |
|--------------|----------|----------|---------|------------------------------------|
| version      | string   | yes      | —       | Full semver (padded)               |
| branch       | string   | yes      | —       | Release branch name                |
| sourceBranch | string   | yes      | —       | Branch or ref the release was based on |
| commit       | string   | yes      | —       | Full commit SHA                    |
| tag          | string   | yes      | —       | Git tag name with `v` prefix       |
| assets       | []string | no       | `[]`    | Paths that were attached           |
| zipGroups    | []string | no       | `[]`    | Zip group names included in release |
| draft        | bool     | no       | `false` | Whether this is a draft release    |
| preRelease   | bool     | no       | `false` | Whether this is a pre-release      |
| createdAt    | string   | yes      | —       | ISO 8601 timestamp                 |
| isLatest     | bool     | yes      | —       | Whether this is the latest stable  |

---

## Latest Pointer: `.gitmap/release/latest.json`

Always points to the highest **stable** (non-pre-release) version.
Updated after each stable release.

```json
{
  "version": "1.2.3",
  "tag": "v1.2.3",
  "branch": "release/v1.2.3"
}
```

| Field   | Type   | Required | Notes                           |
|---------|--------|----------|---------------------------------|
| version | string | yes      | Highest stable semver           |
| tag     | string | yes      | Tag name of the latest release  |
| branch  | string | yes      | Branch name of the latest release |

---

## `.gitmap/release/` Directory Structure

```
.gitmap/release/
├── latest.json
├── v1.0.0.json
├── v1.1.0.json
├── v1.2.0.json
├── v1.2.1.json
└── v2.0.0-rc.1.json
```

### Rules

- One JSON file per release, named by the full padded version.
- `latest.json` is updated only for stable releases.
- Pre-release entries set `preRelease: true` and `isLatest: false`.
- Draft entries set `draft: true`.
- The `isLatest` flag on previous releases is **not** retroactively
  updated — only `latest.json` is the source of truth for which
  version is current.

### Git Tracking

The `.gitmap/release/` directory should **NOT** be committed to the repository.
Release metadata JSON files are local build artifacts, not source code.
Add `.gitmap/release/` to `.gitignore`. Use `gitmap clear-release-json <version>`
to remove individual release files when needed.

---

## Semver Parsing

### Padding Rules

| Input           | Parsed Version   | Pre-Release | Valid |
|-----------------|------------------|-------------|-------|
| `v1`            | `1.0.0`          | no          | yes   |
| `v1.2`          | `1.2.0`          | no          | yes   |
| `v1.2.3`        | `1.2.3`          | no          | yes   |
| `v1.0.0-rc.1`  | `1.0.0-rc.1`     | yes         | yes   |
| `v1.0.0-beta`  | `1.0.0-beta`     | yes         | yes   |
| `abc`           | —                | —           | no    |
| `v`             | —                | —           | no    |

### Comparison

Versions are compared using standard semver ordering:

1. Major > Minor > Patch (numeric comparison)
2. Pre-release versions are always lower than the same stable version
   (`1.0.0-rc.1 < 1.0.0`)
3. `latest.json` is updated only when the new version is strictly
   greater than the current latest stable version.

---

## Package Layout

The release data model is implemented in `release/metadata.go`:

| Function          | Responsibility                                    |
|-------------------|---------------------------------------------------|
| `ReadVersionFile` | Parse `version.json`, return raw version string   |
| `ReadLatest`      | Parse `.gitmap/release/latest.json`                      |
| `WriteLatest`     | Update `.gitmap/release/latest.json` for stable releases |
| `WriteReleaseMeta`| Write `.gitmap/release/vX.Y.Z.json`                      |
| `ReleaseExists`   | Check if `.gitmap/release/vX.Y.Z.json` already exists    |

Cleanup lives in `cmd/clearreleasejson.go`:

| Function              | Responsibility                                    |
|-----------------------|---------------------------------------------------|
| `runClearReleaseJSON` | Remove a `.gitmap/release/vX.Y.Z.json` file by version   |

Version parsing lives in `release/semver.go`:

| Function      | Responsibility                                  |
|---------------|--------------------------------------------------|
| `Parse`       | Parse and pad a version string to full semver    |
| `Bump`        | Increment major, minor, or patch                 |
| `GreaterThan` | Compare two versions for ordering                |

## Cross-References (Generic Specifications)

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Release metadata | [06-release-metadata.md](../07-generic-release/06-release-metadata.md) | `releases.json` manifest, `baseUrl`, version tag verification |
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) | CI trigger, stage sequence, prerelease detection |
| Release command | [12-release-command.md](12-release-command.md) | CLI commands, version resolution, pending workflows |
