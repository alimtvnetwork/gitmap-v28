# Canonical Repository Versioning Specification (Single Source of Truth)

## 1. Overview & Principle
This repository and all managed repositories adhere to a strict **Single Source of Truth** versioning protocol:
- The canonical version of record for this repository is located exclusively in **`version.json`** at the repository root.
- **Never guess or infer versions** from test files, git log history, git tags, or documentation snippets.
- **Changing the version**: To update or bump the version of this repository, simply update the `"version"` field in `version.json`. All automated tools, build pipelines, and release workflows discover their active version directly from `version.json`.

---

## 2. Structure of `version.json`
`version.json` serves as the self-describing manifest for repository release state and tooling extensions:
```json
{
  "$schema_description": "Canonical Single Source of Truth for repository versioning. Every tool, script, and AI agent must read and update this file exclusively.",
  "$documentation": "docs/versioning.md",
  "version": "6.108.0",
  "coding-guidelines": {
    "version": "v24.0.0",
    "installed_at": "2026-08-26T14:30:00Z",
    "status": "active"
  }
}
```

---

## 3. Strict Directive for AI Assistants & Automated Tools
1. **Always Read `version.json` First**: Before performing any release, inspection, or status query, read `version.json`.
2. **Never Treat Test Placeholders as Real Versions**:
   - Synthetic versions (such as `v9.9.0`, `v9.9.1`, `v9.9.2`, `v99.99.99`) used inside test files (e.g. `autocommit_test.go`, test mocks, fixture files) are **strictly for unit test experiments and output assertion testing**.
   - Test versions are completely isolated inside temporary directories (`t.TempDir()`) and MUST NEVER be propagated or treated as the real project version.
3. **Atomic Synchronization**: When a release or version bump is executed, `version.json` is updated first, and all secondary display sites (`readme.md`, `changelog.md`) are synchronized in lock-step.

---

## 4. Cross-Repository Uniformity
Every repository within the ecosystem follows this exact standard. Tooling wrappers such as `gitmap cg prompts-version`, `gitmap cg status`, and `gitmap release` automatically parse and validate the root `version.json` manifest.
