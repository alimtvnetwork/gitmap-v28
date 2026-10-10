// Package cmdmergeai implements root manifest and AI instruction generation for merged repositories.
package cmdmergeai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WriteMergeAIManifest serializes and saves merge-ai-manifest.json into destination repository root.
func WriteMergeAIManifest(destDir string, manifest MergeAIManifest) error {
	filePath := filepath.Join(destDir, "merge-ai-manifest.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal merge-ai-manifest.json: %w", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

// WriteMergeAIInstruction writes instruction.md into destination repository root.
func WriteMergeAIInstruction(destDir string, collisions []FileCollisionRecord) error {
	filePath := filepath.Join(destDir, "instruction.md")

	content := `# AI Merge Consolidation Instruction & Verification Guide

This repository has been prepared by ` + "`gitmap merge-ai`" + `. Multiple source repositories have been assembled into this staging tree.

## Non-Negotiable Instructions for the Consolidating AI Agent:
1. **Preserve All Business Logic:** Do NOT delete features, endpoints, or data models. Consolidate overlapping files (` + "`01_*`" + `, ` + "`02_*`" + `) into the single canonical file name (e.g. synthesize ` + "`01_main.go`" + ` and ` + "`02_main.go`" + ` into ` + "`main.go`" + `).
2. **Flatten Sequence Prefixes:** Once consolidated, remove temporary ` + "`01_`" + ` and ` + "`02_`" + ` files so the repository strictly reflects canonical filenames.
3. **Reference the Manifest:** Open and read ` + "`merge-ai-manifest.json`" + ` to inspect the source repositories, commit histories, and collision map.
4. **Zero Code Deletion:** When resolving duplicate utilities, favor additive composition, interface unions, and modern patterns.

## Detected File Collisions:
`

	if len(collisions) == 0 {
		content += "No conflicting file collisions detected. All files were placed directly.\n\n"
	} else {
		for _, col := range collisions {
			content += fmt.Sprintf("- **%s**:\n", col.CanonicalPath)
			for _, v := range col.Variants {
				content += fmt.Sprintf("  - Variant: `%s` (Source: `%s`, Commit: `%s`)\n", v.File, v.SourceRepo, v.Commit)
			}
		}
		content += "\n"
	}

	content += `## Step-by-Step AI Action Checklist:
- [ ] Read ` + "`merge-ai-manifest.json`" + ` and inventory all entries in ` + "`fileCollisions`" + `.
- [ ] For each collision:
  - [ ] Analyze differences between ` + "`01_<file>`" + ` and ` + "`02_<file>`" + `.
  - [ ] Synthesize unified implementation into ` + "`<file>`" + `.
  - [ ] Delete ` + "`01_<file>`" + ` and ` + "`02_<file>`" + ` after verified synthesis.
- [ ] Ensure non-colliding files maintain correct imports and package names.
- [ ] Verify that no ` + "`01_`" + ` or ` + "`02_`" + ` prefixed files remain in the repository tree.
- [ ] Update ` + "`merge-ai-manifest.json`" + ` setting ` + "`fileCollisions[].status = \"consolidated\"`" + `.
- [ ] Stage and commit the consolidated repository:
      ` + "`gitmap cpf \"feat: consolidate multi-repo merge into unified codebase\"`" + `
`

	return os.WriteFile(filePath, []byte(content), 0644)
}
