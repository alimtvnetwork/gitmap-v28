package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestProfileTitlePrefixAppliedToReplayedSubjects writes a saved
// profile that pins TitlePrefix="[demo] " to the source workspace,
// then runs commit-in with --profile demo. Every replayed commit's
// subject must be `[demo] <original>`.
//
// This proves the full pickProfile → Resolve → message-pipeline chain
// is wired end-to-end.
func TestProfileTitlePrefixAppliedToReplayedSubjects(t *testing.T) {
	src := NewRepo(t, "src")
	input := NewRepo(t, "input")
	input.Commit("a.txt", "1\n", "feat: add a", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	input.Commit("b.txt", "2\n", "fix: tweak b", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))

	saveProfile(t, src.Path, &profile.Profile{
		Name:          "demo",
		SchemaVersion: profile.CurrentSchemaVersion,
		TitlePrefix:   "[demo] ",
	})

	raw := NewRawArgs(src.Path, input.Path)
	raw.ProfileName = "demo"
	res := Run(t, raw)
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s", res.ExitCode, res.Stderr)
	}
	src.AssertHasSubject(t, "[demo] feat: add a")
	src.AssertHasSubject(t, "[demo] fix: tweak b")
}

// TestCliTitlePrefixOverridesProfileValue layers a CLI --title-prefix
// on top of a profile that already sets TitlePrefix. Spec §5.6
// precedence is CLI > profile > defaults, so the CLI value must win.
func TestCliTitlePrefixOverridesProfileValue(t *testing.T) {
	src := NewRepo(t, "src")
	input := NewRepo(t, "input")
	input.Commit("a.txt", "1\n", "seed", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	saveProfile(t, src.Path, &profile.Profile{
		Name:          "demo",
		SchemaVersion: profile.CurrentSchemaVersion,
		TitlePrefix:   "[profile] ",
	})

	raw := NewRawArgs(src.Path, input.Path)
	raw.ProfileName = "demo"
	raw.TitlePrefix = "[cli] " // must beat the profile value
	res := Run(t, raw)
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s", res.ExitCode, res.Stderr)
	}
	src.AssertHasSubject(t, "[cli] seed")
	for _, c := range src.LogFirstParent(t) {
		if strings.HasPrefix(c.Subject, "[profile] ") {
			t.Fatalf("profile prefix leaked through CLI override: %q", c.Subject)
		}
	}
}

// TestMissingNamedProfileExitsWithProfileMissingCode covers the
// "fatal-on-miss" branch of pickProfile: --profile <unknown> with no
// matching file on disk must abort with CommitInExitProfileMissing
// (=6) and produce zero commits.
func TestMissingNamedProfileExitsWithProfileMissingCode(t *testing.T) {
	src := NewRepo(t, "src")
	input := NewRepo(t, "input")
	input.Commit("a.txt", "1\n", "seed", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	raw := NewRawArgs(src.Path, input.Path)
	raw.ProfileName = "does-not-exist"
	res := Run(t, raw)

	if res.ExitCode != constants.CommitInExitProfileMissing {
		t.Fatalf("exit=%d, want CommitInExitProfileMissing (%d)\nstderr=%s",
			res.ExitCode, constants.CommitInExitProfileMissing, res.Stderr)
	}
	src.AssertCommitCount(t, 0)
}

// saveProfile writes `p` to <workspaceRoot>/.gitmap/commit-in/profiles
// using the production SaveToDisk path. Test-private helper.
func saveProfile(t *testing.T, workspaceRoot string, p *profile.Profile) {
	t.Helper()
	if err := profile.SaveToDisk(workspaceRoot, p, true); err != nil {
		t.Fatalf("save profile %q: %v", p.Name, err)
	}
}
