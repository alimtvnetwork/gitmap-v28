// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/pterm/pterm"
)

func prepareAndInspectSources(items []ResolvedSourceItem, stagingBase string) ([]SourceRepoConfig, error) {
	var sources []SourceRepoConfig

	for idx, item := range items {
		slug := deriveSourceSlug(item.RepoURL, idx)
		stagingDir := filepath.Join(stagingBase, slug)

		// Check if source repo already exists in local workspace (e.g. D:\work\<slug>)
		localWorkDir := filepath.Join(resolveWorkBaseDir(), slug)
		if info, errStat := os.Stat(filepath.Join(localWorkDir, ".git")); errStat == nil && info.IsDir() {
			stagingDir = localWorkDir
			pterm.Info.Printf("  ▸ [%d/%d] Adopting existing local repository: %s\n", idx+1, len(items), localWorkDir)
		} else if info, statStaging := os.Stat(filepath.Join(stagingDir, ".git")); statStaging == nil && info.IsDir() {
			pterm.Info.Printf("  ▸ [%d/%d] Reusing cached staging clone: %s\n", idx+1, len(items), stagingDir)
		} else {
			_ = os.RemoveAll(stagingDir)
			pterm.Info.Printf("  ▸ [%d/%d] Cloning source repository: %s -> %s\n", idx+1, len(items), item.RepoURL, stagingDir)
			cmd := exec.Command("git", "clone", "--depth", "50", item.RepoURL, stagingDir)
			if out, errClone := cmd.CombinedOutput(); errClone != nil {
				return nil, fmt.Errorf("failed to clone source %s (%s): %s (%w)", item.Title, item.RepoURL, string(out), errClone)
			}
		}

		cfg := inspectSingleSource(item.RepoURL, stagingDir)
		cfg.Title = item.Title
		sources = append(sources, cfg)
	}

	return sources, nil
}

func resolveWorkBaseDir() string {
	if defDB, err := store.OpenDefault(); err == nil {
		defer defDB.Close()
		if wd, err := defDB.GetDefaultWorkDir(); err == nil && wd != nil && len(wd.AbsolutePath) > 0 {
			if info, errStat := os.Stat(wd.AbsolutePath); errStat == nil && info.IsDir() {
				return wd.AbsolutePath
			}
		}
	}

	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return `D:\work`
	}

	cwd, err := os.Getwd()
	if err == nil && cwd != "" {
		parent := filepath.Dir(cwd)
		if parent != "" && parent != cwd {
			return parent
		}
		return cwd
	}

	return "."
}

func deriveSourceSlug(input string, idx int) string {
	clean := strings.TrimSuffix(input, ".git")
	clean = strings.TrimRight(clean, "/")
	parts := strings.Split(clean, "/")
	if len(parts) > 0 && parts[len(parts)-1] != "" {
		return parts[len(parts)-1]
	}
	return fmt.Sprintf("source-%d", idx+1)
}

func inspectSingleSource(url, dir string) SourceRepoConfig {
	branchOut, _ := gitExec(dir, "rev-parse", "--abbrev-ref", "HEAD")
	epochOut, _ := gitExec(dir, "log", "-1", "--format=%ct")
	commitMsg, _ := gitExec(dir, "log", "-1", "--format=%s")
	tagOut, _ := gitExec(dir, "describe", "--tags", "--abbrev=0")
	headSha, _ := gitExec(dir, "rev-parse", "HEAD")
	rootSha, _ := gitExec(dir, "rev-list", "--max-parents=0", "HEAD")

	var lastDate time.Time
	if epochSec, err := strconv.ParseInt(epochOut, 10, 64); err == nil && epochSec > 0 {
		lastDate = time.Unix(epochSec, 0)
	} else {
		lastDate = time.Now()
	}

	commitRange := headSha
	if rootSha != "" && rootSha != headSha {
		shortRoot := rootSha
		if len(shortRoot) > 7 {
			shortRoot = shortRoot[:7]
		}
		shortHead := headSha
		if len(shortHead) > 7 {
			shortHead = shortHead[:7]
		}
		commitRange = fmt.Sprintf("%s..%s", shortRoot, shortHead)
	}

	tree := enumerateTrackedFiles(dir)

	return SourceRepoConfig{
		RepoURL:        url,
		Branch:         branchOut,
		CommitRange:    commitRange,
		LastCommitDate: lastDate,
		Description:    commitMsg,
		ReleaseTag:     tagOut,
		LocalStaging:   dir,
		FolderTree:     tree,
	}
}

func enumerateTrackedFiles(dir string) []string {
	out, err := gitExec(dir, "ls-files")
	if err != nil || out == "" {
		return nil
	}
	var files []string
	for _, f := range strings.Split(out, "\n") {
		clean := strings.TrimSpace(f)
		if clean != "" && !strings.HasPrefix(clean, ".git") {
			files = append(files, clean)
		}
	}
	return files
}
