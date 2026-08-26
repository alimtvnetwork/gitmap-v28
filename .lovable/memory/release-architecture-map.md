# Release Architecture Map (Single Source of Truth)

## 1. Canonical Version Source
- **File**: `version.json` at repository root.
- **Role**: Sole Single Source of Truth (SSOT) for the entire repository and its sub-components.
- **Format**:
```json
{
  "$schema_description": "Canonical Single Source of Truth for repository versioning. Every tool, script, and AI agent must read and update this file exclusively.",
  "$documentation": "docs/versioning.md",
  "$inheritance_rules": "Sub-components specify 'inherit' to use top-level version or define an explicit SemVer string.",
  "version": "6.111.0",
  "backend": {
    "version": "inherit",
    "status": "active"
  },
  "frontend": {
    "version": "inherit",
    "status": "active"
  }
}
```

## 2. Component Inheritance Protocol
- Sub-components declaring `"version": "inherit"` dynamically resolve to the top-level `"version"`.
- Sub-components never bump independently when set to `"inherit"`.
- Any external service or internal package imports root `version.json` directly.

## 3. Version Propagation Pin Sites
When a MINOR release is executed, the following sites update in lock-step:
1. `version.json`: Canonical version of record.
2. `readme.md`: Header badge and pinned version references.
3. `what-to-read.md`: Header badge and pinned version references.
4. `changelog.md`: Mandatory new version entry with Prompt Architect installer commands.
5. `docs/versioning.md` & `.lovable/versioning.md`: SSOT specification guides.

## 4. Test File Ban
- **Strict Rule**: No test files (`*_test.*`, `*.spec.*`, `test/*`) are ever scanned or modified during a release.
- Test files contain mock/synthetic data (e.g. dummy versions) and must remain isolated to prevent test suite corruption.
