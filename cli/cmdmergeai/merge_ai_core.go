// Package cmdmergeai coordinates multi-repository staging, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdresolver"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/pterm/pterm"
)

func gitExec(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// RunMergeAI executes multi-repository amalgamation and AI merge staging.
func RunMergeAI(args []string) error {
	destTarget, resolvedSources, skippedSources, errParse := parseMergeAIArgs(args)
	if errParse != nil {
		return errParse
	}

	if len(skippedSources) > 0 {
		pterm.Warning.Printf("[Intelligence] Detected %d redundant/file source(s) to automatically skip:\n", len(skippedSources))
		for _, s := range skippedSources {
			if s.IsFile {
				pterm.Warning.Printf("  • Skipping file URL: %s\n    (Repo '%s' is already included; subpath: %s)\n", s.RawInput, s.Repo, s.FilePath)
			} else {
				pterm.Warning.Printf("  • Skipping duplicate repo: %s\n", s.Repo)
			}
		}
	}

	pterm.Info.Printf("[GitMap Merge-AI] Resolving destination target: %s\n", destTarget)
	resolvedDest, errDest := cmdresolver.ResolveRepoFeature(destTarget)
	if errDest != nil {
		return fmt.Errorf("failed to resolve repo feature destination: %w", errDest)
	}

	pterm.Success.Printf("[Repo Feature] Target resolved to: %s (Type: %s)\n",
		resolvedDest.LocalPath, resolvedDest.TargetType)

	stagingBase := filepath.Join(os.TempDir(), "gitmap-merge-ai-staging")
	_ = os.MkdirAll(stagingBase, 0755)

	pterm.Info.Println("[Inspection] Preparing and inspecting source repositories...")
	sources, errPrep := prepareAndInspectSources(resolvedSources, stagingBase)
	if errPrep != nil {
		return errPrep
	}

	if len(sources) == 0 {
		return fmt.Errorf("no valid source repositories found to merge")
	}

	// Sort chronologically (earliest to latest) to establish baseline precedence
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].LastCommitDate.Before(sources[j].LastCommitDate)
	})

	pterm.Info.Println("[Staging] Assembling staging tree and analyzing file collisions...")
	collisions, uniqueCount, errStage := stageAndSequenceFiles(resolvedDest.LocalPath, sources)
	if errStage != nil {
		return errStage
	}

	manifest := assembleManifest(resolvedDest, sources, collisions, uniqueCount, skippedSources)

	if errManifest := WriteMergeAIManifest(resolvedDest.LocalPath, manifest); errManifest != nil {
		return fmt.Errorf("failed to write merge-ai-manifest.json: %w", errManifest)
	}

	if errInstruction := WriteMergeAIInstruction(resolvedDest.LocalPath, collisions, skippedSources...); errInstruction != nil {
		return fmt.Errorf("failed to write instruction.md: %w", errInstruction)
	}

	renderMergeAISummary(resolvedDest, sources, collisions, uniqueCount, skippedSources)
	return nil
}

func parseMergeAIArgs(args []string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	if len(args) == 0 {
		return "", nil, nil, fmt.Errorf("usage: gitmap merge-ai <dest-target> <source-file.txt|source1> [source2...] or gitmap merge-ai <source-file.txt>")
	}

	// Case 1: single argument
	if len(args) == 1 {
		return parseSingleMergeAIArg(args[0])
	}

	// Case 2: two or more arguments
	return parseMultiMergeAIArgs(args)
}

func parseSingleMergeAIArg(arg string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	lower := strings.ToLower(arg)
	if strings.HasSuffix(lower, ".json") {
		return parseJSONConfigFile(arg)
	}

	destHint, resolved, skipped, err := ParseSourceFile(arg)
	if err != nil {
		return "", nil, nil, err
	}

	dest := destHint
	if dest == "" {
		dest = "icon-studio-v1"
		pterm.Info.Printf("[GitMap Merge-AI] No destination target specified; defaulting to: %s\n", dest)
	}

	return dest, resolved, skipped, nil
}

func parseJSONConfigFile(configPath string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg MergeAIConfigFile
	if errJSON := json.Unmarshal(data, &cfg); errJSON != nil {
		return "", nil, nil, fmt.Errorf("failed to parse JSON config %s: %w", configPath, errJSON)
	}

	destHint, resolved, skipped, errParse := ParseSourceManifestContent(strings.Join(cfg.Sources, "\n"))
	if errParse != nil {
		return "", nil, nil, errParse
	}

	dest := cfg.Destination
	if dest == "" {
		dest = destHint
	}
	if dest == "" {
		dest = "icon-studio-v1"
	}

	return dest, resolved, skipped, nil
}

func parseMultiMergeAIArgs(args []string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	// If first arg is a text file and subsequent args form destination name:
	if isTextFilePath(args[0]) {
		otherAreDest := true
		for _, a := range args[1:] {
			if isTextFilePath(a) || isFullURL(a) {
				otherAreDest = false
				break
			}
		}
		if otherAreDest {
			dest := strings.Join(args[1:], " ")
			_, resolved, skipped, err := ParseSourceFile(args[0])
			return dest, resolved, skipped, err
		}
	}

	// If last arg is a text file and preceding args form destination name:
	if isTextFilePath(args[len(args)-1]) {
		precedingAreDest := true
		for _, a := range args[:len(args)-1] {
			if isTextFilePath(a) || isFullURL(a) {
				precedingAreDest = false
				break
			}
		}
		if precedingAreDest {
			dest := strings.Join(args[:len(args)-1], " ")
			_, resolved, skipped, err := ParseSourceFile(args[len(args)-1])
			return dest, resolved, skipped, err
		}
	}

	dest := args[0]
	rawSources := args[1:]

	// If second arg is a text file: parse it
	if len(rawSources) == 1 && isTextFilePath(rawSources[0]) {
		_, resolved, skipped, err := ParseSourceFile(rawSources[0])
		return dest, resolved, skipped, err
	}

	// Combine lines and parse
	var combinedLines []string
	for _, s := range rawSources {
		if isTextFilePath(s) {
			_, resFile, _, errFile := ParseSourceFile(s)
			if errFile == nil {
				for _, r := range resFile {
					combinedLines = append(combinedLines, r.RepoURL)
				}
				continue
			}
		}

		if strings.Contains(s, ",") {
			for _, part := range strings.Split(s, ",") {
				if clean := strings.TrimSpace(part); clean != "" {
					combinedLines = append(combinedLines, clean)
				}
			}
			continue
		}

		if clean := strings.TrimSpace(s); clean != "" {
			combinedLines = append(combinedLines, clean)
		}
	}

	_, resolved, skipped, err := ParseSourceManifestContent(strings.Join(combinedLines, "\n"))
	return dest, resolved, skipped, err
}

func isTextFilePath(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	if strings.HasSuffix(lower, ".txt") || strings.HasSuffix(lower, ".manifest") {
		return true
	}
	if resolved, err := ResolveSourceFilePath(path); err == nil && resolved != "" {
		return true
	}
	return false
}

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

func stageAndSequenceFiles(destDir string, sources []SourceRepoConfig) ([]FileCollisionRecord, int, error) {
	collisionMap := make(map[string]*FileCollisionRecord)
	uniquePlaced := 0

	for seqIdx, src := range sources {
		seqPrefix := fmt.Sprintf("%02d", seqIdx+1)
		sourceSlug := filepath.Base(src.LocalStaging)
		headSha, _ := gitExec(src.LocalStaging, "rev-parse", "HEAD")
		if len(headSha) > 8 {
			headSha = headSha[:8]
		}

		for _, relPath := range src.FolderTree {
			destPath := filepath.Join(destDir, relPath)
			dirName := filepath.Dir(relPath)
			baseName := filepath.Base(relPath)

			if seqIdx == 0 {
				// Base repository files are placed directly
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), destPath); errCopy != nil {
					return nil, 0, errCopy
				}
				uniquePlaced++
				continue
			}

			// Check if this file is already recorded as a known collision
			if col, isKnownCollision := collisionMap[relPath]; isKnownCollision {
				newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
				newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), newSeqDest); errCopy != nil {
					return nil, 0, errCopy
				}

				col.Variants = append(col.Variants, FileCollisionVariant{
					Sequence:   seqPrefix,
					File:       newSeqRel,
					SourceRepo: sourceSlug,
					Commit:     headSha,
				})
				continue
			}

			// If it does not exist on disk, place directly
			if !fileExists(destPath) {
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), destPath); errCopy != nil {
					return nil, 0, errCopy
				}
				uniquePlaced++
				continue
			}

			// First collision detected on relPath:
			// Sequence the pre-existing file as 01_<filename>
			origSeqPath := filepath.Join(destDir, dirName, fmt.Sprintf("01_%s", baseName))
			_ = os.Rename(destPath, origSeqPath)
			if uniquePlaced > 0 {
				uniquePlaced--
			}

			baseHeadSha, _ := gitExec(sources[0].LocalStaging, "rev-parse", "HEAD")
			if len(baseHeadSha) > 8 {
				baseHeadSha = baseHeadSha[:8]
			}

			col := &FileCollisionRecord{
				CanonicalPath: relPath,
				Status:        "pending_ai_consolidation",
				Variants: []FileCollisionVariant{
					{
						Sequence:   "01",
						File:       filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("01_%s", baseName))),
						SourceRepo: filepath.Base(sources[0].LocalStaging),
						Commit:     baseHeadSha,
					},
				},
			}
			collisionMap[relPath] = col

			// Now place the incoming colliding version as seqPrefix_<filename>
			newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
			newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
			if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), newSeqDest); errCopy != nil {
				return nil, 0, errCopy
			}

			col.Variants = append(col.Variants, FileCollisionVariant{
				Sequence:   seqPrefix,
				File:       newSeqRel,
				SourceRepo: sourceSlug,
				Commit:     headSha,
			})
		}
	}

	var collisionList []FileCollisionRecord
	for _, c := range collisionMap {
		collisionList = append(collisionList, *c)
	}

	sort.Slice(collisionList, func(i, j int) bool {
		return collisionList[i].CanonicalPath < collisionList[j].CanonicalPath
	})

	return collisionList, uniquePlaced, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func copyFile(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func assembleManifest(dest *cmdresolver.ResolvedRepoFeature, sources []SourceRepoConfig, collisions []FileCollisionRecord, uniqueCount int, skipped []SkippedSourceItem) MergeAIManifest {
	var m MergeAIManifest
	m.Attributes.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	m.Attributes.GitMapVersion = constants.Version
	m.Attributes.Tool = "gitmap-merge-ai"
	m.Attributes.CollisionStrategy = "prefixed-sequence"

	m.Data.Destination.TargetType = string(dest.TargetType)
	m.Data.Destination.Slug = dest.CanonicalSlug
	m.Data.Destination.ResolvedPath = dest.LocalPath
	m.Data.Destination.OriginURL = dest.RemoteURL

	for i, s := range sources {
		var entry SourceRepoSequenceEntry
		entry.SequenceOrder = i + 1
		entry.Title = s.Title
		entry.RepoURL = s.RepoURL
		entry.Branch = s.Branch
		entry.CommitRange = s.CommitRange
		entry.LastCommitDate = s.LastCommitDate.Format(time.RFC3339)
		entry.Description = s.Description
		entry.ReleaseInfo.LastTag = s.ReleaseTag
		entry.FolderTree = s.FolderTree
		m.Data.RepoSequence = append(m.Data.RepoSequence, entry)
	}

	m.Data.SkippedSources = skipped
	m.Data.FileCollisions = collisions
	m.Data.UniqueFilesDirectlyPlaced = uniqueCount
	return m
}

func renderMergeAISummary(dest *cmdresolver.ResolvedRepoFeature, sources []SourceRepoConfig, collisions []FileCollisionRecord, uniqueCount int, skipped []SkippedSourceItem) {
	fmt.Println()
	pterm.Success.Printf("[GitMap Merge-AI] Multi-repo amalgamation staged successfully!\n")
	fmt.Printf("Destination Workspace: %s\n", pterm.Bold.Sprint(dest.LocalPath))
	fmt.Printf("Source Repositories Merged: %d\n", len(sources))
	for i, s := range sources {
		name := s.Title
		if name == "" {
			name = deriveSourceSlug(s.RepoURL, i)
		}
		fmt.Printf("  [%d] %s (%s, %d files)\n", i+1, pterm.Cyan(name), s.Branch, len(s.FolderTree))
	}
	if len(skipped) > 0 {
		fmt.Printf("\nSources Automatically Skipped: %d\n", len(skipped))
		for _, sk := range skipped {
			fmt.Printf("  - %s (%s)\n", sk.RawInput, pterm.Yellow(sk.Reason))
		}
	}
	fmt.Printf("\nUnique Files Directly Placed: %d\n", uniqueCount)
	fmt.Printf("File Collisions Sequenced: %d\n", len(collisions))

	if len(collisions) > 0 {
		fmt.Println("\nColliding Files Sequenced for AI Consolidation:")
		for _, c := range collisions {
			fmt.Printf("  ! %s -> %d variants (%s)\n", pterm.Yellow(c.CanonicalPath), len(c.Variants), c.Status)
		}
	}

	fmt.Println("\nArtifacts Created:")
	fmt.Println("  ✓ merge-ai-manifest.json (Detailed metadata & sequence mapping)")
	fmt.Println("  ✓ instruction.md (Step-by-step AI checklist for code consolidation)")
	fmt.Println("\nNext Step:")
	fmt.Printf("  An AI coding agent can now open '%s', read 'instruction.md',\n", dest.LocalPath)
	fmt.Println("  and consolidate the sequenced collision files into a unified codebase.")
}
