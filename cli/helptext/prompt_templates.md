# prompts-template

Manage reusable Antigravity (AGY) prompt prefix and verification templates. Supports JSON storage, adding, editing, listing, removing, exporting, and importing single templates or entire template suites.

## Synopsis

```bash
gitmap prompts-template <subcommand> [args] [flags]
gitmap prompt-template <subcommand> [args] [flags]
gitmap pt <subcommand> [args] [flags]
```

## Aliases & Shorthands

- `gitmap prompts-template`
- `gitmap prompt-template`
- `gitmap prompts-templates`
- `gitmap prompt-templates`
- `gitmap prompt_templates`
- `gitmap pt`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ls, list` | List all available prompt templates with ID, name, and preview |
| `add <name> <content>` | Add a new prompt template |
| `edit <name\|id> <content>` | Update content of an existing prompt template |
| `rm, delete <name\|id>` | Remove a prompt template by name or ID |
| `export <name\|id> <file.json>` | Export a specific template to a JSON file |
| `import <file.json>` | Import a template from a JSON file |
| `export-all [file.json]` | Export all registered templates to a single JSON suite file |
| `import-all <file.json>` | Bulk import templates from a JSON suite file |
| `help` | Show this prompt templates command documentation |

---

## Default Built-in Template: `is-done`

GitMap ships with a pre-installed default template named `is-done`:

```text
Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.
```

This template is automatically registered if no templates exist, and is used by default when running `gitmap agy rerun last [N]`.

---

## JSON Storage Schema

Templates are stored in valid JSON format:

### Single Template Format

```json
{
  "id": "is-done",
  "name": "is-done",
  "description": "Standard completion and quality verification check",
  "content": "Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.",
  "createdAt": "2026-09-18T10:00:00Z",
  "updatedAt": "2026-09-18T10:00:00Z"
}
```

### Bulk Suite Format (`export-all` / `import-all`)

```json
{
  "version": "1.0.0",
  "templates": [
    {
      "id": "is-done",
      "name": "is-done",
      "description": "Standard completion and quality verification check",
      "content": "Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes."
    },
    {
      "id": "code-review",
      "name": "code-review",
      "description": "Code review against specifications and coding guidelines",
      "content": "Review the following code against canonical coding guidelines, error wrappers, and line limits:"
    }
  ]
}
```

---

## Integration with Antigravity Rerun (`agy rerun`)

Prompt templates are integrated into `gitmap agy rerun last [N]` using the `-p, --prompt` flag:

```bash
# Uses the default "is-done" template automatically
gitmap agy rerun last 1

# Match template by exact ID or name
gitmap agy rerun last 1 -p is-done

# Match template by prefix
gitmap agy rerun last 1 -p is

# Use custom template
gitmap agy rerun last 3 -p code-review
```

When resolved, the template's content is prepended to the user prompt before dispatching it to Antigravity or copying to the system clipboard.

---

## Remote Delegation & Triad Parity

Prompt templates can be managed and inspected across remote nodes:

```bash
# List templates on remote SSH host
gitmap ssh exec pt ls

# List templates across all cluster nodes
gitmap cluster exec pt ls

# List templates via Servers-Clients (SC) delegation
gitmap sc exec pt ls
```

---

## Examples

```bash
# List all registered prompt templates
gitmap prompts-template ls
gitmap pt ls

# Add a custom verification template
gitmap pt add code-review "Review following code against coding guidelines and error management rules:"

# Edit an existing template
gitmap pt edit code-review "Perform thorough line-by-line review against canonical rules:"

# Export a single template to JSON
gitmap pt export is-done ./templates/is-done.json

# Import a single template from JSON
gitmap pt import ./templates/rca-fix.json

# Export all templates to a single JSON suite file
gitmap pt export-all ./templates/all-templates.json

# Import all templates from a JSON suite file
gitmap pt import-all ./templates/all-templates.json

# Remove a custom template by name or ID
gitmap pt rm code-review
```

---

## See Also

- [agy](agy.md) — Antigravity workspaces, prompts rerun, and VS Code inspection
- [pipeline](pipeline.md) — CI/CD pipeline status, logs, and database management
- [storage](storage.md) — Inspect disk space, pipeline logs, and SQLite database storage
