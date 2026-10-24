// Package cmdmergeai implements root manifest and AI instruction generation for merged repositories.
package cmdmergeai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// WriteMergeAIManifest serializes and saves merge-ai-manifest.json into destination repository root.
func WriteMergeAIManifest(destDir string, manifest MergeAIManifest) error {
	filePath := filepath.Join(destDir, "merge-ai-manifest.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "failed to marshal merge-ai-manifest.json")
	}

	return os.WriteFile(filePath, data, 0644)
}

// WriteMergeAIInstruction writes instruction.md into destination repository root.
func WriteMergeAIInstruction(destDir string, collisions []FileCollisionRecord, skipped ...SkippedSourceItem) error {
	filePath := filepath.Join(destDir, "instruction.md")

	var sb strings.Builder
	sb.WriteString(buildInstructionHeader())
	sb.WriteString(buildCoreRulesSection())
	sb.WriteString(buildFolderStructureAdvice())

	if len(skipped) > 0 {
		sb.WriteString(buildSkippedSection(skipped))
	}

	sb.WriteString(buildCollisionsSection(collisions))
	sb.WriteString(buildActionChecklist())

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

func buildInstructionHeader() string {
	return `# AI Merge Consolidation Instruction & Verification Guide

This repository has been prepared by ` + "`gitmap merge-ai`" + `. Multiple source repositories have been assembled into this staging tree with automated collision sequencing.

`
}

func buildCoreRulesSection() string {
	return `## Non-Negotiable Core Consolidation Rules:
1. **Preserve All Business Logic:** Do NOT delete features, endpoints, components, or data models. Consolidate overlapping files (` + "`01_*`" + `, ` + "`02_*`" + `, ` + "`03_*`" + `) into the single canonical file name (e.g. synthesize ` + "`01_App.tsx`" + ` and ` + "`02_App.tsx`" + ` into ` + "`App.tsx`" + `).
2. **Asset & Media Preservation Invariant:** Never overwrite or delete distinct logos, icons, animations (` + "`.gif`" + `, ` + "`.svg`" + `, ` + "`.png`" + `, ` + "`.ico`" + `), illustrations, or media assets. Preserve distinct assets side-by-side with descriptive semantic names (e.g. ` + "`logo-macromate.png`" + `, ` + "`logo-animated.gif`" + `, ` + "`logo-brand-builder.svg`" + `) or in organized subfolders.
3. **Preserve Side-by-Side & Verify with Humans:** If conflicting assets, logos, brand elements, or themes exhibit stylistic or design differences, preserve both variations side-by-side and clearly annotate them for human stakeholder review rather than arbitrarily choosing one.
4. **Flatten Sequence Prefixes:** Once consolidated into canonical target files or side-by-side asset variants, remove temporary sequence prefixed files (` + "`01_*`" + `, ` + "`02_*`" + `) so the repository strictly contains clean, canonical filenames.
5. **Reference Commit Metadata in Manifest:** Open ` + "`merge-ai-manifest.json`" + ` and inspect ` + "`repoSequence[].commitMetadata`" + ` (SHA, author, date, message, branch) to understand the origin context, recency, and author of each variant before making consolidation decisions.
6. **Zero Code Deletion:** When resolving duplicate utilities or components, favor additive composition, interface unions, and modern patterns.

`
}

func buildFolderStructureAdvice() string {
	return `## Folder Structure & Architectural Consolidation Advice:

### 1. Root Configuration & Manifests
- **` + "`package.json` / `go.mod` / Build Manifests`" + `:**
  - Create an additive union of all dependencies across source repositories. Where package versions differ, adopt the compatible newer version without dropping any required libraries.
  - Merge custom build, test, and dev scripts so all tools (vite, tailwind, linters) remain functional.
- **` + "`README.md`" + `:**
  - Synthesize a comprehensive monorepo or unified product overview documenting all incorporated submodules, features, and asset packages.
- **` + "`tsconfig.json` / Linters / `.gitignore`" + `:**
  - Merge path aliases, compiler flags, and ignore entries to ensure clean compilation across all merged source paths.

### 2. Source Code & Components (` + "`src/`, `components/`, `lib/`" + `)
- **Same-File / Same-Folder Collisions (` + "`01_*`, `02_*`" + `):**
  - Synthesize component hierarchies, routes, and state stores additively into the canonical file path.
  - If multiple sources contain distinct page views or feature workflows, integrate them as coherent routes, subpages, or navigation tabs rather than discarding one.
  - Update local imports to reference canonical consolidated modules.

### 3. Media Assets & Icons (` + "`assets/`, `public/`, `icons/`" + `)
- **Distinct Asset Preservation:**
  - Keep all unique icons, logos, GIFs, and SVGs.
  - For collisions on identical asset paths (e.g. ` + "`01_logo.png`" + ` vs ` + "`02_logo.png`" + `):
    - Name them descriptively (e.g. ` + "`logo-<source>.png`" + `) or place them in source-specific subdirectories under ` + "`assets/branding/`" + ` or ` + "`public/icons/`" + `.
    - Provide a unified asset registry or catalog for downstream consumption.
  - When in doubt about brand or design preference, preserve all variants side-by-side and request human verification.

`
}

func buildCollisionsSection(collisions []FileCollisionRecord) string {
	var sb strings.Builder
	sb.WriteString("## Detected File Collisions & Consolidation Guidance:\n\n")

	if len(collisions) == 0 {
		sb.WriteString("No conflicting file collisions detected. All files were placed directly into their canonical paths.\n\n")
		return sb.String()
	}

	for _, col := range collisions {
		sb.WriteString(formatSingleCollision(col))
	}
	sb.WriteString("\n")

	return sb.String()
}

func formatSingleCollision(col FileCollisionRecord) string {
	category, advice := classifyCollision(col.CanonicalPath)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### `%s`\n", col.CanonicalPath))
	sb.WriteString(fmt.Sprintf("- **Type:** %s\n", category))
	sb.WriteString(fmt.Sprintf("- **Consolidation Advice:** %s\n", advice))
	sb.WriteString("- **Sequenced Variants:**\n")
	for _, v := range col.Variants {
		sb.WriteString(fmt.Sprintf("  - Variant: `%s` (Source: `%s`, Commit: `%s`)\n", v.File, v.SourceRepo, v.Commit))
	}
	sb.WriteString("\n")

	return sb.String()
}

func classifyCollision(path string) (string, string) {
	ext := strings.ToLower(filepath.Ext(path))
	base := strings.ToLower(filepath.Base(path))

	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp", ".bmp", ".tiff", ".avif", ".mp4", ".webm", ".mov", ".wav", ".mp3", ".ttf", ".woff", ".woff2", ".eot":
		return "Media Asset", "Preserve distinct assets side-by-side using descriptive semantic filenames (or categorized subfolders). Do NOT overwrite or delete; verify design preference with human stakeholders."
	}

	if base == "package.json" || base == "package-lock.json" || base == "pnpm-lock.yaml" || base == "yarn.lock" ||
		base == "tsconfig.json" || base == constants.GoModFile || base == "go.sum" || base == "cargo.toml" ||
		base == constants.GitIgnoreFile || base == "vite.config.ts" || base == "vite.config.js" ||
		base == "tailwind.config.ts" || base == "tailwind.config.js" || strings.HasPrefix(base, "readme") {
		return "Configuration / Manifest", "Synthesize configuration additively into canonical file. Merge all unique dependencies, build scripts, and declarations."
	}

	return "Source Code", "Synthesize logic, functions, components, and routes additively into canonical file. Preserve all unique features without deleting business logic."
}

func buildSkippedSection(skipped []SkippedSourceItem) string {
	var sb strings.Builder
	sb.WriteString("## Skipped Sources (Pre-Filter Intelligence):\n")
	sb.WriteString("The following source URLs were automatically skipped during pre-filtering:\n")
	for _, s := range skipped {
		sb.WriteString(formatSkippedSourceItem(s))
	}
	sb.WriteString("\n")

	return sb.String()
}

func formatSkippedSourceItem(s SkippedSourceItem) string {
	if s.IsFile {
		return fmt.Sprintf("- **File URL**: `%s` (Repository `%s` is already included; file path: `%s`)\n", s.RawInput, s.Repo, s.FilePath)
	}

	return fmt.Sprintf("- **Duplicate Repository**: `%s` (%s)\n", s.RawInput, s.Reason)
}

func buildActionChecklist() string {
	return `## Step-by-Step AI Action Checklist:
- [ ] Read ` + "`merge-ai-manifest.json`" + ` and inspect ` + "`repoSequence[].commitMetadata`" + ` and ` + "`fileCollisions`" + `.
- [ ] For each collision:
  - [ ] Analyze differences between sequenced variants (` + "`01_<file>`" + `, ` + "`02_<file>`" + `, etc.).
  - [ ] For code/config: Synthesize unified implementation into ` + "`<file>`" + `.
  - [ ] For visual assets: Preserve distinct media variants side-by-side with descriptive names; verify ambiguities with humans.
  - [ ] Delete temporary sequenced files (` + "`01_*`" + `, ` + "`02_*`" + `) after verified synthesis.
- [ ] Ensure non-colliding files maintain correct imports and package names.
- [ ] Verify that no sequenced files remain in the repository tree.
- [ ] Update ` + "`merge-ai-manifest.json`" + ` setting ` + "`fileCollisions[].status = \"consolidated\"`" + `.
- [ ] Stage and commit the consolidated repository:
      ` + "`gitmap cpf \"merge-ai - consolidate multi-repo merge into unified codebase\"`" + `
`
}
