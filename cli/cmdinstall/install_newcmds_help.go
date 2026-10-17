// Package cmdinstall — install_newcmds_help.go: --help texts and help-example
// registry data for `install verify` and `install urls` (program 275, spec 03).
package cmdinstall

// installVerifyHelpText is the full `--help` text for `install verify`.
const installVerifyHelpText = `gitmap install verify — post-install health checklist for this machine

USAGE
  gitmap install verify

CHECKS
  binary version          the running binary exists; reports v<version>
  PATH registered         the binary's directory is on PATH
  symlinks intact         the executable path resolves; symlink targets exist
  no deploy contamination no run.sh / deploy / staging markers next to the binary
  install source          github release vs local build vs other (informational)

EXAMPLE
  gitmap install verify

NOTES
  - Read-only: inspects the installation, changes nothing.
`

// installURLsUsageLine is the one-line usage for argument errors.
const installURLsUsageLine = "usage: gitmap install urls <tag>\n"

// installURLsHelpText is the full `--help` text for `install urls`.
const installURLsHelpText = `gitmap install urls — audit installer download URLs at a tag

USAGE
  gitmap install urls <tag>

EXAMPLE
  gitmap install urls v6.558.0

NOTES
  - Extracts every download URL from install.sh, install.ps1, and the
    Homebrew cask (when present in the repo) at the given tag, then
    HEAD-checks each one (GET fallback on 405).
  - Scripts load from raw.githubusercontent.com at the tag; falls back to
    the local checkout per script when the fetch fails.
  - Read-only.
`

// InstallCmdHelpExample is the canonical help-example shape, field-for-field
// identical to helpdoc.HelpExample (spec 03 §2). It drops verbatim into
// the cli/helpdoc registry (wired by the lead).
type InstallCmdHelpExample struct {
	WhatItDoes string
	Examples   []InstallCmdExample
}

// InstallCmdExample is one concrete invocation + its expected output (spec 03 §2).
type InstallCmdExample struct {
	Invocation     string
	ExpectedOutput string
}

const (
	installVerifyExampleWhat = "Runs a read-only post-install health checklist: binary version, PATH registration, symlink integrity, deploy contamination, and install source."
	installVerifyExampleInv  = "gitmap install verify"
	installVerifyExampleOut  = "✓ binary version (v6.558.0)\n" +
		"✓ PATH registered (~/.local/bin)\n" +
		"✓ symlinks intact (no symlinks in chain)\n" +
		"✓ no deploy contamination\n" +
		"✓ install source (github release)"

	installURLsExampleWhat = "Extracts every download URL from the release's installer scripts at the given tag and HEAD-checks each one."
	installURLsExampleInv  = "gitmap install urls v6.558.0"
	installURLsExampleOut  = "install.sh (github v6.558.0)\n" +
		"install.ps1 (github v6.558.0)\n" +
		"12 urls checked · 12 × 200 ✓ (0.8s)"
)

// InstallVerifyHelpExamples registers the `install verify` help example.
var InstallVerifyHelpExamples = map[string]InstallCmdHelpExample{
	"install verify": {WhatItDoes: installVerifyExampleWhat, Examples: []InstallCmdExample{{Invocation: installVerifyExampleInv, ExpectedOutput: installVerifyExampleOut}}},
}

// InstallURLsHelpExamples registers the `install urls` help example.
var InstallURLsHelpExamples = map[string]InstallCmdHelpExample{
	"install urls": {WhatItDoes: installURLsExampleWhat, Examples: []InstallCmdExample{{Invocation: installURLsExampleInv, ExpectedOutput: installURLsExampleOut}}},
}
