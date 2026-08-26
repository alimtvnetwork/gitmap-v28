# Pattern 03 — Cross-File Collision Audit

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

The compiler / type checker catches **exact** duplicate declarations
(`redeclared in this block`). It does NOT catch:

1. **Case-insensitive collisions** across files — `HelpFoo` in `a.go`
   and `helpFoo` in `b.go` compile fine but are confusing and a sign
   of inconsistent naming.
2. **Intra-file duplicates** sometimes exist in legacy files where the
   compile error is only triggered when the file is actually built;
   surfacing them in CI gives reviewers context without waiting for
   the build.
3. **Cross-file exact collisions** in pre-build static-analysis where
   the goal is to fail fast (within seconds) instead of waiting for
   the full compile.

---

## Algorithm

```
INPUT:  file_glob (e.g. constants/*.{go,ts,py,rs})
OUTPUT: exit 0 if no collisions, exit 1 with categorized report

1. for each file in glob:
     parse top-level identifier declarations -> (kind, name, lineno)
     skip declarations inside string literals / comments

2. build three index maps:
     exact[name]      -> [(file, line, kind), ...]
     case_insens[name.lower()] -> [(name, file, line, kind), ...]
     intra_file[name][file]    -> [(line, kind), ...]

3. report categories:
     [1] cross_file_exact   = exact[name] spans > 1 file
     [2] case_insensitive   = case_insens[lower] has > 1 distinct name AND spans > 1 file
     [3] intra_file_dupes   = intra_file[name][file] has > 1 entry

4. exit 1 if any category non-empty
```

The parser MUST track string-literal state. The reference Python
implementation tracks **raw-string** (backtick in Go), **regular
string** (double-quote), and `const ( … )` / `var ( … )` block state.

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `python3 check-collisions.py` |
| Exit `0` | No collisions in any of the three categories |
| Exit `1` | One or more categories non-empty (full report on stdout) |
| Exit `2` | No source files matched the glob |
| stdout | Per-category report with `[kind] file:line` for every site |

---

## Configuration Surface

| Variable | Purpose | Example |
|----------|---------|---------|
| `FILE_GLOB` | Files to scan | `gitmap/constants/constants*.go`, `src/**/*.ts` |
| `IDENT_RE` | Regex for top-level identifier names | `^[A-Z][A-Za-z0-9_]*` (Go/TS) or `^[A-Z][A-Z0-9_]*` (Python/Rust SCREAMING_SNAKE) |
| `BLOCK_OPEN_RE` | Optional: regex matching multi-decl block opener | `^\s*const\s*\(` (Go) |
| `BLOCK_CLOSE_RE` | Optional: regex matching block closer | `^\s*\)` |
| `RAWSTRING_DELIM` | Char that toggles raw-string state | `` ` `` (Go), `"""` (Python triple-quote) |

---

## String-Literal Awareness

This is the single most important property. Without it, a constant like:

```go
const SQLInsertUser = `
    INSERT INTO User (Id, Name)
    VALUES (?, ?)
`
```

causes the parser to emit `INSERT`, `INTO`, `User`, `VALUES` as
top-level identifiers — every one a false positive collision with the
rest of the codebase.

The reference parser handles three layers in this order:

1. If currently inside a raw string, scan for the closing delimiter
   on the line; if not found, skip the line.
2. Truncate the line at the first `` ` `` or `"` (whichever comes
   first) so subsequent identifier matching only sees code outside
   strings.
3. Toggle raw-string state if the line contains an odd number of
   backticks.

---

## Adaptations

### Go

- `FILE_GLOB`: `package_dir/*.go`
- `IDENT_RE`: `^[A-Z][A-Za-z0-9_]*`
- Block delimiters: `const (` / `var (` / `)`
- Raw string: backtick

### Node / TypeScript

- `FILE_GLOB`: `src/**/*.ts`
- `IDENT_RE`: `^export +(const|function|class|interface|type|enum) +([A-Za-z_][A-Za-z0-9_]*)`
- Block delimiters: not applicable (one declaration per `export` line)
- Template-literal: backtick — toggle the same way as Go raw strings

### Python

- `FILE_GLOB`: `pkg/*.py`
- `IDENT_RE`: `^([A-Z][A-Z0-9_]*) *=`
- Block delimiters: not applicable
- Triple-quoted strings (`"""` / `'''`): toggle on triple-delimiter
  occurrence; treat the parser state machine the same as Go raw strings

### Rust

- `FILE_GLOB`: `src/**/*.rs`
- `IDENT_RE`: `^pub (const|static) +([A-Z][A-Z0-9_]*)`
- Block delimiters: not applicable
- Raw string: `r#"..."#` — track the `#` count on the opening to know
  the closing delimiter

---

## Output Format

```
Scanned 12 constants_*.go files
Total unique top-level identifiers: 2734

::error::[1] 2 cross-file exact-name collision(s):
  CmdScan
      [const] constants_cmd.go:42
      [const] constants_legacy.go:118

::error::[2] 1 cross-file case-insensitive collision(s):
  'helpfoo':
      [const] HelpFoo @ constants_msg.go:201
      [const] helpFoo @ constants_internal.go:88

::error::[3] 0 intra-file duplicate declaration(s):
```

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| Triple-quoted Python string treated as 3 separate single-quotes | Detect `"""` / `'''` as a single token before the per-quote scan |
| Backtick inside double-quoted string toggles raw-string state | Process strings outside-in: truncate at first quote of either kind, then count backticks only on the truncated portion |
| Identifier extracted from a comment line | Strip `//` and `#` comments before identifier matching |
| Glob with no matches silently passes | Treat empty file list as exit `2` (tool error), not exit `0` |

---

## Reference Implementation

See `.github/scripts/check-constants-collisions.py` in `gitmap-v6`.

---

## Cross-References

- [00-overview.md](./00-overview.md)
- [02-grandfather-baseline-naming.md](./02-grandfather-baseline-naming.md) — Shares the string-literal-aware extractor

---

*Cross-file collision audit — v1.0.0 — 2026-04-21*