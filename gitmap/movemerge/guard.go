package movemerge

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// GuardEndpoints checks LEFT and RIGHT do not collide on disk.
func GuardEndpoints(left, right Endpoint) error {
	lAbs, _ := filepath.Abs(left.WorkingDir)
	rAbs, _ := filepath.Abs(right.WorkingDir)
	if lAbs == rAbs {
		return fmt.Errorf(constants.ErrMMSameFolderFmt, lAbs)
	}
	if isStrictAncestor(lAbs, rAbs) || isStrictAncestor(rAbs, lAbs) {
		return fmt.Errorf(constants.ErrMMNestedFmt, lAbs, rAbs)
	}

	return nil
}

// isStrictAncestor returns true when child is nested inside parent.
func isStrictAncestor(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	if rel == "." || strings.HasPrefix(rel, "..") {
		return false
	}

	return !filepath.IsAbs(rel)
}
