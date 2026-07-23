package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// cloneReplaceResult describes how the replace flow finished.
type cloneReplaceResult struct {
	Strategy string // "direct", "temp-swap", or "" when target was empty.
	Note     string
}

// cloneReplacing implements spec/01-app/96-clone-replace-existing-folder.md.
// It clones url into target, replacing any pre-existing folder via two
// strategies: (1) direct remove + clone, (2) temp-clone then swap-in-place.
func cloneReplacing(url, target string) (cloneReplaceResult, error) {
	res := cloneReplaceResult{}

	// Spec 113: if the user's shell cwd is inside the target folder,
	// chdir to the parent so the upcoming os.RemoveAll can succeed
	// on Windows. No-op when cwd is elsewhere.
	if _, escErr := escapeCwdIfInside(target); escErr != nil {
		return res, escErr
	}

	if _, statErr := os.Stat(target); errors.Is(statErr, fs.ErrNotExist) {
		fmt.Printf(constants.MsgCloneReplaceFree, target)

		if err := runCloneCommand(url, target); err != nil {
			return res, err
		}

		res.Strategy = "direct"

		return res, nil
	}

	fmt.Printf(constants.MsgCloneReplaceExists, target)
	fmt.Println(constants.MsgCloneReplaceStrategy1)

	if removeErr := os.RemoveAll(target); removeErr == nil {
		if err := runCloneCommand(url, target); err != nil {
			return res, err
		}

		res.Strategy = "direct"

		return res, nil
	} else {
		fmt.Printf(constants.MsgCloneReplaceStrat1Fail, removeErr)
	}

	return cloneViaTempSwap(url, target)
}

// cloneViaTempSwap implements strategy 2: clone into a temp sibling folder,
// empty the target's contents, then move every entry across.
func cloneViaTempSwap(url, target string) (cloneReplaceResult, error) {
	res := cloneReplaceResult{}

	fmt.Println(constants.MsgCloneReplaceStrategy2)

	tmp := target + ".gitmap-tmp-" + randSuffix()
	_ = os.RemoveAll(tmp)

	fmt.Printf(constants.MsgCloneReplaceTempClone, tmp)

	if err := runCloneCommand(url, tmp); err != nil {
		return res, fmt.Errorf("git clone into temp failed: %w", err)
	}
	defer os.RemoveAll(tmp)

	if err := emptyDirContents(target); err != nil {
		return res, fmt.Errorf("could not empty target: %w", err)
	}

	if err := moveDirContents(tmp, target); err != nil {
		return res, fmt.Errorf("swap failed: %w", err)
	}

	fmt.Println(constants.MsgCloneReplaceSwapDone)

	res.Strategy = "temp-swap"
	res.Note = "replaced via temp-swap"

	return res, nil
}

// runCloneCommand executes git clone with stdio inherited. Delegates
// to runCloneCommandPretty so every clone path (clone, cfr, cfrp,
// temp-swap fallback) shares the same colorful header / spinner /
// timing / failure-panel formatting and honors --dry-run.
func runCloneCommand(url, dest string) error {
	return runCloneCommandPretty(url, dest)
}

// emptyDirContents removes every entry inside dir, leaving dir itself in place.
// This survives a directory handle held by the caller's shell on Windows.
func emptyDirContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	fmt.Printf(constants.MsgCloneReplaceEmptying, len(entries), dir)

	var failures []string

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if rmErr := os.RemoveAll(path); rmErr != nil {
			fmt.Fprintf(os.Stderr, constants.WarnCloneReplaceEntryFail, entry.Name(), rmErr)

			failures = append(failures, entry.Name())
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("%d entries could not be removed (e.g. %s)", len(failures), failures[0])
	}

	return nil
}

// moveDirContents renames every child of src into dst. dst must exist.
func moveDirContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read src: %w", err)
	}

	fmt.Printf(constants.MsgCloneReplaceMoving, len(entries))

	for _, entry := range entries {
		from := filepath.Join(src, entry.Name())
		to := filepath.Join(dst, entry.Name())

		if mvErr := os.Rename(from, to); mvErr != nil {
			return fmt.Errorf("rename %s -> %s: %w", from, to, mvErr)
		}
	}

	return nil
}

// randSuffix returns 8 hex chars suitable for a temp folder name.
func randSuffix() string {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return "fallback"
	}

	return hex.EncodeToString(buf)
}
