# 06 — Message Pipeline and Function-Intel Detection

## 6.1 Message build pipeline (deterministic order)

Given the source commit's `OriginalMessage`, the final commit message
is built by applying these stages in this exact order. Each stage's
output is the next stage's input.

1. **Strip rules** — for every `MessageRule`, drop any line that
   matches its `Kind`:
   - `StartsWith` — `strings.HasPrefix(line, Value)`
   - `EndsWith`   — `strings.HasSuffix(line, Value)`
   - `Contains`   — `strings.Contains(line, Value)`
   Match is line-level, case-sensitive. After stripping, collapse
   consecutive blank lines to one and trim trailing whitespace.
2. **Override decision** — if `OverrideMessages` is non-empty AND
   (`OverrideOnlyWeak == false` OR title's first word ∈ `WeakWords`):
   replace the entire current message with one randomly picked entry
   from `OverrideMessages` (PRNG seeded per §3.4).
3. **Title prefix / suffix** — first line becomes
   `TitlePrefix + line + TitleSuffix`. Empty values are no-ops.
4. **Body prefix / suffix** — if `MessagePrefix` is a CSV pool, pick
   one entry randomly; prepend it as a new line above the body.
   Symmetric for `MessageSuffix`.
5. **Function-intel block** — when `FunctionIntel.IsEnabled`, append
   the per-file added-function block (see §6.3) separated by one
   blank line.
6. **Empty check** — if the resulting message has zero non-blank
   characters → `SkipLog(EmptyAfterMessageRules)`.

## 6.2 Weak-word matching (resolves the override-only-weak intent)

- Compare against the title's first **word**, lowercased, with
  surrounding punctuation stripped (`.,:;!?`).
- Default `WeakWords`: `change`, `update`, `updates`. Override via
  `--weak-words` or profile.
- Examples:
  - `"Update README"` → first word `update` → matches → override fires.
  - `"Updates: bump deps"` → first word `updates` (colon stripped) → matches.
  - `"Refactor parser"` → first word `refactor` → no match → keep original.

## 6.3 Function-intel block format

Per-file, sorted by relative path ascending. Only files that
(a) were ADDED in this commit OR (b) had ADDED top-level function
declarations vs the immediate parent commit are included.

```
- src/foo.go
  - added: NewParser, parseHeader
- src/bar.ts
  - added: useDebounce
```

## 6.4 Per-language detection rules (`FunctionIntelLanguage`)

Detection is regex / lightweight-AST per language. NO full parser
required for v1. The match is "added top-level declaration line" —
lines present in the new tree's file but not in the parent's.

| Language     | File ext            | "Top-level declaration" pattern (informative)            |
|--------------|---------------------|----------------------------------------------------------|
| `Go`         | `.go`               | `^func (\([^)]*\) )?[A-Z]?[A-Za-z0-9_]+\s*\(`            |
| `JavaScript` | `.js`, `.mjs`, `.cjs` | `^(export\s+)?(async\s+)?function\s+[A-Za-z0-9_$]+\s*\(` |
| `TypeScript` | `.ts`, `.tsx`       | Same as JS plus `^(export\s+)?const\s+[A-Za-z0-9_$]+\s*=\s*(async\s*)?\(`|
| `Rust`       | `.rs`               | `^(pub(\([^)]*\))?\s+)?fn\s+[a-z_][A-Za-z0-9_]*\s*[<(]`  |
| `Python`     | `.py`               | `^def\s+[a-z_][A-Za-z0-9_]*\s*\(`                        |
| `Php`        | `.php`              | `^(public\|protected\|private\s+)?function\s+[A-Za-z0-9_]+\s*\(` |
| `Java`       | `.java`             | `^\s*(public\|protected\|private)?\s*(static\s+)?[A-Za-z0-9_<>\[\]]+\s+[A-Za-z0-9_]+\s*\(.*\)\s*\{` |
| `CSharp`     | `.cs`               | Same as Java.                                            |

- Detection is best-effort. A regex miss is NOT a failure. A regex
  panic / parser crash IS a failure → exit `CommitInExitFunctionIntel`.
- Renamed / moved functions are explicitly **out of scope for v1**
  (per §1.5 #1).
- Multi-language single file (e.g. `.html` with embedded JS): scan as
  the language matching the file extension only.

## 6.5 Implementation hints (non-normative)

- Per-language detector lives in its own file:
  `gitmap-v27/cmd/commitin/funcintel/<language>.go`.
- Each detector exposes `Detect(prevSrc, newSrc string) []string`
  returning new declaration names, sorted ascending, deduped.
- A registry map `FunctionIntelLanguage → Detector` allows
  dispatch without `if`/`switch` chains in the caller.

## 6.6 Author identity precedence (recap from §02)

1. `--author-name` + `--author-email` (both required together).
2. Profile's `Author.Name` + `Author.Email` (both required together).
3. OS Git default (`git config user.name` / `user.email`).
4. Source commit's `AuthorName` / `AuthorEmail`.

`AppliedAuthorName` / `AppliedAuthorEmail` in `RewrittenCommit` records
whichever tier won. `AppliedAuthorDate` and `AppliedCommitterDate`
ALWAYS equal source — author override never changes dates (INV-02).