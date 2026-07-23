package replay

import (
	"fmt"
	"strings"
)

// DetectClobbers returns the subset of plan.Files whose target-side
// blob (at HEAD) differs from the incoming source blob. An empty
// result means the commit is "clean" — applying it cannot lose any
// existing target content.
//
// The check is HEAD-relative (NOT working-tree) because commit-in
// only ever writes via `git update-index --cacheinfo`; the working
// tree is irrelevant.
func DetectClobbers(p Plan) ([]string, error) {
	head, err := readHead(p.TargetRepoDir)
	if err != nil || head == "" {
		// No HEAD = empty repo = nothing to clobber.
		return nil, nil
	}
	var clobbers []string
	for _, rel := range p.Files {
		isClobber, err := oneFileClobbers(p, head, rel)
		if err != nil {
			return nil, fmt.Errorf("clobber check %s: %w", rel, err)
		}
		if isClobber {
			clobbers = append(clobbers, rel)
		}
	}
	return clobbers, nil
}

// oneFileClobbers compares the source blob hash @ p.SourceSha to the
// target blob hash @ HEAD for one path. Missing-on-target → false
// (an add never clobbers). Same hash → false (no-op rewrite).
func oneFileClobbers(p Plan, head, rel string) (bool, error) {
	srcHash, err := blobHashAt(p.SourceRepoDir, p.SourceSha, rel)
	if err != nil {
		return false, fmt.Errorf("source side: %w", err)
	}
	tgtHash, err := blobHashAt(p.TargetRepoDir, head, rel)
	if err != nil {
		// Target has no such path at HEAD => add, not clobber.
		return false, nil
	}
	return srcHash != "" && tgtHash != "" && srcHash != tgtHash, nil
}

// blobHashAt returns the 40-char blob SHA git records for `<sha>:rel`,
// or "" with a non-nil error when the path is absent.
func blobHashAt(dir, sha, rel string) (string, error) {
	out, err := gitRunner(dir, "rev-parse", sha+":"+rel)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
