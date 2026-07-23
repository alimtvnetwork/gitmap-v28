# Release Command

## Overview

The `gitmap release` command automates Git release workflows: creating
release branches, tags, pushing to remote, and tracking release history.
It supports full semver, partial versions (auto-padded), pre-release
suffixes, draft mode, dry-run preview, and auto-increment.

---

## Commands

### `gitmap release [version]` (alias: `r`)

Create a release branch, Git tag, and push to remote.

- Version can be full (`v1.2.3`), partial (`v1.2`, `v1`), or omitted.
- Partial versions are zero-padded: `v1` → `v1.0.0`, `v1.2` → `v1.2.0`.
- If no version is provided, reads from `version.json` in the project root.
- If neither exists, exits with an error.

### `gitmap release-branch <branch>` (alias: `rb`)

Complete a release from an existing `release/vX.Y.Z` branch.
Creates the tag and pushes if not already done.

### `gitmap release-pending` (alias: `rp`)

Release all pending versions from **two sources**:

1. Local `release/v*` branches missing their `vX.Y.Z` tag.
2. `.gitmap/release/vX.Y.Z.json` metadata files where neither the branch
   nor the tag exists — uses the stored `commit` SHA to create
   branch + tag.

See [45-release-pending-metadata.md](./45-release-pending-metadata.md)
for the full metadata-based discovery spec.

---

## Flags

### Release Flags

| Flag                         | Description                                      | Default     |
|------------------------------|--------------------------------------------------|-------------|
| `--assets <path>`            | Directory or file to record as release assets    | (none)      |
| `--bin` / `-b`               | Cross-compile Go binaries and include in assets  | `false`     |
| `--commit <sha>`             | Create release from a specific commit            | (none)      |
| `--branch <name>`            | Create release from latest commit of a branch    | (none)      |
| `--bump major\|minor\|patch` | Auto-increment from the latest released version  | (none)      |
| `--notes <text>` / `-N`     | Release notes or title for the release           | (none)      |
| `--draft`                    | Mark release metadata as draft                   | `false`     |
| `--dry-run`                  | Preview release steps without executing          | `false`     |
| `--no-commit`                | Skip post-release auto-commit and push           | `false`     |
| `--yes` / `-y`               | Auto-confirm all prompts (e.g. commit)           | `false`     |
| `--verbose`                  | Write detailed debug log                         | `false`     |

### Release-Branch Flags

| Flag               | Description                         | Default |
|--------------------|-------------------------------------|---------|
| `--assets <path>`  | Directory or file to record         | (none)  |
| `--notes <text>` / `-N` | Release notes or title         | (none)  |
| `--draft`          | Mark release metadata as draft      | `false` |
| `--dry-run`        | Preview steps without executing     | `false` |
| `--no-commit`      | Skip post-release auto-commit       | `false` |
| `--yes` / `-y`     | Auto-confirm all prompts            | `false` |
| `--verbose`        | Write detailed debug log            | `false` |

### Release-Pending Flags

| Flag               | Description                              | Default |
|--------------------|------------------------------------------|---------|
| `--assets <path>`  | Directory or file to record              | (none)  |
| `--notes <text>` / `-N` | Release notes or title              | (none)  |
| `--draft`          | Mark release metadata as draft           | `false` |
| `--dry-run`        | Preview steps without executing          | `false` |
| `--no-commit`      | Skip post-release auto-commit            | `false` |
| `--yes` / `-y`     | Auto-confirm all prompts                 | `false` |
| `--verbose`        | Write detailed debug log                 | `false` |

---

## Mutual Exclusivity Rules

The following flag combinations are invalid and cause an immediate error:

| Conflict                         | Error Message                                              |
|----------------------------------|------------------------------------------------------------|
| `--bump` + version argument      | `--bump cannot be used with an explicit version argument.` |
| `--commit` + `--branch`          | `--commit and --branch are mutually exclusive.`            |

---

## Version Resolution

Version is resolved in priority order:

1. **CLI argument** — `gitmap release v1.2.3`
2. **`--bump` flag** — reads latest from `.gitmap/release/latest.json`, falls back to git tags
3. **`version.json`** — `{ "version": "1.2.3" }` in project root
4. **Error** — no version source found

### Partial Version Padding

| Input   | Resolved  | Branch              | Tag       |
|---------|-----------|---------------------|-----------|
| `v1`    | `v1.0.0`  | `release/v1.0.0`    | `v1.0.0`  |
| `v1.2`  | `v1.2.0`  | `release/v1.2.0`    | `v1.2.0`  |
| `v1.2.3`| `v1.2.3`  | `release/v1.2.3`    | `v1.2.3`  |

### Pre-Release Versions

Pre-release suffixes are preserved and not padded:

| Input           | Resolved        | Tag             |
|-----------------|-----------------|-----------------|
| `v1.0.0-rc.1`  | `v1.0.0-rc.1`  | `v1.0.0-rc.1`  |
| `v1.0.0-beta`  | `v1.0.0-beta`  | `v1.0.0-beta`  |
| `v2.0.0-alpha.3`| `v2.0.0-alpha.3`| `v2.0.0-alpha.3`|

Pre-release versions are **never** marked as `latest`.

---

## Source Resolution

The commit used as the release base is resolved in order:

1. **`--commit <sha>`** — exact commit
2. **`--branch <name>`** — latest commit on that branch
3. **Current HEAD** — default

---

## Dry-Run Mode

When `--dry-run` is passed, each step is printed with a `[dry-run]`
prefix. No branches, tags, or pushes are created.
Metadata files are not written.

```
  [dry-run] Create branch release/v1.2.3 from main
  [dry-run] Create tag v1.2.3
  [dry-run] Push branch and tag to origin
  [dry-run] Write metadata to .gitmap/release/v1.2.3.json
  [dry-run] Mark v1.2.3 as latest
```

---

## Auto-Increment (`--bump`)

Reads the latest version from `.gitmap/release/latest.json` and increments.
If `latest.json` is missing, falls back to scanning local Git tags
(`v*`) for the highest stable semver version.

| Current Latest | `--bump patch` | `--bump minor` | `--bump major` |
|----------------|----------------|----------------|----------------|
| `1.2.3`        | `1.2.4`        | `1.3.0`        | `2.0.0`        |
| `0.9.1`        | `0.9.2`        | `0.10.0`       | `1.0.0`        |

If no `latest.json` exists and no version tags are found, exits with
an error instructing the user to create an initial release first.

`--bump` is mutually exclusive with a version argument.

---

## Duplicate Detection

Before creating a release, the tool checks:

1. **`.gitmap/release/vX.Y.Z.json`** — if the metadata file exists:
   - Check if the Git tag exists (locally or remote).
   - Check if the release branch exists.
   - If **both** are missing → orphaned metadata (see below).
   - If **either** exists → abort with "already released" error.
2. **Git tags** — if the tag already exists locally or remotely, abort.

Error message:
```
Version v1.2.3 is already released. See .gitmap/release/v1.2.3.json for details.
```

### Orphaned Metadata Recovery

If a `.gitmap/release/vX.Y.Z.json` file exists but neither the Git tag nor
the release branch is found, the metadata is considered **orphaned**
(e.g. from a previously failed or manually cleaned-up release).

Instead of aborting, the tool prompts the user:

```
  ⚠ Release metadata exists for v2.3.10 but no tag or branch was found.
  → Do you want to remove the release JSON and proceed? (y/N):
```

| User Response | Behavior |
|---------------|----------|
| `y` or `yes`  | Deletes the stale `.gitmap/release/vX.Y.Z.json` file and proceeds with the normal release workflow (step 5 onward). |
| `n`, `no`, or Enter | Aborts the release with "release aborted by user". |
| EOF / no input | Aborts with the standard "already released" error. |

Detection logic:

1. Release JSON exists for the target version.
2. Git tag does not exist locally **and** does not exist on remote.
3. Release branch (`release/vX.Y.Z`) does not exist.
4. → Prompt user to remove stale JSON and proceed.

---

## Error Scenarios

| Scenario                        | Behavior                                                    |
|---------------------------------|-------------------------------------------------------------|
| Invalid version string          | `'abc' is not a valid version.`                             |
| `--commit` SHA not found        | `commit abc123 not found.`                                  |
| `--branch` does not exist       | `branch develop does not exist.`                            |
| Push to remote fails            | `failed to push to remote: <detail>`                        |
| Metadata write fails            | `could not write release metadata: <detail>`                |
| `version.json` unreadable       | `could not read version.json: <detail>`                     |

### Rollback Strategy

If a step fails after partial execution:

- **Branch/tag created but push fails**: error is reported; user must
  manually delete the local branch and tag.
- **Push succeeds but metadata write fails**: branch and tag remain on
  remote; user should manually create the `.gitmap/release/` file.

No automatic rollback is performed. The error message includes the
failed step so the user knows exactly what to clean up.

---

## version.json Behavior

- `version.json` is **read-only** from the tool's perspective.
- The tool never auto-updates `version.json` after a release.
- Users manage `version.json` manually or via their own CI scripts.

---

## Workflow: Release from HEAD / Branch / Commit

```
 1. Resolve version (CLI → --bump → version.json → error)
 2. Pad partial version to full semver
 3. Check .gitmap/release/ and git tags for duplicates
 3a. If orphaned metadata detected → prompt to remove and continue
 4. Resolve source commit (--commit / --branch / HEAD)
 5. Create branch release/vX.Y.Z
 6. Create git tag vX.Y.Z (annotated with --notes if provided)
 7. Push branch + tag to origin
 8. If --bin: cross-compile Go binaries, collect --assets contents, upload
 9. Return to original branch
10. Write .gitmap/release/vX.Y.Z.json + update latest.json on original branch
11. Auto-commit .gitmap/release/ metadata files
```

## Workflow: Release from Existing Branch

```
1. Validate release/vX.Y branch exists
2. Extract version from branch name, pad to semver
3. Check if tag already exists → abort if so
4. Checkout the release branch
5. Create tag, push, upload assets
6. Return to original branch
```

**Note:** `release-branch` and `release-pending` skip `.gitmap/release/` metadata
writing and committing. These commands process branches/metadata that
already exist — they only create the tag, push, and upload assets.

---

## Package Layout

```
release/
├── semver.go       # Version parsing, padding, comparison, bumping
├── metadata.go     # Read/write .gitmap/release/*.json, latest.json, version.json
├── gitops.go       # Branch, tag, push, checkout Git operations + git tag fallback
├── github.go       # Asset collection, changelog/readme detection
└── workflow.go     # Orchestration: Execute(), ExecuteFromBranch()
```

Each file stays under 200 lines. `workflow.go` is the entry point;
all other files are pure helpers with no cross-dependencies.

---

## CLI Examples

```bash
# Full semver release from HEAD
gitmap release v1.2.3

# Partial version (padded to v1.0.0)
gitmap release v1

# With assets
gitmap release v2.0.0 --assets ./dist

# With Go binary cross-compilation
gitmap release v2.0.0 --bin
gitmap release v2.0.0 -b --assets ./dist

# Alias
gitmap r v1.5.0

# From specific commit
gitmap release v1.2.3 --commit abc123def

# From specific branch
gitmap release v1.0.0 --branch develop

# Auto-increment
gitmap release --bump patch
gitmap release --bump minor --bin

# Draft release
gitmap release v3.0.0-rc.1 --draft

# Dry-run preview
gitmap release v1.0.0 --dry-run

# No version — reads version.json
gitmap release

# Complete release from existing release branch
gitmap release-branch release/v1.2.0
gitmap rb release/v1.2.0

# Dry-run from branch
gitmap release-branch release/v1.2.0 --dry-run

# Release all untagged release branches
gitmap release-pending
gitmap rp              # alias

# Preview pending releases
gitmap release-pending --dry-run

# Release pending with assets
gitmap release-pending --assets ./dist
```

---

## Acceptance Criteria

- **Given** `gitmap release v1.0.0`, **then** branch `release/v1.0.0`
  and tag `v1.0.0` are created and pushed.
- **Given** `--assets ./dist`, **then** dist folder contents are recorded
  in release metadata.
- **Given** version already released, **then** abort with clear message.
- **Given** no version arg and `version.json` exists, **then** version is
  read from it.
- **Given** `--commit <sha>`, **then** release branch starts from that commit.
- **Given** `--branch main`, **then** latest commit of `main` is used.
- **Given** `gitmap release-branch release/v1.2.0`, **then** tag is
  created from that branch and pushed.
- **Given** multiple releases, **then** `latest.json` points to the highest
  stable semver.
- **Given** `--bump patch` with latest `v1.2.3`, **then** releases `v1.2.4`.
- **Given** `--bump` with no `latest.json` and git tags exist, **then**
  detects highest stable tag and bumps from it.
- **Given** `--draft`, **then** release metadata is marked as draft.
- **Given** pre-release version, **then** it is not marked as latest.
- **Given** `v1`, **then** padded to `v1.0.0`.
- **Given** `--dry-run`, **then** all steps are printed but nothing is
  executed; no branches, tags, or pushes are created.
- **Given** `--bump` with a version argument, **then** abort with conflict
  error.
- **Given** `--commit` with `--branch`, **then** abort with mutual
  exclusivity error.
- **Given** push fails after branch/tag creation, **then** error message
  includes the failed step for manual cleanup.
- **Given** `gitmap release-pending`, **then** all `release/v*` branches
  without matching tags are released.
- **Given** `gitmap release-pending --dry-run`, **then** pending releases
  are listed but no tags or pushes are created.
- **Given** `.gitmap/release/vX.Y.Z.json` exists but no tag or branch, **then**
  user is prompted to remove the orphaned JSON file.
- **Given** orphaned metadata prompt answered `y`, **then** the stale
  JSON is deleted and the release proceeds normally.
- **Given** orphaned metadata prompt answered `n`, **then** the release
  is aborted with "release aborted by user".

---

## CI Release Pipeline (GitHub Actions)

The `release.yml` workflow triggers automatically when:
1. A tag matching `v*` is pushed.
2. A branch matching `release/**` is pushed.

### Steps

1. **Resolve version** — extracted from the tag name or the branch name
   (e.g., `release/v2.49.0` → `v2.49.0`).
2. **Build binaries** — cross-compiles all 6 default targets with version
   baked into the binary via `-ldflags`.
3. **Compress** — Windows binaries are zipped; Linux/macOS are tar.gz'd.
4. **Generate checksums** — SHA256 checksums for all dist files.
5. **Generate install scripts** — version-pinned `install.ps1` (Windows)
   and `install.sh` (Linux/macOS) are created and attached as release
   assets.
6. **Extract changelog** — the matching section from `CHANGELOG.md` is
   extracted for the release body.
7. **Build release body** — combines: changelog entry, release metadata
   table (version, commit, branch, build date, Go version), SHA256
   checksums block, install instructions (PowerShell and Bash
   one-liners), and platform/architecture asset matrix.
8. **Create GitHub Release** — publishes the release with all assets.
   Pre-release versions (containing `-`) are automatically marked as
   prerelease.

### Release Body Format

Each GitHub release body includes:

- **Changelog entry** for the version
- **Release info table**: version, short commit SHA, branch, build date, Go version
- **SHA256 checksums** in a code block
- **Install instructions**: PowerShell one-liner (`install.ps1`) for
  Windows and Bash one-liner (`install.sh`) for Linux/macOS
- **Asset matrix table**: platform, architecture, and filename for each binary

## Cross-References (Generic Specifications)

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) | CI trigger, stage sequence, compression, checksums, publish |
| Install scripts | [03-install-scripts.md](../07-generic-release/03-install-scripts.md) | Version-pinned `install.ps1` / `install.sh`, SHA-256 verification |
| Release metadata | [06-release-metadata.md](../07-generic-release/06-release-metadata.md) | `releases.json` manifest, `baseUrl`, asset maps |
| Release assets | [05-release-assets.md](../07-generic-release/05-release-assets.md) | Asset naming, compression, checksums |
| Cross-compilation | [01-cross-compilation.md](../07-generic-release/01-cross-compilation.md) | Multi-platform Go build targets |
| Release data model | [13-release-data-model.md](13-release-data-model.md) | Per-release metadata, `latest.json`, semver rules |
| Clone-next flatten | [87-clone-next-flatten.md](87-clone-next-flatten.md) | `--flatten` flag, version tracking in DB, RepoVersionHistory |
