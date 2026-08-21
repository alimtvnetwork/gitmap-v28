package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
)

func calculateDestPath(srcPath, destTarget string) (string, error) {
	cleanSrc := filepath.Clean(filepath.FromSlash(fsutil.TrimTrailingSlashes(srcPath)))
	repoName := filepath.Base(cleanSrc)
	t := fsutil.TrimTrailingSlashes(destTarget)

	if t == ".." {
		parent := filepath.Dir(cleanSrc)
		grandParent := filepath.Dir(parent)
		return filepath.Join(grandParent, repoName), nil
	}

	absDest, err := filepath.Abs(t)
	if err != nil {
		return "", err
	}
	if fsutil.DirExists(absDest) {
		return filepath.Join(absDest, repoName), nil
	}
	return absDest, nil
}

func preflightMove(srcPath, destPath string) error {
	if !fsutil.DirExists(srcPath) {
		return fmt.Errorf("source directory %q does not exist", srcPath)
	}
	if fsutil.EqualPaths(srcPath, destPath) {
		return fmt.Errorf("source and destination are identical")
	}
	if fsutil.FileOrDirExists(destPath) {
		return fmt.Errorf("destination %q already exists", destPath)
	}
	if fsutil.IsSubdirectory(srcPath, destPath) {
		return fmt.Errorf("cannot move repository into its own subdirectory")
	}
	return nil
}
