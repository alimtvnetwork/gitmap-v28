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
	destTarget, sourceInputs, errParse := parseMergeAIArgs(args)
	if errParse != nil {
		return errParse
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
	sources, errPrep := prepareAndInspectSources(sourceInputs, stagingBase)
	if errPrep != nil {
		return errPrep
	}

	if len(sources) == 0 {
		return fmt.Errorf("no valid source repositories found to merge")
	}

	// Sort chronologically (earliest to latest)
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].LastCommitDate.Before(sources[j].LastCommitDate)
	})

	pterm.Info.Println("[Staging] Assembling staging tree and analyzing file collisions...")
	collisions, uniqueCount, errStage := stageAndSequenceFiles(resolvedDest.LocalPath, sources)
	if errStage != nil {
		return errStage
	}

	manifest := assembleManifest(resolvedDest, sources, collisions, uniqueCount)

	if errManifest := WriteMergeAIManifest(resolvedDest.LocalPath, manifest); errManifest != nil {
		return fmt.Errorf("failed to write merge-ai-manifest.json: %w", errManifest)
	}

	if errInstruction := WriteMergeAIInstruction(resolvedDest.LocalPath, collisions); errInstruction != nil {
		return fmt.Errorf("failed to write instruction.md: %w", errInstruction)
	}

	renderMergeAISummary(resolvedDest, sources, collisions, uniqueCount)
	return nil
}

func parseMergeAIArgs(args []string) (string, []string, error) {
	if len(args) == 0 {
		return "", nil, fmt.Errorf("usage: gitmap merge-ai <dest-target> <source1> [source2...] or gitmap merge-ai config.json")
	}

	if len(args) == 1 && strings.HasSuffix(strings.ToLower(args[0]), ".json") {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return "", nil, fmt.Errorf("failed to read config file %s: %w", args[0], err)
		}
		var cfg MergeAIConfigFile
		if errJSON := json.Unmarshal(data, &cfg); errJSON != nil {
			return "", nil, fmt.Errorf("failed to parse JSON config %s: %w", args[0], errJSON)
		}
		return cfg.Destination, cfg.Sources, nil
	}

	dest := args[0]
	rawSources := args[1:]
	var sources []string

	for _, s := range rawSources {
		if strings.HasSuffix(strings.ToLower(s), ".txt") {
			lines, errTxt := readLinesFromFile(s)
			if errTxt == nil && len(lines) > 0 {
				sources = append(sources, lines...)
				continue
			}
		}

		if strings.Contains(s, ",") {
			parts := strings.Split(s, ",")
			for _, p := range parts {
				clean := strings.TrimSpace(p)
				if clean != "" {
					sources = append(sources, clean)
				}
			}
			continue
		}

		clean := strings.TrimSpace(s)
		if clean != "" {
			sources = append(sources, clean)
		}
	}

	return dest, sources, nil
}

func readLinesFromFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, l := range strings.Split(string(data), "\n") {
		clean := strings.TrimSpace(l)
		if clean != "" && !strings.HasPrefix(clean, "#") {
			lines = append(lines, clean)
		}
	}
	return lines, nil
}

func prepareAndInspectSources(inputs []string, stagingBase string) ([]SourceRepoConfig, error) {
	var sources []SourceRepoConfig

	for idx, input := range inputs {
		slug := deriveSourceSlug(input, idx)
		stagingDir := filepath.Join(stagingBase, slug)

		if info, errStat := os.Stat(input); errStat == nil && info.IsDir() {
			stagingDir = input
		} else {
			if _, statStaging := os.Stat(stagingDir); os.IsNotExist(statStaging) {
				cmd := exec.Command("git", "clone", "--depth", "50", input, stagingDir)
				if out, errClone := cmd.CombinedOutput(); errClone != nil {
					return nil, fmt.Errorf("failed to clone source %s: %s (%w)", input, string(out), errClone)
				}
			}
		}

		cfg := inspectSingleSource(input, stagingDir)
		sources = append(sources, cfg)
	}

	return sources, nil
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
			destExists := fileExists(destPath)

			if seqIdx == 0 {
				// Base repository files are placed directly
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), destPath); errCopy != nil {
					return nil, 0, errCopy
				}
				uniquePlaced++
				continue
			}

			if !destExists {
				// Non-colliding file placed directly
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), destPath); errCopy != nil {
					return nil, 0, errCopy
				}
				uniquePlaced++
				continue
			}

			// Collision detected
			dirName := filepath.Dir(relPath)
			baseName := filepath.Base(relPath)

			col, ok := collisionMap[relPath]
			if !ok {
				// First collision on this file: sequence the existing base version
				origSeqPath := filepath.Join(destDir, dirName, fmt.Sprintf("01_%s", baseName))
				_ = os.Rename(destPath, origSeqPath)
				if uniquePlaced > 0 {
					uniquePlaced--
				}

				baseHeadSha, _ := gitExec(sources[0].LocalStaging, "rev-parse", "HEAD")
				if len(baseHeadSha) > 8 {
					baseHeadSha = baseHeadSha[:8]
				}

				col = &FileCollisionRecord{
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
			}

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

func assembleManifest(dest *cmdresolver.ResolvedRepoFeature, sources []SourceRepoConfig, collisions []FileCollisionRecord, uniqueCount int) MergeAIManifest {
	var m MergeAIManifest
	m.Attributes.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	m.Attributes.GitMapVersion = "6.520.0"
	m.Attributes.Tool = "gitmap-merge-ai"
	m.Attributes.CollisionStrategy = "prefixed-sequence"

	m.Data.Destination.TargetType = string(dest.TargetType)
	m.Data.Destination.Slug = dest.CanonicalSlug
	m.Data.Destination.ResolvedPath = dest.LocalPath
	m.Data.Destination.OriginURL = dest.RemoteURL

	for i, s := range sources {
		var entry SourceRepoSequenceEntry
		entry.SequenceOrder = i + 1
		entry.RepoURL = s.RepoURL
		entry.Branch = s.Branch
		entry.CommitRange = s.CommitRange
		entry.LastCommitDate = s.LastCommitDate.Format(time.RFC3339)
		entry.Description = s.Description
		entry.ReleaseInfo.LastTag = s.ReleaseTag
		entry.FolderTree = s.FolderTree
		m.Data.RepoSequence = append(m.Data.RepoSequence, entry)
	}

	m.Data.FileCollisions = collisions
	m.Data.UniqueFilesDirectlyPlaced = uniqueCount
	return m
}

func renderMergeAISummary(dest *cmdresolver.ResolvedRepoFeature, sources []SourceRepoConfig, collisions []FileCollisionRecord, uniqueCount int) {
	fmt.Println()
	pterm.Success.Printf("[GitMap Merge-AI] Multi-repo amalgamation staged successfully!\n")
	fmt.Printf("Destination Workspace: %s\n", pterm.Bold.Sprint(dest.LocalPath))
	fmt.Printf("Source Repositories Merged: %d\n", len(sources))
	fmt.Printf("Unique Files Directly Placed: %d\n", uniqueCount)
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
