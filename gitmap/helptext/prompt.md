# prompt — AI Prompt Template Management

Manage, create, export, import, and inject structured AI prompt templates for Antigravity, coding guidelines, and autonomous workflows.

## Subcommands

- `gitmap prompt ls`: List all installed prompt templates with version and descriptions.
- `gitmap prompt show <slug>`: Display full frontmatter metadata and body of a prompt template.
- `gitmap prompt add <slug> <file.md>`: Install or update a prompt template from a markdown file.
- `gitmap prompt rm <slug>`: Delete an installed prompt template.
- `gitmap prompt export [file.zip]`: Export all prompt templates into a portable `.zip` bundle.
- `gitmap prompt import <file.zip|file.md>`: Import prompt templates from a zip archive or single markdown file.
- `gitmap prompt inject <slug> <target>`: Inject a prompt template directly into an AGY project or group.

## Examples

```bash

# List all prompt templates

gitmap prompt ls

# View prompt template

gitmap prompt show code-review

# Add a new custom prompt template

gitmap prompt add security-audit ./prompts/security-audit.md

# Export all templates to a zip archive

gitmap prompt export ./backup/prompts.zip

# Import templates from a zip bundle

gitmap prompt import ./backup/prompts.zip

# Inject a prompt into an AGY project or group

gitmap prompt inject code-review my-project
```
