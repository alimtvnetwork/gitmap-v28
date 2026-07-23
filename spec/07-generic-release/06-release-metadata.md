# 06 — Release Metadata

## Purpose

Define how version numbers are resolved, tags are created, and
changelog entries are extracted for release automation.

---

## Version Resolution

The version is resolved with a 3-tier priority:

| Priority | Source | Example |
|----------|--------|---------|
| 1 | Explicit CLI argument | `release v1.2.0` |
| 2 | Bump flag | `release --bump minor` (1.1.0 → 1.2.0) |
| 3 | Current version from source | Read from `constants.go` or `version.go` |

### Normalization

All versions are normalized to ensure consistency:

```
1.2.0   → v1.2.0     (auto-prefix v)
v1.2.0  → v1.2.0     (already correct)
v01.2.0 → v1.2.0     (strip zero-padding)
```

### Semantic Versioning

Versions follow [SemVer 2.0.0](https://semver.org):

```
v<major>.<minor>.<patch>[-<prerelease>]
```

| Bump | When | Example |
|------|------|---------|
| Major | Breaking changes | `v1.0.0` → `v2.0.0` |
| Minor | New features (backward-compatible) | `v1.1.0` → `v1.2.0` |
| Patch | Bug fixes only | `v1.2.0` → `v1.2.1` |

---

## Tagging

After version resolution, the release system:

1. **Verifies** the tag does not already exist locally or remotely.
2. **Creates** a lightweight Git tag at `HEAD`.
3. **Pushes** the tag to the remote.

```bash
git tag "v${VERSION}"
git push origin "v${VERSION}"
```

If the push fails, the local tag must be deleted to prevent stale state:

```bash
git tag -d "v${VERSION}"
```

---

## Changelog Extraction

The release pipeline extracts the relevant section from `CHANGELOG.md`
for the release body:

```bash
# Extract everything between ## v1.2.0 and the next ## heading
awk '/^## v1\.2\.0/{found=1; next} /^## v/{if(found) exit} found{print}' CHANGELOG.md
```

### Changelog Format

```markdown
## v1.2.0 — Feature Title (2026-04-08)

### Improvements

- Added feature X for better performance.
- Updated Y to support Z.

### Bug Fixes

- Fixed crash when input is empty.
```

### Synchronization Requirement

Three sources must always be in sync:

| Source | Location | Purpose |
|--------|----------|---------|
| `Version` constant | Source code (`constants.go`) | Compiled into binary |
| `CHANGELOG.md` | Repository root | Human-readable history |
| Release metadata | `.release/latest.json` or tags | CI/CD and tooling |

When bumping a version:
1. Update the `Version` constant in source code.
2. Add the new section to `CHANGELOG.md`.
3. Update any metadata files (e.g., `latest.json`).

All three changes must happen in the **same commit** that is tagged.

---

## Release Metadata Files

Optionally maintain a `latest.json` for programmatic version queries:

```json
{
  "version": "1.2.0",
  "tag": "v1.2.0",
  "date": "2026-04-08",
  "commit": "abc123def456"
}
```

This file enables:
- Build scripts to detect the current version without parsing Go source.
- Update commands to check for newer versions.
- CI to verify version synchronization.

---

## Release Branch Strategy

For release preparation:

```
main → release/1.2.0 → tag v1.2.0 → merge back to main
```

1. Create `release/1.2.0` branch from `main`.
2. Bump version, update changelog on the branch.
3. Push the branch — CI runs tests.
4. Tag the branch head as `v1.2.0`.
5. CI publishes the release.
6. Merge back to `main`.

### Rollback

If the release push fails:
1. Switch back to the original branch.
2. Force-delete the local release branch.
3. Delete the local tag.

```bash
git checkout main
git branch -D release/1.2.0
git tag -d v1.2.0
```

---

## Constraints

- Version must be bumped in source code **before** tagging.
- Tags must be lightweight (not annotated) unless signing is required.
- Changelog, version constant, and metadata must be updated atomically.
- Pre-release versions use a hyphen suffix: `v1.2.0-beta.1`.
- Version resolution must never fall back to "latest tag" without
  explicit verification that the tag matches the source `Version`.

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | `releases.json` manifest and `baseUrl` asset maps |
| [13-release-data-model.md](../01-app/13-release-data-model.md) | Per-release metadata, `latest.json`, semver rules |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
