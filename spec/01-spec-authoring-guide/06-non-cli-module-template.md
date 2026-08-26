# Non-CLI Module Template

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Not every specification module is a CLI tool. Many modules are **flat** (no subfolders) or use **custom subfolder layouts** suited to their content. This document covers templates for research modules, utility modules, standards modules, and multi-category modules.

---

## Module Types

| Type | Subfolder Pattern | Example Modules |
|------|-------------------|-----------------|
| Flat | No subfolders | ai-research, license-manager, shared-preset-data |
| Standard CLI | `01-backend/`, `02-frontend/`, `03-deploy/` | gsearch-cli, brun-cli, ai-bridge-cli |
| Multi-category | Custom numbered subfolders | general-spec, coding-guidelines, spec-management-software |
| Combined | Aggregation of other modules | time-log-combined |

---

## Flat Module Structure

Used for small specifications (< 15 files) that don't need subfolder organization.

```
{NN}-{module-name}/
├── 00-overview.md
├── 01-{topic}.md
├── 02-{topic}.md
├── ...
├── 97-acceptance-criteria.md    # Optional
├── 98-changelog.md              # Optional
└── 99-consistency-report.md
```

### Example: AI Research (17-ai-research)

```
17-ai-research/
├── 00-overview.md
├── 01-llm-comparison.md
├── 02-embedding-strategies.md
├── 03-rag-patterns.md
├── ...
└── 99-consistency-report.md
```

---

## Multi-Category Module Structure

Used for large specifications that group content into logical categories.

```
{NN}-{module-name}/
├── 00-overview.md               # MUST list all categories with file counts
│
├── 01-{category}/
│   ├── 00-overview.md
│   ├── 01-{topic}.md
│   └── 99-consistency-report.md
│
├── 02-{category}/
│   ├── 00-overview.md
│   ├── 01-{topic}.md
│   └── 99-consistency-report.md
│
└── 99-consistency-report.md     # Root consistency report lists subfolders
```

### Example: Coding Guidelines (03-coding-guidelines)

```
03-coding-guidelines/
├── 00-overview.md
├── 01-cross-language/           # Language-agnostic rules
│   ├── 00-overview.md
│   ├── 01-issues-and-fixes-log.md
│   ├── 02-boolean-principles.md
│   ├── ...
│   ├── 15-master-coding-guidelines.md
│   └── 99-consistency-report.md
├── 02-typescript/               # TypeScript-specific
│   ├── 00-overview.md
│   └── ...
├── 03-golang/                   # Go-specific (includes enum sub-spec)
│   ├── 00-overview.md
│   ├── 01-enum-specification/   # Nested sub-module
│   │   └── 00-overview.md
│   └── ...
├── 04-php/                      # PHP-specific
│   └── ...
└── 05-rust/                     # Rust-specific
    └── ...
```

**Key observation:** The coding guidelines module demonstrates a **3-level deep** structure where `03-golang/01-enum-specification/` is a nested sub-module. This is the maximum recommended depth.

### Example: General Spec (01-general-spec)

```
01-general-spec/
├── 00-overview.md
├── 01-foundation/               # Architecture foundations
├── 02-systems/                  # System patterns
├── 03-quality/                  # Quality standards
├── 04-advanced/                 # Advanced patterns
├── 05-ux/                       # UX guidelines
├── 06-devops/                   # DevOps practices
├── 07-observability/            # Monitoring & logging
├── 08-data-governance/          # Data governance
├── 09-api-integration/          # API integration patterns
├── 10-wordpress/                # WordPress-specific standards
├── 99-meta/                     # Meta/governance files
├── 97-acceptance-criteria.md
├── 98-changelog.md
└── 99-consistency-report.md
```

**Key observation:** This module has **11 subfolders** (the most of any module) because it covers architecture-wide standards. Each subfolder has its own `00-overview.md`.

---

## Combined/Aggregation Module Structure

Used when multiple related modules share acceptance criteria or summary data.

```
{NN}-{combined-name}/
├── 00-overview.md
├── 00-acceptance-criteria-summary.md   # Aggregated criteria from child modules
└── 99-consistency-report.md
```

### Example: Time Log Combined (36-time-log-combined)

This module aggregates acceptance criteria from `34-time-log-cli` and `35-time-log-ui` into a single summary without duplicating the actual spec content.

---

## Massive Module: Spec Management Software

`02-spec-management-software` is the largest module (500+ files, 60+ subfolders). It uses an extended structure:

```
02-spec-management-software/
├── 00-overview.md
├── 01-ideas/                    # Feature ideas and brainstorming
├── 02-instructions/             # Implementation instructions
├── 03-project-overview/         # Project-level documentation
├── 04-coding-guidelines/        # Module-specific coding rules
├── 05-features/                 # 30+ feature sub-modules (01-30)
│   ├── 01-search/
│   ├── 02-auth/
│   ├── ...
│   └── 30-{feature}/
├── 06-error-management/         # Error handling patterns
├── 07-database-design/          # Database schemas
├── 08-roadmap-overview/         # Roadmap and milestones
├── 09-diagrams/                 # Architecture diagrams
├── 10-research/                 # Research notes
├── 11-skipped-features/         # Intentionally deferred features
├── 12-prompts/                  # AI prompts for code generation
├── 13-shared-packages/          # Shared Go packages
├── 14-microservices/            # Microservice definitions
├── 15-external-tools/           # External tool integrations
├── 19-data-models/              # Data model definitions
└── 99-consistency-report.md
```

**Key lesson:** Even at 500+ files, the same conventions apply — every folder has `00-overview.md`, files are numbered sequentially, and kebab-case is universal.

---

## Checklist: New Non-CLI Module

- [ ] Module number selected (next available)
- [ ] Root `00-overview.md` created with file/subfolder inventory
- [ ] Decided on flat vs. multi-category structure
- [ ] If multi-category: each subfolder has `00-overview.md`
- [ ] `99-consistency-report.md` at root
- [ ] `spec/00-overview.md` master index updated
- [ ] Cross-references added to related modules
