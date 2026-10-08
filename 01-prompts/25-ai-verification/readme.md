# AI Verification Prompts (`25-ai-verification`)

This directory contains autonomous retrospective verification and code quality audit prompts. These prompts inspect recent repository tasks, verify newly authored specifications, check completed plans, audit coding guideline compliance, and monitor CI/CD pipeline health via GitMap telemetry.

## Directory Index

| # | Prompt File | Target Mode | Description | Companion Skill |
| :---: | :--- | :--- | :--- | :--- |
| **01** | [`01-retrospective-ai-verification.md`](01-retrospective-ai-verification.md) | Canonical V6 Workflow | Retrospective audit of the recent 2–3 tasks, specs, code hygiene, and CI/CD status | `ai-verification` |

## Core Invariants

1. **Retrospective Scope:** Audits the last 2–3 completed tasks (~30–40 minutes) to verify consistency and quality.
2. **Acceptance Criteria Verification:** Confirms all touched specification files under `02-spec/21-app/` have structured `## Acceptance Criteria`.
3. **Coding Guideline Enforcement:** Enforces implicit booleans, zero absolute paths, and zero CI/CD disabling.
4. **GitMap Telemetry Integration:** Queries `gitmap pe -t` to ensure all active pipeline workflows remain green.
5. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well (TOTAL BAN on `file:///` URIs and absolute filesystem paths).
