package termhelp

// CommandEntry represents a single CLI command or subcommand line.
type CommandEntry struct {
	Command        string
	Description    string
	HasSubcommands bool
	SubcommandHint string
}

// HelpSection groups related command entries under a styled heading.
type HelpSection struct {
	Title   string
	Color   string
	Entries []CommandEntry
}

// HelpMenu holds complete menu configuration for a command suite.
type HelpMenu struct {
	Title       string
	UsageLines  []string
	Sections    []HelpSection
	FooterFlags []CommandEntry
	Tips        []string
}
