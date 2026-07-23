package orchestrator

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// maybeSaveProfile persists ctx.Resolved as a named profile when
// --save-profile <name> was passed. Honors --save-profile-overwrite
// and --set-default. Called once after setUp and before the per-input
// pipeline so a save failure aborts BEFORE any commits are replayed.
func maybeSaveProfile(ctx *runContext, stderr io.Writer) int {
	if ctx.Raw.SaveProfileName == "" {
		return constants.CommitInExitOk
	}
	name := strings.TrimSpace(ctx.Raw.SaveProfileName)
	if name == "" {
		fmt.Fprintf(stderr, constants.CommitInErrBadArgs, "--save-profile name is empty")
		return constants.CommitInExitBadArgs
	}
	fmt.Fprintf(stderr, constants.CommitInMsgPhaseSaveProfile, name)
	p := profile.BuildFromResolved(profile.BuildArgs{
		Name:           name,
		SourceRepoPath: ctx.Source.Path,
		IsDefault:      ctx.Raw.SetDefault,
		Resolved:       ctx.Resolved,
	})
	return persistProfile(ctx, p, stderr)
}

func persistProfile(ctx *runContext, p *profile.Profile, stderr io.Writer) int {
	if p.IsDefault {
		if err := profile.ClearOtherDefaults(ctx.Paths.SourceRoot, ctx.Source.Path, p.Name); err != nil {
			fmt.Fprintf(stderr, constants.CommitInErrDbWrite, err)
			return constants.CommitInExitDbFailed
		}
	}
	if err := profile.SaveToDisk(ctx.Paths.SourceRoot, p, ctx.Raw.SaveProfileOverwrite); err != nil {
		// Distinguish "exists" (user fixable) from generic IO.
		if isExistsError(err) {
			fmt.Fprintf(stderr, constants.CommitInErrSaveProfileExists+"\n", p.Name)
			return constants.CommitInExitBadArgs
		}
		fmt.Fprintf(stderr, constants.CommitInErrDbWrite, err)
		return constants.CommitInExitDbFailed
	}
	return constants.CommitInExitOk
}

// isExistsError matches the sentinel string SaveToDisk uses when
// refusing to overwrite. Kept as a string check (rather than a typed
// error) to avoid coupling the profile package to orchestrator
// classification — this path is the ONLY caller that needs the split.
func isExistsError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
