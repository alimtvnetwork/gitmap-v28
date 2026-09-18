# Consolidated: Coding Guidelines Research

**Version:** 3.3.0
**Updated:** 2026-04-16

---

## Purpose

Dedicated folder for all exploratory and evaluative work that supports foundational coding guidelines. This is the **single canonical location** for coding-guidelines-scoped research — framework comparisons, language evaluations, pattern investigations, and tooling assessments that directly affect how code is written.

An AI reading only this file must be able to create, organize, and review coding-guidelines research documents correctly.

---

## What Belongs Here

| Content Type | Examples | Deliverable |
|-------------|----------|-------------|
| Comparative studies | Framework X vs Framework Y | Decision matrix with scores |
| Technology evaluations | Assessing a new library or tool | Pros/cons table with recommendation |
| Exploratory technical notes | Proof-of-concept findings | Working code + performance data |
| Game development research | Engine comparisons, architecture patterns | Feature comparison grid |
| Language evaluations | Assessing a new language for the stack | Hello-world benchmark + ecosystem review |
| Pattern research | Investigating design patterns for adoption | Pattern template with usage guide |
| ORM evaluations | Comparing GORM vs sqlx vs sqlc | Feature matrix + migration cost |
| Build tool benchmarks | Vite vs Webpack vs Turbopack | Build-time + bundle-size data |

---

## Placement Rule

All research content scoped to coding guidelines MUST be placed in `02-spec/02-coding-guidelines/10-research/`. Research not scoped to coding guidelines belongs in `02-spec/10-research/`.

### Decision Guide

| Question | If YES → | If NO → |
|----------|----------|---------|
| Would the findings change how you write code? | `02-coding-guidelines/10-research/` | `10-research/` |
| Is it about a language, framework, or library? | `02-coding-guidelines/10-research/` | Check scope |
| Is it about infrastructure, CI/CD, or deployment? | `10-research/` | Check scope |
| Is it about tooling that affects developer workflow? | Could be either — use judgment | Ask human |

---

## Research File Template

````markdown

# Research: [Topic Title]

**Status:** Proposed | In Progress | Complete | Archived
**Created:** YYYY-MM-DD
**Author:** [name]

---

## Question

What specific question is this research answering?

## Context

Why this research is needed now. Link to the spec or feature that triggered it.

## Methodology

How the evaluation was conducted:
- Criteria used for comparison
- Test environment details
- Data sources

## Findings

### Option A: [Name]

| Criterion | Score (1-5) | Notes |
|-----------|-------------|-------|
| Performance | | |
| Ecosystem | | |
| Learning curve | | |
| Maintenance | | |
| Community support | | |
| Type safety | | |

### Option B: [Name]

(Same table structure)

## Recommendation

Clear recommendation with rationale. Include:
- Which option and why
- Migration path if switching
- Risks and mitigations
- Estimated effort (hours/days)

## Decision

Final decision made (filled after team review). Link to any spec updates that result from this decision.
````

---

## Research Lifecycle

| Status | Meaning | Next Action |
|--------|---------|-------------|
| **Proposed** | Topic identified, not yet started | Assign author, define methodology |
| **In Progress** | Actively being researched | Complete findings section |
| **Complete** | Findings documented, recommendation made | Team review → Decision |
| **Archived** | Decision made, no longer active | Reference only |

### Transition Rules

1. A research file MUST NOT skip statuses (Proposed → Complete is invalid)
2. The Decision section is filled ONLY after team review
3. Archived research should link to the spec or code change it produced

---

## Scoring Criteria

When evaluating options, use this standardized rubric:

| Score | Meaning |
|-------|---------|
| 5 | Excellent — best-in-class, no significant downsides |
| 4 | Good — strong choice with minor trade-offs |
| 3 | Adequate — meets requirements but has notable gaps |
| 2 | Below average — significant concerns |
| 1 | Poor — not recommended for this use case |

### Mandatory Criteria (always include)

- **Performance** — runtime speed, memory usage, bundle size
- **Ecosystem** — community size, plugin availability, documentation quality
- **Learning curve** — time to productive use for the team
- **Maintenance** — release cadence, breaking change history, bus factor

### Optional Criteria (add as relevant)

- Type safety, accessibility support, SSR compatibility, mobile support, license terms

---

## Naming Convention

Research files use the standard numbered prefix pattern:

```
02-spec/02-coding-guidelines/10-research/
├── 01-index.md
├── 01-framework-comparison-react-vue.md
├── 02-orm-evaluation-prisma-drizzle.md
├── 03-state-management-patterns.md
├── 04-build-tool-benchmarks.md
└── 99-consistency-report.md
```

---

## Anti-Patterns

| ❌ Don't | ✅ Do |
|----------|-------|
| Write opinion-only research with no data | Include benchmarks, metrics, or structured scoring |
| Skip the methodology section | Document exactly how you evaluated |
| Leave Decision empty after review | Fill it immediately, even if "no change" |
| Put infrastructure research here | Route to `02-spec/10-research/` instead |
| Use inconsistent scoring scales | Always use the 1–5 rubric above |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Root-level research | `02-spec/17-consolidated-guidelines/15-root-research.md` |
| Source folder | `02-spec/02-coding-guidelines/10-research/01-index.md` |
| Coding guidelines | `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md` |
| Spec authoring | `02-spec/17-consolidated-guidelines/04-spec-authoring.md` |

---

## Current Status

No research documents have been added yet. Files will be created as research needs arise.

---

*Consolidated coding guidelines research — v3.3.0 — 2026-04-16*
