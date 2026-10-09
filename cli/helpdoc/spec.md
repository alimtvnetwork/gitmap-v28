# gitmap spec

Issue concurrency-safe spec numbers for the current repository (spec 252).

## Usage

    gitmap spec next [--json]
    gitmap spec --help

## Subcommands

- `next` — Issue the next free spec number and print it on stdout.
- `--help` — Print this help.

## Flags (with next)

- `--json` — Print the issuance record as JSON instead of a bare number.
- `-h`, `--help` — Print usage, exit 0.

## Behavior

- The candidate is `max(highest ^(\d+)- directory prefix under
  02-spec/21-app/, highest number recorded in the Tier-1 agent
  database) + 1`.
- The claim is atomic (`INSERT OR IGNORE` over a `UNIQUE(number)`
  column): two parallel callers each get a distinct number; the loser
  retries to the next free one.
- Plain mode prints ONLY the number on stdout — script-friendly:

      N=$(gitmap spec next) && mkdir "02-spec/21-app/${N}-my-spec"

## JSON mode

`gitmap spec next --json` prints exactly one pretty-printed object:

    {
      "repo": "<absolute repo path>",
      "number": 252,
      "issued_at": "<RFC3339>"
    }

- `repo` is the absolute path of the repository the number was issued for.
- `issued_at` is the issuance timestamp in RFC3339 (UTC).

## Exit codes

- `0` — number issued (plain or `--json`), or usage printed.
- `2` — usage error (unknown flag).
- `1` — issuance failure (DB open, scan, or claim failed). The error
  goes to stderr; stdout contains no number.

## Examples

    gitmap spec next
    gitmap spec next --json
    gitmap spec --help
