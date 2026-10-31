// Package cmdsync provides the 3-Way Hash State Engine and human conflict protection.
package cmdsync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	srcHashCacheMu sync.RWMutex
	srcHashCache   = make(map[string]string)
)

// ClearSourceHashCache resets the in-memory source hash cache.
func ClearSourceHashCache() {
	srcHashCacheMu.Lock()
	srcHashCache = make(map[string]string)
	srcHashCacheMu.Unlock()
}

func getCachedSourceHash(srcPath string) (string, error) {
	srcHashCacheMu.RLock()
	h, ok := srcHashCache[srcPath]
	srcHashCacheMu.RUnlock()
	if ok {
		return h, nil
	}

	h, err := NormalizedFileHash(srcPath)
	if err != nil {
		return "", err
	}

	srcHashCacheMu.Lock()
	srcHashCache[srcPath] = h
	srcHashCacheMu.Unlock()

	return h, nil
}

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

// fastPathThreeWayHash returns the cached source hash when the baseline,
// size, and ModTime fast-path conditions all hold and the hash matches
// the baseline; otherwise it reports not-in-sync. Single-level guard
// only, so callers stay flat for the nested-if linter.
func fastPathThreeWayHash(bHash string, srcInfo, dstInfo os.FileInfo, srcPath string) (string, bool) {
	if bHash == "" || srcInfo.Size() != dstInfo.Size() || !srcInfo.ModTime().Equal(dstInfo.ModTime()) {
		return "", false
	}

	fastHash, _ := getCachedSourceHash(srcPath)

	return fastHash, fastHash == bHash
}

// EvaluateThreeWay performs the 3-way cryptographic hash comparison with ModTime/Size fast-path.
func EvaluateThreeWay(relPath, srcPath, dstPath string, baselineMap map[string]string) ThreeWayResult {
	res := ThreeWayResult{RelPath: relPath}
	bHash := ""
	if baselineMap != nil {
		bHash = baselineMap[filepath.ToSlash(relPath)]
	}
	res.BaselineHash = bHash

	srcInfo, srcErr := os.Stat(srcPath)
	if srcErr != nil {
		res.State = "Error"
		return res
	}

	dstInfo, dstErr := os.Stat(dstPath)
	if os.IsNotExist(dstErr) || dstInfo == nil {
		uHash, _ := getCachedSourceHash(srcPath)
		res.SourceHash = uHash
		res.State = "NewUpstream"
		return res
	}

	// FAST PATH: If baseline is known, size and ModTime match exactly,
	// file is verified InSync without reading destination bytes.
	// WHY helper: the inner bHash==uHash check inside the fast-path
	// guard tripped the nested-if linter on the v6.644.0 tip.
	if fastHash, isFastInSync := fastPathThreeWayHash(bHash, srcInfo, dstInfo, srcPath); isFastInSync {
		res.SourceHash = fastHash
		res.LocalHash = bHash
		res.State = "InSync"
		return res
	}

	uHash, uErr := getCachedSourceHash(srcPath)
	if uErr != nil {
		res.State = "Error"
		return res
	}
	res.SourceHash = uHash

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

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if srcStat, statErr := os.Stat(src); statErr == nil {
		_ = os.Chtimes(dst, srcStat.ModTime(), srcStat.ModTime())
	}

	return nil
}
