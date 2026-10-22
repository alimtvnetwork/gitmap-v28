# Specification 288: Synchronization & Cross-Repository Compatibility

- **Sequence:** 288
- **Slug:** `288-create-prompts-for-codebase-improvements`
- **File:** `02-spec/21-app/288-create-prompts-for-codebase-improvements/03-sync-and-compatibility-spec.md`
- **Status:** APPROVED & SPECIFIED

---

## 1. Synchronization Architecture

GitMap maintains a multi-repository ecosystem where canonical prompts, Antigravity skills, Cursor skills, and AI automation scripts are synchronized across 43 downstream repositories using:
📄 `03-ai-scripts/38-sync-prompts-skills-scripts.py`

### 1.1 Synchronization Boundaries & Target Mappings

| Source Directory | Mirror Destination in Downstream Repositories | Sync Policy |
|---|---|---|
| `01-prompts/30-improvements/` | `01-prompts/30-improvements/` | **Synchronized** (Full directory mirroring) |
| `.agents/skills/find-100-ui-improvements/` | `.agents/skills/find-100-ui-improvements/` | **Synchronized** (Full skill mirroring) |
| `.agents/skills/find-100-codebase-improvements/` | `.agents/skills/find-100-codebase-improvements/` | **Synchronized** (Full skill mirroring) |
| `.agents/skills/find-100-system-ux-improvements/` | `.agents/skills/find-100-system-ux-improvements/` | **Synchronized** (Full skill mirroring) |
| `.cursor/skills/find-100-ui-improvements/` | `.cursor/skills/find-100-ui-improvements/` | **Synchronized** (Full skill mirroring) |
| `.cursor/skills/find-100-codebase-improvements/` | `.cursor/skills/find-100-codebase-improvements/` | **Synchronized** (Full skill mirroring) |
| `.cursor/skills/find-100-system-ux-improvements/` | `.cursor/skills/find-100-system-ux-improvements/` | **Synchronized** (Full skill mirroring) |
| `02-spec/21-app/288-*` | N/A | **Excluded** (Protected domain spec boundary) |

### 1.2 Invariant Verification Gates
1. **Zero Domain Specific Hardcoding:** Prompts must not assume host-specific paths (e.g. `/path/to/repo`). All examples must use workspace-relative paths or abstract identifiers.
2. **Polyglot Stack Applicability:** Prompts must be executable across Go, TypeScript/React, Python, and Rust repositories without modification.
3. **YAML Frontmatter Conformance:** Skills must provide standard frontmatter:
   ```yaml
   ---
   name: find-100-<domain>-improvements
   description: >-
     [Clear, concise description for AI agents and slash command dispatch]
   ---
   ```
4. **Clean UTF-8 & LF Line Endings:** All markdown and YAML files must be encoded in UTF-8 with standard LF (`\n`) line endings.
