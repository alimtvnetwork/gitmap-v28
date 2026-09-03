# Shared Pipeline Conventions

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Platform

All pipelines run on **GitHub Actions** using `ubuntu-latest` runners. No self-hosted runners are required for standard workflows.

---

## Action and Tool Versioning

All GitHub Actions and external tools MUST be pinned to exact version tags. Using `@latest` or `@main` is **prohibited** — it breaks reproducibility and can introduce breaking changes silently.

| Rule | Example |
|------|---------|
| ✅ Pinned tag | `actions/checkout@v6` |
| ✅ Pinned tool | `golangci-lint@v1.64.8` |
| ❌ Floating tag | `actions/checkout@main` |
| ❌ Latest keyword | `go install tool@latest` |

---

## Trigger Patterns

### CI Pipeline (Continuous Integration)

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

Runs on every push or pull request targeting the main branch.

### Release Pipeline

```yaml
on:
  push:
    branches:
      - "release/**"
    tags:
      - "v*"
```

Runs when code is pushed to a `release/*` branch or when a `v*` tag is created.

### Scheduled Scans

```yaml
on:
  schedule:
    - cron: "0 9 * * 1"  # Every Monday at 9:00 UTC
  workflow_dispatch:       # Manual trigger from GitHub UI
```

Used for periodic vulnerability scanning or dependency audits.

---

## Concurrency Control

### CI Pipelines — Cancel Superseded Runs

```yaml
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

**Why**: If two pushes land on the same branch in quick succession, the first run is cancelled to save compute. Release branches are **exempt** — every release commit must run to completion.

### Release Pipelines — Never Cancel

```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
```

**Why**: Every release commit must produce a complete set of artifacts. Cancelling a release build could leave a partial release in an inconsistent state.

---

## Permissions

| Pipeline | Permission | Reason |
|----------|-----------|--------|
| CI | `contents: read` | Only reads source code |
| Release | `contents: write` | Creates GitHub Releases and uploads assets |

Always use the **minimum permissions** required. Never grant `write` access in CI-only workflows.

---

## Version Resolution

The version is derived from the Git ref — never hardcoded in workflow files or source code.

```bash
if [[ "$GITHUB_REF" == refs/tags/* ]]; then
  VERSION="${GITHUB_REF_NAME}"
elif [[ "$GITHUB_REF" == refs/heads/release/* ]]; then
  VERSION="${GITHUB_REF_NAME#release/}"
else
  echo "::error::Unexpected ref: $GITHUB_REF"
  exit 1
fi
echo "version=$VERSION" >> "$GITHUB_OUTPUT"
```

| Ref Pattern | Example | Resolved Version |
|-------------|---------|-----------------|
| Tag | `refs/tags/v1.2.3` | `v1.2.3` |
| Release branch | `refs/heads/release/v1.2.3` | `v1.2.3` |
| CI build | `refs/heads/main` | `dev-<sha>` (short SHA) |

---

## Checksum Generation

All release pipelines MUST generate SHA-256 checksums for every release asset:

```bash
cd dist/
sha256sum * > checksums.txt
```

The `checksums.txt` file is included as a release asset for download verification.

---

## Node.js Compatibility

For workflows using JavaScript-based GitHub Actions, set this environment variable to ensure compatibility with newer Node.js runtimes:

```yaml
env:
  FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
```

---

## Filename Enforcement

Markdown filenames MUST be lowercase kebab-case. The pipeline enforces this:

```bash
VIOLATIONS=$(find . -name '*.md' \
  -not -path './node_modules/*' \
  -not -path './.git/*' \
  | grep '[A-Z]' || true)
if [ -n "$VIOLATIONS" ]; then
  echo "::error::Uppercase .md filenames found"
  exit 1
fi
```

---

## Working Directory Rules

Never use `cd` in CI steps to change directories. Use the `working-directory` key instead:

```yaml

# ✅ Correct

- name: Run lint
  working-directory: my-module
  run: go vet ./...

# ❌ Wrong

- name: Run lint
  run: cd my-module && go vet ./...
```

---

## Artifact Retention

| Context | Retention |
|---------|-----------|
| Inter-job artifacts (CI) | 1 day |
| Test results and coverage | 7 days |
| Build artifacts (CI) | 14 days |
| Release assets | Permanent (attached to GitHub Release) |

---

## Constraints

1. All tool installs use exact version tags — `@latest` is prohibited
2. Version is resolved from the Git ref, never hardcoded
3. Checksums are generated for all release assets
4. Minimum permissions are always used
5. `working-directory` is used instead of `cd`
6. Validate directories before operating: `test -d "$DIR" || exit 1`

---

## Cross-References

- [GitHub Release Standard](./04-github-release-standard.md) — Release body assembly and pre-release detection
- [Vulnerability Scanning](./06-vulnerability-scanning.md) — Scanning classification rules
- [Install Script Generation](./07-install-script-generation.md) — Placeholder strategy, checksum verification
- [Code Signing](./10-code-signing.md) — SignPath integration, feature-flag gating
- [Browser Extension Deploy](./01-browser-extension-deploy/01-index.md) — Node.js pipeline archetype
- [Go Binary Deploy](./02-go-binary-deploy/01-index.md) — Go pipeline archetype

---

*Shared conventions — updated: 2026-04-10*
