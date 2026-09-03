# Pattern 05 — Actionable Lint Suggestions

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

A linter reports `[gocritic] paramTypeCombine: parameters can be
combined`. The contributor sees this and has no idea what change to
apply. They Google the rule, read the docs, infer the fix, then edit.

The actionable-suggestions pattern produces a **PR comment** that
pairs each NEW finding with:

1. A human-readable explanation of the rule
2. A copy-pasteable code-change template (in fenced Markdown)
3. A permalink to the offending line on GitHub

The result: most linter findings can be resolved in seconds without
leaving the PR review UI.

---

## Algorithm

```
INPUT:  current.json, baseline.json (same as Pattern 04)
OUTPUT: Markdown to GITHUB_STEP_SUMMARY (and PR comment)

1. new_findings = current - baseline   # reuse Pattern 04 normalizer
2. cap at 50 entries when seeding (no baseline) so the comment is
   readable; full list still in the artifact

3. group findings by file
4. for each finding:
     (title, fix_block) = SUGGESTERS[linter](message)
     emit Markdown bullet:
        - **L{line}** · `[{linter}]` — {message}
          - 📍 [permalink]
          - 💡 **Suggested fix:** {title}
              {fix_block}

5. wrap output in HTML sentinel comment so the PR-comment action
   can replace the previous comment in place (no stacking)
```

The `SUGGESTERS` table is **the** per-language artifact. It maps a
linter name to a function that takes the raw message and returns a
`(title, fix_block)` tuple. Unknown linters fall back to a generic
"consult the docs" suggestion so the comment is never blank.

---

## Why a Separate Script (Not Merged with Pattern 04)

| Concern | Pattern 04 (gate) | Pattern 05 (suggest) |
|---------|------------------|----------------------|
| Role | Decides exit code | Generates Markdown |
| Auditability | Must stay tiny + trivially reviewable | Suggestion table grows freely |
| Side effects | None | Writes to step summary, posts PR comment |
| Failure tolerance | Strict | Loose (a bad suggestion shouldn't fail the build) |

They share the JSON loader by **duplication** (intentionally — no
Python import dependency between the two scripts so each runs from
any working directory).

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `python3 lint-suggest.py --current cur.json --baseline base.json --out $GITHUB_STEP_SUMMARY` |
| Exit `0` | Always — this script never gates the build |
| Output | Appends Markdown to `--out` (default `$GITHUB_STEP_SUMMARY`) |
| PR comment | Read by `peter-evans/create-or-update-comment` action |
| Sentinel | `<!-- ci-lint-suggestions -->` first line so prior comments are replaced, not stacked |

---

## Suggestion Table — Worked Examples

### golangci-lint family (Go)

| Linter | Title | Fix Template |
|--------|-------|--------------|
| `misspell` | Replace `colour` with `color` (US English) | `- colour\n+ color` |
| `gocritic/paramTypeCombine` | Combine consecutive params of same type | `- func f(a string, b string)\n+ func f(a, b string)` |
| `gocritic/sprintfQuotedString` | Use `%q` instead of `"%s"` | `- fmt.Sprintf(\`KEY="%s"\`, val)\n+ fmt.Sprintf(\`KEY=%q\`, val)` |
| `unused` | Remove the unused symbol | `- // entire declaration` + nolint hint for future-API case |
| `nolintlint` | Remove the unused `//nolint` directive | `- foo() //nolint:gosec\n+ foo()` |
| `gofmt` | Run `gofmt -w` | `gofmt -w path/to/file.go` |
| `errcheck` | Handle the returned error explicitly | `- doThing()\n+ if err := doThing(); err != nil { ... }` |

### ESLint family (Node / TypeScript)

| Rule | Title | Fix Template |
|------|-------|--------------|
| `@typescript-eslint/no-unused-vars` | Remove unused variable | `- const x = …;` |
| `eqeqeq` | Use `===` instead of `==` | `- a == b\n+ a === b` |
| `prefer-const` | Use `const` for never-reassigned variables | `- let x = …\n+ const x = …` |
| `no-console` | Remove or replace with proper logger | `- console.log(...)\n+ logger.debug(...)` |

### Ruff family (Python)

| Rule | Title | Fix Template |
|------|-------|--------------|
| `F401` (unused-import) | Remove unused import | `- from x import y` |
| `E501` (line-too-long) | Wrap line at 88 chars | wrap example |
| `B006` (mutable-argument-default) | Use `None` + `if … is None: …` | template |
| `RUF005` (collection-literal-concat) | Use `[*a, b]` instead of `a + [b]` | template |

### Clippy family (Rust)

| Lint | Title | Fix Template |
|------|-------|--------------|
| `needless_return` | Remove `return` from final expression | `- return x;\n+ x` |
| `redundant_clone` | Drop `.clone()` | `- x.clone()\n+ x` |
| `single_match` | Use `if let` instead of `match` with one arm | template |
| `unused_self` | Convert method to associated function | template |

---

## Output Format

````markdown
<!-- repo-lint-suggestions -->

### golangci-lint — actionable suggestions

Found **2** new finding(s) introduced by this change. Each entry
below maps the warning to a concrete fix you can apply.

#### `pkg/scan/runner.go`

- **L42** · `[errcheck]` — Error return value of `os.Setenv` is not checked
  - 📍 [`pkg/scan/runner.go:42`](https://github.com/owner/repo/blob/<sha>/pkg/scan/runner.go#L42)
  - 💡 **Suggested fix:** Handle the returned error explicitly.

    ```diff
    - doThing()
    + if err := doThing(); err != nil {
    +     return fmt.Errorf("do thing: %w", err)
    + }
    ```

---
_Suggestions are templates, not patches — review before applying. The
full JSON report is attached as the `lint-report` artifact._
````

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| Comment stacks on every push | Use a unique sentinel HTML comment as line 1 + `peter-evans/create-or-update-comment` |
| Markdown table breaks on `\|` in messages | Render bullets, not tables; escape `\|` defensively |
| Suggestion bot crashes on unknown linter | Always wire a generic fallback suggester |
| Comment gets too long | Cap at 50 findings in seeding mode; full data in JSON artifact |
| Permalink uses PR commit not base | Use `github.event.pull_request.head.sha` for the SHA |

---

## Reference Implementation

See `.github/scripts/lint-suggest.py` in reference CLI implementations.

---

## Cross-References

- [01-index.md](./01-index.md)
- [04-baseline-diff-lint-gate.md](./05-baseline-diff-lint-gate.md) — Sibling script that decides which findings are NEW

---

*Actionable lint suggestions — v1.0.0 — 2026-04-21*
