---
name: write-antigravity
description: Author, update, and persist Antigravity skills, rules, and configuration architecture in alignment with repo standards.
---

# Antigravity Customization Architecture & Rule Authoring

Maintains and authors Antigravity agent customizations, including skills (`.agents/skills/<slug>/skill.md`), rules (`.agents/rules/<slug>.md`), and tooling integrations.

## Core Directives

1. **Skill Layout:** Every skill resides in `.agents/skills/<slug>/skill.md` with YAML frontmatter containing `name` and `description`.
2. **Rule Layout:** Coding and architectural rules reside in `.agents/rules/<slug>.md`.
3. **Strict Lowercase:** Filenames must strictly use lowercase naming (e.g. `skill.md`, `agents.md`).
4. **Strict Relative Git Paths:** Total ban on absolute paths or `file:///` URIs.
5. **No Source Code Refactoring:** When authoring agent definitions or skills, do not refactor application source code unless explicitly instructed.

## Verification Checklist

- [ ] Skill contains valid YAML frontmatter (`name` and `description`).
- [ ] Relative Git paths only (no drive letters or `file:///` URIs).
- [ ] Filenames strictly lowercase.
- [ ] Mirrored scripts in `.agents/scripts/` match `03-ai-scripts/`.
