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

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/pterm/pterm"
)

func prepareAndInspectSources(items []ResolvedSourceItem, stagingBase string) ([]SourceRepoConfig, error) {
	var sources []SourceRepoConfig

	for idx, item := range items {
		stagingDir, errStaging := resolveSourceStagingDir(item, idx, len(items), stagingBase)
		if errStaging != nil {
			return nil, errStaging
		}

		cfg := inspectSingleSource(item.RepoURL, stagingDir)
		cfg.Title = item.Title
		sources = append(sources, cfg)
	}

	return sources, nil
}

func resolveSourceStagingDir(item ResolvedSourceItem, idx, total int, stagingBase string) (string, error) {
	slug := deriveSourceSlug(item.RepoURL, idx)
	localWorkDir := filepath.Join(resolveWorkBaseDir(), slug)
	if fsutil.DirExists(filepath.Join(localWorkDir, ".git")) {
		pterm.Info.Printf("  ▸ [%d/%d] Adopting existing local repository: %s\n", idx+1, total, localWorkDir)
		return localWorkDir, nil
	}

	stagingDir := filepath.Join(stagingBase, slug)
	if fsutil.DirExists(filepath.Join(stagingDir, ".git")) {
		pterm.Info.Printf("  ▸ [%d/%d] Reusing cached staging clone: %s\n", idx+1, total, stagingDir)
		return stagingDir, nil
	}

	_ = os.RemoveAll(stagingDir)
	pterm.Info.Printf("  ▸ [%d/%d] Cloning source repository: %s -> %s\n", idx+1, total, item.RepoURL, stagingDir)
	cmd := exec.Command("git", "clone", "--depth", "50", item.RepoURL, stagingDir)
	out, errClone := cmd.CombinedOutput()
	if errClone != nil {
		return "", apperror.WrapSimple(errClone, fmt.Sprintf("failed to clone source %s (%s)", item.Title, item.RepoURL)).WithContext("output", string(out))
	}
	return stagingDir, nil
}

func resolveStoreWorkDir() string {
	defDB, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer defDB.Close()

	wd, errWd := defDB.GetDefaultWorkDir()
	if errWd != nil || wd == nil || len(wd.AbsolutePath) == 0 {
		return ""
	}
	if fsutil.DirExists(wd.AbsolutePath) {
		return wd.AbsolutePath
	}
	return ""
}

func resolveCwdFallback() string {
	cwd, err := os.Getwd()
	if err != nil || cwd == "" {
		return "."
	}
	parent := filepath.Dir(cwd)
	if parent != "" && parent != cwd {
		return parent
	}
	return cwd
}

func resolveWorkBaseDir() string {
	storeDir := resolveStoreWorkDir()
	if storeDir != "" {
		return storeDir
	}

	if fsutil.DirExists(`D:\work`) {
		return `D:\work`
	}

	return resolveCwdFallback()
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

func formatCommitRange(headSha, rootSha string) string {
	if rootSha == "" || rootSha == headSha {
		return headSha
	}
	shortRoot := truncateCommitSha(rootSha)
	shortHead := truncateCommitSha(headSha)
	return fmt.Sprintf("%s..%s", shortRoot, shortHead)
}

func inspectCommitMetadata(dir, branch, headSha, commitMsg string, lastDate time.Time) CommitMetadata {
	authorName, _ := gitExec(dir, "log", "-1", "--format=%an")
	authorEmail, _ := gitExec(dir, "log", "-1", "--format=%ae")
	dateISO, _ := gitExec(dir, "log", "-1", "--format=%cI")

	author := authorName
	if authorEmail != "" {
		author = fmt.Sprintf("%s <%s>", authorName, authorEmail)
	}
	if author == "" {
		author = "unknown"
	}

	if dateISO == "" {
		dateISO = lastDate.UTC().Format(time.RFC3339)
	}

	return CommitMetadata{
		SHA:         headSha,
		Author:      author,
		AuthorEmail: authorEmail,
		Date:        dateISO,
		Message:     commitMsg,
		Branch:      branch,
		ShortSHA:    truncateCommitSha(headSha),
	}
}

func inspectSingleSource(url, dir string) SourceRepoConfig {
	branchOut, _ := gitExec(dir, "rev-parse", "--abbrev-ref", "HEAD")
	epochOut, _ := gitExec(dir, "log", "-1", "--format=%ct")
	commitMsg, _ := gitExec(dir, "log", "-1", "--format=%s")
	tagOut, _ := gitExec(dir, "describe", "--tags", "--abbrev=0")
	headSha, _ := gitExec(dir, "rev-parse", "HEAD")
	rootSha, _ := gitExec(dir, "rev-list", "--max-parents=0", "HEAD")

	tagHash := ""
	if tagOut != "" {
		tagHash, _ = gitExec(dir, "rev-parse", tagOut)
	}

	var lastDate time.Time
	if epochSec, err := strconv.ParseInt(epochOut, 10, 64); err == nil && epochSec > 0 {
		lastDate = time.Unix(epochSec, 0)
	} else {
		lastDate = time.Now()
	}

	commitRange := formatCommitRange(headSha, rootSha)
	tree := enumerateTrackedFiles(dir)
	commitMeta := inspectCommitMetadata(dir, branchOut, headSha, commitMsg, lastDate)

	return SourceRepoConfig{
		RepoURL:        url,
		Branch:         branchOut,
		CommitRange:    commitRange,
		LastCommitDate: lastDate,
		Description:    commitMsg,
		ReleaseTag:     tagOut,
		ReleaseTagHash: tagHash,
		LocalStaging:   dir,
		FolderTree:     tree,
		CommitMetadata: commitMeta,
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
