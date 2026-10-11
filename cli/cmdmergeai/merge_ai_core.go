// Package cmdmergeai coordinates multi-repository staging, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdresolver"
	"github.com/pterm/pterm"
)

func gitExec(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// RunMergeAI executes multi-repository amalgamation and AI merge staging.
func RunMergeAI(args []string) error {
	destTarget, resolvedSources, skippedSources, dryRun, errParse := parseMergeAIArgs(args)
	if errParse != nil {
		return errParse
	}

	if destTarget == "" && len(resolvedSources) == 0 {
		return nil
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

	if dryRun {
		renderDryRunPreview(resolvedDest, resolvedSources, skippedSources)
		return nil
	}

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

	_ = cleanDestStagingDirectory(resolvedDest.LocalPath)

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

	_, _ = gitExec(resolvedDest.LocalPath, "add", "-A")

	renderMergeAISummary(resolvedDest, sources, collisions, uniqueCount, skippedSources)
	return nil
}

func renderDryRunPreview(dest *cmdresolver.ResolvedRepoFeature, sources []ResolvedSourceItem, skipped []SkippedSourceItem) {
	fmt.Println()
	pterm.DefaultHeader.WithFullWidth().Println("GitMap merge-ai: Dry-Run Simulation Plan")
	fmt.Printf("Destination Workspace: %s (Type: %s)\n", pterm.Bold.Sprint(dest.LocalPath), dest.TargetType)
	fmt.Printf("Source Repositories to Merge (%d):\n", len(sources))
	for i, s := range sources {
		localDir := filepath.Join(resolveWorkBaseDir(), deriveSourceSlug(s.RepoURL, i))
		status := "Will Clone"
		if info, err := os.Stat(filepath.Join(localDir, ".git")); err == nil && info.IsDir() {
			status = "Local Workspace Present"
		}
		fmt.Printf("  [%d] %s (%s) -> %s\n", i+1, pterm.Cyan(s.Title), s.RepoURL, pterm.Green(status))
	}
	if len(skipped) > 0 {
		fmt.Printf("\nSources Automatically Skipped by Manifest Intelligence (%d):\n", len(skipped))
		for _, sk := range skipped {
			fmt.Printf("  - %s (%s)\n", sk.RawInput, pterm.Yellow(sk.Reason))
		}
	}
	fmt.Println("\nPlanned Artifacts in Destination Root:")
	fmt.Println("  • merge-ai-manifest.json (Metadata, commit ranges, collision registry)")
	fmt.Println("  • instruction.md (Non-negotiable consolidation guide and checklist for AI)")
	pterm.Success.Println("\n[Dry-Run Simulation Complete] Run without --dry-run to execute staging.")
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
