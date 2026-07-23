# Multi-URL Clone (spec 104)

> **Status:** planned for v3.38.0
> **Depends on:** existing `gitmap clone` direct-URL support

## What it does

`gitmap clone` accepts **multiple repository URLs in a single invocation**,
either space-separated, comma-separated, or both. Each URL is cloned in
turn (with the existing concurrency model), and on success each repo is
optionally registered with GitHub Desktop.

## Syntax

Both forms are accepted, and they can be mixed freely:

```bash
# Space-separated (Unix-idiomatic, plays well with shell expansion)
gitmap clone https://github.com/a/x.git https://github.com/a/y.git https://github.com/a/z.git

# Comma-separated (single arg, useful when piping or scripting)
gitmap clone https://github.com/a/x.git,https://github.com/a/y.git,https://github.com/a/z.git

# Mixed — each positional arg is split on commas and flattened
gitmap clone url1,url2 url3 url4,url5
```

Parser pseudocode:

```
urls = []
for arg in positional_args:
    for piece in arg.split(","):
        piece = piece.strip()
        if piece:
            urls.append(piece)
```

Whitespace around commas is tolerated. Empty fragments (e.g. trailing comma)
are silently skipped.

## Flags (unchanged from single-URL clone)

| Flag | Effect |
|------|--------|
| `--target-dir <dir>` | Parent folder for all clones (default `.`) |
| `--github-desktop` | Register each successful clone with GitHub Desktop |
| `--workers <n>` | Parallel clone count (default 4) |
| `--stop-on-fail` | Halt the batch on first failure (default: continue) |

## Behaviour

1. Parse all URLs into a flat ordered list, deduplicate (case-insensitive,
   trailing `.git` normalised), preserving first-seen order.
2. Validate each URL syntactically (must look like https://, git@, or
   ssh://). Invalid URLs are reported but the rest still run.
3. Clone in parallel up to `--workers` concurrent jobs.
4. For each successful clone:
   - Insert/upsert into the `Repo` table (existing direct-URL behaviour).
   - If `--github-desktop`, call `desktop.RegisterRepo` immediately.
5. Print a per-repo progress line and a final summary:

   ```
   gitmap clone url1 url2,url3 --github-desktop

     [1/3] x ............ cloned (1.2s)  ✓ registered
     [2/3] y ............ cloned (0.9s)  ✓ registered
     [3/3] z ............ FAILED — auth error

   Clone complete: 2/3 in 2.1s
     Cloned: 2 | Failed: 1 | Registered with Desktop: 2
   ```

## Exit codes

| 0 | All URLs cloned and (if requested) registered |
| 1 | One or more clones failed |
| 2 | Lock contention |
| 3 | All URLs invalid (nothing attempted) |

## Idempotency

- A URL whose target folder already contains a healthy clone is treated as
  a no-op success (matches existing single-URL behaviour).
- `--github-desktop` re-registration of an already-registered repo is a
  no-op at the Desktop layer.

## Why both syntaxes

| Syntax | Best for |
|--------|----------|
| Space-separated | Interactive shell use, tab completion, glob expansion |
| Comma-separated | Single-string CI parameters, JSON arrays joined with commas, copy-paste from issue trackers |
| Mixed | Future-proof — no need to re-quote when refactoring scripts |

## Implementation notes

| File | Responsibility |
|------|----------------|
| `cmd/clone.go` | New `flattenURLArgs([]string) []string` helper at top of `runClone` |
| `cmd/clone.go` | Loop over flattened URLs through existing concurrent worker pool |
| `desktop/desktop.go` | `RegisterRepo(absPath)` — already exists, called per success |
| `constants/constants_clone.go` | `MsgCloneInvalidURLFmt`, `MsgCloneSummaryMultiFmt` |

## See also

- [clone direct-URL feature](mem://features/clone-direct-url) — the single-URL behaviour this extends
- [github-desktop / desktop-sync (spec 11)](/desktop-sync)
