# gitmap release-pending

Release all pending versions from two sources: local `release/v*`
branches missing tags, and `.gitmap/release/v*.json` metadata files where
neither the branch nor the tag exists.

## Alias

rp

## Usage

    gitmap release-pending [--assets <path>] [--notes "text"] [--draft] [--dry-run] [--verbose]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --assets | (none) | Directory or file to attach to each release |
| --notes, -N | (none) | Release notes or title for the release |
| --draft | false | Mark releases as drafts |
| --dry-run | false | Preview without executing |
| --verbose | false | Detailed output |
| --no-commit | false | Skip post-release auto-commit and push |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Release pending branches

    gitmap release-pending

**Output:**

    ■ Scanning for pending releases...
    Found 2 pending release branch(es):
      release/v2.20.0 (untagged)
      release/v2.21.0 (untagged)

    [1/2] Releasing v2.20.0...
      Creating tag v2.20.0... done
      Pushing tag... done
      ✓ Release v2.20.0 complete.

    [2/2] Releasing v2.21.0...
      Creating tag v2.21.0... done
      Pushing tag... done
      ✓ Release v2.21.0 complete.

    ✓ 2 pending releases completed

### Example 2: Release from orphaned metadata

    gitmap rp

**Output:**

    ■ Scanning for pending releases...
    Found 0 pending release branch(es).
    → Found 1 unreleased version(s) from .gitmap/release/ metadata:
      v2.19.0 (commit: abc1234, no branch or tag found)

    [1/1] Creating release from metadata: v2.19.0
      Creating branch release/v2.19.0 from commit abc1234... done
      Creating tag v2.19.0... done
      Pushing branch and tag... done
      ✓ Release v2.19.0 complete.

### Example 3: Dry-run preview

    gitmap rp --dry-run

**Output:**

    [DRY RUN] Scanning for pending releases...
    [DRY RUN] Found 2 pending release branch(es):
      release/v2.20.0 (would create tag v2.20.0)
      release/v2.21.0 (would create tag v2.21.0)
    [DRY RUN] Found 1 unreleased from .gitmap/release/ metadata:
      v2.19.0 (would create branch + tag from abc1234)
    No changes made.

### Example 4: Release pending as drafts with assets

    gitmap release-pending --draft --assets ./dist/

**Output:**

    Found 1 pending release branch(es):
      release/v2.22.0 (untagged)
    [1/1] Releasing v2.22.0 (draft)...
      Creating tag v2.22.0... done
      Attaching assets from ./dist/... done
      ✓ Draft release v2.22.0 created

## See Also

- [release](release.md) — Create a release
- [release-branch](release-branch.md) — Complete from existing branch
- [changelog](changelog.md) — View release notes

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter release-pending
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
