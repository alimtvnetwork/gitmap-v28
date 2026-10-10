# gitmap fix release tags

Audit, identify, and remove orphan or broken GitHub releases and Git release tags that lack binary assets, have unpublished release notes, or encountered failed CI/CD pipelines.

## Aliases

- `gitmap fix-release-tags`
- `gitmap frt`
- `gitmap fix release-tags`
- `gitmap fix releasetags`

## Usage

```bash
gitmap fix release tags [flags]
gitmap fix-release-tags [flags]
gitmap frt [flags]
```

## Description

Automated continuous integration and release workflows occasionally create a Git tag before triggering the binary compilation and release publishing steps. If the build runner runs out of disk space, encounters compilation errors, or network limits prevent artifact uploads, the repository is left with:
- An orphan Git tag pointing to an unreleased commit.
- A GitHub release with zero downloadable binary assets.
- A failed CI/CD workflow run attached to the release tag.

`gitmap fix release tags` audits all tags against GitHub Releases API and GitHub Actions workflow runs. If a tag is identified as broken or orphaned, the command securely deletes:
1. The GitHub Release entity (via `gh release delete` or REST API).
2. The remote Git tag on `origin` (`git push origin :refs/tags/<tag>`).
3. The local Git tag (`git tag -d <tag>`).

## Options & Flags

| Flag | Shorthand | Description | Default |
| :--- | :---: | :--- | :--- |
| `--dry-run` | `-n` | Preview detected broken releases and tags without deleting. | `false` |
| `--yes`, `--confirm` | `-y` | Non-interactive execution; bypasses the `[y/N]` confirmation prompt. | `false` |
| `--json` | - | Emit machine-readable JSON output (disables ASCII tables). | `false` |
| `--local-only` | - | Audit and remove only local Git tags; skips remote/GitHub operations. | `false` |
| `--remote-only` | - | Audit and remove only remote origin tags and GitHub releases. | `false` |
| `--repo` | `-r` | Target repository directory or repository name. | `.` |
| `--verbose` | `-v` | Display raw GitHub API responses and Git command outputs. | `false` |
| `--help` | `-h` | Display the two-column interactive help card. | `false` |

## Safety Invariants

The audit engine enforces three strict safety gates:
1. **Healthy Release Immunity:** Tags that possess at least one valid binary asset and a passing/neutral CI/CD status are never deleted.
2. **Grace Window:** Tags created within the last 30 minutes with active "in-progress" or "queued" CI/CD runs are skipped.
3. **Active Binary Protection:** If the current running `gitmap` binary matches the tag version, deletion is blocked.

## Examples

### Example 1: Dry-Run Inspection

Preview all orphan release tags in the current repository without making modifications:

```bash
gitmap fix release tags --dry-run
```

### Example 2: Interactive Remediation

Review candidate tags in a table and confirm deletion when prompted:

```bash
gitmap fix release tags
```

### Example 3: Non-Interactive CI Automation

Clean up failed release tags in CI/CD runners or headless batch scripts:

```bash
gitmap fix-release-tags --confirm
```

### Example 4: JSON Output for External Tooling

Generate structured JSON telemetry for pipeline integrations:

```bash
gitmap frt --dry-run --json
```

## JSON Schema Output Example

```json
{
  "repository": "alimtvnetwork/gitmap-v28",
  "auditTimestamp": "2026-10-10T06:30:00Z",
  "isDryRun": true,
  "candidates": [
    {
      "tag": "v6.528.1",
      "commit": "7e89ab1",
      "releaseStatus": "missing",
      "assetCount": 0,
      "cicdStatus": "failure",
      "action": "delete",
      "reasons": [
        "Release entity missing on GitHub",
        "CI/CD workflow run 1928374 failed"
      ]
    }
  ],
  "summary": {
    "totalAudited": 142,
    "brokenFound": 1,
    "deletedCount": 0
  }
}
```

## See Also

- `gitmap release` — Perform semantic release ceremony and push tags.
- `gitmap list-versions` — Inspect version tags and published metadata.
- `gitmap fix-git` — Self-heal corrupted Git index, permissions, and lockfiles.
