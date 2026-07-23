# Documentation Standards

## Overview

Standards for inline comments, README structure, and API documentation
to ensure every project is navigable by new contributors within minutes.

---

## Inline Comments

### When to Comment

| Situation                        | Action                         |
|----------------------------------|--------------------------------|
| Non-obvious business logic       | Comment the *why*              |
| Complex algorithm or formula     | Comment the intent             |
| Workaround or hack               | Comment with `// HACK:` prefix |
| Self-explanatory code            | Do not comment                 |
| Restating what the code does     | Remove the comment             |

### Comment Style

- Use `//` for single-line comments (Go, TypeScript).
- Place comments on the line above the code, not inline.
- Start with a capital letter. End with a period only for full sentences.
- Keep comments under 80 characters per line.

```go
// Retry three times to handle transient network failures.
for i := 0; i < maxRetries; i++ {
    err = fetchRemote(url)
}
```

### Marker Comments

Use standardized prefixes for actionable items:

| Prefix    | Meaning                              |
|-----------|--------------------------------------|
| `TODO:`   | Planned work, not blocking           |
| `FIXME:`  | Known bug, needs attention           |
| `HACK:`   | Intentional workaround               |
| `NOTE:`   | Context for future readers           |

Always include a brief description after the prefix.

---

## README Structure

### Required Sections

Every project README includes these sections in order:

1. **Title** — Project name and one-line description.
2. **Overview** — What the project does and who it serves (2–3 sentences).
3. **Quick Start** — Minimal steps to install, configure, and run.
4. **Usage** — Common commands, flags, or API calls with examples.
5. **Configuration** — Environment variables, config files, defaults.
6. **Development** — How to build, test, and contribute.
7. **License** — License type and link.

### Optional Sections

Add only when relevant:

- **Architecture** — High-level diagram or description.
- **Troubleshooting** — Common issues and solutions.
- **Changelog** — Link to changelog or release notes.
- **See Also** — Links to related projects or documentation.

### Formatting Rules

- Use a single `#` heading for the project title.
- Use `##` for top-level sections.
- Code blocks specify the language (` ```bash `, ` ```go `).
- Keep the Quick Start under 10 lines of instruction.
- Pin dependency versions in install commands.

---

## API Documentation

### Function and Method Docs

Every exported function includes a doc comment covering:

1. **What** — One-line summary of behavior.
2. **Parameters** — Describe non-obvious inputs.
3. **Return** — What is returned and when errors occur.

```go
// FindBySlug returns the first repo matching the given slug.
// Returns ErrDBNoMatch if no repo exists with that slug.
func FindBySlug(db *sql.DB, slug string) (Repo, error) {
```

### Package-Level Docs

Each package has a doc comment on its primary file explaining the
package's purpose and responsibilities in one to three sentences.

```go
// Package store provides SQLite persistence for repos, groups,
// releases, and command history.
package store
```

### CLI Help Text

Every command and subcommand includes:

- A one-line description (shown in parent command listing).
- A usage line with argument placeholders.
- Flag descriptions with defaults.
- At least one example.

```
Usage:
  gitmap-v28 release <version> [flags]

Flags:
  --draft       Mark as draft release (default: false)
  --dry-run     Preview without executing (default: false)

Examples:
  gitmap-v28 release 2.14.0
  gitmap-v28 release 2.14.0 --draft
```

---

## Specification Documents

### File Naming

Specifications use numeric prefixes for ordering:

```
01-overview.md
02-cli-interface.md
03-scanner.md
```

### Required Structure

Each specification includes:

1. **Title** — `# Feature Name` as the only H1.
2. **Overview** — What the feature does (1–3 sentences).
3. **Behavior** — Detailed rules, flags, and edge cases.
4. **Implementation** — Package structure and file layout.
5. **Constraints** — Line limits, style rules, error handling.

### Tables Over Prose

Use markdown tables for flags, columns, error codes, and mappings.
Tables are faster to scan than paragraphs.

---

## Constraints

- Comments explain *why*, not *what*.
- No commented-out code in production.
- READMEs stay under 300 lines.
- Every exported symbol has a doc comment.
- Spec files follow the numbered naming convention.
- Help text includes at least one example per command.
