package movemerge

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// finalizeURLSides commits + pushes any URL endpoint that should
// receive the change. Direction decides which side(s) are written.
func finalizeURLSides(left, right Endpoint, dir Direction, opts Options) error {
	if opts.NoCommit || opts.DryRun {
		return nil
	}
	if writesRight(dir) && right.Kind == EndpointURL {
		if err := commitAndPushOne(right, otherDisplay(left, right), opts); err != nil {
			return err
		}
	}
	if writesLeft(dir) && left.Kind == EndpointURL {
		if err := commitAndPushOne(left, otherDisplay(right, left), opts); err != nil {
			return err
		}
	}

	return nil
}

// writesRight reports whether the operation modifies RIGHT.
func writesRight(dir Direction) bool {
	return dir == DirBoth || dir == DirRightOnly
}

// writesLeft reports whether the operation modifies LEFT.
func writesLeft(dir Direction) bool {
	return dir == DirBoth || dir == DirLeftOnly
}

// otherDisplay returns the other side's display string for messages.
func otherDisplay(other, _ Endpoint) string {
	return other.DisplayName
}

// commitAndPushOne stages, commits, and pushes a URL endpoint.
func commitAndPushOne(ep Endpoint, otherDisp string, opts Options) error {
	logf(opts.LogPrefix, "committing in %s ...", ep.DisplayName)
	msg := fmt.Sprintf(opts.CommitMsgFmt, otherDisp)
	sha, err := AddCommitPush(ep.WorkingDir, msg, !opts.NoPush)
	if err != nil {
		logErr(opts.LogPrefix, fmt.Sprintf(constants.ErrMMPushFailFmt, sha))

		return err
	}
	if sha != "" {
		logIndent(opts.LogPrefix, "commit %s %q", shortSHA(sha), msg)
	}
	if !opts.NoPush {
		logf(opts.LogPrefix, "pushing %s ...", ep.DisplayName)
		logIndent(opts.LogPrefix, "push OK")
	}

	return nil
}

// shortSHA truncates a 40-char hex SHA to 7 chars for log lines.
func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}

	return sha
}
