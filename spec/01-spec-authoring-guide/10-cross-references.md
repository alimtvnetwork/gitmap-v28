# Cross-References

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Cross-references are the connective tissue of the specification tree. They enable navigation between related modules and are automatically validated by the health dashboard link scanner. This document defines how to write, validate, and maintain cross-references.

---

## Golden Rule

> **All cross-references MUST use file-relative paths.** Never use root-relative paths (starting with `/`).

---

## Path Syntax

> **Note:** All examples below are illustrative templates showing path patterns. The target files may not exist in this project.

### Same folder

```markdown
See: [Architecture](./02-architecture.md)
```

### Sibling file in same directory

```markdown
See: [Error Codes](./07-error-codes.md)
```

### Parent folder

```markdown
See: [Module Overview](../01-index.md)
```

### Another module (one level up, then into sibling)

```markdown
See: [Coding Guidelines](../02-coding-guidelines/01-index.md)
```

### Deep cross-reference (from a subfolder to another module's subfolder)

```markdown
See: [GSearch Backend](../../09-gsearch-cli/01-backend/01-index.md)
```

### From spec to memories

```markdown
See: [Database Standards](../../.lovable/memories/architecture/database-standards.md)
```

---

## Cross-Reference Table Pattern

Every `01-index.md` file should end with a cross-reference table:

```markdown

## Cross-References

| Reference | Location |
|-----------|----------|
| Module A | `../XX-module-a/01-index.md` |
| Module B | `../YY-module-b/01-index.md` |
| Memory X | `../../.lovable/memories/category/file.md` |
```

---

## Common Mistakes

### ❌ Root-relative paths

```markdown
See: [Guidelines](/02-spec/02-coding-guidelines/01-index.md)
```

**Why wrong:** Root-relative paths break depending on the viewer's base URL.

### ❌ Absolute filesystem paths

```markdown
See: [Guidelines](/dev-server/02-spec/02-coding-guidelines/01-index.md)
```

**Why wrong:** Filesystem paths are environment-specific.

### ❌ Mixed case in paths

```markdown
See: [Guidelines](../03-Coding-Guidelines/01-index.md)
```

**Why wrong:** Paths must be entirely lowercase to match the kebab-case convention.

### ❌ Missing file extension

```markdown
See: [Guidelines](../02-coding-guidelines/00-overview)
```

**Why wrong:** Always include `.md` extension for explicit targeting.

---

## Depth Calculation

To write the correct relative path, count directory levels:

| From | To | Path |
|------|----|------|
| `02-spec/09-gsearch-cli/01-backend/01-arch.md` | `02-spec/02-coding-guidelines/01-index.md` | `../../02-coding-guidelines/01-index.md` |
| `02-spec/09-gsearch-cli/01-index.md` | `02-spec/02-coding-guidelines/01-index.md` | `../02-coding-guidelines/01-index.md` |
| `02-spec/09-gsearch-cli/01-backend/01-arch.md` | `02-spec/09-gsearch-cli/01-index.md` | `../01-index.md` |

**Formula:** Count how many directories UP you need to go, then navigate DOWN to the target.

---

## Link Validation

The health dashboard scanner (`linter-scripts/generate-dashboard-data.cjs`) validates all markdown links:

1. **Scans** every `.md` file in `spec/` and `.lovable/memories/`
2. **Extracts** all markdown link patterns (`[text]\(path)` syntax)
3. **Resolves** relative paths from the source file's location
4. **Checks** that the target file exists on disk
5. **Ignores** links inside code-fenced blocks (`` ``` ``)
6. **Reports** broken links with source file, line number, and target path

### Running the Scanner

```bash
node linter-scripts/generate-dashboard-data.cjs
```

The output JSON includes a `BrokenLinks` array. Zero broken links = healthy.

---

## Impact of Module Renaming

When a module's numeric prefix changes (e.g., `09-gsearch-cli` → `09-gsearch-cli`), ALL cross-references pointing to it must be updated. This includes:

1. **Other spec files** — Any `../09-gsearch-cli/...` → `../09-gsearch-cli/...`
2. **Memory files** — Any reference to the old path
3. **Dashboard code** — `SpecFileViewer.tsx` module definitions
4. **Master index** — `02-02-spec/01-index.md` module table
5. **Consistency reports** — If they mention module paths

**This is why module renaming should be done carefully and in batches.**
