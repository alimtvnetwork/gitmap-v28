# `gitmap-v27 replace` — Repo-Wide Text & Version Replace

> **Related specs:**
> - [03-subcommand-architecture.md](03-subcommand-architecture.md) — dispatch pattern
> - [04-flag-parsing.md](04-flag-parsing.md) — per-command FlagSet conventions
> - [../05-coding-guidelines](../05-coding-guidelines) — file/func size, error policy

This document is the **single source of truth** for the `gitmap-v27 replace`
subcommand. It is written so a coding AI can implement or modify the
feature without inferring intent from other files. Every behavior,
edge case, and exit code is fixed here.

---

## 1. Purpose

`gitmap-v27 replace` performs a deterministic, repo-wide find-and-replace
across every text file in the current repository, with two operating
modes:

1. **Literal mode** — `gitmap-v27 replace "<old>" "<new>"`
   Replaces every occurrence of `<old>` with `<new>` in every text file.
2. **Version mode** — `gitmap-v27 replace -N` / `--audit` / `all`
   Bumps occurrences of `<base>-vK` (the project name + version suffix
   parsed from the git remote) to the current version.

Both modes share the same file walker, exclusion rules, and confirmation
flow. The default is **interactive confirm before write** — nothing is
modified until the user types `y` at the diff preview.

---

## 2. Invocation Surface

```text
gitmap-v27 replace "<old>" "<new>"        # literal mode (interactive confirm)
gitmap-v27 replace -N                      # version mode: replace v(current-N)..v(current-1) → vCurrent
gitmap-v27 replace --audit                 # version mode: report only (no writes ever)
gitmap-v27 replace all                     # version mode: replace v1..v(current-1) → vCurrent
gitmap-v27 replace --help | -h             # show help
```

Aliases: `rpl` (short form of `replace`).

Flags accepted in both modes:

| Flag | Type | Default | Meaning |
|------|------|---------|---------|
| `--yes` / `-y` | bool | false | Skip the interactive `y/N` prompt |
| `--dry-run` | bool | false | Force report-only behavior (overrides `--yes`) |
| `--quiet` / `-q` | bool | false | Suppress per-file diff lines; print summary only |
| `--ext` | string | "" | Comma-separated extension allow-list (e.g. `.go,.md`). Leading dot optional, case-insensitive, deduplicated. Empty = all text files. |
| `--ext-case` | string | `insensitive` | Match casing for `--ext`: `sensitive` (byte-exact, preserves user input case) or `insensitive` (lowercased on both sides). Unknown values exit 1. |

Exit codes:

| Code | Condition |
|------|-----------|
| `0`  | Success — replacements applied OR audit completed cleanly |
| `1`  | Bad arguments, no git remote, version unparseable, or user aborted |
| `2`  | I/O failure during walk or write (filesystem, permissions) |

---

## 3. Repository Identity Detection

Version mode parses the **git remote URL** of the current working
directory. No fallback to `go.mod`, no fallback to flag — single source
of truth.

Algorithm (must be implemented exactly):

1. Run `git remote get-url origin` from the current working directory.
2. Trim whitespace and trailing `.git`.
3. Take the last path segment after the final `/`. Call it `slug`.
4. Match `slug` against the regex `^(?P<base>.+)-v(?P<num>\d+)$`.
5. If no match → exit `1` with message
   `replace: cannot detect version from remote: %q (expected suffix -vN)`.
6. `base` is the project base name. `num` (parsed as int) is the
   current version `K`. Both are required for the search pattern.

**Example:** `git@github.com:alimtvnetwork/gitmap-v27.git`
→ `slug = "gitmap-v27"` → `base = "gitmap-v27"`, `K = 7`.

---

## 4. Replace Pattern (Version Mode)

For a given target version `T` (an int < `K`), the search pattern is the
**literal string** `<base>-v<T>`. Replacement is the literal string
`<base>-v<K>`.

By the user's explicit answer ("Any occurrence of the base name +
version"), we also replace `<base>/v<T>` → `<base>/v<K>` in the same
pass (covers Go module import paths like `github.com/x/gitmap-v27` and
`github.com/x/gitmap-v27/v4`).

We do **not** touch bare `vN` tokens that aren't adjacent to `<base>`.
That avoids destroying CSS values, semver references, etc.

---

## 5. File Selection

The walker starts at the repo root (the directory returned by
`git rev-parse --show-toplevel`).

### 5.1 Excluded directories (never descended)

```
.git
.gitmap
.release
node_modules
vendor
```

### 5.2 Excluded path prefixes (relative to repo root)

```
.gitmap/release
.gitmap/release-assets
```

### 5.3 Per-file exclusion: binary detection

Open each candidate file, read up to the first 8192 bytes. If any byte
is `0x00`, treat as binary and skip. This is the same heuristic git
uses and matches the user's "any text file" intent.

### 5.4 Optional extension allow-list (`--ext`)

By default, every file that survives §5.1–5.3 is eligible — there is no
extension whitelist. When the user passes `--ext`, the walker
additionally requires `filepath.Ext(path)` (lowercased) to appear in
the supplied list before the file is considered.

Normalization rules for `--ext` values:

- Split on `,`.
- Trim ASCII whitespace from each piece.
- Lowercase each piece.
- Prepend `.` when the first byte is not already `.`.
- Drop empty pieces and the lone string `"."`.
- Deduplicate while preserving first-seen order.

Examples:

| Input             | Normalized       |
|-------------------|------------------|
| `.go,.md`         | `[.go .md]`      |
| `go,MD`           | `[.go .md]`      |
| `  .go , md ,go ` | `[.go .md]`      |
| `""` (omitted)    | `nil` (no filter) |

Files with no extension never match a non-empty `--ext`. Binary
detection (§5.3) still runs on filtered-in files.

---

## 6. Operating Modes (Detail)

### 6.1 Literal mode

Triggered when **two non-flag positional args** are present.

```
gitmap-v27 replace "old text" "new text"
```

- `old` MUST be non-empty; `new` MAY be empty (deletion).
- Walker scans every eligible file, counts occurrences of `old`.
- Prints a summary table: `<file>: <count>` for each affected file.
- If `--dry-run` → exit 0.
- Otherwise prompt `Apply N replacements across M files? [y/N]:`.
  - `--yes` / `-y` skips the prompt.
- On confirm, rewrite each affected file atomically (write temp → rename).

### 6.2 Version mode `-N`

Triggered when the only positional/flag arg matches `^-(\d+)$` and
`N >= 1`.

- Compute `targets = [K-N, K-N+1, ..., K-1]`. If `K-N < 1`, clamp to 1.
- For each `T` in `targets` (ascending), perform the same scan + replace
  pass as literal mode using the patterns from §4.
- A single confirm prompt covers all passes:
  `Apply replacements for versions v%d..v%d → v%d? [y/N]:`.

### 6.3 Version mode `--audit`

- Same scan as `all` would do.
- Prints, per detected older-version reference:
  `<file>:<line>: <matched-text>`
- **Never writes**, never prompts. `--yes` is ignored. Exit 0 even when
  matches are found.

### 6.4 Version mode `all`

- Equivalent to `-N` with `N = K - 1` (i.e. `targets = [1..K-1]`).
- If `K == 1` → print `replace: already at v1, nothing to upgrade` and
  exit 0.

---

## 7. Atomic Write Contract

Every modified file MUST be written via:

1. `os.CreateTemp(dir, base+".gitmap-replace-*")` in the same directory.
2. Write new contents, `Sync()`, `Close()`.
3. `os.Rename(tmp, original)`.

On any error, remove the temp file and surface the error per the Code
Red zero-swallow policy: `fmt.Fprintf(os.Stderr, "replace: %s: %v\n",
path, err)` and exit `2` after attempting all remaining files.

---

## 8. Output Format

```
replace: scanning 4123 files in /repo/root
replace: src/foo.go: 3 matches (gitmap-v27 → gitmap-v27)
replace: docs/setup.md: 1 match (gitmap-v27 → gitmap-v27)
...
replace: 12 files, 47 replacements
Apply replacements for versions v4..v6 → v7? [y/N]: y
replace: applied 47 replacements across 12 files
```

`--quiet` suppresses the per-file lines but keeps the summary and prompt.

---

## 9. Test Matrix (must pass)

| Case | Expectation |
|------|-------------|
| `replace "a" "b"` in repo with binary file | Binary skipped, text rewritten |
| `replace -3` on v7 repo | Replaces v4, v5, v6 → v7 |
| `replace -10` on v3 repo | Clamps; replaces v1, v2 → v3 |
| `replace all` on v1 repo | No-op message, exit 0 |
| `replace --audit` finds matches | Reports, exits 0, file unchanged |
| Repo with no remote | Exit 1 with clear error |
| Remote without `-vN` suffix | Exit 1 with clear error |
| User answers `n` to prompt | No writes, exit 1 |
| `--yes` flag | No prompt, writes immediately |
| File inside `.gitmap/release-assets/` contains match | Skipped |
| `replace -1 --ext .go,.md` on mixed repo | Only `.go` and `.md` files scanned/written |
| `--ext go` (no leading dot) | Normalized to `.go`, otherwise identical to dotted form |
| `--ext .GO` on `app.go` | Case-insensitive match, file included |

---

## 10. Non-Goals

- No regex mode in v1. Patterns are always literal.
- No git commit/push step. Replace only mutates the working tree.
- No language-aware AST rewriting. This is text replace.
- No multi-repo execution. `replace` runs against the cwd's repo only.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC](https://riseup-asia.com) (2026)
