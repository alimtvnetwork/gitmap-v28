# Prompt Architect: Canonical AI Prompts Library

This directory hosts the canonical, production-grade prompts library for the Prompt Architect meta-repository and all connected repositories.

## Architecture: Version Segregation (V1 & V2)

To maintain 100% backward compatibility while delivering 10x-100x performance acceleration, prompts are segregated into two distinct release tiers:

| Tier | Directory | Primary Acceleration Engine | Fallback Engine | Description |
| :--- | :--- | :--- | :--- | :--- |
| **V1** | [`v1/`](v1/) | Python Toolchain (`03-ai-scripts/`) | Agent Built-ins | Original classic battle-tested prompts |
| **V2** | [`v2/`](v2/) | **GitMap AUM Engine** (`gitmap` CLI) | Python Scripts (`03-ai-scripts/`) | Modernized ultra-fast prompts with GitMap AUM as primary |

---

## V2 Modernization Highlights (GitMap AUM Integration)

1. **High-Speed File & Content Discovery:**
   - Universal Wildcard Search: `gitmap find "<pattern>" [-ext <ext>]` (<10ms across 10,000+ files)
   - Zero-Disk Streaming: `gitmap cat <filepath>`
   - Instant Filesystem Walk: `gitmap search "<symbol>"`
   - Substring Filename Search: `gitmap find-files-any "<str>"` (alias: `gitmap ffa`)
   - Directory Indexing: `gitmap list-files [pattern]` (alias: `gitmap lf`)

2. **Pipeline AI & Zero-Credit-Waste Dynamic Waiting:**
   - Telemetry Query: `gitmap pipeline-ai status --json`
   - Adaptive Sleep: `gitmap pipeline-ai status -t <etaSeconds>` (replaces expensive busy-polling loops)
   - Failure Log Snippets: `gitmap pipeline errors` (alias: `gitmap pe`)
   - Runner Target Diagnostics: `gitmap pipeline details` (alias: `gitmap pd`)

3. **Autonomous Automated Release:**
   - Single-Command Release: `gitmap release --bump <patch|minor> -y` (alias: `gitmap r -y`)
   - Release History: `gitmap changelog` (alias: `gitmap cl`), `gitmap list-versions` (alias: `gitmap lv`)

---

## Directory Index

Both `v1/` and `v2/` mirror the 22 canonical prompt categories:

```text
01-prompts/
├── v1/                                 # Classic prompts (Python-accelerated)
│   ├── 00-folder-structure/
│   ├── 01-prompt-library-setup/
│   ├── 02-core-workflow/
│   ├── 03-read-write/
│   ├── 04-coding-standards/
│   ├── 05-coding-guidelines/
│   ├── 06-testing-and-qa/
│   ├── 07-bug-fix/
│   ├── 08-dry-code/
│   ├── 09-commit-and-multi-agent-code-fix/
│   ├── 10-ui-and-design/
│   ├── 11-content-and-seo/
│   ├── 12-old-plan-prompts/
│   ├── 13-plan-audit/
│   ├── 14-execute/
│   ├── 15-cg-execute/
│   ├── 16-ci-cd/
│   ├── 17-release-management/
│   ├── 18-insults/
│   ├── 19-old-execute-prompts/
│   ├── 20-ai-fix-script-prompts/
│   └── 21-temp-end-to-end-tests/
└── v2/                                 # Modernized prompts (GitMap AUM Primary)
    ├── 00-folder-structure/
    ├── 01-prompt-library-setup/
    ├── ... (all 22 categories)
    └── 21-temp-end-to-end-tests/
```
