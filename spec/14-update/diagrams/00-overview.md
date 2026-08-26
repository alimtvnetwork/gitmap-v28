# Diagrams — Self-Update & App Update

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Index of all Mermaid diagrams supporting the self-update and app update specifications. These diagrams visualize decision trees and workflows described in the parent module.

---

## Diagram Inventory

| # | File | Description | Format |
|---|------|-------------|--------|
| 01 | [01-self-update-workflow.mmd](./01-self-update-workflow.mmd) | Full 9-step `<binary> update` command decision tree with error handling and rollback paths | Mermaid |
| 02 | [02-update-cleanup-workflow.mmd](./02-update-cleanup-workflow.mmd) | 2-phase `<binary> update-cleanup` workflow covering temp copies and `.old` backup removal | Mermaid |

**Total:** 2 diagrams

---

## Rendering

These `.mmd` files use Mermaid flowchart syntax. To render:

- **VS Code**: Install the "Mermaid Preview" extension
- **GitHub**: Mermaid blocks render natively in `.md` files (wrap in ` ```mermaid ` fences)
- **CLI**: Use `mmdc` from `@mermaid-js/mermaid-cli`

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Update Command Workflow (references these diagrams) | [../16-update-command-workflow.md](../22-update-command-workflow.md) |
| Self-Update Overview | [../01-self-update-overview.md](../01-self-update-overview.md) |
| Cleanup Specification | [../06-cleanup.md](../06-cleanup.md) |

---

*Diagrams Overview — v3.2.0 — 2026-04-15*
