// Package installer — ui_help.go prints detailed help text and usage.
package installer

import "strings"

// PrintDetailedHelp returns comprehensive help documentation for installer commands.
func PrintDetailedHelp() string {
	var b strings.Builder
	b.WriteString("Gitmap Installer Management System\n")
	b.WriteString("==================================\n\n")
	b.WriteString("Commands:\n")
	b.WriteString("  gitmap installer create <name>            Create a new installer script\n")
	b.WriteString("  gitmap installer update <slug>            Update and increment version\n")
	b.WriteString("  gitmap installer update-win <slug>        Update Windows configuration\n")
	b.WriteString("  gitmap installer install-win <slug>       Run Windows installer\n")
	b.WriteString("  gitmap installer export <slug>            Export installer to zip\n")
	b.WriteString("  gitmap installer export-all               Export all installers to zip\n")
	b.WriteString("  gitmap installer import [path]            Import from zip or json\n")
	b.WriteString("  gitmap installer reset <slug>             Reset specific installer\n")
	b.WriteString("  gitmap installer undo-version <slug>      Roll back to previous version\n")
	b.WriteString("  gitmap installer redo-version <slug>      Advance to redone version\n")
	b.WriteString("  gitmap installer revert-version <s <v>    Revert to exact version\n")
	b.WriteString("  gitmap installer ls [os]                  List all scripts in table\n")
	return b.String()
}
