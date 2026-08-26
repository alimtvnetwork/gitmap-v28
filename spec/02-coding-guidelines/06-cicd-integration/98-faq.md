# FAQ — Linter Pack Consumer Questions

> **Version:** 1.0.0
> **Updated:** 2026-04-19
> **Status:** Active

Answers to the four questions every team asks within their first week
of adopting the `linters-cicd` pack. Bookmark this doc and link it from
your repo's `CONTRIBUTING.md`.

---

## 1. How do I suppress a single finding?

Findings are suppressed **at the source line** with an inline comment
the orchestrator recognizes. The comment must include the rule ID and
a human-readable reason (no reason = ignored, the suppression is
treated as invalid and the finding still fires).

### Go

```go
if outer { // codeguidelines:disable=CODE-RED-001 — legacy adapter, slated for removal in #482
    if inner {
        doThing()
    }
}
```

### TypeScript

```ts
// codeguidelines:disable-next-line=CODE-RED-003 — vendor SDK string, cannot extract
const url = "https://api.thirdparty.com/v1";
```

### PHP

```php
// codeguidelines:disable=CODE-RED-002 — WordPress hook signature requires this name
function isnt_admin($user) { ... }
```

**Rules:**

- `disable=` suppresses the line the comment is on.
- `disable-next-line=` suppresses the next non-blank line.
- Multiple rule IDs allowed: `disable=CODE-RED-001,CODE-RED-003`.
- Reason text after `—` is **required** and surfaced in the SARIF
  `suppressions[]` block so reviewers can see why.

Suppressions without a reason are themselves reported as a STYLE
finding (`STYLE-099 SuppressionWithoutReason`, severity `warning`).
The original finding still fires — invalid suppressions never silently
hide a violation. STYLE-099 is **synthetic**: it has no per-language
script, it is injected by `scripts/post-process.py` after the merge
step, and it scans every source file under `--path` regardless of
whether other rules produced findings in that file.

---

## 2. How do I baseline existing violations on legacy code?

Adopting the pack on a 200 kLOC codebase will surface thousands of
findings. Don't try to fix them all at once — **baseline** them so only
**new** violations break the build.

```bash
bash linters-cicd/run-all.sh --path . --format sarif \
  --output .codeguidelines-baseline.sarif
git add .codeguidelines-baseline.sarif
```

Then in CI:

```bash
bash linters-cicd/run-all.sh --path . \
  --baseline .codeguidelines-baseline.sarif \
  --format sarif --output coding-guidelines.sarif
```

The orchestrator subtracts every `(ruleId, uri, startLine,
messageDigest)` tuple present in the baseline from the new run. Only
**new** findings appear in the output and only those affect the exit
code.

### Refreshing the baseline

When you fix legacy violations, regenerate the baseline so the count
shrinks over time:

```bash
bash linters-cicd/run-all.sh --path . --refresh-baseline \
  --baseline .codeguidelines-baseline.sarif
```

A weekly CI job that runs `--refresh-baseline` and opens a PR if the
file changed is a healthy pattern — it visualizes legacy debt
shrinking commit by commit.

---

## 3. How do I run only one rule (or one language)?

Three filters are supported, all combinable.

### Single rule

```bash
bash linters-cicd/run-all.sh --path . --rules CODE-RED-001
```

### Multiple rules

```bash
bash linters-cicd/run-all.sh --path . --rules CODE-RED-001,CODE-RED-004
```

### Single language

```bash
bash linters-cicd/run-all.sh --path . --languages typescript
```

### Combine

```bash
bash linters-cicd/run-all.sh --path . \
  --rules CODE-RED-001,CODE-RED-002 \
  --languages go,typescript \
  --format text
```

### Skip a rule globally

For permanent project-wide opt-outs (rare — discuss with your tech
lead first), use `--exclude-rules`:

```bash
bash linters-cicd/run-all.sh --path . \
  --exclude-rules STYLE-002
```

A `.codeguidelines.toml` at repo root pins these defaults so every
contributor and the CI use the same flags:

```toml
[run]
exclude-rules = ["STYLE-002"]
languages = ["go", "typescript", "php"]
fail-on-warning = false
```

---

## 4. How do I pin a specific linter pack version?

Pin in **all three** distribution channels — never use `@latest` or
`@main` in production CI.

### GitHub composite Action

```yaml
- uses: alimtvnetwork/coding-guidelines-v24/linters-cicd@v3.9.0
```

Or pin to the **commit SHA** for maximum supply-chain hardening:

```yaml
- uses: alimtvnetwork/coding-guidelines-v24/linters-cicd@a3c91f2b8d4e5f6789012345678901234567890a
```

### `install.sh` one-liner

```bash
LINTERS_VERSION=v3.9.0 \
  curl -fsSL https://github.com/alimtvnetwork/coding-guidelines-v24/releases/download/v3.9.0/install.sh \
  | bash
```

The installer verifies the SHA-256 of the downloaded ZIP against the
checksum embedded in the release notes — see
[`05-distribution.md`](./05-distribution.md).

### ZIP download

```bash
VERSION=v3.9.0
curl -fsSLO "https://github.com/alimtvnetwork/coding-guidelines-v24/releases/download/${VERSION}/coding-guidelines-linters-${VERSION}.zip"
unzip -q "coding-guidelines-linters-${VERSION}.zip"
```

### Verifying the pinned version at runtime

Every check script honours `--version`:

```bash
$ python3 linters-cicd/checks/nested-if/go.py --version
coding-guidelines/nested-if 3.9.0
```

CI logs print the orchestrator banner with the pack version on every
run so audit trails always show what enforced what:

```
🔍 coding-guidelines linters-cicd v3.9.0
```

### Upgrade policy

- **Patch** (`v3.9.0 → v3.9.1`) — bug fixes only, safe to bump
  automatically via Dependabot.
- **Minor** (`v3.9.0 → v3.10.0`) — new rules or languages added with
  default OFF for existing repos. Manual opt-in required.
- **Major** (`v3.9.0 → v4.0.0`) — rule defaults changed or removed.
  Read the migration guide before bumping.

---

## Cross-References

- [SARIF Contract](./01-sarif-contract.md)
- [Rules Mapping](./06-rules-mapping.md)
- [Distribution](./05-distribution.md)
- [Performance](./07-performance.md)
- [Acceptance Criteria](./97-acceptance-criteria.md)

---

*Part of [CI/CD Integration](./00-overview.md)*
