# CI Templates Inventory

> **Version:** 1.0.0
> **Updated:** 2026-04-19

Ready-to-paste workflow files per CI platform. All call the same
`linters-cicd/run-all.sh` and consume the same SARIF output, so behavior
is identical across platforms.

---

## Shipped templates

| Platform | File | Findings surface as |
|----------|------|---------------------|
| GitHub Actions | `linters-cicd/ci/github-actions.yml` | Code Scanning (Security tab) + PR annotations |
| GitHub composite | `linters-cicd/action.yml` | Same as above, one-liner via `uses:` |
| GitLab CI | `linters-cicd/ci/gitlab-ci.yml` | Code Quality MR widget + SAST report |
| Azure DevOps | `linters-cicd/ci/azure-pipelines.yml` | SARIF SAST extension |
| Bitbucket Pipelines | `linters-cicd/ci/bitbucket-pipelines.yml` | Pipeline log + report artifact |
| Jenkins | `linters-cicd/ci/Jenkinsfile` | Warnings-NG plugin |
| Pre-commit hook | `linters-cicd/ci/pre-commit-hook.sh` | Local block before push |

---

## Other platforms (CircleCI, TeamCity, Drone, …)

Not shipped as templates. Users wire them themselves using the contract:

```bash

# 1. Install

curl -fsSL https://github.com/alimtvnetwork/coding-guidelines-v24/releases/latest/download/install.sh | bash

# 2. Run

./linters-cicd/run-all.sh --path . --format sarif --output coding-guidelines.sarif

# 3. Upload artifact + fail build on exit code 1

```

---

## GitHub composite Action - usage

```yaml
- uses: alimtvnetwork/coding-guidelines-v24/linters-cicd@v3.9.0
  with:
    path: .
    languages: go,typescript,python,rust # specify languages here
    severity: error              # optional, default: error
    fail-on-warning: false       # optional, default: false
```

---

## 🤖 AI Agent Instructions: Connecting Linters for Specific Languages

When an AI is integrating these linters into a CI/CD pipeline, it MUST explicitly define the target languages. This is crucial for avoiding auto-detection errors and speeding up CI runs.

### How an AI should integrate this:

1. **Identify the languages** used in the codebase (e.g., `go`, `typescript`, `python`, `rust`, `php`, `csharp`).
2. **Wire the script** to run explicitly for those languages using the `--languages` argument.
3. **Fail the build** if the script exits with a non-zero code.

**Example: Jenkins / GitLab / Bitbucket (Bash)**
```bash

# 1. Install the tools

curl -fsSL https://github.com/alimtvnetwork/coding-guidelines-v24/releases/latest/download/install.sh | bash

# 2. Execute ONLY for specific languages

./linters-cicd/run-all.sh \
  --path . \
  --languages go,typescript,python \
  --format sarif \
  --output coding-guidelines.sarif
```

**Example: GitHub Actions**
```yaml
      - name: Run Coding Guidelines Linters
        uses: alimtvnetwork/coding-guidelines-v24/linters-cicd@v3.9.0
        with:
          path: .
          languages: 'go,typescript,python'
```

---

*Part of [CI/CD Integration](./00-overview.md)*
