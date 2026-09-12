package movemerge

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"path/filepath"
)

// DiffTrees walks both sides and classifies every relative path.
// Identical files are detected by SHA-256 (computed on demand).
func DiffTrees(leftDir, rightDir string, opts Options) ([]DiffEntry, error) {
	liRes := IndexTree(leftDir, opts)
	if liRes.IsFailure() {
		return nil, liRes.AppError()
	}

	riRes := IndexTree(rightDir, opts)
	if riRes.IsFailure() {
		return nil, riRes.AppError()
	}

	keys := SortedKeys(liRes.Data, riRes.Data)
	out := make([]DiffEntry, 0, len(keys))
	for _, rel := range keys {
		res := classifyOne(rel, liRes.Data, riRes.Data, leftDir, rightDir)
		if res.IsFailure() {
			return nil, res.Err
		}

		entry := res.Value
		out = append(out, entry)
	}

	return out, nil
}

// classifyOne returns the DiffEntry for a single relative path.
func classifyOne(
	rel string,
	li,
	ri map[string]FileMeta,
	leftDir,
	rightDir string,
) DiffEntryResult {
	l, lOK := li[rel]
	r, rOK := ri[rel]
	entry := DiffEntry{RelPath: rel, Left: l, Right: r}
	if lOK && !rOK {
		entry.Kind = DiffMissingRight

		return result.SuccessResult(entry)
	}

	if !lOK && rOK {
		entry.Kind = DiffMissingLeft

		return result.SuccessResult(entry)
	}

	return classifyBoth(entry, leftDir, rightDir)
}

// classifyBoth resolves Identical vs Conflict via SHA-256.
func classifyBoth(entry DiffEntry, leftDir, rightDir string) DiffEntryResult {
	lPath := filepath.Join(leftDir, filepath.FromSlash(entry.RelPath))
	rPath := filepath.Join(rightDir, filepath.FromSlash(entry.RelPath))
	lh, err := HashFile(lPath)
	if err != nil {
		return result.FailureResult[DiffEntry](apperror.WrapSimple(err, "movemerge"))
	}

	rh, err := HashFile(rPath)
	if err != nil {
		return result.FailureResult[DiffEntry](apperror.WrapSimple(err, "movemerge"))
	}

	if lh == rh {
		entry.Kind = DiffIdentical

		return result.SuccessResult(entry)
	}

	entry.Kind = DiffConflict

	return result.SuccessResult(entry)
}
