# Pattern 01 — Forbidden-Name Guard

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

Some packages, modules, or namespaces are **flat** — every identifier
declared in any file shares one symbol table. Examples: a Go `package
cmd` directory, a Python module's `__init__.py` re-exports, a single
Rust `mod foo` file split across `mod.rs` includes, a Node ESM
barrel file.

In flat-namespace packages, two contributors can independently invent
the same generic helper name (`invokeRelease`, `runOne`, `persistAll`,
`processItem`) in different files. The result is a hard build break
(`redeclared in this block` / `duplicate identifier`) that surfaces
only at compile time and blocks every other PR until resolved.

The fix is **prevention via naming policy**: bare or over-generic
helper names are forbidden; they MUST be qualified with a domain
word that narrows the scope.

---

## Algorithm

```
INPUT:  source_dir, forbidden_patterns[], allowed_router_patterns[]
OUTPUT: exit 0 if clean, exit 1 with annotated violations

1. validate source_dir exists
2. for each forbidden_pattern:
     scan source_dir for matching declarations
     if matches found:
        emit ::error annotation per match (file:line)
        increment violation_count
3. exit 1 if violation_count > 0 else 0
```

The patterns are written as **anchored regex** against function/method
declaration syntax of the target language. Two pattern families are
used:

- **Bare verb**: `^func +(invoke|persist) *\(` → `invoke(...)` with no
  domain qualifier.
- **Verb + over-generic noun**: `^func +(invoke|persist|runOne)(Release|Task|Job|Item|All|One|Cmd) *\(`
  → `invokeRelease`, `persistTask` — the noun must be project-specific
  (e.g. `AliasRelease`, `ScanRelease`), not a generic placeholder.

Canonical router patterns (e.g. `executeXxx`, `handleXxx` in a CLI
command dispatcher — exactly one per top-level command) are
**explicitly allowed** and not flagged.

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `bash check-forbidden-names.sh [source_dir]` |
| Exit `0` | No violations |
| Exit `1` | One or more violations (full list printed) |
| Exit `2` | Tool error (source dir missing) |
| stdout | One `::error::` annotation per violation + remediation tip |

---

## Configuration Surface

| Variable | Purpose | Example |
|----------|---------|---------|
| `SOURCE_DIR` | Where to scan | `cmd/`, `src/handlers/`, `lib/` |
| `FILE_GLOB` | Files to inspect | `*.go`, `*.ts`, `*.py`, `*.rs` |
| `FORBIDDEN_PATTERNS` | Anchored regex array | See per-language section |
| `ALLOWED_ROUTER_PATTERNS` | Documentation only — names matching these are not in `FORBIDDEN_PATTERNS` | `executeXxx`, `handleXxx` |
| `REMEDIATION_EXAMPLES` | Lines printed to help the contributor rename | `invokeRelease → invokeAliasRelease` |

---

## Adaptations

### Go

```bash
FILE_GLOB="*.go"
FORBIDDEN_PATTERNS=(
  '^func +(invoke|persist) *\('
  '^func +runOne *\('
  '^func +(invoke|persist|runOne)(Release|Task|Job|Item|All|One|Cmd) *\('
)
```

Run with `grep -HnE "$pat" "$SOURCE_DIR"/*.go`.

### Node / TypeScript

```bash
FILE_GLOB="*.ts"
FORBIDDEN_PATTERNS=(
  '^export +(async +)?function +(invoke|persist) *\('
  '^export +(async +)?function +runOne *\('
  '^export +(async +)?function +(invoke|persist|runOne)(Item|All|One|Task|Job) *\('
)
```

Note `export` anchor — only public helpers are at risk; module-private
functions are scoped to their file.

### Python

```bash
FILE_GLOB="*.py"
FORBIDDEN_PATTERNS=(
  '^def +(invoke|persist) *\('
  '^def +run_one *\('
  '^def +(invoke|persist|run_one)_(release|task|job|item|all|one) *\('
)
```

Use snake_case to match Python's PEP 8 naming.

### Rust

```bash
FILE_GLOB="*.rs"
FORBIDDEN_PATTERNS=(
  '^pub fn +(invoke|persist) *\('
  '^pub fn +run_one *\('
  '^pub fn +(invoke|persist|run_one)_(release|task|job|item|all|one) *\('
)
```

Anchor on `pub fn` — private `fn`s are file-scoped.

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| Pattern matches inside string literals or comments | Use a language-aware parser, OR keep patterns anchored to `^func` / `^def` so leading whitespace excludes string contents |
| Router pattern (`executeRelease`) accidentally matches verb+noun | Whitelist the router prefix in your `FORBIDDEN_PATTERNS` exclusion or anchor the verb list to non-router verbs only |
| `grep` exits non-zero when no matches found, killing `set -e` | Wrap with `|| true` and check the captured variable for non-empty content |
| Pattern only checks one source dir | Loop over multiple dirs OR pass directories as positional args |

---

## Reference Implementation

See `.github/scripts/check-cmd-naming.sh` in `gitmap-v6` for the
production Go implementation that originated this pattern.

---

## Cross-References

- [00-overview.md](./00-overview.md) — Pattern inventory
- [02-grandfather-baseline-naming.md](./02-grandfather-baseline-naming.md) — Companion pattern for legacy code
- [Coding Guidelines — Naming](../../02-coding-guidelines/01-cross-language/15-master-coding-guidelines/01-naming-and-database.md)

---

*Forbidden-name guard — v1.0.0 — 2026-04-21*