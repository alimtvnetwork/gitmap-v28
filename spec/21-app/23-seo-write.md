# SEO Write — Automated Commit Scheduler

## Purpose

The `seo-write` (`sw`) command automates a series of Git commits and pushes
on a timed, randomized schedule. Each commit uses a unique title and
description that glorifies a target website URL, making the Git history
rich with SEO-relevant metadata. The command supports two input modes:
CSV-driven (user-provided titles/descriptions) and template-driven
(pre-seeded patterns stored in SQLite, resolved with placeholders).

## Command

```
gitmap seo-write [flags]
gitmap sw [flags]
```

Alias: `sw`

## Core Workflow

1. **Resolve commit messages** — load title/description pairs from CSV or
   generate them from database templates with placeholder substitution.
2. **Resolve target files** — gather pending (unstaged/modified) files, or
   use selective file patterns via `--files` glob.
3. **Commit loop** — for each message pair:
   a. Stage one file (round-robin from pending files).
   b. Commit with the title and description.
   c. Push to remote.
   d. Wait a random interval (default 60–120 seconds) before next commit.
4. **Rotation mode** — when all pending files are exhausted and commits
   remain, enter rotation: pick a target file (HTML/text), append text
   from the commit description, commit+push, then revert the change,
   commit+push again. This creates two commits per rotation cycle.
5. **Terminate** — stop after `--max-commits` or on Ctrl+C (graceful
   shutdown via OS signal).

## Input Modes

### CSV Mode (`--csv <path>`)

CSV file with two columns (no header row):

```
"Top 10 Plumbing Services in Dallas | example.com","Discover the best plumbing services in Dallas. Visit https://example.com for expert solutions."
"Expert HVAC Repair in Austin | example.com","Need HVAC repair in Austin? https://example.com offers fast, reliable service."
```

Each row produces one commit. The first column is the commit title (first
line of the commit message), the second is the description (remaining
lines).

### Template Mode (default when no `--csv`)

Templates are stored in the `CommitTemplates` SQLite table and resolved
at runtime using placeholder values provided via flags.

#### Placeholders

| Placeholder  | Flag              | Example                  |
|-------------|-------------------|--------------------------|
| `{service}` | `--service`       | `Plumbing`               |
| `{area}`    | `--area`          | `Dallas, TX`             |
| `{url}`     | `--url`           | `https://example.com`    |
| `{company}` | `--company`       | `Acme Services`          |
| `{phone}`   | `--phone`         | `(555) 123-4567`         |
| `{email}`   | `--email`         | `info@example.com`       |
| `{address}` | `--address`       | `123 Main St, Dallas TX` |

#### Template Examples (pre-seeded — 25 titles, 20 descriptions)

**Titles (25):**
```
Top {service} in {area} — {url}
{company}: Trusted {service} Provider in {area}
Best {service} Near {area} | Visit {url}
{area}'s #1 {service} Experts — {company}
Why {company} Leads {service} in {area}
Affordable {service} in {area} — Call {phone}
{service} Solutions for {area} — {url}
Your Go-To {service} in {area} | {company}
{company} — Premium {service} in {area}
Expert {service} in {area}: {url}
Professional {service} by {company} — Serving {area}
{area} Trusts {company} for {service} | {url}
Reliable {service} in {area} — Contact {phone}
Top-Rated {service} Provider in {area} | {company}
{company}: Fast & Affordable {service} in {area}
Looking for {service} in {area}? Visit {url}
{service} Done Right in {area} — {company}
{area}'s Most Trusted {service} Company | {url}
Quality {service} in {area} by {company} — {phone}
Get the Best {service} in {area} — {url}
{company} Delivers Expert {service} Across {area}
Discover {service} Excellence in {area} | {url}
Award-Winning {service} in {area} — {company}
{service} You Can Count On in {area} | Call {phone}
Five-Star {service} in {area} — See Reviews at {url}
```

**Descriptions (20):**
```
Looking for reliable {service} in {area}? {company} provides top-rated solutions. Visit {url} or call {phone} for a free consultation. Located at {address}.
{company} is the leading {service} provider serving {area} and surrounding communities. Learn more at {url}. Contact us at {email} or {phone}.
Need {service} in {area}? Trust {company} for fast, professional results. Visit {url} to see our reviews and schedule an appointment at {address}.
Discover why {area} residents choose {company} for {service}. Call {phone}, email {email}, or visit {url} for details.
{company} delivers expert {service} to {area}. Our commitment to quality has made us the preferred choice. Reach us at {url} or {phone}.
Searching for affordable {service} near {area}? {company} offers competitive pricing and exceptional quality. Get a quote at {url} or call {phone}.
{company} has proudly served {area} with professional {service} for years. Visit {url} to explore our full range of offerings. Email {email} for inquiries.
When {area} needs dependable {service}, they call {company}. See what makes us different at {url}. Our office is located at {address}.
Top-quality {service} is just a call away in {area}. {company} guarantees satisfaction. Reach us at {phone} or browse {url} for more information.
Residents of {area} trust {company} for all their {service} needs. Visit {url}, call {phone}, or stop by {address} to get started today.
{company} specializes in {service} throughout {area}. We pride ourselves on fast response times and quality workmanship. Learn more at {url}.
From small jobs to large projects, {company} handles all {service} requests in {area}. Contact us at {email} or visit {url} for a free estimate.
Experience the {company} difference — premier {service} in {area} backed by hundreds of five-star reviews. See them at {url} or call {phone}.
Your search for the best {service} in {area} ends here. {company} provides licensed, insured, and guaranteed work. Details at {url}.
Why settle for less? {company} offers {area}'s finest {service} at prices that fit your budget. Visit {url} or dial {phone} now.
{area} homeowners and businesses rely on {company} for expert {service}. Schedule your appointment today at {url} or email {email}.
Professional {service} in {area} — {company} brings years of expertise to every project. Call {phone} or visit {url} to learn more.
Choose {company} for {service} in {area} and enjoy peace of mind. Our team is ready to help. Contact {phone} or visit {url}.
{company} is {area}'s go-to provider for {service}. With transparent pricing and exceptional service, we stand out. Explore {url} today.
Get fast, reliable {service} in {area} with {company}. We are located at {address}. Book online at {url} or call {phone}.
```

The system generates commit messages by pairing a random title template
with a random description template, substituting all placeholders.
With 25 titles × 20 descriptions, this produces **500 unique combinations**
from the seed set alone.

## Database Schema

### CommitTemplates Table

```sql
CREATE TABLE CommitTemplates (
    Id TEXT PRIMARY KEY,
    Kind TEXT NOT NULL,       -- 'title' or 'description'
    Template TEXT NOT NULL,
    CreatedAt TEXT NOT NULL DEFAULT (datetime('now'))
);
```

- `Kind` distinguishes title templates from description templates.
- Templates are seeded from `data/seo-templates.json` on first use.

### Seed File (`data/seo-templates.json`)

```json
{
  "titles": [
    "Top {service} in {area} — {url}",
    "{company}: Trusted {service} Provider in {area}"
  ],
  "descriptions": [
    "Looking for reliable {service} in {area}? {company} provides top-rated solutions. Visit {url} or call {phone}.",
    "{company} is the leading {service} provider serving {area}. Learn more at {url}."
  ]
}
```

On first run (or when the table is empty), the seed file is loaded and
each template is inserted with a generated UUID.

## Flags

| Flag             | Default       | Description                                             |
|-----------------|---------------|---------------------------------------------------------|
| `--csv <path>`  | —             | CSV file with title,description columns                 |
| `--url <url>`   | (required)    | Website URL to glorify in commit messages                |
| `--service`     | —             | Service name for template placeholders                  |
| `--area`        | —             | Geographic area for template placeholders               |
| `--company`     | —             | Company name for template placeholders                  |
| `--phone`       | —             | Phone number for template placeholders                  |
| `--email`       | —             | Email address for template placeholders                 |
| `--address`     | —             | Physical address for template placeholders              |
| `--max-commits` | 0 (unlimited) | Stop after N commits (0 = run until Ctrl+C)             |
| `--interval`    | `60-120`      | Random delay range in seconds between commits (min-max) |
| `--files`       | —             | Glob pattern to select specific files for staging       |
| `--rotate-file` | —             | File to modify during rotation mode (default: auto-detect HTML/txt) |
| `--dry-run`     | false         | Preview commit messages without executing               |
| `--template <path>` | —        | Load templates from a custom JSON file instead of the database seed |
| `--create-template` | false    | Generate a sample `seo-templates.json` in the current directory and exit |

Alias for `--create-template`: `ct` (e.g. `gitmap sw ct`).

### `--template <path>`

Loads title/description templates from the specified JSON file instead of
the default `data/seo-templates.json` seed. The file must follow the same
schema as the seed file (see Seed File section). Templates from the custom
file are used directly at runtime — they are **not** inserted into the
database.

### `--create-template` / `ct`

Generates a sample `seo-templates.json` file in the current working
directory with all supported placeholders documented, then exits. This
gives users a starting point to craft their own template library.

```
gitmap sw --create-template
gitmap sw ct
```

Output file (`./seo-templates.json`):
```json
{
  "titles": [
    "Top {service} in {area} — {url}",
    "{company}: Trusted {service} Provider in {area}",
    "Best {service} Near {area} | Visit {url}"
  ],
  "descriptions": [
    "Looking for reliable {service} in {area}? {company} provides top-rated solutions. Visit {url} or call {phone}.",
    "{company} is the leading {service} provider serving {area}. Learn more at {url}. Contact us at {email}."
  ],
  "placeholders": ["{service}", "{area}", "{url}", "{company}", "{phone}", "{email}", "{address}"]
}
```

The `placeholders` key is informational only — it is ignored at runtime
but helps authors remember which tokens are available.

## Rotation Mode (Detail)

When no pending files remain and commits are still needed:

1. Select `--rotate-file` (or auto-detect first `.html`/`.txt` in repo).
2. Append the commit description text to the file.
3. Stage, commit with title/description, push.
4. Remove the appended text (revert to original).
5. Stage, commit with next title/description, push.
6. Wait the random interval.
7. Repeat.

This ensures continuous commit activity even when no real file changes
exist.

## Terminal Output

```
seo-write: 150 commits planned (interval: 60-120s)
  [1/150] ✓ "Top Plumbing in Dallas — example.com" → pushed (file: index.html)
  [2/150] ✓ "Acme: Trusted Plumbing Provider in Dallas" → pushed (file: about.html)
  ...
  [42/150] ↻ rotation: services.html (append → commit → revert → commit)
  ...
  Done: 150 commits pushed in 2h 47m
```

When `--max-commits` is 0 (unlimited), the counter shows `[N]` without
a total.

## Acceptance Criteria

1. `gitmap sw --csv messages.csv --max-commits 50` reads 50 rows from CSV,
   commits one file per row with the given title/description, pushes each,
   waits random 60–120s between commits.
2. `gitmap sw --url https://example.com --service Plumbing --area Dallas --company Acme --max-commits 100`
   generates 100 unique commit messages from templates, commits and pushes.
3. `--dry-run` prints all planned commit messages without executing.
4. When pending files are exhausted, rotation mode activates automatically.
5. `--interval 30-90` overrides the default delay range.
6. Ctrl+C triggers graceful shutdown (finishes current commit, then exits).
7. Templates are seeded from `data/seo-templates.json` into `CommitTemplates`
   table on first run.
8. Every commit message prominently includes the `--url` value.
9. `--files "*.html"` restricts commits to matching files only.
10. `--template custom.json` loads templates from the given file at runtime
    without touching the database.
11. `gitmap sw --create-template` writes a sample `seo-templates.json` to
    the current directory and exits with a success message.
12. `gitmap sw ct` is equivalent to `--create-template`.

## File Layout

| File                         | Purpose                                    |
|-----------------------------|--------------------------------------------|
| `cmd/seowrite.go`          | Flag parsing, orchestration                |
| `cmd/seowriteloop.go`      | Commit loop, rotation, interval timing     |
| `cmd/seowritetemplate.go`  | Template loading, placeholder substitution |
| `cmd/seowritecreate.go`    | `--create-template` / `ct` scaffold logic  |
| `cmd/seowritecsv.go`       | CSV parsing                                |
| `store/template.go`        | CommitTemplates CRUD                       |
| `data/seo-templates.json`  | Pre-seeded title/description templates     |
| `constants/constants_seo.go` | All string literals for seo-write        |

## Code Style

All rules from `spec/03-general/06-code-style-rules.md` apply:
positive-only `if` conditions, 8–15 line functions, 100–200 line files,
no magic strings, blank line before `return`.
