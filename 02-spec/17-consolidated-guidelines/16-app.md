# Consolidated: App Specifications

**Version:** 3.3.0
**Updated:** 2026-04-16

---

## Purpose

Root-level folder for **application-specific specification content** — feature definitions, workflows, architecture decisions, and implementation specs that are specific to a particular application rather than cross-cutting coding principles.

An AI reading only this file must be able to create, organize, and review app-specific spec files correctly.

---

## What Belongs Here

| Content Type | Examples | File Pattern |
|-------------|----------|-------------|
| Feature specs | User authentication flow, dashboard layout | `01-feature-name.md` |
| App architecture | State management strategy, API design | `02-architecture-topic.md` |
| Workflow definitions | Onboarding flow, data import pipeline | `03-workflow-name.md` |
| UI/UX specifications | Screen-by-screen interaction specs | `04-screen-name.md` |
| Business rules | Pricing logic, access control matrix | `05-business-domain.md` |
| Integration specs | Third-party API integrations, webhooks | `06-integration-name.md` |
| Data flow specs | ETL pipelines, event processing | `07-data-flow-name.md` |

---

## App Spec File Template

````markdown

# App Spec: [Feature/Topic Title]

**Version:** 1.0.0
**Status:** Draft | Review | Approved | Implemented
**Created:** YYYY-MM-DD

---

## Overview

What this feature/component does and why it exists.

## Requirements

### Functional Requirements

- FR-1: [requirement]
- FR-2: [requirement]

### Non-Functional Requirements

- NFR-1: [performance, security, accessibility]

## Design

### User Flow

1. Step 1
2. Step 2
3. Step 3

### Data Model

Tables, fields, relationships affected. Reference `23-app-db/` for schema.

### API Endpoints

Routes needed (reference REST conventions from `15-rest-api-conventions.md`).

## Edge Cases

| Scenario | Expected Behavior |
|----------|------------------|
| [edge case] | [what should happen] |

## Dependencies

- Which spec modules this feature relies on
- External services or APIs needed

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
````

---

## Spec Lifecycle

| Status | Meaning | Who Can Transition |
|--------|---------|-------------------|
| **Draft** | Initial write-up, incomplete or unreviewed | Author |
| **Review** | Ready for team review | Author → Reviewer |
| **Approved** | Reviewed and accepted for implementation | Reviewer |
| **Implemented** | Code is written and merged | Developer |

### Transition Rules

1. Draft → Review: All required sections must be filled
2. Review → Approved: At least one reviewer sign-off
3. Approved → Implemented: All acceptance criteria checked
4. Any status can revert to Draft if requirements change

---

## Placement Decision Guide

### Place content in `21-app/` when:

- It defines a specific app feature or workflow
- It documents app architecture decisions not reusable across projects
- It specifies UI/UX flows tied to a particular application
- It contains business logic rules unique to the app

### Do NOT place here:

| Content Type | Correct Location |
|-------------|-----------------|
| Cross-cutting coding principles | Core fundamentals (`01–20`) |
| Bug analysis or fixes | `22-app-issues/` |
| Research or evaluations | `10-research/` or `02-coding-guidelines/10-research/` |
| Database schema (core conventions) | `04-database-conventions/` |
| Database schema (app-specific) | `23-app-db/` |
| Design system tokens (core) | `07-design-system/` |
| Design system (app-specific) | `24-app-ui-design-system/` |

---

## Required Files

Every `21-app/` module must contain:

| File | Purpose | Required |
|------|---------|----------|
| `01-index.md` | Module overview, file inventory | **Always** |
| `97-acceptance-criteria.md` | Consolidated acceptance criteria | When 3+ features |
| `98-changelog.md` | Version history | When versioned |
| `99-consistency-report.md` | Cross-reference validation | When 3+ files |

---

## Naming Convention

```
02-spec/21-app/
├── 01-index.md
├── 01-user-authentication.md
├── 02-dashboard-layout.md
├── 03-data-import-pipeline.md
├── 04-notification-system.md
├── 97-acceptance-criteria.md
├── 98-changelog.md
└── 99-consistency-report.md
```

### Rules

- Lowercase kebab-case only
- Numeric prefix mandatory for sequencing
- Feature name should be 2–4 words, descriptive
- No abbreviations unless universally understood (API, UI, DB)

---

## Relationship to Other App Folders

The `21+` range forms a cohesive app-specific layer:

| Folder | Purpose | Relationship |
|--------|---------|-------------|
| `21-app/` | Feature specs and workflows | **This folder** — defines what to build |
| `22-app-issues/` | Bug tracking and resolution | References features from `21-app/` |
| `23-app-db/` | Schema and migration specs | Implements data models from `21-app/` |
| `24-app-ui-design-system/` | UI tokens and component specs | Implements visual specs from `21-app/` |

### Dependency Direction

```
21-app/ (defines features)
  ├── 23-app-db/ (implements data layer)
  ├── 24-app-ui-design-system/ (implements UI layer)
  └── 22-app-issues/ (tracks problems)
```

Features in `21-app/` drive work in the other three folders. Not the reverse.

---

## Anti-Patterns

| ❌ Don't | ✅ Do |
|----------|-------|
| Write feature specs with no acceptance criteria | Always include testable acceptance criteria |
| Mix core conventions with app-specific logic | Keep app logic in `21-app/`, conventions in `01–20` |
| Skip the data model section | Reference `23-app-db/` tables explicitly |
| Write implementation details (code snippets) | Write requirements and behavior, not code |
| Create specs for trivial features (< 1 day work) | Only spec features with multiple requirements or edge cases |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Issues | `02-spec/17-consolidated-guidelines/17-app-issues.md` |
| App Database | `02-spec/17-consolidated-guidelines/25-app-database.md` |
| App Design System | `02-spec/17-consolidated-guidelines/19-app-design-system-and-ui.md` |
| Coding Guidelines | `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md` |
| Spec Authoring Guide | `02-spec/17-consolidated-guidelines/04-spec-authoring.md` |

---

## Current Status

No app-specific specs have been added yet. Files will be created as features are specified.

---

*Consolidated app specs — v3.3.0 — 2026-04-16*
