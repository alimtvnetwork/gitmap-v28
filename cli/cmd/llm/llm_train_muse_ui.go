// Package llm — llm_train_muse_ui.go implements terminal summary,
// watermark reports, checklist generation, and relative path display formatting.
package llm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func printDigestSummary(path, digestDir, digestName string, days int, files []string, commitCount int, since sinceRead) {
	sequenced := listSequencedDigestFiles(digestDir, digestName)
	fmt.Println("======================================================================")
	fmt.Printf("GITMAP LLM TRAIN — MUSE MEMORY SYNC (last %d days)\n", days)
	fmt.Println("======================================================================")
	fmt.Printf("Digest root: %s\n", relForDisplay(digestDir))
	fmt.Printf("Commits: %d | Files changed: %d | Sequenced digests: %d\n", commitCount, len(files), len(sequenced))
	printSinceLastRead(since)
	fmt.Println()
	printDigestFileList("SEQUENCED DIGEST FILES", sequenced)
	printDigestFileList("CHANGED FILES", files)
	printTraversalChecklist(path, digestDir, sequenced)
}

func printSinceLastRead(since sinceRead) {
	fmt.Println("SINCE LAST READ (tracked by gitmap — no MD file):")
	switch since.state {
	case "baseline":
		fmt.Println("  First run: this sync is the baseline watermark.")
	case "unchanged":
		fmt.Printf("  Vault unchanged since last read at %s — 0 new commits.\n", since.readAt)
	case "fallback":
		fmt.Println("  Checkpoint SHA not in history — showing the --days window instead.")
	case "diff":
		fmt.Printf("  %d new commit(s) since last read at %s:\n", since.count, since.readAt)
		if len(since.files) == 0 {
			fmt.Println("    (file list empty)")
		}
		for _, f := range since.files {
			fmt.Printf("    %s\n", f)
		}
	}
}

func printDigestFileList(title string, files []string) {
	fmt.Println(title + ":")
	if len(files) == 0 {
		fmt.Println("  (none)")
		fmt.Println()
		return
	}
	for _, f := range files {
		fmt.Printf("  %s\n", f)
	}
	fmt.Println()
}

func printTraversalChecklist(path, digestDir string, sequenced []string) {
	fmt.Println("======================================================================")
	fmt.Println("AI TRAVERSAL CHECKLIST (execute in order)")
	fmt.Println("======================================================================")
	fmt.Printf("1. Read the digest root contract first (why: it defines the say-READ protocol and the file order):\n   gitmap cat %s\n", relForDisplay(filepath.Join(digestDir, constants.ReadmeFile)))
	fmt.Println("2. Walk the sequenced digests in order (why: lexical order is the authored narrative):")
	for _, f := range sequenced {
		fmt.Printf("   gitmap cat %s\n", relForDisplay(filepath.Join(path, filepath.FromSlash(f))))
	}
	fmt.Printf("3. Scope searches to the digest tree (why: stay inside the digest, avoid repo-wide noise):\n   gitmap aum search \"<keyword>\" %s --ext .md\n", relForDisplay(digestDir))
	fmt.Printf("4. Inspect chat seeds when a conversation needs full context (why: a seed rebuilds one chat without re-reading everything):\n   gitmap lf %s\n", relForDisplay(filepath.Join(digestDir, "chat-seeds")))
	fmt.Printf("5. Learn GitMap autonomous skills (why: Muse can execute GitMap CLI primitives, clustering, and split-DB ops):\n   Skill Spec: .agents/skills/gitmap/SKILL.md\n   gitmap cat .agents/skills/gitmap/SKILL.md\n")
	fmt.Println()
}

func printDigestMissingGuidance(path, digestName string) {
	fmt.Println("======================================================================")
	fmt.Println("GITMAP LLM TRAIN — MUSE MEMORY SYNC")
	fmt.Println("======================================================================")
	fmt.Printf("[LLM-TRAIN-MUSE] No %s/ digest directory at %s (No Muse/ digest directory).\n", digestName, relForDisplay(path))
	fmt.Println("The 7-day Muse digest is authored by the Muse memory-IO program (another worker).")
	fmt.Println("Nothing to traverse yet — once the digest lands, re-run:")
	fmt.Println("  gitmap llm train muse --days 7")
	fmt.Println()
	fmt.Println("In the meantime, learn GitMap autonomous skills:")
	fmt.Println("  Skill Spec: .agents/skills/gitmap/SKILL.md")
	fmt.Println("  gitmap cat .agents/skills/gitmap/SKILL.md")
	fmt.Println()
}

func relForDisplay(abs string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.ToSlash(abs)
	}
	rel, relErr := filepath.Rel(cwd, abs)
	if relErr != nil {
		return filepath.ToSlash(abs)
	}
	return filepath.ToSlash(rel)
}
