# gitmap release-self

Release gitmap itself from any directory.

## Alias

rs, rself

## Usage

    gitmap release-self [version] [flags]

## Flags

All flags from `release` are supported. See `gitmap release --help`.

## Prerequisites

- The gitmap binary must be located inside or beside its source Git repo
- GitHub CLI (`gh`) recommended for publishing

## Examples

### Example 1: Self-release with bump from another repo

    $ cd ~/projects/other-repo
    $ gitmap rself --bump patch

**Output:**

    → Self-release: switching to /home/user/go/src/gitmap
    v2.45.0 → v2.45.1
    Creating release v2.45.1...
      ✓ Created branch release/v2.45.1
      ✓ Created tag v2.45.1
      ✓ Pushed branch and tag to origin
      ✓ Auto-committed release metadata

      ── Release v2.45.1 complete ──

    ✓ Returned to /home/user/projects/other-repo
    $

### Example 2: Auto-fallback when not in a Git repo

    $ cd /tmp
    $ gitmap release --bump minor

**Output:**

    → Self-release: switching to /home/user/go/src/gitmap
    v2.45.0 → v2.46.0
    Creating release v2.46.0...
      ✓ Created branch release/v2.46.0
      ✓ Created tag v2.46.0
      ✓ Pushed branch and tag to origin
      ✓ Auto-committed release metadata

      ── Release v2.46.0 complete ──

    ✓ Returned to /tmp
    $

## See Also

- [release](release.md) — Full release workflow
- [release-branch](release-branch.md) — Complete release from existing branch

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter release-self
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
