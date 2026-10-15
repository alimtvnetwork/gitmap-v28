package cmdmerge

// merge_write_help.go — help text and help-example registry entry for
// `gitmap merge write` / `gitmap mr` (program 273 WS3, spec 03 + spec 05).

import "strings"

// isHelpRequested reports whether the user asked for help.
func isHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

// mwHelpText is the full `--help` text for `merge write`.
const mwHelpText = `gitmap merge write — open a GitHub pull request from the CLI

USAGE
  gitmap merge write --title <title> [--body <body>] [--source <branch>] [--target <branch>] [--draft] [--yes]
  gitmap mr             (alias; identical behavior)

FLAGS
  --title <text>     PR title (required)
  --body <text>      PR body, markdown allowed (default: "")
  --source <branch>  Head branch (default: current git branch; must already be pushed)
  --target <branch>  Base branch (default: the repo's default branch on GitHub)
  --draft            Open the PR as a draft
  --yes, -y          Skip the confirmation prompt

EXAMPLE
  gitmap merge write --title "Migrate Lovable app" --body "Ports the export." --source lovable-export --target main --draft

NOTES
  - v1 supports github.com remotes only; other hosts get a clear error.
  - Auth goes through the 'gh' CLI: 'gh auth login', or GITHUB_TOKEN / GH_TOKEN.
  - The source branch must already exist on the remote (push it first).
`

// HelpExample is the canonical help-example shape, field-for-field
// identical to helpdoc.HelpExample (spec 03 §2; Worker 2's registry at
// cli/helpdoc/examples.go uses WhatItDoes + Examples). MergeWriteHelpExamples
// drops verbatim into that registry.
type HelpExample struct {
	WhatItDoes string
	Examples   []Example
}

// Example is one concrete invocation + its expected output (spec 03 §2).
type Example struct {
	Invocation     string
	ExpectedOutput string
}

const (
	mwExampleWhat = "Opens a pull request from a source branch into a target branch on the repo's host (GitHub), optionally as a draft."
	mwExampleInv  = `gitmap merge write --title "Migrate Lovable app to Vite shell" --body "Ports the Lovable export into the TanStack Start shell." --source lovable-export --target main --draft`
	mwExampleOut  = "pr #42 opened: lovable-export -> main (draft)\nhttps://github.com/alimtvnetwork/my-app/pull/42"
)

// MergeWriteHelpExamples registers the `merge write` / `mr` help example,
// keyed by dispatcher verb AND alias (spec 03 §3: every command and alias
// gets an entry).
var MergeWriteHelpExamples = map[string]HelpExample{
	"merge write": {WhatItDoes: mwExampleWhat, Examples: []Example{{Invocation: mwExampleInv, ExpectedOutput: mwExampleOut}}},
	"mr":          {WhatItDoes: mwExampleWhat, Examples: []Example{{Invocation: strings.Replace(mwExampleInv, "gitmap merge write", "gitmap mr", 1), ExpectedOutput: mwExampleOut}}},
}
