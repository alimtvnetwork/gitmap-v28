# CI/CD Prompts (`16-ci-cd`)

This collection contains canonical prompts for diagnosing, creating, fixing, and maintaining CI/CD pipelines, workflows, and automated runners across polyglot stacks.

## Directory Index

| # | Prompt File | Workflow | Description |
| :---: | :--- | :--- | :--- |
| **01** | [`01-ci-cd-fix-tweak.md`](01-ci-cd-fix-tweak.md) | Targeted Smart Testing | Fast surgical CI fixes with 4-part RCA |
| **02** | [`02-ci-cd-fix-with-release-tweak.md`](02-ci-cd-fix-with-release-tweak.md) | Fix & Release Tweak | CI fixes followed by minor release ceremony |
| **03** | [`03-ci-cd-fix.md`](03-ci-cd-fix.md) | Standard CI Fix | Diagnostic and healing pipeline |
| **04** | [`04-create-run-ps1-file.md`](04-create-run-ps1-file.md) | Run & Install Scripts | Dynamic runner (run.ps1/run.sh) and local-install script architecture |
| **05** | [`05-fix-ci-cd-and-run-scripts.md`](05-fix-ci-cd-and-run-scripts.md) | Scripts Remediation | Fixing local CI runner scripts |
| **06** | [`06-ci-cd-fix-with-release.md`](06-ci-cd-fix-with-release.md) | Fix with Release | Complete CI fix and release ceremony |
| **07** | [`07-cicd-pipeline-create.md`](07-cicd-pipeline-create.md) | Pipeline Architecture | Cross-platform Python pipeline automation |
| **08** | [`08-ci-cd-fix-with-n-steps.md`](08-ci-cd-fix-with-n-steps.md) | N-Step Loop | Continuous N-step loop CI/CD remediation |
| **09** | [`09-ci-cd-fix-with-release-n-steps.md`](09-ci-cd-fix-with-release-n-steps.md) | N-Step Release Loop | Continuous N-step remediation with release |
| **10** | [`10-zero-storage-actions-purge.md`](10-zero-storage-actions-purge.md) | Zero-Storage Actions | Purging GitHub Actions artifacts and cache |
| **11** | [`11-ci-cd-fix-gitmap-release.md`](11-ci-cd-fix-gitmap-release.md) | GitMap Self-Healing Release Loop | Autonomous loop with `gitmap pe -t`, 4-part RCA, and minor bump |

## Core Invariants

1. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well (TOTAL BAN on `file:///` URIs and absolute filesystem paths).
2. **Zero-Storage CI/CD Mandate:** Zero routine artifact uploads to GitHub Actions storage.
