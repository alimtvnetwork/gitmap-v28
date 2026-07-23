---
name: CI Release Pipeline Untouchable
description: .github/workflows/release.yml and all release-pipeline scripts are STRICTLY off-limits. Never edit, refactor, "improve", or add steps. User explicitly forbade any modification after a perceived artifact-publishing regression.
type: constraint
---

# STRICTLY-PROHIBITED: Touching the CI Release Pipeline

**Scope (all forbidden):**
- `.github/workflows/release.yml`
- Any new workflow that triggers on `release/**` branches or `v*` tags
- `.github/scripts/smoke-installer.*` (release-mode smoke harness)
- Anything under `.gitmap/release/` and `.gitmap/release-assets/` (already in Core)

**Rule:** Do NOT edit, refactor, reorder, "clean up", add steps to, or remove steps from any of the above — even if a lint, style, or consistency rule would otherwise apply. Even a "harmless" whitespace change is forbidden.

**Why:** The release pipeline is the single source of truth for publishing GitHub Release artifacts (6 cross-compiled binaries + checksums + install scripts). Any change risks silently breaking artifact upload, which is invisible until a user hits an empty release page. User has explicitly stated this is broken-everything-class severity and revoked permission.

**If a release job is failing:**
- Diagnose by reading run logs only.
- Fix the upstream cause (failing test, broken build, compile error) in the application code so the gate passes.
- NEVER "fix" it by editing the workflow itself.

**If the user explicitly asks for a release-pipeline edit in the future:**
- Confirm explicitly in chat that they are overriding this constraint.
- Quote this rule back to them and require a second confirmation before editing.
