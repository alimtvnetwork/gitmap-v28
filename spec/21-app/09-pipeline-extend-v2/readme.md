# Pipeline Extend V2

This folder contains the complete V2 specification for extending the CI/CD pipeline and guiding AI Agents on proper release procedures, error handling, and coding standards. 

If this folder is provided to an AI Agent context, it must implement these exact standards without deviation.

## Table of Contents

| # | File | Topic |
|---|------|-------|
| 01 | [01-ai-release-synchronization.md](01-ai-release-synchronization.md) | PowerShell script and AI branch sync protocols for releases. |
| 02 | [02-changelog-awk-integration.md](02-changelog-awk-integration.md) | How GitHub Actions uses `awk` to extract patch notes. |
| 03 | [03-query-wrapper-python-ts.md](03-query-wrapper-python-ts.md) | The generic `isFail` query wrapper spec for TS and Python. |
| 04 | [04-strict-enum-enforcement.md](04-strict-enum-enforcement.md) | Requirements for eliminating TS string unions in favor of Enums. |
| 05 | [05-rca-release-skew.md](05-rca-release-skew.md) | The deep-dive post-mortem into the `v6.87.1` tag mismatch CI failure. |
