package cmd

// clonetermblock_golden_test.go — golden-file fixtures verifying the
// per-repo `--output terminal` block produced by every clone-related
// command (clone, clone-next, clone-now, clone-from, clone-pick) is
// byte-identical to a checked-in expected file.
//
// Why goldens vs. inline strings: the cmd: line concatenates many
// pieces (binary, subcommand, optional --filter / --branch / --depth
// flags, URL, dest). A regression that re-orders or drops a token is
// trivial to introduce and hard to spot in inline test strings. A
// golden file makes the diff obvious in PR review and lets CI fail
// loudly if the format drifts.
//
// Update procedure: run `go test ./gitmap/cmd -run TestCloneTermBlock
// _Golden -update` (the -update flag rewrites the .golden files from
// the current output). The flag is intentionally local to this file
// so a typo elsewhere can't accidentally regenerate fixtures.

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/goldenguard"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/render"
)

// updateGolden, when true, rewrites .golden files instead of asserting.
// Local to this test file by design — see file-header rationale.
var updateGolden = flag.Bool("update", false,
	"rewrite clonetermblock .golden fixtures from current output")

// goldenCase pairs a fixture file name with the CloneTermBlockInput
// the production code path would build for that command. Each input
// is constructed using the SAME helpers (buildCloneCommand,
// pickCmdBranch, RenderRepoTermBlock) the live code uses — the test
// only differs in WHERE the input comes from (hand-written here vs.
// derived from a real Plan/Row at runtime).
type goldenCase struct {
	name    string // logical command name, used in failure messages
	fixture string // file name under testdata/
	input   CloneTermBlockInput
}

// cloneTermGoldenCases enumerates fixtures for every clone-related
// command in BOTH URL modes (HTTPS and SSH). Two-mode coverage matters
// because URL formatting flows untouched through the cmd: line — a
// regression that mangles SSH-style `git@host:path` (e.g. by URL-
// parsing it as if it were https://) wouldn't show up in an HTTPS-only
// fixture. Keeping one fixture per (command, url-mode) pair makes the
// failure message in CI point straight at the broken combination.
//
// Inputs intentionally use realistic values (real-looking URLs,
// branches, dest paths) so a reviewer can sanity-check the cmd: line
// by reading the golden alone.
func cloneTermGoldenCases() []goldenCase {
	const (
		repoName = "scripts-fixer"
		httpsURL = "https://github.com/owner/scripts-fixer.git"
		sshURL   = "git@github.com:owner/scripts-fixer.git"
	)

	return []goldenCase{
		{
			name:    "clone-https",
			fixture: "clonetermblock_clone.golden",
			// URL-driven `gitmap clone <url>` — matches clonetermurl.go:
			// non-nil empty CmdExtraArgsPre = explicit "no -b" sentinel.
			input: CloneTermBlockInput{
				Index:           1,
				Name:            repoName,
				Branch:          "main",
				BranchSource:    "remote HEAD",
				OriginalURL:     httpsURL,
				TargetURL:       httpsURL,
				Dest:            repoName,
				CmdBranch:       "",
				CmdExtraArgsPre: []string{},
			},
		},
		{
			name:    "clone-ssh",
			fixture: "clonetermblock_clone_ssh.golden",
			// Same shape as clone-https but with scp-style SSH URL —
			// guards against a future change that URL-parses inputs
			// (which would mangle `git@host:path` into nonsense).
			input: CloneTermBlockInput{
				Index:           1,
				Name:            repoName,
				Branch:          "main",
				BranchSource:    "remote HEAD",
				OriginalURL:     sshURL,
				TargetURL:       sshURL,
				Dest:            repoName,
				CmdBranch:       "",
				CmdExtraArgsPre: []string{},
			},
		},
		{
			name:    "clone-next-https",
			fixture: "clonetermblock_clonenext.golden",
			// clone-next routes through the same URL-driven helper as
			// `gitmap clone <url>`, so the fixture matches the clone
			// case shape (no -b, remote-HEAD branch source).
			input: CloneTermBlockInput{
				Index:           1,
				Name:            repoName,
				Branch:          "main",
				BranchSource:    "remote HEAD",
				OriginalURL:     httpsURL,
				TargetURL:       httpsURL,
				Dest:            repoName,
				CmdBranch:       "",
				CmdExtraArgsPre: []string{},
			},
		},
		{
			name:    "clone-next-ssh",
			fixture: "clonetermblock_clonenext_ssh.golden",
			// SSH counterpart of clone-next-https. clone-next reuses
			// the URL-driven helper, so the same SSH-pass-through
			// guarantee applies.
			input: CloneTermBlockInput{
				Index:           1,
				Name:            repoName,
				Branch:          "main",
				BranchSource:    "remote HEAD",
				OriginalURL:     sshURL,
				TargetURL:       sshURL,
				Dest:            repoName,
				CmdBranch:       "",
				CmdExtraArgsPre: []string{},
			},
		},
		{
			name:    "clone-now-ssh",
			fixture: "clonetermblock_clonenow.golden",
			// clone-now manifest row with explicit branch + SSH URL.
			// CmdBranch=row.Branch ⇒ -b is rendered; no extra args.
			input: CloneTermBlockInput{
				Index:        3,
				Name:         repoName,
				Branch:       "develop",
				BranchSource: "manifest",
				OriginalURL:  sshURL,
				TargetURL:    sshURL,
				Dest:         "repos/" + repoName,
				CmdBranch:    "develop",
			},
		},
		{
			name:    "clone-now-https",
			fixture: "clonetermblock_clonenow_https.golden",
			// HTTPS counterpart of clone-now-ssh. Same row shape,
			// same -b rendering — only the URL token changes.
			input: CloneTermBlockInput{
				Index:        3,
				Name:         repoName,
				Branch:       "develop",
				BranchSource: "manifest",
				OriginalURL:  httpsURL,
				TargetURL:    httpsURL,
				Dest:         "repos/" + repoName,
				CmdBranch:    "develop",
			},
		},
		{
			name:    "clone-from-https",
			fixture: "clonetermblock_clonefrom.golden",
			// clone-from with both pinned branch and depth=1. The
			// executor places --depth AFTER -b, so CmdExtraArgsPost
			// carries `--depth=1`.
			input: CloneTermBlockInput{
				Index:            2,
				Name:             repoName,
				Branch:           "main",
				BranchSource:     "manifest",
				OriginalURL:      httpsURL,
				TargetURL:        httpsURL,
				Dest:             repoName,
				CmdBranch:        "main",
				CmdExtraArgsPost: []string{"--depth=1"},
			},
		},
		{
			name:    "clone-from-ssh",
			fixture: "clonetermblock_clonefrom_ssh.golden",
			// SSH counterpart with the same depth+branch combo —
			// also guards against drift in the locked `--depth=N`
			// (joined) format across URL modes.
			input: CloneTermBlockInput{
				Index:            2,
				Name:             repoName,
				Branch:           "main",
				BranchSource:     "manifest",
				OriginalURL:      sshURL,
				TargetURL:        sshURL,
				Dest:             repoName,
				CmdBranch:        "main",
				CmdExtraArgsPost: []string{"--depth=1"},
			},
		},
		{
			name:    "clone-pick-https",
			fixture: "clonetermblock_clonepick.golden",
			// clone-pick uses partial-clone flags + long-form
			// --branch/--depth, all in CmdExtraArgsPre. CmdBranch
			// stays empty so no `-b` is rendered.
			input: CloneTermBlockInput{
				Index:        1,
				Name:         repoName,
				Branch:       "main",
				BranchSource: "manifest",
				OriginalURL:  httpsURL,
				TargetURL:    httpsURL,
				Dest:         repoName,
				CmdBranch:    "",
				CmdExtraArgsPre: []string{
					"--filter=blob:none", "--no-checkout",
					"--branch", "main",
					"--depth", "1",
				},
			},
		},
		{
			name:    "clone-pick-ssh",
			fixture: "clonetermblock_clonepick_ssh.golden",
			// SSH counterpart — same long-form flag block, only the
			// positional URL token changes.
			input: CloneTermBlockInput{
				Index:        1,
				Name:         repoName,
				Branch:       "main",
				BranchSource: "manifest",
				OriginalURL:  sshURL,
				TargetURL:    sshURL,
				Dest:         repoName,
				CmdBranch:    "",
				CmdExtraArgsPre: []string{
					"--filter=blob:none", "--no-checkout",
					"--branch", "main",
					"--depth", "1",
				},
			},
		},
	}
}

// renderGoldenBlock reproduces the exact byte sequence the production
// `--output terminal` path writes for one repo: build the cmd via
// buildCloneCommand, then render via render.RenderRepoTermBlock.
// Kept here (vs. calling maybePrintCloneTermBlock) so the test
// doesn't depend on os.Stdout and produces deterministic output.
func renderGoldenBlock(t *testing.T, in CloneTermBlockInput) []byte {
	t.Helper()
	var buf bytes.Buffer
	err := render.RenderRepoTermBlock(&buf, render.RepoTermBlock{
		Index:        in.Index,
		Name:         in.Name,
		Branch:       in.Branch,
		BranchSource: in.BranchSource,
		OriginalURL:  in.OriginalURL,
		TargetURL:    in.TargetURL,
		CloneCommand: buildCloneCommand(in),
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	return buf.Bytes()
}

// TestCloneTermBlock_Golden is the CI guard: every clone-related
// command's per-repo block must match its checked-in fixture. A diff
// indicates either an intentional format change (update the golden
// with `-update` AND GITMAP_ALLOW_GOLDEN_UPDATE=1, then call it out
// in the PR) or a regression. The dual gate (flag + env) prevents a
// stray `-update` in a CI invocation from silently rewriting bytes.
func TestCloneTermBlock_Golden(t *testing.T) {
	for _, tc := range cloneTermGoldenCases() {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", tc.fixture)
			got := renderGoldenBlock(t, tc.input)

			if goldenguard.AllowUpdate(t, *updateGolden) {
				if err := os.WriteFile(path, got, 0o644); err != nil {
					t.Fatalf("update golden %s: %v", path, err)
				}
				t.Logf("updated %s", path)

				return
			}

			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden %s: %v (run with -update "+
					"AND GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)",
					path, err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("%s mismatch\n--- want (%s) ---\n%s"+
					"\n--- got ---\n%s", tc.name, path, want, got)
			}
		})
	}
}
