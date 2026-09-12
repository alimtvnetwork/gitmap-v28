// Package constants — constants_templates_cli.go: CLI identifiers for the
// `gitmap templates ...` discovery command and its subcommands.
package constants

// Top-level command. // gitmap:cmd top-level
const (
	CmdTemplates      = "templates"
	CmdTemplatesAlias = "tpl"
)

// Templates subcommand aliases that should appear in shell tab-completion
// alongside the top-level `templates` / `tpl` commands. Users frequently
// type the alias directly after `templates` (e.g. `gitmap tpl td --lang go`),
// so surfacing them in completion makes the alias discoverable.
//
// The full subcommand IDs (`diff`, `init`) are intentionally skipped via
// `// gitmap:cmd skip` because:
//   - `diff` is already registered as a top-level command (folder-tree
//     diff via gitmap/cmd/diff.go); re-listing it here would be a no-op
//     for the union but conceptually incorrect.
//   - `init` is not a top-level gitmap command at all; surfacing it
//     standalone would mislead users into typing `gitmap init`.
//
// gitmap:cmd top-level
const (
	CmdTemplatesDiff      = "diff" // gitmap:cmd skip
	CmdTemplatesDiffAlias = "td"
	CmdTemplatesInit      = "init" // gitmap:cmd skip
	CmdTemplatesInitAlias = "ti"
)
