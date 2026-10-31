// Package cmdsync provides the 3-Way Hash State Engine and human conflict protection.
package cmdsync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ThreeWayResult captures the evaluation decision for a single file.
type ThreeWayResult struct {
	RelPath          string
	State            string // InSync, CleanUpdate, LocalEdited, Converged, Collision, NewUpstream, UntrackedConflict
	SourceHash       string
	LocalHash        string
	BaselineHash     string
	IsHumanPreserved bool
	HasConflict      bool
}

// EvaluateThreeWay performs the 3-way cryptographic hash comparison.
func EvaluateThreeWay(relPath, srcPath, dstPath string, baselineMap map[string]string) ThreeWayResult {
	res := ThreeWayResult{RelPath: relPath}
	uHash, uErr := NormalizedFileHash(srcPath)
	if uErr != nil {
		res.State = "Error"
		return res
	}
	res.SourceHash = uHash

	bHash := ""
	if baselineMap != nil {
		bHash = baselineMap[filepath.ToSlash(relPath)]
	}
	res.BaselineHash = bHash

	dstInfo, dstErr := os.Stat(dstPath)
	if os.IsNotExist(dstErr) || dstInfo == nil {
		res.State = "NewUpstream"
		return res
	}

	lHash, lErr := NormalizedFileHash(dstPath)
	if lErr != nil {
		res.State = "Error"
		return res
	}
	res.LocalHash = lHash

	return resolveThreeWayState(res, bHash, uHash, lHash)
}

func resolveThreeWayState(res ThreeWayResult, b, u, l string) ThreeWayResult {
	if b == "" {
		return resolveUntrackedThreeWay(res, u, l)
	}
	if b == u && b == l {
		res.State = "InSync"
		return res
	}
	if b == l && b != u {
		res.State = "CleanUpdate"
		return res
	}
	if b == u && b != l {
		res.State = "LocalEdited"
		res.IsHumanPreserved = true
		return res
	}
	if u == l {
		res.State = "Converged"
		return res
	}

	res.State = "Collision"
	res.HasConflict = true
	res.IsHumanPreserved = true
	return res
}

func resolveUntrackedThreeWay(res ThreeWayResult, u, l string) ThreeWayResult {
	if u == l {
		res.State = "Converged"
		return res
	}

	res.State = "UntrackedConflict"
	res.HasConflict = true
	res.IsHumanPreserved = true
	return res
}

// ApplyThreeWayDecision applies the file update, sidecar creation, or preservation.
func ApplyThreeWayDecision(eval ThreeWayResult, srcPath, dstPath string, isDryRun bool) error {
	switch eval.State {
	case "InSync":
		return nil
	case "NewUpstream", "CleanUpdate":
		if isDryRun {
			return nil
		}
		return copyFileContents(srcPath, dstPath)
	case "LocalEdited":
		return nil
	case "Collision", "UntrackedConflict":
		if isDryRun {
			return nil
		}
		sidecarPath := dstPath + ".upstream"
		return copyFileContents(srcPath, sidecarPath)
	case "Converged":
		return nil
	default:
		return fmt.Errorf("unknown three-way state: %s", eval.State)
	}
}

func copyFileContents(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
