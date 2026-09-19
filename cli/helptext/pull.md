# gitmap pull

Pull a specific tracked repository by slug, group, or all at once with interactive progress bar tracking.

## Alias

p

## Usage

    gitmap pull [<repo-name> | all] [flags]
    gitmap pull-all [flags]
    gitmap pa [flags]
    gitmap git pull [<repo-name>] [flags]
    gitmap git pull-all [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -A, --alias \<name\> | — | Target a repo by its alias |
| --group \<name\> | — | Pull all repos in a group |
| --all | false | Pull all tracked repos |
| --raw | false | Stream raw git output directly instead of using progress bar |
| --verbose | false | Enable verbose logging |
| --parallel \<N\> | 1 | Run up to N pulls concurrently (worker pool) |
| --only-available | false | Skip repos whose latest probe reports no new tag |
| --stop-on-fail | false | Halt the batch after the first failure |

## Key Features

- **Active 80ms Ticker**: Background ticker loop rendering an animated Braille spinner (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`), elapsed time, and live percentages.
- **Single-Repo 4-Step UX**: Explicit milestone tracking: `[Step 1/4] Inspecting`, `[Step 2/4] Fetching remote objects`, `[Step 3/4] Fast-forwarding / Merging`, `[Step 4/4] Complete`.
- **Live Stream Parser**: Automatically passes `--progress` to git pull subprocesses and extracts object events (`Counting`, `Compressing`, `Receiving`, `Resolving deltas`, `Fast-forward`).
- **Concurrent Worker Slots**: Deterministic worker slot mapping (`0..limit-1`) ensuring thread-safe progress updates without parallel terminal line clobbering.
- **Alphabetical Summary Table**: Final batch pull summary tables are consistently sorted alphabetically by repository name.
- **Git Command Routing**: Transparent CLI command passthrough supports `gitmap git pull` and `gitmap git pull-all`.

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Pull a single repo by slug (with animated progress bar)

    gitmap pull my-api

**Output:**

    ⠋ [████████████████████] 100% [Step 4/4] Complete (up-to-date) | ✔ my-api: up-to-date (0.4s)

      REPO      BRANCH   RANGE     CHANGES   STATUS
      -----------------------------------------------
      my-api    main     a1b2c3d   -         ✔ active
      -----------------------------------------------

### Example 2: Pull all tracked repos using the `all` keyword

    gitmap pull all

**Output:**

    ⠋ [████████████░░░░░░░░]  60% (3/5 repos) | [W0: auth-gateway: up-to-date], [W1: payments-api: Receiving (45%)] (1.8s)

      REPO             BRANCH   COMMIT RANGE       CHANGES      PR    STATUS
      -----------------------------------------------------------------------
      auth-gateway     main     9ff44cf            -            -     ✔ active
      billing-svc      main     a1b2c3d..e5f6g7h   +24/-5 (3)   01    ✔ updated
      notification-svc main     77c6edf            -            -     ✔ active
      payments-api     main     03be798..f1d94df   +12/-2 (2)   02    ✔ updated
      user-svc         develop  5e7599b            -            -     ✔ active
      -----------------------------------------------------------------------

### Example 3: Pull all repos in a group

    gitmap p --group backend

### Example 4: Pull in raw streaming mode (unbuffered git output)

    gitmap pull --raw

### Example 5: Parallel pull, only what's actually new

First refresh the probe so `--only-available` has fresh data:

    gitmap probe --all
    gitmap pull all --only-available --parallel 4

### Example 6: Git passthrough command aliases

Execute native GitMap pull workflows directly using standard `git` subcommands:

    gitmap git pull
    gitmap git pull-all

## See Also

- [scan](scan.md) — Scan directories to populate the database
- [clone](clone.md) — Clone repos from output files
- [status](status.md) — Check repo statuses before pulling
- [group](group.md) — Manage groups for targeted pulls
- [alias](alias.md) — Manage repo aliases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pull
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
