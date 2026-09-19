# 23 — SSH Fleet Management Specification

**Version:** 1.0.0  
**Updated:** 2026-09-20  
**Status:** Canonical  
**AI Confidence:** Production-Ready  
**Ambiguity:** None  

---

## Purpose & Scope

The SSH Fleet Management module specifies the architectural protocols, public key authorization strategies, remote execution runtimes, node enrollment, and registry lifecycle management for SSH nodes and clusters across Windows, Linux, and macOS platforms in GitMap.

---

## Mandatory Module Scoring

| Criterion | Status |
|-----------|--------|
| `01-index.md` present in module | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Keywords

`ssh` · `fleet-management` · `authorized-keys` · `windows-ssh` · `administrators-authorized-keys` · `dual-table-purge`

---

## Document Inventory

| File | Type | Description |
|------|------|-------------|
| [`01-index.md`](./01-index.md) | Entry Point | Module architecture, inventory, scoring, and cross-references |
| [`01-windows-ssh-authorized-keys-spec.md`](./01-windows-ssh-authorized-keys-spec.md) | Specification | Windows OpenSSH authorized keys, ACLs, deduplication, and service lifecycle |

---

## Cross-References

- [SSH Keys Specification](../21-app/50-ssh-keys.md)
- [Golang Coding Guidelines](../02-coding-guidelines/03-golang/00-overview.md)
- [Error Management Architecture](../03-error-manage/02-error-architecture/00-overview.md)
