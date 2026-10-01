# Subtask 03: ANSI Colored Rich Output & Interactive Tree View

> **Parent Plan:** [62-devtools-cache-discovery-tree-and-split-db](../../pending/62-devtools-cache-discovery-tree-and-split-db.md)  
> **Tracking Spec:** [02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md](../../../../02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md)  
> **Primary File Targets:** `cli/cmdos/os_dev_clean_render.go`, `cli/cmdos/os_dev_clean_tree.go`  
> **Status:** `PENDING`  
> **Owner:** Worker 01  

---

## 1. Objective & User Pain Points

### 1.1 Problem Statement
The user reported:
> *"Now, coming to the point that the output is very terrible. As you can see, no coloring, nothing. Okay. And then you didn't do it right... And these paths needs to be listed like which ones you are removing. You did not list out, and you could list out as a tree view as well."*

1. **Uncolored & Flat Output:** Plain monochrome text table with no visual hierarchy, making it hard to distinguish dry-run previews from actual deletions or cached paths from fresh scans.
2. **Opaque Removals:** No itemized path list showing which directories on disk are being cleaned.
3. **Missing Hierarchical Tree View:** No directory tree visualization to show parent-child relationships, sub-caches (`cache`, `testcache`, `pkg/mod`), and individual file counts.

### 1.2 Core Mission
1. Build an ANSI-colored rich terminal rendering engine with colored cards, status badges, and itemized path displays.
2. Implement a dedicated hierarchical tree renderer triggered by `--tree` (`-t`) displaying expandable directory structures with branch connectors (`├──`, `└──`, `│   `), file counts, and color-coded byte sizes.

---

## 2. Terminal UI Mockups

### 2.1 Rich Category Cards with Status Badges & Itemized Paths

```text
  ┌── 📦 Go (Golang) [⚡ Reclaimable] ──────────────────────────────────────────┐
  │  Total: 1.46 GB across 17,302 files in 2 directories                        │
  │                                                                             │
  │  ├── 📁 C:\dev-tool\go\cache                                                │
  │  │   ├── Size: 752.68 MB  │ Files: 4,812  │ Source: [Heuristic / Split-DB]   │
  │  └── 📁 C:\dev-tool\go\pkg\mod                                              │
  │      ├── Size: 711.71 MB  │ Files: 12,490 │ Source: [CLI probe: GOMODCACHE] │
  └─────────────────────────────────────────────────────────────────────────────┘

  ┌── 📦 pnpm Store [⚡ Reclaimable] ───────────────────────────────────────────┐
  │  Total: 840.50 MB across 6,104 files in 1 directory                         │
  │                                                                             │
  │  └── 📁 C:\dev-tool\pnpm\store\v10                                          │
  │      ├── Size: 840.50 MB  │ Files: 6,104  │ Source: [CLI probe: store path] │
  └─────────────────────────────────────────────────────────────────────────────┘

  ✔ Dry-Run Reclaimable: 2.30 GB across 23,406 files and 3 directories (42ms)
```

### 2.2 Hierarchical Directory Tree View (`--tree` / `-t`)

```text
  🌳 Developer Tools Cache Tree (Total: 2.30 GB, 23,406 files)
  ├── 📦 Go (Golang) [1.46 GB - 17,302 files]
  │   └── 📁 C:\dev-tool\go [1.46 GB]
  │       ├── 📁 cache [752.68 MB - 4,812 files]
  │       └── 📁 pkg [711.71 MB - 12,490 files]
  │           └── 📁 mod [711.71 MB - 12,490 files]
  └── 📦 pnpm [840.50 MB - 6,104 files]
      └── 📁 C:\dev-tool\pnpm [840.50 MB]
          └── 📁 store [840.50 MB]
              └── 📁 v10 [840.50 MB - 6,104 files]
```

---

## 3. Data Structures & Types

All data structures use affirmative booleans (`Is...`, `Has...`, `Can...`) complying with repository coding guidelines.

```go
package cmdos

// DevTreeRenderOptions controls tree formatting and colorization.
type DevTreeRenderOptions struct {
	HasColorEnabled bool   `json:"hasColorEnabled"`
	HasCompactMode  bool   `json:"hasCompactMode"`
	MaxDepth        int    `json:"maxDepth"`
	HighlightPath   string `json:"highlightPath,omitempty"`
}

// DevCacheTreeNode represents a node in the hierarchical cache tree.
type DevCacheTreeNode struct {
	Name        string              `json:"name"`
	FullPath    string              `json:"fullPath"`
	SizeBytes   int64               `json:"sizeBytes"`
	FilesCount  int                 `json:"filesCount"`
	IsDirectory bool                `json:"isDirectory"`
	IsCustom    bool                `json:"isCustom"`
	Children    []*DevCacheTreeNode `json:"children,omitempty"`
}

// DevRenderCard models a colored visual card for a single ecosystem.
type DevRenderCard struct {
	Category   string
	Label      string
	Badge      string
	BadgeColor string
	Paths      []DiscoveredCachePath
	TotalBytes int64
	TotalFiles int
}

// DevTreeSummary encapsulates the root tree node and aggregated stats.
type DevTreeSummary struct {
	RootNode       *DevCacheTreeNode
	TotalSizeBytes int64
	TotalFiles     int
	TotalEcosystems int
}
```

---

## 4. Function Signatures & Modular Design (<= 15 Lines per Function)

### 4.1 Rich ANSI Output & Category Cards (`cli/cmdos/os_dev_clean_render.go`)

```go
// RenderDevCleanAnsiCards displays colored cards and itemized path details.
func RenderDevCleanAnsiCards(paths []DiscoveredCachePath, isDryRun bool) {
	grouped := groupPathsByEcosystem(paths)
	for eco, ecoPaths := range grouped {
		card := buildDevRenderCard(eco, ecoPaths, isDryRun)
		renderSingleDevCard(card)
	}
	renderDevCleanTotals(paths, isDryRun)
}

// buildDevRenderCard constructs the card metadata and aggregated metrics.
func buildDevRenderCard(eco string, paths []DiscoveredCachePath, isDryRun bool) DevRenderCard {
	var totalBytes int64
	var totalFiles int
	for _, p := range paths {
		totalBytes += p.SizeBytes
		totalFiles += p.FilesCount
	}
	return DevRenderCard{
		Category:   eco,
		Label:      resolveEcosystemDisplayName(eco),
		Badge:      resolveBadgeText(isDryRun),
		BadgeColor: resolveBadgeColor(isDryRun),
		Paths:      paths,
		TotalBytes: totalBytes,
		TotalFiles: totalFiles,
	}
}

// renderSingleDevCard prints the framed ANSI card with path details.
func renderSingleDevCard(c DevRenderCard) {
	border := resolveBorderColor()
	reset := constants.ColorReset
	fmt.Printf("\n  %s┌── 📦 %s %s%s%s %s%s\n", border, c.Label, c.BadgeColor, c.Badge, reset, border, strings.Repeat("─", 36))
	fmt.Printf("  %s│%s  Total: %s across %d files in %d paths\n", border, reset, formatAnsiSize(c.TotalBytes), c.TotalFiles, len(c.Paths))
	fmt.Printf("  %s│%s\n", border, reset)
	for i, p := range c.Paths {
		isLast := i == len(c.Paths)-1
		renderCardPathRow(p, isLast, border)
	}
	fmt.Printf("  %s└──%s%s\n", border, strings.Repeat("─", 68), reset)
}

// renderCardPathRow prints an individual itemized path within a card.
func renderCardPathRow(p DiscoveredCachePath, isLast bool, border string) {
	connector := "├──"
	if isLast {
		connector = "└──"
	}
	pathColor := constants.ColorCyan
	sizeColor := constants.ColorYellow
	dimColor := constants.ColorDim
	reset := constants.ColorReset
	fmt.Printf("  %s│%s  %s 📁 %s%s%s\n", border, reset, connector, pathColor, p.Path, reset)
	fmt.Printf("  %s│%s      ├── Size: %s%s%s │ Files: %d │ [%s%s%s]\n",
		border, reset, sizeColor, formatAnsiSize(p.SizeBytes), reset, p.FilesCount, dimColor, p.DiscoverySource, reset)
}

// resolveBadgeText returns the status label depending on execution mode.
func resolveBadgeText(isDryRun bool) string {
	if isDryRun {
		return "[⚡ Reclaimable]"
	}
	return "[✔ Cleaned]"
}

// resolveBadgeColor returns bright yellow for dry-run and bright green for clean.
func resolveBadgeColor(isDryRun bool) string {
	if isDryRun {
		return constants.ColorYellow
	}
	return constants.ColorGreen
}

// formatAnsiSize formats byte counts with vivid ANSI highlighting.
func formatAnsiSize(bytes int64) string {
	const gb = 1024 * 1024 * 1024
	const mb = 1024 * 1024
	color := constants.ColorCyan
	if bytes >= gb {
		color = constants.ColorGreen
		return fmt.Sprintf("%s%.2f GB%s", color, float64(bytes)/float64(gb), constants.ColorReset)
	}
	if bytes >= 50*mb {
		color = constants.ColorYellow
	}
	return fmt.Sprintf("%s%.2f MB%s", color, float64(bytes)/float64(mb), constants.ColorReset)
}
```

### 4.2 Interactive Directory Tree View (`cli/cmdos/os_dev_clean_tree.go`)

```go
// RenderDevTreeOutput builds and prints the hierarchical tree to terminal.
func RenderDevTreeOutput(paths []DiscoveredCachePath, opts DevTreeRenderOptions) {
	tree := BuildDevCacheTree(paths)
	output := RenderDevCacheTree(tree, opts)
	fmt.Println(output)
}

// BuildDevCacheTree constructs the hierarchical tree nodes from discovered paths.
func BuildDevCacheTree(paths []DiscoveredCachePath) *DevCacheTreeNode {
	root := &DevCacheTreeNode{
		Name:        "Developer Tools Cache Tree",
		IsDirectory: true,
	}
	for _, p := range paths {
		attachPathToTree(root, p)
	}
	calculateTreeAggregates(root)
	return root
}

// RenderDevCacheTree formats the tree hierarchy with UTF-8 branch connectors.
func RenderDevCacheTree(root *DevCacheTreeNode, opts DevTreeRenderOptions) string {
	var sb strings.Builder
	titleColor := constants.ColorWhite
	reset := constants.ColorReset
	sb.WriteString(fmt.Sprintf("\n  🌳 %s%s%s (%s, %d files)\n",
		titleColor, root.Name, reset, formatAnsiSize(root.SizeBytes), root.FilesCount))
	renderSubtreeChildren(&sb, root.Children, "  ", opts)
	return sb.String()
}

// renderSubtreeChildren recursively appends formatted branch lines.
func renderSubtreeChildren(sb *strings.Builder, children []*DevCacheTreeNode, prefix string, opts DevTreeRenderOptions) {
	for i, child := range children {
		isLast := i == len(children)-1
		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}
		sb.WriteString(formatTreeLine(prefix+branch, child, opts) + "\n")
		if len(child.Children) > 0 {
			renderSubtreeChildren(sb, child.Children, nextPrefix, opts)
		}
	}
}

// formatTreeLine composes a single node line with icon, name, and metrics.
func formatTreeLine(prefix string, node *DevCacheTreeNode, opts DevTreeRenderOptions) string {
	icon := "📁"
	if !node.IsDirectory {
		icon = "📄"
	}
	nameColor := constants.ColorCyan
	if node.IsCustom {
		nameColor = constants.ColorOrange
	}
	reset := constants.ColorReset
	sizeStr := formatAnsiSize(node.SizeBytes)
	return fmt.Sprintf("%s%s %s%s%s [%s - %d files]",
		prefix, icon, nameColor, node.Name, reset, sizeStr, node.FilesCount)
}

// attachPathToTree inserts a discovered cache path into the hierarchy.
func attachPathToTree(root *DevCacheTreeNode, path DiscoveredCachePath) {
	ecoNode := findOrCreateChild(root, path.Ecosystem, true)
	ecoNode.IsCustom = path.IsCustom
	pathParts := splitPathComponents(path.Path)
	curr := ecoNode
	for _, part := range pathParts {
		curr = findOrCreateChild(curr, part, true)
	}
	curr.SizeBytes += path.SizeBytes
	curr.FilesCount += path.FilesCount
}

// calculateTreeAggregates bubbles byte and file counts up the tree branches.
func calculateTreeAggregates(node *DevCacheTreeNode) (int64, int) {
	if len(node.Children) == 0 {
		return node.SizeBytes, node.FilesCount
	}
	var totalBytes int64
	var totalFiles int
	for _, child := range node.Children {
		b, f := calculateTreeAggregates(child)
		totalBytes += b
		totalFiles += f
	}
	node.SizeBytes = totalBytes
	node.FilesCount = totalFiles
	return totalBytes, totalFiles
}
```

---

## 5. Status Badges & Palette Guide

| Badge / Element | Color Constant | Hex / ANSI Sequence | Purpose |
| :--- | :--- | :--- | :--- |
| `[✔ Cleaned]` | `constants.ColorGreen` | `\033[1;92m` | Successful deletion in active clean mode |
| `[⚡ Reclaimable]` | `constants.ColorYellow` | `\033[1;93m` | Non-destructive dry-run preview |
| `[💾 Cached]` | `constants.ColorCyan` | `\033[1;96m` | Loaded from Split-DB without rescan |
| `[🔄 Fresh Scan]` | `constants.ColorMagenta` | `\033[1;95m` | Forced deep discovery triggered via `--force` |
| `[📂 Custom Path]` | `constants.ColorOrange` | `\033[38;5;208m` | Non-standard drive or user directory detected |
| Directory Nodes | `constants.ColorCyan` | `\033[1;96m` | Directory name in tree view |
| Sizes (> 500 MB) | `constants.ColorGreen` | `\033[1;92m` | High reclaimed storage highlight |

---

## 6. Verification & Acceptance Criteria

- **AC-1:** Executing `gitmap clear devtools --dry-run` prints rich colored cards with `[⚡ Reclaimable]` badges instead of flat monochrome text.
- **AC-2:** Every cleaned directory is itemized with its full path, byte size, file count, and discovery origin.
- **AC-3:** Passing `--tree` (`-t`) generates a hierarchical tree structure with UTF-8 branches (`├──`, `└──`, `│   `) and aggregated byte counts per folder node.
- **AC-4:** Large directories (> 500 MB, e.g. `C:\dev-tool\go\cache`) are highlighted in bold green or orange.
- **AC-5:** All functions adhere strictly to `<= 15 lines` and affirmative boolean naming.
