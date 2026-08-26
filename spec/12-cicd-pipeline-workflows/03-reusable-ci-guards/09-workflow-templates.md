# Pattern 09 — GitHub Actions Workflow Templates

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

Patterns 01–08 give a repo the *building blocks* (guards, runner,
config). Most teams still copy-paste 60+ lines of YAML to wire them
into GitHub Actions, then drift over time as different repos edit
their copies. The workflow templates collapse that boilerplate into
**one `uses:` line** and expose every reasonable knob as a typed input.

---

## Three Layers

| Layer | File | When to use |
|-------|------|-------------|
| Composite action | `.github/workflow-templates/action.yml` | Drop into an existing job as one step |
| Reusable workflow | `.github/workflow-templates/ci-guards.reusable.yml` | Call from another workflow with `workflow_call` |
| Starter workflows | `.github/workflow-templates/starters/{go,node,python,rust}.yml` | Copy/paste a complete workflow per language |

All three ultimately invoke `bash scripts/ci-runner.sh` (Pattern 07),
which loads `ci-guards.yaml` (Pattern 08) and dispatches Patterns 01–06.

---

## Composite Action Inputs

| Input | Type | Default | Maps to runner flag |
|-------|------|---------|---------------------|
| `phase` | string | `all` | `--phase` |
| `guard` | string | _(empty)_ | `--guard` |
| `config` | string | `ci-guards.yaml` | `--config` (skipped if missing) |
| `source-dir` | string | _(empty)_ | `--source-dir` |
| `baseline` | string | _(empty)_ | `--baseline` |
| `results-dir` | string | _(empty)_ | `--results-dir` |
| `scripts-dir` | string | `.github/scripts` | `--scripts-dir` |
| `json-summary` | string | _(empty)_ | `--json` |
| `fail-on-violation` | bool | `true` | Post-process; downgrades exit `1` → `0` if `false` |
| `node-version` | string | `20` | `actions/setup-node` |

Outputs: `exit-code`, `summary-path`.

## Reusable Workflow Inputs

The reusable workflow accepts every composite-action input plus three
language-agnostic hints exposed to the underlying scripts as env vars:

| Input | Env var | Purpose |
|-------|---------|---------|
| `languages` | `CI_GUARDS_LANGUAGES` | Comma list (`go,typescript`); consumed by `linters-cicd` and language-aware guards |
| `file-glob` | `CI_GUARDS_FILE_GLOB` | Override per-call file glob |
| `rule-set` | `CI_GUARDS_RULESET` | Tag (`go-default`, `python-strict`, …) the guards can branch on |

Plus `runs-on` (default `ubuntu-latest`) and `node-version`.

---

## Caller Examples

### Composite action (one job)

```yaml
jobs:
  guards:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: ./.github/workflow-templates
        with:
          phase: all
          config: ci-guards.yaml
```

### Reusable workflow (cross-repo)

```yaml
jobs:
  guards:
    uses: alimtvnetwork/coding-guidelines-v24/.github/workflow-templates/ci-guards.reusable.yml@v3.44.0
    with:
      phase: all
      languages: "go"
      file-glob: "*.go"
      rule-set: "go-default"
```

### Single-guard mode

```yaml
- uses: ./.github/workflow-templates
  with:
    guard: forbidden_names
    config: ci-guards.yaml
```

---

## Starter Catalog

| File | Languages | File glob | Rule-set tag |
|------|-----------|-----------|--------------|
| `starters/go.yml` | `go` | `*.go` | `go-default` |
| `starters/node.yml` | `typescript,javascript` | `*.{ts,tsx,js,jsx}` | `node-default` |
| `starters/python.yml` | `python` | `*.py` | `python-default` |
| `starters/rust.yml` | `rust` | `src/**/*.rs` | `rust-default` |

Copy any starter into `.github/workflows/`, edit the four hint inputs
if needed, and you have a working CI in under a minute.

---

## Versioning

Templates ship inside the same repo, version-locked to `package.json`.
Callers SHOULD pin to a release tag (`@v3.44.0`) — not `@main` — so
input contracts cannot drift under them.

| Stability | Reference |
|-----------|-----------|
| **Stable** | Pinned tag (`@v3.44.0`, `@v3`) |
| **Tracking** | `@main` — only for early adopters / CI of the templates themselves |

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| `uses: ./.github/workflow-templates` resolves to nothing | The composite `action.yml` must live in that exact directory; rename = breakage |
| `config: ci-guards.yaml` always loads even when missing | The action checks `[ -f "$config" ]` first; missing config simply skips `--config` |
| Caller wants to ignore violations during migration | Set `fail-on-violation: false` to downgrade exit `1` → `0` (tool errors `2`/`64` still fail) |
| Output artifact missing | The upload step uses `if-no-files-found: ignore`; provide `json-summary` to actually emit one |

---

## Cross-References

- [00-overview.md](./00-overview.md)
- [07-shared-cli-wrapper.md](./07-shared-cli-wrapper.md) — Underlying runner
- [08-config-schema.md](./08-config-schema.md) — `ci-guards.yaml` schema
- [99-ai-implementation-guide.md](./99-ai-implementation-guide.md) — End-to-end decision tree

---

*Workflow templates — v1.0.0 — 2026-04-21*