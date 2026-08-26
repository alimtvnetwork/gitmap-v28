# Pattern 02 — Grandfather-Baseline Naming

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

You want to enforce a naming convention (e.g. "all constants must
begin with `Cmd*`, `Msg*`, `Err*`, `Flag*`, or `Default*`") in a
mature codebase that already contains thousands of identifiers
violating it. Renaming everything in one PR is risky and noisy.

The grandfather-baseline pattern lets you **freeze the existing
population** as exempt and enforce the rule **only on identifiers
added after the freeze date**. Over time the codebase converges on
the convention as old names are renamed (or deleted) and the baseline
shrinks.

---

## Algorithm

```
INPUT:  source_dir, baseline_file, allowed_prefix_regex
OUTPUT: exit 0 if no new violations, exit 1 with annotated list

1. extract_identifiers(source_dir) -> sorted unique list (current)
2. read baseline_file -> sorted unique list (baseline)
3. new = current - baseline   # set difference
4. for each name in new:
     if name does not match allowed_prefix_regex:
        find file:line of declaration
        record violation
5. exit 1 if any violations else 0

SUPPORT MODE:
    --regenerate-baseline -> overwrite baseline_file from current
```

The **identifier extractor** must be string-literal aware (see Pattern
03 for details) so that tokens inside multi-line raw strings (SQL
keywords like `WHERE`, `FROM`) are NOT treated as identifiers.

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `bash check-naming.sh` (gate) or `bash check-naming.sh --regenerate-baseline` |
| Exit `0` | No new violations OR regeneration completed |
| Exit `1` | New constants/identifiers that violate the prefix rule |
| Exit `2` | Source dir or baseline file missing |
| stdout | One `::error::` annotation per violation, with file:line and rename example |

---

## Configuration Surface

| Variable | Purpose | Example |
|----------|---------|---------|
| `SOURCE_DIR` | Where identifiers live | `gitmap/constants/`, `src/lib/constants/` |
| `BASELINE_FILE` | Path to grandfathered list | `.github/scripts/constants-baseline.txt` |
| `ALLOWED_PREFIX_REGEX` | Regex matching legal new names | `^(Cmd\|Msg\|Err\|Flag\|Default)` |
| `EXTRACT_FN` | Language-specific identifier extractor | See per-language section |
| `RENAME_EXAMPLES` | Block printed on violation to guide renames | `ScanTimeout → DefaultScanTimeout` |

---

## Baseline File Format

- One identifier per line
- ASCII sorted (`LC_ALL=C sort -u`) — required so `comm -23` works
- No comments, no blank lines, no trailing whitespace
- Committed to the repo and updated only via `--regenerate-baseline`
- Reviewers should treat baseline shrinkage (lines removed) as
  positive signal; growth (lines added) requires justification

---

## Adaptations

### Go (constants)

Extract top-level `const Name = ...` and `const ( Name1 = ... )` block
entries. Use `awk` to track `const (` block state and skip raw-string
regions. See `check-constants-naming.sh` reference.

### Node / TypeScript (exported names)

```bash
EXTRACT_FN() {
  grep -hE '^export +(const|function|class|interface|type|enum) +[A-Z]' \
    "$SOURCE_DIR"/**/*.ts \
    | sed -E 's/^export +(const|function|class|interface|type|enum) +([A-Za-z0-9_]+).*/\2/' \
    | LC_ALL=C sort -u
}
ALLOWED_PREFIX_REGEX='^(use|create|fetch|build|render)'
```

### Python (module-level names)

```bash
EXTRACT_FN() {
  grep -hE '^[A-Z][A-Z0-9_]+ *=' "$SOURCE_DIR"/*.py \
    | sed -E 's/^([A-Z][A-Z0-9_]+).*/\1/' \
    | LC_ALL=C sort -u
}
ALLOWED_PREFIX_REGEX='^(CFG_|MSG_|ERR_|FLAG_|DEFAULT_)'
```

### Rust (`pub const` / `pub static`)

```bash
EXTRACT_FN() {
  grep -hE '^pub (const|static) +[A-Z]' "$SOURCE_DIR"/**/*.rs \
    | sed -E 's/^pub (const|static) +([A-Z][A-Z0-9_]+).*/\2/' \
    | LC_ALL=C sort -u
}
ALLOWED_PREFIX_REGEX='^(CMD_|MSG_|ERR_|FLAG_|DEFAULT_)'
```

---

## Why `awk` State Tracking Matters

The Go reference uses `awk` to maintain three pieces of state per file:

- `in_const` — are we inside a `const ( … )` block?
- `in_rawstr` — are we inside a backtick raw string?
- a one-pass identifier extraction with portability across **gawk and
  mawk** (GitHub Actions runners ship `mawk` as `/usr/bin/awk`)

The mawk constraint forbids `match(string, regex, array)` (a 3-arg
gawk extension). Use only POSIX-portable `match()` + `RSTART` /
`RLENGTH` + `substr()`. A silent parse failure here makes the whole
CI step exit `1` with no `::error::` output.

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| Baseline sorted in non-C locale → `comm -23` produces phantom diffs | Force `LC_ALL=C sort -u` on **both** sides of the diff |
| Awk extracts SQL keywords inside raw-string constants | Track raw-string toggle on every backtick; skip lines while inside |
| Renaming a baselined identifier appears as "new violation" | Document the `--regenerate-baseline` workflow in error output |
| 3-arg `match()` works locally (gawk) but fails on CI (mawk) | Use only 2-arg `match()` + `RSTART` / `RLENGTH` |

---

## Reference Implementation

See `.github/scripts/check-constants-naming.sh` in `gitmap-v6`.

---

## Cross-References

- [00-overview.md](./00-overview.md)
- [01-forbidden-name-guard.md](./01-forbidden-name-guard.md) — Companion for new-code naming
- [03-cross-file-collision-audit.md](./03-cross-file-collision-audit.md) — Detect duplicates across files
- [04-baseline-diff-lint-gate.md](./04-baseline-diff-lint-gate.md) — Same "grandfather + diff" idea applied to lint findings

---

*Grandfather-baseline naming — v1.0.0 — 2026-04-21*