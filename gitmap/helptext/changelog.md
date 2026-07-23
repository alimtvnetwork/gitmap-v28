# gitmap changelog

Display release notes for a specific version or the latest release.

## Alias

cl

## Usage

    gitmap changelog [version] [--open] [--source]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --latest | false | Show latest release notes |
| --open | false | Open changelog in browser |
| --source | false | Show source file path |

## Prerequisites

- Must be inside a Git repository with release metadata

## Rendering

Bullet text is routed through the same pretty markdown renderer used by
`gitmap help`, so formatting stays consistent across the CLI:

- `**bold**` and `` `inline code` `` render with ANSI styling.
- `"double quotes"` highlight in cyan; single quotes / apostrophes are
  left untouched.
- Headers, version, and bullet markers keep their depth-based coloring
  (green at depth 0, cyan at depth 1, dim at depth 2+).

Set `GITMAP_NO_PRETTY=1` or pipe to a non-TTY to get raw markdown. Use
`--pretty` to force ANSI rendering (e.g. `gitmap cl --latest --pretty | less -R`)
or `--no-pretty` to strip every ANSI escape (e.g. for redirects/grep).

## Examples

### Example 1: Show latest changelog

    gitmap changelog --latest

**Output:**

    ═══════════════════════════════════════════
    v2.22.0 — 2025-03-10
    ═══════════════════════════════════════════
    - Add interactive TUI with dashboard view
    - Add zip-group support for release assets
    - Add alias suggest --apply flag
    - Fix watch interval minimum validation
    - Fix cd picker numbering off-by-one
    ═══════════════════════════════════════════
    5 changes | Branch: release/v2.22.0

### Example 2: Show changelog for a specific version

    gitmap cl v2.20.0

**Output:**

    ═══════════════════════════════════════════
    v2.20.0 — 2025-02-28
    ═══════════════════════════════════════════
    - Add clear-release-json command with --dry-run
    - Add release-pending metadata source recovery
    - Fix orphaned metadata detection prompt
    ═══════════════════════════════════════════
    3 changes | Branch: release/v2.20.0

### Example 3: Show changelog with source path

    gitmap changelog v2.20.0 --source

**Output:**

    Source: .gitmap/release/v2.20.0.json
    v2.20.0 — 2025-02-28
    - Add clear-release-json command with --dry-run
    - Add release-pending metadata source recovery

### Example 4: Open changelog in browser

    gitmap changelog --latest --open

**Output:**

    v2.22.0 — 2025-03-10
    ...
    Opening in browser... done

## See Also

- [release](release.md) — Create a release
- [list-versions](list-versions.md) — List release tags
- [list-releases](list-releases.md) — List stored release metadata
- [release-pending](release-pending.md) — Show unreleased commits

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter changelog
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
