# Pattern 08 — Unified Config Schema

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

Patterns 01–06 each have their own configuration surface — file globs,
regex patterns, baseline file paths, identifier prefixes, test markers.
The shared CLI wrapper (Pattern 07) unifies their **invocation**, but
each repo still needs to set the right inputs for each guard.

Hand-wiring those inputs across multiple workflow YAML files,
per-script env vars, and CLI flags scatters configuration in a way
that's hard to audit and impossible to share. The unified config
schema collapses every per-guard knob into **one file** at the repo
root that both humans and AI assistants can read in seconds.

---

## File Format

- Filename: `ci-guards.yaml` (canonical) or `ci-guards.yml` / `ci-guards.json`
- Location: repository root (override with `--config <path>`)
- Encoding: UTF-8
- Comments: `# …` (YAML only)

The YAML parser in `scripts/ci-config.mjs` accepts a deliberately
small subset (top-level keys, nested maps, scalar values, simple
sequences). Anchors, aliases, multi-line scalars, flow style, and
merge keys are NOT supported. The schema is designed so you never
need them.

---

## Schema (canonical)

```yaml
version: 1                     # required; integer schema version
language: go                   # optional; one of: go|node|python|rust|polyglot
scripts_dir: .github/scripts   # optional; default: .github/scripts

forbidden_names:               # Pattern 01 (optional section)
  enabled: true
  source_dir: <path>           # required when enabled
  file_glob: "*.go"            # required when enabled
  patterns: [<regex>, ...]     # required when enabled
  remediation_examples: [<line>, ...]

naming_baseline:               # Pattern 02
  enabled: true
  source_dir: <path>
  baseline_file: <path>
  allowed_prefix_regex: "^(Cmd|Msg|Err|Flag|Default)"
  rename_examples: [<line>, ...]

collisions:                    # Pattern 03
  enabled: true
  file_glob: <glob>
  ident_regex: "^[A-Z][A-Za-z0-9_]*"
  rawstring_delim: "`"         # one char
  block_open_regex: "^\\s*const\\s*\\("
  block_close_regex: "^\\s*\\)"

lint_diff:                     # Pattern 04
  enabled: true
  current_report: <path>
  baseline_report: <path>
  normalizer: golangci         # golangci|eslint|ruff|clippy
  cap_seeding_findings: 50

lint_suggest:                  # Pattern 05
  enabled: true
  current_report: <path>
  baseline_report: <path>
  comment_marker: "<!-- repo-lint-suggestions -->"
  suggester_table: golangci

test_summary:                  # Pattern 06
  enabled: true
  results_dir: <path>
  shard_pattern: "test-results-*"
  output_filename: test-output.txt
  pass_marker: "^--- PASS:"
  fail_marker: "^--- FAIL:"
  reason_keywords: "expected|got|Error|FAIL|panic|undefined|mismatch"
  reason_max_lines: 10
```

Every guard section is **optional**; an absent section means that
guard runs with documented defaults (or is skipped if no defaults
make sense for the project).

---

## Validation Rules

| Rule | Failure |
|------|---------|
| `version` is an integer | exit `1` with `top-level 'version' must be an integer` |
| `language` (if present) is one of the supported values | exit `1` with `'language' must be one of …` |
| Each guard section is a map (not list/scalar) | exit `1` with `'<guard>' must be a map` |
| `enabled` (if present) is `true`/`false` | exit `1` with `'<guard>.enabled' must be true or false` |
| File path itself unreadable | exit `2` with `config file not found` |
| Unsupported extension | exit `2` with `unsupported config extension` |

Validation is intentionally loose on the per-guard fields — those are
enforced by the underlying guard scripts at run time so that the
config loader stays small and language-agnostic.

---

## Loader Contract (`scripts/ci-config.mjs`)

| Aspect | Value |
|--------|-------|
| Invocation | `node scripts/ci-config.mjs --config <path> [--emit env\|json\|summary]` |
| `--emit env` | Shell `export VAR=…` lines (default; safe to `eval`) |
| `--emit json` | Resolved config as pretty-printed JSON |
| `--emit summary` | Human-readable per-guard table |
| Exit `0` | Config valid, output emitted |
| Exit `1` | Validation errors printed to stderr |
| Exit `2` | Tool error (file unreadable, parse failure) |
| Exit `64` | Usage error (bad flags) |

### Env Variable Naming Convention

Every config field maps to one env var:

```
CI_GUARDS_<SECTION>_<KEY>
```

Examples:

| Config field | Env variable |
|--------------|--------------|
| `version` | `CI_GUARDS_VERSION` |
| `language` | `CI_GUARDS_LANGUAGE` |
| `scripts_dir` | `CI_GUARDS_SCRIPTS_DIR` |
| `forbidden_names.source_dir` | `CI_GUARDS_FORBIDDEN_NAMES_SOURCE_DIR` |
| `naming_baseline.baseline_file` | `CI_GUARDS_NAMING_BASELINE_BASELINE_FILE` |
| `test_summary.fail_marker` | `CI_GUARDS_TEST_SUMMARY_FAIL_MARKER` |

Sequence values become newline-joined strings; booleans become `1`/`0`.

---

## Wiring Into `ci-runner.sh`

The runner accepts `--config <path>`. When provided, it shells out to
the loader, evaluates the env output, and uses those values as
defaults — explicit CLI flags still take precedence:

```bash

# Use only the config

./scripts/ci-runner.sh --phase all --config ci-guards.yaml

# Override one field from the config

./scripts/ci-runner.sh --phase test --config ci-guards.yaml \
  --results-dir ./alt-test-artifacts

# No config — pure CLI flags (Pattern 07 baseline)

./scripts/ci-runner.sh --phase check --source-dir src/lib
```

CI workflow snippet:

```yaml
- name: Run all CI guards
  run: bash scripts/ci-runner.sh --phase all --config ci-guards.yaml --json $RUNNER_TEMP/ci.json
```

---

## Per-Language Templates

### Go

```yaml
version: 1
language: go
forbidden_names:
  enabled: true
  source_dir: cmd
  file_glob: "*.go"
  patterns:
    - "^func +(invoke|persist) *\\("
    - "^func +runOne *\\("
naming_baseline:
  enabled: true
  source_dir: constants
  baseline_file: .baselines/constants.txt
  allowed_prefix_regex: "^(Cmd|Msg|Err|Flag|Default)"
lint_diff:
  enabled: true
  current_report: golangci.json
  baseline_report: .cache/golangci.json
  normalizer: golangci
test_summary:
  enabled: true
  results_dir: ./test-artifacts
  pass_marker: "^--- PASS:"
  fail_marker: "^--- FAIL:"
```

### Node / TypeScript

```yaml
version: 1
language: node
forbidden_names:
  enabled: true
  source_dir: src/handlers
  file_glob: "*.ts"
  patterns:
    - "^export +(async +)?function +(invoke|persist) *\\("
lint_diff:
  enabled: true
  current_report: eslint.json
  baseline_report: .cache/eslint.json
  normalizer: eslint
test_summary:
  enabled: true
  results_dir: ./vitest-results
  pass_marker: "^ ✓ "
  fail_marker: "^ ✗ "
```

### Python

```yaml
version: 1
language: python
naming_baseline:
  enabled: true
  source_dir: pkg
  baseline_file: .baselines/constants.txt
  allowed_prefix_regex: "^(CFG_|MSG_|ERR_|FLAG_|DEFAULT_)"
lint_diff:
  enabled: true
  current_report: ruff.json
  baseline_report: .cache/ruff.json
  normalizer: ruff
test_summary:
  enabled: true
  results_dir: ./pytest-results
  pass_marker: " PASSED$"
  fail_marker: " FAILED$"
```

### Rust

```yaml
version: 1
language: rust
collisions:
  enabled: true
  file_glob: "src/**/*.rs"
  ident_regex: "^pub (const|static) +([A-Z][A-Z0-9_]*)"
lint_diff:
  enabled: true
  current_report: clippy.json
  baseline_report: .cache/clippy.json
  normalizer: clippy
test_summary:
  enabled: true
  results_dir: ./cargo-test-results
  pass_marker: " ... ok$"
  fail_marker: " ... FAILED$"
```

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| YAML feature outside the supported subset (anchors, multi-line scalars) | Use the documented schema only — file an issue if you need more |
| Regex backslashes consumed by YAML string parsing | Always use double-quoted strings and escape `\\` (the example file shows this) |
| Env var name collision with existing shell vars | All emitted vars share the `CI_GUARDS_` prefix to avoid clashes |
| Config edits don't take effect | Confirm you passed `--config <path>` to `ci-runner.sh`; without it the loader is never invoked |
| Per-field validation silently passes | Per-field validation lives in the underlying guards by design; loader only checks structural shape |

---

## Reference Implementation

- Loader: `scripts/ci-config.mjs` (zero-dep Node, ~340 lines)
- Example config: `ci-guards.example.yaml`
- Wrapper integration: `scripts/ci-runner.sh --config <path>`

---

## Cross-References

- [00-overview.md](./00-overview.md)
- [07-shared-cli-wrapper.md](./07-shared-cli-wrapper.md) — Wrapper that consumes the env vars
- [99-ai-implementation-guide.md](./99-ai-implementation-guide.md) — Decision tree referencing this schema

---

*Unified config schema — v1.0.0 — 2026-04-21*