// Package cmd — llmdocssections.go writes the remaining LLM.md sections.
package cmd

import (
	"strings"
)

// writeLLMGlobalFlags writes the global flags section.
func writeLLMGlobalFlags(sb *strings.Builder) {
	sb.WriteString("## Global Flags\n\n")
	sb.WriteString("These flags work with most commands:\n\n")
	sb.WriteString("| Flag | Description |\n")
	sb.WriteString("|------|-------------|\n")
	sb.WriteString("| `--help`, `-h` | Show help text for any command |\n")
	sb.WriteString("| `--verbose` | Enable debug logging |\n")
	sb.WriteString("| `--json` | JSON output (where supported) |\n")
	sb.WriteString("| `-A`, `--alias` | Use a repo alias instead of slug |\n")
	sb.WriteString("| `--group <name>` | Filter by group name |\n")
	sb.WriteString("| `--all` | Run against all tracked repos |\n")
	sb.WriteString("| `--stop-on-fail` | Halt batch operations after first failure |\n\n")
}

// writeLLMCodingConventions writes the coding conventions section.
func writeLLMCodingConventions(sb *strings.Builder) {
	sb.WriteString("## Coding Conventions\n\n")
	sb.WriteString("When modifying GitMap code, follow these rules:\n\n")
	sb.WriteString("1. **No magic strings** — All literals go in `constants/` package\n")
	sb.WriteString("2. **Functions ≤ 15 lines** — Extract helpers liberally\n")
	sb.WriteString("3. **Files ≤ 200 lines** — Split when approaching limit\n")
	sb.WriteString("4. **PascalCase** for exported constants\n")
	sb.WriteString("5. **`is`/`has` prefix** for boolean variables and functions\n")
	sb.WriteString("6. **Blank line before `return`** statements\n")
	sb.WriteString("7. **Chained `if` + `return`** for dispatch (not switch)\n")
	sb.WriteString("8. **No swallowed errors** — every error return must be checked and logged\n")
	sb.WriteString("9. **Group same-type parameters** — `func(a, b bool)` not `func(a bool, b bool)`\n")
	sb.WriteString("10. **Positive logic** in `if` conditions\n\n")
}

// writeLLMProjectStructure writes the project structure section.
func writeLLMProjectStructure(sb *strings.Builder) {
	sb.WriteString("## Project Structure\n\n")
	sb.WriteString("```\n")
	sb.WriteString("/\n")
	sb.WriteString("├── gitmap/                    # Main CLI (Go module)\n")
	sb.WriteString("│   ├── cmd/                   # Command handlers\n")
	sb.WriteString("│   ├── constants/             # All string constants\n")
	sb.WriteString("│   ├── model/                 # Data types\n")
	sb.WriteString("│   ├── store/                 # SQLite database\n")
	sb.WriteString("│   ├── release/               # Version/tag management\n")
	sb.WriteString("│   ├── cloner/                # Clone operations\n")
	sb.WriteString("│   ├── dashboard/             # HTML dashboard\n")
	sb.WriteString("│   ├── verbose/               # Debug logging\n")
	sb.WriteString("│   ├── completion/            # Shell completions\n")
	sb.WriteString("│   └── helptext/              # Embedded help (go:embed)\n")
	sb.WriteString("├── gitmap-updater/            # Standalone updater (Go module)\n")
	sb.WriteString("├── spec/                      # Specifications & design docs\n")
	sb.WriteString("├── docs-site/                 # Documentation website\n")
	sb.WriteString("├── CHANGELOG.md\n")
	sb.WriteString("├── README.md\n")
	sb.WriteString("└── LLM.md                     # This file\n")
	sb.WriteString("```\n\n")
}

// writeLLMDatabase writes the database schema section.
func writeLLMDatabase(sb *strings.Builder) {
	sb.WriteString("## Database Schema (Conceptual)\n\n")
	sb.WriteString("GitMap uses SQLite with tables for:\n\n")
	sb.WriteString("- **repos** — scanned repository records (path, slug, URLs, branch, type)\n")
	sb.WriteString("- **groups** / **group_members** — named groups of repos\n")
	sb.WriteString("- **aliases** — short names pointing to repo slugs\n")
	sb.WriteString("- **bookmarks** — saved command configurations\n")
	sb.WriteString("- **amendments** — author amendment audit trail\n")
	sb.WriteString("- **releases** — release metadata\n")
	sb.WriteString("- **ssh_keys** — SSH key records\n")
	sb.WriteString("- **zip_groups** / **zip_group_items** — release archive configurations\n")
	sb.WriteString("- **history** — command execution history\n")
	sb.WriteString("- **tasks** — file-sync task definitions\n")
	sb.WriteString("- **temp_releases** — temporary branch tracking\n\n")
}

// writeLLMInstallation writes the installation section.
func writeLLMInstallation(sb *strings.Builder) {
	sb.WriteString("## Installation\n\n")
	sb.WriteString("| Method | Command |\n")
	sb.WriteString("|--------|---------|\n")
	sb.WriteString("| Windows one-liner | `irm <install-url> \\| iex` |\n")
	sb.WriteString("| Unix one-liner | `curl -fsSL <install-url> \\| sh` |\n")
	sb.WriteString("| From source | `cd gitmap && go build -o ../gitmap .` |\n\n")
}

// writeLLMPatterns writes the common patterns section.
func writeLLMPatterns(sb *strings.Builder) {
	sb.WriteString("## Common Patterns for LLM Assistance\n\n")
	writeLLMPattern(sb, "Find/navigate to a repo",
		"gitmap cd <repo-name>\ngitmap cd repos                    # interactive picker\ngitmap cd -A <alias>               # via alias")
	writeLLMPattern(sb, "Update all repos",
		"gitmap g work                      # activate group\ngitmap g pull                      # pull all\n# or\ngitmap exec pull                   # pull across ALL repos")
	writeLLMPattern(sb, "Organize repos",
		"gitmap scan ~/projects             # discover repos\ngitmap group create work            # create group\ngitmap group add work api web       # add repos\ngitmap alias suggest --apply        # auto-name repos")
	writeLLMPattern(sb, "Release a project",
		"gitmap release-pending             # what's unreleased?\ngitmap changelog-generate --write   # generate changelog\ngitmap release --bump patch --bin   # release with binaries")
	writeLLMPattern(sb, "Clone a project iteration",
		"gitmap cn v++                      # next version\ngitmap cn v++ --delete              # and remove old")
	writeLLMPattern(sb, "Check repo health",
		"gitmap doctor                      # diagnose issues\ngitmap status --all                 # all repo statuses\ngitmap hau                         # check for unpulled commits")
}

// writeLLMPattern writes a single pattern example block.
func writeLLMPattern(sb *strings.Builder, title, commands string) {
	sb.WriteString("### \"I want to " + title + "\"\n")
	sb.WriteString("```bash\n" + commands + "\n```\n\n")
}
