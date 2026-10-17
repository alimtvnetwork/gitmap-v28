// Package cmdhelp — help_examples.go: --help texts and help-example registry
// data for `help topic` / `ht` and `howto` (program 275, spec 03 format).
package cmdhelp

// isHelpRequested reports whether the user asked for help.
func isHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}

	return false
}

// helpTopicHelpText is the full `--help` text for `help topic` / `ht`.
const helpTopicHelpText = `gitmap help topic — topic index: the commands for one workflow, in order

USAGE
  gitmap help topic <name>     (alias: gitmap ht <name>)
  gitmap help topic            (list all topics)

TOPICS
  releases    ship a version: summary → release-dry → release-notes → release → pe
  pipeline    triage a red pipeline: pe → pe triage → pe rerun
  fleet       fleet machines: nodes ping → nodes version → nodes jump
  migration   move a repo set: scan export → repo import → scan
  scanner     index truth: scan → scan verify → scan diff
  themes      palettes: theme list → theme preview → --theme

EXAMPLE
  gitmap help topic releases
  gitmap ht pipeline

NOTES
  - Read-only: prints a learning path, changes nothing.
  - Unknown topic names print the valid topic list.
`

// howtoUsageLine is the one-line usage for missing-argument errors.
const howtoUsageLine = "usage: gitmap howto <task description>\n"

// howtoHelpText is the full `--help` text for `howto`.
const howtoHelpText = `gitmap howto — keyword task→recipe matching: "how do I …?" in one shot

USAGE
  gitmap howto <task description...>

EXAMPLE
  gitmap howto "ship a new version"
  gitmap howto fix red pipeline

NOTES
  - Pure keyword matching, no AI calls.
  - Unknown tasks print the topic index; continue with: gitmap help topic <name>
`

// HelpExample is the canonical help-example shape, field-for-field
// identical to helpdoc.HelpExample (spec 03 §2). It drops verbatim into
// the cli/helpdoc registry (wired by the lead).
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
	helpTopicExampleWhat = "Prints the ordered command list for one workflow topic (releases, pipeline, fleet, migration, scanner, themes); with no topic, lists all topics."
	helpTopicExampleInv  = "gitmap help topic releases"
	helpTopicExampleOut  = "TOPIC releases: Ship a version: check state, dry-run, notes, release, watch CI.\n" +
		"  1) gitmap summary — repo state before you ship\n" +
		"  2) gitmap release-dry — preview the release with no side effects\n" +
		"  3) gitmap release-notes — changelog between tags\n" +
		"  4) gitmap release vX.Y.Z — cut the release\n" +
		"  5) gitmap pe — watch the release pipeline"

	howtoExampleWhat = "Matches a plain-language task to an ordered command recipe via keyword matching (no AI calls); unknown tasks print the topic index."
	howtoExampleInv  = `gitmap howto "ship a new version"`
	howtoExampleOut  = "RECIPE — ship a new version:\n" +
		"  1) gitmap summary\n" +
		"  2) gitmap release-dry\n" +
		"  3) gitmap release-notes\n" +
		"  4) gitmap release vX.Y.Z\n" +
		"  5) gitmap pe\n" +
		"see: gitmap help topic releases"
)

// HelpTopicHelpExamples registers the `help topic` / `ht` help example,
// keyed by dispatcher verb AND alias (spec 03 §3).
var HelpTopicHelpExamples = map[string]HelpExample{
	"help topic": {WhatItDoes: helpTopicExampleWhat, Examples: []Example{{Invocation: helpTopicExampleInv, ExpectedOutput: helpTopicExampleOut}}},
	"ht":         {WhatItDoes: helpTopicExampleWhat, Examples: []Example{{Invocation: "gitmap ht releases", ExpectedOutput: helpTopicExampleOut}}},
}

// HowtoHelpExamples registers the `howto` help example.
var HowtoHelpExamples = map[string]HelpExample{
	"howto": {WhatItDoes: howtoExampleWhat, Examples: []Example{{Invocation: howtoExampleInv, ExpectedOutput: howtoExampleOut}}},
}
