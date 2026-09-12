// Package cmdpurge — types.go defines domain payload models and single reusable Result envelopes.
package cmdpurge

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type (
	// TrackedLovableFilesMapResult is the canonical single reusable result envelope for tracked lovable files.
	TrackedLovableFilesMapResult = result.ResultMap[string, bool]
)
