package cmdos

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type DevTreeRenderOptions struct {
	HasColorEnabled bool   `json:"hasColorEnabled"`
	HasCompactMode  bool   `json:"hasCompactMode"`
	MaxDepth        int    `json:"maxDepth"`
	HighlightPath   string `json:"highlightPath,omitempty"`
}

type DevCacheTreeNode struct {
	Name        string              `json:"name"`
	FullPath    string              `json:"fullPath"`
	SizeBytes   int64               `json:"sizeBytes"`
	FilesCount  int                 `json:"filesCount"`
	IsDirectory bool                `json:"isDirectory"`
	IsCustom    bool                `json:"isCustom"`
	Children    []*DevCacheTreeNode `json:"children,omitempty"`
}

type CacheTreeNode = DevCacheTreeNode

func RenderDevTreeOutput(paths []DiscoveredCachePath, opts DevTreeRenderOptions) {
	opts.HasColorEnabled = true
	tree := BuildDevCacheTree(paths)
	output := RenderDevCacheTree(tree, opts)
	fmt.Print(output)
}

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

func BuildCacheTree(paths []DiscoveredCachePath) *CacheTreeNode {
	return BuildDevCacheTree(paths)
}

func RenderDevCacheTree(root *DevCacheTreeNode, opts DevTreeRenderOptions) string {
	var sb strings.Builder
	sizeStr := formatAnsiSize(root.SizeBytes)
	if !opts.HasColorEnabled {
		sizeStr = FormatCleanSize(root.SizeBytes)
		sb.WriteString(fmt.Sprintf("\n  🌳 %s (%s, %d files)\n", root.Name, sizeStr, root.FilesCount))
	} else {
		sb.WriteString(fmt.Sprintf("\n  🌳 %s%s%s (%s, %d files)\n",
			constants.ColorWhite, root.Name, constants.ColorReset, sizeStr, root.FilesCount))
	}
	renderSubtreeChildren(&sb, root.Children, "  ", opts)
	return sb.String()
}

func RenderCacheTree(root *CacheTreeNode, opts DevTreeRenderOptions) string {
	return RenderDevCacheTree(root, opts)
}

func renderSubtreeChildren(sb *strings.Builder, children []*DevCacheTreeNode, prefix string, opts DevTreeRenderOptions) {
	for i, child := range children {
		branch, nextPrefix := resolveBranchConnectors(prefix, i == len(children)-1)
		sb.WriteString(formatTreeLine(branch, child, opts) + "\n")
		if len(child.Children) > 0 {
			renderSubtreeChildren(sb, child.Children, nextPrefix, opts)
		}
	}
}

func resolveBranchConnectors(prefix string, isLast bool) (string, string) {
	if isLast {
		return prefix + "└── ", prefix + "    "
	}
	return prefix + "├── ", prefix + "│   "
}

func resolveNodeNameColor(isCustom bool) string {
	if isCustom {
		return constants.ColorOrange
	}
	return constants.ColorCyan
}

func formatTreeLine(prefix string, node *DevCacheTreeNode, opts DevTreeRenderOptions) string {
	icon := "📁"
	if !node.IsDirectory {
		icon = "📄"
	}
	if !opts.HasColorEnabled {
		return fmt.Sprintf("%s%s %s [%s - %d files]",
			prefix, icon, node.Name, FormatCleanSize(node.SizeBytes), node.FilesCount)
	}
	nameColor := resolveNodeNameColor(node.IsCustom)
	return fmt.Sprintf("%s%s %s%s%s [%s - %d files]",
		prefix, icon, nameColor, node.Name, constants.ColorReset, formatAnsiSize(node.SizeBytes), node.FilesCount)
}

func attachPathToTree(root *DevCacheTreeNode, path DiscoveredCachePath) {
	ecoName := resolveEcosystemDisplayName(getPathEcosystemKey(path))
	ecoNode := findOrCreateChild(root, ecoName, true)
	ecoNode.IsCustom = path.IsCustom
	curr := ecoNode
	for _, part := range splitPathComponents(path.Path) {
		curr = findOrCreateChild(curr, part, true)
	}
	curr.SizeBytes += path.SizeBytes
	curr.FilesCount += path.FilesCount
}

func splitPathComponents(p string) []string {
	clean := filepath.Clean(p)
	vol := filepath.VolumeName(clean)
	rest := strings.TrimPrefix(clean, vol)
	rest = strings.TrimPrefix(rest, string(filepath.Separator))
	parts := strings.Split(rest, string(filepath.Separator))
	if len(vol) > 0 {
		return append([]string{vol + `\`}, parts...)
	}
	return parts
}

func findOrCreateChild(parent *DevCacheTreeNode, name string, isDir bool) *DevCacheTreeNode {
	for _, child := range parent.Children {
		if strings.EqualFold(child.Name, name) {
			return child
		}
	}
	child := &DevCacheTreeNode{
		Name:        name,
		IsDirectory: isDir,
	}
	parent.Children = append(parent.Children, child)
	return child
}

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
