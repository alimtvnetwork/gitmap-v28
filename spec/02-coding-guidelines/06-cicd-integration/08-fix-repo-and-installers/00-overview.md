# Fix-Repo & Installer Scripts — CI/CD Spec

> **Version:** 1.0.0
> **Status:** Active
> **Updated:** 2026-04-28
> **AI Confidence:** Production-Ready
> **Ambiguity:** None — every flag, exit code, and side effect is normative below.

---

## Keywords

`installer`, `fix-repo`, `visibility-change`, `release-install`, `pinned-mode`,
`implicit-mode`, `ci`, `cd`, `bash`, `powershell`, `cross-platform`,
`coding-guidelines`, `code-red`

---

## Purpose

This subfolder is the **canonical, AI-shareable contract** for every
top-level shell/PowerShell script the project ships:

| Script (Bash) | Script (PowerShell) | Role |
|---|---|---|
| `install.sh` | `install.ps1` | Implicit/pinned installer (`raw.githubusercontent.com` entry point) |
| `release-install.sh` | `release-install.ps1` | Pinned-only installer (GitHub Release-asset entry point) |
| `fix-repo.sh` | `fix-repo.ps1` | Versioned-repo-name token rewriter |
| `visibility-change.sh` | `visibility-change.ps1` | Toggle / set GitHub or GitLab repo visibility |

Each script has a Bash and PowerShell sibling that MUST be behaviorally
identical (only sigil + exit-code reporting differ). They are bundled by
the linter pack release pipeline (`linters-cicd/`) and consumed via the
one-liner installer in `linters-cicd/install.sh`.

---

## Why this lives in `06-cicd-integration/`

These scripts are how a downstream pipeline **bootstraps** the linter
pack and **keeps a fork's branding consistent** between repo renames
(`coding-guidelines-v24` → `coding-guidelines-v24`). They are CI/CD
infrastructure, not application code. Sibling docs in this folder
(SARIF contract, plugin model, distribution) all assume the install
contract defined here.

---

## Document inventory

| File | Purpose |
|------|---------|
| `00-overview.md` | This file — scope, authority, cross-refs |
| `01-fix-repo-contract.md` | Normative spec for `fix-repo.sh` / `fix-repo.ps1` |
| `02-installer-contract.md` | Normative spec for `install.sh` / `install.ps1` and `release-install.*` |
| `03-visibility-change-contract.md` | Normative spec for `visibility-change.sh` / `.ps1` |
| `97-acceptance-criteria.md` | Binary AC list + test-matrix table referencing `tests/installer/*` |
| `98-faq.md` | Common questions: pinning, dry-run, offline, log dir |

---

## Authority & precedence

If this folder ever contradicts an older `spec-authoring/22-fix-repo/`
or `spec-authoring/23-visibility-change/` doc, **this folder wins**.
The `spec-authoring/` versions remain as historical design notes; this
folder is the contract every implementation and AI integration must
honor.

---

## Cross-references

- [CI/CD Integration overview](../00-overview.md)
- [SARIF contract](../01-sarif-contract.md)
- [Distribution](../05-distribution.md) — how installers are packaged
- [Generic installer behavior](../../../14-update/27-generic-installer-behavior.md)
- [Release pinned installer](../../../14-update/25-release-pinned-installer.md)
- [Code Red Guidelines](../../../17-consolidated-guidelines/00-overview.md)
- [`tests/installer/`](../../../../tests/installer/) — executable AC

---

## Contributors

- **Md. Alim Ul Karim** — Creator & Lead Architect
- **Riseup Asia LLC** — Sponsor
