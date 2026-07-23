# gitmap diff-profiles

Compare repositories across two profiles to find differences.

## Alias

dp

## Usage

    gitmap diff-profiles <profileA> <profileB> [--all] [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --all | false | Show all repos, not just differences |
| --json | false | Output as structured JSON |

## Prerequisites

- At least two profiles must exist (see profile.md)

## Examples

### Example 1: Compare two profiles

    gitmap diff-profiles home work

**Output:**

    Comparing profiles: home vs work
    ═══════════════════════════════════════════
    Only in 'home' (3):
      personal-blog
      side-project
      dotfiles
    Only in 'work' (2):
      billing-svc
      internal-tools
    Common repos: 12
    ═══════════════════════════════════════════
    Summary: 3 unique to home, 2 unique to work, 12 shared

### Example 2: Full comparison showing all repos

    gitmap dp home work --all

**Output:**

    REPO              HOME    WORK
    my-api            ✓       ✓
    web-app           ✓       ✓
    billing-svc       —       ✓
    personal-blog     ✓       —
    shared-lib        ✓       ✓
    ...
    17 repos total (12 common, 3 home-only, 2 work-only)

### Example 3: JSON output for scripting

    gitmap dp home work --json

**Output:**

    {
      "only_a": ["personal-blog", "side-project", "dotfiles"],
      "only_b": ["billing-svc", "internal-tools"],
      "common": ["my-api", "web-app", "shared-lib", ...],
      "summary": {"only_a": 3, "only_b": 2, "common": 12}
    }

## See Also

- [profile](profile.md) — Create and manage profiles
- [list](list.md) — View repos in the current profile
- [group](group.md) — Organize repos within a profile

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter diff-profiles
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
