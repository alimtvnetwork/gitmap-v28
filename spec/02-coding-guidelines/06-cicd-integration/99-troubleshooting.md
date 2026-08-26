# Troubleshooting — Linter Pack Operations

> **Version:** 1.0.0
> **Updated:** 2026-04-19
> **Status:** Active

The five issues that account for ~95% of support tickets opened against
the `linters-cicd` pack, with diagnosis steps, copy-pasteable fixes,
and root-cause notes so the same problem does not resurface in your
pipeline a quarter later.

If your symptom is not listed here, file an issue with the output of
`bash linters-cicd/run-all.sh --path . --format text` attached.

---

## 1. `python3: command not found`

### Symptom

```
::error::python3 is required (>= 3.10)
##[error]Process completed with exit code 2.
```

The orchestrator hard-fails before any check runs. CI logs show the
banner was never printed.

### Root cause

The runner image does not have a Python 3 binary on `$PATH`, or the
binary is named `python` (Windows / Alpine) instead of `python3`.

### Fix — GitHub-hosted runners

`ubuntu-latest`, `ubuntu-22.04`, and `ubuntu-24.04` all ship Python
3.10+ pre-installed. If you see this error on a hosted runner you are
almost certainly on `windows-latest` or `macos-latest` without the
setup step:

```yaml
- uses: actions/setup-python@v5
  with:
    python-version: '3.11'
- uses: alimtvnetwork/coding-guidelines-v24/linters-cicd@v3.11.0
```

### Fix — self-hosted runners

```bash

# Debian / Ubuntu

sudo apt-get update && sudo apt-get install -y python3 python3-venv

# RHEL / Amazon Linux

sudo dnf install -y python3.11

# Alpine (note: 'python3' symlink may be missing)

apk add --no-cache python3 && ln -sf /usr/bin/python3 /usr/local/bin/python3
```

### Fix — Windows runners

The pack's `run-all.sh` requires bash; on Windows runners use
`shell: bash` and `python` is auto-aliased by `actions/setup-python`:

```yaml
- shell: bash
  run: bash linters-cicd/run-all.sh --path .
```

### Verification

```bash
python3 --version    # must print 3.10.0 or higher
which python3        # must resolve, not error
```

---

## 2. `tree-sitter` install failures (Phase 2+)

### Symptom

```
ERROR: Failed building wheel for tree-sitter
fatal error: 'Python.h' file not found
```

Phase 1 (Go + TypeScript + PHP) is **pure-Python regex** and has no
native dependencies. This error only appears in **Phase 2+** when a
check imports the `tree-sitter` Python binding for AST-grade parsing.

### Root cause

`tree-sitter` is a C extension. Building from source needs:

- A C compiler (`gcc` or `clang`)
- Python development headers (`python3-dev` / `python3-devel`)
- A modern `pip` (≥ 23.0) so the prebuilt wheel is preferred over a
  source build

### Fix — install the wheel, not source

```bash
python3 -m pip install --upgrade pip
python3 -m pip install --only-binary=:all: tree-sitter==0.21.3
```

The `--only-binary=:all:` flag refuses to fall back to a source build —
if no wheel exists for your platform, `pip` errors out fast instead of
hanging on a 90-second compile that ultimately fails.

### Fix — install build tools (last resort)

```bash

# Debian / Ubuntu

sudo apt-get install -y build-essential python3-dev

# RHEL / Fedora

sudo dnf install -y gcc python3-devel

# Alpine

apk add --no-cache build-base python3-dev
```

### Fix — pin the binding version per language pack

`tree-sitter` changed its public API at 0.21 and again at 0.22. Pin
the binding **and** every grammar to mutually-compatible majors in
`linters-cicd/requirements.txt`:

```
tree-sitter==0.21.3
tree-sitter-go==0.21.0
tree-sitter-typescript==0.21.2
tree-sitter-php==0.22.5
```

### Fix — vendor the wheels

For air-gapped runners, drop the wheels into
`linters-cicd/vendor/wheels/` and install offline:

```bash
python3 -m pip install --no-index --find-links linters-cicd/vendor/wheels \
  tree-sitter tree-sitter-go tree-sitter-typescript
```

### Verification

```bash
python3 -c "from tree_sitter import Language, Parser; print('ok')"
```

---

## 3. SARIF too large for GitHub upload (10 MB limit)

### Symptom

```
##[error]Failed to upload SARIF results: file size 12.4 MB exceeds the 10 MB limit
##[warning]The 'github/codeql-action/upload-sarif' step did not upload results
```

GitHub's Code Scanning API caps a single SARIF upload at **10 MB** and
**5 000 results per run**. Large monorepos blow past both on first
adoption.

### Root cause

Three multipliers compound:

1. Every finding is a JSON object ~600 bytes serialized.
2. The `partialFingerprints` and `properties` blocks add ~200 bytes
   each when present.
3. Pretty-printing with `indent=2` (the default in
   `scripts/merge-sarif.py`) inflates by ~30%.

10 000 findings × 800 bytes × 1.3 ≈ 10.4 MB → over the limit.

### Fix #1 — baseline first, fix later

The fastest path back to green CI is to baseline existing violations
so only **new** findings flow to GitHub:

```bash
bash linters-cicd/run-all.sh --path . \
  --refresh-baseline .codeguidelines-baseline.sarif
git add .codeguidelines-baseline.sarif && git commit -m "chore: baseline existing violations"
```

CI then runs:

```bash
bash linters-cicd/run-all.sh --path . \
  --baseline .codeguidelines-baseline.sarif \
  --output coding-guidelines.sarif
```

Typical drop: 10 000 → < 50 findings on the first PR. See
[`98-faq.md`](./98-faq.md#2-how-do-i-baseline-existing-violations-on-legacy-code).

### Fix #2 — split the SARIF by tool

Each `runs[]` entry can be uploaded separately. The `category` field
on `upload-sarif` namespaces them so they merge in the UI:

```yaml
- uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: coding-guidelines-code-red.sarif
    category: code-red
- uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: coding-guidelines-style.sarif
    category: style
```

A future `--split-by severity` flag on `run-all.sh` will automate this
(tracked in `97-acceptance-criteria.md` as `AC-CI-013`).

### Fix #3 — minify the SARIF

Drop pretty-printing for the upload artifact only:

```bash
python3 -c "import json,sys; \
  d=json.load(open('coding-guidelines.sarif')); \
  json.dump(d, open('coding-guidelines.min.sarif','w'), separators=(',',':'))"
```

Reduces size by ~30% with zero schema impact. Keep the pretty version
for human review and upload the minified one.

### Fix #4 — exclude generated code

Generated files (protobuf, GraphQL, OpenAPI clients) flood findings
with no signal. Add to `.codeguidelines.toml`:

```toml
[run]
exclude-paths = ["**/*.pb.go", "**/generated/**", "**/__generated__/**"]
```

### Verification

```bash
ls -lh coding-guidelines.sarif       # < 10 MB
python3 -c "import json; d=json.load(open('coding-guidelines.sarif')); \
  print(sum(len(r['results']) for r in d['runs']))"   # < 5000
```

---

## 4. False-positive triage workflow

### Symptom

A check fires on code you believe is correct. Examples:

- `CODE-RED-003 MagicStrings` flags an HTTP status code constant.
- `CODE-RED-001 NoNestedIf` fires inside a generated `switch`-equivalent.
- `CODE-RED-002 BooleanNaming` flags a domain term (`isnt_admin` for a
  WordPress hook signature).

### Triage decision tree

```
                    ┌───────────────────────────┐
                    │ Is the code actually wrong │
                    │ by the spec's definition?  │
                    └────────┬───────────────────┘
                             │
              ┌──── yes ─────┴──── no ─────┐
              │                             │
              ▼                             ▼
    Fix the code (preferred).      Is the spec wrong, or is
    Closes the finding without     the check wrong?
    debt.                                   │
                                  ┌────── spec ──────── check ──────┐
                                  ▼                                 ▼
                          File a spec issue with             File a check issue
                          repro + proposed wording.          with the source line
                          Once merged, regenerate            and expected behaviour.
                          baseline.                          Suppress locally until
                                                             a fix ships.
```

### Suppress (only after the triage above)

Inline suppression with a **mandatory reason** — see
[`98-faq.md` §1](./98-faq.md#1-how-do-i-suppress-a-single-finding):

```ts
// codeguidelines:disable-next-line=CODE-RED-003 — RFC 7231 status code, not a magic string
const NOT_FOUND = 404;
```

Reasonless suppressions are themselves reported as
`STYLE-099 SuppressionWithoutReason` and the original finding still
fires — no silent hiding is possible.

### Bulk suppress for a known-noisy directory

Use `.codeguidelines.toml` instead of sprinkling comments:

```toml
[run]
exclude-paths = ["legacy/wordpress-hooks/**"]
```

### Filing a check bug

When you open the issue, include:

1. The finding's full SARIF result block (copy from the merged file).
2. A 5-line minimal repro saved as a fixture under
   `linters-cicd/checks/<rule-slug>/_fixtures/false-positive-NNN.<ext>`.
3. The expected outcome (`should not fire because …`).
4. The pack version (`cat linters-cicd/VERSION`).

A check bug without a fixture file is closed automatically — the
fixture becomes the regression test once the bug is fixed.

### Verification

After applying a suppression or exclude, re-run targeting just that
rule to confirm green:

```bash
bash linters-cicd/run-all.sh --path . --rules CODE-RED-003 --format text
```

---

## 5. `.codeguidelines.toml` parse errors

### Symptom

```
::warning::Failed to parse .codeguidelines.toml — falling back to CLI defaults
```

The orchestrator silently falls back to built-in defaults, so checks
still run, but **none of your repo-level overrides are applied**. CI
appears to work but is not enforcing what you configured.

### Root causes (in order of frequency)

#### a) Python < 3.11 on the runner

`scripts/load-config.py` uses `tomllib`, which was added in Python
3.11. On 3.10 the import silently returns `{}` and your config is
ignored.

**Fix:** bump the runner Python (see §1) or pin
`actions/setup-python` to `'3.11'` or newer.

#### b) Quoting — values must be valid TOML strings or arrays

```toml

# ❌ WRONG — bare list, comma-less, unquoted

[run]
languages = go typescript php

# ❌ WRONG — Python-style list with single quotes (TOML allows but mixing types breaks)

exclude-rules = ['STYLE-002', "STYLE-099"]

# ✅ RIGHT

[run]
languages = ["go", "typescript", "php"]
exclude-rules = ["STYLE-002"]
```

#### c) Section header typo

The pack reads only the `[run]` table. These are silently ignored:

```toml
[runs]            # ❌ extra 's'
[ run ]           # ❌ spaces inside brackets
[Run]             # ❌ TOML is case-sensitive
```

#### d) Trailing comma in an inline array (TOML 1.0)

```toml
exclude-rules = ["STYLE-002",]   # ❌ TOML 1.0 forbids — bumped to 1.1 spec only
```

Drop the trailing comma or upgrade to a pack version that targets TOML
1.1 (`v4.0.0+`).

#### e) BOM or CRLF line endings on Windows-edited files

Some editors save UTF-8 with a BOM that confuses `tomllib`. Convert:

```bash
sed -i '1s/^\xEF\xBB\xBF//' .codeguidelines.toml
dos2unix .codeguidelines.toml
```

### Diagnose

The `load-config.py` helper has no `--strict` mode yet (tracked as
`AC-CI-014`). For now, validate manually with a one-liner:

```bash
python3 -c "import tomllib; print(tomllib.loads(open('.codeguidelines.toml').read()))"
```

A `tomllib.TOMLDecodeError` traceback pinpoints the exact line and
column.

### Reference config (copy + customise)

```toml

# .codeguidelines.toml — repo-level defaults for the linter pack.

# CLI flags always override these values.

[run]
languages         = ["go", "typescript", "php"]
rules             = []                    # empty = run all rules in registry
exclude-rules     = ["STYLE-002"]         # warnings we accept project-wide
exclude-paths     = ["**/generated/**"]   # globs, evaluated relative to --path
fail-on-warning   = false                 # only `error` level breaks the build
```

### Verification

```bash
python3 linters-cicd/scripts/load-config.py --config .codeguidelines.toml

# Should print the effective LANGUAGES=, RULES=, EXCLUDE_RULES=, FAIL_ON_WARNING= values

```

---

## Cross-References

- [SARIF Contract](./01-sarif-contract.md) — wire format every check emits
- [FAQ](./98-faq.md) — suppression syntax, baseline workflow, version pinning
- [Performance](./07-performance.md) — timeouts, parallel jobs, walker order
- [Acceptance Criteria](./97-acceptance-criteria.md) — testable AC list
- [Distribution](./05-distribution.md) — install methods, version pinning

---

## Contributors

- **Md. Alim Ul Karim** — Creator & Lead Architect
- **Riseup Asia LLC** — Sponsor

---

*Part of [CI/CD Integration](./00-overview.md)*
