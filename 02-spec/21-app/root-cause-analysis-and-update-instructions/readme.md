# Specification Index: Antigravity IDE RCA & Update Instructions

- **Module:** `02-spec/21-app/root-cause-analysis-and-update-instructions/`
- **Parent Task:** `root-cause-analysis-and-update-instructions`
- **Status:** `APPROVED`

---

## Documents in this Directory

| Specification Document | Purpose & Scope |
| :--- | :--- |
| **`01-root-cause-analysis.md`** | Comprehensive Root Cause Analysis of the version gap between Windows (`2.22.0`) and Linux (`2.19.1`), in-app update failure in Electron on bare tarballs, absence of Google APT repo, and tooling version staleness. |
| **`02-dynamic-manifest-and-update-architecture.md`** | GitMap dynamic manifest resolution engine, Cloud Run updater integration, fallback defaults (`2.22.0` / `5446056071266304`), and launcher wrapper preservation. |
| **`03-scripts-fixer-and-agm-parity.md`** | Modifications for companion tools (`scripts-fixer`, `antigravity-manager`) to achieve version 2.22.0 parity across the fleet. |

---

## Key Invariants

- **Zero Absolute Paths:** All references strictly utilize relative Git paths.
- **Data Protection:** Absolute prohibition against deleting or altering `~/.gemini/` or `~/.config/Antigravity/`.
- **Runtime Sandbox:** Always preserve `antigravity-ide.run` enforcing `--no-sandbox` for unprivileged user namespaces.
