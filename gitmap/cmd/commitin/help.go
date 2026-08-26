// Package commitin — help.go prints detailed multi-section usage documentation.
package commitin

import (
	"strings"
)

// PrintCommitInHelp returns detailed documentation and JSON examples for commit-in.
func PrintCommitInHelp() string {
	var b strings.Builder
	b.WriteString("gitmap commit-in / commit-write: Automated Commit Engine\n")
	b.WriteString("========================================================\n\n")
	b.WriteString("Features:\n")
	b.WriteString("  - Author Rotation: automatically cycles author metadata\n")
	b.WriteString("  - SEO Commit Scheduling: generates formatted commit messages with custom templates\n")
	b.WriteString("  - FuncIntel: AST parsing of changed functions to draft intelligent messages\n")
	b.WriteString("  - Deduplication: prevents repeated messages across repository history\n\n")
	b.WriteString("Usage:\n")
	b.WriteString("  gitmap commit-in --seo-url https://example.com --interval 60-120\n")
	b.WriteString("  gitmap commit-in --profile work --dry-run\n")
	return b.String()
}
