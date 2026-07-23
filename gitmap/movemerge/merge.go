package movemerge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RunMerge executes merge-both / merge-left / merge-right.
func RunMerge(left, right Endpoint, dir Direction, opts Options) error {
	if err := GuardEndpoints(left, right); err != nil {
		return err
	}
	logf(opts.LogPrefix, "diffing trees ...")
	entries, err := DiffTrees(left.WorkingDir, right.WorkingDir, opts)
	if err != nil {
		return err
	}
	resolver := NewResolver(effectivePolicy(dir, opts), os.Stdin, os.Stdout)
	for _, e := range entries {
		if applyErr := applyEntry(e, left, right, dir, resolver, opts); applyErr != nil {
			return applyErr
		}
	}
	if finErr := finalizeURLSides(left, right, dir, opts); finErr != nil {
		return finErr
	}
	logf(opts.LogPrefix, "done")

	return nil
}

// effectivePolicy returns the bypass policy when -y is set.
func effectivePolicy(dir Direction, opts Options) PreferPolicy {
	if !opts.Yes {
		return PreferNone
	}
	if opts.Prefer != PreferNone {
		return opts.Prefer
	}
	switch dir {
	case DirBoth:
		return PreferNewer
	case DirRightOnly:
		return PreferLeft
	case DirLeftOnly:
		return PreferRight
	}

	return PreferNewer
}

// applyEntry handles one DiffEntry per the requested direction.
func applyEntry(e DiffEntry, l, r Endpoint, dir Direction, res *Resolver, opts Options) error {
	switch e.Kind {
	case DiffIdentical:
		return nil
	case DiffMissingLeft:
		return applyMissing(e, l, r, dir, opts, false)
	case DiffMissingRight:
		return applyMissing(e, l, r, dir, opts, true)
	case DiffConflict:
		return applyConflict(e, l, r, dir, res, opts)
	}

	return applyConflict(e, l, r, dir, res, opts)
}

// applyMissing copies a file present on only one side to the other.
// fromLeft=true means LEFT has it; copy to RIGHT (when allowed).
func applyMissing(e DiffEntry, l, r Endpoint, dir Direction, opts Options, fromLeft bool) error {
	if fromLeft && (dir == DirBoth || dir == DirRightOnly) {
		return copyOne(l.WorkingDir, r.WorkingDir, e.RelPath, e.Left.Info, opts)
	}
	if !fromLeft && (dir == DirBoth || dir == DirLeftOnly) {
		return copyOne(r.WorkingDir, l.WorkingDir, e.RelPath, e.Right.Info, opts)
	}

	return nil
}

// applyConflict resolves and applies one conflicting path.
func applyConflict(e DiffEntry, l, r Endpoint, dir Direction, res *Resolver, opts Options) error {
	choice, err := res.Resolve(e.RelPath, e.Left, e.Right)
	if err != nil {
		return err
	}
	if choice == ChoiceQuit {
		return fmt.Errorf("%s", constants.ErrMMQuit)
	}
	if choice == ChoiceSkip {
		logIndent(opts.LogPrefix, "conflict %s -> skipped", e.RelPath)

		return nil
	}

	return writeChoice(choice, e, l, r, dir, opts)
}

// writeChoice writes the chosen side onto the destination(s).
func writeChoice(c Choice, e DiffEntry, l, r Endpoint, dir Direction, opts Options) error {
	if c == ChoiceLeft && (dir == DirBoth || dir == DirRightOnly) {
		logIndent(opts.LogPrefix, "conflict %s -> took LEFT", e.RelPath)

		return copyOne(l.WorkingDir, r.WorkingDir, e.RelPath, e.Left.Info, opts)
	}
	if c == ChoiceRight && (dir == DirBoth || dir == DirLeftOnly) {
		logIndent(opts.LogPrefix, "conflict %s -> took RIGHT", e.RelPath)

		return copyOne(r.WorkingDir, l.WorkingDir, e.RelPath, e.Right.Info, opts)
	}
	logIndent(opts.LogPrefix, "conflict %s -> no-op (direction)", e.RelPath)

	return nil
}

// copyOne copies a single relative path between working dirs.
func copyOne(srcDir, dstDir, rel string, info os.FileInfo, opts Options) error {
	if opts.DryRun {
		logIndent(opts.LogPrefix, "[dry-run] copy %s", rel)

		return nil
	}
	src := filepath.Join(srcDir, filepath.FromSlash(rel))
	dst := filepath.Join(dstDir, filepath.FromSlash(rel))

	return CopyFile(src, dst, info)
}

// silenceUnused is here only to retain io import for future hooks.
var _ = io.Discard
