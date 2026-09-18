# gitmap prompts-template

Manage reusable Antigravity (AGY) prompt prefix and verification templates. Supports JSON storage, adding, editing, listing, removing, exporting, and importing single templates or entire template suites.

## Synopsis

```bash
gitmap prompts-template <subcommand> [args] [flags]
```

## Aliases & Shorthands

- `gitmap prompts-template`
- `gitmap prompt-template`
- `gitmap prompts-templates`
- `gitmap prompt-templates`
- `gitmap pt`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ls`, `list` | List all available prompt templates |
| `add <name> <content>` | Add a new prompt template |
| `edit <name\|id> <content>` | Update content of an existing prompt template |
| `rm`, `delete <name\|id>` | Remove a prompt template by name or ID |
| `export <name\|id> <file.json>` | Export a specific template to a JSON file |
| `import <file.json>` | Import a template from a JSON file |
| `export-all [file.json]` | Export all registered templates to JSON |
| `import-all <file.json>` | Bulk import templates from a JSON suite |

---

## Default Built-in Template: `is-done`

GitMap ships with a pre-installed default template named `is-done`:

```text
Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.
```

This template is automatically registered if no templates exist, and is used by default when running `gitmap agy rerun last [N]`.

---

## JSON Format

Templates are stored in valid JSON format:

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

Bulk export (`export-all`) format:

```json
{
  "version": "1.0.0",
  "templates": [
    {
      "id": "is-done",
      "name": "is-done",
      "description": "Standard completion and quality verification check",
      "content": "Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes."
    }
  ]
}
```

---

## Examples

```bash
# List all registered prompt templates
gitmap prompts-template ls

# Add a custom code review template
gitmap prompts-template add code-review "Review following code against coding guidelines and error management rules:"

# Edit an existing template
gitmap prompts-template edit code-review "Perform thorough line-by-line review against canonical rules:"

# Export a single template to JSON
gitmap prompts-template export is-done ./templates/is-done.json

# Import a single template from JSON
gitmap prompts-template import ./templates/rca-fix.json

# Export all templates to a file
gitmap prompts-template export-all ./templates/all-templates.json

# Import all templates from a file
gitmap prompts-template import-all ./templates/all-templates.json

# Remove a custom template
gitmap prompts-template rm code-review
```
