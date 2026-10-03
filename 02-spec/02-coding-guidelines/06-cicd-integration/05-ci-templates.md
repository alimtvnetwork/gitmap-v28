# CI Templates Inventory (AI Execution Prompt)

> **/goal** Provide standardized, copy-paste ready CI/CD workflow templates across GitHub Actions, GitLab CI, Azure DevOps, Bitbucket, and Jenkins.
> **/learn** Master single-command composite actions, language target filters, unified SARIF PR annotations, and zero-storage CI compliance.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Provide verified CI template configurations for all major CI platforms calling `run-all.sh`.
- [ ] `/learn` Never trigger auto-detection blindly; pass explicit `--languages` flags in CI workflow definitions.
- [ ] `/goal` Configure SARIF report ingestion so rule violations render inline as code annotations on pull requests.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

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

*Part of [CI/CD Integration](./readme.md)*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CI-005: Cross-Platform CI Templates and Invocation Standards

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.
