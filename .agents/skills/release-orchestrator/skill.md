---
name: release-orchestrator
description: >-
  Execute full automated release orchestration, semantic version bumping, branch management, and tag creation using Python scripts.
---

# Automated Release Orchestrator & Branch Lifecycle — Release Management

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## RULE 0: Version Calculation Standard

1. Read canonical version source (`version.json` or `package.json`).
2. Default bump tier is MINOR: `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`. PATCH MUST reset to 0.
3. Only bump PATCH if user explicitly specified `patch`.
4. Only bump MAJOR if user explicitly specified `major` or breaking change.
5. State previous and new versions explicitly in output before touching files.

## Master Architecture: Heavy-Lifting Release Script (`03-ai-scripts/29-release-orchestrator.py`)

All release operations (version bumping, commit creation, release branching, git tagging, and branch reversion) MUST be executed through the centralized heavy-lifting script in the AI scripts directory:

```bash
python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major>
```

## Pre-Release Validation & Quality Gates (Mandatory)

1. **Mandatory Pre-Release Unit Tests & CI/CD Verification:** Execute `python 03-ai-scripts/06-cicd-local-runner.py --run-tests` and verify all unit test suites, AST checks, and quality gates pass 100% green (`exit 0`).
2. **Test Inventory Validation:** Validate `.lovable/temp/recent-file-changes.json` against `.lovable/test-inventory.json` ensuring all associated tests pass prior to cutting the release.

## Git Release Lifecycle & Original Branch Invariant

1. Record original branch (`git rev-parse --abbrev-ref HEAD`).
2. Bump SemVer across `version.json`, `package.json`, `changelog.md`, `readme.md`.
3. Stage and commit on current branch: `release: vX.Y.Z <scope>`.
4. Create release branch: `release/vX.Y.Z`.
5. Create annotated git tag: `vX.Y.Z`.
6. Push release branch and tag to remote.
7. MANDATORY REVERT: Checkout back to `original_branch`.
