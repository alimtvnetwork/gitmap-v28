# gitmap seo-write

Auto-generate and commit SEO-optimized messages to a repository on a schedule.

## Alias

sw

## Usage

    gitmap seo-write [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --csv \<path\> | — | CSV file with SEO data |
| --url \<url\> | — | Target website URL |
| --service \<name\> | — | Service name |
| --area \<name\> | — | Geographic area |
| --company \<name\> | — | Company name |
| --phone \<number\> | — | Phone number |
| --email \<addr\> | — | Contact email |
| --address \<addr\> | — | Physical address |
| --max-commits \<N\> | 10 | Maximum commits per run |
| --interval \<secs\> | 60 | Seconds between commits |
| --files \<list\> | — | Files to modify |
| --rotate | false | Rotate through templates |
| --dry-run | false | Preview without committing |
| --template \<path\> | — | Custom template file |
| --create-template | false | Generate a starter template |
| --author-name \<n\> | — | Commit author name |
| --author-email \<e\> | — | Commit author email |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Run SEO writes from CSV

    gitmap seo-write --csv data.csv --max-commits 5 --interval 30

**Output:**

    Loading SEO data from data.csv...
    Found 12 messages (using first 5)
    [1/5] "Best plumber in Seattle — 24/7 emergency"... committed
          Waiting 30s...
    [2/5] "Licensed plumbing contractor in Seattle WA"... committed
          Waiting 30s...
    [3/5] "Affordable drain cleaning services Seattle"... committed
          Waiting 30s...
    [4/5] "Emergency pipe repair — call now"... committed
          Waiting 30s...
    [5/5] "Seattle's top-rated plumbing company"... committed
    ✓ 5 commits created (30s intervals)

### Example 2: Dry-run preview

    gitmap sw --csv data.csv --dry-run --max-commits 3

**Output:**

    [DRY RUN] Loading SEO data from data.csv...
    [DRY RUN] Found 12 messages (would use first 3)
    [DRY RUN] 1. "Best plumber in Seattle — 24/7 emergency"
    [DRY RUN] 2. "Licensed plumbing contractor in Seattle WA"
    [DRY RUN] 3. "Affordable drain cleaning services Seattle"
    No changes made.

### Example 3: Create a starter template

    gitmap seo-write --create-template

**Output:**

    ✓ Template created at ./data/seo-templates.json
    → Edit the template, then run:
      gitmap seo-write --url https://example.com --service "Plumbing"

### Example 4: SEO writes from template with custom author

    gitmap sw --url https://plumber.com --service "Plumbing" --area "Seattle" \
              --company "AcePlumb" --author-name "SEO Bot" --author-email "seo@aceplumb.com"

**Output:**

    Loading template for https://plumber.com...
    Generated 10 messages for "Plumbing" in "Seattle"
    Author: SEO Bot <seo@aceplumb.com>
    [1/10] "AcePlumb — Best Plumbing in Seattle"... committed
    ...
    ✓ 10 commits created

## See Also

- [scan](scan.md) — Scan directories to find target repos
- [history](history.md) — View command execution history

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter seo-write
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
