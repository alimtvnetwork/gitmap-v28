# Consolidated: Root-Level Research

**Version:** 3.3.0
**Updated:** 2026-04-16

---

## Purpose

Dedicated root-level folder for all exploratory and evaluative work that supports the spec system but is **not tied to a specific coding guideline**. This covers cross-cutting investigations, architectural explorations, infrastructure evaluations, and tooling assessments.

An AI reading only this file must be able to create, organize, and review root-level research documents correctly.

---

## What Belongs Here

| Content Type | Examples | Deliverable |
|-------------|----------|-------------|
| Architecture research | Evaluating deployment strategies, DB architectures | Architecture Decision Record (ADR) |
| Cross-cutting evaluations | CI/CD tool comparison, monitoring solutions | Decision matrix |
| Tooling investigations | IDE plugins, linter configurations, AI assistants | Setup guide + pros/cons |
| Infrastructure research | Hosting providers, CDN options, security tools | Cost/feature comparison |
| Process improvements | Workflow optimizations, automation opportunities | Process proposal |
| DevOps evaluations | Container orchestration, logging platforms | Feature matrix + cost analysis |
| Security assessments | Auth providers, vulnerability scanners, WAF options | Risk matrix + recommendation |

---

## Placement Rules

| Scope | Location |
|-------|----------|
| Root-level research (not tied to coding guidelines) | `02-spec/10-research/` |
| Coding-guidelines-specific research | `02-spec/02-coding-guidelines/10-research/` |

### Decision Guide

| Question | If YES → | If NO → |
|----------|----------|---------|
| Would the findings change how you write code (language, framework, pattern)? | Coding guidelines research | Check next |
| Does it affect infrastructure, deployment, tooling, or process? | Root-level research (here) | Check next |
| Is it about project architecture that crosses multiple domains? | Root-level research (here) | Ask human |

---

## Research File Template

Uses the same template as `10-research.md` (coding guidelines research). See that file for the complete template structure with Question, Context, Methodology, Findings, Recommendation, and Decision sections.

### ADR Template (Architecture Decisions)

For architecture-specific research, use this extended template:

````markdown

# ADR: [Decision Title]

**Status:** Proposed | Accepted | Deprecated | Superseded
**Created:** YYYY-MM-DD
**Author:** [name]
**Supersedes:** [ADR number, if applicable]

---

## Context

What is the issue that we're seeing that motivates this decision?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

### Positive

- [benefit 1]
- [benefit 2]

### Negative

- [trade-off 1]
- [trade-off 2]

### Risks

- [risk 1 with mitigation strategy]

## Alternatives Considered

| Option | Pros | Cons | Why Not |
|--------|------|------|---------|
| [alt 1] | ... | ... | ... |
| [alt 2] | ... | ... | ... |
````

---

## Research Lifecycle

| Status | Meaning | Next Action |
|--------|---------|-------------|
| **Proposed** | Topic identified, not yet started | Assign author, define methodology |
| **In Progress** | Actively being researched | Complete findings section |
| **Complete** | Findings documented, recommendation made | Team review → Decision |
| **Archived** | Decision made, no longer active | Reference only |

### ADR Lifecycle

| Status | Meaning |
|--------|---------|
| **Proposed** | Under discussion |
| **Accepted** | Decision adopted and in effect |
| **Deprecated** | No longer applies but kept for history |
| **Superseded** | Replaced by a newer ADR (link to it) |

---

## Naming Convention

```
02-spec/10-research/
├── 01-index.md
├── 01-deployment-strategy-evaluation.md
├── 02-monitoring-tool-comparison.md
├── 03-cicd-platform-assessment.md
├── 04-adr-database-hosting.md
└── 99-consistency-report.md
```

- Prefix ADR files with `NN-adr-` for quick identification
- Standard research uses `NN-descriptive-name.md`

---

## Scoring Criteria

Use the same standardized 1–5 rubric defined in `10-research.md`. For infrastructure/tooling research, add these criteria:

| Criterion | What to Evaluate |
|-----------|-----------------|
| **Cost** | Monthly/annual cost at projected scale |
| **Vendor lock-in** | Difficulty of migration away |
| **Compliance** | GDPR, SOC2, data residency support |
| **Integration** | How well it fits existing stack |
| **Scalability** | Behavior under 10x, 100x load |

---

## Anti-Patterns

| ❌ Don't | ✅ Do |
|----------|-------|
| Put language/framework comparisons here | Route to `02-coding-guidelines/10-research/` |
| Make architecture decisions without an ADR | Use the ADR template above |
| Skip cost analysis for infrastructure research | Always include cost estimates |
| Ignore vendor lock-in assessment | Evaluate migration difficulty |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding guidelines research | `02-spec/17-consolidated-guidelines/10-research.md` |
| Source folder | `02-spec/10-research/01-index.md` |
| CI/CD pipeline specs | `02-spec/17-consolidated-guidelines/18-cicd-pipeline-workflows.md` |
| Split DB architecture | `02-spec/17-consolidated-guidelines/08-split-db-architecture.md` |
| Spec authoring | `02-spec/17-consolidated-guidelines/04-spec-authoring.md` |

---

## Current Status

No research documents have been added yet. Files will be created as research needs arise.

---

*Consolidated root-level research — v3.3.0 — 2026-04-16*
