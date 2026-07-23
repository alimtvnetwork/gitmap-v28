package movemerge

import "path/filepath"

// DiffKind classifies a path across LEFT and RIGHT.
type DiffKind int

const (
	// DiffMissingLeft = present on RIGHT only.
	DiffMissingLeft DiffKind = iota
	// DiffMissingRight = present on LEFT only.
	DiffMissingRight
	// DiffConflict = present on both with different content.
	DiffConflict
	// DiffIdentical = present on both with byte-equal content.
	DiffIdentical
)

// DiffEntry is one classified path with both sides' metadata.
type DiffEntry struct {
	RelPath string
	Kind    DiffKind
	Left    FileMeta
	Right   FileMeta
}

// DiffTrees walks both sides and classifies every relative path.
// Identical files are detected by SHA-256 (computed on demand).
func DiffTrees(leftDir, rightDir string, opts Options) ([]DiffEntry, error) {
	li, err := IndexTree(leftDir, opts)
	if err != nil {
		return nil, err
	}
	ri, err := IndexTree(rightDir, opts)
	if err != nil {
		return nil, err
	}
	keys := SortedKeys(li, ri)
	out := make([]DiffEntry, 0, len(keys))
	for _, rel := range keys {
		entry, derr := classifyOne(rel, li, ri, leftDir, rightDir)
		if derr != nil {
			return nil, derr
		}
		out = append(out, entry)
	}

	return out, nil
}

// classifyOne returns the DiffEntry for a single relative path.
func classifyOne(rel string, li, ri map[string]FileMeta, leftDir, rightDir string) (DiffEntry, error) {
	l, lOK := li[rel]
	r, rOK := ri[rel]
	entry := DiffEntry{RelPath: rel, Left: l, Right: r}
	if lOK && !rOK {
		entry.Kind = DiffMissingRight

		return entry, nil
	}
	if !lOK && rOK {
		entry.Kind = DiffMissingLeft

		return entry, nil
	}

	return classifyBoth(entry, leftDir, rightDir)
}

// classifyBoth resolves Identical vs Conflict via SHA-256.
func classifyBoth(entry DiffEntry, leftDir, rightDir string) (DiffEntry, error) {
	lPath := filepath.Join(leftDir, filepath.FromSlash(entry.RelPath))
	rPath := filepath.Join(rightDir, filepath.FromSlash(entry.RelPath))
	lh, err := HashFile(lPath)
	if err != nil {
		return entry, err
	}
	rh, err := HashFile(rPath)
	if err != nil {
		return entry, err
	}
	if lh == rh {
		entry.Kind = DiffIdentical

		return entry, nil
	}
	entry.Kind = DiffConflict

	return entry, nil
}
