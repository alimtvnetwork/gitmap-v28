# GitMap Release Method

Canonical documentation for GitMap release orchestration, version bumping, and changelog synchronization.

## Release Process
1. Bump minor version in `version.json` using `python 03-ai-scripts/37-bump-version.py`.
2. Sync version pins across `package.json`, `readme.md`, `what-to-read.md`, and Go constants.
3. Validate all quality gates with `python 03-ai-scripts/06-cicd-local-runner.py`.
4. Create release branch and tag via `python 03-ai-scripts/29-release-orchestrator.py`.
