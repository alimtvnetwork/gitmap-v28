// Package cmdrelease — release_newcmds_help.go: --help texts and help-example
// registry data for `release verify`, `release winscan`, and `release train`
// / `rtrain` (program 275, spec 03 format).
package cmdrelease

// isReleaseNewCmdHelpRequested reports whether the user asked for help.
func isReleaseNewCmdHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}

	return false
}

// releaseVerifyUsageLine is the one-line usage for argument errors.
const releaseVerifyUsageLine = "usage: gitmap release verify <repo> <tag>\n"

// releaseVerifyHelpText is the full `--help` text for `release verify`.
const releaseVerifyHelpText = `gitmap release verify — post-release integrity audit for one tag

USAGE
  gitmap release verify <repo> <tag>

  <repo>  owner/repo, or a bare repo name (resolves to alimtvnetwork/<repo>)
  <tag>   release tag, e.g. v6.558.0

CHECKS
  assets     published asset count
  draft      draft flag must be false
  checksums  a checksum/sha256 asset is present
  tag→HEAD   the tag's commit equals the repo default-branch HEAD

EXAMPLE
  gitmap release verify gitmap-v28 v6.558.0

NOTES
  - Shells out to the 'gh' CLI; read-only, changes nothing.
`

// releaseWinScanUsageLine is the one-line usage for argument errors.
const releaseWinScanUsageLine = "usage: gitmap release winscan <tag>\n"

// releaseWinScanHelpText is the full `--help` text for `release winscan`.
const releaseWinScanHelpText = `gitmap release winscan — pre-publish Defender heuristic scan of the Windows asset

USAGE
  gitmap release winscan <tag>

EXAMPLE
  gitmap release winscan v6.558.0

NOTES
  - Downloads the *windows*amd64*.zip asset for the tag into a temp dir
    (repo resolves from the current directory via 'gh'), extracts the exe,
    and byte-scans it against the 7 Windows Defender heuristic trigger
    patterns from RCA 273 (commit 385aba38e):
    P1 reg.exe Run-key persistence · P2 vssadmin shadow copy ·
    P3 Winlogon auto-login credential · P4 hidden PowerShell WMI enum ·
    P5 spawned taskkill · P6 delayed self-deletion · P7 reg.exe Run-key delete
  - Any pattern hit fails the command — review before publishing.
  - Read-only.
`

// releaseTrainUsageLine is the one-line usage for argument errors.
const releaseTrainUsageLine = "usage: gitmap release train <repo...> [--order <csv>] [--yes]\n"

// releaseTrainHelpText is the full `--help` text for `release train`.
const releaseTrainHelpText = `gitmap release train — ordered multi-repo release (dry-run by default)

USAGE
  gitmap release train <repo...> [--order <csv>] [--yes]
  gitmap rtrain <repo...>        (alias)

  <repo...>  owner/repo or bare repo names (bare → alimtvnetwork/<repo>)

FLAGS
  --order <csv>   repo order, overrides positional order
  --yes, -y       execute the train (WITHOUT it: dry-run prints the plan only)

PER-REPO STEPS
  1) version bump  next patch tag from the latest release (vX.Y.Z → vX.Y.Z+1)
  2) CI gate       latest workflow run on the default branch must be success
  3) tag           created via gh release create
  4) release notes generated notes from commits since the previous tag

EXAMPLE
  gitmap release train agm gitmap-v28
  gitmap release train agm gitmap-v28 --order agm,gitmap-v28 --yes

NOTES
  - DEFAULT IS DRY-RUN: nothing is created until --yes is passed.
  - Execution aborts at the first repo whose CI gate is not green.
`

// ReleaseCmdHelpExample is the canonical help-example shape, field-for-field
// identical to helpdoc.HelpExample (spec 03 §2). It drops verbatim into
// the cli/helpdoc registry (wired by the lead).
type ReleaseCmdHelpExample struct {
	WhatItDoes string
	Examples   []ReleaseCmdExample
}

// ReleaseCmdExample is one concrete invocation + its expected output (spec 03 §2).
type ReleaseCmdExample struct {
	Invocation     string
	ExpectedOutput string
}

const (
	releaseVerifyExampleWhat = "Post-release integrity audit via gh release view: asset count, draft flag, checksums presence, and tag→commit == repo HEAD."
	releaseVerifyExampleInv  = "gitmap release verify gitmap-v28 v6.558.0"
	releaseVerifyExampleOut  = "v6.558.0 (alimtvnetwork/gitmap-v28):\n" +
		"  assets:    16\n" +
		"  draft:     false\n" +
		"  checksums: present ✓\n" +
		"  tag→abc0f2fa = HEAD ✓"

	releaseWinScanExampleWhat = "Pre-publish scan of the Windows release asset against the 7 Windows Defender heuristic trigger patterns from RCA 273."
	releaseWinScanExampleInv  = "gitmap release winscan v6.558.0"
	releaseWinScanExampleOut  = "gitmap_6.558.0_windows_amd64.zip (gitmap.exe, 18.4 MB):\n" +
		"  0/7 trigger patterns ✓"

	releaseTrainExampleWhat = "Ordered multi-repo release: version bump → CI gate → tag → release notes per repo. Dry-run by default; --yes executes."
	releaseTrainExampleInv  = "gitmap release train agm gitmap-v28"
	releaseTrainExampleOut  = "RELEASE TRAIN [DRY-RUN] — 2 repos:\n" +
		"  [1/2] alimtvnetwork/agm\n" +
		"    bump: v4.216.2 → v4.216.3 · CI gate: green ✓ (completed/success) · tag: v4.216.3 · notes: 12 commits since v4.216.2\n" +
		"  [2/2] alimtvnetwork/gitmap-v28\n" +
		"    bump: v6.558.0 → v6.559.0 · CI gate: green ✓ (completed/success) · tag: v6.559.0 · notes: 8 commits since v6.558.0\n" +
		"dry-run: no changes made — pass --yes to execute"
)

// ReleaseVerifyHelpExamples registers the `release verify` help example.
var ReleaseVerifyHelpExamples = map[string]ReleaseCmdHelpExample{
	"release verify": {WhatItDoes: releaseVerifyExampleWhat, Examples: []ReleaseCmdExample{{Invocation: releaseVerifyExampleInv, ExpectedOutput: releaseVerifyExampleOut}}},
}

// ReleaseWinScanHelpExamples registers the `release winscan` help example.
var ReleaseWinScanHelpExamples = map[string]ReleaseCmdHelpExample{
	"release winscan": {WhatItDoes: releaseWinScanExampleWhat, Examples: []ReleaseCmdExample{{Invocation: releaseWinScanExampleInv, ExpectedOutput: releaseWinScanExampleOut}}},
}

// ReleaseTrainHelpExamples registers the `release train` / `rtrain` help example,
// keyed by dispatcher verb AND alias (spec 03 §3).
var ReleaseTrainHelpExamples = map[string]ReleaseCmdHelpExample{
	"release train": {WhatItDoes: releaseTrainExampleWhat, Examples: []ReleaseCmdExample{{Invocation: releaseTrainExampleInv, ExpectedOutput: releaseTrainExampleOut}}},
	"rtrain":        {WhatItDoes: releaseTrainExampleWhat, Examples: []ReleaseCmdExample{{Invocation: "gitmap rtrain agm gitmap-v28", ExpectedOutput: releaseTrainExampleOut}}},
}
